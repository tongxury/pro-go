package videoz

import (
	"fmt"
	"os"

	"github.com/google/uuid"
	ffmpeg "github.com/u2takey/ffmpeg-go"
)

type EditParams struct {
	From  float64
	To    float64
	Speed float64
}

type DelogoParams struct {
	X      int
	Y      int
	Width  int
	Height int
}

// RemoveWatermark 使用 FFmpeg 的 delogo 滤镜通过模糊指定区域来移除水印
func RemoveWatermark(source []byte, params DelogoParams) ([]byte, error) {
	ffmpeg.LogCompiledCommand = false

	tempInputFile, err := os.CreateTemp("", fmt.Sprintf("input_%s_*.mp4", uuid.NewString()))
	if err != nil {
		return nil, fmt.Errorf("failed to create temp input file: %v", err)
	}
	defer os.Remove(tempInputFile.Name())
	defer tempInputFile.Close()

	if _, err = tempInputFile.Write(source); err != nil {
		return nil, fmt.Errorf("failed to write temp input file: %v", err)
	}
	tempInputFile.Close()

	tempOutputFile, err := os.CreateTemp("", fmt.Sprintf("output_%s_*.mp4", uuid.NewString()))
	if err != nil {
		return nil, fmt.Errorf("failed to create temp output file: %v", err)
	}
	defer os.Remove(tempOutputFile.Name())
	tempOutputFile.Close()

	outputArgs := ffmpeg.KwArgs{
		"format":   "mp4",
		"c:v":      "libx264",
		"c:a":      "copy",
		"preset":   "ultrafast",
		"movflags": "faststart",
	}

	input := ffmpeg.Input(tempInputFile.Name())
	v := input.Video()
	a := input.Audio()

	// 应用 delogo 滤镜
	v = v.Filter("delogo", ffmpeg.Args{
		fmt.Sprintf("x=%d:y=%d:w=%d:h=%d", params.X, params.Y, params.Width, params.Height),
	}, ffmpeg.KwArgs{})

	err = ffmpeg.Output([]*ffmpeg.Stream{v, a}, tempOutputFile.Name(), outputArgs).OverWriteOutput().Run()
	if err != nil {
		// 兼容无音轨视频
		err = ffmpeg.Output([]*ffmpeg.Stream{v}, tempOutputFile.Name(), outputArgs).OverWriteOutput().Run()
		if err != nil {
			return nil, fmt.Errorf("ffmpeg delogo failed: %v", err)
		}
	}

	return os.ReadFile(tempOutputFile.Name())
}

type CropParams struct {
	Width  int
	Height int
	X      int
	Y      int
}

// CropVideo 裁剪视频区域，可用于移除边缘水印
func CropVideo(source []byte, params CropParams) ([]byte, error) {
	ffmpeg.LogCompiledCommand = false

	tempInputFile, err := os.CreateTemp("", fmt.Sprintf("input_%s_*.mp4", uuid.NewString()))
	if err != nil {
		return nil, fmt.Errorf("failed to create temp input file: %v", err)
	}
	defer os.Remove(tempInputFile.Name())
	defer tempInputFile.Close()

	if _, err = tempInputFile.Write(source); err != nil {
		return nil, fmt.Errorf("failed to write temp input file: %v", err)
	}
	tempInputFile.Close()

	tempOutputFile, err := os.CreateTemp("", fmt.Sprintf("output_%s_*.mp4", uuid.NewString()))
	if err != nil {
		return nil, fmt.Errorf("failed to create temp output file: %v", err)
	}
	defer os.Remove(tempOutputFile.Name())
	tempOutputFile.Close()

	outputArgs := ffmpeg.KwArgs{
		"format":   "mp4",
		"c:v":      "libx264",
		"c:a":      "copy",
		"preset":   "ultrafast",
		"movflags": "faststart",
	}

	input := ffmpeg.Input(tempInputFile.Name())
	v := input.Video()
	a := input.Audio()

	// 应用 crop 滤镜
	v = v.Filter("crop", ffmpeg.Args{
		fmt.Sprintf("%d:%d:%d:%d", params.Width, params.Height, params.X, params.Y),
	}, ffmpeg.KwArgs{})

	err = ffmpeg.Output([]*ffmpeg.Stream{v, a}, tempOutputFile.Name(), outputArgs).OverWriteOutput().Run()
	if err != nil {
		err = ffmpeg.Output([]*ffmpeg.Stream{v}, tempOutputFile.Name(), outputArgs).OverWriteOutput().Run()
		if err != nil {
			return nil, fmt.Errorf("ffmpeg crop failed: %v", err)
		}
	}

	return os.ReadFile(tempOutputFile.Name())
}

