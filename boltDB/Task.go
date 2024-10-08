package boltdb

type TaskStatus string

const (
	DEFAULT_CATEGORY = "miscellaneous"
)

const (
	ONGOING   TaskStatus = "ongoing"
	COMPLETED            = "completed"
	PAUSED               = "paused"
)

type Task struct {
	ID        int        `json:"id,string"`
	Name      string     `json:"name"`
	Category  string     `json:"category"`
	StartTime int64      `json:"startTime"`
	UpdatedAt int64      `json:"updatedAt"`
	EndTime   int64      `json:"endTime"`
	Duration  int64      `json:"duration"`
	Status    TaskStatus `json:"status"`
}
