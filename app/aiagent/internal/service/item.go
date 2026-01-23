package service

import (
	"context"
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/golang/protobuf/ptypes/empty"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	trackerpb "store/api/aiagent"
	"store/app/aiagent/configs"
	"store/pkg/sdk/conv"
	"store/pkg/sdk/third/gemini"
	"store/pkg/sdk/third/tikhub"
	"strings"
)

func (t TrackerService) UpdateCategory(ctx context.Context) {

	list, err := t.Data.Mongo.Item.List(ctx, bson.M{"cat": &bson.M{"$exists": false}}, options.Find().SetLimit(1))
	if err != nil {
		return
	}

	for _, x := range list {

		name := map[string]string{
			"财经":   "finance",
			"AI":   "ai",
			"大健康":  "health",
			"职场":   "workplace",
			"情感":   "emotion",
			"科技数码": "emotion",
			"家居":   "home",
			"美食":   "food",
			"母婴育儿": "parenting",
			"医美护肤": "skincare",
			"珠宝":   "jewelry",
			"移民留学": "studyabroad",
			"财富保险": "insurance",
			"房地产":  "realestate",
		}[x.Category]

		f, err := t.Data.Mongo.Item.UpdateFieldsById(ctx, x.XId,
			bson.M{
				"cat": bson.M{"value": name, "name": x.Category},
			})

		fmt.Println(f)

		if err != nil {
			log.Error(err)
			return
		}
	}

}

func (t TrackerService) UpdateItemId(ctx context.Context) {

	list, err := t.Data.Mongo.Item.List(ctx, bson.M{"_id": bson.M{"$not": primitive.Regex{
		Pattern: "^.{24}$",
		Options: "",
	}}}, options.Find().SetLimit(1))
	if err != nil {
		return
	}

	for _, x := range list {

		log.Debugw("update item id", x)

		url := "http://xhslink.com/a/" + x.XId
		note, err := t.Data.Tikhub.XhsGetNoteByShareUrl(ctx, url)
		if err != nil {
			url = "http://xhslink.com/m/" + x.XId
			note, err = t.Data.Tikhub.XhsGetNoteByShareUrl(ctx, url)
		}

		if err != nil {
			continue
		}

		oldId := x.XId

		x.XId = note.NoteId
		_, err = t.Data.Mongo.Item.InsertIfNotExistsByID(ctx, note.NoteId, x)
		if err != nil {
			log.Error(err)
			continue
		}

		err = t.Data.Mongo.Item.DestroyById(ctx, oldId)
		if err != nil {
			return
		}
	}

}

func (t TrackerService) UpdateOssUrl(ctx context.Context) {

	list, err := t.Data.Mongo.Item.List(ctx, bson.M{"a": &bson.M{"$exists": false}}, options.Find().SetLimit(1))
	if err != nil {
		return
	}

	for _, x := range list {
		coverUrl, err := t.uploadToS3(ctx, x.Cover, "image/jpeg")
		if err != nil {
			log.Error(err)
			return
		}

		f, err := t.Data.Mongo.Item.UpdateFieldsById(ctx, x.XId,
			bson.M{
				"cover": coverUrl,
				"a":     "1",
			})

		fmt.Println(f)

		if err != nil {
			log.Error(err)
			return
		}
	}

}

