package boltdb

import (
	"fmt"
	"slices"
	"strings"
	"time"
)

type CategoryDuration struct {
	Name     string
	Duration int64
}

type TaskFreq struct {
	Name string
	Freq int
}

func CalculateTotalTime(tasks []Task) int64 {
	// BUG format time duration to return human friendly string
	var total int64
	for _, task := range tasks {
		total += task.Duration
	}
	return total
}

func FindLongestAndShortestTasks(tasks []Task) (Task, Task) {
	if len(tasks) == 0 {
		return Task{}, Task{}
	}
	longest, shortest := tasks[0], tasks[0]
	for _, task := range tasks {
		if task.Duration > longest.Duration {
			longest = task
		}
		if task.Duration < shortest.Duration {
			shortest = task
		}
	}
	return longest, shortest
}

func FindLongestAndShortestCategories(tasks []Task) (CategoryDuration, CategoryDuration, []CategoryDuration) {
	if len(tasks) == 0 {
		return CategoryDuration{}, CategoryDuration{}, nil
	}
	categories := make(map[string]int64)

	for _, task := range tasks {
		if _, ok := categories[task.Category]; ok {
			categories[task.Category] += task.Duration
		} else {
			categories[task.Category] = task.Duration
		}
	}

	var longest, shortest CategoryDuration
	var result []CategoryDuration
	first := true

	for k, v := range categories {
		catMap := CategoryDuration{Name: k, Duration: v}
		if first {
			longest, shortest = catMap, catMap
			first = false
			continue
		}
		if v > longest.Duration {
			longest = catMap
		}
		if v < shortest.Duration {
			shortest = catMap
		}
		result = append(result, catMap)

		slices.SortStableFunc(result, func(a, b CategoryDuration) int {
			if a.Duration > b.Duration {
				return 1
			} else if a.Duration < b.Duration {
				return -1
			}
			return 0
		})
	}

	return longest, shortest, result
}

func CalculateTaskCompletionStats(tasks []Task) (int, int, int) {
	var completed, ongoing, paused int
	for _, task := range tasks {
		switch task.Status {
		case COMPLETED:
			completed++
		case ONGOING:
			ongoing++
		case PAUSED:
			paused++
		}
	}
	return completed, ongoing, paused
}

func AnalyzeTimeByDay(tasks []Task) map[string]time.Duration {
	result := make(map[string]time.Duration)

	for _, task := range tasks {
		startTime := time.Unix(task.StartTime, 0)
		dateKey := startTime.Format("2006-01-02 15:00")
		duration := time.Duration(task.EndTime-task.StartTime) * time.Second
		result[dateKey] += duration
	}

	return result
}

func AnalyzeTimeByHour(tasks []Task) map[string]time.Duration {
	result := make(map[string]time.Duration)

	for _, task := range tasks {
		startTime := time.Unix(task.StartTime, 0)
		hourKey := startTime.Format("2006-01-02 15:00")
		duration := time.Duration(task.EndTime-task.StartTime) * time.Second
		result[hourKey] += duration
	}

	return result
}

func MostFrequentTaskName(tasks []Task) []TaskFreq {

	length := 3
	if len(tasks) < 3 {
		length = len(tasks)
	}
	freqMap := make(map[string]int)
	for _, task := range tasks[:length] {
		freqMap[task.Name]++
	}

	result := make([]TaskFreq, 0, len(freqMap))
	for k, v := range freqMap {
		result = append(result, TaskFreq{Name: k, Freq: v})
	}

	slices.SortStableFunc(result, func(a, b TaskFreq) int {
		if a.Freq > b.Freq {
			return 1
		} else if a.Freq < b.Freq {
			return -1
		}
		return 0
	})

	return result
}

func FindLongestStreak() string               { return "" }
func FindFirstAndLastTasks() (string, string) { return "", "" }

func FormatMultipleTaskStatus(tasks []Task) string {
	result := ""
	for _, task := range tasks {
		result += FormatTaskStatus(task)
	}
	return result
}

func FormatTaskStatus(task Task) string {
	startTime := time.Unix(task.StartTime, 0).Format("2006-01-02 15:04:05")
	endTime := "N/A"
	if task.EndTime > 0 {
		endTime = time.Unix(task.EndTime, 0).Format("2006-01-02 15:04:05")
	}

	var duration int64
	if strings.Compare(string(task.Status), string(TaskStatus(ONGOING))) == 0 {
		duration = task.Duration + time.Now().Unix() - task.UpdatedAt
	} else if strings.Compare(string(task.Status), string(TaskStatus(PAUSED))) == 0 {
		duration = task.Duration
	} else {
		duration = task.Duration
	}

	return fmt.Sprintf("\nTask Status:\n"+
		"---------------------------------\n"+
		"ID        : %d\n"+
		"Name      : %s\n"+
		"Category  : %s\n"+
		"Status    : %s\n"+
		"Start Time: %s\n"+
		"End Time  : %s\n"+
		"Duration  : %s\n"+
		"---------------------------------\n",
		task.ID,
		task.Name,
		task.Category,
		task.Status,
		startTime,
		endTime,
		FormatDuration(duration),
	)
}

func FormatTasksNamesAndIDs(tasks []Task) string {
	result := ""
	for _, task := range tasks {
		result += fmt.Sprintf("\nTask Status:\n"+
			"---------------------------------\n"+
			"ID        : %d\n"+
			"Name      : %s\n"+
			"Category  : %s\n"+
			"Status    : %s\n"+
			"---------------------------------\n",
			task.ID,
			task.Name,
			task.Category,
			task.Status,
		)
	}
	return result
}

func FormatTopCategories(items []CategoryDuration) string {
	length := 3
	if len(items) < 3 {
		length = len(items)
	}
	result := "\n"

	for _, item := range items[:length] {
		result += fmt.Sprintf("%s - %s", item.Name, FormatDuration(item.Duration))
	}
	result += "\n"
	return result
}

func FormatTopTaskName(items []TaskFreq) string {
	length := 3
	if len(items) < 3 {
		length = len(items)
	}
	result := "\n"

	for _, item := range items[:length] {
		result += fmt.Sprintf("%s - %d", item.Name, item.Freq)
	}
	result += "\n"
	return result
}

func FormatDuration(d int64) string {
	hours := d / 3600
	minutes := (d % 3600) / 60
	seconds := d % 60
	return fmt.Sprintf("%d hours %d minutes %d seconds", hours, minutes, seconds)
}

func FormatTimeByDay(m map[string]time.Duration) string  { return "" }
func FormatTimeByHour(m map[string]time.Duration) string { return "" }
