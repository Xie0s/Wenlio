/**
 * 文件：scripts/postbuild.js
 * 作用：构建完成后处理脚本，用于生成可追溯版本号并注入产物，配合 CDN 缓存策略。
 * 职责边界：
 * 1) 优先使用 CI 构建号 + Commit SHA 生成版本号。
 * 2) 在无 CI 环境时，使用 scripts/.build-counter 作为本地兜底构建号。
 * 3) 将 dist/sw.js 中的 __SW_BUILD_VERSION__ 占位符替换为实际版本号。
 * 4) 生成 dist/version.json 版本清单，便于排查与后续版本提示。
 * 对外接口：
 * - CLI 执行入口：node scripts/postbuild.js
 */
import { execSync } from 'node:child_process'
import { existsSync, readFileSync, writeFileSync } from 'node:fs'
import { dirname, isAbsolute, join, resolve } from 'node:path'
import { fileURLToPath } from 'node:url'

const __dirname = dirname(fileURLToPath(import.meta.url))
const rootDir = join(__dirname, '..')
const rawDistDir = process.env.DOCMATHER_FRONTEND_OUT_DIR?.trim() || 'dist'
const distDir = isAbsolute(rawDistDir) ? rawDistDir : resolve(rootDir, rawDistDir)
const swPath = join(distDir, 'sw.js')
const versionManifestPath = join(distDir, 'version.json')
const pkgPath = join(rootDir, 'package.json')
const buildCounterPath = join(__dirname, '.build-counter')
const PLACEHOLDER = '__SW_BUILD_VERSION__'
const UNKNOWN_SHA = 'nogit'

const MANUAL_BUILD_VERSION_KEYS = ['APP_BUILD_VERSION', 'VITE_APP_BUILD_VERSION']
const CI_BUILD_SERIAL_KEYS = [
  'CI_PIPELINE_IID',
  'CI_PIPELINE_ID',
  'GITHUB_RUN_NUMBER',
  'BUILD_NUMBER',
  'CI_BUILD_ID',
]
const COMMIT_SHA_KEYS = [
  'CI_COMMIT_SHA',
  'CI_COMMIT_SHORT_SHA',
  'GITHUB_SHA',
  'VERCEL_GIT_COMMIT_SHA',
]

function firstNonEmptyEnv(keys) {
  for (const key of keys) {
    const value = process.env[key]
    if (typeof value === 'string' && value.trim()) {
      return value.trim()
    }
  }

  return null
}

function normalizeSegment(rawValue, fallbackValue) {
  const normalized = String(rawValue ?? '')
    .trim()
    .toLowerCase()
    .replace(/[^a-z0-9._-]/g, '-')
    .replace(/-+/g, '-')
    .replace(/^[.-]+|[.-]+$/g, '')

  return normalized || fallbackValue
}

function normalizeBuildVersion(rawValue, fallbackValue) {
  const sanitized = String(rawValue ?? '')
    .trim()
    .replace(/\s+/g, '')
    .replace(/[^0-9a-zA-Z.+_-]/g, '-')

  if (!sanitized) {
    return fallbackValue
  }

  return sanitized.startsWith('v') ? sanitized : `v${sanitized}`
}

function getShortSha(rawSha) {
  return normalizeSegment(String(rawSha ?? '').slice(0, 12), UNKNOWN_SHA)
}

function readPackageVersion() {
  const pkg = JSON.parse(readFileSync(pkgPath, 'utf-8'))
  return normalizeSegment(pkg.version ?? '0.0.0', '0.0.0')
}

function nextBuildNumber(counterPath) {
  // 第一次构建不存在计数器文件时，从 1 开始
  if (!existsSync(counterPath)) {
    return 1
  }

  const rawValue = readFileSync(counterPath, 'utf-8').trim()
  const parsedValue = Number.parseInt(rawValue, 10)

  if (Number.isNaN(parsedValue) || parsedValue < 0) {
    return 1
  }

  return parsedValue + 1
}

function getGitCommitSha() {
  const envSha = firstNonEmptyEnv(COMMIT_SHA_KEYS)
  if (envSha) {
    return normalizeSegment(envSha, UNKNOWN_SHA)
  }

  try {
    return normalizeSegment(
      execSync('git rev-parse HEAD', {
        cwd: rootDir,
        stdio: ['ignore', 'pipe', 'ignore'],
      })
        .toString()
        .trim(),
      UNKNOWN_SHA,
    )
  } catch {
    return UNKNOWN_SHA
  }
}

function resolveBuildSerial() {
  const ciSerial = firstNonEmptyEnv(CI_BUILD_SERIAL_KEYS)
  if (ciSerial) {
    return {
      serial: normalizeSegment(ciSerial, '0'),
      source: 'ci',
    }
  }

  // 本地开发环境兜底：每次构建自动 +1，并持久化到 scripts/.build-counter
  const buildNumber = nextBuildNumber(buildCounterPath)
  writeFileSync(buildCounterPath, String(buildNumber), 'utf-8')

  return {
    serial: `local.${buildNumber}`,
    source: 'local-counter',
  }
}

function resolveBuildMeta() {
  const packageVersion = readPackageVersion()
  const commitSha = getGitCommitSha()
  const commitShortSha = getShortSha(commitSha)
  const manualBuildVersion = firstNonEmptyEnv(MANUAL_BUILD_VERSION_KEYS)

  if (manualBuildVersion) {
    return {
      buildVersion: normalizeBuildVersion(manualBuildVersion, `v${packageVersion}`),
      packageVersion,
      buildSerial: 'manual',
      commitSha,
      commitShortSha,
      source: 'manual',
      builtAt: new Date().toISOString(),
    }
  }

  const { serial, source } = resolveBuildSerial()

  return {
    buildVersion: `v${packageVersion}+${serial}.${commitShortSha}`,
    packageVersion,
    buildSerial: serial,
    commitSha,
    commitShortSha,
    source,
    builtAt: new Date().toISOString(),
  }
}

function injectVersionToServiceWorker(buildVersion) {
  if (!existsSync(swPath)) {
    console.warn(`[postbuild] File not found, skip SW version injection: ${swPath}`)
    return false
  }

  const swContent = readFileSync(swPath, 'utf-8')

  // 若占位符不存在，则保留原文件并输出告警
  if (!swContent.includes(PLACEHOLDER)) {
    console.warn(`[postbuild] Placeholder not found in sw.js: ${PLACEHOLDER}`)
    return false
  }

  const nextContent = swContent.replace(new RegExp(PLACEHOLDER, 'g'), buildVersion)
  writeFileSync(swPath, nextContent, 'utf-8')

  return true
}

function writeVersionManifest(buildMeta) {
  writeFileSync(versionManifestPath, `${JSON.stringify(buildMeta, null, 2)}\n`, 'utf-8')
}

function main() {
  const buildMeta = resolveBuildMeta()
  const injected = injectVersionToServiceWorker(buildMeta.buildVersion)
  writeVersionManifest(buildMeta)

  if (injected) {
    console.log(`[postbuild] Service Worker version injected: ${buildMeta.buildVersion}`)
  }

  console.log(`[postbuild] Version manifest generated: ${versionManifestPath}`)
  console.log(`[postbuild] Build source: ${buildMeta.source}`)
}

main()
