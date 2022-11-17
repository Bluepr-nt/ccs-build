package main

import (
	"bytes"
	"flag"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"k8s.io/klog/v2"
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
	klog.InitFlags(nil)

	flag.Set("logtostderr", "false")
	flag.Set("alsologtostderr", "false")
	flag.Parse()
	// Set favorite-color with the config file
	t.Run("config file", func(t *testing.T) {
		// Copy the config file into our temporary test directory
		configB, err := ioutil.ReadFile(filepath.Join(testDir, "ccs.yaml"))
		require.NoError(t, err, "error reading test config file")
		err = ioutil.WriteFile(filepath.Join(tmpDir, "ccs.yaml"), configB, 0644)
		require.NoError(t, err, "error writing test config file")
		defer os.Remove(filepath.Join(tmpDir, "ccs.yaml"))

		output := new(bytes.Buffer)
		cmd := NewRootCommand(output)
		cmd.SetArgs([]string{"--dry-run"})
		cmd.SetOut(output)
		cmd.Execute()
		klog.Flush()
		gotOutput := output.String()
		wantOutput := `Successful login of gimlee to https://index.docker.io/v1`
		assert.Contains(t, gotOutput, wantOutput)

		// assert.Equal(t, wantOutput, gotOutput, "expected the color from the config file and the number from the flag default")
	})

	// Set favorite-color with an environment variable
	t.Run("env var", func(t *testing.T) {
		os.Setenv("CCS_CONTAINER_REGISTRY_USERNAME", "sarouman")
		defer os.Unsetenv("CCS_CONTAINER_REGISTRY_USERNAME")

		output := new(bytes.Buffer)
		cmd := NewRootCommand(output)
		cmd.SetArgs([]string{"--dry-run"})
		cmd.SetOut(output)
		cmd.Execute()
		klog.Flush()

		gotOutput := output.String()
		wantOutput := `Successful login of sarouman to https://index.docker.io/v1`
		assert.Contains(t, gotOutput, wantOutput)
		// assert.Equal(t, wantOutput, gotOutput, "expected the color to use the environment variable value and the number to use the flag default")
	})

	// Set number with a flag
	t.Run("flag", func(t *testing.T) {
		// Run ./stingoftheviper --number 2
		output := &bytes.Buffer{}
		cmd := NewRootCommand(output)
		cmd.SetOut(output)
		cmd.SetArgs([]string{"--container-registry-username", "sauron", "--dry-run"})
		cmd.Execute()
		klog.Flush()
		gotOutput := output.String()
		wantOutput := `Successful login of sauron to https://index.docker.io/v1`
		assert.Contains(t, gotOutput, wantOutput)

		// assert.Equal(t, wantOutput, gotOutput, "expected the number to use the flag value and the color to use the flag default")
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
