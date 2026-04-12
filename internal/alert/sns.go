package alert

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sns"
	"github.com/user/portwatch/internal/monitor"
)

// snsPublisher abstracts the SNS Publish call for testing.
type snsPublisher interface {
	Publish(ctx context.Context, params *sns.PublishInput, optFns ...func(*sns.Options)) (*sns.PublishOutput, error)
}

// SNSHandler sends port change alerts to an AWS SNS topic.
type SNSHandler struct {
	client   snsPublisher
	topicARN string
}

// NewSNSHandler creates an SNSHandler that publishes to the given SNS topic ARN.
// It loads AWS credentials from the default credential chain (env, ~/.aws, IAM role).
func NewSNSHandler(topicARN, region string) (*SNSHandler, error) {
	if topicARN == "" {
		return nil, fmt.Errorf("sns: topic ARN must not be empty")
	}
	cfg, err := config.LoadDefaultConfig(context.Background(),
		config.WithRegion(region),
	)
	if err != nil {
		return nil, fmt.Errorf("sns: load aws config: %w", err)
	}
	return &SNSHandler{
		client:   sns.NewFromConfig(cfg),
		topicARN: topicARN,
	}, nil
}

// Handle publishes each change as a separate SNS message.
func (h *SNSHandler) Handle(changes []monitor.Change) error {
	if len(changes) == 0 {
		return nil
	}
	for _, c := range changes {
		msg := formatSNSMsg(c)
		_, err := h.client.Publish(context.Background(), &sns.PublishInput{
			TopicArn: aws.String(h.topicARN),
			Message:  aws.String(msg),
			Subject:  aws.String(fmt.Sprintf("portwatch: port %s %s", c.Binding.Port, c.Kind)),
		})
		if err != nil {
			return fmt.Errorf("sns: publish: %w", err)
		}
	}
	return nil
}

// Drain is a no-op for SNS; messages are sent immediately.
func (h *SNSHandler) Drain() error { return nil }

func formatSNSMsg(c monitor.Change) string {
	host := c.Binding.Hostname
	if host == "" {
		host = c.Binding.Addr
	}
	proc := c.Binding.Process
	if proc == "" {
		proc = "unknown"
	}
	return fmt.Sprintf("[%s] %s:%s proto=%s pid=%d process=%s",
		c.Kind, host, c.Binding.Port, c.Binding.Proto, c.Binding.PID, proc)
}
