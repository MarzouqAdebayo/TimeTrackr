package boltdb

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/boltdb/bolt"
)

var fileName = "trackr.db"
var bucketName = []byte("task")

func SaveTask(task *Task) error {
	db, openErr := bolt.Open(fileName, 0600, &bolt.Options{Timeout: 1 * time.Second})
	if openErr != nil {
		log.Fatal(openErr)
	}
	defer db.Close()
	updateErr := db.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists(bucketName)
		if err != nil {
			return err
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
	db, openErr := bolt.Open(fileName, 0600, &bolt.Options{Timeout: 1 * time.Second})
	if openErr != nil {
		log.Fatal(openErr)
	}
	defer db.Close()

	var task Task
	err := db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(bucketName)
		if bucket == nil {
			return fmt.Errorf("Bucket %q not found!", bucketName)
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
	db, openErr := bolt.Open(fileName, 0600, &bolt.Options{Timeout: 1 * time.Second})
	if openErr != nil {
		log.Fatal(openErr)
	}
	defer db.Close()

	var task Task
	err := db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(bucketName)
		if bucket == nil {
			return fmt.Errorf("Bucket %q not found!", bucketName)
		}
		c := bucket.Cursor()
		for k, v := c.First(); k != nil && bytes.Contains(v, []byte(taskStatus)); k, v = c.Next() {
			fmt.Println(string(k), string(v))
			err := json.Unmarshal(v, &task)
			if err != nil {
				return err
			}
			return nil
		}
		return fmt.Errorf("Did not find any %s task", taskStatus)
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(task)
	return task, nil
}
