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
	"context"
	"os"
	osexec "os/exec"
	"runtime"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestRuntime(t *testing.T) {
	execer := NewDefaultExecer()
	assert.Equal(t, runtime.GOOS, execer.OS())
	assert.Equal(t, runtime.GOARCH, execer.Arch())
}

func TestDefaultLookPath(t *testing.T) {
	tests := []struct {
		name string
		arg  string
	}{{
		name: "ls",
		arg:  "ls",
	}, {
		name: "unknown",
		arg:  "unknown",
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			execer := NewDefaultExecer()

			expectPath, expectErr := osexec.LookPath(tt.arg)
			resultPath, resultErr := execer.LookPath(tt.arg)

			assert.Equal(t, expectPath, resultPath)
			assert.Equal(t, expectErr, resultErr)
		})
	}
}

func TestDefaultExecer(t *testing.T) {
	tests := []struct {
		name      string
		cmd       string
		args      []string
		expectErr bool
		verify    func(t *testing.T, out string)
	}{{
		name:      "go version",
		cmd:       "go",
		args:      []string{"version"},
		expectErr: false,
		verify: func(t *testing.T, out string) {
			assert.Contains(t, out, "go version")
		},
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ex := NewDefaultExecer()
			out, err := ex.Command(tt.cmd, tt.args...)
			assert.Equal(t, tt.expectErr, err != nil, err)
			if tt.verify != nil {
				tt.verify(t, string(out))
			}
			err = ex.RunCommand(tt.cmd, tt.args...)
			assert.Equal(t, tt.expectErr, err != nil, err)

			arch := ex.Arch()
			assert.Equal(t, runtime.GOARCH, arch)

			var outStr string
			outStr, err = ex.RunCommandAndReturn(tt.cmd, "", tt.args...)
			assert.Equal(t, tt.expectErr, err != nil, err)
			if tt.verify != nil {
				tt.verify(t, outStr)
			}

			err = ex.RunCommandWithEnv(tt.cmd, tt.args, nil, os.Stdout, os.Stderr)
			assert.Equal(t, tt.expectErr, err != nil, err)

			err = ex.RunCommandWithIO(tt.cmd, os.TempDir(), os.Stdout, os.Stderr, tt.args...)
			assert.Equal(t, tt.expectErr, err != nil, err)

			err = ex.RunCommandInDir(tt.cmd, "", tt.args...)
			assert.Equal(t, tt.expectErr, err != nil, err)

			err = ex.RunCommandWithBuffer(tt.cmd, "", nil, nil, tt.args...)
			assert.Equal(t, tt.expectErr, err != nil, err)
		})
	}
}

func TestRunCommandAndReturn(t *testing.T) {
	ex := NewDefaultExecer()
	result, err := ex.RunCommandAndReturn("go", "", "fake")
	assert.NotNil(t, err)
	assert.Equal(t, "go fake: unknown command\nRun 'go help' for usage.\n", result)
}

func TestMkdirAll(t *testing.T) {
	ex := NewDefaultExecer()
	err := ex.MkdirAll(os.TempDir(), 0755)
	assert.NoError(t, err)
}

func TestCommandWithContext(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.TODO(), time.Second)
	defer cancel()

	execer := NewDefaultExecerWithContext(ctx)
	err := execer.RunCommand("sleep 3")
	assert.Error(t, err)
}
