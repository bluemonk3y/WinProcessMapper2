package main


import (
	"fmt"
	"log"
	"os/exec"
	"bufio"
	"io"
	"strings"
)

func main() {

	handle := exec.Command("./etc/Handle.exe")
	handlesPipe,err := handle.StdoutPipe()
	handle.Start()
	handlesProcessReader := bufio.NewReader(handlesPipe)


	var processMap = make(map[string]map[string][]string)
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

//		fmt.Print(line)
		if (strings.Contains(line, "pid:")) {

			var parts = strings.Fields(line)
			currentPIdMap = make(map[string][]string)

			currentPIdMap["name"] = []string { parts[0] }
			currentPIdMap["pid"] = []string { parts[2] }
			currentPIdMap["owner"] = []string { parts[3] }
			currentPIdMap["files"] = []string{}
			processMap["pid"] = currentPIdMap

		} else if (strings.Contains(line, " File ")){
			var parts = strings.Fields(line)
			var sofile = parts[3]
			var fileIndex = strings.LastIndex(line, sofile)
			var ff = line[fileIndex:len(line)]
			currentPIdMap["files"] = append(currentPIdMap["files"],ff)
		}
	}


	// Wait for the result of the command; also closes our end of the pipe
	err = handle.Wait()
	//fmt.Printf("HANDLES %s\n", handles)

	netout := exec.Command("netstat", "-a", "-n", "-o")
	netstatPipe,err2 := netout.StdoutPipe()

	netout.Start()
	netstatProcessReader := bufio.NewReader(netstatPipe)
	for {
		var line string
		line, err2 = netstatProcessReader.ReadString('\n')
		if err2 == io.EOF  {
			// Good end of file with no partial line
			break
		}
		if err2 == io.EOF {
			err2 := fmt.Errorf("Last line not terminated: %q", line)
			panic(err2)
		}
		if err2 != nil {
			panic(err2)
		}

		fmt.Print(line)
	 }



	if err2 != nil {
		log.Fatal(err)
	}
//	fmt.Printf("NETSTAT %s\n", netout)
}