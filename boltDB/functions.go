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
	longestTask, shortestTask := FindLongestAndShortestTasks(tasks)
	longestCategory, shortestCategory, topCategories := FindLongestAndShortestCategories(tasks)
	completedTasks, ongoingTasks, pausedTasks := CalculateTaskCompletionStats(tasks)
	timeByDay := AnalyzeTimeByDay(tasks)
	timeByHour := AnalyzeTimeByHour(tasks)
	mostFrequentTask := FormatTopTaskName(MostFrequentTaskName(tasks))
	longestStreak := FindLongestStreak()
	firstTask, lastTask := FindFirstAndLastTasks()
	// Generate the report as a formatted string
	report := fmt.Sprintf(`
TimeTrackr Report: Detailed Analysis

Overview:
  Report Period: %s - %s
  Total Tasks Tracked: %d
  Total Time Spent: %s

Top Tasks:
  Task with Longest Duration: "%s" - %s (Category: %s)
  Task with Shortest Duration: "%s" - %s (Category: %s)

Category Analysis:
  Category with Longest Total Duration: "%s" - %s
  Category with Shortest Total Duration: "%s" - %s
  Top 3 Most Time-Consuming Categories: %s

Task Completion Statistics:
  Number of Completed Tasks: %d
  Number of Ongoing Tasks: %d
  Number of Paused Tasks: %d
  Average Time Spent per Task: %s

Time Distribution by Day:
  %s

Time Distribution by Hour:
  %s

Task Status Summary:
  Percentage of Completed Tasks: %.2f%%
  Percentage of Ongoing Tasks: %.2f%%
  Percentage of Paused Tasks: %.2f%%

Miscellaneous Insights:
  Most Frequent Task Name: "%s" - Tracked %d times
  Longest Streak of Task Tracking: "%s" - %d consecutive days
  First Task of the Report Period: "%s" - Started at %s
  Last Task of the Report Period: "%s" - Ended at %s
`,
		time.Unix(startDate, 0).Format("Jan 2, 2006 3:04 PM"), time.Unix(endDate, 0).Format("Jan 2, 2006 3:04PM"),
		totalTasks, FormatDuration(totalTimeSpent),
		longestTask.Name, FormatDuration(longestTask.Duration), longestTask.Category,
		shortestTask.Name, FormatDuration(shortestTask.Duration), shortestTask.Category,
		longestCategory.Name, FormatDuration(longestCategory.Duration),
		shortestCategory.Name, FormatDuration(shortestCategory.Duration),
		FormatTopCategories(topCategories),
		completedTasks, ongoingTasks, pausedTasks, FormatDuration(totalTimeSpent/int64(totalTasks)),
		FormatTimeByDay(timeByDay), FormatTimeByHour(timeByHour),
		float64(completedTasks)/float64(totalTasks)*100, float64(ongoingTasks)/float64(totalTasks)*100, float64(pausedTasks)/float64(totalTasks)*100,
		mostFrequentTask, 1,
		longestStreak, 2,
		firstTask, 3,
		lastTask, 4)
	return report, nil
}
