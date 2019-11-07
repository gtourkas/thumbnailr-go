module thumbnailr/stores_s3

go 1.12

require (
	github.com/aws/aws-sdk-go v1.25.19
	github.com/stretchr/testify v1.4.0 // indirect
	golang.org/x/net v0.0.0-20191021144547-ec77196f6094 // indirect
	thumbnailr/app v0.0.0
)

replace thumbnailr/app v0.0.0 => ../app
