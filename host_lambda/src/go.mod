module thumbnailr/host_lambda

go 1.12

require (
	github.com/aws/aws-lambda-go v1.13.2
	github.com/aws/aws-sdk-go v1.25.26
	github.com/pkg/errors v0.8.1
	golang.org/x/net v0.0.0-20191105084925-a882066a44e0 // indirect
	thumbnailr/app v0.0.0
	thumbnailr/repos_dynamodb v0.0.0
	thumbnailr/stores_s3 v0.0.0
)

replace thumbnailr/app v0.0.0 => ../../app

replace thumbnailr/stores_s3 v0.0.0 => ../../stores_s3

replace thumbnailr/repos_dynamodb v0.0.0 => ../../repos_dynamodb
