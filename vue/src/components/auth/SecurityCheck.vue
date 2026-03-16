<!--
  SecurityCheck.vue - 自动安全检查组件

  职责：通过收集用户行为信号（停留时长、焦点变化、可见性等）自动完成人机验证，
       成功后返回 captcha_token 供后续接口使用。支持失败重试与 token 过期自动重置。
  对外暴露：
    - Props: modelValue (captcha_token), scene (验证场景), disabled (禁用状态)
    - Emits: update:modelValue (token 更新)
    - Methods: reset() (重置验证状态)
-->
<script setup lang="ts">
import { computed, onBeforeUnmount, onMounted, ref } from 'vue'
import { createCaptchaChallenge, verifyCaptcha } from '@/lib/auth-captcha'
import type { CaptchaScene, CaptchaSignalSummary } from '@/types/captcha'

const props = defineProps<{
  modelValue: string
  scene: CaptchaScene
  disabled?: boolean
}>()

const emit = defineEmits<{
  'update:modelValue': [value: string]
}>()

type Status = 'checking' | 'verified' | 'failed'
type RequestError = Error & { status?: number }

const status = ref<Status>('checking')
const tokenExpiresAt = ref(0)
const failureMessage = ref('')

const mountedAt = ref(0)
const focusStartedAt = ref(0)
const visibleStartedAt = ref(0)
const focusAccumulatedMs = ref(0)
const visibleAccumulatedMs = ref(0)
const visibilityChanges = ref(0)
const focusChanges = ref(0)

let tokenExpiryTimer: number | null = null
let retryTimer: number | null = null
let retryCount = 0
let disposed = false
let requestSerial = 0

const MAX_RETRY = 3

const statusText = computed(() => {
  if (status.value === 'verified') return '验证通过'
  if (status.value === 'failed') return failureMessage.value || (retryCount < MAX_RETRY ? '验证失败，准备重试...' : '验证失败，请刷新重试')
  return '正在验证当前环境...'
})

const hintText = computed(() => {
  if (status.value === 'verified') return '已完成自动安全检查，可继续提交'
  if (status.value === 'failed') return retryCount < MAX_RETRY ? '系统正在重新检测当前访问环境' : '请稍后刷新页面后重试'
  return '系统正在分析浏览器环境与访问行为'
})

const leftDisplayText = computed(() => {
  if (status.value === 'checking') return statusText.value
  if (status.value === 'verified') return '验证通过，可继续提交'
  return hintText.value
})

const rootClass = computed(() => {
  if (status.value === 'verified') {
    return 'border-[#d7d7d7] bg-transparent dark:border-slate-700 dark:bg-transparent'
  }
  if (status.value === 'failed') {
    return 'border-[#d7d7d7] bg-transparent dark:border-slate-700 dark:bg-transparent'
  }
  return 'border-[#d7d7d7] bg-transparent dark:border-slate-700 dark:bg-transparent'
})

const ringClass = computed(() => {
  if (status.value === 'verified') return 'text-[#16a34a]'
  if (status.value === 'failed') return 'text-[#dc2626]'
  return 'text-[#16a34a]'
})

function clearTokenTimer() {
  if (tokenExpiryTimer !== null) {
    window.clearTimeout(tokenExpiryTimer)
    tokenExpiryTimer = null
  }
}

function clearRetryTimer() {
  if (retryTimer !== null) {
    window.clearTimeout(retryTimer)
    retryTimer = null
  }
}

function normalizeError(error: unknown): RequestError {
  if (error instanceof Error) {
    return error as RequestError
  }
  return new Error('安全检查失败，请稍后再试') as RequestError
}

function shouldRetry(error: RequestError) {
  if (error.status === undefined) return true
  if (error.status === 0) return false
  if (error.status === 429) return false
  if (error.status >= 500) return false
  return error.status >= 400 && error.status < 500
}

function isStale(serial: number) {
  return disposed || serial !== requestSerial
}

