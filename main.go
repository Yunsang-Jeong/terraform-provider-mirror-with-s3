package main

import "github.com/Yunsang-Jeong/terraform-provider-mirror-with-s3/cmd"

func main() {
	run := &cmd.RunCmd{}

	cmd.RootCmd.AddCommand(run.Init())
	cmd.Execute()
}
