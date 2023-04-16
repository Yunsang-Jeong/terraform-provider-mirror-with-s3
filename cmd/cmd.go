package cmd

type stringFlag struct {
	shorten      string
	defaultValue string
	description  string
	requirement  bool
}

const (
	BucketName = "bucket-name"
)
