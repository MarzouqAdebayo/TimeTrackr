package boltdb

import (
	"errors"
	"fmt"
	"log"
	"strings"
	"time"
)

const (
	NO_ONGOING_ERR_MSG       = "No ongoing time tracking session"
	NO_PAUSED_ERR_MSG        = "No paused time tracking session"
	NO_FILTER_RESULT_ERR_MSG = "No tasks matched the filters\n"
)

var ErrTaskFilterNotFound = errors.New(NO_FILTER_RESULT_ERR_MSG)

type FilterObject struct {
	Name        string
	Category    string
	Status      TaskStatus
	StartDate   int64
	EndDate     int64
	MinDuration int64
	MaxDuration int64
}

func NewFilterObject(filters ...interface{}) *FilterObject {
	filterObj := &FilterObject{}
	statusMap := map[string]TaskStatus{
		"ongoing":   TaskStatus(ONGOING),
		"completed": TaskStatus(COMPLETED),
		"paused":    TaskStatus(PAUSED),
	}

	for i := 0; i < len(filters); i += 2 {
		if i+1 >= len(filters) {
			fmt.Println("Warning: Filter value for", filters[i], "is missing.")
			break
		}
		key := filters[i]
		value := filters[i+1]

		switch key {
		case "name":
			if v, ok := value.(string); ok {
				filterObj.Name = v
			}
		case "category":
			if v, ok := value.(string); ok {
				filterObj.Category = v
			}
		case "status":
			if v, ok := value.(string); ok {
				filterObj.Status = statusMap[v]
			}
		case "startDate":
			if v, ok := value.(int64); ok {
				filterObj.StartDate = v
			}
		case "endDate":
			if v, ok := value.(int64); ok {
				filterObj.EndDate = v
			}
		case "minDuration":
			if v, ok := value.(int64); ok {
				filterObj.MinDuration = v
			}
		case "maxDuration":
			if v, ok := value.(int64); ok {
				filterObj.MaxDuration = v
			}
		default:
			log.Fatalf("Warning: Unknown filter key", key)
			// fmt.Println("Warning: Unknown filter key", key)
		}
	}

	return filterObj
}

func Setup() error {
	return CreateBucket()
}

func StartTask(taskName string, category string) error {
	ongoingErr := OngoingExists()
	if ongoingErr != nil {
		return ongoingErr
	}
	newTask := &Task{
		Name:      taskName,
		Category:  category,
		StartTime: time.Now().Unix(),
		UpdatedAt: time.Now().Unix(),
		Status:    TaskStatus(ONGOING),
	}
	saveErr := SaveTask(newTask)
	if saveErr != nil {
		return saveErr
	}
	return nil
}

func StopCurrentTask(taskID *int) (string, error) {
	var task Task
	if taskID != nil {
		item, err := GetTask(*taskID)
		if err != nil {
			return "", err
		}
		task = item
	} else {
		filterObj := FilterObject{Status: TaskStatus(ONGOING)}
		tasks, err := FilterTasks(filterObj)
		if err != nil {
			return "", err
		}
		if len(tasks) == 0 {
			return "", fmt.Errorf(NO_ONGOING_ERR_MSG)
		}
		if len(tasks) > 1 {
			return "These are currently ongoing tasks, pass the ID of the task you want to stop to the stop command using the --id or -i flag\n" + FormatTasksNamesAndIDs(tasks), nil
		}
		task = tasks[0]
	}
	if strings.Compare(string(task.Status), string(TaskStatus(ONGOING))) == 0 {
		task.Duration += time.Now().Unix() - task.UpdatedAt
	}
	task.UpdatedAt = time.Now().Unix()
	task.EndTime = time.Now().Unix()
	task.Status = TaskStatus(COMPLETED)
	saveErr := UpdateTask(&task)
	if saveErr != nil {
		return "", saveErr
	}
	return FormatTaskStatus(task), nil
}

func PauseCurrentTask(taskID *int) (string, error) {
	var task Task
	if taskID != nil {
		item, err := GetTask(*taskID)
		if err != nil {
			return "", err
		}
		task = item
	} else {

		filterObj := FilterObject{Status: TaskStatus(ONGOING)}
		tasks, err := FilterTasks(filterObj)
		if err != nil {
			return "", err
		}
		if len(tasks) == 0 {
			return "", fmt.Errorf(NO_ONGOING_ERR_MSG)
		}
		if len(tasks) > 1 {
			return "These are currently ongoing tasks, pass the ID of the task you want to pause to the pause command using the --id or -i flag\n" + FormatTasksNamesAndIDs(tasks), nil
		}
		task = tasks[0]
	}
	task.Status = TaskStatus(PAUSED)
	task.Duration += time.Now().Unix() - task.UpdatedAt
	task.UpdatedAt = time.Now().Unix()
	saveErr := UpdateTask(&task)
	if saveErr != nil {
		return "", saveErr
	}
	return FormatTaskStatus(task), nil
}

