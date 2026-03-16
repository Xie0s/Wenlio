/**
 * auth-captcha.ts - 认证安全检查业务逻辑层
 * 职责：封装认证页验证码 challenge 获取与校验请求，供 auth 相关组件复用。
 * 对外暴露：createCaptchaChallenge、verifyCaptcha。
 */
import { http } from '@/utils/http'
import type {
  CaptchaChallengeResp,
  CaptchaScene,
  CaptchaSignalSummary,
  VerifyCaptchaResp,
} from '@/types/captcha'

export async function createCaptchaChallenge(scene: CaptchaScene): Promise<CaptchaChallengeResp> {
  const res = await http.post<CaptchaChallengeResp>('/auth/captcha/challenge', { scene }, { throwHttpError: true, suppressErrorToast: true })
  return res.data
}

export async function verifyCaptcha(
  scene: CaptchaScene,
  challengeId: string,
  durationMs: number,
  signals: CaptchaSignalSummary,
): Promise<VerifyCaptchaResp> {
  const res = await http.post<VerifyCaptchaResp>('/auth/captcha/verify', {
    scene,
    challenge_id: challengeId,
    duration_ms: durationMs,
    signals,
  }, { throwHttpError: true, suppressErrorToast: true })
  return res.data
}
