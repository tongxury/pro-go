package biz

const (
	WorkflowStatusCreated   = "created"
	WorkflowStatusRunning   = "running"
	WorkflowStatusFailed    = "failed"
	WorkflowStatusCompleted = "completed"
	WorkflowStatusCanceled  = "canceled"
	WorkflowStatusPaused    = "paused"
)

const (
	JobStatusWaiting    = "waiting"
	JobStatusRunning    = "running"
	JobStatusConfirming = "confirming"
	JobStatusFailed     = "failed"
	JobStatusCompleted  = "completed"
	JobStatusCanceled   = "canceled"
)

const (
	ExecuteStatusWaiting   = "waiting"
	ExecuteStatusCompleted = "completed"
	ExecuteStatusFailed    = "failed"
	ExecuteStatusRunning   = "running"
	ExecuteStatusReviewing = "reviewing"
)
