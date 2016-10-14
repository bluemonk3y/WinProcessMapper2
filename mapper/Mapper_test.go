package mapper


import (
	"testing"
	"fmt"
)
/**
* go test -v
 */
func TestHandlerAndNetstat_Integration(t *testing.T) {

	//if (true) {
	//	return
	//}

	stats := new(ServerStats)

	t.Log("stuff")

	fmt.Println("\n\nMapper: Starting FAT Integration test ======================================================")
	var processMap = make(map[int]PidMap)

	processHandles(stats, processMap)
	processNetstat(stats, processMap)


	// find process with largest number of file handles
	var largestFiles = 0
	var foundPid = 0
	for _, v := range processMap{
		//fmt.Printf("%d %d\n", v.pid, len(v.clients))
		if (len(v.clients) > largestFiles) {
			//fmt.Println("======:", v.pid, len(v.clients))
			foundPid = v.pid
			largestFiles = len(v.clients)
		}
	}

	fmt.Println("\t\t Most Clients: ", processMap[foundPid].pid, processMap[foundPid].fileHandles)
}

