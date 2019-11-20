package repos_dynamodb

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/pkg/errors"
	"thumbnailr/app"
)

type ThumbnailRepo struct {
	db        *dynamodb.DynamoDB
	tableName *string
}

func NewThumbnailRepo(sess *session.Session) *ThumbnailRepo {
	return &ThumbnailRepo{
		db:        dynamodb.New(sess),
		tableName: aws.String("thumbnailr-thumbnails"),
	}
}

func (tr *ThumbnailRepo) Get(id string) (*app.Thumbnail, error) {
	input := &dynamodb.GetItemInput{
		TableName: tr.tableName,
		Key: map[string]*dynamodb.AttributeValue{
			"ID": {
				S: aws.String(id),
			},
		},
	}

	tmp := app.Thumbnail{}
	if res, err := tr.db.GetItem(input); err == nil {
		if res.Item == nil {
			return nil, nil
		}
		if err := dynamodbattribute.UnmarshalMap(res.Item, &tmp); err != nil {
			return nil, errors.Wrapf(err, "cannot unmarshal thumbnail %s", id)
		}
		return &tmp, nil
	} else {
		return nil, errors.Wrapf(err, "cannot get thumbnail %s", id)
	}
}

func (tr *ThumbnailRepo) Save(thumbnail *app.Thumbnail) error {
	item, e := dynamodbattribute.MarshalMap(thumbnail)
	if e != nil {
		return e
	}

	input := &dynamodb.PutItemInput{
		TableName: tr.tableName,
		Item:      item,
	}

	if _, err := tr.db.PutItem(input); err != nil {
		return errors.Wrapf(err, "cannot save thumbnail %s", thumbnail.ID)
	}

	return nil
}
