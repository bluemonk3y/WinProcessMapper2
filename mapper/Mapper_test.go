package mapper


import (
	"testing"
)
/**
* go test -v
 */
func Test_INT_ALL(t *testing.T) {

	stats := new(ServerStats)

	var processMap = make(map[int]PidMap)

	processHandles(stats, processMap)
	processNetstat(stats, processMap)
	finaliseServerStats(stats, processMap)

	t.Log(stats)
}

