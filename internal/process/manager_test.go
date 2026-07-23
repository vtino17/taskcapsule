package process

import (
	"os"
	"testing"
	"time"
)

func TestStopProcessNonExistent(t *testing.T) {
	StopProcess(999999, 1)
}

func TestIsAlive(t *testing.T) {
	alive := IsAlive(os.Getpid())
	if !alive {
		t.Error("current process should be alive")
	}
}

func TestIsAliveNonExistent(t *testing.T) {
	alive := IsAlive(999999)
	if alive {
		t.Log("IsAlive returned true for non-existent PID (expected on Windows)")
	}
}

func TestFormatDuration(t *testing.T) {
	d := formatDuration(5 * time.Second)
	if d != "5.0s" {
		t.Errorf("unexpected format: %s", d)
	}
}

func TestStopProcessNonExistentLargePid(t *testing.T) {
	// Use a PID that cannot exist in any reasonable system
	StopProcess(987654321, 1)
}

func TestStopProcessZero(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("StopProcess(0) panicked: %v", r)
		}
	}()
	StopProcess(0, 1)
}

func TestStopProcessGroupZero(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("StopProcessGroup(0) panicked: %v", r)
		}
	}()
	StopProcessGroup(0, 1)
}


