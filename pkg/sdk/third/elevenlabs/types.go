package elevenlabs

// --- Common Types ---

type Cursor string

type Pagination struct {
	HasMore    bool    `json:"has_more"`
	NextCursor *Cursor `json:"next_cursor"`
}

type AccessInfo struct {
	IsCreator    bool   `json:"is_creator"`
	CreatorName  string `json:"creator_name"`
	CreatorEmail string `json:"creator_email"`
	Role         string `json:"role"`
}

// --- Agent Types ---

type Agent struct {
	AgentID            string             `json:"agent_id"`
	Name               string             `json:"name"`
	ConversationConfig ConversationConfig `json:"conversation_config"`
	Metadata           AgentMetadata      `json:"metadata"`
	PlatformSettings   *PlatformSettings  `json:"platform_settings,omitempty"`
	Workflow           *Workflow          `json:"workflow,omitempty"`
	Tags               []string           `json:"tags,omitempty"`
	AccessInfo         *AccessInfo        `json:"access_info,omitempty"`
	Archived           *bool              `json:"archived,omitempty"`
}

type AgentMetadata struct {
	CreatedAtUnixSecs int64 `json:"created_at_unix_secs,omitempty"`
	UpdatedAtUnixSecs int64 `json:"updated_at_unix_secs,omitempty"`
}

type ConversationConfig struct {
	Agent    AgentConfig `json:"agent"`
	ASR      *ASRConfig  `json:"asr,omitempty"`
	Turn     *TurnConfig `json:"turn,omitempty"`
	TTS      *TTSConfig  `json:"tts,omitempty"`
	Language string      `json:"language,omitempty"`
}

type AgentConfig struct {
	Prompt                           *PromptSettings `json:"prompt,omitempty"`
	FirstMessage                     *string         `json:"first_message,omitempty"`
	Language                         *string         `json:"language,omitempty"`
	DisableFirstMessageInterruptions *bool           `json:"disable_first_message_interruptions,omitempty"`
}

type PromptSettings struct {
	Prompt        string          `json:"prompt,omitempty"`
	KnowledgeBase []KnowledgeBase `json:"knowledge_base,omitempty"`
}

type ASRConfig struct {
	Quality              string   `json:"quality,omitempty"`
	Provider             string   `json:"provider,omitempty"`
	UserInputAudioFormat string   `json:"user_input_audio_format,omitempty"`
	Keywords             []string `json:"keywords,omitempty"`
}

type TurnConfig struct {
	TurnTimeout           int                `json:"turn_timeout,omitempty"`
	SilenceEndCallTimeout int                `json:"silence_end_call_timeout,omitempty"`
	TurnEagerness         float64            `json:"turn_eagerness,omitempty"`
	SoftTimeoutConfig     *SoftTimeoutConfig `json:"soft_timeout_config,omitempty"`
}

type SoftTimeoutConfig struct {
	TimeoutSeconds int    `json:"timeout_seconds,omitempty"`
	Message        string `json:"message,omitempty"`
}

type TTSConfig struct {
	ModelID                  string   `json:"model_id,omitempty"`
	VoiceID                  string   `json:"voice_id,omitempty"`
	AgentOutputAudioFormat   string   `json:"agent_output_audio_format,omitempty"`
	OptimizeStreamingLatency *int     `json:"optimize_streaming_latency,omitempty"`
	Stability                *float64 `json:"stability,omitempty"`
	SimilarityBoost          *float64 `json:"similarity_boost,omitempty"`
	Speed                    *float64 `json:"speed,omitempty"`
}

type Tool struct {
	ToolID      string         `json:"tool_id,omitempty"`
	Name        string         `json:"name"`
	Description string         `json:"description,omitempty"`
	Parameters  map[string]any `json:"parameters,omitempty"`
	Type        string         `json:"type,omitempty"` // e.g., "webhook", "client"
}

type KnowledgeBase struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Type string `json:"type,omitempty"`
}

type PlatformSettings struct {
	AuthEnabled bool `json:"auth_enabled"`
}

type Workflow struct {
	Nodes                map[string]any `json:"nodes"`
	Edges                map[string]any `json:"edges"`
	PreventSubagentLoops bool           `json:"prevent_subagent_loops"`
}

type ListAgentsParams struct {
	PageSize      int    `url:"page_size,omitempty"`
	Search        string `url:"search,omitempty"`
	Archived      bool   `url:"archived,omitempty"`
	ShowOnlyOwned bool   `url:"show_only_owned_agents,omitempty"`
	SortBy        string `url:"sort_by,omitempty"`
	SortDirection string `url:"sort_direction,omitempty"`
	Cursor        string `url:"cursor,omitempty"`
}

