package boltdb

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

func Setup() error {
	return CreateBucket()
}

func StartTask(taskName string, category string) error {
	ongoingErr := OngoingExists()
	if ongoingErr != nil {
		return ongoingErr
	}
	taskExistsErr := TaskExists(taskName)
	if taskExistsErr != nil {
		return taskExistsErr
	}
	newTask := &Task{
		ID:        "00001",
		Name:      taskName,
		Category:  category,
		StartTime: time.Now().Unix(),
		EndTime:   0,
		Status:    TaskStatus(ONGOING),
	}
	saveErr := SaveTask(newTask)
	if saveErr != nil {
		return saveErr
	}
	return nil
}

func StopCurrentTask() (string, error) {
	task, err := GetTaskByValue(TaskStatus(ONGOING))
	if err != nil && errors.Is(err, ErrTaskNotFound) {
		return "", fmt.Errorf("No ongoing task")
	}
	if err != nil {
		return "", err
	}
	task.EndTime = time.Now().Unix()
	task.Status = TaskStatus(COMPLETED)
	saveErr := SaveTask(&task)
	if saveErr != nil {
		return "", err
	}
	return formatTaskStatus(task), nil
}

func PauseCurrentTask() (string, error) {
	task, err := GetTaskByValue(TaskStatus(ONGOING))
	if err != nil && errors.Is(err, ErrTaskNotFound) {
		return "", fmt.Errorf("No ongoing task")
	}
	if err != nil {
		return "", err
	}
	task.Status = TaskStatus(PAUSED)
	task.Duration = task.Duration + time.Now().Unix() - task.StartTime
	saveErr := SaveTask(&task)
	if saveErr != nil {
		return "", err
	}
	return formatTaskStatus(task), nil
}

func ContinuePausedTask(taskName string) (string, error) {
	return "", nil
}

//FIX Create a once and for all gettask function that get a task based on different criteria
func FindTask() {}

func formatTaskStatus(task Task) string {
	startTime := time.Unix(task.StartTime, 0).Format("2006-01-02 15:04:05")
	endTime := "N/A"
	if task.EndTime > 0 {
		endTime = time.Unix(task.EndTime, 0).Format("2006-01-02 15:04:05")
	}

	var duration int64
	if strings.Compare(string(task.Status), string(TaskStatus(ONGOING))) == 0 {
		duration = time.Now().Unix() - task.StartTime
	} else if strings.Compare(string(task.Status), string(TaskStatus(PAUSED))) == 0 {
		duration = task.Duration + time.Now().Unix() - task.StartTime
	} else {
		duration = task.EndTime - task.StartTime
	}

	return fmt.Sprintf(
		"\nTask Status:\n"+
			"---------------------------------\n"+
			"ID        : %s\n"+
			"Name      : %s\n"+
			"Category  : %s\n"+
			"Status    : %s\n"+
			"Start Time: %s\n"+
			"End Time  : %s\n"+
			"Duration  : %d seconds\n"+
			"---------------------------------\n",
		task.ID,
		task.Name,
		task.Category,
		task.Status,
		startTime,
		endTime,
		duration,
	)
}

func Status(status string) (string, error) {
	statusMap := map[string]TaskStatus{
		"ongoing":   TaskStatus(ONGOING),
		"completed": TaskStatus(COMPLETED),
		"paused":    TaskStatus(PAUSED),
	}
	task, err := GetTaskByValue(statusMap[status])
	if err != nil {
		return "", err
	}
	return formatTaskStatus(task), nil
}
