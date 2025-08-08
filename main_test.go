package main

import (
	"os"
	"testing"
	"time"
)

func runAppAsync() {
	go func() {
		main()
	}()
	time.Sleep(100 * time.Millisecond)
}

func TestRunAppDevAndProd(t *testing.T) {
	origStage := os.Getenv("STAGE_STATUS")
	defer os.Setenv("STAGE_STATUS", origStage)

	os.Setenv("STAGE_STATUS", "dev")
	runAppAsync()

	os.Setenv("STAGE_STATUS", "prod")
	runAppAsync()
}
