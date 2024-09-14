package boltdb

import (
	"fmt"
	"strings"
	"time"
)

const (
	NO_ONGOING_ERR_MSG = "No ongoing time tracking session"
	NO_PAUSED_ERR_MSG  = "No paused time tracking session"
)

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
		filterObj := Task{Status: TaskStatus(ONGOING)}
		tasks, err := FindTasks(filterObj)
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

		filterObj := Task{Status: TaskStatus(ONGOING)}
		tasks, err := FindTasks(filterObj)
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
		filterObj := Task{
			Name:   *taskName,
			Status: TaskStatus(PAUSED),
		}
		tasks, err := FindTasks(filterObj)
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

func Status(status string) (string, error) {
	statusMap := map[string]TaskStatus{
		"ongoing":   TaskStatus(ONGOING),
		"completed": TaskStatus(COMPLETED),
		"paused":    TaskStatus(PAUSED),
	}
	filterObj := Task{
		Status: statusMap[status],
	}
	tasks, err := FindTasks(filterObj)
	if err != nil {
		return "", err
	}
	if len(tasks) == 0 {
		return fmt.Sprintf("There are no %s tasks\n", status), nil
	}
	return FormatMultipleTaskStatus(tasks), nil
}

func GenerateReport(startDate, endDate int64) (string, error) {
	tasks, err := FilterTasks(&startDate, &endDate, nil, nil)
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