func (t TrackerService) FillItemRecord(ctx context.Context, _ *empty.Empty) (*empty.Empty, error) {
	ctx = context.Background()

	logger := log.NewHelper(log.With(log.DefaultLogger,
		"func", "FillItemRecord",
	))

	list, err := t.Data.Mongo.Item.List(ctx, bson.M{"status": "pending"}, options.Find().SetLimit(1))
	if err != nil {
		return nil, err
	}

	if len(list) == 0 {
		return &empty.Empty{}, nil
	}

	x := list[0]

	logger.Debugw("fill item record", "FillItemRecord", "x", x.XId)

	url, err := t.uploadToS3(ctx, x.Url, "video/mp4")
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	coverUrl, err := t.uploadToS3(ctx, x.Cover, "image/jpeg")
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	genaiClient := t.Data.GenaiFactory.Get()

	parts, err := t.packSystemPromptV2(ctx, genaiClient, "analysis",
		&trackerpb.Prompt{Id: "analysis"},
		trackerpb.Resources{
			{
				Url:      url,
				MimeType: "video/mp4",
			},
			{
				Category: "authorProfile",
				MimeType: "text/plain",
				Content:  conv.S2J(x.Profile),
				Meta: map[string]string{
					"username":      x.Profile.Username,
					"sign":          x.Profile.Sign,
					"tags":          strings.Join(x.Profile.Tags, ","),
					"followerCount": x.Profile.FollowerCount,
					"likedCount":    x.Profile.LikedCount,
					"noteCount":     x.Profile.NoteCount,
					"platform":      "xiaohongshu",
				},
			},
		},
	)

	if err != nil {
		return nil, err
	}

	//model := service.Data.GenaiClient.GenerativeModel("gemini-2.0-flash")
	//model := t.Data.GenaiClient.GenerativeModel("gemini-2.5-pro-exp-03-25")

	for i := 0; i < 6; i++ {
		modelName := configs.ModelFlash
		answer, err := genaiClient.GenerateContent(ctx, gemini.GenerateContentRequest{
			Model: modelName,
			Parts: parts,
		})
		if err != nil {
			log.Error(err)

			if strings.Contains(err.Error(), "blocked") {
				t.Data.Mongo.Item.DestroyById(ctx, x.XId)
				return nil, err
			}

			continue

		}

		if answer != "" {
			_, err = t.Data.Mongo.Item.UpdateFieldsById(ctx, x.XId,
				bson.M{
					"url":              url,
					"cover":            coverUrl,
					"reports.analysis": answer,
					"status":           "active",
				})
			if err != nil {
				log.Error(err)
				return nil, err
			}
			break
		}

	}

	return &empty.Empty{}, nil
	//t.Data.Mongo.
}

func (t TrackerService) ListItems(ctx context.Context, params *trackerpb.ListItemsParams) (*trackerpb.ItemList, error) {

	size := int64(24)

	if params.Keyword != "" {
		return t.listItems(ctx, params)
	}

	filters := bson.M{"status": "active"}
	if params.Category != "" {
		filters["category"] = params.Category
	}

	if params.Cat != "" {
		filters["cat.value"] = params.Cat
	}

	list, count, err := t.Data.Mongo.Item.ListAndCount(ctx, filters,
		options.Find().
			SetSort(bson.M{"_id": 1}).
			SetLimit(size).
			SetSkip((params.Page-1)*size))
	if err != nil {
		log.Error(err)
		return nil, err
	}

	return &trackerpb.ItemList{
		List:  trackerpb.Items(list).Safe(),
		Total: count,
		//HasMore: true,
		HasMore: params.Page*size <= count,
		Page:    params.Page,
	}, nil
}
func (t TrackerService) listItems(ctx context.Context, params *trackerpb.ListItemsParams) (*trackerpb.ItemList, error) {

	notes, err := t.Data.Tikhub.XhsSearchNotes(ctx, tikhub.XhsSearchNotesParams{
		Keyword:  params.Keyword,
		Sort:     "",
		NoteType: "",
	})
	if err != nil {
		return nil, err
	}

	var items []*trackerpb.Item
	for _, x := range notes {
		items = append(items, &trackerpb.Item{
			XId:   x.Id,
			Title: x.Title,
			Profile: &trackerpb.Profile{
				XId:      x.User.Id,
				Avatar:   x.User.Avatar,
				Username: x.User.Name,
			},
			InteractInfo: &trackerpb.Item_InteractInfo{
				LikedCount:     conv.String(x.LikedCount),
				CollectedCount: conv.String(x.CollectedCount),
				CommentCount:   conv.Str(x.CommentsCount),
				SharedCount:    conv.Str(x.SharedCount),
			},
			Cover: x.Cover,
			Desc:  x.Desc,
		})
	}

	return &trackerpb.ItemList{
		List: trackerpb.Items(items).Safe(),
	}, nil
}
