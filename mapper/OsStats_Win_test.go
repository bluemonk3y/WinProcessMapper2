package mapper

import "testing"

func TestWinCPUGetter_Integration(t *testing.T) {

	t.Log("CPU Load:", GetWinStats())
}