// Copyright 2023 codestation. All rights reserved.
// Use of this source code is governed by a MIT-license
// that can be found in the LICENSE file.

package s3

import (
	"context"
	"fmt"
	"io"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

type Client interface {
	Upload(ctx context.Context, key string, r io.Reader) (string, error)
	Download(ctx context.Context, key string, w io.WriterAt) (int64, error)
	ListObjects(ctx context.Context, opts ListOptions) (*ListOutput, error)
}

type ClientImpl struct {
	uploader   *s3manager.Uploader
	downloader *s3manager.Downloader
	svc        *s3.S3
	bucket     string
}

type ClientOpts struct {
	// Endpoint Optional hostname or URI for the S3 service
	Endpoint string
	// Region of the S3 service
	Region string
	// Bucket name where the files are read/written to
	Bucket string
	// AccessKey credential for the S3 service
	AccessKey string
	// SecretKey credential for the S3 service
	SecretKey string
	// ForcePathStyle set to `true` to force the request to use path-style addressing
	ForcePathStyle bool
}

// NewClient creates a new SFTP client
func NewClient(opts ClientOpts) (*ClientImpl, error) {
	sess, err := session.NewSession(&aws.Config{
		Credentials:      credentials.NewStaticCredentials(opts.AccessKey, opts.SecretKey, ""),
		Endpoint:         aws.String(opts.Endpoint),
		Region:           aws.String(opts.Region),
		S3ForcePathStyle: aws.Bool(opts.ForcePathStyle),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create session: %w", err)
	}

	// Create an uploader with the session and default options
	uploader := s3manager.NewUploader(sess)

	// Create a downloader with the session and default options
	downloader := s3manager.NewDownloader(sess)

	svc := s3.New(sess)

	return &ClientImpl{
		uploader:   uploader,
		downloader: downloader,
		svc:        svc,
		bucket:     opts.Bucket,
	}, nil
}

func (c *ClientImpl) Upload(ctx context.Context, key string, r io.Reader) (string, error) {
	// Upload the file to S3.
	result, err := c.uploader.UploadWithContext(ctx, &s3manager.UploadInput{
		Bucket: aws.String(c.bucket),
		Key:    aws.String(key),
		Body:   r,
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload file, %v", err)
	}

	return result.Location, nil
}

func (c *ClientImpl) Download(ctx context.Context, key string, w io.WriterAt) (int64, error) {
	// Write the contents of S3 Object to the file
	n, err := c.downloader.DownloadWithContext(ctx, w, &s3.GetObjectInput{
		Bucket: aws.String(c.bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return 0, fmt.Errorf("failed to download file: %w", err)
	}

	return n, nil
}

type ListItem struct {
	// Key is the name of the object
	Key string
	// Size in bytes of the object
	Size int64
}

type ListOutput struct {
	// Items Metadata about each object returned
	Items []ListItem
	// IsTruncated indicates if there are more results to be returned
	IsTruncated bool
	// ContinuationToken is set when IsTruncated is true
	ContinuationToken string
}

type ListOptions struct {
	// Prefix limits the response to keys that starts with the specified prefix
	Prefix string
	// Set the maximum number of keys to be returned. By default, the action
	// returns up to 1000 keys but never more
	MaxKeys int64
	// ContinuationToken is used to continue the list operation in the bucket
	ContinuationToken string
}

func (c *ClientImpl) ListObjects(ctx context.Context, opts ListOptions) (*ListOutput, error) {
	input := &s3.ListObjectsV2Input{
		Bucket: aws.String(c.bucket),
		Prefix: aws.String(opts.Prefix),
	}

	if opts.MaxKeys > 0 && opts.MaxKeys <= 1000 {
		input.MaxKeys = aws.Int64(opts.MaxKeys)
	}

	if opts.ContinuationToken != "" {
		input.ContinuationToken = aws.String(opts.ContinuationToken)
	}

	result, err := c.svc.ListObjectsV2WithContext(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to list objects: %w", err)
	}

	output := ListOutput{
		Items:             make([]ListItem, 0, len(result.Contents)),
		IsTruncated:       aws.BoolValue(result.IsTruncated),
		ContinuationToken: aws.StringValue(result.NextContinuationToken),
	}

	for _, item := range result.Contents {
		output.Items = append(output.Items, ListItem{
			Key:  aws.StringValue(item.Key),
			Size: aws.Int64Value(item.Size),
		})
	}

	return &output, nil
}
