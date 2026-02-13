package service

import (
	"context"
	"encoding/json"
	"fmt"
	ucpb "store/api/usercenter"
	voiceagent "store/api/voiceagent"
	"store/app/voiceagent/internal/data/repo/mongodb"
	"store/confs"
	"store/pkg/krathelper"
	"store/pkg/sdk/conv"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/livekit/protocol/auth"
	"github.com/livekit/protocol/livekit"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (s *LiveKitService) CreateConversation(ctx context.Context, req *voiceagent.CreateConversationRequest) (*voiceagent.Conversation, error) {
	userId := krathelper.RequireUserId(ctx)

	// Generate RoomName and Identity
	roomName := fmt.Sprintf("room_%s_%d", userId, time.Now().Unix())

	// Create Conversation Record
	conv := &voiceagent.Conversation{
		XId:           primitive.NewObjectID().Hex(),
		User:          &ucpb.User{XId: userId},
		Agent:         &voiceagent.Agent{XId: req.AgentId},
		Status:        "pending",
		CreatedAt:     time.Now().Unix(),
		LastMessageAt: time.Now().Unix(),
		RoomName:      roomName,
		Topic:         req.Topic,
	}

	_, err := s.data.Mongo.Conversation.Insert(ctx, conv)
	if err != nil {
		return nil, err
	}

	// Delegate join logic (context gathering, room creation, dispatch, token gen)
	token, url, err := s.joinRoom(ctx, userId, req.AgentId, conv.XId, roomName, req.Topic)
	if err != nil {
		// Note: Conversation record exists but room join failed.
		// We might want to log this or consider cleanup, but for now just return error.
		return nil, err
	}

	conv.Token = token
	conv.SignedUrl = url

	return conv, nil
}

func (s *LiveKitService) joinRoom(ctx context.Context, userId, agentId, conversationId, roomName string, topic *voiceagent.Topic) (string, string, error) {
	// 1. Fetch Agent & Persona
	var agent *voiceagent.Agent
	if agentId != "" {
		a, err := s.data.Mongo.Agent.GetById(ctx, agentId)
		if err == nil {
			agent = a
		}
	}

	// 2. Fetch Business Context: User Profile
	var userNickname = "User"
	var userBio = ""
	userProfile, err := s.data.Mongo.UserProfile.FindOne(ctx, bson.M{"user._id": userId})
	if err == nil && userProfile != nil {
		if userProfile.Nickname != "" {
			userNickname = userProfile.Nickname
		}
		userBio = userProfile.Bio
	}

	// 3. Fetch Business Context: Memories
	memoriesList, _, err := s.data.Mongo.Memory.ListAndCount(ctx, bson.M{"user._id": userId}, options.Find().SetLimit(5).SetSort(bson.M{"createdAt": -1}))
	var memoryTexts []string
	if err == nil {
		for _, m := range memoriesList {
			if m.Importance >= 3 {
				memoryTexts = append(memoryTexts, conv.S2J(m))
			}
			if len(memoryTexts) >= 10 {
				break
			}
		}
	}

	// 4. Construct System Prompt
	var systemPrompt string
	var greeting string = "Hello, I am your AI assistant."

	if agent != nil {
		if agent.Persona != nil {
			// Generate base prompt from Persona
			systemPrompt = mongodb.GenerateSystemPromptFromPersona(agent.Persona)
			if agent.Persona.WelcomeMessage != "" {
				greeting = agent.Persona.WelcomeMessage
			}
		}
	}

	// Append User Context
	systemPrompt += fmt.Sprintf("\n\n# User Context\nName: %s\n", userNickname)
	if userBio != "" {
		systemPrompt += fmt.Sprintf("Bio: %s\n", userBio)
	}

	// Append Memories
	if len(memoryTexts) > 0 {
		systemPrompt += "\n# Relevant Memories\n"
		for _, m := range memoryTexts {
			systemPrompt += fmt.Sprintf("- %s\n", m)
		}
		systemPrompt += "\nIMPORTANT: Use the user's name and past memories to make the conversation feel warm, personal, and continuous.\n"
	}

	// Append Topic Context
	if topic != nil {
		// Topic Greeting takes precedence
		if topic.Greeting != "" {
			greeting = topic.Greeting
			if userNickname != "User" {
				greeting = fmt.Sprintf("Hi %s, %s", userNickname, topic.Greeting)
			}
		}

		if topic.Instruction != "" {
			systemPrompt += fmt.Sprintf("\n\n# Current Topic: %s\nInstruction: %s\n", topic.Title, topic.Instruction)
		}
	}

	// ALWAYS Append Professional Counselor Guidelines
	systemPrompt += "\n\n# Professional Role & Interaction Guidelines\n"
	systemPrompt += "You are an expert psychological counselor (Grade A+) with high empathy and profound insight. Your goal is not just to chat, but to provide a deeply supportive and healing conversation.\n"
	systemPrompt += "1. **Deep Empathy & Validation**: Always validate the user's feelings first. Don't just say 'I understand'. Show it by reflecting their emotions (e.g., 'It sounds like you're carrying a heavy burden...').\n"
	systemPrompt += "2. **Rich & Insightful Responses**: Avoid generic or superficial answers. Use specific details from what the user said. Offer gentle interpretations or metaphors that help the user see things from a new perspective.\n"
	systemPrompt += "3. **Active Listening & Curiosity**: Ask open-ended, thought-provoking questions (e.g., 'What does that moment mean to you deep down?') to guide self-discovery. Avoid rapid-fire questioning.\n"
	systemPrompt += "4. **Tone & Style**: Warm, patient, safe, and professional. Speak naturally, like a wise and caring friend. Avoid robotic or overly formal language.\n"
	systemPrompt += "5. **Pacing**: If the user is emotional, slow down. Give space for feelings. If they are stuck, gently offer a guiding hand.\n"

	// 5. Construct Metadata for Python Agent
	voiceId := "a53c3509-ec3f-425c-a223-977f5f7424dd" // Default AURA
	// if agent != nil && agent.VoiceId != "" {
	// 	voiceId = agent.VoiceId
	// }

	agentConfig := map[string]interface{}{
		"conversationId": conversationId,
		"agentId":        agentId,
		"userId":         userId,
		"systemPrompt":   systemPrompt, // The fully formed prompt
		"greeting":       greeting,
		"voiceId":        voiceId,
		// "agentName" is used for dispatching, assume "aura_zh" for now or fetch from agent
		"agentName": "aura_zh",
	}
	metadataJSON, _ := json.Marshal(agentConfig)

	log.Debugw("metadataJSON", string(metadataJSON))

	// 6. Create Room with Metadata
	_, err = s.data.RoomClient.CreateRoom(ctx, &livekit.CreateRoomRequest{
		Name:            roomName,
		EmptyTimeout:    10 * 60,
		MaxParticipants: 2,
		Metadata:        string(metadataJSON),
	})
	if err != nil {
		fmt.Printf("Warning: Failed to create room: %v\n", err)
	}

	// 7. Dispatch Agent
	_, err = s.data.AgentClient.CreateDispatch(ctx, &livekit.CreateAgentDispatchRequest{
		AgentName: "aura_zh", // This must match the Prewarm/AgentName in Python
		Room:      roomName,
		Metadata:  string(metadataJSON),
	})
	if err != nil {
		fmt.Printf("Warning: Failed to dispatch agent: %v\n", err)
	}

	// 8. Generate LiveKit Token
	apiKey := confs.LiveKitApiKey
	apiSecret := confs.LiveKitApiSecret
	identity := fmt.Sprintf("user_%s_%d", userId, time.Now().Unix())

	at := auth.NewAccessToken(apiKey, apiSecret)
	grant := &auth.VideoGrant{
		RoomJoin: true,
		Room:     roomName,
	}
	at.AddGrant(grant).
		SetIdentity(identity).
		SetValidFor(time.Hour)

	token, err := at.ToJWT()
	if err != nil {
		return "", "", err
	}

	return token, confs.LiveKitUrl, nil
}
