package service

import (
	"context"
	voiceagent "store/api/voiceagent"
	"store/pkg/clients/mgz"

	"go.mongodb.org/mongo-driver/bson"
)

func (s *LiveKitService) ListTranscriptEntries(ctx context.Context, req *voiceagent.ListTranscriptEntriesRequest) (*voiceagent.TranscriptEntryList, error) {
	list, total, err := s.data.Mongo.Transcript.ListAndCount(ctx, bson.M{"conversation._id": req.ConversationId}, mgz.Paging(req.Page, req.Size))
	if err != nil {
		return nil, err
	}

	return &voiceagent.TranscriptEntryList{
		List:  list,
		Total: total,
	}, nil
}
