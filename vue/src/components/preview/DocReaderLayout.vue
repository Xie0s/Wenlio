<!-- DocReaderLayout.vue - 文档阅读界面布局组件
     职责：提供文档阅读的顶部导航栏（面包屑、搜索、主题切换、首页按钮）和内容区域插槽
     对外接口：
       Props: tenantId
       Emits: navigate(path: string)
       Slots: default（内容区域） -->
<script setup lang="ts">
import { ref, computed, onMounted, watch, toRef, provide } from 'vue'
import { useRouter } from 'vue-router'
import { applyTenantBrowserBranding } from '@/lib/browser-branding'
import { useReaderStore } from '@/stores/reader'
import { ChevronRight } from 'lucide-vue-next'
import { SearchDialog } from '@/components/search'
import UserMenu from '@/components/auth/UserMenu.vue'
import { useAuthStore } from '@/stores/auth'
import { Badge } from '@/components/ui/badge'
import { ACCESS_MODE_LABEL, ACCESS_MODE_COLOR } from '@/utils/types'
import type { AccessMode } from '@/utils/types'

const props = defineProps<{
  tenantId: string
}>()

const emit = defineEmits<{
  navigate: [path: string]
}>()

const store = useReaderStore()
const router = useRouter()
const authStore = useAuthStore()
const tenantIdRef = toRef(props, 'tenantId')

const currentSection = computed(() => {
  if (!store.currentPage) return null
  return store.tree.find(s => s.pages.some(p => p.id === store.currentPage!.id)) ?? null
})

const showSearch = ref(false)

async function loadTenant() {
  if (!tenantIdRef.value) return
  if (store.tenant?.id === tenantIdRef.value) return
  try {
    await store.loadTenant(tenantIdRef.value)
  } catch {
    // 网络异常（fetch 本身抛错），兜底跳转 404
    router.replace({ name: 'NotFound' })
    return
  }
  // http.get 对 404 不抛出，需额外判断：租户加载后仍为 null 说明租户不存在
  if (!store.tenant) {
    router.replace({ name: 'NotFound' })
  }
}

function openSearch() {
  showSearch.value = true
}

function goTenantHome() {
  emit('navigate', `/${tenantIdRef.value}`)
}

function goToThemes() {
  emit('navigate', `/${tenantIdRef.value}/themes`)
}

provide('openSearch', openSearch)
provide('goTenantHome', goTenantHome)

onMounted(loadTenant)
watch(tenantIdRef, loadTenant)
watch(
  () => [
    store.tenant?.name ?? '',
    store.tenant?.browser_title ?? '',
    store.tenant?.browser_icon_url ?? '',
  ],
  () => {
    applyTenantBrowserBranding({
      tenantName: store.tenant?.name,
      browserTitle: store.tenant?.browser_title,
      browserIconUrl: store.tenant?.browser_icon_url,
    })
  },
  { immediate: true },
)
</script>

<template>
  <div class="h-screen flex flex-col overflow-hidden bg-background">
    <!-- 顶部导航栏 -->
    <header class="shrink-0 z-40 border-b bg-background">
      <div class="flex h-14 w-full items-center gap-2 px-4">

        <!-- 桌面端：完整面包屑（lg+） -->
        <div class="hidden lg:flex items-center gap-1.5 min-w-0">
          <a class="flex items-center gap-2 shrink-0 hover:opacity-80 transition-opacity cursor-pointer"
            @click.prevent="goTenantHome">
            <img v-if="store.tenant?.logo_url" :src="store.tenant.logo_url" :alt="store.tenant.name"
              class="h-8 w-8 rounded-lg object-cover reader-layout-img" draggable="false" @contextmenu.prevent
              @dragstart.prevent />
            <span class="text-xl font-medium">{{ store.tenantName }}</span>
          </a>
          <template v-if="store.currentTheme">
            <ChevronRight class="h-5 w-5 shrink-0 text-muted-foreground" />
            <span class="text-base text-muted-foreground truncate">{{ store.currentTheme.name }}</span>
            <Badge v-if="store.currentTheme.access_mode && store.currentTheme.access_mode !== 'public'"
              variant="outline" class="shrink-0 text-[10px] px-1.5 py-0 rounded-full border-current"
              :style="{ color: ACCESS_MODE_COLOR[store.currentTheme.access_mode as AccessMode] }">
              {{ ACCESS_MODE_LABEL[store.currentTheme.access_mode as AccessMode] }}
            </Badge>
          </template>
          <template v-if="store.currentVersion">
            <ChevronRight class="h-5 w-5 shrink-0 text-muted-foreground" />
            <span class="text-base text-muted-foreground truncate">{{ store.currentVersion.label ||
              store.currentVersion.name }}</span>
          </template>
          <template v-if="currentSection">
            <ChevronRight class="h-5 w-5 shrink-0 text-muted-foreground" />
            <span class="text-base text-muted-foreground truncate">{{ currentSection.title }}</span>
          </template>
          <template v-if="store.currentPage">
            <ChevronRight class="h-5 w-5 shrink-0 text-muted-foreground" />
            <span class="text-base text-foreground truncate font-medium">{{ store.currentPage.title }}</span>
          </template>
        </div>

        <!-- 移动端：Logo + 当前章节名（< lg） -->
        <a class="flex lg:hidden items-center gap-2.5 min-w-0 flex-1 cursor-pointer hover:opacity-80 transition-opacity"
          @click.prevent="goTenantHome">
          <img v-if="store.tenant?.logo_url" :src="store.tenant.logo_url" :alt="store.tenant.name"
            class="h-8 w-8 shrink-0 rounded-lg object-cover reader-layout-img" draggable="false" @contextmenu.prevent
            @dragstart.prevent />
          <span class="text-base font-medium truncate">
            {{ store.currentPage?.title || store.tenantName }}
          </span>
        </a>

        <!-- 桌面端弹性占位 -->
        <div class="hidden lg:block flex-1" />

        <!-- 移动端：汉堡菜单注入点（由 DocMobileToolbar Teleport 填充） -->
        <div id="reader-mobile-actions" class="flex items-center lg:hidden" />

        <!-- 桌面端：操作按钮组 -->
        <div class="hidden lg:flex items-center gap-2">
          <button v-if="store.currentTheme"
            class="rounded-full border border-foreground/25 bg-transparent text-foreground px-4 py-2.5 text-sm font-medium transition-colors hover:bg-foreground/5"
            @click="goToThemes">
            前往文档主题列表
          </button>
          <UserMenu v-if="authStore.isLoggedIn" :tenant-id="tenantId" />
        </div>

      </div>
    </header>

    <!-- 内容区 -->
    <main id="reader-scroll-container" class="scrollbar-visible flex-1 min-h-0 overflow-y-auto w-full">
      <slot />
    </main>

    <!-- 搜索弹窗 -->
    <SearchDialog v-model:open="showSearch" :tenant-id="tenantId" @navigate="emit('navigate', $event)" />
  </div>
</template>

<style scoped>
.reader-layout-img {
  user-select: none;
  -webkit-user-drag: none;
}
</style>
