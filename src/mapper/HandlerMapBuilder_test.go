package mapper

import (
	"testing"
)

func TestHMB_LineMatching(t *testing.T) {

	t.Errorf("TestLineMatching:Should be PID")

	if (getLineType("10: File  (---)   C:\\Wi11ndows") != LINE_FILE) {
		t.Errorf("TestLineMatching:Should be FILE")
	}

	if (getLineType("GoogleCrashHandler64.exe pid: 4256 ") != LINE_PID) {
		t.Errorf("TestLineMatching:Should be PID")
	}

	if (getLineType("GoogleCrashHandler64.exe 4256 <unable to open process>") != LINE_OTHER) {
		t.Errorf("TestLineMatching:Should be OTHER")
	}
}

/**
 *>go test -v  -run TestHMB_Integration
 */
func TestHMB_Integration(t *testing.T) {

	stats := new(ServerStats)

	var processMap = make(map[int]PidMap)
	processHandles(stats, processMap)

	// find process with largest number of file handles
	var largestFiles = 0
	var foundPid = 0
	for _, v := range processMap{
		if (len(v.files) > largestFiles) {
			foundPid = v.pid
			largestFiles = len(v.files)
		}
	}


	t.Log("tMost Files: ", foundPid, processMap[foundPid].name, len(processMap[foundPid].files))
}