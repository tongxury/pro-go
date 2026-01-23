package aiagentpb

import (
	"fmt"
	"store/pkg/sdk/conv"
	"store/pkg/sdk/helper"
	"strings"
)

func (t *AnswerChunks) FullText() string {

	var b strings.Builder
	for _, x := range t.Chunks {
		b.WriteString(x.Text)
	}

	return b.String()
}

type Resources []*ResourceV2

func (ts Resources) IsAllMediaImages() bool {
	for _, t := range ts {

		// 空代表 是媒体文件
		if t.Category != "" {
			continue
		}

		if !strings.HasPrefix(t.MimeType, "image/") {
			return false
		}
	}

	return true
}

func (t *ResourceV2) ContentMap() map[string]any {

	contentMap := make(map[string]any)

	err := conv.J2M(t.Content, &contentMap)
	if err != nil {
		return contentMap
	}

	return contentMap
}

func (t *ResourceV2) PromptText() string {

	return fmt.Sprintf("[category='%s']: <DATA>%s</DATA>", t.Category, t.Content)

	//switch t.Category {
	//case "title":
	//	return fmt.Sprintf("标题: %s", t.Content)
	//default:
	//	return fmt.Sprintf("[category='%s']: <DATA>%s</DATA>", t.Category, t.Content)
	//
	//}
}

func (t *ResourceV2) IsVideo() bool {
	return strings.HasPrefix(t.MimeType, "video/")
}

func (t *ResourceV2) IsImage() bool {
	return strings.HasPrefix(t.MimeType, "image/")
}

func (ts Resources) FindOneByCategory(category string) *ResourceV2 {

	for _, t := range ts {
		if t.Category == category {
			return t
		}
	}

	return nil
}

func (t *Item) Safe() *Item {

	t.GenaiUri = ""
	t.Reports = nil
	return t
}

type Items []*Item

func (ts Items) Safe() Items {
	for _, t := range ts {
		t.Safe()
	}

	return ts
}

func (t *Question) Safe() *Question {

	t.Session.Safe()

	if t.GetPrompt() != nil {
		t.Prompt.SystemContent = ""
	}

	return t
}

//func (t *Prompt) Text() string {
//
//	if t == nil {
//		return ""
//	}
//
//	if t.SystemContent != "" {
//		return t.SystemContent
//	}
//
//	return t.Content
//}

type Sessions []*Session

func (ts Sessions) Safe() Sessions {
	for _, t := range ts {
		t.Safe()
	}

	return ts
}

func (t *Session) Safe() *Session {

	if t.GetResource() != nil {
		//t.Resource.GenaiUri = ""
		t.Resource.PlatformUrl = ""
	}

	for _, x := range t.GetResources() {
		//x.GenaiUri = ""
		x.PlatformUrl = ""
	}

	return t
}

type Questions []*Question

func (ts Questions) Safe() Questions {
	for _, t := range ts {
		t.Safe()
	}

	return ts
}

func (t *Profile) Text() string {

	if t == nil {
		return ""
	}

	var profileInfo = fmt.Sprintf(`
%s
%s
%s关注
%s粉丝
%s获赞与收藏
个性签名: %s
个人标签: %s
视频平台: %s
`,
		t.Username,
		t.IpAddress,
		t.FollowingCount,
		t.FollowerCount,
		t.LikedCount,
		t.Sign,
		strings.Join(t.Tags, ","),
		helper.OrString(t.Platform, "xiaohongshu"),
	)

	return profileInfo
}
