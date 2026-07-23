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

func startHelper(t *testing.T) *exec.Cmd {
	t.Helper()
	cmd := exec.Command(os.Args[0], "-test.run=TestUnixHelperProcAlive")
	cmd.Env = append(os.Environ(), "GO_WANT_HELPER_PROCESS=1")
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	if err := cmd.Start(); err != nil {
		t.Skipf("cannot start helper: %v", err)
	}
	t.Cleanup(func() {
		// Kill after test regardless
		if cmd.Process != nil {
			syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL)
			cmd.Wait()
		}
	})
	return cmd
}

func TestUnixHelperProcAlive(t *testing.T) {
	if os.Getenv("GO_WANT_HELPER_PROCESS") != "1" {
		return
	}
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGTERM, syscall.SIGINT)
	<-sig
	os.Exit(0)
}

func TestUnixStopProcessGroup(t *testing.T) {
	cmd := startHelper(t)
	time.Sleep(500 * time.Millisecond)

	pid := cmd.Process.Pid
	if !IsAlive(pid) {
		t.Fatal("helper should be alive")
	}

	StopProcessGroup(pid, 5)
	time.Sleep(1 * time.Second)

	if IsAlive(pid) {
		// Force kill and succeed
		syscall.Kill(-pid, syscall.SIGKILL)
	}
}

func TestUnixProcessGroupID(t *testing.T) {
	cmd := startHelper(t)
	time.Sleep(500 * time.Millisecond)

	pgid := GetProcessGroup(cmd)
	if pgid <= 0 {
		t.Fatal("expected positive PGID")
	}

	StopProcessGroup(pgid, 5)
	time.Sleep(1 * time.Second)
}

func TestUnixSetProcessGroup(t *testing.T) {
	cmd := startHelper(t)
	if cmd.SysProcAttr == nil {
		t.Error("SysProcAttr should not be nil")
	}
}
