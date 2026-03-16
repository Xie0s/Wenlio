<!--
  CodeBlockMenu.vue - 代码块浮动菜单组件
  职责：当光标在代码块内时显示浮动工具栏，提供语言选择和删除操作
  对外暴露：Props: editor (Tiptap Editor 实例)
-->
<script setup lang="ts">
import { ref, computed } from 'vue'
import { Trash2, Check, ChevronDown, Search } from 'lucide-vue-next'
import { useMenuScheduler, isRectInView, isTopInView } from './useMenuScheduler'
import type { Editor } from '@tiptap/vue-3'
import { Button } from '@/components/ui/button'
import { Tooltip, TooltipTrigger, TooltipContent } from '@/components/ui/tooltip'
import { Popover, PopoverContent, PopoverTrigger } from '@/components/ui/popover'

const props = defineProps<{ editor: Editor }>()

const LANGUAGES = [
  { value: 'plain', label: '纯文本' },
  { value: 'javascript', label: 'JavaScript' },
  { value: 'typescript', label: 'TypeScript' },
  { value: 'vue', label: 'Vue' },
  { value: 'html', label: 'HTML' },
  { value: 'css', label: 'CSS' },
  { value: 'scss', label: 'SCSS' },
  { value: 'less', label: 'Less' },
  { value: 'json', label: 'JSON' },
  { value: 'yaml', label: 'YAML' },
  { value: 'toml', label: 'TOML' },
  { value: 'ini', label: 'INI' },
  { value: 'markdown', label: 'Markdown' },
  { value: 'bash', label: 'Bash' },
  { value: 'shell', label: 'Shell' },
  { value: 'powershell', label: 'PowerShell' },
  { value: 'go', label: 'Go' },
  { value: 'python', label: 'Python' },
  { value: 'java', label: 'Java' },
  { value: 'kotlin', label: 'Kotlin' },
  { value: 'swift', label: 'Swift' },
  { value: 'rust', label: 'Rust' },
  { value: 'cpp', label: 'C++' },
  { value: 'c', label: 'C' },
  { value: 'csharp', label: 'C#' },
  { value: 'ruby', label: 'Ruby' },
  { value: 'php', label: 'PHP' },
  { value: 'sql', label: 'SQL' },
  { value: 'graphql', label: 'GraphQL' },
  { value: 'docker', label: 'Dockerfile' },
  { value: 'nginx', label: 'Nginx' },
  { value: 'xml', label: 'XML' },
  { value: 'diff', label: 'Diff' },
  { value: 'makefile', label: 'Makefile' },
  { value: 'proto', label: 'Protobuf' },
]

const isVisible = ref(false)
const position = ref({ top: 0, left: 0 })
const menuRef = ref<HTMLDivElement | null>(null)
const currentLanguage = ref('plain')
const popoverOpen = ref(false)
const searchQuery = ref('')

const currentLabel = computed(() =>
  LANGUAGES.find(l => l.value === currentLanguage.value)?.label ?? '纯文本',
)

const filteredLanguages = computed(() => {
  const q = searchQuery.value.toLowerCase().trim()
  if (!q) return LANGUAGES
  return LANGUAGES.filter(l =>
    l.label.toLowerCase().includes(q) || l.value.toLowerCase().includes(q),
  )
})

const lastState = { isVisible: false, top: 0, left: 0 }
let cachedMenuWidth = 0

function applyState(nextVisible: boolean, nextTop: number, nextLeft: number) {
  const roundedTop = Math.round(nextTop)
  const roundedLeft = Math.round(nextLeft)
  if (lastState.isVisible !== nextVisible) {
    lastState.isVisible = nextVisible
    isVisible.value = nextVisible
  }
  if (!nextVisible) return
  if (lastState.top !== roundedTop || lastState.left !== roundedLeft) {
    lastState.top = roundedTop
    lastState.left = roundedLeft
    position.value = { top: roundedTop, left: roundedLeft }
  }
}

