package boltdb

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/boltdb/bolt"
)

var fileName = "trackr.db"
var bucketName = []byte("task")

var ErrTaskNotFound = errors.New("task not found")
var ErrBucketDoesNotExist = errors.New("bucket does not exist")

type CustomErr string

func (e CustomErr) Error() string {
	return string(e)
}

func opendb() *bolt.DB {
	db, openErr := bolt.Open(fileName, 0600, &bolt.Options{Timeout: 1 * time.Second})
	if openErr != nil {
		log.Fatal(openErr)
	}
	return db
}

func CreateBucket() error {
	db := opendb()
	defer db.Close()
	updateErr := db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(bucketName)
		if err != nil {
			return err
		}
		return nil
	})
	if updateErr != nil {
		return updateErr
	}
	return nil
}

func TaskExists(taskName string) error {
	db := opendb()
	defer db.Close()

	err := db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(bucketName)
		if bucket == nil {
			return fmt.Errorf("Please run trackr setup to start using TimeTrackr")
		}
		c := bucket.Cursor()
		for k, v := c.First(); k != nil && bytes.Contains(v, []byte(taskName)); k, v = c.Next() {
			return CustomErr(fmt.Sprintf("Task '%s' already exists, please provide a new task name", taskName))
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func OngoingExists() error {
	db := opendb()
	defer db.Close()

	err := db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(bucketName)
		if bucket == nil {
			return fmt.Errorf("Please run trackr setup to start using TimeTrackr")
		}
		c := bucket.Cursor()
		for k, v := c.First(); k != nil && bytes.Contains(v, []byte(TaskStatus(ONGOING))); k, v = c.Next() {
			var task Task
			err := json.Unmarshal(v, &task)
			if err != nil {
				return err
			}
			return CustomErr(fmt.Sprintf("Task '%s' is currently ongoing, stop/pause it before starting a new one", task.Name))
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

// FIX Create a once and for all gettask function that get a task based on different criteria
func FindTasks(obj Task) ([]Task, error) {
	db := opendb()
	defer db.Close()

	var matchingTasks []Task
	err := db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(bucketName)
		if bucket == nil {
			return fmt.Errorf("Please run trackr setup to start using TimeTrackr")
		}
		c := bucket.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			var task Task
			err := json.Unmarshal(v, &task)
			if err != nil {
				return err
			}
			if (obj.Name == "" || task.Name == obj.Name) &&
				(obj.Category == "" || task.Category == obj.Category) &&
				(obj.Status == "" || task.Status == obj.Status) {
				matchingTasks = append(matchingTasks, task)
			}
		}
		return nil
	})
	if err != nil {
		return matchingTasks, err
	}
	if len(matchingTasks) == 0 {
		return nil, ErrTaskNotFound
	}
	return matchingTasks, nil
}

func FilterTasks(startDate, endDate *time.Time, minDuration, maxDuration *int64) ([]Task, error) {
	db := opendb()
	defer db.Close()

	var matchingTasks []Task
	err := db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(bucketName)
		if bucket == nil {
			return fmt.Errorf("Please run trackr setup to start using TimeTrackr")
		}
		c := bucket.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			var task Task
			err := json.Unmarshal(v, &task)
			if err != nil {
				return err
			}
			if (startDate == nil || task.StartTime >= startDate.Unix()) &&
				(endDate == nil || task.EndTime <= endDate.Unix()) &&
				(minDuration == nil || task.Duration >= *minDuration) &&
				(maxDuration == nil || task.Duration <= *maxDuration) {
				matchingTasks = append(matchingTasks, task)
			}
		}
		return nil
	})
	if err != nil {
		return matchingTasks, err
	}
	if len(matchingTasks) == 0 {
		return nil, ErrTaskNotFound
	}
	return matchingTasks, nil
}

func SaveTask(task *Task) error {
	db := opendb()
	defer db.Close()
	updateErr := db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(bucketName)
		if bucket == nil {
			return fmt.Errorf("Please run trackr setup to start using TimeTrackr")
		}
		key := []byte(task.ID)
		taskBuf, err := json.Marshal(task)
		if err != nil {
			return err
		}
		err = bucket.Put(key, taskBuf)
		if err != nil {
			return err
		}
		return nil
	})
	if updateErr != nil {
		log.Fatal(updateErr)
	}
	return nil
}

func GetTask(taskID string) (Task, error) {
	db := opendb()
	defer db.Close()

	var task Task
	err := db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(bucketName)
		if bucket == nil {
			return fmt.Errorf("Please run trackr setup to start using TimeTrackr")
		}
		key := []byte(taskID)
		taskJson := bucket.Get(key)
		if taskJson == nil {
			return nil
		}

		err := json.Unmarshal(taskJson, &task)
		if err != nil {
			return err
		}
		// fmt.Println(task)
		// bucket.ForEach(func(k, v []byte) error {
		// 	t := &Task{}
		// 	err := json.Unmarshal(v, &t)
		// 	if err != nil {
		// 		return nil
		// 	}
		// 	fmt.Printf("key=%s, value=%v with time=%v\n", k, t, t.StartTime)
		// 	return nil
		// })

		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
	return task, nil
}

func GetTaskByValue(taskStatus TaskStatus) (Task, error) {
	db := opendb()
	defer db.Close()

	var task Task
	err := db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(bucketName)
		if bucket == nil {
			return fmt.Errorf("Please run trackr setup to start using TimeTrackr")
		}
		c := bucket.Cursor()
		for k, v := c.First(); k != nil && bytes.Contains(v, []byte(taskStatus)); k, v = c.Next() {
			err := json.Unmarshal(v, &task)
			if err != nil {
				return err
			}
			return nil
		}
		return ErrTaskNotFound
	})
	if err != nil {
		return task, err
	}
	return task, nil
}
