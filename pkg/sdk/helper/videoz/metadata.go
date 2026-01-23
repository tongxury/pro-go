package videoz

import (
	"bytes"
	"encoding/json"
	"io"
	"os/exec"
	"strconv"
)

type Metadata struct {
	Duration  float64
	StartTime float64
	Width     int
	Height    int
}

func GetMetadata(url string) (*Metadata, error) {
	cmd := exec.Command(
		"ffprobe",
		"-v", "error",
		"-show_entries", "format=duration,start_time:stream=width,height",
		"-of", "json",
		url,
	)

	var stdout bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = io.Discard

	if err := cmd.Run(); err != nil {
		return nil, err
	}

	var result struct {
		Streams []struct {
			Width  int `json:"width"`
			Height int `json:"height"`
		} `json:"streams"`
		Format struct {
			Duration  string `json:"duration"`
			StartTime string `json:"start_time"`
		} `json:"format"`
	}

	if err := json.Unmarshal(stdout.Bytes(), &result); err != nil {
		return nil, err
	}

	duration, _ := strconv.ParseFloat(result.Format.Duration, 64)
	startTime, _ := strconv.ParseFloat(result.Format.StartTime, 64)

	width, height := 0, 0
	if len(result.Streams) > 0 {
		width = result.Streams[0].Width
		height = result.Streams[0].Height
	}

	return &Metadata{
		StartTime: startTime,
		Duration:  duration,
		Width:     width,
		Height:    height,
	}, nil
}
