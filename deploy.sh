#!/bin/bash

rm -Rf lambda-update-sqs-queue-depth.zip main

BUCKET_NAME=public-aws-serverless-repo
GOOS=linux go build main.go

zip lambda-update-sqs-queue-depth.zip ./main

aws s3 cp lambda-update-sqs-queue-depth.zip s3://$BUCKET_NAME/lambda-update-sqs-queue-depth.zip

aws s3api put-object-tagging --bucket $BUCKET_NAME --key lambda-update-sqs-queue-depth.zip --tagging 'TagSet={Key=public,Value=yes}'
