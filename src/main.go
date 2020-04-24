package main

import (
	"errors"
	"fmt"
	"os"
	"path"

	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/mkideal/cli"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
)

type argT struct {
	cli.Helper
	OutputDir string `cli:"output_dir" usage:"output directory" dft:"/neuroide"`

	ModelBucket string `cli:"model_bucket" usage:"model s3 bucket"`
	ModelPath string `cli:"model_path" usage:"model s3 path"`

	UserInputBucket string `cli:"user_input_bucket" usage:"user-input s3 bucket"`
	UserInputPath string `cli:"user_input_path" usage:"user-input s3 path"`
}

func download(downloader *s3manager.Downloader, bucket, item string, outputFile *os.File) error {
	numBytes, err := downloader.Download(outputFile,
		&s3.GetObjectInput{
			Bucket: aws.String(bucket),
			Key:    aws.String(item),
		})
	if err != nil {
		return errors.New(fmt.Sprintf("Unable to download item %q, %v", item, err))
	}

	fmt.Println("Downloaded", outputFile.Name(), numBytes, "bytes")
	return nil
}

func Download(args *argT) error {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("eu-central-1"),
		Credentials: credentials.NewEnvCredentials(),
	})
	if err != nil {
		return err
	}
	downloader := s3manager.NewDownloader(sess)

	files := []struct {
		Bucket string
		Path string
	}{
		{ Bucket: args.ModelBucket, Path: args.ModelPath },
		{ Bucket: args.UserInputBucket, Path: args.UserInputPath },
	}

	for _, f := range files {
		if f.Bucket == "" || f.Path == "" {
			continue
		}
		outputFile, err := os.Create(path.Join(args.OutputDir, f.Path))
		if err != nil {
			return err
		}
		if err := download(downloader, f.Bucket, f.Path, outputFile); err != nil {
			return err
		}
	}

	return nil
}

func main() {
	os.Exit(cli.Run(new(argT), func(ctx *cli.Context) error {
		argv := ctx.Argv().(*argT)
		return Download(argv)
	}))
}