package client

import (
	"context"
	"go-aws-sdk-retry/retryer"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/iam"
)

const SleepTimeSec = 5

type Iam struct {
	client *iam.Client
}

func NewIam(client *iam.Client) *Iam {
	return &Iam{
		client,
	}
}

func (i *Iam) RetryByOptionsSimpleParams(ctx context.Context, roleName *string) error {
	input := &iam.DeleteRoleInput{
		RoleName: roleName,
	}

	optFn := func(o *iam.Options) {
		o.RetryMaxAttempts = 3
		o.RetryMode = aws.RetryModeStandard
	}

	_, err := i.client.DeleteRole(ctx, input, optFn)

	return err
}

func (i *Iam) RetryByOptionsRetryer(ctx context.Context, roleName *string) error {
	input := &iam.DeleteRoleInput{
		RoleName: roleName,
	}

	retryable := func(err error) bool {
		return strings.Contains(err.Error(), "api error Throttling: Rate exceeded")
	}
	optFn := func(o *iam.Options) {
		o.Retryer = retryer.NewRetryer(retryable, SleepTimeSec)
	}

	_, err := i.client.DeleteRole(ctx, input, optFn)

	return err
}

func (i *Iam) RetryByGenerics(ctx context.Context, roleName *string) error {
	input := &iam.DeleteRoleInput{
		RoleName: roleName,
	}

	retryable := func(err error) bool {
		return strings.Contains(err.Error(), "api error Throttling: Rate exceeded")
	}

	_, err := retryer.Retry(
		&retryer.RetryInput[iam.DeleteRoleInput, iam.DeleteRoleOutput, iam.Options]{
			Ctx:              ctx,
			SleepTimeSec:     SleepTimeSec,
			TargetResource:   roleName,
			Input:            input,
			ApiCaller:        i.client.DeleteRole,
			RetryableChecker: retryable,
		},
	)

	return err
}
