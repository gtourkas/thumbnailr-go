module thumbnailr/host_lambda

go 1.12

require (
	github.com/aws/aws-lambda-go v1.13.2
	github.com/aws/aws-sdk-go v1.25.27
	github.com/awslabs/aws-lambda-go-api-proxy v0.5.0
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/golang/protobuf v1.3.2 // indirect
	github.com/mattn/go-isatty v0.0.9 // indirect
	github.com/pkg/errors v0.8.1
	golang.org/x/net v0.0.0-20191105084925-a882066a44e0 // indirect
	gopkg.in/oauth2.v3 v3.12.0
	thumbnailr/app v0.0.0
	thumbnailr/bus_sns v0.0.0
	thumbnailr/repos_dynamodb v0.0.0
	thumbnailr/stores_s3 v0.0.0
)

replace thumbnailr/app v0.0.0 => ../../app

replace thumbnailr/stores_s3 v0.0.0 => ../../stores_s3

replace thumbnailr/repos_dynamodb v0.0.0 => ../../repos_dynamodb

replace thumbnailr/bus_sns v0.0.0 => ../../bus_sns
