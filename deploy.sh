#!/bin/bash

rm -Rf lambda-update-sqs-queue-depth.zip main

BUCKET_PREFIX=pub-cfn-cust-res-pocs
GOOS=linux go build main.go

REGIONS=$(aws ec2 describe-regions --output text --query Regions[*].RegionName)

# Currently, StackSets are not supported in these regions. This may change over time.
REGIONS=${REGIONS//eu-north-1/}
REGIONS=${REGIONS//ap-northeast-3/}

zip lambda-update-sqs-queue-depth.zip ./main

for REGION in $REGIONS; do
  aws s3 cp lambda-update-sqs-queue-depth.zip s3://$BUCKET_PREFIX-$REGION/lambda-update-sqs-queue-depth.zip --region $REGION
  aws s3api put-object-tagging --region $REGION --bucket $BUCKET_PREFIX-$REGION --key lambda-update-sqs-queue-depth.zip --tagging 'TagSet={Key=public,Value=yes}'
done





