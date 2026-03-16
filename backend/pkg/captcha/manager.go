// Package captcha 提供可直接迁入业务项目的验证码核心能力。
//
// 职责：管理 challenge 创建、环境信号校验、一次性 token 签发与消费、频控与失败惩罚。
// 对外暴露：Manager、ClientContext、SignalSummary、Challenge、VerifyInput、VerifyResult 以及相关错误定义。
package captcha

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"math/big"
	"net"
	"strings"
	"sync"
	"time"
)

const (
	SceneLogin         = "login"
	SceneRegister      = "register"
	SceneResetPassword = "reset_password"
)

var (
	ErrInvalidScene      = errors.New("invalid captcha scene")
	ErrChallengeNotFound = errors.New("captcha challenge not found")
	ErrChallengeExpired  = errors.New("captcha challenge expired")
	ErrChallengeFailed   = errors.New("captcha challenge failed")
	ErrTokenInvalid      = errors.New("captcha token invalid")
	ErrTokenExpired      = errors.New("captcha token expired")
	ErrTooManyRequests   = errors.New("captcha too many requests")
	ErrClientBlocked     = errors.New("captcha client temporarily blocked")

	defaultManager *Manager
	once           sync.Once
)

type ClientContext struct {
	IP        string
	UserAgent string
}

type SignalSummary struct {
	DwellMs             int64
	VisibleMs           int64
	FocusedMs           int64
	VisibilityChanges   int
	FocusChanges        int
	PointerEvents       int
	KeyEvents           int
	TrustedClick        bool
	Language            string
	Platform            string
	ScreenWidth         int
	ScreenHeight        int
	TimezoneOffset      int
	TouchPoints         int
	HardwareConcurrency int
	Webdriver           bool
}

type Challenge struct {
	ID            string
	Scene         string
	Mode          string
	Prompt        string
	ExpiresAt     int64
	MinDecisionMs int64
}

type VerifyInput struct {
	ChallengeID string
	Scene       string
	DurationMs  int64
	Signals     SignalSummary
}

type VerifyResult struct {
	Token     string
	ExpiresAt int64
}

type Manager struct {
	mu             sync.Mutex
	secret         []byte
	challenges     map[string]*challengeState
	createAttempts map[string][]time.Time
	verifyAttempts map[string][]time.Time
	verifyFailures map[string]*failureState
}

type challengeState struct {
	ID            string
	Scene         string
	IssuedAt      time.Time
	ExpiresAt     time.Time
	MinDecisionMs int64
	ContextHash   string
	TokenNonce    string
	VerifiedAt    time.Time
	Consumed      bool
}

type failureState struct {
	Count     int
	LastAt    time.Time
	BlockedAt time.Time
}

type tokenPayload struct {
	ChallengeID string `json:"challenge_id"`
	Scene       string `json:"scene"`
	Nonce       string `json:"nonce"`
	ContextHash string `json:"context_hash"`
	ExpiresAt   int64  `json:"expires_at"`
}

func DefaultManager() *Manager {
	once.Do(func() {
		secret, err := randomBytes(32)
		if err != nil {
			panic(err)
		}
		defaultManager = NewManager(secret)
	})
	return defaultManager
}

func NewManager(secret []byte) *Manager {
	if len(secret) == 0 {
		panic("captcha secret is required")
	}
	copiedSecret := make([]byte, len(secret))
	copy(copiedSecret, secret)
	return &Manager{
		secret:         copiedSecret,
		challenges:     make(map[string]*challengeState),
		createAttempts: make(map[string][]time.Time),
		verifyAttempts: make(map[string][]time.Time),
		verifyFailures: make(map[string]*failureState),
	}
}

