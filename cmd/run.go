package cmd

import (
	"github.com/Yunsang-Jeong/terraform-provider-mirror-with-s3/internal/server"
	"github.com/spf13/cobra"
)

type RunCmd struct{}

var runStringFlags = map[string]flag{
	"Bukcet": {
		shorten:     "b",
		description: "[req] The name of AWS S3 bucket to search terraform provider",
		requirement: true,
	},
}

func (r *RunCmd) Init() *cobra.Command {
	c := &cobra.Command{
		Use:   "run",
		Short: "Run terraform provider mirror",
		RunE: func(cmd *cobra.Command, args []string) error {
			bucket, _ := cmd.Flags().GetString("Bukcet")

			s := server.NewProviderMirrorServer(bucket, "server.crt", "ca.key", true, true)
			if err := s.Start(); err != nil {
				return err
			}

			return nil
		},
	}

	cobraFlagRegister(c, runStringFlags)

	return c
}
