package helper

import (
	"fmt"
	"os"
	"strings"
)

func ExitOnError(err error, args ...interface{}) {
	// always check args
	if len(args) > 0 {
		if _, ok := args[0].(string); !ok {
			panic("first arg must be string")
		}
	}

	if err != nil {
		msg := "ERR:"
		if len(args) > 0 {
			msg += " " + fmt.Sprintf(args[0].(string), args[1:]...)
		}
		fmt.Println(msg, err.Error())
		os.Exit(1)
	}
}

func ExitWithMessage(format string, a ...interface{}) {
	fmt.Println(fmt.Sprintf(format, a...))
	os.Exit(1)
}

func ReadLines(file string) []string {
	data, err := os.ReadFile(file)
	ExitOnError(err)
	lines := strings.Split(string(data), "\n")
	for i := range lines {
		lines[i] = strings.Trim(lines[i], "\r")
	}
	return lines
}

func ReadString(file string) string {
	data, err := os.ReadFile(file)
	ExitOnError(err)
	return string(data)
}
