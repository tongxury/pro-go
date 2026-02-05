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
