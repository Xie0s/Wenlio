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

func TestCreateChallengeInvalidScene(t *testing.T) {
	manager := newTestManager()
	ctx := ClientContext{IP: "10.0.0.1", UserAgent: "test"}

	_, err := manager.CreateChallenge("invalid_scene", ctx)
	if err != ErrInvalidScene {
		t.Fatalf("CreateChallenge() error = %v, want ErrInvalidScene", err)
	}
}

func TestCreateChallengeAllScenes(t *testing.T) {
	scenes := []string{SceneLogin, SceneRegister, SceneResetPassword}
	for _, scene := range scenes {
		t.Run(scene, func(t *testing.T) {
			manager := newTestManager()
			ctx := ClientContext{IP: "10.0.0.1", UserAgent: "test"}

			challenge, err := manager.CreateChallenge(scene, ctx)
			if err != nil {
				t.Fatalf("CreateChallenge(%q) error = %v", scene, err)
			}
			if challenge.Scene != scene {
				t.Errorf("Challenge.Scene = %q, want %q", challenge.Scene, scene)
			}
			if challenge.ID == "" {
				t.Error("Challenge.ID is empty")
			}
			if challenge.MinDecisionMs < 500 || challenge.MinDecisionMs > 800 {
				t.Errorf("Challenge.MinDecisionMs = %d, want [500, 800]", challenge.MinDecisionMs)
			}
			if challenge.Mode != "managed_auto" {
				t.Errorf("Challenge.Mode = %q, want %q", challenge.Mode, "managed_auto")
			}
		})
	}
}

func TestVerifyRejectsNilInput(t *testing.T) {
	manager := newTestManager()
	ctx := ClientContext{IP: "10.0.0.1", UserAgent: "test"}

	_, err := manager.Verify(nil, ctx)
	if err != ErrInvalidScene {
		t.Fatalf("Verify(nil) error = %v, want ErrInvalidScene", err)
	}
}

func TestVerifyRejectsInvalidChallengeID(t *testing.T) {
	manager := newTestManager()
	ctx := ClientContext{IP: "10.0.0.1", UserAgent: "test"}

	_, err := manager.Verify(&VerifyInput{
		ChallengeID: "nonexistent",
		Scene:       SceneLogin,
		DurationMs:  1000,
		Signals:     validSignals(1000),
	}, ctx)
	if err != ErrChallengeNotFound {
		t.Fatalf("Verify() error = %v, want ErrChallengeNotFound", err)
	}
}

func TestVerifyRejectsSceneMismatch(t *testing.T) {
	manager := newTestManager()
	ctx := ClientContext{IP: "10.0.0.1", UserAgent: "test"}

	challenge, err := manager.CreateChallenge(SceneLogin, ctx)
	if err != nil {
		t.Fatalf("CreateChallenge() error = %v", err)
	}

	duration := challenge.MinDecisionMs + 200
	manager.challenges[challenge.ID].IssuedAt = time.Now().Add(-time.Duration(duration) * time.Millisecond)

	_, err = manager.Verify(&VerifyInput{
		ChallengeID: challenge.ID,
		Scene:       SceneRegister, // wrong scene
		DurationMs:  duration,
		Signals:     validSignals(duration),
	}, ctx)
	if err != ErrChallengeFailed {
		t.Fatalf("Verify() error = %v, want ErrChallengeFailed", err)
	}
}

func TestConsumeTokenRejectsInvalidScene(t *testing.T) {
	manager := newTestManager()
	ctx := ClientContext{IP: "10.0.0.1", UserAgent: "test"}

	err := manager.ConsumeToken("invalid_scene", "some-token", ctx)
	if err != ErrTokenInvalid {
		t.Fatalf("ConsumeToken() error = %v, want ErrTokenInvalid", err)
	}
}

func TestConsumeTokenRejectsEmptyToken(t *testing.T) {
	manager := newTestManager()
	ctx := ClientContext{IP: "10.0.0.1", UserAgent: "test"}

	err := manager.ConsumeToken(SceneLogin, "", ctx)
	if err != ErrTokenInvalid {
		t.Fatalf("ConsumeToken() error = %v, want ErrTokenInvalid", err)
	}
}

func TestConsumeTokenRejectsWhitespaceToken(t *testing.T) {
	manager := newTestManager()
	ctx := ClientContext{IP: "10.0.0.1", UserAgent: "test"}

	err := manager.ConsumeToken(SceneLogin, "   ", ctx)
	if err != ErrTokenInvalid {
		t.Fatalf("ConsumeToken() error = %v, want ErrTokenInvalid", err)
	}
}

func TestTokenCannotBeConsumedTwice(t *testing.T) {
	manager := newTestManager()
	ctx := ClientContext{IP: "192.168.1.50", UserAgent: "test-agent"}

	challenge, err := manager.CreateChallenge(SceneLogin, ctx)
	if err != nil {
		t.Fatalf("CreateChallenge() error = %v", err)
	}

	duration := challenge.MinDecisionMs + 300
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

	// First consume should succeed
	if err := manager.ConsumeToken(SceneLogin, result.Token, ctx); err != nil {
		t.Fatalf("First ConsumeToken() error = %v", err)
	}

	// Second consume should fail
	if err := manager.ConsumeToken(SceneLogin, result.Token, ctx); err == nil {
		t.Fatal("Second ConsumeToken() should fail (token already consumed)")
	}
}

func TestVerifyRejectsZeroScreenDimensions(t *testing.T) {
	manager := newTestManager()
	ctx := ClientContext{IP: "10.0.0.1", UserAgent: "test"}

	challenge, err := manager.CreateChallenge(SceneLogin, ctx)
	if err != nil {
		t.Fatalf("CreateChallenge() error = %v", err)
	}

	duration := challenge.MinDecisionMs + 200
	manager.challenges[challenge.ID].IssuedAt = time.Now().Add(-time.Duration(duration) * time.Millisecond)

	signals := validSignals(duration)
	signals.ScreenWidth = 0
	signals.ScreenHeight = 0

	_, err = manager.Verify(&VerifyInput{
		ChallengeID: challenge.ID,
		Scene:       SceneLogin,
		DurationMs:  duration,
		Signals:     signals,
	}, ctx)
	if err != ErrChallengeFailed {
		t.Fatalf("Verify() error = %v, want ErrChallengeFailed", err)
	}
}

func TestVerifyRejectsTooShortDwell(t *testing.T) {
	manager := newTestManager()
	ctx := ClientContext{IP: "10.0.0.1", UserAgent: "test"}

	challenge, err := manager.CreateChallenge(SceneLogin, ctx)
	if err != nil {
		t.Fatalf("CreateChallenge() error = %v", err)
	}

	// Set proper IssuedAt but use too-short dwell time in signals
	duration := challenge.MinDecisionMs + 200
	manager.challenges[challenge.ID].IssuedAt = time.Now().Add(-time.Duration(duration) * time.Millisecond)

	signals := validSignals(duration)
	signals.DwellMs = 100 // too short

	_, err = manager.Verify(&VerifyInput{
		ChallengeID: challenge.ID,
		Scene:       SceneLogin,
		DurationMs:  duration,
		Signals:     signals,
	}, ctx)
	if err != ErrChallengeFailed {
		t.Fatalf("Verify() error = %v, want ErrChallengeFailed", err)
	}
}

func TestNewManagerPanicsOnEmptySecret(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("NewManager(nil) should panic")
		}
	}()
	NewManager(nil)
}

func TestNewManagerPanicsOnZeroLenSecret(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("NewManager([]byte{}) should panic")
		}
	}()
	NewManager([]byte{})
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
