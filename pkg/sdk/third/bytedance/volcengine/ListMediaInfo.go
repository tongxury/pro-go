package volcengine

import (
	"context"
	"encoding/json"
	"errors"
	"store/pkg/sdk/conv"
)

// {"ResponseMetadata":{"RequestId":"2025091218261822503CC3C7D04C797E80","Action":"CreateUrlMaterial","Version":"2022-02-01","Service":"iccloud_muse","Region":"cn-north","Code":0},"Result":{"MediaId":"7549146192932864051"}}
// 7549146192932864051
func (t *Client) ListMediaInfo(ctx context.Context, params ListMediaInfoParams) (*ListMediaInfoResult, error) {

	bytes, err := t.doRequest(ctx,
		Req{
			Version: "2022-02-01",
			Action:  "ListMediaInfo",
			Method:  "POST",
			Body: conv.M2B(map[string]any{
				"MediaType": 1,
			}),
		},
	)
	var resp Response[ListMediaInfoResult]
	err = json.Unmarshal(bytes, &resp)
	if err != nil {
		return nil, err
	}

	if resp.Code != 0 {
		return nil, errors.New(resp.Message)
	}

	return &resp.Result, nil
}

type ListMediaInfoParams struct {
}

type ListMediaInfoResult struct {
	MediaInfos []MediaInfo `json:"MediaInfos"`
	PageInfo   PageInfo    `json:"PageInfo"`
}

type MediaInfo struct {
	BasicInfo struct {
		MediaId   string `json:"MediaId"`
		MediaType string `json:"MediaType"`
		/*
			1
			上传失败
			2
			上传完成
			3
			转码中
			4
			转码完成
			5
			转码失败
			6
			上传中
		*/
		MediaStatus        int    `json:"MediaStatus"`
		MediaFirstCategory string `json:"MediaFirstCategory"`
		Owner              struct {
			Id           int64  `json:"Id"`
			Type         string `json:"Type"`
			VolcUsername string `json:"VolcUsername"`
		} `json:"Owner"`
		Name        string        `json:"Name"`
		CreateTime  int           `json:"CreateTime"`
		UpdateTime  int           `json:"UpdateTime"`
		Tags        []interface{} `json:"Tags"`
		PreviewUrl  string        `json:"PreviewUrl"`
		SourceFrom  string        `json:"SourceFrom"`
		Description string        `json:"Description"`
	} `json:"BasicInfo"`
	VideoMedia struct {
		VideoId string `json:"VideoId"`
		Name    string `json:"Name"`
		Owner   struct {
			Id           int64  `json:"Id"`
			Type         string `json:"Type"`
			VolcUsername string `json:"VolcUsername"`
		} `json:"Owner"`
		CreateTime            int    `json:"CreateTime"`
		UpdateTime            int    `json:"UpdateTime"`
		DownloadUrl           string `json:"DownloadUrl"`
		TranscodeDownloadUrls struct {
			DynamicPoster string `json:"dynamic_poster"`
			Hls480PH264   string `json:"hls_480p_h264"`
			Icon          string `json:"icon"`
			Mp41080PH264  string `json:"mp4_1080p_h264"`
			Mp4480PH264   string `json:"mp4_480p_h264"`
			Mp4540PH264   string `json:"mp4_540p_h264"`
			Origin        string `json:"origin"`
			Poster        string `json:"poster"`
			Preview       string `json:"preview"`
		} `json:"TranscodeDownloadUrls"`
		Layout          int    `json:"Layout"`
		Resolution      string `json:"Resolution"`
		CoverUrl        string `json:"CoverUrl"`
		MediaStatus     int    `json:"MediaStatus"`
		PreviewVideo    string `json:"PreviewVideo"`
		SecondaryEdited bool   `json:"SecondaryEdited"`
		MediaMetaInfo   struct {
			Size               int    `json:"Size"`
			Md5                string `json:"Md5"`
			Height             int    `json:"Height"`
			Width              int    `json:"Width"`
			Duration           int    `json:"Duration"`
			Bitrate            int    `json:"Bitrate"`
			Mime               string `json:"Mime"`
			VideoStreamInfoSet []struct {
				Bitrate int    `json:"Bitrate"`
				Height  int    `json:"Height"`
				Width   int    `json:"Width"`
				Codec   string `json:"Codec"`
				Size    int    `json:"Size"`
				Md5     string `json:"Md5"`
				Fps     int    `json:"Fps"`
			} `json:"VideoStreamInfoSet"`
			AudioStreamInfoSet []interface{} `json:"AudioStreamInfoSet"`
		} `json:"MediaMetaInfo"`
		FolderId   int         `json:"FolderId"`
		FolderName string      `json:"FolderName"`
		FullPath   interface{} `json:"FullPath"`
	} `json:"VideoMedia"`
}
type PageInfo struct {
	PageNum   int `json:"PageNum"`
	PageSize  int `json:"PageSize"`
	TotalNum  int `json:"TotalNum"`
	TotalPage int `json:"TotalPage"`
}

