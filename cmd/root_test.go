package cmd

import (
	"bytes"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExecuteCommand(t *testing.T) {
	cmd := rootCmd
	b := bytes.NewBufferString("")
	cmd.SetOut(b)
	cmd.SetArgs([]string{"help", "worker"})
	err := cmd.Execute()
	assert.NoError(t, err)
	out, err := ioutil.ReadAll(b)
	assert.NoError(t, err)
	assert.Contains(t, string(out), "Cobra is a CLI library for Go that empowers applications.")
}
