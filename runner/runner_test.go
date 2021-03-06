package runner

import (
	"errors"
	"os"
	"os/exec"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestRunnerWalk(t *testing.T) {
	t.Run("No files change should return last modify time", func(t *testing.T) {
		task := &Runner{dir: "."}
		now := time.Now()

		mod := task.Walk(now)

		assert.True(t, mod.Equal(now), "should return last modify time.")
	})

	t.Run("File chagne", func(t *testing.T) {
		task := &Runner{dir: "."}
		form := "Mon Jan _2 15:04:05 2006"
		lastMod, _ := time.Parse(form, "Sat Feb 08 07:00:00 1992")

		mod := task.Walk(lastMod)

		assert.True(t, mod.After(lastMod), "should return lastest modify time.")
	})
}

type TRunner struct {
	isKillCalled bool
	killReturn   error

	isStartCalled bool
	startReturn   error
}

func (r *TRunner) KillCommand() error {
	r.isKillCalled = true
	return r.killReturn
}

func (r *TRunner) Start() error {
	r.isStartCalled = true
	return r.startReturn
}

func TestRunnerRun(t *testing.T) {
	t.Run("kill command success then should call Start and return nil", func(t *testing.T) {
		tr := &TRunner{
			killReturn:  nil,
			startReturn: nil,
		}

		err := run(tr)

		assert.Nil(t, err, "should run comamnd success but it have error")
		assert.True(t, tr.isKillCalled, "should have been called Kill command but it not.")
		assert.True(t, tr.isStartCalled, "should have been called Start command but it not.")
	})

	t.Run("should return error when can't start the command", func(t *testing.T) {
		errMsg := "start error"
		tr := &TRunner{
			killReturn:  nil,
			startReturn: errors.New(errMsg),
		}

		err := run(tr)

		assert.Error(t, err, "should return an error but it not.")
	})

	t.Run("should return error when can't kill the command", func(t *testing.T) {
		errMsg := "kill command error"
		tr := &TRunner{
			killReturn:  errors.New(errMsg),
			startReturn: nil,
		}

		err := run(tr)

		assert.Error(t, err, "should return an error but it not.")
	})
}

func TestRunnerStart(t *testing.T) {
	t.Run("should return nil when command execute successfully", func(t *testing.T) {
		task := &Runner{
			prog:   "echo",
			args:   []string{"This is working"},
			dir:    ".",
			stderr: os.Stderr,
			stdout: os.Stdout,
		}

		expectedCmd := exec.Command("echo", "This is working")

		err := task.Start()
		assert.NoError(t, err, "should run comamnd success but it have error")
		assert.Equal(t, expectedCmd.Args, task.cmd.Args, "should run the same command with the initiated one but it doesn't")
	})

	t.Run("should return error when command fail to execute", func(t *testing.T) {
		task := &Runner{
			prog:   "fakecommand",
			args:   []string{"This is working"},
			dir:    ".",
			stderr: os.Stderr,
			stdout: os.Stdout,
		}

		err := task.Start()
		assert.Error(t, err, "should return an error but it not.")
	})
}
