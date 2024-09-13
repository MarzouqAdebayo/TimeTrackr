package boltdb

import (
	"fmt"
	"strings"
	"testing"
	"time"
)

func TestSetupTask(t *testing.T) {
	t.Run("Test to setup test bucket", func(t *testing.T) {
		err := Setup()
		if err != nil {
			t.Errorf("%s", err)
		}
	})
}

func TestFindTasks(t *testing.T) {
	t.Run("Test if it finds task with passed non-id params", func(t *testing.T) {
		test := Task{
			ID:       1,
			Name:     "bowling",
			Category: "sports",
		}
		err := SaveTask(&test)
		if err != nil {
			t.Errorf("%q\n", err.Error())
		}
		findParams := Task{
			Name:     "bowling",
			Category: "sports",
		}
		fmt.Println("saved")
		tasks, err := FindTasks(findParams)
		if err != nil {
			t.Errorf("%q\n", err.Error())
		}
		fmt.Println("found task")
		if strings.Compare(test.Name, tasks[0].Name) != 0 || strings.Compare(test.Category, tasks[0].Category) != 0 {
			got := Task{
				Name:     tasks[0].Name,
				Category: tasks[0].Category,
			}
			t.Errorf("Expected %v, got %v\n", test, got)
		}
	})
}

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
		taskStatus := TaskStatus(ONGOING)
		got, err := GetTaskByValue(taskStatus)
		if err != nil {
			t.Errorf("Failed to get task\n")
			return
		}
		if got.Status != taskStatus {
			t.Errorf("Task Ids are not the same")
			return
		}
	})
}
