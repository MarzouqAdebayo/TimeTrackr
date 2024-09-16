package boltdb

import (
	"testing"
	"time"
)

func TestStartTask(t *testing.T) {
	t.Run("Creates new task", func(t *testing.T) {
		err := StartTask("testTask", "fun")

		if err != nil {
			t.Error("Error in start task test")
		}
	})
}

func TestStopCurrentTask(t *testing.T) {
	t.Run("Stops currently running task", func(t *testing.T) {
		_, err := StopCurrentTask()

		if err != nil {
			t.Error(err.Error())
		}
	})
}

func TestTaskStatus(t *testing.T) {
	t.Run("Get the status of the currently running task", func(t *testing.T) {
		_, err := Filter("ongoing")

		if err != nil {
			// t.Error("Error: Could not fetch task status")
		}
	})
}

func TestReportGen(t *testing.T) {
	t.Run("Get status generation report", func(t *testing.T) {
	})
}
