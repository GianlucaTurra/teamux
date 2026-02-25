package tmux

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/GianlucaTurra/teamux/common"
)

// TODO: move errors to a proper file/package

type NoSuchDirectoryError struct {
	msg string
}

func (e NoSuchDirectoryError) Error() string {
	return e.msg
}

type Warning struct{ msg string }

func NewWarning(msg string) error {
	return Warning{msg}
}

func (e Warning) Error() string {
	return e.msg
}

// executeCommand Runs the given command and returns the stdErr or the err if
// one is returned from the command execution
func executeCommand(command string) error {
	args := strings.Fields(command)
	common.GetLogger().Info("tmux " + command)
	cmd := exec.Command("tmux", args...)
	if output, err := cmd.CombinedOutput(); err != nil {
		if strings.TrimSpace(string(output)) == "" {
			return err
		} else {
			return errors.New(string(output))
		}
	} else {
		return nil
	}
}

func checkDirectory(directory string) error {
	if strings.TrimSpace(directory) == "" {
		return nil
	}
	var path string
	if strings.TrimSpace(directory) == "~" {
		return nil
	}
	if strings.HasPrefix(directory, "~/") {
		home, err := os.UserHomeDir()
		if err != nil {
			return err
		}
		path = filepath.Join(home, directory[2:])
	}
	if strings.HasPrefix(directory, "$HOME") {
		path = os.ExpandEnv(directory)
	}
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return err
	}
	return nil
}

func commandWithWorkDir(workingDirectory string, cmd string, extra string) error {
	var nsdErr error
	if err := checkDirectory(workingDirectory); err != nil {
		nsdErr = NoSuchDirectoryError{"working directory doesn't exist"}
	}
	if strings.TrimSpace(workingDirectory) != "" {
		cmd += fmt.Sprintf(" -c %s", workingDirectory)
	}
	if strings.TrimSpace(extra) != "" {
		cmd += fmt.Sprintf(" \"%s\"", extra)
	}
	if err := executeCommand(cmd); err != nil {
		return err
	} else {
		return nsdErr
	}
}
