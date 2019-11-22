#!/usr/bin/env bash

ACCOUNTID=$(aws sts get-caller-identity | jq -r '. | .Account')
if [[ -z "$ACCOUNTID" ]]; then
  exit 1
fi
sam package --template-file template.yaml --s3-bucket "sam-deployment-$ACCOUNTID" --output-template-file packaged.yaml
rc=$?; if [[ $rc != 0 ]]; then exit $rc; fi
sam deploy --stack-name thumbnailr --template-file packaged.yaml --capabilities CAPABILITY_IAM