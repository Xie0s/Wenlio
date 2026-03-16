/**
 * 文件：public/sw.js
 * 作用：Service Worker 脚本，负责字体资源运行时缓存。
 * 职责边界：
 * 1) 基于构建版本号维护字体缓存命名空间。
 * 2) 在 fetch 阶段对字体资源执行缓存优先策略。
 * 对外接口：
 * - 浏览器 Service Worker 生命周期事件：install / activate / fetch。
 */

// 该占位符会在 `vite build` 后由 scripts/postbuild.js 替换。
const CACHE_VERSION =
  '__SW_BUILD_VERSION__' !== '__SW_BUILD_' + 'VERSION__'
    ? '__SW_BUILD_VERSION__'
    : `v${Date.now()}`

const FONT_CACHE = `font-cache-${CACHE_VERSION}`
const FONT_FILE_PATTERN = /\.(woff2?|ttf|otf|eot)$/i

self.addEventListener('install', (event) => {
  event.waitUntil(self.skipWaiting())
})

self.addEventListener('activate', (event) => {
  event.waitUntil(
    caches
      .keys()
      .then((keys) =>
        Promise.all(
          // 仅保留当前版本的字体缓存，清理旧版本缓存
          keys
            .filter((key) => key.startsWith('font-cache-') && key !== FONT_CACHE)
            .map((key) => caches.delete(key)),
        ),
      )
      .then(() => self.clients.claim()),
  )
})

self.addEventListener('fetch', (event) => {
  const { request } = event

  if (request.method !== 'GET') {
    return
  }

  const requestUrl = new URL(request.url)
  const isFontRequest =
    request.destination === 'font' || FONT_FILE_PATTERN.test(requestUrl.pathname)

  if (!isFontRequest) {
    return
  }

  // 字体请求走缓存优先，减少重复下载
  event.respondWith(cacheFirstFont(request))
})

async function cacheFirstFont(request) {
  const cache = await caches.open(FONT_CACHE)
  const cachedResponse = await cache.match(request)

  if (cachedResponse) {
    return cachedResponse
  }

  const networkResponse = await fetch(request)

  if (networkResponse.ok) {
    await cache.put(request, networkResponse.clone())
  }

  return networkResponse
}