function updateFocusClock() {
  if (document.hasFocus()) {
    if (focusStartedAt.value === 0) focusStartedAt.value = Date.now()
    return
  }
  if (focusStartedAt.value > 0) {
    focusAccumulatedMs.value += Date.now() - focusStartedAt.value
    focusStartedAt.value = 0
  }
}

function updateVisibleClock() {
  if (document.visibilityState === 'visible') {
    if (visibleStartedAt.value === 0) visibleStartedAt.value = Date.now()
    return
  }
  if (visibleStartedAt.value > 0) {
    visibleAccumulatedMs.value += Date.now() - visibleStartedAt.value
    visibleStartedAt.value = 0
  }
}

function currentFocusedMs() {
  return focusAccumulatedMs.value + (focusStartedAt.value > 0 ? Date.now() - focusStartedAt.value : 0)
}

function currentVisibleMs() {
  return visibleAccumulatedMs.value + (visibleStartedAt.value > 0 ? Date.now() - visibleStartedAt.value : 0)
}

function collectSignals(): CaptchaSignalSummary {
  return {
    dwell_ms: Date.now() - mountedAt.value,
    visible_ms: currentVisibleMs(),
    focused_ms: currentFocusedMs(),
    visibility_changes: visibilityChanges.value,
    focus_changes: focusChanges.value,
    pointer_events: 0,
    key_events: 0,
    trusted_click: false,
    language: navigator.language || '',
    platform: navigator.platform || '',
    screen_width: window.screen.width || 0,
    screen_height: window.screen.height || 0,
    timezone_offset: new Date().getTimezoneOffset(),
    touch_points: navigator.maxTouchPoints || 0,
    hardware_concurrency: navigator.hardwareConcurrency || 0,
    webdriver: Boolean(navigator.webdriver),
  }
}

function scheduleTokenExpiry(expiresAt: number) {
  clearTokenTimer()
  tokenExpiresAt.value = expiresAt
  const delay = expiresAt * 1000 - Date.now()
  if (delay <= 0) {
    reset()
    return
  }
  tokenExpiryTimer = window.setTimeout(() => {
    reset()
  }, delay)
}

async function runAutoVerify() {
  if (props.disabled || disposed) return

  const serial = ++requestSerial
  status.value = 'checking'
  failureMessage.value = ''
  emit('update:modelValue', '')

  try {
    const challengeCreatedAt = Date.now()
    const challenge = await createCaptchaChallenge(props.scene)
    if (isStale(serial)) return

    const remaining = challenge.min_decision_ms - (Date.now() - challengeCreatedAt)
    if (remaining > 0) {
      await new Promise<void>(resolve => setTimeout(resolve, remaining))
    }
    if (isStale(serial)) return

    const durationMs = Date.now() - challengeCreatedAt
    const result = await verifyCaptcha(
      props.scene,
      challenge.challenge_id,
      durationMs,
      collectSignals(),
    )
    if (isStale(serial)) return

    retryCount = 0
    status.value = 'verified'
    emit('update:modelValue', result.captcha_token)
    scheduleTokenExpiry(result.expires_at)
  } catch (error) {
    if (isStale(serial)) return
    const requestError = normalizeError(error)
    status.value = 'failed'
    failureMessage.value = requestError.message || '安全检查失败，请稍后再试'
    emit('update:modelValue', '')
    if (shouldRetry(requestError)) {
      scheduleRetry()
    }
  }
}

function scheduleRetry() {
  if (retryCount >= MAX_RETRY) return
  retryCount++
  const serial = requestSerial
  retryTimer = window.setTimeout(() => {
    if (isStale(serial)) return
    clearRetryTimer()
    tokenExpiresAt.value = 0
    failureMessage.value = ''
    emit('update:modelValue', '')
    void runAutoVerify()
  }, 3000)
}

function reset() {
  requestSerial++
  clearTokenTimer()
  clearRetryTimer()
  retryCount = 0
  tokenExpiresAt.value = 0
  status.value = 'checking'
  failureMessage.value = ''
  emit('update:modelValue', '')
  void runAutoVerify()
}

