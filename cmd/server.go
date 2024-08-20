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
	"CertFile": {
		_type:       "string",
		shorten:     "c",
		description: "[req] The https certificate file",
		requirement: true,
	},
	"KeyFile": {
		_type:       "string",
		shorten:     "k",
		description: "[req] The https private-key file",
		requirement: true,
	},
}

func (r *ServerCmd) Init() *cobra.Command {
	c := &cobra.Command{
		Use:   "server",
		Short: "Run the server serving the terraform providers",
		RunE: func(cmd *cobra.Command, args []string) error {
			flags := cmd.Flags()

			bucketName, _ := flags.GetString("BukcetName")
			certFile, _ := flags.GetString("CertFile")
			keyFile, _ := flags.GetString("KeyFile")

			s := server.NewProviderMirrorServer(bucketName, certFile, keyFile)
			if err := s.Start(); err != nil {
				return err
			}

			return nil
		},
	}

	cobraFlagRegister(c, serverCmdFlags)

	return c
}
