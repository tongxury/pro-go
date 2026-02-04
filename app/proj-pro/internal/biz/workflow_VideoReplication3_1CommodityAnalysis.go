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

	_, err = t.data.Mongo.Commodity.UpdateByIDIfExists(ctx, commodity.XId,
		mgz.Op().
			Set("name", analyzeResult.Name).
			Set("tags", analyzeResult.Tags).
			Set("brand", analyzeResult.Brand).
			Set("description", analyzeResult.Description).
			Set("chances", analyzeResult.Chances).
			Set("character", analyzeResult.Character).
			Set("packagingAttributes", analyzeResult.PackagingAttributes).
			Set("status", "completed"),
	)
	if err != nil {
		logger.Errorw("update commodity err", err)
		return nil, err
	}

	segments, err := t.searchTemplateSegments(ctx, strings.Join(analyzeResult.Tags, ","), "", 1)
	if err != nil {
		return nil, err
	}

	segments = helper.Filter(segments, func(param *projpb.ResourceSegment) bool {
		if dataBus.GetSettings().GetDuration() == 0 {
			return true
		}
		return int64(param.TimeEnd-param.TimeStart) < dataBus.GetSettings().GetDuration()
	})

	if len(segments) == 0 {
		logger.Errorw("searchTemplateSegments err", "no segment found")
		return &ExecuteResult{
			Status: ExecuteStatusFailed,
			Error:  "no segment found",
		}, errors.New("no segment template found")
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

	generationConfig := &genai.GenerateContentConfig{
		ResponseMIMEType: "application/json",
		ResponseSchema: &genai.Schema{
			Type:     genai.TypeObject,
			Required: []string{"name", "tags", "brand", "character", "packagingAttributes", "description", "chances"},
			Properties: map[string]*genai.Schema{
				"name": {
					Type:        genai.TypeString,
					Description: projpb.PromptCommodityName,
				},
				"tags": {
					Type: genai.TypeArray,
					Items: &genai.Schema{
						Type: genai.TypeString,
					},
					Description: projpb.PromptCommodityTags,
				},
				"character": {
					Type: genai.TypeString,
					Description: `
分析产品质地与属性:
示例:
视觉质地:	表面特征（光滑度、光泽度、孔隙等）	面包的孔洞、牛奶的浓稠度、粉底的光泽感
触觉质地:	物体在接触时的物理感受	润肤霜的丝滑感、饼干的酥脆感、牛奶的浓稠度
听觉质地:	声音反馈（脆性、易碎感）	薯片的咔嚓声、饼干的碎裂声
`,
				},
				"packagingAttributes": {
					Type: genai.TypeString,
					Description: `
外包装属性。
示例:
纸质包装	环保、亲和力强、易回收	"塑料太不环保了"	特写：手撕纸盒+回收标志+孩子开心笑
玻璃瓶	纯净感、高端、可重复使用	"玻璃瓶装的才放心"	慢镜头：阳光透过玻璃瓶，牛奶清澈透亮
金属罐	保鲜、防潮、高端	"罐装比袋装更保质"	实验：金属罐vs塑料袋牛奶对比
塑料瓶	轻便、防摔、成本低	"出门带不破"	场景：孩子蹦跳，瓶子完好无损
复合材料	多功能（阻氧+保鲜+环保）	"既要保质又要环保"	对比：普通包装vs复合包装保质期测试
`,
				},
				"brand": {
					Type:        genai.TypeString,
					Description: "帮我分析图片中的 brand 品牌名称",
				},
				"description": {
					Type:        genai.TypeString,
					Description: `帮我分析图片中的产品品牌，产品卖点对应推理出核心人群及具体需求场景`,
				},
				"chances": {
					Type:        genai.TypeArray,
					Description: projpb.PromptCommodityChance,
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
