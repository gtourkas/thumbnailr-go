#!/usr/bin/env bash

ACCOUNTID=`aws sts get-caller-identity | jq -r '. | .Account'`
if [[ -z "$ACCOUNTID" ]]; then
  exit 1
fi
aws s3 sync . "s3://thumbnailr-photostore-$ACCOUNTID" --exclude "*.sh"