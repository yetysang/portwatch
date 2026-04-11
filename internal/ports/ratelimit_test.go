package ports

import (
	"testing"
	"time"
)

func fixedClock(t time.Time) func() time.Time {
	return func() time.Time { return t }
}

func TestRateLimiter_FirstCallAlwaysAllowed(t *testing.T) {
	rl := NewRateLimiter(10 * time.Second)
	if !rl.Allow("tcp:8080") {
		t.Fatal("expected first call to be allowed")
	}
}

func TestRateLimiter_SecondCallWithinCooldownBlocked(t *testing.T) {
	base := time.Now()
	rl := NewRateLimiter(10 * time.Second)
	rl.now = fixedClock(base)

	rl.Allow("tcp:8080")
	rl.now = fixedClock(base.Add(5 * time.Second))

	if rl.Allow("tcp:8080") {
		t.Fatal("expected second call within cooldown to be blocked")
	}
}

func TestRateLimiter_AllowedAfterCooldown(t *testing.T) {
	base := time.Now()
	rl := NewRateLimiter(10 * time.Second)
	rl.now = fixedClock(base)

	rl.Allow("tcp:8080")
	rl.now = fixedClock(base.Add(11 * time.Second))

	if !rl.Allow("tcp:8080") {
		t.Fatal("expected call after cooldown to be allowed")
	}
}

func TestRateLimiter_DifferentKeysIndependent(t *testing.T) {
	rl := NewRateLimiter(10 * time.Second)
	rl.Allow("tcp:8080")

	if !rl.Allow("tcp:9090") {
		t.Fatal("expected different key to be allowed")
	}
}

func TestRateLimiter_Reset(t *testing.T) {
	base := time.Now()
	rl := NewRateLimiter(10 * time.Second)
	rl.now = fixedClock(base)

	rl.Allow("tcp:8080")
	rl.Reset()

	if !rl.Allow("tcp:8080") {
		t.Fatal("expected allow after reset")
	}
}

func TestRateLimiter_Len(t *testing.T) {
	rl := NewRateLimiter(10 * time.Second)
	if rl.Len() != 0 {
		t.Fatalf("expected 0, got %d", rl.Len())
	}
	rl.Allow("a")
	rl.Allow("b")
	rl.Allow("a") // duplicate — should not increase len
	if rl.Len() != 2 {
		t.Fatalf("expected 2, got %d", rl.Len())
	}
}
