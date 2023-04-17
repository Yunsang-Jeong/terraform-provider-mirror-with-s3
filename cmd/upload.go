package cmd

import (
	"github.com/Yunsang-Jeong/terraform-provider-mirror-with-s3/internal/uploader"
	"github.com/spf13/cobra"
)

type UploadCmd struct{}

var uploadCmdFlags = map[string]flag{
	"Bukcet": {
		_type:       "string",
		shorten:     "b",
		description: "[req] The name of AWS S3 bucket to search terraform provider",
		requirement: true,
	},
	"Providers": {
		_type:       "list",
		shorten:     "p",
		description: "[req] The list of terraform provider to upload",
		requirement: true,
	},
	"OS": {
		_type:       "string",
		shorten:     "o",
		description: "[req] The os-type in the terraform CLI runtime",
		requirement: true,
	},
	"Architecture": {
		_type:       "string",
		shorten:     "a",
		description: "[req] The architecture-type in the terraform CLI runtime",
		requirement: true,
	},
}

func (r *UploadCmd) Init() *cobra.Command {
	c := &cobra.Command{
		Use:   "upload",
		Short: "Upload the latest version of A to AWS S3 bucket",
		RunE: func(cmd *cobra.Command, args []string) error {
			bucket, _ := cmd.Flags().GetString("Bukcet")
			providers, _ := cmd.Flags().GetStringSlice("Providers")
			os, _ := cmd.Flags().GetString("OS")
			architecture, _ := cmd.Flags().GetString("Architecture")

			s := uploader.NewUploader(bucket, providers, os, architecture)
			if err := s.Start(); err != nil {
				return err
			}

			return nil
		},
	}

	cobraFlagRegister(c, uploadCmdFlags)

	return c
}
