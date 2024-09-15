package boltdb

import (
	"fmt"
	"slices"
	"strconv"
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

func TopTasksByDuration(tasks []Task) string {
	if len(tasks) == 0 {
		return ""
	}
	slices.SortStableFunc(tasks, func(a, b Task) int {
		if a.Duration < b.Duration {
			return 1
		} else if a.Duration > b.Duration {
			return -1
		}
		return 0
	})
	result := "\n"
	length := len(tasks)
	for index, task := range tasks {
		if index == 3 {
			break
		}
		result += fmt.Sprintf("    %d. %s (%s): %s\n", index+1, task.Name, task.Category, FormatDuration(int64(task.Duration)))
	}
	if length > 4 {
		result += "    ...\n"
		result += fmt.Sprintf("    %d. %s (%s): %s", len(tasks), tasks[length-1].Name, tasks[length-1].Category, FormatDuration(int64(tasks[length-1].Duration)))
	}
	return result
}

func TopCategoriesByDuration(tasks []Task) string {
	if len(tasks) == 0 {
		return ""
	}
	categories := make(map[string]int64)

	for _, task := range tasks {
		if _, ok := categories[task.Category]; ok {
			categories[task.Category] += task.Duration
		} else {
			categories[task.Category] = task.Duration
		}
	}

	var result []CategoryDuration
	index := 0
	for k, v := range categories {
		if index > 2 {
			break
		}
		index++
		catMap := CategoryDuration{Name: k, Duration: v}
		result = append(result, catMap)
	}
	slices.SortStableFunc(result, func(a, b CategoryDuration) int {
		if a.Duration < b.Duration {
			return 1
		} else if a.Duration > b.Duration {
			return -1
		}
		return 0
	})
	report := "\n"
	length := len(result)
	for index, category := range result {
		if index == 3 {
			break
		}
		report += fmt.Sprintf("    %d. %s: %s\n", index+1, category.Name, FormatDuration(category.Duration))
	}
	if length > 4 {
		report += "    ...\n"
		report += fmt.Sprintf("    %d. %s: %s\n", index+1, result[length-1].Name, FormatDuration(result[length-1].Duration))
	}
	return report
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

func MostFrequentTaskName(tasks []Task) string {
	freqMap := make(map[string]int)
	for _, task := range tasks {
		freqMap[task.Name]++
	}

	result := make([]TaskFreq, 0, len(freqMap))
	for k, v := range freqMap {
		result = append(result, TaskFreq{Name: k, Freq: v})
	}

	slices.SortStableFunc(result, func(a, b TaskFreq) int {
		if a.Freq < b.Freq {
			return 1
		} else if a.Freq > b.Freq {
			return -1
		}
		return 0
	})

	report := "\n"
	length := len(result)
	for index, category := range result {
		if index == 3 {
			break
		}
		report += fmt.Sprintf("    %d. %s: Tracked %d time(s)\n", index+1, category.Name, category.Freq)
	}
	if length > 4 {
		report += "    ...\n"
		report += fmt.Sprintf("    %d. %s: Tracked %d time(s)\n", length, result[length-1].Name, result[length-1].Freq)
	}
	return report
}

func FindLongestStreak() string { return "" }

func ParseIntoRows(tasks []Task) [][]string {
	rows := make([][]string, 0, len(tasks))
	for _, v := range tasks {
		id := fmt.Sprintf("%d", v.ID)
		name := v.Name
		status := string(v.Status)
		category := v.Category
		duration := FormatDuration(v.Duration)
		// startTime := FormatDate(v.StartTime)
		// updatedAt := FormatDate(v.UpdatedAt)
		// endTime := FormatDate(v.EndTime)

		row := []string{
			strings.Trim(id, " \n\t"),
			strings.Trim(name, " \n\t"),
			strings.Trim(status, " \n\t"),
			strings.Trim(category, " \n\t"),
			strings.Trim(duration, " \n\t")}
		// strings.Trim(startTime, " \n\t"),
		// strings.Trim(updatedAt, " \n\t"),
		// strings.Trim(endTime, "\n\t")}
		rows = append(rows, row)
	}
	return rows
}

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

func FormatDate(d int64) string {
	return time.Unix(d, 0).Format("01-02-2006 03:04 PM")
}

func FormatDuration(d int64) string {
	hours := d / 3600
	minutes := (d % 3600) / 60
	seconds := d % 60

	return fmt.Sprintf("%d:%d:%d\n", hours, minutes, seconds)
}

func ParseDateCommand(startDateStr, endDateStr string) (int64, int64, error) {
	loc := time.Local

	var startTimestamp int64
	if startDateStr == "" {
		startTimestamp = 0
	} else {
		startDate, err := parseDateWithTimezone(startDateStr, loc)
		if err != nil {
			return 0, 0, fmt.Errorf("Invalid start date format. Please use YYYY-MM-DD or YYYY-MM-DD HH:MM:SS")
		}
		startTimestamp = startDate.Unix()
	}

	var endTimestamp int64
	if endDateStr == "" {
		endTimestamp = time.Now().In(loc).Unix()
	} else {
		endDate, err := parseDateWithTimezone(endDateStr, loc)
		if err != nil {
			return 0, 0, fmt.Errorf("Invalid end date format. Please use 'YYYY-MM-DD' or 'DD-MM-YYYY' or 'YYYY-MM-DD' or 'DD-MM-YYYY HH:MM' or 'HH:MM AM/PM'")
		}
		endTimestamp = endDate.Unix()
	}

	if startTimestamp >= endTimestamp {
		return 0, 0, fmt.Errorf("Error: Start date must be before end date")
	}

	if startTimestamp >= endTimestamp {
		return 0, 0, fmt.Errorf("Error: Start date must be before end date")
	}
	return startTimestamp, endTimestamp, nil
}

func parseDateWithTimezone(dateStr string, loc *time.Location) (time.Time, error) {
	layouts := []string{
		"02-01-2006",
		"2006-01-02",
		"02-01-2006 15:04",
		"2006-01-02 15:04",
		"02-01-2006 03:04 PM",
		"2006-01-02 03:04 PM",
	}
	var parsedDate time.Time
	var err error
	for _, layout := range layouts {
		parsedDate, err = time.ParseInLocation(layout, dateStr, loc)
		if err == nil {
			return parsedDate, nil
		}
	}
	return time.Time{}, fmt.Errorf("Invalid date format")
}

func DurationParser(durationStr string) (int64, error) {
	var duration int64

	i := 0
	cursor := 0
	for i < len(durationStr) {
		char := durationStr[i]
		switch char {
		case 'h', 'H':
			hours, err := strconv.ParseInt(durationStr[cursor:i], 10, 64)
			if err != nil {
				return 0, fmt.Errorf("Invalid duration format. Use '2h3m30s' or '2H3M30S'")
			}
			duration += hours * 3600
			cursor = i + 1
		case 'm', 'M':
			minutes, err := strconv.ParseInt(durationStr[cursor:i], 10, 64)
			if err != nil {
				return 0, fmt.Errorf("Invalid duration format. Use '2h3m30s' or '2H3M30S'")
			}
			duration += minutes * 60
			cursor = i + 1
		case 's', 'S':
			seconds, err := strconv.ParseInt(durationStr[cursor:i], 10, 64)
			if err != nil {
				return 0, fmt.Errorf("Invalid duration format. Use '2h3m30s' or '2H3M30S'")
			}
			duration += seconds
			cursor = i + 1
		}
		i++
	}
	return duration, nil
}