func EditVideo(source []byte, params EditParams) ([]byte, error) {
	if params.Speed <= 0 {
		params.Speed = 1.0
	}
	// 如果没有任何修改逻辑，直接返回原数据
	if params.Speed == 1.0 && params.From <= 0 && params.To <= 0 {
		return source, nil
	}

	ffmpeg.LogCompiledCommand = false // 开启日志以便排查

	tempInputFile, err := os.CreateTemp("", fmt.Sprintf("input_%s_*.mp4", uuid.NewString()))
	if err != nil {
		return nil, fmt.Errorf("failed to create temp input file: %v", err)
	}
	defer os.Remove(tempInputFile.Name())
	defer tempInputFile.Close()

	if _, err = tempInputFile.Write(source); err != nil {
		return nil, fmt.Errorf("failed to write temp input file: %v", err)
	}
	tempInputFile.Close()

	tempOutputFile, err := os.CreateTemp("", fmt.Sprintf("output_%s_*.mp4", uuid.NewString()))
	if err != nil {
		return nil, fmt.Errorf("failed to create temp output file: %v", err)
	}
	defer os.Remove(tempOutputFile.Name())
	tempOutputFile.Close()

	inputArgs := ffmpeg.KwArgs{}
	if params.From > 0 {
		inputArgs["ss"] = params.From
	}

	outputArgs := ffmpeg.KwArgs{
		"format":   "mp4",
		"c:v":      "libx264",
		"c:a":      "aac",
		"preset":   "ultrafast",
		"movflags": "faststart", // 优化网络播放
	}

	if params.To > 0 && params.To > params.From {
		duration := params.To - params.From
		if params.Speed != 1.0 {
			duration = duration / params.Speed
		}
		outputArgs["t"] = duration
	}

	input := ffmpeg.Input(tempInputFile.Name(), inputArgs)

	// 分离视频流和音频流独立处理，这是 ffmpeg-go 处理滤镜的推荐方式
	v := input.Video()
	a := input.Audio()

	// 处理变速
	if params.Speed != 1.0 {
		v = v.Filter("setpts", ffmpeg.Args{fmt.Sprintf("%f*PTS", 1.0/params.Speed)}, ffmpeg.KwArgs{})

		s := params.Speed
		for s > 2.0 {
			a = a.Filter("atempo", ffmpeg.Args{"2.0"}, ffmpeg.KwArgs{})
			s /= 2.0
		}
		for s < 0.5 {
			a = a.Filter("atempo", ffmpeg.Args{"0.5"}, ffmpeg.KwArgs{})
			s /= 0.5
		}
		a = a.Filter("atempo", ffmpeg.Args{fmt.Sprintf("%f", s)}, ffmpeg.KwArgs{})
	}

	// 合并流并输出
	err = ffmpeg.Output([]*ffmpeg.Stream{v, a}, tempOutputFile.Name(), outputArgs).OverWriteOutput().Run()
	if err != nil {
		// 如果带音频处理失败，尝试仅处理视频（兼容无音轨视频）
		err = ffmpeg.Output([]*ffmpeg.Stream{v}, tempOutputFile.Name(), outputArgs).OverWriteOutput().Run()
		if err != nil {
			return nil, fmt.Errorf("ffmpeg run failed: %v", err)
		}
	}

	return os.ReadFile(tempOutputFile.Name())
}
