package mapper

import (
	"fmt"
	"os/exec"
	"bufio"
	"io"
	"strings"
	"strconv"
)
const (
	LINE_FILE = iota
	LINE_PID
	LINE_OTHER
)
var IGNORE_PATHS = []string { "C:\\Windows",  "C:\\Program Files" }

func processHandles(path string, stats *ServerStats, processMap map[int]PidMap) {

	handle := exec.Command(path)
	handlesPipe,err := handle.StdoutPipe()
	err = handle.Start()
	if (err != nil) {
		err := fmt.Errorf("Failed to run process: %q", err)
		panic(err)
	}
	handlesProcessReader := bufio.NewReader(handlesPipe)

	var currentPIdMap = PidMap{name: "", owner: "", pid: 0, files: []string{} }

	for {
		var line string
		line, err = handlesProcessReader.ReadString('\n')
		if err == io.EOF  {
			// Good end of file with no partial line
			break
		}
		if err == io.EOF {
			err := fmt.Errorf("Last line not terminated: %q", line)
			panic(err)
		}
		if err != nil {
			panic(err)
		}
		var lineType = getLineType(line)

		if (lineType == LINE_PID) {
			if (currentPIdMap.pid != 0) {
				processMap[currentPIdMap.pid] = currentPIdMap
			}
			var parts = strings.Fields(line)

			var pid, _ =   strconv.Atoi(parts[2])
			currentPIdMap = PidMap{name: parts[0], owner: parts[3], pid: pid, files: []string{"aaa"} }

			stats.processes += 1

		} else if (lineType == LINE_FILE) {
			//fmt.Printf("BBB %+v \n", currentPIdMap)
			if (strings.Contains(line, "log")) {
				var parts = strings.Fields(line)
				var sofile = parts[3]
				var fileIndex = strings.LastIndex(line, sofile)
				var ff = line[fileIndex:len(line)]
				currentPIdMap.files = append(currentPIdMap.files, ff)
			}
			currentPIdMap.fileHandles += 1
			stats.file_handles += 1
		}
	}

	// Wait for the result of the command; also closes our end of the pipe
	err = handle.Wait()
}
func getLineType(line string) int {
	if (strings.Contains(line, ": File")) {
		for _, v := range IGNORE_PATHS {
			if (strings.Contains(line, v)) {
				return LINE_OTHER
			}
		}

		return LINE_FILE
	}
	if (strings.Contains(line, "pid: ")) {
		if (strings.Contains(line, "unable to open process")) {
			return  LINE_OTHER

		}
		return LINE_PID
	}
	return LINE_OTHER

}