func (m *Manager) CreateChallenge(scene string, ctx ClientContext) (*Challenge, error) {
	if !isValidScene(scene) {
		return nil, ErrInvalidScene
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	now := time.Now()
	m.pruneLocked(now)
	if !m.allowAttemptLocked(m.createAttempts, m.clientKey(ctx), now, 6, 20*time.Second) {
		return nil, ErrTooManyRequests
	}

	id, err := randomString(24)
	if err != nil {
		return nil, err
	}
	minDecisionMs, err := randomRange(500, 800)
	if err != nil {
		return nil, err
	}

	expiresAt := now.Add(2 * time.Minute)
	m.challenges[id] = &challengeState{
		ID:            id,
		Scene:         scene,
		IssuedAt:      now,
		ExpiresAt:     expiresAt,
		MinDecisionMs: minDecisionMs,
		ContextHash:   m.contextHash(ctx),
	}

	return &Challenge{
		ID:            id,
		Scene:         scene,
		Mode:          "managed_auto",
		Prompt:        "正在验证...",
		ExpiresAt:     expiresAt.Unix(),
		MinDecisionMs: minDecisionMs,
	}, nil
}

func (m *Manager) Verify(input *VerifyInput, ctx ClientContext) (*VerifyResult, error) {
	if input == nil || !isValidScene(input.Scene) {
		return nil, ErrInvalidScene
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	now := time.Now()
	m.pruneLocked(now)
	clientKey := m.clientKey(ctx)
	if m.isBlockedLocked(clientKey, now) {
		return nil, ErrClientBlocked
	}
	if !m.allowAttemptLocked(m.verifyAttempts, clientKey, now, 8, 20*time.Second) {
		return nil, ErrTooManyRequests
	}

	state, ok := m.challenges[input.ChallengeID]
	if !ok {
		m.recordFailureLocked(clientKey, now)
		return nil, ErrChallengeNotFound
	}
	if state.Scene != input.Scene {
		m.recordFailureLocked(clientKey, now)
		return nil, ErrChallengeFailed
	}
	if state.ContextHash != m.contextHash(ctx) {
		m.recordFailureLocked(clientKey, now)
		return nil, ErrChallengeFailed
	}
	if now.After(state.ExpiresAt) {
		delete(m.challenges, input.ChallengeID)
		return nil, ErrChallengeExpired
	}
	if state.Consumed {
		delete(m.challenges, input.ChallengeID)
		m.recordFailureLocked(clientKey, now)
		return nil, ErrTokenInvalid
	}
	if !validateSignals(input, state, now) {
		m.recordFailureLocked(clientKey, now)
		return nil, ErrChallengeFailed
	}

	nonce, err := randomString(18)
	if err != nil {
		return nil, err
	}
	tokenExpiresAt := now.Add(3 * time.Minute)
	state.TokenNonce = nonce
	state.VerifiedAt = now

	token, err := m.signToken(tokenPayload{
		ChallengeID: state.ID,
		Scene:       state.Scene,
		Nonce:       nonce,
		ContextHash: state.ContextHash,
		ExpiresAt:   tokenExpiresAt.Unix(),
	})
	if err != nil {
		return nil, err
	}

	delete(m.verifyFailures, clientKey)
	return &VerifyResult{Token: token, ExpiresAt: tokenExpiresAt.Unix()}, nil
}

func validateSignals(input *VerifyInput, state *challengeState, now time.Time) bool {
	serverElapsed := now.Sub(state.IssuedAt).Milliseconds()
	if input.DurationMs < state.MinDecisionMs || serverElapsed+150 < state.MinDecisionMs {
		return false
	}

	signals := input.Signals
	if signals.Webdriver {
		return false
	}
	if signals.DwellMs < state.MinDecisionMs {
		return false
	}
	if signals.VisibleMs < minInt64(state.MinDecisionMs, 400) {
		return false
	}
	if signals.ScreenWidth <= 0 || signals.ScreenHeight <= 0 {
		return false
	}
	if signals.HardwareConcurrency == 0 && signals.TouchPoints == 0 {
		return false
	}
	if signals.VisibilityChanges > 6 || signals.FocusChanges > 8 {
		return false
	}

	return true
}

func (m *Manager) ConsumeToken(scene, token string, ctx ClientContext) error {
	if !isValidScene(scene) || strings.TrimSpace(token) == "" {
		return ErrTokenInvalid
	}

	payload, err := m.parseToken(token)
	if err != nil {
		if errors.Is(err, ErrTokenExpired) {
			return ErrTokenExpired
		}
		return ErrTokenInvalid
	}
	if payload.Scene != scene || payload.ContextHash != m.contextHash(ctx) {
		return ErrTokenInvalid
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	now := time.Now()
	m.pruneLocked(now)

	state, ok := m.challenges[payload.ChallengeID]
	if !ok {
		return ErrTokenInvalid
	}
	if now.After(state.ExpiresAt) || now.Unix() > payload.ExpiresAt {
		delete(m.challenges, payload.ChallengeID)
		return ErrTokenExpired
	}
	if state.Scene != scene || state.ContextHash != payload.ContextHash || state.TokenNonce == "" || state.TokenNonce != payload.Nonce {
		return ErrTokenInvalid
	}
	if state.Consumed {
		delete(m.challenges, payload.ChallengeID)
		return ErrTokenInvalid
	}

	state.Consumed = true
	delete(m.challenges, payload.ChallengeID)
	return nil
}

func (m *Manager) signToken(payload tokenPayload) (string, error) {
	body, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}
	encoded := base64.RawURLEncoding.EncodeToString(body)
	sig := m.sign([]byte(encoded))
	return encoded + "." + base64.RawURLEncoding.EncodeToString(sig), nil
}

func (m *Manager) parseToken(token string) (*tokenPayload, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 2 {
		return nil, ErrTokenInvalid
	}

	body := parts[0]
	sig, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return nil, ErrTokenInvalid
	}
	if !hmac.Equal(sig, m.sign([]byte(body))) {
		return nil, ErrTokenInvalid
	}

	bodyBytes, err := base64.RawURLEncoding.DecodeString(body)
	if err != nil {
		return nil, ErrTokenInvalid
	}

	var payload tokenPayload
	if err := json.Unmarshal(bodyBytes, &payload); err != nil {
		return nil, ErrTokenInvalid
	}
	if time.Now().Unix() > payload.ExpiresAt {
		return nil, ErrTokenExpired
	}
	return &payload, nil
}

func (m *Manager) sign(data []byte) []byte {
	h := hmac.New(sha256.New, m.secret)
	_, _ = h.Write(data)
	return h.Sum(nil)
}

func (m *Manager) contextHash(ctx ClientContext) string {
	normalized := normalizeIP(ctx.IP) + "|" + normalizeUA(ctx.UserAgent)
	return base64.RawURLEncoding.EncodeToString(m.sign([]byte(normalized)))
}

func normalizeIP(raw string) string {
	ip := net.ParseIP(strings.TrimSpace(raw))
	if ip == nil {
		return "unknown"
	}
	if v4 := ip.To4(); v4 != nil {
		return net.IPv4(v4[0], v4[1], v4[2], 0).String()
	}
	v6 := ip.To16()
	if v6 == nil {
		return "unknown"
	}
	masked := make(net.IP, len(v6))
	copy(masked, v6)
	for i := 8; i < len(masked); i++ {
		masked[i] = 0
	}
	return masked.String()
}

func normalizeUA(raw string) string {
	ua := strings.ToLower(strings.TrimSpace(raw))
	if ua == "" {
		return "unknown"
	}
	if len(ua) > 160 {
		return ua[:160]
	}
	return ua
}

func (m *Manager) pruneLocked(now time.Time) {
	for id, state := range m.challenges {
		if state.Consumed || now.After(state.ExpiresAt.Add(3*time.Minute)) {
			delete(m.challenges, id)
		}
	}
	m.pruneAttemptMapLocked(m.createAttempts, now, 20*time.Second)
	m.pruneAttemptMapLocked(m.verifyAttempts, now, 20*time.Second)
	for key, state := range m.verifyFailures {
		if state.Count == 0 || now.Sub(state.LastAt) > 10*time.Minute {
			delete(m.verifyFailures, key)
		}
	}
}

func (m *Manager) allowAttemptLocked(bucket map[string][]time.Time, key string, now time.Time, limit int, window time.Duration) bool {
	windowStart := now.Add(-window)
	requests := bucket[key]
	valid := requests[:0]
	for _, ts := range requests {
		if ts.After(windowStart) {
			valid = append(valid, ts)
		}
	}
	if len(valid) >= limit {
		if len(valid) == 0 {
			delete(bucket, key)
		} else {
			bucket[key] = valid
		}
		return false
	}
	valid = append(valid, now)
	bucket[key] = valid
	return true
}

func (m *Manager) pruneAttemptMapLocked(bucket map[string][]time.Time, now time.Time, window time.Duration) {
	windowStart := now.Add(-window)
	for key, requests := range bucket {
		valid := requests[:0]
		for _, ts := range requests {
			if ts.After(windowStart) {
				valid = append(valid, ts)
			}
		}
		if len(valid) == 0 {
			delete(bucket, key)
			continue
		}
		bucket[key] = valid
	}
}

func (m *Manager) recordFailureLocked(key string, now time.Time) {
	state, ok := m.verifyFailures[key]
	if !ok || now.Sub(state.LastAt) > 5*time.Minute {
		m.verifyFailures[key] = &failureState{
			Count:  1,
			LastAt: now,
		}
		return
	}
	state.Count++
	state.LastAt = now
	if state.Count >= 5 {
		state.BlockedAt = now.Add(2 * time.Minute)
	}
}

func (m *Manager) isBlockedLocked(key string, now time.Time) bool {
	state, ok := m.verifyFailures[key]
	if !ok {
		return false
	}
	if !state.BlockedAt.IsZero() && now.Before(state.BlockedAt) {
		return true
	}
	if !state.BlockedAt.IsZero() && !now.Before(state.BlockedAt) {
		delete(m.verifyFailures, key)
	}
	return false
}

func (m *Manager) clientKey(ctx ClientContext) string {
	return normalizeIP(ctx.IP) + "|" + normalizeUA(ctx.UserAgent)
}

func isValidScene(scene string) bool {
	return scene == SceneLogin || scene == SceneRegister || scene == SceneResetPassword
}

func minInt64(a, b int64) int64 {
	if a < b {
		return a
	}
	return b
}

func randomString(length int) (string, error) {
	buf, err := randomBytes(length)
	if err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(buf)[:length], nil
}

func randomBytes(length int) ([]byte, error) {
	buf := make([]byte, length)
	_, err := rand.Read(buf)
	return buf, err
}

func randomRange(min, max int64) (int64, error) {
	if max <= min {
		return min, nil
	}
	n, err := rand.Int(rand.Reader, big.NewInt(max-min+1))
	if err != nil {
		return 0, err
	}
	return min + n.Int64(), nil
}
