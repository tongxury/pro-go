package videoz

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/google/uuid"
	ffmpeg "github.com/u2takey/ffmpeg-go"
)

type Segment struct {
	From    float64
	To      float64
	Content []byte
	//FirstFrame []byte
	//LastFrame  []byte
}

type SegmentParams struct {
	From            float64
	To              float64
	FirstValidFrame int
}

func GetSegment(ctx context.Context, source []byte, from, to float64) (*Segment, error) {

	batch, err := SplitInBatch(ctx, source, []SegmentParams{{
		From: from,
		To:   to,
	}})
	if err != nil {
		return nil, err
	}

	if len(batch) == 0 {
		return nil, fmt.Errorf("failed to split")
	}

	return &batch[0], nil
}

func GetBytes(url string) ([]byte, error) {
	v, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	all, err := io.ReadAll(v.Body)
	if err != nil {
		return nil, err
	}

	return all, nil
}

func GetSegmentByUrl(ctx context.Context, url string, from, to float64) (*Segment, error) {

	v, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	all, err := io.ReadAll(v.Body)
	if err != nil {
		return nil, err
	}

	batch, err := SplitInBatch(ctx, all, []SegmentParams{{
		From: from,
		To:   to,
	}})
	if err != nil {
		return nil, err
	}

	if len(batch) == 0 {
		return nil, fmt.Errorf("failed to split")
	}

	return &batch[0], nil
}

func Split(ctx context.Context, source []byte, from, to float64) (*Segment, error) {

	batch, err := SplitInBatch(ctx, source, []SegmentParams{{
		From: from,
		To:   to,
	}})
	if err != nil {
		return nil, err
	}

	if len(batch) == 0 {
		return nil, fmt.Errorf("failed to split")
	}

	return &batch[0], nil
}

func SplitInBatch(ctx context.Context, source []byte, segments []SegmentParams) ([]Segment, error) {

	ffmpeg.LogCompiledCommand = false

	tempInputFile, err := os.CreateTemp("", fmt.Sprintf("input_%s_*.mp4", uuid.NewString()))
	if err != nil {
		return nil, fmt.Errorf("创建临时文件失败: %v", err)
	}
	defer os.Remove(tempInputFile.Name())
	defer tempInputFile.Close()

	//fmt.Println("tempInputFile", tempInputFile.Name())

	_, err = tempInputFile.Write(source)
	if err != nil {
		return nil, fmt.Errorf("写入临时文件失败: %v", err)
	}
	tempInputFile.Close()

	var result []Segment
	for _, x := range segments {

		// 完整视频流
		tempOutputFile, err := os.CreateTemp("", fmt.Sprintf("output_%s_*.mp4", uuid.NewString()))
		if err != nil {
			return nil, fmt.Errorf("创建临时输出文件失败: %v", err)
		}
		defer os.Remove(tempOutputFile.Name())
		defer tempOutputFile.Close()

		//fmt.Println("tempOutputFile", tempOutputFile.Name())

		duration := x.To - x.From
		err = ffmpeg.Input(
			tempInputFile.Name(),
			ffmpeg.KwArgs{"ss": x.From}).
			Output(
				tempOutputFile.Name(),
				//fmt.Sprintf("%d.mp4", i),

				ffmpeg.KwArgs{
					"t":      duration, // 持续时间
					"c":      "copy",   // 复制流
					"format": "mp4",    // 输出格式
				}).
			OverWriteOutput().
			Run()

		if err != nil {
			return nil, err
		}

		file, err := os.ReadFile(tempOutputFile.Name())
		if err != nil {
			return nil, err
		}

		result = append(result, Segment{
			From:    x.From,
			To:      x.To,
			Content: file,
			//FirstFrame: firstFrameFile,
			//LastFrame:  lastFrameFile,
		})
	}

	return result, nil
}
