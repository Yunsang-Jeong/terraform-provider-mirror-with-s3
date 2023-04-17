package main

import "github.com/Yunsang-Jeong/terraform-provider-mirror-with-s3/cmd"

func main() {
	run := &cmd.RunCmd{}
	upload := &cmd.UploadCmd{}

	cmd.RootCmd.AddCommand(run.Init())
	cmd.RootCmd.AddCommand(upload.Init())
	cmd.Execute()
}
