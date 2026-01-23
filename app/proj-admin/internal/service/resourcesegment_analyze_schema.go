package service

type AnalyzeResponse struct {
	Commodity *AnalyzeCommodity `json:"commodity" jsonschema_description:"commodity"`
	Segments  []*AnalyzeSegment `json:"segments" jsonschema_description:"segments"`
}

type AnalyzeCommodity struct {
	Name string   `json:"name" jsonschema_description:"name"`
	Tags []string `json:"tags" jsonschema_description:"tags"`
}

type AnalyzeSegment struct {
	Formula         string                   `json:"formula" jsonschema_description:"formula"`
	Script          string                   `json:"script" jsonschema_description:"script"`
	Style           string                   `json:"style" jsonschema_description:"style"`
	ContentStyle    string                   `json:"contentStyle" jsonschema_description:"contentStyle"`
	SceneStyle      string                   `json:"sceneStyle" jsonschema_description:"sceneStyle"`
	HighlightFrames []*AnalyzeHighlightFrame `json:"highlightFrames" jsonschema_description:"highlightFrames"`
	Description     string                   `json:"description" jsonschema_description:"description"`
	ShootingStyle   string                   `json:"shootingStyle" jsonschema_description:"shootingStyle"`
	TypedTags       *AnalyzeTypedTags        `json:"typedTags" jsonschema_description:"typedTags"`
	Subtitle        string                   `json:"subtitle" jsonschema_description:"subtitle"`
}

type AnalyzeTypedTags struct {
	FocusOn       []string `json:"focusOn" jsonschema_description:"focusOn"`
	Picture       []string `json:"picture" jsonschema_description:"picture"`
	ShootingStyle []string `json:"shootingStyle" jsonschema_description:"shootingStyle"`
	Scene         []string `json:"scene" jsonschema_description:"scene"`
	Action        []string `json:"action" jsonschema_description:"action"`
	Person        []string `json:"person" jsonschema_description:"person"`
	Text          []string `json:"text" jsonschema_description:"text"`
	Emotion       []string `json:"emotion" jsonschema_description:"emotion"`
}

type AnalyzeHighlightFrame struct {
	Timestamp float64 `json:"timestamp" jsonschema_description:"timestamp"`
	Desc      string  `json:"desc" jsonschema_description:"desc"`
}
