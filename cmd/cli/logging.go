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

const debugLogLevel = 4

func klogInit(c *cli.Context) error {
	_ = flag.CommandLine.Parse([]string{})
	logLevel := c.Int("v")
	if debug && logLevel < debugLogLevel {
		logLevel = debugLogLevel
	}
	klogFlags := map[string]string{
		"v":                fmt.Sprint(logLevel),
		"logtostderr":      fmt.Sprint(c.Bool("logtostderr")),
		"stderrthreshold":  fmt.Sprint(c.Int("stderrthreshold")),
		"alsologtostderr":  fmt.Sprint(c.Bool("alsologtostderr")),
		"vmodule":          c.String("vmodule"),
		"log_dir":          c.String("log_dir"),
		"log_backtrace_at": c.String("log_backtrace_at"),
	}
	klog.InitFlags(nil)
	flag.VisitAll(func(fl *flag.Flag) {
		if val, ok := klogFlags[fl.Name]; ok {
			_ = fl.Value.Set(val)
		}
	})
	return nil
}
