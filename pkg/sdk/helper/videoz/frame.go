package videoz

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"

	"github.com/google/uuid"
	ffmpeg_go "github.com/u2takey/ffmpeg-go"
)

func RemoveFirstFrame(src []byte) ([]byte, error) {
	// 创建临时输入文件
	tmpInputFile, err := os.CreateTemp("", fmt.Sprintf("input_%s_*.mp4", uuid.NewString()))
	if err != nil {
		return nil, err
	}
	defer os.Remove(tmpInputFile.Name())
	defer tmpInputFile.Close()

	if _, err := tmpInputFile.Write(src); err != nil {
		return nil, err
	}
	tmpInputFile.Close()

	// 创建临时输出文件
	tmpOutputFile, err := os.CreateTemp("", fmt.Sprintf("output_%s_*.mp4", uuid.NewString()))
	if err != nil {
		return nil, err
	}
	tmpOutputName := tmpOutputFile.Name()
	tmpOutputFile.Close()
	defer os.Remove(tmpOutputName)

	// 使用 ffmpeg 跳过前 0.1 秒
	cmd := exec.Command(
		"ffmpeg",
		"-i", tmpInputFile.Name(),
		"-ss", "0.1",
		"-c:v", "libx264",
		"-c:a", "copy",
		"-y",
		tmpOutputName,
	)

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("ffmpeg error: %v, stderr: %s", err, stderr.String())
	}

	return os.ReadFile(tmpOutputName)
}

func GetFrameByUrl(url string, timestamp float64) ([]byte, error) {
	v, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	all, err := io.ReadAll(v.Body)
	if err != nil {
		return nil, err
	}

	return GetFrame(all, timestamp)
}

func GetFrame(source []byte, timestamp float64) ([]byte, error) {

	ffmpeg_go.LogCompiledCommand = false

	// 创建临时输入文件
	tmpInputFile, err := os.CreateTemp("", fmt.Sprintf("input_%s_*.mp4", uuid.NewString()))

	if err != nil {
		return nil, err
	}
	defer os.Remove(tmpInputFile.Name())
	defer tmpInputFile.Close()

	if _, err := tmpInputFile.Write(source); err != nil {
		return nil, err
	}
	tmpInputFile.Close()

	duration := 0.0
	startTime := 0.0
	// 获取视频时长，确保时间戳不越界
	if metadata, err := GetMetadata(tmpInputFile.Name()); err == nil {
		duration = metadata.Duration
		startTime = metadata.StartTime
		if timestamp < startTime {
			timestamp = startTime
		}
		if duration > 0 && timestamp > duration+startTime {
			timestamp = duration + startTime
		}
	}

	return getFrameRobustly(tmpInputFile.Name(), timestamp, duration+startTime)
}

func getFrameRobustly(inputPath string, timestamp float64, maxTime float64) ([]byte, error) {
	// 尝试多次提取，以防 len(res) == 0
	// 1. 快速定位方式 (Fast seeking)
	res, _ := runFFmpegGetFrame(inputPath, timestamp, true)
	if len(res) > 0 {
		return res, nil
	}

	// 2. 精确提取方式 (Accurate seeking, -ss 放在 -i 之后)
	res, _ = runFFmpegGetFrame(inputPath, timestamp, false)
	if len(res) > 0 {
		return res, nil
	}

	// 3. 容错提取：如果时间点由于某些原因无法提取，尝试前后偏移（寻找最近的一帧）
	// 优先尝试向后偏，再尝试向前偏
	offsets := []float64{0.033, -0.033, 0.1, -0.1, 0.2, -0.2, 0.5, -0.5, 1.0, -1.0}
	for _, offset := range offsets {
		ts := timestamp + offset
		if ts < 0 {
			ts = 0
		}
		if maxTime > 0 && ts > maxTime {
			ts = maxTime
		}
		res, _ = runFFmpegGetFrame(inputPath, ts, false)
		if len(res) > 0 {
			return res, nil
		}
	}

	// 4. 终极兜底：获取视频的第一帧 (按开始时间)
	res, _ = runFFmpegGetFrame(inputPath, 0, true)
	if len(res) > 0 {
		return res, nil
	}

	return nil, fmt.Errorf("empty frame extracted at timestamp %f after all retries", timestamp)
}

func runFFmpegGetFrame(inputPath string, timestamp float64, fast bool) ([]byte, error) {
	var args []string

	// "暴力"切片法核心逻辑：
	// 1. 如果时间点较后，先快速定位到前 1 秒，减少解码开销
	// 2. 使用 select 滤镜逐帧扫描，确保 100% 命中
	// 3. 使用 -vsync 0 禁用同步，防止 VFR 视频丢帧
	seekTime := 0.0
	if timestamp > 1.0 {
		seekTime = timestamp - 1.0
	}

	args = []string{
		"-ss", fmt.Sprintf("%.3f", seekTime),
		"-i", inputPath,
		"-vf", fmt.Sprintf("select='gte(t,%.3f)'", timestamp-seekTime),
		"-frames:v", "1",
		"-vsync", "0",
		"-an",
		"-f", "image2pipe",
		"-vcodec", "mjpeg",
		"-q:v", "2",
		"-",
	}

	cmd := exec.Command("ffmpeg", args...)

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("ffmpeg run error: %v, stderr: %s", err, stderr.String())
	}

	if stdout.Len() == 0 {
		return nil, fmt.Errorf("ffmpeg output empty, stderr: %s", stderr.String())
	}

	return stdout.Bytes(), nil
}

func GetFrames(source []byte, timestamps []float64) ([][]byte, error) {
	// 创建临时输入文件
	tmpInputFile, err := os.CreateTemp("", fmt.Sprintf("input_%s_*.mp4", uuid.NewString()))

	if err != nil {
		return nil, err
	}
	defer os.Remove(tmpInputFile.Name())
	defer tmpInputFile.Close()

	if _, err := tmpInputFile.Write(source); err != nil {
		return nil, err
	}
	tmpInputFile.Close()

	var duration float64
	if metadata, err := GetMetadata(tmpInputFile.Name()); err == nil {
		duration = metadata.Duration
	}

	frames := make([][]byte, 0, len(timestamps))

	for _, timestamp := range timestamps {
		if duration > 0 && timestamp > duration {
			timestamp = duration
		}
		// 使用更加鲁棒的提取逻辑
		res, _ := getFrameRobustly(tmpInputFile.Name(), timestamp, duration)
		if len(res) > 0 {
			frames = append(frames, res)
		}
	}

	return frames, nil
}
