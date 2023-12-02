/*
MIT License

Copyright (c) 2023 Rick

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/

package exec

import (
	"bytes"
	"context"
	"io"
	"os"
	"os/exec"
	"runtime"
	"sync"
	"syscall"
)

// Execer is an interface for OS-related operations
type Execer interface {
	LookPath(string) (string, error)
	Command(name string, arg ...string) ([]byte, error)
	RunCommand(name string, arg ...string) (err error)
	RunCommandWithEnv(name string, argv, envv []string, stdout, stderr io.Writer) (err error)
	RunCommandInDir(name, dir string, args ...string) error
	RunCommandAndReturn(name, dir string, args ...string) (result string, err error)
	RunCommandWithSudo(name string, args ...string) (err error)
	RunCommandWithBuffer(name, dir string, stdout, stderr *bytes.Buffer, args ...string) error
	RunCommandWithIO(name, dir string, stdout, stderr io.Writer, args ...string) (err error)
	SystemCall(name string, argv []string, envv []string) (err error)
	MkdirAll(path string, perm os.FileMode) error
	OS() string
	Arch() string
}

const (
	// OSLinux is the alias of Linux
	OSLinux = "linux"
	// OSDarwin is the alias of Darwin
	OSDarwin = "darwin"
	// OSWindows is the alias of Windows
	OSWindows = "windows"
)

// DefaultExecer is a wrapper for the OS exec
type defaultExecer struct {
	ctx context.Context
}

func NewDefaultExecer() Execer {
return NewDefaultExecerWithContext(nil)
}

func NewDefaultExecerWithContext(ctx context.Context) Execer {
	if ctx == nil {
		ctx = context.TODO()
	}
	return &defaultExecer{
		ctx: ctx,
	}
}

// LookPath is the wrapper of os/exec.LookPath
func (e *defaultExecer) LookPath(file string) (string, error) {
	return exec.LookPath(file)
}

// Command is the wrapper of os/exec.Command
func (e *defaultExecer) Command(name string, arg ...string) ([]byte, error) {
	return exec.CommandContext(e.ctx, name, arg...).CombinedOutput()
}

// RunCommand runs a command
func (e *defaultExecer) RunCommand(name string, arg ...string) error {
	return e.RunCommandWithIO(name, "", os.Stdout, os.Stderr, arg...)
}

// RunCommandWithEnv runs a command with given Env
func (e *defaultExecer) RunCommandWithEnv(name string, argv, envv []string, stdout, stderr io.Writer) (err error) {
	command := exec.CommandContext(e.ctx, name, argv...)
	command.Env = envv
	//var stdout []byte
	//var errStdout error
	stdoutIn, _ := command.StdoutPipe()
	stderrIn, _ := command.StderrPipe()
	err = command.Start()
	if err == nil {
		// cmd.Wait() should be called only after we finish reading
		// from stdoutIn and stderrIn.
		// wg ensures that we finish
		var wg sync.WaitGroup
		wg.Add(1)
		go func() {
			_, _ = copyAndCapture(stdout, stdoutIn)
			wg.Done()
		}()

		_, _ = copyAndCapture(stderr, stderrIn)

		wg.Wait()

		err = command.Wait()
	}
	return
}

// RunCommandWithIO runs a command with given IO
func (e *defaultExecer) RunCommandWithIO(name, dir string, stdout, stderr io.Writer, args ...string) (err error) {
	command := exec.CommandContext(e.ctx, name, args...)
	if dir != "" {
		command.Dir = dir
	}
	//var stdout []byte
	//var errStdout error
	stdoutIn, _ := command.StdoutPipe()
	stderrIn, _ := command.StderrPipe()
	err = command.Start()
	if err == nil {
		// cmd.Wait() should be called only after we finish reading
		// from stdoutIn and stderrIn.
		// wg ensures that we finish
		var wg sync.WaitGroup
		wg.Add(1)
		go func() {
			_, _ = copyAndCapture(stdout, stdoutIn)
			wg.Done()
		}()

		_, _ = copyAndCapture(stderr, stderrIn)

		wg.Wait()

		err = command.Wait()
	}
	return
}

// MkdirAll is the wrapper of os.MkdirAll
func (e *defaultExecer) MkdirAll(path string, perm os.FileMode) error {
	return os.MkdirAll(path, perm)
}

// OS returns the os name
func (e *defaultExecer) OS() string {
	return runtime.GOOS
}

// Arch returns the os arch
func (e *defaultExecer) Arch() string {
	return runtime.GOARCH
}

// RunCommandAndReturn runs a command, then returns the output
func (e *defaultExecer) RunCommandAndReturn(name, dir string, args ...string) (result string, err error) {
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	if err = e.RunCommandWithBuffer(name, dir, stdout, stderr, args...); err == nil {
		result = stdout.String()
	} else {
		result = stdout.String()
		result += stderr.String()
	}
	return
}

// RunCommandWithBuffer runs a command with buffer
// stdout and stderr could be nil
func (e *defaultExecer) RunCommandWithBuffer(name, dir string, stdout, stderr *bytes.Buffer, args ...string) error {
	if stdout == nil {
		stdout = &bytes.Buffer{}
	}
	if stderr == nil {
		stderr = &bytes.Buffer{}
	}
	return e.RunCommandWithIO(name, dir, stdout, stderr, args...)
}

// RunCommandInDir runs a command
func (e *defaultExecer) RunCommandInDir(name, dir string, args ...string) error {
	return e.RunCommandWithIO(name, dir, os.Stdout, os.Stderr, args...)
}

// RunCommandWithSudo runs a command with sudo
func (e *defaultExecer) RunCommandWithSudo(name string, args ...string) (err error) {
	newArgs := make([]string, 0)
	newArgs = append(newArgs, name)
	newArgs = append(newArgs, args...)
	return e.RunCommand("sudo", newArgs...)
}

// SystemCall is the wrapper of syscall.Exec
func (e *defaultExecer) SystemCall(name string, argv []string, envv []string) (err error) {
	return syscall.Exec(name, argv, envv)
}

func copyAndCapture(w io.Writer, r io.Reader) ([]byte, error) {
	var out []byte
	buf := make([]byte, 1024)
	for {
		n, err := r.Read(buf[:])
		if n > 0 {
			d := buf[:n]
			out = append(out, d...)
			_, err := w.Write(d)
			if err != nil {
				return out, err
			}
		}
		if err != nil {
			// Read returns io.EOF at the end of file, which is not an error for us
			if err == io.EOF {
				err = nil
			}
			return out, err
		}
	}
}
