package captcha

import (
	"testing"
	"time"
)

func TestManagerVerifyAndConsume(t *testing.T) {
	manager := newTestManager()
	ctx := ClientContext{IP: "192.168.1.23", UserAgent: "test-agent"}

	challenge, err := manager.CreateChallenge(SceneLogin, ctx)
	if err != nil {
		t.Fatalf("CreateChallenge() error = %v", err)
	}

	duration := challenge.MinDecisionMs + 220
	manager.challenges[challenge.ID].IssuedAt = time.Now().Add(-time.Duration(duration) * time.Millisecond)

	result, err := manager.Verify(&VerifyInput{
		ChallengeID: challenge.ID,
		Scene:       SceneLogin,
		DurationMs:  duration,
		Signals:     validSignals(duration),
	}, ctx)
	if err != nil {
		t.Fatalf("Verify() error = %v", err)
	}

	if err := manager.ConsumeToken(SceneLogin, result.Token, ctx); err != nil {
		t.Fatalf("ConsumeToken() error = %v", err)
	}
}

func TestManagerRejectsContextMismatch(t *testing.T) {
	manager := newTestManager()
	createCtx := ClientContext{IP: "10.1.2.88", UserAgent: "agent-a"}
	verifyCtx := ClientContext{IP: "10.1.3.99", UserAgent: "agent-b"}

	challenge, err := manager.CreateChallenge(SceneRegister, createCtx)
	if err != nil {
		t.Fatalf("CreateChallenge() error = %v", err)
	}

	duration := challenge.MinDecisionMs + 200
	manager.challenges[challenge.ID].IssuedAt = time.Now().Add(-time.Duration(duration) * time.Millisecond)

	_, err = manager.Verify(&VerifyInput{
		ChallengeID: challenge.ID,
		Scene:       SceneRegister,
		DurationMs:  duration,
		Signals:     validSignals(duration),
	}, verifyCtx)
	if err == nil {
		t.Fatalf("Verify() should fail on context mismatch")
	}
}

func TestManagerRejectsWebdriverSignals(t *testing.T) {
	manager := newTestManager()
	ctx := ClientContext{IP: "172.16.10.6", UserAgent: "agent-c"}

	challenge, err := manager.CreateChallenge(SceneResetPassword, ctx)
	if err != nil {
		t.Fatalf("CreateChallenge() error = %v", err)
	}

	duration := challenge.MinDecisionMs + 160
	manager.challenges[challenge.ID].IssuedAt = time.Now().Add(-time.Duration(duration) * time.Millisecond)
	signals := validSignals(duration)
	signals.Webdriver = true

	_, err = manager.Verify(&VerifyInput{
		ChallengeID: challenge.ID,
		Scene:       SceneResetPassword,
		DurationMs:  duration,
		Signals:     signals,
	}, ctx)
	if err == nil {
		t.Fatalf("Verify() should reject webdriver signals")
	}
}

func newTestManager() *Manager {
	return &Manager{
		secret:         []byte("0123456789abcdef0123456789abcdef"),
		challenges:     make(map[string]*challengeState),
		createAttempts: make(map[string][]time.Time),
		verifyAttempts: make(map[string][]time.Time),
		verifyFailures: make(map[string]*failureState),
	}
}

func validSignals(duration int64) SignalSummary {
	return SignalSummary{
		DwellMs:             duration,
		VisibleMs:           duration,
		FocusedMs:           duration,
		VisibilityChanges:   0,
		FocusChanges:        1,
		PointerEvents:       2,
		KeyEvents:           0,
		TrustedClick:        true,
		Language:            "zh-CN",
		Platform:            "Win32",
		ScreenWidth:         1920,
		ScreenHeight:        1080,
		TimezoneOffset:      -480,
		TouchPoints:         0,
		HardwareConcurrency: 8,
		Webdriver:           false,
	}
}
