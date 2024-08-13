package cmd

import (
	_ "embed"
	"fmt"
	"os"
	"strings"

	"buffalo-cli/internal/cmd/build"
	"buffalo-cli/internal/cmd/destroy"
	"buffalo-cli/internal/cmd/dev"
	"buffalo-cli/internal/cmd/fix"
	"buffalo-cli/internal/cmd/generate"
	"buffalo-cli/internal/cmd/info"
	"buffalo-cli/internal/cmd/new"
	cmdplugins "buffalo-cli/internal/cmd/plugins"
	"buffalo-cli/internal/cmd/routes"
	"buffalo-cli/internal/cmd/setup"
	"buffalo-cli/internal/cmd/task"
	"buffalo-cli/internal/cmd/test"
	"buffalo-cli/internal/cmd/version"
	"buffalo-cli/internal/plugins"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

const (
	dbNotFound  = `unknown command "db"`
	popNotFound = `unknown command "pop"`
)

var (
	anywhereCommands = []string{"completion", "help", "info", "new", "version"}

	//go:embed popinstructions.txt
	popInstallInstructions string
)

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	rootCmd := cmd()
	err := rootCmd.Execute()
	if err == nil {
		return
	}

	if strings.Contains(err.Error(), dbNotFound) || strings.Contains(err.Error(), popNotFound) {
		logrus.Errorf(popInstallInstructions)
		os.Exit(-1)
	}
	logrus.Errorf("Error: %s", err)
	if strings.Contains(err.Error(), dbNotFound) || strings.Contains(err.Error(), popNotFound) {
		fmt.Println(popInstallInstructions)
		os.Exit(-1)
	}
	os.Exit(-1)
}

// preRunCheck confirms that the command being called can be
// invoked within the current folder. Some commands like generators
// require to be called from the root of the project.
func preRunCheck(cmd *cobra.Command, args []string) error {
	if err := plugins.Load(); err != nil {
		return err
	}

	cmdName := cmd.Name()
	if cmd.Parent() != cmd.Root() {
		// e.g. `completion` for `buffalo completion bash` instead of `bash`
		cmdName = cmd.Parent().Name()
	}

	for _, freeCmd := range anywhereCommands {
		if freeCmd == cmdName {
			return nil
		}
	}

	// check if command is being run from inside a buffalo project
	_, err := os.Stat(".buffalo.dev.yml")
	if err != nil {
		return fmt.Errorf("you need to be inside your buffalo project path to run this command")
	}

	return nil
}

func cmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:               "buffalo",
		Short:             "Build Buffalo applications with ease",
		PersistentPreRunE: preRunCheck,
		TraverseChildren:  true,
		SilenceErrors:     true,
	}

	newCmd := new.Cmd()
	setupCmd := setup.Cmd()
	generateCmd := generate.Cmd()
	destroyCmd := destroy.Cmd()
	versionCmd := version.Cmd()
	testCmd := test.Cmd()
	devCmd := dev.Cmd()
	taskCmd := task.Cmd()
	routesCmd := routes.Cmd()

	decorate("new", newCmd)
	decorate("info", rootCmd)
	decorate("fix", rootCmd)
	decorate("update", rootCmd)
	decorate("setup", setupCmd)
	decorate("generate", generateCmd)
	decorate("destroy", destroyCmd)
	decorate("version", versionCmd)
	decorate("test", testCmd)
	decorate("dev", devCmd)
	decorate("task", taskCmd)
	decorate("routes", routesCmd)

	rootCmd.AddCommand(newCmd)
	rootCmd.AddCommand(build.Cmd())
	rootCmd.AddCommand(cmdplugins.PluginsCmd)
	rootCmd.AddCommand(info.Cmd())
	rootCmd.AddCommand(setupCmd)
	rootCmd.AddCommand(generateCmd)
	rootCmd.AddCommand(fix.Cmd())
	rootCmd.AddCommand(destroyCmd)
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(testCmd)
	rootCmd.AddCommand(devCmd)
	rootCmd.AddCommand(taskCmd)
	rootCmd.AddCommand(routesCmd)

	decorate("root", rootCmd)

	return rootCmd
}
