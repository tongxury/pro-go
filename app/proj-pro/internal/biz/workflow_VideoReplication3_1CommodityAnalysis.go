package biz

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sort"
	projpb "store/api/proj"
	"store/app/proj-pro/internal/data"
	"store/pkg/clients/mgz"
	"store/pkg/sdk/helper"
	"store/pkg/sdk/third/bytedance/vikingdb"
	"store/pkg/sdk/third/gemini"
	"strings"

	"github.com/go-kratos/kratos/v2/log"
	"go.mongodb.org/mongo-driver/bson"
	"google.golang.org/genai"
)

type VideoReplication3_CommodityAnalysisJob struct {
	data *data.Data
}

func (t VideoReplication3_CommodityAnalysisJob) Initialize(ctx context.Context, options Options) error {
	return nil
}

func (t VideoReplication3_CommodityAnalysisJob) GetName() string {
	return "commodityAnalysisJob"
}

func (t VideoReplication3_CommodityAnalysisJob) Execute(ctx context.Context, jobState *projpb.Job, wfState *projpb.Workflow) (status *ExecuteResult, err error) {

	dataBus := GetDataBus(wfState)

	logger := log.NewHelper(log.With(log.DefaultLogger,
		"func", "VideoReplication3_CommodityAnalysisJob.Execute",
		"workflowId ", wfState.XId,
		"jobState.Name ", jobState.Name,
		"jobState.Index ", jobState.Index,
	))

	commodity := dataBus.Commodity
	if commodity == nil {
		return nil, errors.New("commodity not found")
	}

	if commodity.Name != "" && len(commodity.Chances) > 0 {
		return &ExecuteResult{
			Status: ExecuteStatusCompleted,
		}, nil
	}

	logger.Debugw("start analyze commodity", commodity.XId)

	analyzeResult, err := t.doAnalyzeChances(ctx, commodity)
	if err != nil {
		logger.Errorw("doAnalyzeChances err", err)
		return nil, err
	}

	// 更新商品信息
	commodity.Name = analyzeResult.Name
	commodity.Tags = analyzeResult.Tags
	commodity.Brand = analyzeResult.Brand
	commodity.Description = analyzeResult.Description
	commodity.Character = analyzeResult.Character
	commodity.PackagingAttributes = analyzeResult.PackagingAttributes
	commodity.Chances = analyzeResult.Chances

	for _, x := range commodity.Images {
		commodity.Medias = append(commodity.Medias, &projpb.Media{
			MimeType: "image/jpeg",
			Url:      x,
		})
	}

	// _, err = t.data.Mongo.Commodity.UpdateByIDIfExists(ctx, commodity.XId,
	// 	mgz.Op().
	// 		Set("name", analyzeResult.Name).
	// 		Set("tags", analyzeResult.Tags).
	// 		Set("brand", analyzeResult.Brand).
	// 		Set("description", analyzeResult.Description).
	// 		Set("chances", analyzeResult.Chances).
	// 		Set("character", analyzeResult.Character).
	// 		Set("packagingAttributes", analyzeResult.PackagingAttributes).
	// 		Set("status", "completed"),
	// )
	// if err != nil {
	// 	logger.Errorw("update commodity err", err)
	// 	return nil, err
	// }

	// 更新 workflow 中的 dataBus
	_, err = t.data.Mongo.Workflow.UpdateByIDIfExists(ctx, wfState.XId, mgz.Op().
		Set(fmt.Sprintf("jobs.%d.dataBus.commdity", jobState.Index), commodity))
	if err != nil {
		logger.Errorw("update workflow commodity fail", "err", err)
		return nil, err
	}

	segments, err := t.searchTemplateSegments(ctx, strings.Join(analyzeResult.Tags, ","), "", 1)
	if err != nil {
		return nil, err
	}

	segments = helper.Filter(segments, func(param *projpb.ResourceSegment) bool {

		if int64(param.TimeEnd-param.TimeStart) < 8 {
			return false
		}

		return true
	})

	if len(segments) == 0 {
		logger.Errorw("searchTemplateSegments err", "no segment found", "tags", strings.Join(analyzeResult.Tags, ","))

		count, err := t.data.Redis.Incr(ctx, "video_replication_3_1_no_segment_found:retry:"+wfState.XId).Result()
		if err != nil {
			logger.Errorw("update workflow segment fail", "err", err)
			return nil, err
		}

		if count > 5 {
			return &ExecuteResult{
				Status: ExecuteStatusFailed,
				Error:  "no segment found",
			}, nil
		}

		return nil, nil
	}

	// 更新 workflow 中的 dataBus
	_, err = t.data.Mongo.Workflow.UpdateByIDIfExists(ctx, wfState.XId, mgz.Op().
		Set(fmt.Sprintf("jobs.%d.dataBus.segment", jobState.Index), segments[0]))
	if err != nil {
		logger.Errorw("update workflow segment fail", "err", err)
		return nil, err
	}

	return &ExecuteResult{
		Status: ExecuteStatusCompleted,
	}, nil
}

type AnalyzeResult struct {
	Name                string
	Tags                []string
	PackagingAttributes string
	Character           string
	Brand               string
	Description         string
	Chances             []*projpb.Chance
}

