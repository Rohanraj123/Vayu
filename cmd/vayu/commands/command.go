package commands

import (
	vayuserver "github.com/Rohanraj123/vayu/cmd/vayu/commands/listen"
	"github.com/spf13/cobra"
)

var (
	example = `
			# Start server with defaults
			vayu run

			# Start with TLS
			vayu run --addr :8443 --tls-cert cert.pem --tls-key key.pem
`
)

func RootCommand() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:          "vayu",
		Short:        "vayu - High Performance Api Gateway",
		Long:         "vayu - High Performance Api Gateway built in Go for cloud native envrionments",
		Example:      example,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return cmd.Help()
		},
	}

	rootCmd.AddCommand(vayuserver.Command())

	return rootCmd
}
