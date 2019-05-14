#!/bin/bash

STACK_NAME=test-0

aws cloudformation create-stack \
  --stack-name $STACK_NAME \
  --template-body file://test.cfn.yml \
  --capabilities CAPABILITY_NAMED_IAM \

