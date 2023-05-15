package log

import (
	"fmt"
	"github.com/H-kang-better/msgo/internal/msstrings"
	"io"
	"log"
	"os"
	"path"
	"strings"
	"time"
)

const (
	LevelDebug LoggerLevel = iota
	LevelInfo
	LevelError
)

const (
	greenBg   = "\033[97;42m"
	whiteBg   = "\033[90;47m"
	yellowBg  = "\033[90;43m"
	redBg     = "\033[97;41m"
	blueBg    = "\033[97;44m"
	magentaBg = "\033[97;45m"
	cyanBg    = "\033[97;46m"
	green     = "\033[32m"
	white     = "\033[37m"
	yellow    = "\033[33m"
	red       = "\033[31m"
	blue      = "\033[34m"
	magenta   = "\033[35m"
	cyan      = "\033[36m"
	reset     = "\033[0m"
)

type Fields map[string]any

type LoggerLevel int

type LoggerFormatter struct {
	IsColor      bool // 日志是否要颜色
	Level        LoggerLevel
	LoggerFields Fields
}

type LoggingFormatter interface {
	Formatter(param *LoggingFormatterParam) string
}

type LoggingFormatterParam struct {
	IsColor      bool
	Level        LoggerLevel
	Msg          any
	LoggerFields Fields
}

type Logger struct {
	Formatter    LoggingFormatter
	Outs         []*LoggerWriter // 放置多个输出方式
	Level        LoggerLevel
	LoggerFields Fields
	logPath      string
	LogFileSize  int64
}
type LoggerWriter struct {
	Level LoggerLevel
	Out   io.Writer
}

func New() *Logger {
	return &Logger{}
}

func Default() *Logger {
	logger := New()
	w := &LoggerWriter{
		Level: LevelDebug,
		Out:   os.Stdout,
	}
	logger.Outs = append(logger.Outs, w)
	logger.Level = LevelDebug
	logger.Formatter = &TextFormatter{}
	return logger
}

func (l *Logger) Info(msg any) {
	l.Print(LevelInfo, msg)
}

func (l *Logger) Debug(msg any) {
	l.Print(LevelDebug, msg)
}

func (l *Logger) Error(msg any) {
	l.Print(LevelError, msg)
}

func (l *Logger) Print(level LoggerLevel, msg any) {
	if l.Level > level {
		//级别不满足 不打印日志
		return
	}
	param := &LoggingFormatterParam{
		Level:        level,
		Msg:          msg,
		LoggerFields: l.LoggerFields,
	}
	formatter := l.Formatter.Formatter(param)
	for _, out := range l.Outs {
		if out.Out == os.Stdout {
			param.IsColor = true
			formatter = l.Formatter.Formatter(param)
			fmt.Fprintln(out.Out, formatter)
		}
		if out.Level == -1 || out.Level == level {
			fmt.Fprintln(out.Out, formatter)
			l.CheckFileSize(out)
		}
	}
}

func (l *Logger) WithFields(fields Fields) *Logger {
	return &Logger{
		Formatter:    l.Formatter,
		Outs:         l.Outs,
		Level:        l.Level,
		LoggerFields: fields,
	}
}

func (l *Logger) SetLogPath(logPath string) {
	l.logPath = logPath
	l.Outs = append(l.Outs, &LoggerWriter{
		Level: -1,
		Out:   FileWriter(path.Join(logPath, "all.log")),
	})
	l.Outs = append(l.Outs, &LoggerWriter{
		Level: LevelDebug,
		Out:   FileWriter(path.Join(logPath, "debug.log")),
	})
	l.Outs = append(l.Outs, &LoggerWriter{
		Level: LevelInfo,
		Out:   FileWriter(path.Join(logPath, "info.log")),
	})
	l.Outs = append(l.Outs, &LoggerWriter{
		Level: LevelError,
		Out:   FileWriter(path.Join(logPath, "error.log")),
	})
}

func (l *Logger) CheckFileSize(out *LoggerWriter) {
	osFile := out.Out.(*os.File)
	if osFile != nil {
		stat, err := osFile.Stat()
		if err != nil {
			log.Println("logger checkFileSize error info :", err)
			return
		}
		size := stat.Size()
		//这里要检查大小，如果满足条件 就重新创建文件，并且更换logger中的输出
		if l.LogFileSize <= 0 {
			//默认100M
			l.LogFileSize = 100 << 20
		}
		if size >= l.LogFileSize {
			_, fileName := path.Split(osFile.Name())
			name := fileName[0:strings.Index(fileName, ".")]
			w := FileWriter(path.Join(l.logPath, msstrings.JoinStrings(name, ".", time.Now().UnixMilli(), ".log")))
			if err != nil {
				log.Println("logger checkFileSize error info :", err)
				return
			}
			out.Out = w
		}
	}

}

func (f *LoggerFormatter) formatter(msg any) string {
	now := time.Now()
	if f.IsColor {
		//要带颜色  error的颜色 为红色 info为绿色 debug为蓝色
		levelColor := f.LevelColor()
		msgColor := f.MsgColor()
		return fmt.Sprintf("%s [msgo] %s %s%v%s | level= %s %s %s | msg=%s %#v %s | fields=%s\n",
			yellow, reset, blue, now.Format("2006/01/02 - 15:04:05"), reset,
			levelColor, f.Level.Level(), reset, msgColor, msg, reset, f.LoggerFields,
		)
	}
	return fmt.Sprintf("[msgo] %v | level=%s | msg=%#v | fields=%s \n",
		now.Format("2006/01/02 - 15:04:05"),
		f.Level.Level(), msg, f.LoggerFields,
	)
}

func (level LoggerLevel) Level() string {
	switch level {
	case LevelDebug:
		return "DEBUG"
	case LevelInfo:
		return "INFO"
	case LevelError:
		return "ERROR"
	default:
		return ""
	}
}

func (f *LoggerFormatter) MsgColor() string {
	switch f.Level {
	case LevelDebug:
		return ""
	case LevelInfo:
		return ""
	case LevelError:
		return red
	default:
		return cyan
	}
}

func (f *LoggerFormatter) LevelColor() interface{} {
	switch f.Level {
	case LevelDebug:
		return blue
	case LevelInfo:
		return green
	case LevelError:
		return red
	default:
		return cyan
	}
}

// FileWriter 日志写入文件
func FileWriter(name string) io.Writer {
	w, _ := os.OpenFile(name, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	return w
}
