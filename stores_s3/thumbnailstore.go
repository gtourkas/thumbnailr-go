package stores_s3

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/pkg/errors"
	"io"
)

type ThumbnailStore struct {
	Bucket string

	uploader *s3manager.Uploader
	downloader *s3manager.Downloader
}

func NewThumbnailStore(sess *session.Session, bucket string) *ThumbnailStore {
	return &ThumbnailStore{
		Bucket: bucket,

		uploader: s3manager.NewUploader(sess),
		downloader: s3manager.NewDownloader(sess),
	}
}

func (ts *ThumbnailStore) getKey(userID string, id string) string {
	return fmt.Sprintf("%s-%s", userID, id)
}

func (ts *ThumbnailStore) Put(userID string, id string, reader io.Reader) error {
	key := ts.getKey(userID, id)

	if _, err := ts.uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(ts.Bucket),
		Key:    aws.String(key),
		Body:   reader,
	}); err != nil {
		return errors.Wrapf(err, "cannot upload thumbnail %s of user %s", userID, id)
	}
	return nil
}

func (ts *ThumbnailStore) Get(userID string, id string, writer io.Writer) error {
	key := ts.getKey(userID, id)

	var buf []byte
	waBuf := aws.NewWriteAtBuffer(buf)

	if _, err := ts.downloader.Download(waBuf, &s3.GetObjectInput{
		Bucket: aws.String(ts.Bucket),
		Key:    aws.String(key),
	}); err != nil {
		return errors.Wrapf(err, "cannot download thumbnail %s of user %s", userID, id)
	}

	if _, err := writer.Write(waBuf.Bytes()); err != nil {
		return errors.Wrapf(err, "cannot write the contents of thumbnail %s of user %s to the supplied writer", userID, id)
	}

	return nil
}