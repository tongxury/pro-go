package service

import (
	"context"
	responsepb "store/api/public/response"
	toolkitpb "store/api/toolkit"
	"store/app/usercenter/internal/biz"
	"store/app/usercenter/internal/data"
	"store/pkg/sdk/third/bytedance/tos"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
)

type ToolkitService struct {
	toolkitpb.UnimplementedToolkitServiceServer
	data *data.Data
	item *biz.ItemBiz
}

func NewToolkitService(data *data.Data, item *biz.ItemBiz) *ToolkitService {
	return &ToolkitService{
		data: data,
		item: item,
	}
}

func (t *ToolkitService) ToolKits(ctx context.Context, request *toolkitpb.ToolKitRequest) (*responsepb.MapResponse, error) {

	ctx = context.Background()

	var d map[string]string

	switch request.Method {
	case "getDouyinVideoUrl":
		url := request.Params["url"]
		video, err := t.data.Tikhub.DouyinGetVideoByShareUrl(ctx, url)
		if err != nil {
			return nil, errors.BadRequest("invalidUrl", "")
		}
		urls := video.GetPlayAddr()

		d = map[string]string{
			"url": urls[0],
		}
		return &responsepb.MapResponse{Data: d}, nil

	case "uploadDouyinVideo":

		url := request.Params["url"]

		video, err := t.data.Tikhub.DouyinGetVideoByShareUrl(ctx, url)
		if err != nil {
			return nil, err
		}
		url = video.AwemeDetail.Video.DownloadAddr.UrlList[0]

		log.Debugw("url", video.AwemeDetail.Video.DownloadAddr)

		//url = "https://v5-hl-szyd-ov.zjcdn.com/5c0da246ed7b20a90ceeea29219b3e12/691d4ddc/video/tos/cn/tos-cn-ve-15/ogAhGLGY7BBUfRH7zkIg1pCADFyA2rvK5ReU0f/?a=6383&ch=26&cr=3&dr=0&lr=all&cd=0%7C0%7C0%7C3&cv=1&br=2440&bt=2440&cs=0&ds=3&ft=tQiBtKJqRv.sn1C~pWHghwy35fnobLPhFV-U_4MP26bJNv7TGW&mime_type=video_mp4&qs=0&rc=PDs6Zmg7aDk1Omg7NWYzN0BpamdkOnQ5cmkzNjMzNGkzM0AwMmIuYS5hNmExLzEzMGJjYSMvbV5wMmRjXm1hLS1kLTBzcw%3D%3D&btag=80000e00030000&cquery=100w_100B_100x_100z_100o&dy_q=1763517051&feature_id=dc6e471fd69cf67451806384e04f2b47&l=2025111909505116892F85B7A7F154AE8A"

		//fmt.Println("url:", url)

		//url = "https://yoozy.oss-cn-hangzhou.aliyuncs.com/a18afd25c51b10d50e5c1dd581de027c.mp4"
		//url := "https://api.amemv.com/aweme/v1/play/?video_id=v03033g10000d4bb7mvog65v4jeuu34g&line=1&ratio=720p&watermark=0&media_type=4&vr_type=0&improve_bitrate=0&biz_sign=NRPI1lC4AlkgB1YZ_o3JxvAcRCRKvePEcuaXpNxjeWuwTcZpH81bld8CoeDT9VaaOzLa7CA0Ou-V1V2BgVkRfNQW5E7GDGbNHaUZIZuaRMUGyO7s5wbVW7et0I2xs1hA&logo_name=aweme&source=PackSourceEnum_AWEME_DETAIL"

		url, err = t.data.TOS.PutByUrl(ctx, tos.PutByUrlRequest{
			Bucket: "yoozyres",
			Url:    url,
			Suffix: ".mp4",
		})
		if err != nil {
			return nil, err
		}

		d = map[string]string{
			"url": url,
		}
	}

	return &responsepb.MapResponse{Data: d}, nil
}
