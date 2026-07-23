//go:build unix

package process

import (
	"os"
	"os/exec"
	"os/signal"
	"syscall"
	"testing"
	"time"
)

func startHelperProcess(t *testing.T) *exec.Cmd {
	t.Helper()
	cmd := exec.Command(os.Args[0], "-test.run=TestHelperProcess")
	cmd.Env = append(os.Environ(), "GO_WANT_HELPER_PROCESS=1")
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	if err := cmd.Start(); err != nil {
		t.Fatalf("failed to start helper process: %v", err)
	}
	t.Cleanup(func() {
		StopProcessGroup(cmd.Process.Pid, 3)
		cmd.Wait()
	})
	return cmd
}

func TestHelperProcess(t *testing.T) {
	if os.Getenv("GO_WANT_HELPER_PROCESS") != "1" {
		return
	}
	// Remain alive until signaled
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGTERM, syscall.SIGINT)
	<-sig
	os.Exit(0)
}

func TestStopHelperViaProcessGroup(t *testing.T) {
	cmd := startHelperProcess(t)
	time.Sleep(200 * time.Millisecond)

	pid := cmd.Process.Pid
	if !IsAlive(pid) {
		t.Fatal("helper should be alive after starting")
	}

	StopProcessGroup(pid, 3)
	time.Sleep(500 * time.Millisecond)

	if IsAlive(pid) {
		t.Error("helper should have been stopped")
	}
}

func TestProcessGroupID(t *testing.T) {
	cmd := startHelperProcess(t)
	time.Sleep(200 * time.Millisecond)

	pgid := GetProcessGroup(cmd)
	if pgid <= 0 {
		t.Fatalf("expected positive PGID, got %d", pgid)
	}

	StopProcessGroup(pgid, 3)
	time.Sleep(500 * time.Millisecond)

	if IsAlive(cmd.Process.Pid) {
		t.Error("helper should have been stopped via process group")
	}
}
