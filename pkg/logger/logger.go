package logger

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"runtime"
	"time"
)

type Level int8

type Fields map[string]interface{}

const (
	LevelDebug Level = iota
	LevelInfo
	LevelWarn
	LevelError
	LevelFatal
	LevelPanic
)

func (l Level) String() string {
	switch l {
	case LevelDebug:
		return "debug"
	case LevelInfo:
		return "info"
	case LevelWarn:
		return "warn"
	case LevelError:
		return "error"
	case LevelFatal:
		return "fatal"
	case LevelPanic:
		return "panic"
	}
	return ""
}

type Logger struct {
	newLogger *log.Logger
	ctx context.Context
	level Level
	fields Fields
	callers []string
}

func NewLogger(w io.Writer, prefix string, flag int) *Logger {
	//参数w设置日志信息写入的目的地。参数prefix会添加到生成的每一条日志前面。参数flag定义日志的属性（时间、文件等等）
	l := log.New(w, prefix, flag)
	return &Logger{newLogger: l}
}

func (l *Logger)clone() *Logger {
	nl := *l
	return &nl
}

//设置日志等级
func (l *Logger) WithLevel(lvl Level) *Logger {
	ll := l.clone()
	ll.level = lvl
	return ll
}

//设置日志公共字段
func (l *Logger) WithFields(f Fields) *Logger {
	ll := l.clone()
	if ll.fields == nil {
		ll.fields = make(Fields)
	}
	for k, v := range f {
		ll.fields[k] = v
	}
	return ll
}


//设置日志上下文属性
func (l *Logger) WithContext(ctx context.Context) *Logger {
	ll := l.clone()
	ll.ctx = ctx
	return ll
}

//设置当前某一层调用栈的信息(程序计数器、文件信息和行号)
//实参skip为上溯的栈帧数，0表示Caller的调用者（Caller所在的调用栈）
func (l *Logger) WithCaller(skip int) *Logger {
	ll := l.clone()
	//函数的返回值为调用栈标识符、文件名、该调用在文件中的行号。如果无法获得信息，ok会被设为false。
	pc, file, line, ok := runtime.Caller(skip)
	if ok {
		//FuncForPC返回一个表示调用栈标识符pc对应的调用栈的*Func；如果该调用栈标识符没有对应的调用栈，函数会返回nil。每一个调用栈必然是对某个函数的调用。
		f := runtime.FuncForPC(pc)
		//f.Name()返回该调用栈所调用的函数的名字
		ll.callers = []string{fmt.Sprintf("%s: %d %s", file, line, f.Name())}
	}

	return ll
}

//设置当前的整个调用栈信息
func (l *Logger) WithCallersFrames() *Logger {
	maxCallerDepth := 25
	minCallerDepth := 1
	callers := []string{}
	//uintptr是一个整数类型，它的大小足以容纳任何指针的位模式
	pcs := make([]uintptr, maxCallerDepth)
	//函数把当前go程调用栈上的调用栈标识符填入切片pc中，返回写入到pc中的项数。
	//实参skip为开始在pc中记录之前所要跳过的栈帧数，0表示Callers自身的调用栈，1表示Callers所在的调用栈。
	depth := runtime.Callers(minCallerDepth, pcs)
	//CallersFrames获取调用者返回的PC值的一部分，并准备返回函数/文件/行信息。
	//Frames是Frames为每个调用帧返回的信息
	frames := runtime.CallersFrames(pcs[:depth])
	//Next返回下一个调用方的帧信息。
	//如果more为false，则不再有调用方（帧值有效）。
	for frame, more := frames.Next(); more; frame, more = frames.Next() {
		callers = append(callers, fmt.Sprintf("%s: %d %s",frame.File, frame.Line, frame.Function))
		if !more {
			break
		}
	}

	ll := l.clone()
	ll.callers = callers
	return ll
}

//日志格式化
func (l *Logger) JSONFormat(message string) map[string]interface{} {
	data := make(Fields, len(l.fields)+4)
	data["level"] = l.level.String()
	data["time"] = time.Now().Local().UnixNano()
	data["message"] = message
	data["callers"] = l.callers
	if len(l.fields) > 0 {
		for k, v := range l.fields {
			if _, ok := data[k]; !ok {
				data[k] = v
			}
		}
	}

	return data
}

//日志输出
func (l *Logger) Output(message string) {
	body, _ := json.Marshal(l.JSONFormat(message))
	content := string(body)
	switch l.level {
	case LevelDebug:
		l.newLogger.Print(content)
	case LevelInfo:
		l.newLogger.Print(content)
	case LevelWarn:
		l.newLogger.Print(content)
	case LevelError:
		l.newLogger.Print(content)
	case LevelFatal:
		l.newLogger.Print(content)
	case LevelPanic:
		l.newLogger.Print(content)
	}
}

//日志分级输出,以及格式化输出
func (l *Logger) Debug(v ...interface{}) {
	l.WithLevel(LevelDebug).Output(fmt.Sprint(v...))
}

func (l *Logger) Debugf(format string, v ...interface{}) {
	l.WithLevel(LevelDebug).Output(fmt.Sprintf(format, v...))
}

func (l *Logger) Info(v ...interface{}) {
	l.WithLevel(LevelInfo).Output(fmt.Sprint(v...))
}

func (l *Logger) Infof(format string, v ...interface{}) {
	l.WithLevel(LevelInfo).Output(fmt.Sprintf(format, v...))
}

func (l *Logger) Fatal(v ...interface{}) {
	l.WithLevel(LevelFatal).Output(fmt.Sprint(v...))
}

func (l *Logger) Fatalf(format string, v ...interface{}) {
	l.WithLevel(LevelFatal).Output(fmt.Sprintf(format, v...))
}

func (l *Logger) Warn(v ...interface{}) {
	l.WithLevel(LevelWarn).Output(fmt.Sprint(v...))
}

func (l *Logger) Warnf(format string, v ...interface{}) {
	l.WithLevel(LevelWarn).Output(fmt.Sprintf(format, v...))
}

func (l *Logger) Error(v ...interface{}) {
	l.WithLevel(LevelError).Output(fmt.Sprint(v...))
}

func (l *Logger) Errorf(format string, v ...interface{}) {
	l.WithLevel(LevelError).Output(fmt.Sprintf(format, v...))
}

func (l *Logger) Panic(v ...interface{}) {
	l.WithLevel(LevelPanic).Output(fmt.Sprint(v...))
}

func (l *Logger) Panicf(format string, v ...interface{}) {
	l.WithLevel(LevelPanic).Output(fmt.Sprintf(format, v...))
}