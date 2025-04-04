package cmd

import "github.com/spf13/cobra"

// A cli for output all environment variables
// If build into app.exe then app outenv
// Run runs the actual command
var outEnvCmd = &cobra.Command{
    Use:   "outenv",
    Short: "Output all environment variables to std",
    Run: func(cmd *cobra.Command, args []string) {
        newService().OutEnv()
    },
}