func ContinuePausedTask(taskName *string, taskID *int) (string, error) {
	var task Task
	if taskID != nil {
		item, err := GetTask(*taskID)
		if err != nil {
			return "", err
		}
		task = item
	} else {
		filterObj := FilterObject{
			Name:   *taskName,
			Status: TaskStatus(PAUSED),
		}
		tasks, err := FilterTasks(filterObj)
		if err != nil {
			return "", err
		}
		if len(tasks) == 0 {
			return "", fmt.Errorf("%s with the name '%s'\n", NO_PAUSED_ERR_MSG, *taskName)
		}
		if len(tasks) > 1 {
			return "These are currently paused tasks, pass the ID of the task you want to continue to the continue command using the --id or -i flag\n" + FormatTasksNamesAndIDs(tasks), nil
		}
		task = tasks[0]
	}
	task.UpdatedAt = time.Now().Unix()
	task.Status = ONGOING
	return "", UpdateTask(&task)
}

func Filter(filter *FilterObject) ([][]string, error) {
	tasks, err := FilterTasks(*filter)
	if err != nil {
		return nil, err
	}
	if len(tasks) == 0 {
		return nil, fmt.Errorf(NO_FILTER_RESULT_ERR_MSG)
	}
	return ParseIntoRows(tasks), nil
}

func Status(id *int) (string, error) {
	if id != nil {
		task, err := GetTask(*id)
		if err != nil {
			return "", err
		}
		return FormatTaskStatus(task), nil
	} else {
		filter := FilterObject{
			Status: TaskStatus(ONGOING),
		}
		tasks, err := FilterTasks(filter)
		if err != nil {
			return "", err
		}
		if len(tasks) == 0 {
			return "", fmt.Errorf("No current time tracking session\n")
		}
		return FormatTaskStatus(tasks[0]), nil
	}
}

func GenerateReport(startDate, endDate int64) (string, error) {
	filterObj := FilterObject{StartDate: startDate, EndDate: endDate}
	tasks, err := FilterTasks(filterObj)
	if err != nil {
		return "", err
	}
	// TODO use goroutines for all these function calls :)
	totalTasks := len(tasks)
	totalTimeSpent := CalculateTotalTime(tasks)
	topTasks := TopTasksByDuration(tasks)
	topCategories := TopCategoriesByDuration(tasks)
	completedTasks, ongoingTasks, pausedTasks := CalculateTaskCompletionStats(tasks)
	mostFrequentTasks := MostFrequentTaskName(tasks)
	longestStreak := FindLongestStreak()

	report := fmt.Sprintf(`
TimeTrackr Report: Detailed Analysis

Overview:
    Report Period: %s - %s
    Total Tasks Tracked: %d
    Total Time Spent: %s

Top Time Consuming Tasks:%s
Top Time Consuming Categories:%s
Task Completion Statistics:
    Number of Completed Tasks: %d
    Number of Ongoing Tasks: %d
    Number of Paused Tasks: %d
    Average Time Spent per Task: %s

Task Status Summary:
    Percentage of Completed Tasks: %.2f%%
    Percentage of Ongoing Tasks: %.2f%%
    Percentage of Paused Tasks: %.2f%%

Miscellaneous Insights:
    Most Frequent Task Name:%s
    Longest Streak of Task Tracking: "%s" - %d consecutive days
`,
		time.Unix(startDate, 0).Format("Jan 2, 2006 3:04 PM"), time.Unix(endDate, 0).Format("Jan 2, 2006 3:04PM"),
		totalTasks, FormatDuration(totalTimeSpent),
		topTasks, topCategories,
		completedTasks, ongoingTasks, pausedTasks, FormatDuration(totalTimeSpent/int64(totalTasks)),
		float64(completedTasks)/float64(totalTasks)*100, float64(ongoingTasks)/float64(totalTasks)*100, float64(pausedTasks)/float64(totalTasks)*100,
		mostFrequentTasks,
		longestStreak, 2)
	return report, nil
}
