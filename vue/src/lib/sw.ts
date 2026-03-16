/**
 * src/lib/sw.ts
 * What: Frontend bootstrap helper for Service Worker registration.
 * Responsibility boundary:
 * 1) Register /sw.js only in production mode.
 * 2) Keep registration failure isolated from app bootstrap.
 * Exposed interface:
 * - registerServiceWorker(): void
 */

export function registerServiceWorker(): void {
  if (!('serviceWorker' in navigator)) {
    return
  }

  if (!import.meta.env.PROD) {
    return
  }

  window.addEventListener('load', () => {
    navigator.serviceWorker.register('/sw.js').catch((error) => {
      console.warn('[sw] registration failed:', error)
    })
  })
}
