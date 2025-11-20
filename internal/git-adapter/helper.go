package gitadapter

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/cheetahbyte/centra/internal/config"
	"github.com/cheetahbyte/centra/internal/logger"
)

func UpdateRepo(dir string) error {
	pullOrigin := "origin"
	pullBranch := "main"

	privateKey := filepath.Join(config.GetKeysDir(), "id_ed25519")

	cmd := exec.Command("git", "pull", pullOrigin, pullBranch, "--ff-only")
	cmd.Dir = dir
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	env := os.Environ()
	env = append(env,
		"GIT_SSH_COMMAND=ssh -i "+privateKey+" -o IdentitiesOnly=yes -o StrictHostKeyChecking=accept-new",
	)
	cmd.Env = env

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("git pull failed: %v\nstdout: %s\nstderr: %s",
			err, stdout.String(), stderr.String())
	}
	changed := parseChangedFiles(stdout.String())
	logger := logger.AcquireLogger()
	logger.Info().Int("changed_files", changed).Msg("git pull suceeded")
	return nil
}

func CloneRepo(repoUrl string, dist string) error {
	privateKey := filepath.Join(config.GetKeysDir(), "id_ed25519")

	cmd := exec.Command("git", "clone", repoUrl, dist)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	env := os.Environ()
	env = append(env, "GIT_SSH_COMMAND=ssh -i "+privateKey+" -o IdentitiesOnly=yes -o StrictHostKeyChecking=accept-new")
	cmd.Env = env

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("git clone failed: %v\nstdout: %s\nstderr: %s",
			err, stdout.String(), stderr.String())
	}

	return nil
}

func parseChangedFiles(output string) int {
	lines := strings.Split(output, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasSuffix(line, "files changed") ||
			strings.HasSuffix(line, "file changed") {
			parts := strings.Split(line, " ")
			if len(parts) > 0 {
				n, err := strconv.Atoi(parts[0])
				if err == nil {
					return n
				}
			}
		}
	}
	return 0
}