function handleVisibilityChange() {
  visibilityChanges.value += 1
  updateVisibleClock()
}

function handleWindowFocus() {
  focusChanges.value += 1
  updateFocusClock()
}

function handleWindowBlur() {
  focusChanges.value += 1
  updateFocusClock()
}

onMounted(() => {
  disposed = false
  const now = Date.now()
  mountedAt.value = now
  focusStartedAt.value = document.hasFocus() ? now : 0
  visibleStartedAt.value = document.visibilityState === 'visible' ? now : 0

  document.addEventListener('visibilitychange', handleVisibilityChange)
  window.addEventListener('focus', handleWindowFocus)
  window.addEventListener('blur', handleWindowBlur)

  void runAutoVerify()
})

onBeforeUnmount(() => {
  disposed = true
  requestSerial++
  updateFocusClock()
  updateVisibleClock()
  clearTokenTimer()
  clearRetryTimer()
  document.removeEventListener('visibilitychange', handleVisibilityChange)
  window.removeEventListener('focus', handleWindowFocus)
  window.removeEventListener('blur', handleWindowBlur)
})

defineExpose({ reset })
</script>

<template>
  <div
    class="relative flex min-h-[52px] w-full select-none items-center justify-between rounded-3xl border px-3 py-1.5 transition-colors duration-300"
    :class="rootClass"
  >
    <div class="flex min-w-0 flex-1 items-center gap-2">
      <div
        class="relative flex h-6 w-6 flex-shrink-0 items-center justify-center transition-colors duration-300"
        :class="ringClass"
      >
        <svg
          v-if="status === 'checking'"
          class="sc-spin h-6 w-6"
          viewBox="0 0 28 28"
          fill="none"
        >
          <circle
            cx="14"
            cy="14"
            r="10"
            stroke="currentColor"
            stroke-width="3"
            stroke-dasharray="1.4 5.2"
            stroke-linecap="round"
            opacity="0.95"
          />
        </svg>

        <svg
          v-else-if="status === 'verified'"
          class="sc-pop h-5 w-5"
          viewBox="0 0 20 20"
          fill="none"
        >
          <path
            d="M4.5 10.5l3.3 3.1L15.5 6"
            stroke="currentColor"
            stroke-width="2.4"
            stroke-linecap="round"
            stroke-linejoin="round"
          />
        </svg>

        <svg
          v-else
          class="h-[18px] w-[18px]"
          viewBox="0 0 20 20"
          fill="none"
        >
          <path
            d="M6 6l8 8M14 6l-8 8"
            stroke="currentColor"
            stroke-width="2.2"
            stroke-linecap="round"
          />
        </svg>
      </div>

      <div class="min-w-0 flex-1">
        <div
          class="truncate text-[14px] font-medium leading-[1.1] transition-colors duration-300 text-[#111111] dark:text-slate-100"
        >
          {{ leftDisplayText }}
        </div>
      </div>
    </div>

    <div class="ml-3 mr-1 flex w-[108px] flex-shrink-0 translate-y-[1px] flex-col items-end justify-center text-right leading-[1.02]">
      <span class="whitespace-nowrap text-[10.5px] font-semibold tracking-[0.03em] text-[#111111] dark:text-slate-100">Microswift Core</span>
      <span class="mt-0.5 whitespace-nowrap text-[9px] tracking-[0.01em] text-[#555555] dark:text-slate-400">微讯云信息安全验证</span>
    </div>
  </div>
</template>

<style scoped>
.sc-spin {
  animation: sc-spin 1s linear infinite;
  transform-origin: center;
}

.sc-pop {
  animation: sc-pop 0.4s cubic-bezier(0.34, 1.56, 0.64, 1) forwards;
}

@keyframes sc-spin {
  to {
    transform: rotate(360deg);
  }
}

@keyframes sc-pop {
  0% {
    transform: scale(0) rotate(-20deg);
    opacity: 0;
  }

  60% {
    transform: scale(1.2) rotate(5deg);
  }

  100% {
    transform: scale(1) rotate(0deg);
    opacity: 1;
  }
}
</style>
