package mapper

import (
	"testing"
	"fmt"
)

func TestNetstatStruct_Integration(t *testing.T) {

	//if (true) {
	//	return
	//}

	stats := new(ServerStats)

	var processMap = make(map[int]PidMap)

	processNetstat(stats, processMap)


	// find process with largest number of file handles
	var largestFiles = 0
	var foundPid = 0
	for _, v := range processMap{
		if (len(v.clients) > largestFiles) {
			foundPid = v.pid
			largestFiles = len(v.clients)
		}
	}

	fmt.Println("\t\tMost Clients: ", foundPid, len(processMap[foundPid].clients))
}
