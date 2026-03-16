export type CaptchaScene = 'login' | 'register' | 'reset_password'

export interface CaptchaChallengeResp {
  challenge_id: string
  scene: CaptchaScene
  mode: string
  prompt: string
  expires_at: number
  min_decision_ms: number
}

export interface CaptchaSignalSummary {
  dwell_ms: number
  visible_ms: number
  focused_ms: number
  visibility_changes: number
  focus_changes: number
  pointer_events: number
  key_events: number
  trusted_click: boolean
  language: string
  platform: string
  screen_width: number
  screen_height: number
  timezone_offset: number
  touch_points: number
  hardware_concurrency: number
  webdriver: boolean
}

export interface VerifyCaptchaResp {
  captcha_token: string
  expires_at: number
}
