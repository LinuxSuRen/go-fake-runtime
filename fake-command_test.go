package exec

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLookPath(t *testing.T) {
	fake := FakeExecer{
		ExpectLookPathError: errors.New("fake"),
		ExpectOutput:        "output",
		ExpectErrOutput:     "error",
		ExpectOS:            "os",
		ExpectArch:          "arch",
	}
	_, err := fake.LookPath("")
	assert.NotNil(t, err)

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
	assert.Nil(t, fake.RunCommandWithIO("", "", nil, nil))
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
}
