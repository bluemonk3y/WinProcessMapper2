package main

import (
	"github.com/bluemonk3y/WinProcessMapper2/mapper"
)

func main() {

	stats := new(mapper.ServerStats)

	var processMap = make(map[int]mapper.PidMap)

	mapper.processHandles(stats, processMap)
	mapper.processNetstat(stats, processMap)
	mapper.finaliseServerStats(stats, processMap)

	//t.Log(stats)

	mapper.writeServerStatsToInflux(stats)
}