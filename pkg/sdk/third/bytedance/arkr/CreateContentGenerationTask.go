package arkr

import (
	"context"

	"github.com/volcengine/volcengine-go-sdk/service/arkruntime/model"
	"github.com/volcengine/volcengine-go-sdk/volcengine"
)

func (t *Client) CreateContentGenerationTask(ctx context.Context) (string, error) {

	req := model.CreateContentGenerationTaskRequest{
		Model: "doubao-seedance-1-0-pro-250528",
		//Model: "doubao-seedance-1-0-lite-i2v-250428",
		Content: []*model.CreateContentGenerationContentItem{
			{
				Type: "text",
				Text: volcengine.String("特写镜头，一只手将一颗透明的DHA胶囊和一颗绿色的AD胶囊放入白色勺子中。背景虚化，可以看到多个伊可新产品包装盒和一个可爱的玩偶，营造出温馨的家庭氛围。 --rs 720p  --dur 5 --cf false"),
			},
			{
				Type: "image_url",
				ImageURL: &model.ImageURL{
					URL: "https://yoozy.oss-cn-hangzhou.aliyuncs.com/%E6%88%AA%E5%B1%8F2025-09-26%2003.50.08.png",
				},
				Role: volcengine.String("first_frame"),
			},
		},
	}

	resp, err := t.c.CreateContentGenerationTask(ctx, req)
	if err != nil {
		return "", err
	}

	return resp.ID, nil
}
