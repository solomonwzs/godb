package lmdb

import (
	"fmt"
	"runtime"
	"time"
)

var (
	Debug = true
)

func __debug(message string) {
	_, file, line, ok := runtime.Caller(2)
	if !ok {
		file = "???"
		line = 0
	}

	short := file
	for i := len(file) - 1; i > 0; i-- {
		if file[i] == '/' {
			short = file[i+1:]
			break
		}
	}

	fmt.Printf("\033[33m[%s %s:%d] \033[0m%s",
		time.Now().Format("15:04:05"), short, line, message)
}

func debugln(argv ...interface{}) {
	if !Debug {
		return
	}
	__debug(fmt.Sprintln(argv...))
}

func debugf(format string, argv ...interface{}) {
	if !Debug {
		return
	}
	__debug(fmt.Sprintf(format, argv...))
}
