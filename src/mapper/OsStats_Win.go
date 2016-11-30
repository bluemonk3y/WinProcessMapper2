package mapper

import (
	"bufio"
	"os/exec"
	"fmt"
	"io"
	"strings"
	"strconv"
)

var WIN_CPU_EXEC = []string { "cscript.exe", "//NoLogo", "..\\etc\\cpuTime.vbs" }


func Win_CpuLoad() float32 {

	handle := exec.Command(WIN_CPU_EXEC[0], WIN_CPU_EXEC[1], WIN_CPU_EXEC[2])
	handlesPipe,err := handle.StdoutPipe()
	handle.Start()
	handlesProcessReader := bufio.NewReader(handlesPipe)


	var cpu = float32(0)
	for {
		var line string
		line, err = handlesProcessReader.ReadString('\n')
		if err == io.EOF {
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

		var parts = strings.Split(line, ",")
		if (len(parts) > 1) {
			var acpu, _ =   strconv.ParseFloat(parts[2], 32)
			cpu = float32(acpu)
		}
	}
	// Wait for the result of the command; also closes our end of the pipe
	err = handle.Wait()
	return cpu
}