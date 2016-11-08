package mapper

import "testing"

func TestOSW_WinCPUGetter_Integration(t *testing.T) {

	t.Log("CPU Load:", Win_CpuLoad())
}