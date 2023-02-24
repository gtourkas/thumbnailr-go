module thumbnailr/stores_s3

go 1.12

require (
	github.com/aws/aws-sdk-go v1.25.19
	github.com/pkg/errors v0.8.1
	github.com/stretchr/testify v1.4.0 // indirect
	golang.org/x/net v0.7.0 // indirect
)

replace thumbnailr/app v0.0.0 => ../app
