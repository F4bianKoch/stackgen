package cmd

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

type level int

const (
	ok level = iota
	warn
	fail
)

func (l level) icon() string {
	switch l {
	case ok:
		return "✓"
	case warn:
		return "⚠"
	default:
		return "✗"
	}
}

func (l level) label() string {
	switch l {
	case ok:
		return "OK"
	case warn:
		return "WARN"
	default:
		return "FAIL"
	}
}

type check struct {
	name string
	run  func() result
}

type result struct {
	lvl   level
	detail string
	hint   string
}

var doctorCmd = &cobra.Command{
	Use:   "doctor",
	Short: "Check whether this system can run stackgen (Docker, Compose, permissions, filesystem)",
	RunE: func(cmd *cobra.Command, args []string) error {
		return runDoctor()
	},
}

func init() {
	rootCmd.AddCommand(doctorCmd)
}

func runDoctor() error {
	fmt.Println("stackgen doctor")
	fmt.Println("--------------")

	checks := []check{
		{
			name: "OS",
			run: func() result {
				return result{
					lvl:    ok,
					detail: fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH),
				}
			},
		},
		{
			name: "Docker CLI",
			run: func() result {
				out, err := runCmd(3*time.Second, "docker", "--version")
				if err != nil {
					return result{
						lvl:    fail,
						detail: "docker not found in PATH",
						hint:   "Install Docker and ensure the `docker` command is available in your terminal.",
					}
				}
				return result{lvl: ok, detail: strings.TrimSpace(out)}
			},
		},
		{
			name: "Docker daemon",
			run: func() result {
				_, err := runCmd(6*time.Second, "docker", "info")
				if err != nil {
					return result{
						lvl:    fail,
						detail: "cannot reach docker daemon",
						hint:   "Start Docker (Docker Desktop on Windows/macOS, or the docker service on Linux).",
					}
				}
				return result{lvl: ok, detail: "reachable"}
			},
		},
		{
			name: "Docker permissions",
			run: func() result {
				// On Linux this catches: "permission denied" when user isn't in docker group.
				_, err := runCmd(6*time.Second, "docker", "ps")
				if err != nil {
					return result{
						lvl:    fail,
						detail: "docker commands fail (likely permissions)",
						hint:   "Ensure your user can run Docker. On Linux you may need to be in the `docker` group or use sudo.",
					}
				}
				return result{lvl: ok, detail: "can run docker commands"}
			},
		},
		{
			name: "Docker Compose v2",
			run: func() result {
				out, err := runCmd(3*time.Second, "docker", "compose", "version")
				if err != nil {
					return result{
						lvl:    fail,
						detail: "`docker compose` not available",
						hint:   "Update Docker / Docker Desktop to include Docker Compose v2 (`docker compose ...`).",
					}
				}
				return result{lvl: ok, detail: strings.TrimSpace(out)}
			},
		},
		{
			name: "Filesystem write access",
			run: func() result {
				cwd, err := os.Getwd()
				if err != nil {
					return result{
						lvl:    fail,
						detail: "cannot determine current directory",
						hint:   "Check filesystem permissions and try again.",
					}
				}
				tmpDir := filepath.Join(cwd, ".stackgen-doctor-tmp")
				tmpFile := filepath.Join(tmpDir, "test.txt")

				if err := os.MkdirAll(tmpDir, 0o755); err != nil {
					return result{
						lvl:    fail,
						detail: "cannot create directories in current location",
						hint:   "Run stackgen in a directory where you have write permissions.",
					}
				}
				if err := os.WriteFile(tmpFile, []byte("ok\n"), 0o644); err != nil {
					_ = os.RemoveAll(tmpDir)
					return result{
						lvl:    fail,
						detail: "cannot write files in current location",
						hint:   "Run stackgen in a directory where you have write permissions.",
					}
				}
				_ = os.RemoveAll(tmpDir)
				return result{lvl: ok, detail: "can create and write files"}
			},
		},
		{
			name: "Git (recommended)",
			run: func() result {
				out, err := runCmd(3*time.Second, "git", "--version")
				if err != nil {
					return result{
						lvl:    warn,
						detail: "git not found",
						hint:   "Install Git to enable repo initialization and common workflows (recommended).",
					}
				}
				return result{lvl: ok, detail: strings.TrimSpace(out)}
			},
		},
	}

	var hasFail bool
	for _, c := range checks {
		r := c.run()

		// Format: statusicon NAME [LEVEL] - detail
		fmt.Printf("%s %-22s [%s] %s\n", r.lvl.icon(), c.name, r.lvl.label(), r.detail)
		if r.hint != "" && (r.lvl == warn || r.lvl == fail) {
			fmt.Printf("  → %s\n", r.hint)
		}
		if r.lvl == fail {
			hasFail = true
		}
	}

	if hasFail {
		fmt.Println("\nstackgen cannot run on this system until the FAIL items are fixed.")
		return fmt.Errorf("doctor found failing checks")
	}

	fmt.Println("\nSystem looks ready for stackgen.")
	return nil
}

// runCmd executes a command with a timeout and returns stdout.
// If it fails, stderr is included in the returned error.
func runCmd(timeout time.Duration, name string, args ...string) (string, error) {
	cmd := exec.Command(name, args...)

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if timeout > 0 {
		timer := time.AfterFunc(timeout, func() {
			_ = cmd.Process.Kill()
		})
		defer timer.Stop()
	}

	err := cmd.Run()
	if err != nil {
		errText := strings.TrimSpace(stderr.String())
		if errText == "" {
			return "", err
		}
		return "", fmt.Errorf("%v: %s", err, errText)
	}
	return stdout.String(), nil
}
