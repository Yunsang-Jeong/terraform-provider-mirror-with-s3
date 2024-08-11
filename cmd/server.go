package cmd

import (
	"github.com/Yunsang-Jeong/terraform-provider-mirror-with-s3/internal/server"
	"github.com/spf13/cobra"
)

type ServerCmd struct{}

var serverCmdFlags = map[string]flag{
	"BukcetName": {
		_type:       "string",
		shorten:     "b",
		description: "[req] The name of AWS S3 bucket to search terraform provider",
		requirement: true,
	},
}

func (r *ServerCmd) Init() *cobra.Command {
	c := &cobra.Command{
		Use:   "server",
		Short: "Run the server serving the terraform providers",
		RunE: func(cmd *cobra.Command, args []string) error {
			bucketName, _ := cmd.Flags().GetString("BukcetName")

			s := server.NewProviderMirrorServer(bucketName)
			if err := s.Start(); err != nil {
				return err
			}

			return nil
		},
	}

	cobraFlagRegister(c, serverCmdFlags)

	return c
}
