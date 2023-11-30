package main

import (
	"aoc/helper"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func main() {
	dirs, err := os.ReadDir(".")
	helper.ExitOnError(err)

	for _, d := range dirs {
		if d.IsDir() && strings.HasPrefix(strings.ToLower(d.Name()), "puzzle-") {
			fmt.Println("## Puzzle " + d.Name()[7:] + ": ##")
			cmd := exec.Command("go", "run", ".")
			absDir, err := filepath.Abs("./" + d.Name())
			helper.ExitOnError(err, "get abs filepath")
			cmd.Dir = absDir
			out, err := cmd.CombinedOutput()
			helper.ExitOnError(err, "exec %s", d.Name())
			fmt.Print(string(out))
			fmt.Println()
		}
	}
}
