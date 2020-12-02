package main

import (
	"flag"
	"fmt"
	"github.com/urfave/cli/v2"
	"k8s.io/klog/v2"
)

var klogCliFlags = []cli.Flag{
	&cli.IntFlag{
		Name: "v", Value: 0, Usage: "log level for V logs",
	},
	&cli.BoolFlag{
		Name: "logtostderr", Usage: "log to standard error instead of files",
	},
	&cli.IntFlag{
		Name:  "stderrthreshold",
		Usage: "logs at or above this threshold go to stderr",
	},
	&cli.BoolFlag{
		Name: "alsologtostderr", Usage: "log to standard error as well as files",
	},
	&cli.StringFlag{
		Name:  "vmodule",
		Usage: "comma-separated list of pattern=N settings for file-filtered logging",
	},
	&cli.StringFlag{
		Name: "log_dir", Usage: "If non-empty, write log files in this directory",
	},
	&cli.StringFlag{
		Name:  "log_backtrace_at",
		Usage: "when logging hits line file:N, emit a stack trace",
		Value: ":0",
	},
}

func klogInit(c *cli.Context) error {
	_ = flag.CommandLine.Parse([]string{})
	klogFlags := map[string]string{
		"v":                fmt.Sprint(c.Int("v")),
		"logtostderr":      fmt.Sprint(c.Bool("logtostderr")),
		"stderrthreshold":  fmt.Sprint(c.Int("stderrthreshold")),
		"alsologtostderr":  fmt.Sprint(c.Bool("alsologtostderr")),
		"vmodule":          c.String("vmodule"),
		"log_dir":          c.String("log_dir"),
		"log_backtrace_at": c.String("log_backtrace_at"),
	}
	flag.VisitAll(func(fl *flag.Flag) {
		if val, ok := klogFlags[fl.Name]; ok {
			fl.Value.Set(val)
		}
	})
	klog.InitFlags(nil)
	flag.Parse()
	return nil
}

// var toString = map[logrus.Level]string{
// 	logrus.PanicLevel: "panic",
// 	logrus.FatalLevel: "fatal",
// 	logrus.ErrorLevel: "error",
// 	logrus.WarnLevel:  "warn",
// 	logrus.InfoLevel:  "info",
// 	logrus.DebugLevel: "debug",
// 	logrus.TraceLevel: "trace",
// }
//
// var fromString = map[string]logrus.Level{
// 	"panic":   logrus.PanicLevel,
// 	"fatal":   logrus.FatalLevel,
// 	"error":   logrus.ErrorLevel,
// 	"warn":    logrus.WarnLevel,
// 	"info":    logrus.InfoLevel,
// 	"debug":   logrus.DebugLevel,
// 	"trace":   logrus.TraceLevel,
// 	"warning": logrus.WarnLevel,
// 	"tracing": logrus.TraceLevel,
// }
//
// const defaultLogLevel = logrus.WarnLevel
//
// func logrusInit(c *cli.Context) error {
//
// 	logrus.SetFormatter(&logrus.TextFormatter{
// 		FullTimestamp: true,
// 	})
//
// 	logrus.SetOutput(os.Stdout)
//
// 	logLevel := c.String("log-level")
// 	if verbose {
// 		logrus.SetLevel(logrus.TraceLevel)
// 	} else if debug {
// 		logrus.SetLevel(logrus.DebugLevel)
// 	} else {
// 		SetLevelWithString(logLevel)
// 	}
// 	return nil
// }
//
// func getLevelFromNum(num int) logrus.Level {
// 	if num < int(logrus.PanicLevel) || num > int(logrus.TraceLevel) {
// 		log.Printf("Unsupported log level: %d, fallback to default level: %d", num, defaultLogLevel)
// 		return defaultLogLevel
// 	}
// 	return logrus.Level(num)
// }
//
// func getLevelFromString(str string) logrus.Level {
// 	key := strings.ToLower(str)
// 	if level, ok := fromString[key]; ok {
// 		return level
// 	}
// 	log.Printf("Unsupported log level: %s, fallback to default level: %s", str, toString[defaultLogLevel])
// 	return defaultLogLevel
// }
//
// func SetLevelWithNum(num int) {
// 	level := getLevelFromNum(num)
// 	logrus.SetLevel(level)
// }
//
// func SetLevelWithString(str string) {
// 	level := getLevelFromString(str)
// 	logrus.SetLevel(level)
// }
//
