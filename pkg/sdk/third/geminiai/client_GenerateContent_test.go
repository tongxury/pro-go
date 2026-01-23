package geminiai

import (
	"context"
	"os"
	"testing"

	"store/confs"

	"google.golang.org/genai"
)

func TestClient_GenerateContent(t *testing.T) {

	ctx := context.Background()

	c := NewGenaiFactory(FactoryConfig{
		Configs: []Config{
			{Proxy: "http://proxy:strOngPAssWOrd@aa404deaba3e54490a68959040f22566-685580509.ap-southeast-1.elb.amazonaws.com:6060", Key: confs.AQKey},
		},
	})

	//url := "https://yoozy.oss-cn-hangzhou.aliyuncs.com/0175810acead1db41b5d2894b126e2ba.jpg"

	client := c.Get()

	//file, err := os.ReadFile("./0175810acead1db41b5d2894b126e2ba.jpg")
	file, err := os.ReadFile("./990f58c9bc3a2dbb58ea5ccfdcdd44.jpg")
	//file, err := os.ReadFile("./c46a3e1ed18a48f916a03a67f40b53.jpg")
	//file, err := os.ReadFile("./ykx.jpg")
	if err != nil {
		return
	}

	url, err := client.UploadBlob(ctx, file, "image/jpeg")
	if err != nil {
		t.Fatal(err)
	}

	//file2, err := os.ReadFile("./02175871564508800000000000000000000ffffac1596066a3260_last-frame.png")
	//if err != nil {
	//	return
	//}
	//
	//url2, err := client.UploadBlob(ctx, file2, "image/jpeg")
	//if err != nil {
	//	t.Fatal(err)
	//}

	content, err := client.GenerateBlob(ctx, GenerateContentRequest{
		Model: "gemini-2.5-flash-image-preview",
		Parts: []genai.Part{
			genai.FileData{
				MIMEType: "image/jpeg",
				URI:      url,
			},
			//genai.FileData{
			//	MIMEType: "image/jpeg",
			//	URI:      url2,
			//},
			//genai.Text("请将图中的商品推广的女子手中的商品换成给定的图中的商品"),
			//genai.Text("请将图中的商品推广的女子的脸和衣服换掉"),
			genai.Text("请将在图中添加一个帅哥"),
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	of, err := os.Create("./out.jpg")
	if err != nil {
		return
	}

	of.Write(content.Data)
}
