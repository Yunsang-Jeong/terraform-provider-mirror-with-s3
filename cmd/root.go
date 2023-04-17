package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

type flag struct {
	_type        string
	shorten      string
	defaultValue interface{}
	description  string
	requirement  bool
}

var RootCmd = &cobra.Command{
	Use:           "tpm",
	Short:         "Terraform provider mirror with s3",
	Long:          "Terraform provider mirror with s3",
	SilenceErrors: true,
}

func Execute() {
	// Remove help for root command
	RootCmd.SetHelpCommand(&cobra.Command{Hidden: true})

	// Remove shell completion
	RootCmd.CompletionOptions = cobra.CompletionOptions{
		DisableDefaultCmd:   true,
		DisableNoDescFlag:   true,
		DisableDescriptions: true,
		HiddenDefaultCmd:    true,
	}

	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
	}
}

func cobraFlagRegister(c *cobra.Command, flags map[string]flag) {
	for name, flag := range flags {

		switch flag._type {
		case "string":
			c.Flags().StringP(
				name,
				flag.shorten,
				flag.defaultValue.(string),
				flag.description,
			)
			if flag.requirement {
				c.MarkFlagRequired(name)
			}
		case "list":
			c.Flags().StringSliceP(
				name,
				flag.shorten,
				flag.defaultValue.([]string),
				flag.description,
			)
			if flag.requirement {
				c.MarkFlagRequired(name)
			}
		default:
			panic(fmt.Sprintf("unsupported flag type: %s", flag._type))
		}
	}
}
