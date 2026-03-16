<!--
  VersionUpdatePrompt.vue - 版本更新提示组件
  职责：
  1) 监听 app:version-update 事件并展示可见提示。
  2) 提供“立即刷新”与“稍后处理”交互闭环。
  3) 仅负责 UI 与用户交互，不耦合版本探针实现细节。
  对外接口：
  - Props：无
  - Emits：无
  - 依赖：@/lib/version 的 VERSION_UPDATE_EVENT
-->
<script setup lang="ts">
import { computed, onMounted, onUnmounted, ref } from 'vue'
import { RefreshCw, Sparkles, X } from 'lucide-vue-next'
import { Button } from '@/components/ui/button'
import { Tooltip, TooltipContent, TooltipTrigger } from '@/components/ui/tooltip'
import { VERSION_UPDATE_EVENT, type BuildVersionUpdateDetail } from '@/lib/version'

const isVisible = ref(false)
const isRefreshing = ref(false)
const latestVersion = ref('')
const previousVersion = ref('')

const latestVersionLabel = computed(() => latestVersion.value || '未知版本')
const previousVersionLabel = computed(() => previousVersion.value || '未知版本')

const handleVersionUpdate: EventListener = (event) => {
  const customEvent = event as CustomEvent<BuildVersionUpdateDetail>
  const detail = customEvent.detail

  if (!detail?.latest?.buildVersion) {
    return
  }

  latestVersion.value = detail.latest.buildVersion
  previousVersion.value = detail.previousVersion
  isRefreshing.value = false
  isVisible.value = true
}

function dismissPrompt(): void {
  isVisible.value = false
}

function refreshNow(): void {
  if (isRefreshing.value) {
    return
  }

  isRefreshing.value = true
  window.location.reload()
}

onMounted(() => {
  window.addEventListener(VERSION_UPDATE_EVENT, handleVersionUpdate)
})

onUnmounted(() => {
  window.removeEventListener(VERSION_UPDATE_EVENT, handleVersionUpdate)
})
</script>

<template>
  <Transition name="version-update-prompt">
    <section
      v-if="isVisible"
      class="fixed right-4 top-4 z-[80] w-[min(92vw,26rem)]"
      aria-live="polite"
      aria-atomic="true"
    >
      <div class="rounded-2xl border border-border/80 bg-card/95 p-4 text-card-foreground shadow-xl backdrop-blur supports-[backdrop-filter]:bg-card/85">
        <div class="flex items-start gap-3">
          <div class="mt-0.5 flex size-9 shrink-0 items-center justify-center rounded-full bg-primary/12 text-primary">
            <Sparkles class="size-4" />
          </div>

          <div class="min-w-0 flex-1 space-y-1">
            <p class="text-sm font-semibold leading-none">发现新版本，立即刷新</p>
            <p class="text-xs text-muted-foreground">当前 {{ previousVersionLabel }} → 最新 {{ latestVersionLabel }}</p>
          </div>

          <div class="flex shrink-0 items-center gap-2">
            <Tooltip>
              <TooltipTrigger as-child>
                <Button
                  size="icon"
                  class="rounded-full"
                  aria-label="立即刷新"
                  @click="refreshNow"
                >
                  <RefreshCw class="size-4" :class="{ 'animate-spin': isRefreshing }" />
                </Button>
              </TooltipTrigger>
              <TooltipContent>
                <p>立即刷新</p>
              </TooltipContent>
            </Tooltip>

            <Tooltip>
              <TooltipTrigger as-child>
                <Button
                  variant="outline"
                  size="icon"
                  class="rounded-full"
                  aria-label="稍后处理"
                  @click="dismissPrompt"
                >
                  <X class="size-4" />
                </Button>
              </TooltipTrigger>
              <TooltipContent>
                <p>稍后处理</p>
              </TooltipContent>
            </Tooltip>
          </div>
        </div>
      </div>
    </section>
  </Transition>
</template>

<style scoped>
.version-update-prompt-enter-active,
.version-update-prompt-leave-active {
  transition: opacity 0.22s ease, transform 0.22s ease;
}

.version-update-prompt-enter-from,
.version-update-prompt-leave-to {
  opacity: 0;
  transform: translate3d(0, -8px, 0) scale(0.98);
}
</style>
