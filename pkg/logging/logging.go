package logging

import (
	"github.com/sirupsen/logrus"
	"os"
	"strings"
)

var (
	log *logrus.Logger
)

const (
	PanicLevel = logrus.PanicLevel
	FatalLevel = logrus.FatalLevel
	ErrorLevel = logrus.ErrorLevel
	WarnLevel  = logrus.WarnLevel
	InfoLevel  = logrus.InfoLevel
	DebugLevel = logrus.DebugLevel
	TraceLevel = logrus.TraceLevel
)

var toString = map[logrus.Level]string{
	PanicLevel: "panic",
	FatalLevel: "fatal",
	ErrorLevel: "error",
	WarnLevel:  "warn",
	InfoLevel:  "info",
	DebugLevel: "debug",
	TraceLevel: "trace",
}

var fromString = map[string]logrus.Level{
	"panic":   PanicLevel,
	"fatal":   FatalLevel,
	"error":   ErrorLevel,
	"warn":    WarnLevel,
	"info":    InfoLevel,
	"debug":   DebugLevel,
	"trace":   TraceLevel,
	"warning": WarnLevel,
	"tracing": TraceLevel,
}

const defaultLogLevel = TraceLevel

func init() {

	log = logrus.New()

	log.SetReportCaller(true)

	log.SetOutput(os.Stdout)
	log.SetLevel(defaultLogLevel)
}

func SetLevel(level logrus.Level) {
	log.SetLevel(level)
}

func getLevelFromNum(num int) logrus.Level {
	if num < int(PanicLevel) || num > int(TraceLevel) {
		log.Printf("Unsupported log level: %d, fallback to default level: %d", num, defaultLogLevel)
		return defaultLogLevel
	}
	return logrus.Level(num)
}

func getLevelFromString(str string) logrus.Level {
	key := strings.ToLower(str)
	if level, ok := fromString[key]; ok {
		return level
	}
	log.Printf("Unsupported log level: %s, fallback to default level: %s", str, toString[defaultLogLevel])
	return defaultLogLevel
}

func SetLevelWithNum(num int) {
	level := getLevelFromNum(num)
	SetLevel(level)
}

func SetLevelWithString(str string) {
	level := getLevelFromString(str)
	SetLevel(level)
}

func Trace(format string, v ...interface{}) {
	log.Tracef(format, v...)
}

func Debug(format string, v ...interface{}) {
	log.Debugf(format, v...)
}

func Info(format string, v ...interface{}) {
	log.Infof(format, v...)
}

func Warn(format string, v ...interface{}) {
	log.Warnf(format, v...)
}

func Error(format string, v ...interface{}) {
	log.Errorf(format, v...)
}

func Fatal(format string, v ...interface{}) {
	log.Fatalf(format, v...)
}

func Panic(format string, v ...interface{}) {
	log.Panicf(format, v...)
}
