/**
 * 文件：src/lib/version.ts
 * 作用：前端版本探针，周期性读取 /version.json 用于检测新版本发布。
 * 职责边界：
 * 1) 仅在生产环境启用，不干扰本地开发体验。
 * 2) 在页面可见或定时轮询时检查版本，发现变更后分发事件。
 * 3) 不直接强制刷新页面，交由上层 UI 决定提示与刷新策略。
 * 对外接口：
 * - VERSION_UPDATE_EVENT：版本更新事件名
 * - setupBuildVersionProbe(options?)：初始化版本探针
 */

export const VERSION_UPDATE_EVENT = 'app:version-update'

export type BuildVersionManifest = {
  buildVersion: string
  packageVersion: string
  buildSerial: string
  commitSha: string
  commitShortSha: string
  source: string
  builtAt: string
}

export type BuildVersionUpdateDetail = {
  previousVersion: string
  latest: BuildVersionManifest
}

export type BuildVersionProbeOptions = {
  pollIntervalMs?: number
  onVersionChange?: (detail: BuildVersionUpdateDetail) => void
}

const VERSION_MANIFEST_URL = '/version.json'
const DEFAULT_POLL_INTERVAL_MS = 5 * 60 * 1000
const MIN_POLL_INTERVAL_MS = 60 * 1000

let hasInitialized = false
let currentBuildVersion: string | null = null

function isBuildVersionManifest(value: unknown): value is BuildVersionManifest {
  if (!value || typeof value !== 'object') {
    return false
  }

  const data = value as Partial<BuildVersionManifest>
  return typeof data.buildVersion === 'string' && data.buildVersion.length > 0
}

async function fetchBuildVersionManifest(): Promise<BuildVersionManifest | null> {
  try {
    const response = await fetch(VERSION_MANIFEST_URL, {
      method: 'GET',
      cache: 'no-store',
      headers: {
        'cache-control': 'no-cache',
      },
    })

    if (!response.ok) {
      return null
    }

    const payload = (await response.json()) as unknown
    if (!isBuildVersionManifest(payload)) {
      return null
    }

    return payload
  } catch {
    return null
  }
}

export function setupBuildVersionProbe(options: BuildVersionProbeOptions = {}): void {
  if (hasInitialized || !import.meta.env.PROD) {
    return
  }

  hasInitialized = true

  const normalizedInterval = Math.max(
    options.pollIntervalMs ?? DEFAULT_POLL_INTERVAL_MS,
    MIN_POLL_INTERVAL_MS,
  )

  const checkForNewVersion = async () => {
    const manifest = await fetchBuildVersionManifest()
    if (!manifest) {
      return
    }

    if (!currentBuildVersion) {
      currentBuildVersion = manifest.buildVersion
      return
    }

    if (manifest.buildVersion === currentBuildVersion) {
      return
    }

    const previousVersion = currentBuildVersion
    currentBuildVersion = manifest.buildVersion

    const detail: BuildVersionUpdateDetail = {
      previousVersion,
      latest: manifest,
    }

    window.dispatchEvent(new CustomEvent<BuildVersionUpdateDetail>(VERSION_UPDATE_EVENT, { detail }))

    if (options.onVersionChange) {
      options.onVersionChange(detail)
    }

    console.info(
      `[version] detected new build version: ${previousVersion} -> ${manifest.buildVersion}`,
    )
  }

  void checkForNewVersion()

  window.setInterval(() => {
    void checkForNewVersion()
  }, normalizedInterval)

  document.addEventListener('visibilitychange', () => {
    if (document.visibilityState === 'visible') {
      void checkForNewVersion()
    }
  })
}
