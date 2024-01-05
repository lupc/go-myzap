package myzap

import (
	"fmt"
	"io"
	"os"

	"github.com/lupc/go-rollingwriter"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type MyZapConfig struct {
	FileLevel          zapcore.LevelEnabler //文件日记级别
	ConsoleLevel       zapcore.LevelEnabler //控制台日记级别
	LogDirFormat       string               //日志目录格式 (年月日yyyyMMdd，时分秒HHmmss，使用{}包含) 如：./logs/{yyyy-MM}/
	LogNameFormat      string               //日志文件格式 (年月日yyyyMMdd，时分秒HHmmss，使用{}包含) 如：test_{HH}
	TimeFormat         string               //时间显示格式
	LogExt             string               //日志文件后缀 如：log
	IsLogToFile        bool                 //是否输出到文件
	IsLogToConsole     bool                 //是否输出到控制台
	StackLevel         zapcore.LevelEnabler //调用栈输出级别
	IsLogCaller        bool                 //是否输出调用位置
	RollingTimePattern string               //滚动时间匹配 如：0 0 0 * * * 表示每天00:00滚动日志
	RollingSize        string               //滚动大小，数字+单位(K,M,G,KB,MB,GB)，如“50MB”或者“50M”
	MaxAge             int                  //日志保留天数，0一直保留
	ClearTimePattern   string               //执行清理任务时间，如：0 0 2 * * * 表示每天02:00执行清理任务
}

// 构建一个logger
func (c *MyZapConfig) BuildLogger() (logger *zap.Logger) {
	if c == nil {
		return
	}
	var zc zapcore.Core = nil
	var encCfg = zap.NewDevelopmentEncoderConfig()
	if c.TimeFormat != "" {
		encCfg.EncodeTime = zapcore.TimeEncoderOfLayout(c.TimeFormat)
	}
	var encoder = zapcore.NewConsoleEncoder(encCfg)

	if c.IsLogToFile {
		fileWriter := c.getFileWriter()
		var fc = zapcore.NewCore(encoder, zapcore.AddSync(fileWriter), c.FileLevel)
		zc = fc
	}
	if c.IsLogToConsole {
		var cc = zapcore.NewCore(encoder, os.Stdout, c.ConsoleLevel)
		zc = zapcore.NewTee(zc, cc)
	}

	logger = zap.New(zc, zap.AddStacktrace(zap.ErrorLevel), zap.AddCaller())
	return
}

// 获取滚动文件writer
func (c *MyZapConfig) getFileWriter() (writer io.Writer) {

	var cfg = rollingwriter.Config{
		LogPath:                c.LogDirFormat,
		FileName:               c.LogNameFormat,
		FileExtension:          c.LogExt,
		MaxRemain:              0,                    // disable auto delete
		RollingPolicy:          2,                    // TimeRotate by default
		RollingTimePattern:     c.RollingTimePattern, // Rolling at 00:00 AM everyday
		RollingVolumeSize:      c.RollingSize,
		WriterMode:             "async",
		BufferWriterThershould: 64,
		Compress:               false,
		MaxAge:                 c.MaxAge,
		ClearTimePattern:       c.ClearTimePattern,
	}
	var w, err = rollingwriter.NewWriterFromConfig(&cfg)

	if err != nil {
		panic(err)
	}
	writer = w
	return
}

// 根据自定义名称新建一个默认日志配置
func NewConfigByName(logName string) (cfg *MyZapConfig) {
	cfg = &MyZapConfig{
		FileLevel:          zapcore.DebugLevel,
		ConsoleLevel:       zapcore.InfoLevel,
		LogDirFormat:       "./logs/{yyyy-MM}/",
		LogNameFormat:      fmt.Sprintf("%v_{yyyy-MM-dd}", logName),
		TimeFormat:         "2006-01-02 15:04:05.0000",
		LogExt:             "log",
		IsLogToFile:        true,
		IsLogToConsole:     true,
		StackLevel:         zapcore.ErrorLevel,
		IsLogCaller:        false,
		RollingTimePattern: "0 0 0 * * *",
		MaxAge:             90,
		ClearTimePattern:   "0 0 2 * * *",
	}
	return
}

// 新建名称为“log”的默认日志配置
func NewDefaultConfig() (cfg *MyZapConfig) {
	return NewConfigByName("log")
}

// 新建默认logger
func NewDefaultLogger() (logger *zap.Logger) {

	var cfg = NewDefaultConfig()
	if cfg != nil {
		logger = cfg.BuildLogger()
	}
	return
}
