package main

import (
	"fmt"
	"os"
	"strings"

	"ccs-build.thephoenixhomelab.com/services"
	"github.com/go-sanitize/sanitize"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	klog "k8s.io/klog/v2"
)

// parse flags
const (
	defaultConfigFilename = "ccs"
	envPrefix             = "CCS"
)

type environmentRegistry struct {
	Username string `san:"trim,max=64"`
	Password string `san:"trim,max=256"`
	Registry string `san:"trim"`
}

type config struct {
	envCreds environmentRegistry
	dryRun   bool
}

func main() {
	cmd := NewRootCommand()
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func NewRootCommand() *cobra.Command {

	config := config{
		envCreds: environmentRegistry{},
		dryRun:   false,
	}
	// klog.InitFlags(nil)

	// "For Frodo." - Aragorn II
	rootCmd := &cobra.Command{
		Use:   "build",
		Short: "CICD Standard build task implementation.",
		Long:  `CICD Standard build task implementation with docker as container runtime engine.`,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			// You can bind cobra and viper in a few locations, but PersistencePreRunE on the root command works well
			return initializeConfig(cmd)
		},
		Run: func(cmd *cobra.Command, args []string) {
			if err := SanitizeInputs(&config.envCreds); err != nil {
				klog.Errorf("CLI argument error: %w", err)
				panic(config.envCreds)
			}
			// container service should take arg and have docker as default
			ctnrSvc, err := newCntrSvc(config.dryRun, "")
			if err != nil {
				klog.Errorf("error initializing Container service: %w", err)
				panic("error initializing Container service")
			}
			// Working with OutOrStdout/OutOrStderr allows us to unit test our command easier
			klog.SetOutput(cmd.OutOrStdout())

			if config.envCreds.Username != "" {
				ctnrSvc.Login(config.envCreds.Username, config.envCreds.Password, config.envCreds.Registry)
			} else {
				klog.Info("no container registry")
			}
		},
	}

	rootCmd.Flags().StringVarP(&config.envCreds.Username, "container-registry-username", "u", "galadriel",
		"the username to log into the container registry")
	rootCmd.Flags().StringVarP(&config.envCreds.Password, "container-registry-password", "p", "indis",
		"the password to log into the container registry")
	rootCmd.Flags().StringVarP(&config.envCreds.Registry, "container-registry", "r", "https://index.docker.io/v1",
		"the password to log into the container registry")
	rootCmd.Flags().BoolVar(&config.dryRun, "dry-run", false,
		"Execute the command in dry-run mode")
	return rootCmd
}

func newCntrSvc(dryRun bool, engine string) (services.CntrSvcI, error) {
	if len(engine) == 0 {
		engine = "docker"
	}
	if dryRun {
		engine = "dry-run"
	}
	ctnrSvc, err := services.NewCntrSvc(engine)

	return ctnrSvc, err
}

func SanitizeInputs(envCreds *environmentRegistry) (err error) {
	sanitzr, err := sanitize.New()
	if err != nil {
		return err
	}

	err = sanitzr.Sanitize(envCreds)
	return err
}

func initializeConfig(cmd *cobra.Command) error {
	v := viper.New()

	// Set the base name of the config file, without the file extension.
	v.SetConfigName(defaultConfigFilename)

	v.AddConfigPath(".")

	// Attempt to read the config file, gracefully ignoring errors
	// caused by a config file not being found. Return an error
	// if we cannot parse the config file.
	if err := v.ReadInConfig(); err != nil {
		// It's okay if there isn't a config file
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return err
		}
	}

	// When we bind flags to environment variables expect that the
	// environment variables are prefixed, e.g. a flag like --number
	// binds to an environment variable STING_NUMBER. This helps
	// avoid conflicts.
	v.SetEnvPrefix(envPrefix)

	// Bind to environment variables
	// Works great for simple config names, but needs help for names
	// like --favorite-color which we fix in the bindFlags function
	v.AutomaticEnv()

	// Bind the current command's flags to viper
	bindFlags(cmd, v)
	// v.BindPFlags(cmd.Flags())

	return nil
}

// Bind each cobra flag to its associated viper configuration (config file and environment variable)
func bindFlags(cmd *cobra.Command, v *viper.Viper) {
	cmd.Flags().VisitAll(func(f *pflag.Flag) {
		// Environment variables can't have dashes in them, so bind them to their equivalent
		// keys with underscores, e.g. --favorite-color to STING_FAVORITE_COLOR
		if strings.Contains(f.Name, "-") {
			envVarSuffix := strings.ToUpper(strings.ReplaceAll(f.Name, "-", "_"))
			v.BindEnv(f.Name, fmt.Sprintf("%s_%s", envPrefix, envVarSuffix))
		}

		// Apply the viper config value to the flag when the flag is not set and viper has a value
		if !f.Changed && v.IsSet(f.Name) {
			val := v.Get(f.Name)
			cmd.Flags().Set(f.Name, fmt.Sprintf("%v", val))
		}
	})
}

// launch docker container
