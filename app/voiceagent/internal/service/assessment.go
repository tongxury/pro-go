package service

import (
	"context"
	ucpb "store/api/usercenter"
	pb "store/api/voiceagent"
	"store/app/voiceagent/internal/data"
	"store/pkg/krathelper"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type AssessmentService struct {
	pb.UnimplementedAssessmentServiceServer

	Data *data.Data
}

func NewAssessmentService(data *data.Data) *AssessmentService {
	return &AssessmentService{
		Data: data,
	}
}

func (s *AssessmentService) CreateAssessment(ctx context.Context, req *pb.CreateAssessmentRequest) (*pb.Assessment, error) {
	userId := krathelper.RequireUserId(ctx)
	// TODO: Validate input

	assessment := &pb.Assessment{
		User:      &ucpb.User{XId: userId},
		Type:      req.Type,
		Score:     req.Score,
		Level:     req.Level,
		CreatedAt: time.Now().UnixMilli(),
		Details:   req.Details,
	}

	if _, err := s.Data.Mongo.Assessment.Insert(ctx, assessment); err != nil {
		return nil, err
	}

	return assessment, nil
}

func (s *AssessmentService) ListAssessments(ctx context.Context, req *pb.ListAssessmentsRequest) (*pb.AssessmentList, error) {
	userId := krathelper.RequireUserId(ctx)

	// Using $or to match potential BSON field names for User ID
	filter := bson.M{
		"$or": []bson.M{
			{"user._id": userId},
		},
	}

	page := req.Page
	if page < 1 {
		page = 1
	}
	size := req.Size
	if size < 1 {
		size = 20
	}

	opts := options.Find().
		SetSort(bson.D{{Key: "createdAt", Value: -1}}).
		SetSkip((page - 1) * size).
		SetLimit(size)

	assessments, total, err := s.Data.Mongo.Assessment.ListAndCount(ctx, filter, opts)
	if err != nil {
		return nil, err
	}

	return &pb.AssessmentList{
		List:  assessments,
		Total: total,
	}, nil
}
