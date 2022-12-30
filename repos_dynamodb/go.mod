module thumbnailr/repos_dynamodb

go 1.12

require (
	github.com/aws/aws-sdk-go v1.33.0
	github.com/pkg/errors v0.9.1
	thumbnailr/app v0.0.0
)

replace thumbnailr/app v0.0.0 => ../app
