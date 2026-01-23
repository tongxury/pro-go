package videoz

import (
	"fmt"
	"os"

	"github.com/google/uuid"
	ffmpeg "github.com/u2takey/ffmpeg-go"
)

type ConcatVideoSegment struct {
	Source    []byte
	Speed     float64
	TimeStart float64
	TimeEnd   float64
}

// ConcatVideos 将多段视频二进制数据串联成一段视频，支持对每一段进行裁剪和变速
func ConcatVideos(segments []ConcatVideoSegment) ([]byte, error) {
	if len(segments) == 0 {
		return nil, fmt.Errorf("no segments provided")
	}
	if len(segments) == 1 {
		seg := segments[0]
		// 如果只有一段且需要编辑，则调用 EditVideo
		if (seg.Speed > 0 && seg.Speed != 1.0) || seg.TimeStart > 0 || seg.TimeEnd > 0 {
			return EditVideo(seg.Source, EditParams{
				From:  seg.TimeStart,
				To:    seg.TimeEnd,
				Speed: seg.Speed,
			})
		}
		return seg.Source, nil
	}

	ffmpeg.LogCompiledCommand = false

	// 创建临时目录存放输入文件
	tempFiles := make([]string, 0, len(segments))
	defer func() {
		for _, f := range tempFiles {
			os.Remove(f)
		}
	}()

	var streams []*ffmpeg.Stream

	for i, seg := range segments {
		source := seg.Source
		// 如果该片段需要编辑（裁剪或变速），先进行预处理
		if (seg.Speed > 0 && seg.Speed != 1.0) || seg.TimeStart > 0 || seg.TimeEnd > 0 {
			var err error
			source, err = EditVideo(seg.Source, EditParams{
				From:  seg.TimeStart,
				To:    seg.TimeEnd,
				Speed: seg.Speed,
			})
			if err != nil {
				return nil, fmt.Errorf("failed to edit segment %d: %v", i, err)
			}
		}

		tmpFile, err := os.CreateTemp("", fmt.Sprintf("concat_in_%d_%s_*.mp4", i, uuid.NewString()))
		if err != nil {
			return nil, fmt.Errorf("failed to create temp input file %d: %v", i, err)
		}
		tempFiles = append(tempFiles, tmpFile.Name())

		if _, err = tmpFile.Write(source); err != nil {
			tmpFile.Close()
			return nil, fmt.Errorf("failed to write source %d: %v", i, err)
		}
		tmpFile.Close()

		input := ffmpeg.Input(tmpFile.Name())
		streams = append(streams, input.Video(), input.Audio())
	}

	tempOutputFile, err := os.CreateTemp("", fmt.Sprintf("concat_out_%s_*.mp4", uuid.NewString()))
	if err != nil {
		return nil, fmt.Errorf("failed to create temp output file: %v", err)
	}
	defer os.Remove(tempOutputFile.Name())
	tempOutputFile.Close()

	// 使用 concat 滤镜串联
	// v=1, a=1 表示每个段落有1个视频流和1个音频流
	concatStream := ffmpeg.Concat(streams, ffmpeg.KwArgs{"v": 1, "a": 1})

	outputArgs := ffmpeg.KwArgs{
		"format":   "mp4",
		"c:v":      "libx264",
		"c:a":      "aac",
		"preset":   "ultrafast",
		"movflags": "faststart",
	}

	err = concatStream.Output(tempOutputFile.Name(), outputArgs).OverWriteOutput().Run()
	if err != nil {
		// 尝试不带音频串联（防止某些片段没有音频导致失败）
		var videoOnlyStreams []*ffmpeg.Stream
		for _, f := range tempFiles {
			videoOnlyStreams = append(videoOnlyStreams, ffmpeg.Input(f).Video())
		}
		concatStream = ffmpeg.Concat(videoOnlyStreams, ffmpeg.KwArgs{"v": 1, "a": 0})
		err = concatStream.Output(tempOutputFile.Name(), outputArgs).OverWriteOutput().Run()
		if err != nil {
			return nil, fmt.Errorf("ffmpeg concat failed: %v", err)
		}
	}

	return os.ReadFile(tempOutputFile.Name())
}

// ConcatVideosByUrls 通过 URL 列表将多段视频串联成一段视频
func ConcatVideosByUrls(urls []string) ([]byte, error) {
	if len(urls) == 0 {
		return nil, fmt.Errorf("no urls provided")
	}

	var segments []ConcatVideoSegment
	for _, url := range urls {
		source, err := GetBytes(url)
		if err != nil {
			return nil, fmt.Errorf("failed to download video from %s: %v", url, err)
		}
		segments = append(segments, ConcatVideoSegment{
			Source: source,
		})
	}

	return ConcatVideos(segments)
}
