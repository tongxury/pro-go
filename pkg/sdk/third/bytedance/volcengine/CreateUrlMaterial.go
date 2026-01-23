package volcengine

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"store/pkg/sdk/conv"
	"store/pkg/sdk/helper"
	"time"
)

// {"ResponseMetadata":{"RequestId":"202509262141244824E809D49354050785","Action":"GetMediaInfo","Version":"2022-02-01","Service":"iccloud_muse","Region":"cn-north","Code":0},"Result":{"MediaInfos":[{"BasicInfo":{"MediaId":"7257375916521865268","MediaType":"material","MediaStatus":4,"MediaFirstCategory":"video","Owner":{"Id":7061823634582585357,"Type":"PERSON","VolcUsername":"wujiashuo.tutu"},"Name":"3-3 (1).mp4","CreateTime":1689739505,"UpdateTime":1689739505,"Tags":[],"PreviewUrl":"https://museapaas.aigc-cloud.com/api/storage/objects/media/7257375916521865268_transcode_video_mp4_480p_h264.mp4?infer_mime=ext\u0026x-muse-token=ChYIARABGIWw2sYGIIXT38YGKgQ2MTI0EqYBCg4IroCBxrvFkNVoEKjFDxINCgt2b2xjX2VuZ2luZRoCGAIiCAoCbGYSAmNuKgYI2jQQhgMyKQoneyJiaXpfaWQiOiIxMDAwMDAwMDEiLCJkZXJpdmVkX3VpZCI6IiJ9QhQKEmljYXJjaC5pYW0uYXV0aHJwY0ouL2FwaS9zdG9yYWdlL29iamVjdHMvbWVkaWEvNzI1NzM3NTkxNjUyMTg2NTI2OBrYAkoyaGtiYnd4OVVmenRCUW1CKzcydUFKWHpBUkNSSGVONEszN1VBc1RUaHY2Mi9DN1BscFpKQTUvZG1jNWdSSGJTMTVVUDMwZ2RCKzBraFhCbHNrcWxKRytTRXd1dmtrUzFIZ2pqbFZjemcraU9DbFY5emRjbHdVbC9LckVJOG1jRkVhYm5nd1ZhSnFlL1NJQ05NRjY5bFQrZ1plNTRmczBqbU9tYis2ZnlWeDdvWWNQN0dxVkhoemFEaElQWWdHaTlubStOOUZBSGZkZ3pqVEVLb3NRbzA0ZklVZGwwcDJtakNJZEJYdllaNjJnMWw0WVdYN2doc0kxQnhzSkMwN2xaTnI4RHFwNEVRRlpTSmUzZmNpSUVpR0ZiMDV2ekpsVEVFM3lWUE5HamVhbGtlSG0ybVZqa2NGVDVaQ0NtYVBLVUhnYkdDbXNwTjR4ZDVuYlRjWWtMQT09","SourceFrom":"copy_from_private","Description":"","Copyright":"商用版权"},"VideoMedia":{"VideoId":"7257375916521865268","Name":"3-3 (1).mp4","Owner":{"Id":7061823634582585357,"Type":"PERSON","VolcUsername":"wujiashuo.tutu"},"CreateTime":1689739505,"UpdateTime":1689739505,"DownloadUrl":"https://museapaas.aigc-cloud.com/api/storage/objects/media/7257375916521865268_origin.mp4?infer_mime=ext\u0026x-muse-token=ChYIARABGIWw2sYGIIXT38YGKgQ2MTI0EqYBCg4IroCBxrvFkNVoEKjFDxINCgt2b2xjX2VuZ2luZRoCGAIiCAoCbGYSAmNuKgYI2jQQhgMyKQoneyJiaXpfaWQiOiIxMDAwMDAwMDEiLCJkZXJpdmVkX3VpZCI6IiJ9QhQKEmljYXJjaC5pYW0uYXV0aHJwY0ouL2FwaS9zdG9yYWdlL29iamVjdHMvbWVkaWEvNzI1NzM3NTkxNjUyMTg2NTI2OBrYAkoyaGtiYnd4OVVmenRCUW1CKzcydUFKWHpBUkNSSGVONEszN1VBc1RUaHY2Mi9DN1BscFpKQTUvZG1jNWdSSGJTMTVVUDMwZ2RCKzBraFhCbHNrcWxKRytTRXd1dmtrUzFIZ2pqbFZjemcraU9DbFY5emRjbHdVbC9LckVJOG1jRkVhYm5nd1ZhSnFlL1NJQ05NRjY5bFQrZ1plNTRmczBqbU9tYis2ZnlWeDdvWWNQN0dxVkhoemFEaElQWWdHaTlubStOOUZBSGZkZ3pqVEVLb3NRbzA0ZklVZGwwcDJtakNJZEJYdllaNjJnMWw0WVdYN2doc0kxQnhzSkMwN2xaTnI4RHFwNEVRRlpTSmUzZmNpSUVpR0ZiMDV2ekpsVEVFM3lWUE5HamVhbGtlSG0ybVZqa2NGVDVaQ0NtYVBLVUhnYkdDbXNwTjR4ZDVuYlRjWWtMQT09","TranscodeDownloadUrls":{"attach":"","certificate":"","dynamic_poster":"","hls_480p_h264":"https://museapaas.aigc-cloud.com/api/storage/objects/media/7257375916521865268_transcode_video_hls_480p_h264/main.m3u8?infer_mime=ext\u0026x-muse-token=ChYIARABGIWw2sYGIIXT38YGKgQ2MTI0EqYBCg4IroCBxrvFkNVoEKjFDxINCgt2b2xjX2VuZ2luZRoCGAIiCAoCbGYSAmNuKgYI2jQQhgMyKQoneyJiaXpfaWQiOiIxMDAwMDAwMDEiLCJkZXJpdmVkX3VpZCI6IiJ9QhQKEmljYXJjaC5pYW0uYXV0aHJwY0ouL2FwaS9zdG9yYWdlL29iamVjdHMvbWVkaWEvNzI1NzM3NTkxNjUyMTg2NTI2OBrYAkoyaGtiYnd4OVVmenRCUW1CKzcydUFKWHpBUkNSSGVONEszN1VBc1RUaHY2Mi9DN1BscFpKQTUvZG1jNWdSSGJTMTVVUDMwZ2RCKzBraFhCbHNrcWxKRytTRXd1dmtrUzFIZ2pqbFZjemcraU9DbFY5emRjbHdVbC9LckVJOG1jRkVhYm5nd1ZhSnFlL1NJQ05NRjY5bFQrZ1plNTRmczBqbU9tYis2ZnlWeDdvWWNQN0dxVkhoemFEaElQWWdHaTlubStOOUZBSGZkZ3pqVEVLb3NRbzA0ZklVZGwwcDJtakNJZEJYdllaNjJnMWw0WVdYN2doc0kxQnhzSkMwN2xaTnI4RHFwNEVRRlpTSmUzZmNpSUVpR0ZiMDV2ekpsVEVFM3lWUE5HamVhbGtlSG0ybVZqa2NGVDVaQ0NtYVBLVUhnYkdDbXNwTjR4ZDVuYlRjWWtMQT09","icon":"","jpeg_1080p":"","jpeg_480p":"","mp3_b64":"","mp4_1080p_h264":"","mp4_480p_h264":"https://museapaas.aigc-cloud.com/api/storage/objects/media/7257375916521865268_transcode_video_mp4_480p_h264.mp4?infer_mime=ext\u0026x-muse-token=ChYIARABGIWw2sYGIIXT38YGKgQ2MTI0EqYBCg4IroCBxrvFkNVoEKjFDxINCgt2b2xjX2VuZ2luZRoCGAIiCAoCbGYSAmNuKgYI2jQQhgMyKQoneyJiaXpfaWQiOiIxMDAwMDAwMDEiLCJkZXJpdmVkX3VpZCI6IiJ9QhQKEmljYXJjaC5pYW0uYXV0aHJwY0ouL2FwaS9zdG9yYWdlL29iamVjdHMvbWVkaWEvNzI1NzM3NTkxNjUyMTg2NTI2OBrYAkoyaGtiYnd4OVVmenRCUW1CKzcydUFKWHpBUkNSSGVONEszN1VBc1RUaHY2Mi9DN1BscFpKQTUvZG1jNWdSSGJTMTVVUDMwZ2RCKzBraFhCbHNrcWxKRytTRXd1dmtrUzFIZ2pqbFZjemcraU9DbFY5emRjbHdVbC9LckVJOG1jRkVhYm5nd1ZhSnFlL1NJQ05NRjY5bFQrZ1plNTRmczBqbU9tYis2ZnlWeDdvWWNQN0dxVkhoemFEaElQWWdHaTlubStOOUZBSGZkZ3pqVEVLb3NRbzA0ZklVZGwwcDJtakNJZEJYdllaNjJnMWw0WVdYN2doc0kxQnhzSkMwN2xaTnI4RHFwNEVRRlpTSmUzZmNpSUVpR0ZiMDV2ekpsVEVFM3lWUE5HamVhbGtlSG0ybVZqa2NGVDVaQ0NtYVBLVUhnYkdDbXNwTjR4ZDVuYlRjWWtMQT09","mp4_540p_h264":"","origin":"https://museapaas.aigc-cloud.com/api/storage/objects/media/7257375916521865268_origin.mp4?infer_mime=ext\u0026x-muse-token=ChYIARABGIWw2sYGIIXT38YGKgQ2MTI0EqYBCg4IroCBxrvFkNVoEKjFDxINCgt2b2xjX2VuZ2luZRoCGAIiCAoCbGYSAmNuKgYI2jQQhgMyKQoneyJiaXpfaWQiOiIxMDAwMDAwMDEiLCJkZXJpdmVkX3VpZCI6IiJ9QhQKEmljYXJjaC5pYW0uYXV0aHJwY0ouL2FwaS9zdG9yYWdlL29iamVjdHMvbWVkaWEvNzI1NzM3NTkxNjUyMTg2NTI2OBrYAkoyaGtiYnd4OVVmenRCUW1CKzcydUFKWHpBUkNSSGVONEszN1VBc1RUaHY2Mi9DN1BscFpKQTUvZG1jNWdSSGJTMTVVUDMwZ2RCKzBraFhCbHNrcWxKRytTRXd1dmtrUzFIZ2pqbFZjemcraU9DbFY5emRjbHdVbC9LckVJOG1jRkVhYm5nd1ZhSnFlL1NJQ05NRjY5bFQrZ1plNTRmczBqbU9tYis2ZnlWeDdvWWNQN0dxVkhoemFEaElQWWdHaTlubStOOUZBSGZkZ3pqVEVLb3NRbzA0ZklVZGwwcDJtakNJZEJYdllaNjJnMWw0WVdYN2doc0kxQnhzSkMwN2xaTnI4RHFwNEVRRlpTSmUzZmNpSUVpR0ZiMDV2ekpsVEVFM3lWUE5HamVhbGtlSG0ybVZqa2NGVDVaQ0NtYVBLVUhnYkdDbXNwTjR4ZDVuYlRjWWtMQT09","origin_m3u8":"","png_1080p":"","png_480p":"","poster":"https://museapaas.aigc-cloud.com/api/storage/objects/media/7257375112331083812_transcode_image_jpeg_480p.jpeg?infer_mime=ext\u0026x-muse-token=ChYIARABGIWw2sYGIIXT38YGKgQ2MTI0EqYBCg4IroCBxrvFkNVoEKjFDxINCgt2b2xjX2VuZ2luZRoCGAIiCAoCbGYSAmNuKgYI2jQQhgMyKQoneyJiaXpfaWQiOiIxMDAwMDAwMDEiLCJkZXJpdmVkX3VpZCI6IiJ9QhQKEmljYXJjaC5pYW0uYXV0aHJwY0ouL2FwaS9zdG9yYWdlL29iamVjdHMvbWVkaWEvNzI1NzM3NTExMjMzMTA4MzgxMhrYAlpyS1ZuWnh3dURoaFVWeVI5U1Q3Ym52NS9KQjVLc2tBaFdBc1k3NnFBMDlOVXEwNnlaclBuUHRGOENrNE9Qams1ZmM4SWc0SWF5ZS9GSG1lamUydk9ZUW1TWDhYMXBuajhlNDIwcm0xRHgrWjBOajRMaXU1TFZzYyt6c0JyZUVxMGQxMmRWM2hyYnVpR1dZMU9EQkd6R2N4Yk9kRFAveEYrSjZHRFo2aThJUnBCeGFpNm8wY0VVZUpTLzhaZy9oeXd5bEE4WG5EN3Ava1ZyTEg4UTk5RnBqQWt2TDBrN24wNGdVb2RiaXVHYlQvdUQvWDVabnR6enJxU2FQb1l2cGYzQnRrOGF0Y29DM2YwQk83K1lEWjdYSWFiMzVWUW1WVHorNnlSakc4Mkg0M3NiNDV0Q1diOSs0TXQ1S2VHeGFBOWpxNm5oSHVLeDRUdzJma2daMDNZQT09","preview":"https://museapaas.aigc-cloud.com/api/storage/objects/media/7257375916521865268_transcode_video_mp4_480p_h264.mp4?infer_mime=ext\u0026x-muse-token=ChYIARABGIWw2sYGIIXT38YGKgQ2MTI0EqYBCg4IroCBxrvFkNVoEKjFDxINCgt2b2xjX2VuZ2luZRoCGAIiCAoCbGYSAmNuKgYI2jQQhgMyKQoneyJiaXpfaWQiOiIxMDAwMDAwMDEiLCJkZXJpdmVkX3VpZCI6IiJ9QhQKEmljYXJjaC5pYW0uYXV0aHJwY0ouL2FwaS9zdG9yYWdlL29iamVjdHMvbWVkaWEvNzI1NzM3NTkxNjUyMTg2NTI2OBrYAkoyaGtiYnd4OVVmenRCUW1CKzcydUFKWHpBUkNSSGVONEszN1VBc1RUaHY2Mi9DN1BscFpKQTUvZG1jNWdSSGJTMTVVUDMwZ2RCKzBraFhCbHNrcWxKRytTRXd1dmtrUzFIZ2pqbFZjemcraU9DbFY5emRjbHdVbC9LckVJOG1jRkVhYm5nd1ZhSnFlL1NJQ05NRjY5bFQrZ1plNTRmczBqbU9tYis2ZnlWeDdvWWNQN0dxVkhoemFEaElQWWdHaTlubStOOUZBSGZkZ3pqVEVLb3NRbzA0ZklVZGwwcDJtakNJZEJYdllaNjJnMWw0WVdYN2doc0kxQnhzSkMwN2xaTnI4RHFwNEVRRlpTSmUzZmNpSUVpR0ZiMDV2ekpsVEVFM3lWUE5HamVhbGtlSG0ybVZqa2NGVDVaQ0NtYVBLVUhnYkdDbXNwTjR4ZDVuYlRjWWtMQT09","zip":""},"Layout":0,"Resolution":"","CoverUrl":"https://museapaas.aigc-cloud.com/api/storage/objects/media/7257375112331083812_transcode_image_jpeg_480p.jpeg?infer_mime=ext\u0026x-muse-token=ChYIARABGIWw2sYGIIXT38YGKgQ2MTI0EqYBCg4IroCBxrvFkNVoEKjFDxINCgt2b2xjX2VuZ2luZRoCGAIiCAoCbGYSAmNuKgYI2jQQhgMyKQoneyJiaXpfaWQiOiIxMDAwMDAwMDEiLCJkZXJpdmVkX3VpZCI6IiJ9QhQKEmljYXJjaC5pYW0uYXV0aHJwY0ouL2FwaS9zdG9yYWdlL29iamVjdHMvbWVkaWEvNzI1NzM3NTExMjMzMTA4MzgxMhrYAlpyS1ZuWnh3dURoaFVWeVI5U1Q3Ym52NS9KQjVLc2tBaFdBc1k3NnFBMDlOVXEwNnlaclBuUHRGOENrNE9Qams1ZmM4SWc0SWF5ZS9GSG1lamUydk9ZUW1TWDhYMXBuajhlNDIwcm0xRHgrWjBOajRMaXU1TFZzYyt6c0JyZUVxMGQxMmRWM2hyYnVpR1dZMU9EQkd6R2N4Yk9kRFAveEYrSjZHRFo2aThJUnBCeGFpNm8wY0VVZUpTLzhaZy9oeXd5bEE4WG5EN3Ava1ZyTEg4UTk5RnBqQWt2TDBrN24wNGdVb2RiaXVHYlQvdUQvWDVabnR6enJxU2FQb1l2cGYzQnRrOGF0Y29DM2YwQk83K1lEWjdYSWFiMzVWUW1WVHorNnlSakc4Mkg0M3NiNDV0Q1diOSs0TXQ1S2VHeGFBOWpxNm5oSHVLeDRUdzJma2daMDNZQT09","MediaStatus":4,"PreviewVideo":"https://museapaas.aigc-cloud.com/api/storage/objects/media/7257375916521865268_transcode_video_mp4_480p_h264.mp4?infer_mime=ext\u0026x-muse-token=ChYIARABGIWw2sYGIIXT38YGKgQ2MTI0EqYBCg4IroCBxrvFkNVoEKjFDxINCgt2b2xjX2VuZ2luZRoCGAIiCAoCbGYSAmNuKgYI2jQQhgMyKQoneyJiaXpfaWQiOiIxMDAwMDAwMDEiLCJkZXJpdmVkX3VpZCI6IiJ9QhQKEmljYXJjaC5pYW0uYXV0aHJwY0ouL2FwaS9zdG9yYWdlL29iamVjdHMvbWVkaWEvNzI1NzM3NTkxNjUyMTg2NTI2OBrYAkoyaGtiYnd4OVVmenRCUW1CKzcydUFKWHpBUkNSSGVONEszN1VBc1RUaHY2Mi9DN1BscFpKQTUvZG1jNWdSSGJTMTVVUDMwZ2RCKzBraFhCbHNrcWxKRytTRXd1dmtrUzFIZ2pqbFZjemcraU9DbFY5emRjbHdVbC9LckVJOG1jRkVhYm5nd1ZhSnFlL1NJQ05NRjY5bFQrZ1plNTRmczBqbU9tYis2ZnlWeDdvWWNQN0dxVkhoemFEaElQWWdHaTlubStOOUZBSGZkZ3pqVEVLb3NRbzA0ZklVZGwwcDJtakNJZEJYdllaNjJnMWw0WVdYN2doc0kxQnhzSkMwN2xaTnI4RHFwNEVRRlpTSmUzZmNpSUVpR0ZiMDV2ekpsVEVFM3lWUE5HamVhbGtlSG0ybVZqa2NGVDVaQ0NtYVBLVUhnYkdDbXNwTjR4ZDVuYlRjWWtMQT09","SecondaryEdited":false,"MediaMetaInfo":{"Size":5922791,"Md5":"ECE0ACF993670179E26AAB1308D12883","Height":1280,"Width":720,"Duration":10380,"Bitrate":4564771,"Mime":"","VideoStreamInfoSet":[],"AudioStreamInfoSet":[]},"FolderId":0,"FolderName":"","FullPath":null}}]}}
func (t *Client) CreateUrlMaterial(ctx context.Context, params CreateMaterialParams) (*CreateMaterialResponse, error) {

	bytes, err := t.doRequest(ctx,
		Req{
			Version: "2022-02-01",
			Action:  "CreateUrlMaterial",
			Method:  "POST",
			Body: conv.M2B(map[string]any{
				"Owner": map[string]any{
					"Type": "PERSON",
					"Id":   7541913281535950894,
				},
				"CreateUrlMaterialInfo": map[string]any{
					"Title":              params.Title,
					"MediaFirstCategory": params.MediaFirstCategory,
					"MaterialUrl":        params.MaterialUrl,
				},
			}),
		},
	)

	var resp Response[CreateMaterialResponse]
	err = json.Unmarshal(bytes, &resp)
	if err != nil {
		return nil, err
	}

	if resp.Code != 0 || resp.ResponseMetadata.Code != 0 {
		return nil, errors.New(resp.Message)
	}

	if params.Wait {

		loop := 1

		for {

			time.Sleep(1 * time.Second)
			info, err := t.GetMediaInfo(ctx, GetMediaInfoParams{
				MediaIds:  []string{resp.Result.MediaId},
				MediaType: 1,
			})
			if err != nil {
				return nil, err
			}

			//fmt.Println("checking info", info.MediaInfos[0].BasicInfo.MediaId, info.MediaInfos[0].BasicInfo.MediaStatus)

			medias := helper.Filter(info.MediaInfos, func(x MediaInfo) bool {
				return x.BasicInfo.MediaId == resp.Result.MediaId
			})

			if len(medias) == 0 {
				return nil, fmt.Errorf("no media info for %s", resp.Result.MediaId)
			}

			if medias[0].BasicInfo.MediaStatus == 1 {
				return nil, fmt.Errorf("media upload faild: %v, %v", conv.S2J(info), conv.S2J(resp))
			}

			if medias[0].BasicInfo.MediaStatus == 2 {
				//return &CreateMaterialResponse{medias[0].BasicInfo.MediaId}, nil
				continue
			}

			if medias[0].BasicInfo.MediaStatus == 3 {
				continue
			}

			if medias[0].BasicInfo.MediaStatus == 4 {
				return &CreateMaterialResponse{medias[0].BasicInfo.MediaId}, nil
			}

			if medias[0].BasicInfo.MediaStatus == 6 {
				continue
			}

			if loop > 300 {
				return nil, fmt.Errorf("waiting for media faild: %v, %v", conv.S2J(info), conv.S2J(resp))
			}
			loop++
		}

	}

	return &resp.Result, nil
}

type CreateMaterialParams struct {
	MediaFirstCategory string
	Title              string
	MaterialUrl        string
	Wait               bool
}

type CreateMaterialResponse struct {
	MediaId string
}
