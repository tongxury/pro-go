package biz

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	projpb "store/api/proj"
	"store/app/proj-pro/internal/data"
	"store/pkg/clients/mgz"
	"store/pkg/sdk/helper"
	"store/pkg/sdk/helper/imagez"
	"store/pkg/sdk/helper/videoz"
	"store/pkg/sdk/third/bytedance/tos"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/openai/openai-go/v3"
)

type VideoGeneration_VideoGenerationJob struct {
	data *data.Data
}

func (t VideoGeneration_VideoGenerationJob) Initialize(ctx context.Context, options Options) error {
	return nil
}

func (t VideoGeneration_VideoGenerationJob) GetName() string {
	return "videoGenerationJob"
}

func (t VideoGeneration_VideoGenerationJob) Execute(ctx context.Context, jobState *projpb.Job, wfState *projpb.Workflow) (s *ExecuteResult, err error) {

	dataBus := GetDataBus(wfState)

	logger := log.NewHelper(log.With(log.DefaultLogger,
		"func", "VideoGeneration_VideoGenerationJob.Execute",
		"workflowId ", wfState.XId,
		"jobState.Name ", jobState.Name,
		"jobState.Index ", jobState.Index,
		"jobState.Status ", jobState.Status,
	))

	segmentScript := dataBus.SegmentScript
	if segmentScript == nil {
		return nil, errors.New("segmentScript not found")
	}

	if len(dataBus.VideoGenerations) == 0 {

		t.data.Mongo.Workflow.UpdateByIDIfExists(ctx, wfState.XId,
			mgz.Op().Set(
				fmt.Sprintf("jobs.%d.dataBus.videoGenerations", jobState.Index),
				[]*projpb.VideoGeneration{
					{
						Status: "running",
					},
				},
			),
		)

		return nil, nil
	}

	videoGeneration := dataBus.VideoGenerations[0]

	if videoGeneration.Url != "" {
		return &ExecuteResult{
			Status: ExecuteStatusCompleted,
			//SkipConfirm: true,
		}, nil
	}

	if videoGeneration.TaskId != "" {

		v, err := t.data.OpenAI.Videos().Get(ctx, videoGeneration.TaskId)
		if err != nil {
			logger.Errorw("OpenAI.PollStatus err", err, "processingVideoId", videoGeneration.TaskId)
			return nil, err
		}

		logger.Debugw("PollStatus ", "", "videoId", v.ID, "status", v.Status, "process", v.Progress)

		// 失败了重新生成视频
		if v.Status == openai.VideoStatusFailed {
			return &ExecuteResult{
				Status: ExecuteStatusFailed,
			}, nil
		}

		// 成功了结束任务
		if v.Status == openai.VideoStatusCompleted {
			cr, err := t.data.OpenAI.Videos().DownloadContent(ctx, v.ID, openai.VideoDownloadContentParams{
				Variant: "video",
			})
			if err != nil {
				return nil, err
			}

			content, err := io.ReadAll(cr.Body)
			if err != nil {
				return nil, err
			}

			content, err = videoz.RemoveFirstFrame(content)
			if err != nil {
				return nil, err
			}
			//url, err := t.data.Alioss.UploadBytes(ctx, helper.MD5(content)+".mp4", content)
			url, err := t.data.TOS.Put(ctx, tos.PutRequest{
				Bucket:  "yoozyres",
				Content: content,
				Key:     helper.MD5(content) + ".mp4",
			})
			if err != nil {
				return nil, err
			}

			logger.Debugw("Upload video ", url)

			coverFrame, err := videoz.GetFrame(content, 1)
			if err != nil {
				log.Errorw("GetFrame err", err)
				return nil, err
			}

			logger.Debugw("GetCoverFrame url", url)

			coverUrl, err := t.data.TOS.Put(ctx, tos.PutRequest{
				Bucket:  "yoozyres",
				Content: coverFrame,
				Key:     helper.MD5(content) + ".jpg",
			})
			if err != nil {
				log.Errorw("Put err", err)
				return nil, err
			}

			_, err = t.data.Mongo.Workflow.UpdateByIDIfExists(ctx,
				wfState.XId,
				mgz.Op().
					SetListItem(fmt.Sprintf("jobs.%d.dataBus.videoGenerations", jobState.Index), 0, "url", url).
					SetListItem(fmt.Sprintf("jobs.%d.dataBus.videoGenerations", jobState.Index), 0, "status", ExecuteStatusCompleted).
					SetListItem(fmt.Sprintf("jobs.%d.dataBus.videoGenerations", jobState.Index), 0, "coverUrl", coverUrl),
			)
			if err != nil {
				log.Errorw("UpdateByIDXX err", err)
				return nil, err
			}

			logger.Debugw("Update status ", "success")

			return &ExecuteResult{
				Status:      ExecuteStatusCompleted,
				SkipConfirm: true,
			}, nil
		}

		return nil, nil
	}

	seconds := "12"

	var inputReference io.Reader
	if len(dataBus.SegmentScript.GetImages()) > 0 {

		imgBytes, err := helper.ReadBytesByUrl(dataBus.SegmentScript.GetImages()[0])
		if err != nil {
			return nil, err
		}

		resizedImage, err := imagez.ResizeKeepRatio(imgBytes, 720, 1280)
		if err != nil {
			return nil, err
		}

		inputReference = openai.File(bytes.NewReader(resizedImage), helper.CreateUUID(), "image/jpeg")
	}

	video, err := t.data.OpenAI.Videos().New(ctx, openai.VideoNewParams{
		Model: openai.VideoModelSora2,
		//Prompt:         dataBus.SegmentScript.GetScript(),
		Prompt:         helper.OrString(videoGeneration.Prompt, dataBus.SegmentScript.GetScript()),
		InputReference: inputReference,
		Seconds:        openai.VideoSeconds(seconds),
	})

	if err != nil {
		logger.Errorw("start generate error", err)

		if video != nil && video.Error.Code == "rate_limit_exceeded" {
			time.Sleep(1 * time.Minute)
		}

		return nil, err
	}

	_, err = t.data.Mongo.Workflow.UpdateByIDIfExists(ctx, wfState.XId, mgz.Op().
		SetListItem(fmt.Sprintf("jobs.%d.dataBus.videoGenerations", jobState.Index), 0, "taskId", video.ID))
	if err != nil {
		return nil, err
	}

	return nil, nil
}
