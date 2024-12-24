/*
MIT License

Copyright (c) 2023-2024 Rick

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
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLookPath(t *testing.T) {
	fake := &FakeExecer{
		ExpectLookPathError: errors.New("fake"),
		ExpectOutput:        "output",
		ExpectErrOutput:     "error",
		ExpectOS:            "os",
		ExpectArch:          "arch",
		ExpectLookPath:      "lookpath",
	}
	fake.WithContext(context.Background())
	targetPath, err := fake.LookPath("")
	assert.NotNil(t, err)
	assert.Equal(t, "lookpath", targetPath)

	fake.ExpectLookPathError = nil
	_, err = fake.LookPath("")
	assert.Nil(t, err)

	var output []byte
	output, err = fake.Command("fake")
	assert.Equal(t, "output", string(output))
	assert.Nil(t, err)
	assert.Equal(t, "os", fake.OS())
	assert.Equal(t, "arch", fake.Arch())
	assert.Nil(t, fake.RunCommand("", ""))
	err = fake.RunCommandWithIO("", "", nil, nil, nil)
	assert.Nil(t, err)
	assert.Nil(t, fake.RunCommandWithEnv("", nil, nil, nil, nil))
	assert.Nil(t, fake.RunCommandInDir("", ""))

	var result string
	result, err = fake.RunCommandAndReturn("", "")
	assert.Equal(t, "output", result)
	assert.Nil(t, err)
	assert.Nil(t, fake.RunCommandWithSudo("", ""))
	assert.Nil(t, fake.RunCommandWithBuffer("", "", nil, nil))
	assert.Nil(t, fake.SystemCall("", nil, nil))

	fakeWithErr := FakeExecer{
		ExpectError:     errors.New("fake"),
		ExpectOutput:    "output",
		ExpectErrOutput: "error",
	}
	result, err = fakeWithErr.RunCommandAndReturn("", "")
	assert.Equal(t, "outputerror", result)
	assert.NotNil(t, err)
	assert.Error(t, fakeWithErr.MkdirAll("", 0))

	assert.Nil(t, fake.Setenv("key", "value"))
	assert.Equal(t, "value", fake.Getenv("key"))
}