func (t VideoReplication3_CommodityAnalysisJob) doAnalyzeChances(ctx context.Context, commodity *projpb.Commodity) (*AnalyzeResult, error) {

	genaiClient := t.data.GenaiFactory.Get()

	var parts []*genai.Part

	images := helper.Filter(commodity.Images, func(param string) bool {
		return !strings.HasSuffix(param, ".webp")
	})

	for i, x := range images {

		if i > 10 {
			break
		}

		part, err := gemini.NewImagePart(x)
		if err != nil {
			return nil, err
		}
		parts = append(parts, part)
	}

	for _, x := range commodity.GetMedias() {
		part, err := gemini.NewMediaPart(x.Url, x.MimeType)
		if err != nil {
			return nil, err
		}
		parts = append(parts, part)
	}

	if commodity.Description != "" {
		parts = append(parts, gemini.NewTextPart(commodity.Description))
	}

	prompt, err := t.data.Mongo.Settings.GetPrompt(ctx, "product_analysis")
	if err != nil {
		return nil, err
	}

	parts = append(parts, gemini.NewTextPart(prompt.Content))

	generationConfig := &genai.GenerateContentConfig{
		ResponseMIMEType: "application/json",
		ResponseSchema: &genai.Schema{
			Type:     genai.TypeObject,
			Required: []string{"name", "tags", "brand", "character", "packagingAttributes", "description", "chances"},
			Properties: map[string]*genai.Schema{
				"name": {
					Type: genai.TypeString,
				},
				"tags": {
					Type: genai.TypeArray,
					Items: &genai.Schema{
						Type: genai.TypeString,
					},
				},
				"character": {
					Type:        genai.TypeString,
					Description: `产品质地与属性`,
				},
				"packagingAttributes": {
					Type:        genai.TypeString,
					Description: `外包装属性`,
				},
				"brand": {
					Type:        genai.TypeString,
					Description: "帮我分析图片中的品牌名称",
				},
				"description": {
					Type:        genai.TypeString,
					Description: `帮我分析图片中的产品品牌，产品卖点对应推理出核心人群及具体需求场景`,
				},
				"chances": {
					Type: genai.TypeArray,
					//Description: projpb.PromptCommodityChance,
					Items: &genai.Schema{
						Type:     genai.TypeObject,
						Required: []string{"targetAudience", "sellingPoints"},
						Properties: map[string]*genai.Schema{
							"targetAudience": {
								Type:     genai.TypeObject,
								Required: []string{"description", "tags"},
								Properties: map[string]*genai.Schema{
									"description": {
										Type:        genai.TypeString,
										Description: "受众描述",
									},
									"tags": {
										Type: genai.TypeArray,
										Items: &genai.Schema{
											Type:        genai.TypeString,
											Description: "受众标签",
										},
									},
								},
							},
							"sellingPoints": {
								Type: genai.TypeArray,
								Items: &genai.Schema{
									Type:     genai.TypeObject,
									Required: []string{"description", "tags"},
									Properties: map[string]*genai.Schema{
										"description": {
											Type:        genai.TypeString,
											Description: "卖点描述",
										},
										"tags": {
											Type: genai.TypeArray,
											Items: &genai.Schema{
												Type:        genai.TypeString,
												Description: "卖点标签",
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}

	response, err := genaiClient.GenerateContent(ctx, gemini.GenerateContentRequest{
		Config: generationConfig,
		Parts:  parts,
	})

	if err != nil {
		log.Errorw("GenerateContent err", err)
		return nil, err
	}

	var result AnalyzeResult
	err = json.Unmarshal([]byte(response), &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (t VideoReplication3_CommodityAnalysisJob) searchTemplateSegments(ctx context.Context, keyword, by string, size int) ([]*projpb.ResourceSegment, error) {

	size = helper.Select(size > 0, size, 24)

	var err error
	var items *vikingdb.SearchResponse

	if by == "video" {
		items, err = t.data.VikingDB.SearchByKeywords(ctx, vikingdb.SearchByKeywordsRequest{
			SearchRequest: vikingdb.SearchRequest{
				CollectionName: "segment_video_coll",
				IndexName:      "segment_video_idx",
				Limit:          int(size),
			},
			Keywords: []string{keyword},
		})
	} else {
		items, err = t.data.VikingDB.SearchByKeywords(ctx, vikingdb.SearchByKeywordsRequest{
			SearchRequest: vikingdb.SearchRequest{
				CollectionName: "segment_commodity_coll",
				IndexName:      "segment_commodity_idx",
				Limit:          int(size),
			},
			Keywords: []string{keyword},
		})
	}

	if err != nil {
		return nil, err
	}

	if len(items.Data) == 0 {
		return nil, nil
	}

	var ids []string
	idSort := map[string]int{}
	for i, item := range items.Data {

		if helper.InSlice(item.Id, ids) {
			continue
		}

		ids = append(ids, item.Id)
		idSort[item.Id] = i
	}

	list, _, err := t.data.Mongo.TemplateSegment.ListAndCount(ctx,
		bson.M{"_id": bson.M{"$in": mgz.ObjectIds(ids)}},
		mgz.Find().
			Paging(0, 10).
			B(),
	)
	if err != nil {
		return nil, err
	}

	sort.Slice(list, func(i, j int) bool {
		return idSort[list[i].XId] < idSort[list[j].XId]
	})

	return list, nil
}
