package executor

import "time"

type Job struct {
	ID       int64
	Name     string
	Image    string
	WorkDir  string
	MountDir string
	Commands []string
	Env      []string
}

type Result struct {
	ExitCode int64
	Success  bool
	Duration time.Duration
	Logs     string
}
