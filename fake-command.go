package exec

import (
	"bytes"
	"io"
	"os"
)

// FakeExecer is for the unit test purposes
type FakeExecer struct {
	ExpectError         error
	ExpectLookPathError error
	ExpectOutput        string
	ExpectErrOutput     string
	ExpectOS            string
	ExpectArch          string
	ExpectLookPath      string
}

// LookPath is a fake method
func (f FakeExecer) LookPath(path string) (string, error) {
	return f.ExpectLookPath, f.ExpectLookPathError
}

// Command is a fake method
func (f FakeExecer) Command(name string, arg ...string) ([]byte, error) {
	return []byte(f.ExpectOutput), f.ExpectError
}

// RunCommand runs a command
func (f FakeExecer) RunCommand(name string, arg ...string) error {
	return f.ExpectError
}

// RunCommandInDir is a fake method
func (f FakeExecer) RunCommandInDir(name, dir string, args ...string) error {
	return f.ExpectError
}

// RunCommandAndReturn is a fake method
func (f FakeExecer) RunCommandAndReturn(name, dir string, args ...string) (result string, err error) {
	if err = f.ExpectError; err == nil {
		result = f.ExpectOutput
	} else {
		result = f.ExpectOutput
		result += f.ExpectErrOutput
	}
	return
}

// RunCommandWithSudo is a fake method
func (f FakeExecer) RunCommandWithSudo(name string, args ...string) (err error) {
	return f.ExpectError
}

// RunCommandWithBuffer is a fake method
func (f FakeExecer) RunCommandWithBuffer(name, dir string, stdout, stderr *bytes.Buffer, args ...string) error {
	return f.ExpectError
}

// RunCommandWithIO is a fake method
func (f FakeExecer) RunCommandWithIO(name, dir string, stdout, stderr io.Writer, args ...string) error {
	return f.ExpectError
}

// RunCommandWithEnv is a fake method
func (f FakeExecer) RunCommandWithEnv(name string, argv, envv []string, stdout, stderr io.Writer) error {
	return f.ExpectError
}

// SystemCall is a fake method
func (f FakeExecer) SystemCall(name string, argv []string, envv []string) error {
	return f.ExpectError
}

// MkdirAll is the wrapper of os.MkdirAll
func (f FakeExecer) MkdirAll(path string, perm os.FileMode) error {
	return f.ExpectError
}

// OS returns the os name
func (f FakeExecer) OS() string {
	return f.ExpectOS
}

// Arch returns the os arch
func (f FakeExecer) Arch() string {
	return f.ExpectArch
}
