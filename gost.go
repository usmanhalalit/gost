package gost

import "github.com/usmanhalalit/gost/adapter"
import "github.com/usmanhalalit/gost/adapter/s3"

func New() adapter.Directory {
	fs := s3.New(s3.Config{
		Id: "AKIAJBRFB4PEZIKTETJQ",
		Secret: "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",
		Region: "us-east-1",
		Bucket: "usman-gost",
	})

	return fs
}
