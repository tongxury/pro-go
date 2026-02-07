package biz

import (
	"context"
	"encoding/json"
	"fmt"
	voiceagent "store/api/voiceagent"
	"store/pkg/clients/mgz"
	"store/pkg/sdk/helper"
	"strings"
	"time"

	"store/pkg/sdk/third/gemini"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/livekit/protocol/livekit"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/genai"
)

func (b *AgentBiz) SummarizeEndedConversations(ctx context.Context) error {
	// 1. Find conversations that are ended but not yet analyzed
	filter := mgz.Filter().EQ("status", "completed").NotEQ("extra.analysisStatus", "completed").B()
	// Limit to process a batch at a time to avoid long running loops
	opts := options.Find().SetLimit(10)

	conversations, err := b.data.Mongo.Conversation.Find(ctx, filter, opts)
	if err != nil {
		log.Errorw("msg", "failed to find conversations to summarize", "err", err)
		return err
	}

	if len(conversations) == 0 {
		return nil
	}

	log.Infow("msg", "found conversations to summarize", "count", len(conversations))

	//wg.WaitGroup(ctx, conversations, b.processConversation)
	for _, conversation := range conversations {
		b.processConversation(ctx, conversation)
	}

	return nil
}

func (b *AgentBiz) processConversation(ctx context.Context, conv *voiceagent.Conversation) error {
	// 2. Fetch transcript entries
	transcripts, err := b.data.Mongo.Transcript.Find(ctx, mgz.Filter().EQ("conversation._id", conv.XId).B())
	if err != nil {
		return fmt.Errorf("failed to fetch transcripts: %w", err)
	}

	if len(transcripts) == 0 {
		// No transcripts, mark as completed
		_, _ = b.data.Mongo.Conversation.UpdateByIDIfExists(ctx, conv.XId, mgz.Op().Set("extra.analysisStatus", "completed"))
		return nil
	}

	// 3. Construct prompt
	var transcriptText strings.Builder
	for _, t := range transcripts {
		role := "User"
		if t.Role == "agent" {
			role = "AI"
		}
		transcriptText.WriteString(fmt.Sprintf("%s: %s\n", role, t.Message))
	}

	prompt := fmt.Sprintf(`Extract user facts, preferences, events, and relationship updates from this transcript:
%s`, transcriptText.String())

	// 4. Call Gemini
	client := b.data.Gemini.Get()
	model := "gemini-2.0-flash-exp"

	responseText, err := client.GenerateContent(ctx, gemini.GenerateContentRequest{
		Model: model,
		Parts: []*genai.Part{{Text: prompt}},
		Config: &genai.GenerateContentConfig{
			ResponseMIMEType: "application/json",
			ResponseSchema: &genai.Schema{
				Type: "object",
				Properties: map[string]*genai.Schema{
					"subject": {
						Type:        genai.TypeString,
						Description: "会话主题，精简描述，10个字左右",
					},
					"memories": {
						Type: genai.TypeArray,
						Items: &genai.Schema{
							Type:     genai.TypeObject,
							Required: []string{"type", "content", "importance", "tags"},
							Properties: map[string]*genai.Schema{
								"type": {
									Type:        genai.TypeString,
									Description: "One of 'fact', 'preference', 'event', 'relationship'",
									Enum:        []string{"fact", "preference", "event", "relationship"},
								},
								"content": {
									Type:        genai.TypeString,
									Description: "Concise description of the extracted information",
								},
								"importance": {
									Type:        genai.TypeInteger,
									Minimum:     helper.Pointer[float64](1),
									Maximum:     helper.Pointer[float64](5),
									Description: "Importance level from 1 to 5",
								},
								"tags": {
									Type: genai.TypeArray,
									Items: &genai.Schema{
										Type: genai.TypeString,
									},
								},
							},
						},
					},
				},
			},
		},
	})
	if err != nil {
		return fmt.Errorf("gemini generation failed: %w", err)
	}

	var result struct {
		Memories []*voiceagent.Memory
		Subject  string
	}

	if err := json.Unmarshal([]byte(responseText), &result); err != nil {
		log.Errorw("msg", "failed to unmarshal gemini response", "response", responseText, "err", err)
		return fmt.Errorf("json unmarshal failed: %w", err)
	}

	// 5. Save memories
	for _, x := range result.Memories {
		x.CreatedAt = time.Now().Unix()
		x.UpdatedAt = time.Now().Unix()
		x.Extra = &voiceagent.Memory_Extra{Source: conv.XId}
		x.User = conv.User
	}
	_, err = b.data.Mongo.Memory.InsertMany(ctx, result.Memories...)
	if err != nil {
		log.Errorw("msg", "failed to insert memory", "err", err)
		return err
	}

	// 6. Update conversation status
	_, err = b.data.Mongo.Conversation.UpdateByIDIfExists(ctx, conv.XId,
		mgz.Op().
			Set("extra.analysisStatus", "completed").
			Set("subject", result.Subject),
	)
	if err != nil {
		return fmt.Errorf("failed to update conversation status: %w", err)
	}

	log.Infow("msg", "successfully summarized conversation", "id", conv.XId, "memories_count", len(result.Memories), "subject", result.Subject)

	return nil
}

func (b *AgentBiz) CleanupOrphanedConversations(ctx context.Context) error {
	// 1. Find active conversations created more than 30 minutes ago
	thirtyMinsAgo := time.Now().Unix() - 1800
	filter := mgz.Filter().EQ("status", "ongoing").LT("createdAt", thirtyMinsAgo).B()

	conversations, err := b.data.Mongo.Conversation.Find(ctx, filter, mgz.Find().Limit(10).B())
	if err != nil {
		return err
	}
	if len(conversations) == 0 {
		return nil
	}

	// 2. Fetch active rooms from LiveKit
	var roomNames []string
	for _, c := range conversations {
		if c.RoomName != "" {
			roomNames = append(roomNames, c.RoomName)
		}
	}

	if len(roomNames) == 0 {
		return nil
	}

	resp, err := b.data.RoomClient.ListRooms(ctx, &livekit.ListRoomsRequest{
		Names: roomNames,
	})
	if err != nil {
		log.Errorw("msg", "failed to list rooms from LiveKit", "err", err)
		return err
	}

	activeRooms := make(map[string]bool)
	for _, r := range resp.Rooms {
		activeRooms[r.Name] = true
	}

	// 3. Mark orphaned conversations as completed
	for _, c := range conversations {
		if !activeRooms[c.RoomName] {
			log.Infow("msg", "cleaning up orphaned conversation", "convId", c.XId, "roomName", c.RoomName)
			_, err := b.data.Mongo.Conversation.UpdateByIDIfExists(ctx, c.XId, mgz.Op().Set("status", "completed").Set("endedAt", time.Now().Unix()))
			if err != nil {
				log.Errorw("msg", "failed to update orphaned conversation", "convId", c.XId, "err", err)
			}
		}
	}

	return nil
}
