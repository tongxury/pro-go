package biz

import (
	"context"
	"encoding/json"
	"fmt"
	voiceagent "store/api/voiceagent"
	"store/pkg/clients/mgz"
	"strings"
	"time"

	"store/pkg/sdk/helper"
	"store/pkg/sdk/third/gemini"

	"github.com/go-kratos/kratos/v2/log"
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
		_, _ = b.data.Mongo.Conversation.UpdateOne(ctx, mgz.Filter().EQ("_id", conv.XId).B(), mgz.Op().Set("extra.analysisStatus", "completed"))
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
	})
	if err != nil {
		return fmt.Errorf("gemini generation failed: %w", err)
	}

	var memories []*voiceagent.Memory
	if err := json.Unmarshal([]byte(responseText), &memories); err != nil {
		log.Errorw("msg", "failed to unmarshal gemini response", "response", responseText, "err", err)
		return fmt.Errorf("json unmarshal failed: %w", err)
	}

	// 5. Save memories
	for _, x := range memories {
		x.CreatedAt = time.Now().Unix()
		x.UpdatedAt = time.Now().Unix()
		x.Extra = &voiceagent.Memory_Extra{Source: conv.XId}
		x.User = conv.User
	}
	_, err = b.data.Mongo.Memory.InsertMany(ctx, memories...)
	if err != nil {
		log.Errorw("msg", "failed to insert memory", "err", err)
		return err
	}

	// 6. Update conversation status
	_, err = b.data.Mongo.Conversation.UpdateOne(ctx, mgz.Filter().EQ("_id", conv.XId).B(), mgz.Op().Set("extra.analysisStatus", "completed"))
	if err != nil {
		return fmt.Errorf("failed to update conversation status: %w", err)
	}

	log.Infow("msg", "successfully summarized conversation", "id", conv.XId, "memories_count", len(memories))

	return nil
}
