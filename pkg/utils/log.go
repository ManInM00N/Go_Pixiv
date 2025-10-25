package utils

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"
)

type contextKey string

const traceIDKey contextKey = "traceID"

// Logger 封装了标准logger并添加调用链追踪
type Logger struct {
	logger *log.Logger
	level  string
}

var (
	DebugLog   *log.Logger
	InfoLog    *log.Logger
	ErrorLog   *Logger
	WarnLog    *Logger
	logf       *os.File
	mu         sync.Mutex
	LogPositon = "errorlog"
)

func Log_init() {
	T := time.Now()
	_, err := os.Stat("errorlog")
	if err != nil {
		os.Mkdir("errorlog", 0777)
		os.Chmod("errorlog", 0777)
	}

	logfile := fmt.Sprintf("errorlog/%04d-%02d-%02d.log", T.Year(), T.Month(), T.Day())
	//log.SetFlags(log.Ltime)
	logf, _ = os.OpenFile(logfile, os.O_CREATE|os.O_APPEND|os.O_RDWR, os.ModePerm)
	log.SetOutput(logf)
	InfoLog = log.New(logf, "[Info] - ", log.Ltime)
	DebugLog = log.New(logf, "[Debug] - ", log.Ltime)
	ErrorLog = &Logger{logger: log.New(logf, "", 0), level: "ERROR"}
	WarnLog = &Logger{logger: log.New(logf, "", 0), level: "WARN"}
}

// getGoroutineID 获取当前goroutine ID
func getGoroutineID() uint64 {
	b := make([]byte, 64)
	b = b[:runtime.Stack(b, false)]
	b = bytes.TrimPrefix(b, []byte("goroutine "))
	b = b[:bytes.IndexByte(b, ' ')]
	n, _ := strconv.ParseUint(string(b), 10, 64)
	return n
}

// getCaller 获取调用者信息
func getCaller(skip int) string {
	pc, file, line, ok := runtime.Caller(skip)
	if !ok {
		return "unknown:0"
	}

	// 获取函数名
	fn := runtime.FuncForPC(pc)
	funcName := "unknown"
	if fn != nil {
		funcName = fn.Name()
		// 只保留包名和函数名
		if idx := strings.LastIndex(funcName, "/"); idx != -1 {
			funcName = funcName[idx+1:]
		}
	}

	// 只保留相对路径
	file = filepath.Base(file)

	return fmt.Sprintf("%s:%d:%s", file, line, funcName)
}

// formatLog 格式化日志消息
func (l *Logger) formatLog(traceID string, format string, v ...interface{}) string {
	now := time.Now().Format("15:04:05.000")
	caller := getCaller(3) // skip: formatLog -> Log方法 -> 实际调用者
	gid := getGoroutineID()

	msg := fmt.Sprintf(format, v...)

	if traceID != "" {
		return fmt.Sprintf("[%s] [%s] [G:%d] [TraceID:%s] [%s] %s",
			l.level, now, gid, traceID, caller, msg)
	}
	return fmt.Sprintf("[%s] [%s] [G:%d] [%s] %s",
		l.level, now, gid, caller, msg)
}

// Printf 打印日志（不带TraceID）
func (l *Logger) Printf(format string, v ...interface{}) {
	mu.Lock()
	defer mu.Unlock()
	l.logger.Println(l.formatLog("", format, v...))
}

// PrintfWithTrace 打印带TraceID的日志
func (l *Logger) PrintfWithTrace(traceID string, format string, v ...interface{}) {
	mu.Lock()
	defer mu.Unlock()
	l.logger.Println(l.formatLog(traceID, format, v...))
}

// PrintfWithContext 从context中提取TraceID并打印日志
func (l *Logger) PrintfWithContext(ctx context.Context, format string, v ...interface{}) {
	traceID := GetTraceID(ctx)
	l.PrintfWithTrace(traceID, format, v...)
}

// Context相关辅助函数

// WithTraceID 在context中设置TraceID
func WithTraceID(ctx context.Context, traceID string) context.Context {
	return context.WithValue(ctx, traceIDKey, traceID)
}

// GetTraceID 从context中获取TraceID
func GetTraceID(ctx context.Context) string {
	if ctx == nil {
		return ""
	}
	if traceID, ok := ctx.Value(traceIDKey).(string); ok {
		return traceID
	}
	return ""
}

// GenerateTraceID 生成一个简单的TraceID（可以替换为更复杂的实现）
func GenerateTraceID() string {
	return fmt.Sprintf("%d-%d", time.Now().UnixNano(), getGoroutineID())
}
