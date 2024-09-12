package boltdb

import (
	"errors"
	"fmt"
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
	newTask := &Task{
		Name:      taskName,
		Category:  category,
		StartTime: time.Now().Unix(),
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
	return FormatTaskStatus(task), nil
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
	return FormatTaskStatus(task), nil
}

func ContinuePausedTask(taskName string) error {
	taskObj := Task{
		Name:   taskName,
		Status: TaskStatus(PAUSED),
	}
	tasks, err := FindTasks(taskObj)
	if err != nil {
		return err
	}
	// BUG if len > 1, return to user and ask them to pass the ID instead
	task := tasks[0]
	task.Status = ONGOING
	return SaveTask(&task)
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
	return FormatTaskStatus(task), nil
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