type T struct {
	ResponseMetadata struct {
		RequestId string `json:"RequestId"`
		Action    string `json:"Action"`
		Version   string `json:"Version"`
		Service   string `json:"Service"`
		Region    string `json:"Region"`
		Code      int    `json:"Code"`
	} `json:"ResponseMetadata"`
	Result struct {
		MediaInfos []struct {
			BasicInfo struct {
				MediaId            string `json:"MediaId"`
				MediaType          string `json:"MediaType"`
				MediaStatus        int    `json:"MediaStatus"`
				MediaFirstCategory string `json:"MediaFirstCategory"`
				Owner              struct {
					Id           int64  `json:"Id"`
					Type         string `json:"Type"`
					VolcUsername string `json:"VolcUsername"`
				} `json:"Owner"`
				Name        string        `json:"Name"`
				CreateTime  int           `json:"CreateTime"`
				UpdateTime  int           `json:"UpdateTime"`
				Tags        []interface{} `json:"Tags"`
				PreviewUrl  string        `json:"PreviewUrl"`
				SourceFrom  string        `json:"SourceFrom"`
				Description string        `json:"Description"`
			} `json:"BasicInfo"`
			VideoMedia struct {
				VideoId string `json:"VideoId"`
				Name    string `json:"Name"`
				Owner   struct {
					Id           int64  `json:"Id"`
					Type         string `json:"Type"`
					VolcUsername string `json:"VolcUsername"`
				} `json:"Owner"`
				CreateTime            int    `json:"CreateTime"`
				UpdateTime            int    `json:"UpdateTime"`
				DownloadUrl           string `json:"DownloadUrl"`
				TranscodeDownloadUrls struct {
					DynamicPoster string `json:"dynamic_poster"`
					Hls480PH264   string `json:"hls_480p_h264"`
					Icon          string `json:"icon"`
					Mp41080PH264  string `json:"mp4_1080p_h264"`
					Mp4480PH264   string `json:"mp4_480p_h264"`
					Mp4540PH264   string `json:"mp4_540p_h264"`
					Origin        string `json:"origin"`
					Poster        string `json:"poster"`
					Preview       string `json:"preview"`
				} `json:"TranscodeDownloadUrls"`
				Layout          int    `json:"Layout"`
				Resolution      string `json:"Resolution"`
				CoverUrl        string `json:"CoverUrl"`
				MediaStatus     int    `json:"MediaStatus"`
				PreviewVideo    string `json:"PreviewVideo"`
				SecondaryEdited bool   `json:"SecondaryEdited"`
				MediaMetaInfo   struct {
					Size               int     `json:"Size"`
					Md5                string  `json:"Md5"`
					Height             int     `json:"Height"`
					Width              int     `json:"Width"`
					Duration           float64 `json:"Duration"`
					Bitrate            int     `json:"Bitrate"`
					Mime               string  `json:"Mime"`
					VideoStreamInfoSet []struct {
						Bitrate int    `json:"Bitrate"`
						Height  int    `json:"Height"`
						Width   int    `json:"Width"`
						Codec   string `json:"Codec"`
						Size    int    `json:"Size"`
						Md5     string `json:"Md5"`
						Fps     int    `json:"Fps"`
					} `json:"VideoStreamInfoSet"`
					AudioStreamInfoSet []interface{} `json:"AudioStreamInfoSet"`
				} `json:"MediaMetaInfo"`
				FolderId   int         `json:"FolderId"`
				FolderName string      `json:"FolderName"`
				FullPath   interface{} `json:"FullPath"`
			} `json:"VideoMedia"`
		} `json:"MediaInfos"`
		PageInfo struct {
			PageNum   int `json:"PageNum"`
			PageSize  int `json:"PageSize"`
			TotalNum  int `json:"TotalNum"`
			TotalPage int `json:"TotalPage"`
		} `json:"PageInfo"`
	} `json:"Result"`
}
