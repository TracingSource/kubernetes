package workloads

import (
	"fmt"
	"os/exec"
	"regexp"
	"strings"
)

func runCmd(cmd string, args []string) error {
	err := exec.Command(cmd, args...).Run()
	return err
}

func getMatchingLineFromLog(log string, pattern string) (line string, err error) {
	regex, err := regexp.Compile(pattern)
	if err != nil {
		return line, fmt.Errorf("failed to compile regexp %v: %v", pattern, err)
	}

	logLines := strings.Split(log, "\n")
	for _, line := range logLines {
		if regex.MatchString(line) {
			return line, nil
		}
	}

	return line, fmt.Errorf("line with pattern %v not found in log", pattern)
}
