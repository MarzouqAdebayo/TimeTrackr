package boltdb

import (
	"testing"
	"time"
)

func TestSave(t *testing.T) {
	t.Run("Saves task to db", func(t *testing.T) {
		newTask := Task{
			ID:        "00001",
			Name:      "Task00001",
			Category:  "",
			StartTime: time.Now().Unix(),
			EndTime:   time.Now().Unix(),
			Status:    TaskStatus(ONGOING),
		}
		err := SaveTask(&newTask)
		if err != nil {
			t.Errorf("Failed to save task\n")
		}
	})

	t.Run("Updates task in db", func(t *testing.T) {
		taskID := "00001"
		got, err := GetTask(taskID)
		if err != nil {
			t.Errorf("Task with ID %s does not exist", taskID)
		}
		got.Status = TaskStatus(PAUSED)
		saveErr := SaveTask(&got)
		if saveErr != nil {
			t.Errorf("Failed to save task\n")
		}
	})
}

func TestGet(t *testing.T) {
	t.Run("Gets task from db", func(t *testing.T) {
		taskID := "00001"
		got, err := GetTask(taskID)
		if err != nil {
			t.Errorf("Failed to get task\n")
		}
		if got.ID != taskID {
			t.Errorf("Task Ids are not the same")
		}
		// fmt.Println(got)
	})

	t.Run("Gets task from db by status", func(t *testing.T) {
		taskStatus := TaskStatus(PAUSED)
		got, err := GetTaskByValue(taskStatus)
		if err != nil {
			t.Errorf("Failed to get task\n")
		}
		if got.Status != taskStatus {
			t.Errorf("Task Ids are not the same")
		}
	})
}
