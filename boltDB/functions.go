package boltdb

import (
	"errors"
	"fmt"
	"time"
)

func Setup() error {
	return CreateBucket()
}

func StartTask(taskName string) error {
	_, err := GetTaskByValue(TaskStatus(ONGOING))
	if err != nil && !errors.Is(err, ErrTaskNotFound) {
		return err
	}
	newTask := &Task{
		ID:        "00001",
		Name:      taskName,
		StartTime: time.Now().Unix(),
		EndTime:   time.Now().Unix(),
		Status:    TaskStatus(ONGOING),
	}
	saveErr := SaveTask(newTask)
	if saveErr != nil {
		return saveErr
	}
	return nil
}

func StopCurrentTask() error {
	task, err := GetTaskByValue(TaskStatus(ONGOING))
	if err != nil {
		return err
	}
	task.Status = TaskStatus(COMPLETED)
	saveErr := SaveTask(&task)
	if saveErr != nil {
		return err
	}
	return nil
}

func PauseCurrentTask() {
}

func formatTaskStatus(task Task) string {
	startTime := time.Unix(task.StartTime, 0).Format("2006-01-02 15:04:05")
	endTime := "N/A"
	if task.EndTime > 0 {
		endTime = time.Unix(task.EndTime, 0).Format("2006-01-02 15:04:05")
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
		task.Duration,
	)
}

func Status() (string, error) {
	task, err := GetTaskByValue(TaskStatus(PAUSED))
	if err != nil {
		return "", err
	}
	return formatTaskStatus(task), nil
}
