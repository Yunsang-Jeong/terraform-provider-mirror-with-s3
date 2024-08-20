package main

import "github.com/Yunsang-Jeong/terraform-provider-mirror-with-s3/cmd"

func main() {
	server := &cmd.ServerCmd{}

	cmd.RootCmd.AddCommand(server.Init())

	cmd.Execute()
}
