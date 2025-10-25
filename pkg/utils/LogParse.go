package utils

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"
)

// LogEntry 日志条目结构
type LogEntry struct {
	Type      string `json:"type"`      // 日志类型: info, debug, error, warn 等
	Time      string `json:"time"`      // 时间: HH:MM:SS
	Message   string `json:"message"`   // 日志内容
	RawText   string `json:"rawText"`   // 原始日志文本
	Timestamp string `json:"timestamp"` // 完整时间戳用于排序
}

// 从日志目录读取最近的N条日志
func ReadRecentLogs(logDir string, maxLines int) ([]string, error) {
	// 读取目录下的所有文件
	files, err := os.ReadDir(logDir)
	if err != nil {
		return nil, fmt.Errorf("读取目录失败: %v", err)
	}

	// 收集所有符合日期格式的日志文件
	var logFiles []string
	for _, file := range files {
		if file.IsDir() {
			continue
		}

		filename := file.Name()
		// 检查是否符合 YYYY-MM-DD.log 格式
		if strings.HasSuffix(filename, ".log") {
			datePart := strings.TrimSuffix(filename, ".log")
			// 尝试解析日期
			_, err := time.Parse("2006-01-02", datePart)
			if err == nil {
				logFiles = append(logFiles, filename)
			}
		}
	}

	if len(logFiles) == 0 {
		return nil, fmt.Errorf("未找到符合格式的日志文件")
	}

	// 按文件名（日期）降序排序，最新的在前面
	sort.Slice(logFiles, func(i, j int) bool {
		return logFiles[i] > logFiles[j]
	})

	// 从最新的文件开始读取日志
	var logs []string
	res := maxLines
	for _, filename := range logFiles {
		if len(logs) >= maxLines {
			break
		}

		filePath := filepath.Join(logDir, filename)
		fileLines, err := readFileLines(filePath, res)
		res -= len(fileLines)
		if err != nil {
			fmt.Printf("警告: 读取文件 %s 失败: %v\n", filename, err)
			continue
		}

		// 从文件末尾开始添加日志（最新的日志在文件末尾）
		for i := len(fileLines) - 1; i >= 0; i-- {
			if len(logs) >= maxLines {
				break
			}
			logs = append(logs, fileLines[i])
		}
	}

	// 反转logs数组，使最新的日志在前面
	for i := 0; i < len(logs)/2; i++ {
		logs[i], logs[len(logs)-1-i] = logs[len(logs)-1-i], logs[i]
	}

	return logs, nil
}

func readFileLines(filePath string, readLines int) ([]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.TrimSpace(line) != "" { // 跳过空行
			lines = append(lines, line)
		}
		if len(lines) >= readLines {
			lines = lines[1:]
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return lines, nil
}

// 日志格式的正则表达式: [info/debug] - 19:20:23 ....
var logPattern = regexp.MustCompile(`^\[(\w+)\]\s*-\s*(\d{2}:\d{2}:\d{2})\s+(.*)$`)

// 解析单行日志
func parseLogLine(line string, date string) (*LogEntry, error) {
	matches := logPattern.FindStringSubmatch(line)
	if matches == nil {
		// 如果不匹配格式，返回nil
		return nil, fmt.Errorf("日志格式不匹配")
	}

	return &LogEntry{
		Type:      strings.ToLower(matches[1]),
		Time:      matches[2],
		Message:   strings.TrimSpace(matches[3]),
		RawText:   line,
		Timestamp: fmt.Sprintf("%s %s", date, matches[2]),
	}, nil
}