type ListAgentsResp struct {
	Agents []Agent `json:"agents"`
	Pagination
}

type CreateAgentRequest struct {
	Name               string             `json:"name"`
	ConversationConfig ConversationConfig `json:"conversation_config"`
	PlatformSettings   *PlatformSettings  `json:"platform_settings,omitempty"`
}

// --- Voice Types ---

type Voice struct {
	VoiceID  string         `json:"voice_id"`
	Name     string         `json:"name"`
	Samples  []VoiceSample  `json:"samples,omitempty"`
	Category string         `json:"category,omitempty"`
	Settings *VoiceSettings `json:"settings,omitempty"`
	Sharing  *VoiceSharing  `json:"sharing,omitempty"`
}

type VoiceSample struct {
	SampleID  string `json:"sample_id"`
	FileName  string `json:"file_name"`
	MimeType  string `json:"mime_type"`
	SizeBytes int    `json:"size_bytes"`
	Hash      string `json:"hash"`
}

type VoiceSettings struct {
	Stability       float64 `json:"stability"`
	SimilarityBoost float64 `json:"similarity_boost"`
	Style           float64 `json:"style,omitempty"`
	UseSpeakerBoost bool    `json:"use_speaker_boost,omitempty"`
}

type VoiceSharing struct {
	Status string `json:"status"`
}

type ListVoicesResp struct {
	Voices []Voice `json:"voices"`
}

// --- Conversation Types ---

type ConversationSummary struct {
	AgentID        string `json:"agent_id"`
	ConversationID string `json:"conversation_id"`
	StartTimeUnix  int64  `json:"start_time_unix"`
	Status         string `json:"status"`
	DurationSecs   int    `json:"duration_secs"`
}

type ListConversationsParams struct {
	AgentID         string `url:"agent_id,omitempty"`
	CallStartAfter  int64  `url:"call_start_after,omitempty"`
	CallStartBefore int64  `url:"call_start_before,omitempty"`
	PageSize        int    `url:"page_size,omitempty"`
	Cursor          string `url:"cursor,omitempty"`
}

type ListConversationsResp struct {
	Conversations []ConversationSummary `json:"conversations"`
	Pagination
}

type ConversationDetail struct {
	ConversationID string            `json:"conversation_id"`
	AgentID        string            `json:"agent_id"`
	Status         string            `json:"status"`
	Transcript     []TranscriptEntry `json:"transcript"`
	Metadata       map[string]any    `json:"metadata"`
	Analysis       *ConvAnalysis     `json:"analysis,omitempty"`
}

type ConvAnalysis struct {
	SuccessCriteria map[string]any `json:"success_criteria"`
	DataSummary     string         `json:"data_summary"`
	Transcript      string         `json:"transcript_summary"`
}

type TranscriptEntry struct {
	Role      string `json:"role"`
	Message   string `json:"message"`
	TimeInSec int    `json:"time_in_call_secs"`
}

// --- Knowledge Base Types ---

type KnowledgeDocument struct {
	ID              string      `json:"id"`
	Name            string      `json:"name"`
	Type            string      `json:"type"`
	Metadata        DocMetadata `json:"metadata"`
	SupportedUsages []string    `json:"supported_usages"`
	AccessInfo      AccessInfo  `json:"access_info"`
}

type DocMetadata struct {
	CreatedAtUnixSecs     int64 `json:"created_at_unix_secs"`
	LastUpdatedAtUnixSecs int64 `json:"last_updated_at_unix_secs"`
	SizeBytes             int64 `json:"size_bytes"`
}

type DependentAgent struct {
	AgentID               string   `json:"agent_id"`
	Name                  string   `json:"name"`
	CreatedAtUnixSecs     int64    `json:"created_at_unix_secs"`
	ReferencedResourceIDs []string `json:"referenced_resource_ids"`
}

type ListKnowledgeBaseResp struct {
	Documents []KnowledgeDocument `json:"documents"`
	Pagination
}

type ListDependentAgentsResp struct {
	Agents []DependentAgent `json:"agents"`
	Pagination
}

// --- Speech-to-Text Types ---

type TranscribeRequest struct {
	ModelID string `json:"model_id"`
	File    string `json:"file,omitempty"`
	URL     string `json:"url,omitempty"`
}

type TranscriptionResult struct {
	Text     string  `json:"text"`
	Language string  `json:"language"`
	Duration float64 `json:"duration"`
}
