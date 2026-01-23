package service

import (
	"context"
	"encoding/json"
	projpb "store/api/proj"
	"store/pkg/clients/mgz"
	"store/pkg/sdk/helper/wg"
	"store/pkg/sdk/third/gemini"

	"github.com/go-kratos/kratos/v2/log"
	"go.mongodb.org/mongo-driver/bson"
	"google.golang.org/genai"
)

func (t ProjService) AnalyzeChances() {

	ctx := context.Background()

	list, err := t.data.Mongo.Commodity.List(ctx, bson.M{"status": "created"})
	if err != nil {
		log.Errorw("List err", err)
		return
	}

	if len(list) == 0 {
		return
	}

	wg.WaitGroup(ctx, list, t.analyzeChances)
}

func (t ProjService) analyzeChances(ctx context.Context, x *projpb.Commodity) error {

	log.Debugw("preparing", x.Url, "x", x)

	// 分析商品
	analyzeResult, err := t.doAnalyzeChances(ctx, x)
	if err != nil {
		return err
	}

	//map[string]any{
	//	"name": analyzeResult.Name,
	//	"tags": analyzeResult.Tags,
	//	//"title":       metadata.Title,
	//	"brand":       analyzeResult.Brand,
	//	"description": analyzeResult.Description,
	//	//"images":      helper.SubSlice(metadata.Images, 10),
	//
	//	"chances": analyzeResult.Chances,
	//	"status":  "completed",
	//})

	_, _ = t.data.Mongo.Commodity.UpdateByIDIfExists(ctx, x.XId,
		mgz.Op().
			Set("name", analyzeResult.Name).
			Set("tags", analyzeResult.Tags).
			Set("brand", analyzeResult.Brand).
			Set("description", analyzeResult.Description).
			Set("chances", analyzeResult.Chances).
			Set("character", analyzeResult.Character).
			Set("packagingAttributes", analyzeResult.PackagingAttributes).
			Set("chances", analyzeResult.Chances).
			Set("status", "completed"),
	)

	return nil
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

func (t ProjService) doAnalyzeChances(ctx context.Context, commodity *projpb.Commodity) (*AnalyzeResult, error) {

	log.Debugw("analyze", "", "commodity", commodity)

	genaiClient := t.data.GenaiFactory.Get()

	var parts []*genai.Part
	for i, x := range commodity.Images {

		if i > 10 {
			break
		}
		//
		//genaiUrl, err := genaiClient.UploadFile(ctx, x, "image/jpeg")
		//if err != nil {
		//	return nil, err
		//}

		part, err := gemini.NewImagePart(x)
		if err != nil {
			return nil, err
		}

		parts = append(parts, part) //&genai.Part{
		//	gemini.New
		//	//FileData: &genai.FileData{
		//	//	MIMEType: "image/jpeg",
		//	//	FileURI:  genaiUrl,
		//	//},
		//},

	}

	for _, x := range commodity.GetMedias() {
		//genaiUrl, err := genaiClient.UploadFile(ctx, x.Url, x.MimeType)
		//if err != nil {
		//	return nil, err
		//}

		part, err := gemini.NewMediaPart(x.Url, x.MimeType)
		if err != nil {
			return nil, err
		}

		parts = append(parts, part)
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
					Description: "帮我分析图片中的品牌名称",
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
		//Model:  "gemini-2.5-pro",
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
