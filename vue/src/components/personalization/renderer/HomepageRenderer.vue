<!-- HomepageRenderer.vue - 首页配置驱动渲染器
     职责：接收 HomepageLayout 配置，按 sections 顺序动态渲染各区块组件。
           未配置时（layout 为 null）回退到默认首页。
     对外接口：
       Props: layout, tenantId, tenantName, tenantLogo, themes
       Emits: navigate(url: string) -->
<script setup lang="ts">
import { computed } from 'vue'
import type { Theme } from '@/utils/types'
import type {
  HomepageLayout,
  HomepageSection,
  NavbarConfig,
  HeroConfig,
  IntroductionConfig,
  ThemeListConfig,
  CtaConfig,
  FooterConfig,
} from '../types'
import PNavbar from '../navbar/PNavbar.vue'
import PHero from '../hero/PHero.vue'
import PIntroduction from '../introduction/PIntroduction.vue'
import PThemeList from '../theme-list/PThemeList.vue'
import PCta from '../cta/PCta.vue'
import PFooter from '../footer/PFooter.vue'
import { useThemeStore } from '@/stores/theme'

const props = defineProps<{
  layout: HomepageLayout
  tenantId: string
  tenantName: string
  tenantLogo: string
  themes: Theme[]
}>()

const emit = defineEmits<{
  navigate: [url: string]
}>()

const themeStore = useThemeStore()

const visibleSections = computed(() =>
  props.layout.sections.filter(s => s.visible),
)

const rootStyle = computed(() => {
  const g = props.layout.global
  const vars: Record<string, string> = {}
  if (g.max_width) vars['--hp-max-width'] = g.max_width
  if (g.font_family) vars.fontFamily = g.font_family
  vars['--hp-section-py'] = `${g.section_spacing ?? 96}px`

  const isDark = themeStore.resolvedTheme === 'dark'
  if (isDark) {
    const darkBg = g.dark_background_color || '#1a1a1a'
    vars.backgroundColor = darkBg
    vars['--hp-bg'] = darkBg
  } else {
    const lightBg = g.background_color || '#ffffff'
    if (g.background_color) vars.backgroundColor = g.background_color
    vars['--hp-bg'] = lightBg
  }

  return vars
})

function sectionConfig<T>(section: HomepageSection): T {
  return section.config as T
}

/** 首页图片防护：禁止右键另存、拖拽保存、复制 */
function onContextMenu(e: MouseEvent) {
  if ((e.target as HTMLElement)?.closest?.('img')) e.preventDefault()
}

function onDragStart(e: DragEvent) {
  if ((e.target as HTMLElement)?.closest?.('img')) e.preventDefault()
}

function onCopy(e: ClipboardEvent) {
  const container = (e.target as HTMLElement)?.closest?.('.homepage-renderer-root')
  if (!container) return
  const sel = document.getSelection()
  if (!sel?.rangeCount) return
  const imgs = container.querySelectorAll('img')
  for (let i = 0; i < sel.rangeCount; i++) {
    const range = sel.getRangeAt(i)
    for (const img of imgs) {
      try {
        if (range.intersectsNode(img)) {
          e.preventDefault()
          if (e.clipboardData) e.clipboardData.clearData()
          return
        }
      } catch {
        /* ignore */
      }
    }
  }
}
</script>

<template>
  <div
    :style="rootStyle"
    class="homepage-renderer-root reader-img-protect overflow-y-auto"
    @contextmenu="onContextMenu"
    @dragstart="onDragStart"
    @copy="onCopy"
  >
    <template v-for="section in visibleSections" :key="section.id">
      <PNavbar
        v-if="section.type === 'navbar'"
        :config="sectionConfig<NavbarConfig>(section)"
        :tenant-id="tenantId"
        :tenant-name="tenantName"
        :tenant-logo="tenantLogo"
        @navigate="emit('navigate', $event)"
      />
      <PHero
        v-else-if="section.type === 'hero'"
        :config="sectionConfig<HeroConfig>(section)"
        @navigate="emit('navigate', $event)"
      />
      <PIntroduction
        v-else-if="section.type === 'introduction'"
        :config="sectionConfig<IntroductionConfig>(section)"
      />
      <PThemeList
        v-else-if="section.type === 'theme_list'"
        :config="sectionConfig<ThemeListConfig>(section)"
        :themes="themes"
        :tenant-id="tenantId"
        @navigate="emit('navigate', $event)"
      />
      <PCta
        v-else-if="section.type === 'cta'"
        :config="sectionConfig<CtaConfig>(section)"
        @navigate="emit('navigate', $event)"
      />
      <PFooter
        v-else-if="section.type === 'footer'"
        :config="sectionConfig<FooterConfig>(section)"
        :tenant-name="tenantName"
        :tenant-logo="tenantLogo"
      />
    </template>
  </div>
</template>

<style scoped>
.reader-img-protect :deep(img) {
  user-select: none;
  -webkit-user-select: none;
  -webkit-user-drag: none;
  pointer-events: auto;
}
</style>
