module thumbnailr/host_api_standalone

go 1.12

require (
	github.com/aws/aws-sdk-go v1.33.0
	github.com/pkg/errors v0.9.1
	thumbnailr/app v0.0.0
	thumbnailr/bus_sns v0.0.0
	thumbnailr/repos_dynamodb v0.0.0
	thumbnailr/stores_s3 v0.0.0
)

replace thumbnailr/app v0.0.0 => ../app

replace thumbnailr/stores_s3 v0.0.0 => ../stores_s3

replace thumbnailr/repos_dynamodb v0.0.0 => ../repos_dynamodb

replace thumbnailr/bus_sns v0.0.0 => ../bus_sns
