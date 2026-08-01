package executor

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"

	"github.com/moby/moby/api/pkg/stdcopy"
	"github.com/moby/moby/api/types/container"
	"github.com/moby/moby/client"
)

type DockerExecutor struct {
	cli *client.Client
}

func NewDockerExecutor() (*DockerExecutor, error) {
	cli, err := client.New(client.FromEnv)
	if err != nil {
		return nil, fmt.Errorf("failed to create docker client: %w", err)
	}
	return &DockerExecutor{cli: cli}, nil
}

func (d *DockerExecutor) Execute(ctx context.Context, job Job) (*Result, error) {
	startTime := time.Now()

	// 1. Pull Image if missing
	imageRef := job.Image
	if !strings.Contains(imageRef, "/") {
		imageRef = "docker.io/library/" + imageRef
	}

	log.Printf("Pulling image %s...", imageRef)
	reader, err := d.cli.ImagePull(ctx, imageRef, client.ImagePullOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to pull image %s: %w", imageRef, err)
	}
	defer reader.Close()

	if _, err := io.Copy(io.Discard, reader); err != nil {
		return nil, fmt.Errorf("error reading image pull stream: %w", err)
	}

	// 2. Set default working directory
	workDir := job.WorkDir
	if workDir == "" {
		workDir = "/workspace"
	}

	// 3. Build container execution command
	var cmd []string
	if len(job.Commands) > 0 {
		joinedCmds := strings.Join(job.Commands, " && ")
		cmd = []string{"/bin/sh", "-c", joinedCmds}
	}

	// 4. Create container
	log.Println("Creating Docker container...")
	config := &container.Config{
		Image:      job.Image,
		Cmd:        cmd,
		WorkingDir: workDir,
		Env:        job.Env,
	}

	hostConfig := &container.HostConfig{}
	if job.MountDir != "" {
		hostConfig.Binds = []string{fmt.Sprintf("%s:%s", job.MountDir, workDir)}
	}

	created, err := d.cli.ContainerCreate(ctx, client.ContainerCreateOptions{
		Config:     config,
		HostConfig: hostConfig,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create container: %w", err)
	}

	containerID := created.ID
	defer func() {
		log.Println("Removing container...")
		_, removeErr := d.cli.ContainerRemove(context.Background(), containerID, client.ContainerRemoveOptions{Force: true})
		if removeErr != nil {
			log.Printf("Warning: failed to remove container %s: %v", containerID, removeErr)
		}
	}()

	// 5. Start container
	log.Println("Container started")
	log.Println("Executing build...")
	if _, err := d.cli.ContainerStart(ctx, containerID, client.ContainerStartOptions{}); err != nil {
		return nil, fmt.Errorf("failed to start container: %w", err)
	}

	// 6. Wait for container execution completion
	waitRes := d.cli.ContainerWait(ctx, containerID, client.ContainerWaitOptions{
		Condition: container.WaitConditionNotRunning,
	})

	var exitCode int64 = -1
	select {
	case err := <-waitRes.Error:
		if err != nil {
			return nil, fmt.Errorf("error waiting for container: %w", err)
		}
	case status := <-waitRes.Result:
		exitCode = status.StatusCode
	case <-ctx.Done():
		return nil, ctx.Err()
	}

	duration := time.Since(startTime)
	log.Printf("Build completed. Exit Code: %d", exitCode)

	// 7. Stream container stdout and stderr logs
	var logBuf bytes.Buffer
	outStream, logErr := d.cli.ContainerLogs(ctx, containerID, client.ContainerLogsOptions{
		ShowStdout: true,
		ShowStderr: true,
	})
	if logErr == nil && outStream != nil {
		defer outStream.Close()
		log.Println("--- CONTAINER STDOUT & STDERR LOGS ---")
		stdoutWriter := io.MultiWriter(os.Stdout, &logBuf)
		stderrWriter := io.MultiWriter(os.Stderr, &logBuf)
		_, _ = stdcopy.StdCopy(stdoutWriter, stderrWriter, outStream)
		log.Println("--- END CONTAINER LOGS ---")
	} else if logErr != nil {
		log.Printf("Warning: failed to retrieve container logs: %v", logErr)
	}

	return &Result{
		ExitCode: exitCode,
		Success:  exitCode == 0,
		Duration: duration,
		Logs:     logBuf.String(),
	}, nil
}
