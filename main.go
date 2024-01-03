package main

import "github.com/Yunsang-Jeong/terraform-provider-mirror-with-s3/cmd"

func main() {
	server := &cmd.ServerCmd{}
	upload := &cmd.UploadCmd{}

	cmd.RootCmd.AddCommand(server.Init())
	cmd.RootCmd.AddCommand(upload.Init())
	cmd.Execute()
}
