package mapper

import (
	"bufio"
	"fmt"
	"strings"
	"os/exec"
	"io"
	"strconv"
)

const (
	PORT_LISTENING = iota
	PORT_ESTABLISHED
	PORT_OTHER
)

func processNetstat(stats *ServerStats, processMap map[int]PidMap) {
	//netout := exec.Command("netstat", "-a", "-n", "-o")
	netout := exec.Command("netstat", "-a", "-o")
	netstatPipe, err2 := netout.StdoutPipe()

	netout.Start()
	netstatProcessReader := bufio.NewReader(netstatPipe)
	for {
		var line string
		line, err2 = netstatProcessReader.ReadString('\n')
		if err2 == io.EOF {
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

		var lineType = getNetstatLineType(line)

		var parts = strings.Fields(line)
		// TCP 5
		if (len(parts) == 5) {


			var pid, _ = strconv.Atoi(parts[4])
			var address = parts[1]
			//fmt.Print(pid)

			vval, exists := processMap[pid]
			if (!exists) {
				vval = PidMap{pid: pid}
			}

			if (lineType == PORT_LISTENING) {
				//fmt.Print(line)
				vval.listening = append(vval.listening, address)
				stats.server_ports += 1
			}
			if (lineType == PORT_ESTABLISHED) {
				vval.clients = append(vval.clients, address)
				processMap[pid] = vval
				stats.client_ports += 1
			}
			processMap[pid] = vval
		}
	}
}

func getNetstatLineType(line string) int {
	if (strings.Contains(line, "LISTENING")) {
		return PORT_LISTENING
	}
	if (strings.Contains(line, "ESTABLISHED")) {
		return PORT_ESTABLISHED
	}
	return PORT_OTHER

}
