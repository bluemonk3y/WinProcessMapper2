package mapper

import (
	"testing"
	"fmt"
)

func TestLineMatching(t *testing.T) {

	if (getLineType("10: File  (---)   C:\\Wi11ndows") != LINE_FILE) {
		t.Errorf("Should be FILE")
	}

	if (getLineType("GoogleCrashHandler64.exe pid: 4256 <unable to open process>") != LINE_PID) {
		t.Errorf("Should be PID")
	}
	if (getLineType("GoogleCrashHandler64.exe 4256 <unable to open process>") != LINE_OTHER) {
		t.Errorf("Should be OTHER")
	}


	//t.Logf("does it work:%d", PID)
	//t.Errorf("\t\tShould receive a \"%d\" status. ", OTHER)
}
func TestSmellyIntegration(t *testing.T) {

	stats := new(ServerStats)

	fmt.Println("\n\nStarting Integration test ======================================================")
	var processMap = make(map[int]PidMap)
	processHandles(stats, processMap)
	fmt.Println(len(processMap))

	//t.Logf("does it work:%d", PID)
	//t.Errorf("\t\tShould receive a \"%d\" status. ", OTHER)
}