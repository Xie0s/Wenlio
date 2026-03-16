<!--
  hero-demo/HeroDemo.vue
  职责：Hero 区域演示窗口，包含浏览器 Chrome 外壳、右上角场景切换标签
  子组件：ViewerDemo（查看场景）、EditorDemo（编辑场景）
-->
<script setup lang="ts">
import { ref } from 'vue'
import ViewerDemo from './ViewerDemo.vue'
import EditorDemo from './EditorDemo.vue'
import { BookOpen, PencilLine } from 'lucide-vue-next'

type Scene = 'viewer' | 'editor'

const activeScene = ref<Scene>('viewer')

const scenes = [
  { id: 'viewer' as Scene, label: '查看',  icon: BookOpen  },
  { id: 'editor' as Scene, label: '编辑',  icon: PencilLine },
]

const urlMap: Record<Scene, string> = {
  viewer: 'docs.wenlio.com/guide/quickstart',
  editor: 'admin.wenlio.com/editor/user-guide',
}
</script>

<template>
  <div
    class="relative mx-auto w-full max-w-none overflow-hidden rounded-xl
           border border-black/[0.12] dark:border-white/[0.12]
           bg-white dark:bg-[#1e1e1e]"
  >

    <!-- 浏览器 Chrome 顶栏 -->
    <div
      class="flex h-10 items-center gap-3 shrink-0
             border-b border-black/[0.06] dark:border-white/[0.06]
             bg-[#f5f5f5] dark:bg-[#2a2a2a] px-4"
    >

      <!-- 交通灯按钮 -->
      <div class="flex gap-1.5 shrink-0">
        <span class="h-2.5 w-2.5 rounded-full bg-[#ff5f57]" />
        <span class="h-2.5 w-2.5 rounded-full bg-[#febc2e]" />
        <span class="h-2.5 w-2.5 rounded-full bg-[#28c840]" />
      </div>

      <!-- URL 地址栏 -->
      <div
        class="flex-1 flex items-center justify-center"
      >
        <div
          class="flex h-5.5 items-center rounded-md bg-black/[0.05] dark:bg-white/[0.07]
                 px-3 text-[11px] text-muted-foreground/60 select-none
                 min-w-0 max-w-[260px] overflow-hidden"
        >
          <span class="truncate">{{ urlMap[activeScene] }}</span>
        </div>
      </div>

      <!-- 场景切换标签（右上角） -->
      <div class="flex items-center gap-1 shrink-0">
        <button
          v-for="scene in scenes"
          :key="scene.id"
          class="flex items-center gap-1 px-2.5 h-6 rounded-full text-[11px] font-medium transition-colors select-none"
          :class="activeScene === scene.id
            ? 'bg-primary text-primary-foreground'
            : 'text-muted-foreground hover:bg-black/[0.07] dark:hover:bg-white/[0.1]'"
          @click="activeScene = scene.id"
        >
          <component :is="scene.icon" class="h-3 w-3 shrink-0" />
          {{ scene.label }}
        </button>
      </div>

    </div>

    <!-- 演示内容区 -->
    <div class="h-[720px] relative">
      <Transition
        enter-active-class="transition-opacity duration-250 ease-out"
        enter-from-class="opacity-0"
        enter-to-class="opacity-100"
        leave-active-class="transition-opacity duration-150 ease-in absolute inset-0"
        leave-from-class="opacity-100"
        leave-to-class="opacity-0"
        mode="out-in"
      >
        <ViewerDemo v-if="activeScene === 'viewer'" class="h-full" />
        <EditorDemo v-else class="h-full" />
      </Transition>
    </div>

  </div>
</template>
