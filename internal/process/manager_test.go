package process

import (
	"os"
	"os/exec"
	"runtime"
	"testing"
	"time"
)

func TestSetProcessGroup(t *testing.T) {
	cmd := exec.Command("go", "version")
	SetProcessGroup(cmd)
	if cmd.SysProcAttr == nil {
		t.Fatal("SysProcAttr should not be nil after SetProcessGroup")
	}
}

func TestStopProcessRunning(t *testing.T) {
	cmd := exec.Command("go", "version")
	if err := cmd.Start(); err != nil {
		t.Fatal(err)
	}

	StopProcess(cmd.Process.Pid, 2)
	// If we reach here without hanging, the stop worked
}

func TestStopProcessNonExistent(t *testing.T) {
	// This should not panic
	StopProcess(999999, 1)
}

func TestStopProcessZero(t *testing.T) {
	// This should not panic
	StopProcess(0, 1)
}

func TestGetProcessGroup(t *testing.T) {
	cmd := exec.Command("go", "version")
	SetProcessGroup(cmd)
	if err := cmd.Start(); err != nil {
		t.Fatal(err)
	}
	defer cmd.Wait()

	pgid := GetProcessGroup(cmd)
	if pgid <= 0 && runtime.GOOS != "windows" {
		t.Errorf("expected positive PGID, got %d", pgid)
	}
}

func TestStopProcessGroup(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("process group management is experimental on Windows")
	}
	cmd := exec.Command("sleep", "60")
	SetProcessGroup(cmd)
	if err := cmd.Start(); err != nil {
		t.Skip("sleep not available")
	}
	defer cmd.Process.Kill()

	pgid := GetProcessGroup(cmd)
	StopProcessGroup(pgid, 2)
}

func TestStopProcessGroupZero(t *testing.T) {
	StopProcessGroup(0, 1)
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
		t.Log("IsAlive returned true for non-existent PID (possible on Windows)")
	}
}

func TestFormatDuration(t *testing.T) {
	d := formatDuration(5 * time.Second)
	if d != "5.0s" {
		t.Errorf("unexpected format: %s", d)
	}
}
