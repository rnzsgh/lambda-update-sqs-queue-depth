package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"github.com/aws/aws-sdk-go/service/sqs"
	log "github.com/golang/glog"
)

func init() {
	flag.Parse()
	flag.Lookup("logtostderr").Value.Set("true")
}

func main() {
	lambda.Start(handler)
}

func handler(ctx context.Context, event events.CloudWatchEvent) error {
	if err := update(); err != nil {
		log.Errorf("Failed to update CloudWatch - reason: %v", err)
		return err
	}
	return nil
}

func update() error {
	queueUrl := aws.String(os.Getenv("QueueUrl"))
	metricName := aws.String(os.Getenv("CloudWatchMetricName"))
	namespace := aws.String(os.Getenv("CloudWatchMetricNamespace"))

	sqsSvc := sqs.New(session.New())

	depth, err := retry(func() (*float64, error) {
		out, err := sqsSvc.GetQueueAttributes(&sqs.GetQueueAttributesInput{
			QueueUrl:       queueUrl,
			AttributeNames: []*string{aws.String("ApproximateNumberOfMessages")},
		})

		if err != nil {
			return nil, err
		}

		v, err := strconv.ParseFloat(*(out.Attributes["ApproximateNumberOfMessages"]), 64)

		if err == nil {
			return aws.Float64(v), nil
		}

		return nil, err

	})

	if err != nil {
		return fmt.Errorf("Unable to get queue depth - reason: %v", err)
	}

	cloudWatch := cloudwatch.New(session.New())

	_, err = retry(func() (*float64, error) {
		_, err = cloudWatch.PutMetricData(&cloudwatch.PutMetricDataInput{
			Namespace: namespace,
			MetricData: []*cloudwatch.MetricDatum{
				&cloudwatch.MetricDatum{
					MetricName: metricName,
					Value:      depth,
				},
			},
		})

		return nil, err
	})

	if err != nil {
		return fmt.Errorf("Unable to put metric - reason: %v", err)
	}

	return nil
}

func retry(call func() (*float64, error)) (*float64, error) {
	var err error
	var v *float64
	for i := 0; i < 3; i++ {
		if v, err = call(); err != nil {
			time.Sleep(2 * time.Second)
		} else {
			return v, nil
		}
	}
	return nil, err
}
