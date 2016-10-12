package mapper

import (
	"fmt"
	"log"
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
var IGNORE_PATHS = []string { "C:\\Windows",  }

func logIt(msg string) {
	fmt.Println(msg)
	log.Println(msg)
}

func processHandles(stats *ServerStats, processMap map[int]PidMap) {

	logIt("process handles")

	handle := exec.Command("../etc/Handle.exe")
	handlesPipe,err := handle.StdoutPipe()
	handle.Start()
	handlesProcessReader := bufio.NewReader(handlesPipe)

	var currentPIdMap = new(PidMap)

	logIt("process handles================================")

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
		var lineType = getLineType(line);

		if (lineType == LINE_PID) {
			var parts = strings.Fields(line)

			if (currentPIdMap != nil) {
				logIt(fmt.Sprintf("%+v", currentPIdMap))
			}

			var aaa, _ =   strconv.Atoi(parts[2])
			currentPIdMap = &PidMap{name: parts[0], owner: parts[3], pid: aaa, files: []string{} }

			processMap[aaa] = *currentPIdMap

		} else if (lineType == LINE_FILE){
			//fmt.Printf("BBB %+v \n", currentPIdMap)

			var parts = strings.Fields(line)
			var sofile = parts[3]
			var fileIndex = strings.LastIndex(line, sofile)
			var ff = line[fileIndex:len(line)]
			currentPIdMap.files = append(currentPIdMap.files,ff)
			stats.file_handles += 1
		}
	}

	// Wait for the result of the command; also closes our end of the pipe
	err = handle.Wait()
	fmt.Printf("HANDLES >> process map size %d\n", len(processMap))

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