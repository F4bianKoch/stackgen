package doctor

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"
)

func TestLevelFormatting(t *testing.T) {
	tests := []struct {
		level level
		icon  string
		label string
	}{
		{level: ok, icon: "✓", label: "OK"},
		{level: warn, icon: "⚠", label: "WARN"},
		{level: fail, icon: "✗", label: "FAIL"},
		{level: level(99), icon: "✗", label: "FAIL"},
	}

	for _, test := range tests {
		if got := test.level.icon(); got != test.icon {
			t.Errorf("level(%d).icon() = %q, want %q", test.level, got, test.icon)
		}
		if got := test.level.label(); got != test.label {
			t.Errorf("level(%d).label() = %q, want %q", test.level, got, test.label)
		}
	}
}

func TestRunCmd(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("stackgen development is Linux-first and this test uses sh")
	}

	t.Run("returns stdout", func(t *testing.T) {
		got, err := runCmd(time.Second, "sh", "-c", `printf 'hello'`)
		if err != nil {
			t.Fatalf("runCmd() error = %v", err)
		}
		if got != "hello" {
			t.Fatalf("runCmd() = %q, want %q", got, "hello")
		}
	})

	t.Run("includes stderr on failure", func(t *testing.T) {
		_, err := runCmd(time.Second, "sh", "-c", `printf 'specific failure' >&2; exit 7`)
		if err == nil || !strings.Contains(err.Error(), "specific failure") {
			t.Fatalf("runCmd() error = %v, want stderr context", err)
		}
	})

	t.Run("enforces timeout", func(t *testing.T) {
		started := time.Now()
		_, err := runCmd(50*time.Millisecond, "sh", "-c", "exec sleep 5")
		if err == nil {
			t.Fatal("runCmd() error = nil, want timeout failure")
		}
		if elapsed := time.Since(started); elapsed > time.Second {
			t.Fatalf("runCmd() took %v, timeout was not enforced", elapsed)
		}
	})
}

func TestRun_WithAvailableDependencies(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("fake executables use POSIX shell scripts")
	}
	binDir := t.TempDir()
	writeExecutable(t, filepath.Join(binDir, "docker"), "#!/bin/sh\nexit 0\n")
	writeExecutable(t, filepath.Join(binDir, "git"), "#!/bin/sh\nprintf 'git version test'\n")
	t.Setenv("PATH", binDir)
	workingDir := t.TempDir()
	chdir(t, workingDir)

	if err := Run(); err != nil {
		t.Fatalf("Run() error = %v", err)
	}
	if _, err := os.Stat(filepath.Join(workingDir, ".stackgen-doctor-tmp")); !os.IsNotExist(err) {
		t.Fatalf("filesystem probe directory was not cleaned up, stat error = %v", err)
	}
}

func TestRun_FailsWithoutDocker(t *testing.T) {
	t.Setenv("PATH", t.TempDir())
	chdir(t, t.TempDir())

	if err := Run(); err == nil {
		t.Fatal("Run() error = nil, want failed Docker checks")
	}
}

func writeExecutable(t *testing.T, path, contents string) {
	t.Helper()
	if err := os.WriteFile(path, []byte(contents), 0o755); err != nil {
		t.Fatal(err)
	}
}

func chdir(t *testing.T, dir string) {
	t.Helper()
	original, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := os.Chdir(original); err != nil {
			t.Errorf("restore working directory: %v", err)
		}
	})
}
