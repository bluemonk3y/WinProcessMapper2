package mapper


import (
	"testing"
)
/**
* go test -v
* go test -v  -run TestMapper_INT_ALL
 */
func TestMapper_INT_ALL(t *testing.T) {

	stats := new(ServerStats)

	var processMap = make(map[int]PidMap)

	processHandles(stats, processMap)
	processNetstat(stats, processMap)
	finaliseServerStats(stats, processMap)

	t.Log(stats)

	writeServerStatsToInflux(stats)
}

