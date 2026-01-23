package volcengine

import (
	"context"
	"fmt"
	"hash/crc32"
	"os"
	"store/pkg/sdk/helper"
	"testing"
)

func TestClient_UploadFlow(t *testing.T) {
	ctx := context.Background()
	c := NewClient()

	// 1. 准备测试数据
	content := []byte("hello world volcengine upload test " + fmt.Sprint(helper.MD5([]byte("test"))))
	size := int64(len(content))
	md5Str := helper.MD5(content)
	crcValue := crc32.ChecksumIEEE(content)

	owner := &Owner{
		Type: "PERSON",
		Id:   7348381768925446179, // 使用 Python 示例中的 ID
	}

	fmt.Printf("File Info: Size=%d, MD5=%s, CRC=%d\n", size, md5Str, crcValue)

	// 2. 获取上传状态
	stateRes, err := c.GetUploadState(ctx, GetUploadStateRequest{
		Owner: owner,
		Md5:   md5Str,
		Crc:   crcValue,
		Size:  size,
		Start: 0,
		End:   size - 1,
	})
	if err != nil {
		t.Fatalf("GetUploadState failed: %v", err)
	}

	fmt.Printf("Upload State: SkipDataComplete=%v\n", stateRes.SkipDataComplete)

	// 3. 如果需要上传数据
	if !stateRes.SkipDataComplete {
		err = c.StreamUploadData(ctx, StreamUploadDataRequest{
			Md5:       md5Str,
			Size:      size,
			Offset:    0,
			OwnerId:   owner.Id,
			OwnerType: owner.Type,
			Data:      content,
		})
		if err != nil {
			t.Fatalf("StreamUploadData failed: %v", err)
		}
		fmt.Println("StreamUploadData success")
	}

	// 4. 创建媒资
	createRes, err := c.CreateMaterial(ctx, CreateMaterialRequest{
		Owner: owner,
		StoreItem: &StoreItem{
			Md5:              md5Str,
			Size:             size,
			SkipDataComplete: true,
			Filename:         "test_upload_go.txt",
			FileExtension:    "txt",
		},
		CreateMaterialInfo: &CreateMaterialInfo{
			Title:              "test_upload_go",
			MediaType:          1,
			MediaFirstCategory: "docx",
			Tags:               []string{"test", "go"},
			MediaExtension:     "txt",
		},
	})
	if err != nil {
		t.Fatalf("CreateMaterial failed: %v", err)
	}

	fmt.Printf("CreateMaterial success, MediaId: %s\n", createRes.MediaId)
}

func TestClient_CreateBytesMaterial(t *testing.T) {
	ctx := context.Background()
	c := NewClient()

	content, err := os.ReadFile("7257375112331083812_transcode_image_jpeg_480p.jpeg")
	if err != nil {
		t.Fatalf("ReadFile failed: %v", err)
	}

	info := CreateMaterialInfo{
		Title:              "test_image_upload",
		MediaType:          2, // 2 通常代表图片
		MediaFirstCategory: "image",
		Tags:               []string{"test", "image"},
		MediaExtension:     "jpeg",
	}

	res, err := c.CreateBytesMaterial(ctx, content, info)
	if err != nil {
		t.Fatalf("CreateBytesMaterial failed: %v", err)
	}

	fmt.Printf("CreateBytesMaterial success, MediaId: %s\n", res.MediaId)
}