function computeMenu() {
  const editor = props.editor
  if (!editor.isActive('codeBlock')) {
    if (popoverOpen.value) return
    applyState(false, 0, 0)
    return
  }
  if (popoverOpen.value) return

  if (menuRef.value) {
    const w = menuRef.value.offsetWidth
    if (w > 0) cachedMenuWidth = w
  }

  const attrs = editor.getAttributes('codeBlock')
  currentLanguage.value = attrs?.language || 'plain'

  try {
    const { view, state } = editor
    const domAtPos = view.domAtPos(state.selection.$anchor.pos)
    const node = domAtPos.node
    const el = (node.nodeType === Node.TEXT_NODE ? node.parentElement : node) as HTMLElement
    const preEl = el?.closest('pre')
    if (!preEl) { applyState(false, 0, 0); return }

    const preRect = preEl.getBoundingClientRect()
    const vp = getViewport()
    if (!isRectInView(preRect, vp)) { applyState(false, 0, 0); return }

    let top = preRect.top - 45
    if (top < 8) top = preRect.bottom + 8
    if (!isTopInView(top, vp)) { applyState(false, 0, 0); return }

    const centerLeft = preRect.left + preRect.width / 2
    const menuW = cachedMenuWidth
    if (menuW > 0) {
      const left = Math.round(Math.max(12, Math.min(window.innerWidth - menuW - 12, centerLeft - menuW / 2)))
      applyState(true, top, left)
    } else {
      applyState(true, top, Math.round(Math.max(12, centerLeft)))
      window.requestAnimationFrame(computeMenu)
    }
  } catch {
    applyState(false, 0, 0)
  }
}

const { getViewport } = useMenuScheduler(props.editor, computeMenu)

function selectLanguage(value: string) {
  currentLanguage.value = value
  props.editor.chain().focus().updateAttributes('codeBlock', { language: value === 'plain' ? '' : value }).run()
  popoverOpen.value = false
  searchQuery.value = ''
}

function deleteCodeBlock() {
  props.editor.chain().focus().clearNodes().run()
}
</script>

<template>
  <div v-if="isVisible" ref="menuRef" class="glass fixed z-50 flex items-center gap-0.5 p-1.5 rounded-full antialiased"
    :style="{ top: `${position.top}px`, left: `${position.left}px` }">
    <!-- 语言选择触发器 -->
    <Popover v-model:open="popoverOpen">
      <PopoverTrigger as-child>
        <button
          class="flex items-center gap-1 h-8 px-3 text-xs font-medium rounded-full transition-colors hover:bg-foreground/10">
          {{ currentLabel }}
          <ChevronDown class="h-3 w-3 opacity-50" />
        </button>
      </PopoverTrigger>
      <PopoverContent align="start" :side-offset="8" class="w-48 p-0 rounded-2xl">
        <!-- 搜索栏 -->
        <div class="flex items-center gap-2 px-3 py-2 border-b">
          <Search class="h-3.5 w-3.5 text-muted-foreground shrink-0" />
          <input v-model="searchQuery"
            class="flex-1 text-xs bg-transparent outline-none placeholder:text-muted-foreground"
            placeholder="搜索语言..." />
        </div>
        <!-- 语言列表 -->
        <div class="max-h-52 overflow-y-auto p-1">
          <button v-for="lang in filteredLanguages" :key="lang.value"
            class="flex items-center justify-between w-full px-2.5 py-1.5 text-xs rounded-lg transition-colors hover:bg-accent"
            :class="currentLanguage === lang.value && 'text-primary font-medium'" @click="selectLanguage(lang.value)">
            {{ lang.label }}
            <Check v-if="currentLanguage === lang.value" class="h-3 w-3" />
          </button>
          <div v-if="filteredLanguages.length === 0" class="px-2.5 py-4 text-center text-xs text-muted-foreground">
            未找到匹配语言
          </div>
        </div>
      </PopoverContent>
    </Popover>

    <div class="w-px h-5 bg-border" />

    <!-- 删除代码块 -->
    <Tooltip>
      <TooltipTrigger as-child>
        <Button variant="ghost" size="icon" class="h-8 w-8 rounded-full hover:bg-destructive/10 hover:text-destructive"
          @click="deleteCodeBlock">
          <Trash2 class="h-4 w-4" />
        </Button>
      </TooltipTrigger>
      <TooltipContent>删除代码块</TooltipContent>
    </Tooltip>
  </div>
</template>
