package mapper

import (
	"fmt"
	"log"
	"os/exec"
	"bufio"
	"io"
	"strings"
)


func processHandles(stats *ServerStats, processMap map[string]map[string][]string) {

	log.Println("process handles...")

	handle := exec.Command("./etc/Handle.exe")
	handlesPipe,err := handle.StdoutPipe()
	handle.Start()
	handlesProcessReader := bufio.NewReader(handlesPipe)

	var currentPIdMap = make(map[string][]string)


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

		//fmt.Print(line)
		if (isPidLine(line)) {

			if (!strings.Contains(line, "unable to open process")) {

				var parts = strings.Fields(line)
				currentPIdMap = make(map[string][]string)

				currentPIdMap["name"] = []string{parts[0] }
				currentPIdMap["pid"] = []string{parts[2] }
				currentPIdMap["owner"] = []string{parts[3] }
				currentPIdMap["files"] = []string{}
				processMap[parts[2]] = currentPIdMap
			}

		} else if (strings.Contains(line, " File ")){
			// B4: File  (---)   C:\Windows\System32\en-US\user32.dll.mui

			var parts = strings.Fields(line)
			var sofile = parts[3]
			var fileIndex = strings.LastIndex(line, sofile)
			var ff = line[fileIndex:len(line)]
			currentPIdMap["files"] = append(currentPIdMap["files"],ff)
			stats.file_handles += 1
		}
	}


	// Wait for the result of the command; also closes our end of the pipe
	err = handle.Wait()
	fmt.Printf("HANDLES >> process map size %d\n", len(processMap))

}
func isPidLine(line string) bool {
	return strings.Contains(line, "pid:")
}