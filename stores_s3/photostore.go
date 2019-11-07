package stores_s3

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/pkg/errors"
	"io"
)
import "github.com/aws/aws-sdk-go/service/s3"

type PhotoStore struct {
	Bucket string

	downloader *s3manager.Downloader
}

func NewPhotoStore(sess *session.Session, bucket string) *PhotoStore {
	return &PhotoStore{
		Bucket: bucket,

		downloader: s3manager.NewDownloader(sess),
	}
}

func (ps *PhotoStore) Get(photoID string, writer io.Writer) error {

	var buf []byte
	waBuf := aws.NewWriteAtBuffer(buf)

	if _, err := ps.downloader.Download(waBuf, &s3.GetObjectInput{
		Bucket: aws.String(ps.Bucket),
		Key:    aws.String(photoID),
	}); err != nil {
		return errors.Wrapf(err, "cannot download photo %s", photoID)
	}

	if _, err := writer.Write(waBuf.Bytes()); err != nil {
		return errors.Wrapf(err, "cannot write the contents of photo %s to the supplied writer", photoID)
	}

	return nil
}
