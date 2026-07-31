package executor

import "time"

type Job struct {
	Image    string
	MountDir string
	WorkDir  string
	Commands []string
	Env      []string
}

type Result struct {
	ExitCode int64
	Success  bool
	Duration time.Duration
}
