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

func cobraFlagRegister(c *cobra.Command, flags map[string]flag) {
	for name, flag := range flags {

		switch flag._type {
		case "string":
			if flag.defaultValue != nil {
				c.Flags().StringP(
					name,
					flag.shorten,
					flag.defaultValue.(string),
					flag.description,
				)
			} else {
				c.Flags().StringP(
					name,
					flag.shorten,
					"",
					flag.description,
				)
			}
			if flag.requirement {
				c.MarkFlagRequired(name)
			}

		case "list":
			if flag.defaultValue != nil {
				c.Flags().StringSliceP(
					name,
					flag.shorten,
					flag.defaultValue.([]string),
					flag.description,
				)
			} else {
				c.Flags().StringSliceP(
					name,
					flag.shorten,
					[]string{},
					flag.description,
				)
			}
			if flag.requirement {
				c.MarkFlagRequired(name)
			}
		default:
			panic(fmt.Sprintf("unsupported flag type: %s", flag._type))
		}
	}
}
