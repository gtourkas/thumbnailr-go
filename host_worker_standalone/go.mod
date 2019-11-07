module host_worker_standalone

go 1.12

require (
	thumbnailr/app v0.0.0
	thumbnailr/repos_dynamodb v0.0.0
	thumbnailr/stores_s3 v0.0.0
)

replace thumbnailr/app v0.0.0 => ../app

replace thumbnailr/stores_s3 v0.0.0 => ../stores_s3

replace thumbnailr/repos_dynamodb v0.0.0 => ../repos_dynamodb
