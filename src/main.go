package main

import (
	"errors"
	"fmt"

	"os"
	"path"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/mkideal/cli"
)

type argT struct {
	cli.Helper
	S3Region string `cli:"s3_region" usage:"AWS S3 region" dft:"eu-central-1"`
	S3Endpoint string `cli:"s3_endpoint" usage:"AWS S3 endpoint"`

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
		Region: aws.String(args.S3Region),
		Endpoint: aws.String(args.S3Endpoint),
		Credentials: credentials.AnonymousCredentials,
	})
	if err != nil {
		return errors.New(fmt.Sprintf("Failed to create session: %v", err))
	}

	downloader := s3manager.NewDownloader(sess)

	files := []struct {
		Bucket       string
		Path         string
		IsCompressed bool
	}{
		{ Bucket: args.ModelBucket, Path: args.ModelPath, IsCompressed: true },
		{ Bucket: args.UserInputBucket, Path: args.UserInputPath, IsCompressed: false },
	}

	for _, f := range files {
		if f.Bucket == "" || f.Path == "" {
			continue
		}
		fileName := path.Join(args.OutputDir, f.Path)
		outputFile, err := os.Create(fileName)
		if err != nil {
			return err
		}
		if err := download(downloader, f.Bucket, f.Path, outputFile); err != nil {
			return err
		}
		if f.IsCompressed {
			if err := Untar(outputFile, args.OutputDir); err != nil {
				return err
			}
			if err := os.Remove(fileName); err != nil {
				return err
			}
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