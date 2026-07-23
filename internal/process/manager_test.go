package process

import (
	"os"
	"os/exec"
	"runtime"
	"testing"
	"time"
)

func TestSetProcessGroup(t *testing.T) {
	cmd := exec.Command("sh", "-c", "sleep 60")
	if err := cmd.Start(); err != nil {
		t.Skip("sh not available")
	}
	SetProcessGroup(cmd)
	t.Cleanup(func() {
		StopProcess(cmd.Process.Pid, 2)
	})
	if cmd.SysProcAttr == nil && runtime.GOOS != "windows" {
		t.Error("SysProcAttr should not be nil")
	}
}

func TestStopProcessRunning(t *testing.T) {
	cmd := exec.Command("sleep", "60")
	if err := cmd.Start(); err != nil {
		t.Skip("sleep not available")
	}

	pid := cmd.Process.Pid

	// Wait a moment for the process to start
	time.Sleep(100 * time.Millisecond)

	StopProcess(pid, 2)

	// Wait and check it's stopped
	time.Sleep(500 * time.Millisecond)
	alive := IsAlive(pid)
	if alive && runtime.GOOS != "windows" {
		t.Error("process should have been stopped")
	}
}

func TestStopProcessNonExistent(t *testing.T) {
	StopProcess(999999, 1)
}

func TestGetProcessGroup(t *testing.T) {
	cmd := exec.Command("sleep", "60")
	if err := cmd.Start(); err != nil {
		t.Skip("sleep not available")
	}
	SetProcessGroup(cmd)
	t.Cleanup(func() { StopProcess(cmd.Process.Pid, 2) })

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
	pid := cmd.Process.Pid

	pgid := GetProcessGroup(cmd)
	StopProcessGroup(pgid, 2)

	time.Sleep(500 * time.Millisecond)
	alive := IsAlive(pid)
	if alive {
		t.Error("process group should have been stopped")
	}
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
		t.Log("IsAlive unexpectedly returned true (possible on Windows)")
	}
}

func TestStartStopChildProcess(t *testing.T) {
	cmd := exec.Command("sh", "-c", "trap 'exit 0' TERM INT; while true; do sleep 1; done")
	if err := cmd.Start(); err != nil {
		t.Skip("long-lived child not supported on this platform")
	}
	pid := cmd.Process.Pid

	time.Sleep(200 * time.Millisecond)
	if !IsAlive(pid) {
		t.Fatal("child should be alive after starting")
	}

	StopProcess(pid, 3)

	time.Sleep(500 * time.Millisecond)
	alive := IsAlive(pid)
	if alive && runtime.GOOS != "windows" {
		t.Error("child should have been stopped")
	}
	cmd.Process.Kill()
	cmd.Wait()
}

func TestStopProcessGroupZero(t *testing.T) {
	StopProcessGroup(0, 1)
}

func TestFormatDuration(t *testing.T) {
	d := formatDuration(5 * time.Second)
	if d != "5.0s" {
		t.Errorf("unexpected format: %s", d)
	}
}
