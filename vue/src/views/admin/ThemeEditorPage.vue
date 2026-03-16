<!-- ThemeEditorPage.vue - 主题集成编辑页
     职责：纯布局壳，创建编辑器 composable 并 provide 给子组件树
     子组件：EditorToolbar、EditorWorkspace/VisualEditorWorkspace、EditorSidebar、EditorSettingsPanel
     对外暴露：路由页面组件 -->
<script setup lang="ts">
import { ref, provide, watchEffect, onMounted, onUnmounted } from 'vue'
import { storeToRefs } from 'pinia'
import { useRoute } from 'vue-router'
import { useThemeEditor, EDITOR_KEY } from '@/composables/useThemeEditor'
import { PAGE_STATUS_LABEL, PAGE_STATUS_COLOR } from '@/utils/types'
import { Badge } from '@/components/ui/badge'
import { useUploadStore } from '@/stores/upload'

import AdminHeader from '@/components/admin/AdminHeader.vue'
import EditorToolbar from '@/components/editor/edit-page/tool/EditorToolbar.vue'
import EditorWorkspace from '@/components/editor/edit-page/EditorWorkspace.vue'
import VisualEditorWorkspace from '@/components/editor/edit-page/visual-editor/VisualEditorWorkspace.vue'
import EditorSidebar from '@/components/editor/edit-page/sidebar/EditorSidebar.vue'
import EditorSettingsPanel from '@/components/editor/edit-page/tool/EditorSettingsPanel.vue'
import { UploadDropZone, UploadProgressPanel, UploadPreviewDialog } from '@/components/editor/media/update'
import { ResizablePanelGroup, ResizablePanel, ResizableHandle } from '@/components/ui/resizable'

const route = useRoute()
const themeId = route.params.themeId as string

const uploadStore = useUploadStore()
const { showPreviewDialog, pendingPreviewFiles, pendingAppendFiles } = storeToRefs(uploadStore)

const editor = useThemeEditor(themeId)
provide(EDITOR_KEY, editor)

const eyeCareMode = ref(false)
provide('eyeCareMode', eyeCareMode)
provide('toggleEyeCare', () => { eyeCareMode.value = !eyeCareMode.value })

watchEffect(() => {
  document.body.classList.toggle('eye-care', eyeCareMode.value)
})

onMounted(() => editor.init())
onUnmounted(() => {
  editor.destroy()
  document.body.classList.remove('eye-care')
})
</script>

<template>
  <div class="h-svh flex flex-col overflow-hidden bg-background">
    <AdminHeader :title-override="editor.theme.value?.name"
      :active-title="editor.page.value ? editor.title.value : undefined" hide-user-info>
      <template #title-suffix>
        <Badge v-if="editor.page.value" :class="PAGE_STATUS_COLOR[editor.page.value.status]" class="text-xs">
          {{ PAGE_STATUS_LABEL[editor.page.value.status] }}
        </Badge>
      </template>
      <template #actions>
        <EditorToolbar />
      </template>
    </AdminHeader>

    <div class="fixed top-16 inset-x-0 h-0 border-b border-border pointer-events-none" style="z-index: 29;" />
    <UploadDropZone class="flex-1 min-h-0">
    <ResizablePanelGroup direction="horizontal" class="h-full pt-16">
      <ResizablePanel :default-size="75" :min-size="50">
        <EditorWorkspace v-if="editor.editorMode.value === 'source'" />
        <VisualEditorWorkspace v-else />
      </ResizablePanel>
      <ResizableHandle class="w-[0.5px]" />
      <ResizablePanel :default-size="25" :min-size="15" :max-size="40">
        <EditorSidebar />
      </ResizablePanel>
    </ResizablePanelGroup>
    </UploadDropZone>

    <EditorSettingsPanel />
    <UploadProgressPanel />
    <UploadPreviewDialog
      v-model:open="showPreviewDialog"
      :initial-files="pendingPreviewFiles"
      :append-files="pendingAppendFiles"
      :auto-open-picker="false"
      title="确认上传文件"
      @confirm="uploadStore.confirmPreview"
      @update:open="val => !val && uploadStore.cancelPreview()"
    />
  </div>
</template>
