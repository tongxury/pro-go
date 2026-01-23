package service

import (
	"context"
	voiceagent "store/api/voiceagent"
	"store/pkg/clients/mgz"

	"go.mongodb.org/mongo-driver/bson"
)

func (s *VoiceAgentService) ListTranscriptEntries(ctx context.Context, req *voiceagent.ListTranscriptEntriesRequest) (*voiceagent.TranscriptEntryList, error) {
	list, total, err := s.Data.Mongo.Transcript.ListAndCount(ctx, bson.M{"conversationId": req.ConversationId}, mgz.Paging(req.Page, req.Size))
	if err != nil {
		return nil, err
	}

	return &voiceagent.TranscriptEntryList{
		List:  list,
		Total: total,
	}, nil
}
