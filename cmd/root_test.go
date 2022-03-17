package cmd

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"testing"
)

func Test_ExecuteCommand(t *testing.T) {
	cmd := rootCmd
	b := bytes.NewBufferString("")
	cmd.SetOut(b)
	cmd.SetArgs([]string{"help", "serve"})
	err := cmd.Execute()
	assert.NoError(t, err)
	out, err := ioutil.ReadAll(b)
	assert.NoError(t, err)
	assert.Contains(t, string(out), "Cobra is a CLI library for Go that empowers applications.")
}
