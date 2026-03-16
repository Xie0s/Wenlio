<!--
  EditorSidebar.vue - 编辑器右侧侧边栏
  职责：合并章节/文档导航（含增删改管理）与文章 TOC 到统一的右侧面板
  对外暴露：无 props/emits，通过 inject 消费编辑器上下文
-->
<script setup lang="ts">
import { useEditorInject } from '@/composables/useThemeEditor'
import { useDocTreeCrud } from '@/composables/useDocTreeCrud'
import { slugHint } from '@/lib/validation'
import DocSectionTree from '@/components/editor/edit-page/sidebar/DocSectionTree.vue'
import EditorSidebarTocList from '@/components/editor/edit-page/sidebar/EditorSidebarTocList.vue'

const editor = useEditorInject()
const crud = useDocTreeCrud(editor)
</script>

<template>
  <aside class="w-full h-full overflow-hidden bg-background flex flex-col">

    <!-- 章节导航 + 管理区（弹性占满剩余空间） -->
    <div class="flex-1 min-h-0 overflow-hidden flex flex-col">
      <div v-if="editor.loadingTree.value" class="flex items-center justify-center py-8">
        <div class="h-4 w-4 animate-spin rounded-full border-2 border-primary border-t-transparent" />
      </div>
      <DocSectionTree v-else :sections="editor.sections.value" :active-page-id="editor.activePageId.value"
        :slug-hint="slugHint" :editor="editor" :crud="crud" />
    </div>

    <!-- TOC 区固定在底部，含展开/折叠过渡动画 -->
    <Transition enter-active-class="transition-all duration-300 ease-out overflow-hidden"
      leave-active-class="transition-all duration-200 ease-in overflow-hidden" enter-from-class="opacity-0 max-h-0"
      enter-to-class="opacity-100 max-h-[50vh]" leave-from-class="opacity-100 max-h-[50vh]"
      leave-to-class="opacity-0 max-h-0">
      <div v-if="editor.tocItems.value.length > 0" class="shrink-0 flex flex-col h-[38%] min-h-0">
        <div class="shrink-0 h-px bg-border" />
        <div class="flex-1 min-h-0 overflow-hidden">
          <EditorSidebarTocList :items="editor.tocItems.value" :active-id="editor.tocActiveId.value"
            @scroll-to="editor.tocScrollTo" />
        </div>
      </div>
    </Transition>

  </aside>
</template>
