package config

import "testing"

func TestDefaultSNSConfig_Values(t *testing.T) {
	cfg := DefaultSNSConfig()
	if cfg.Enabled {
		t.Error("expected Enabled=false by default")
	}
	if cfg.TopicARN != "" {
		t.Errorf("expected empty TopicARN, got %q", cfg.TopicARN)
	}
	if cfg.Region != "us-east-1" {
		t.Errorf("expected Region=us-east-1, got %q", cfg.Region)
	}
}

func TestSNSConfig_ValidateDisabledSkipsChecks(t *testing.T) {
	cfg := SNSConfig{Enabled: false, TopicARN: "", Region: ""}
	if err := cfg.Validate(); err != nil {
		t.Errorf("expected no error for disabled config, got %v", err)
	}
}

func TestSNSConfig_ValidateEnabledRequiresTopicARN(t *testing.T) {
	cfg := SNSConfig{Enabled: true, TopicARN: "", Region: "us-east-1"}
	if err := cfg.Validate(); err == nil {
		t.Error("expected error when topic_arn is empty")
	}
}

func TestSNSConfig_ValidateEnabledRequiresRegion(t *testing.T) {
	cfg := SNSConfig{Enabled: true, TopicARN: "arn:aws:sns:us-east-1:123456789012:portwatch", Region: ""}
	if err := cfg.Validate(); err == nil {
		t.Error("expected error when region is empty")
	}
}

func TestSNSConfig_ValidateEnabledWithBothFields(t *testing.T) {
	cfg := SNSConfig{
		Enabled:  true,
		TopicARN: "arn:aws:sns:us-east-1:123456789012:portwatch",
		Region:   "us-east-1",
	}
	if err := cfg.Validate(); err != nil {
		t.Errorf("expected no error for valid config, got %v", err)
	}
}

func TestSNSConfig_ValidateErrorMessage(t *testing.T) {
	cfg := SNSConfig{Enabled: true, TopicARN: "", Region: "us-east-1"}
	err := cfg.Validate()
	if err == nil {
		t.Fatal("expected an error")
	}
	if err.Error() != "sns: topic_arn is required when enabled" {
		t.Errorf("unexpected error message: %v", err)
	}
}
