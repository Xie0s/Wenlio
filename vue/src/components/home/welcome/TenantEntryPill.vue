<!--
  src/components/home/welcome/TenantEntryPill.vue
  职责：首页入口胶囊组件，输入租户 ID 后跳转租户门户，或直接进入管理后台
-->
<script setup lang="ts">
import { computed, ref, watch } from 'vue'
import { useRouter } from 'vue-router'
import { ArrowUpRight } from 'lucide-vue-next'
import { toast } from 'vue-sonner'
import Tabs from '@/components/ui/tabs/Tabs.vue'
import TabsList from '@/components/ui/tabs/TabsList.vue'
import TabsTrigger from '@/components/ui/tabs/TabsTrigger.vue'

type Mode = 'home' | 'admin'

const router = useRouter()
const mode = ref<Mode>('home')
const tenantId = ref('')
const inputError = ref(false)
const focused = ref(false)
const hasTenantInput = computed(() => tenantId.value.trim().length > 0)
const goButtonActive = computed(() => mode.value === 'admin' || hasTenantInput.value)
const tenantIdRuleText = '租户 ID 格式错误：仅支持小写字母、数字和中划线，且长度至少 3 位。'

watch(mode, () => {
  tenantId.value = ''
  inputError.value = false
})

function handleGo() {
  if (mode.value === 'admin') {
    router.push('/admin/login')
    return
  }
  const id = tenantId.value.trim().toLowerCase()
  if (!id || !/^[a-z0-9][a-z0-9-]+[a-z0-9]$/.test(id)) {
    inputError.value = true
    toast.error(tenantIdRuleText)
    return
  }
  inputError.value = false
  router.push(`/${id}`)
}

function handleInput() {
  if (inputError.value) inputError.value = false
}

function handleKeydown(e: KeyboardEvent) {
  if (e.key === 'Enter') handleGo()
}
</script>

<template>
  <div class="inline-flex flex-col items-start gap-1.5">
    <div
      class="inline-flex h-16 w-[420px] items-center gap-1 rounded-full border bg-transparent px-2 transition-colors"
      :class="inputError ? 'border-red-400 dark:border-red-400' : 'border-black/15 dark:border-white/20'"
    >
      <!-- 左侧文字区（固定宽度，内容淡入淡出） -->
      <div
        class="relative flex h-11 flex-1 items-center overflow-hidden rounded-full transition-colors"
        :class="mode !== 'home'
          ? 'border-transparent'
          : inputError
            ? 'border border-red-400 dark:border-red-400'
            : focused
              ? 'border border-black/45 dark:border-white/45'
              : 'border border-black/20 dark:border-white/30'"
      >
      <Transition
        enter-active-class="transition-opacity duration-200 ease-out"
        enter-from-class="opacity-0"
        enter-to-class="opacity-100"
        leave-active-class="transition-opacity duration-150 ease-in absolute inset-0"
        leave-from-class="opacity-100"
        leave-to-class="opacity-0"
      >
        <input
          v-if="mode === 'home'"
          v-model="tenantId"
          type="text"
          placeholder="输入租户 ID"
          class="h-full w-full bg-transparent text-base text-foreground placeholder:text-muted-foreground outline-none pl-4 pr-3"
          :class="inputError ? 'placeholder:text-red-400' : ''"
          autocomplete="off"
          spellcheck="false"
          @input="handleInput"
          @keydown="handleKeydown"
          @focus="focused = true"
          @blur="focused = false"
        />
      </Transition>

      <Transition
        enter-active-class="transition-opacity duration-200 ease-out"
        enter-from-class="opacity-0"
        enter-to-class="opacity-100"
        leave-active-class="transition-opacity duration-150 ease-in absolute inset-0"
        leave-from-class="opacity-100"
        leave-to-class="opacity-0"
      >
        <span
          v-if="mode === 'admin'"
          class="flex h-full items-center w-full pl-4 pr-3 text-base text-muted-foreground select-none truncate"
        >
          直接点击前进跳转管理后台
        </span>
      </Transition>
      </div>

      <!-- 模式切换 Tabs -->
      <Tabs v-model="mode" class="flex-row gap-0 items-center shrink-0">
        <TabsList class="h-11">
          <TabsTrigger
            value="home"
          class="px-4 text-sm data-[state=active]:bg-black data-[state=active]:text-white dark:data-[state=active]:bg-foreground/[0.07] dark:data-[state=active]:text-foreground"
          >
            首页
          </TabsTrigger>
          <TabsTrigger
            value="admin"
          class="px-4 text-sm data-[state=active]:bg-black data-[state=active]:text-white dark:data-[state=active]:bg-foreground/[0.07] dark:data-[state=active]:text-foreground"
          >
            管理
          </TabsTrigger>
        </TabsList>
      </Tabs>

      <!-- 前进按钮 -->
      <button
        type="button"
        class="ml-1 flex h-11 w-11 shrink-0 items-center justify-center rounded-full border transition-all active:scale-95"
        :class="[
          inputError ? 'border-red-400 dark:border-red-400' : 'border-black/15 dark:border-white/20',
          goButtonActive
            ? 'bg-primary text-primary-foreground hover:opacity-90'
            : 'text-foreground hover:bg-foreground hover:text-background',
        ]"
        @click="handleGo"
      >
        <ArrowUpRight class="h-5 w-5" />
      </button>
    </div>

    <p v-if="mode === 'home' && inputError" class="pl-3 text-xs text-red-500">
      {{ tenantIdRuleText }}
    </p>
  </div>
</template>
