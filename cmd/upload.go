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
		description: "[req] The Bucket Name for uploading the terraform provider",
		requirement: true,
	},
	"Hostname": {
		_type:       "string",
		shorten:     "r",
		description: "[opt] The Hostname of the Terraform Provider Address",
		requirement: false,
		defaultValue: "registry.terraform.io",
	},
	"Namespace": {
		_type:       "string",
		shorten:     "n",
		description: "[req] The Namespace of the Terraform Provider Address",
		requirement: true,
	},
	"Type": {
		_type:       "string",
		shorten:     "t",
		description: "[req] The Type of the Terraform Provider Address",
		requirement: true,
	},
	"OS": {
		_type:       "string",
		shorten:     "o",
		description: "[req] The OS type used in the terraform CLI runtime",
		requirement: true,
	},
	"Arch": {
		_type:       "string",
		shorten:     "a",
		description: "[req] The Architecture type used in the terraform CLI runtime",
		requirement: true,
	},
	"Version": {
		_type:       "string",
		shorten:     "v",
		description: "[opt] The version of the Terraform Provider to upload",
		requirement: false,
		defaultValue: "",
	},
}

func (r *UploadCmd) Init() *cobra.Command {
	c := &cobra.Command{
		Use:   "upload",
		Short: "Upload the latest version of A to AWS S3 bucket",
		RunE: func(cmd *cobra.Command, args []string) error {
			bucket, _ := cmd.Flags().GetString("Bukcet")
			providerHostname, _ := cmd.Flags().GetString("Hostname")
			providerNamespace, _ := cmd.Flags().GetString("Namespace")
			providerType, _ := cmd.Flags().GetString("Type")
			providerOS, _ := cmd.Flags().GetString("OS")
			providerArch, _ := cmd.Flags().GetString("Arch")
			providerVersion, _ := cmd.Flags().GetString("Version")

			s := uploader.NewUploader(bucket, providerHostname, providerNamespace, providerType, providerOS, providerArch, providerVersion)
			if err := s.Start(); err != nil {
				return err
			}

			return nil
		},
	}

	cobraFlagRegister(c, uploadCmdFlags)

	return c
}
