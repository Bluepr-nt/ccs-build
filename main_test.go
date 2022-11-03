package main

import (
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPrecedence(t *testing.T) {
	// Run the tests in a temporary directory
	tmpDir, err := ioutil.TempDir("", "ccs")
	require.NoError(t, err, "error creating a temporary test directory")
	testDir, err := os.Getwd()
	require.NoError(t, err, "error getting the current working directory")
	defer os.Chdir(testDir)
	err = os.Chdir(tmpDir)
	require.NoError(t, err, "error changing to the temporary test directory")

	// Set favorite-color with the config file
	t.Run("config file", func(t *testing.T) {
		// Copy the config file into our temporary test directory
		configB, err := ioutil.ReadFile(filepath.Join(testDir, "ccs.yaml"))
		require.NoError(t, err, "error reading test config file")
		err = ioutil.WriteFile(filepath.Join(tmpDir, "ccs.yaml"), configB, 0644)
		require.NoError(t, err, "error writing test config file")
		defer os.Remove(filepath.Join(tmpDir, "ccs.yaml"))

		// Run ./stingoftheviper
		cmd := NewRootCommand()
		output := &bytes.Buffer{}
		cmd.SetArgs([]string{})
		cmd.SetOut(output)
		cmd.Execute()

		gotOutput := output.String()
		wantOutput := `My name is: gimlee
The mother's name is: morgul-blade
I live here: https://index.docker.io/v1
`
		assert.Equal(t, wantOutput, gotOutput, "expected the color from the config file and the number from the flag default")
	})

	// Set favorite-color with an environment variable
	t.Run("env var", func(t *testing.T) {
		os.Setenv("CCS_CONTAINER_REGISTRY_USERNAME", "sarouman")
		defer os.Unsetenv("CCS_CONTAINER_REGISTRY_USERNAME")

		cmd := NewRootCommand()
		output := &bytes.Buffer{}
		cmd.SetOut(output)
		cmd.Execute()

		gotOutput := output.String()
		wantOutput := `My name is: sarouman
The mother's name is: indis
I live here: https://index.docker.io/v1
`
		assert.Equal(t, wantOutput, gotOutput, "expected the color to use the environment variable value and the number to use the flag default")
	})

	// Set number with a flag
	t.Run("flag", func(t *testing.T) {
		// Run ./stingoftheviper --number 2
		cmd := NewRootCommand()
		output := &bytes.Buffer{}
		cmd.SetOut(output)
		cmd.SetArgs([]string{"--container-registry-username", "sauron"})
		cmd.Execute()

		gotOutput := output.String()
		wantOutput := `My name is: sauron
The mother's name is: indis
I live here: https://index.docker.io/v1
`
		assert.Equal(t, wantOutput, gotOutput, "expected the number to use the flag value and the color to use the flag default")
	})
}

func TestSanitizeInputs(t *testing.T) {
	type args struct {
		envCreds *environmentRegistry
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := SanitizeInputs(tt.args.envCreds); (err != nil) != tt.wantErr {
				t.Errorf("SanitizeInputs() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
