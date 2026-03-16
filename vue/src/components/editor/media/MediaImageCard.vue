<!--
  MediaImageCard.vue - 媒体图片卡组件
  职责：以 HeroUI Card 呈现单个图片媒体项，底部毛玻璃信息栏，悬浮显示操作按钮和标签
  对外暴露：Props: item, deleting, usageLoaded, usageRefs
            Events: delete, open-usage
-->
<script setup lang="ts">
import { computed } from 'vue'
import type { MediaItem, MediaUsageRef } from '@/lib/media'
import { formatBytes } from '@/lib/media'
import { Card, CardFooter } from '@/components/ui/a-heroui/card'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import { Tooltip, TooltipContent, TooltipTrigger } from '@/components/ui/tooltip'
import {
  Trash2, ExternalLink, Cloud, HardDrive,
  Unlink, BookOpen,
} from 'lucide-vue-next'

const props = defineProps<{
  item: MediaItem
  deleting: string | null
  usageLoaded: boolean
  usageRefs?: MediaUsageRef[]
}>()

const emit = defineEmits<{
  delete: [id: string]
  'open-usage': [item: MediaItem]
}>()

const hasUsage = computed(() => (props.usageRefs?.length ?? 0) > 0)

function openInNewTab() {
  window.open(props.item.file_url, '_blank', 'noopener')
}
</script>

<template>
  <Card class="group border-none" radius="2xl" style="content-visibility: auto; contain-intrinsic-size: auto 200px">
    <!-- 图片 -->
    <img
      :src="item.file_url"
      :alt="item.file_name"
      class="w-full aspect-square object-cover"
      loading="lazy"
    />

    <!-- 右上角操作按钮（悬浮显示） -->
    <div class="absolute top-2 right-2 flex gap-1 opacity-0 group-hover:opacity-100 transition-opacity">
      <Tooltip>
        <TooltipTrigger as-child>
          <Button
            variant="outline" size="icon"
            class="size-7 rounded-full shadow bg-background/80 backdrop-blur-sm"
            @click="openInNewTab"
          >
            <ExternalLink class="size-3.5" :stroke-width="1.5" />
          </Button>
        </TooltipTrigger>
        <TooltipContent>新标签页打开</TooltipContent>
      </Tooltip>
      <Tooltip>
        <TooltipTrigger as-child>
          <Button
            variant="outline" size="icon"
            class="size-7 rounded-full shadow bg-background/80 backdrop-blur-sm text-destructive hover:text-destructive hover:bg-destructive/10 border-destructive/30"
            :disabled="deleting === item.id"
            @click="emit('delete', item.id)"
          >
            <Trash2 class="size-3.5" :stroke-width="1.5" />
          </Button>
        </TooltipTrigger>
        <TooltipContent>删除文件</TooltipContent>
      </Tooltip>
    </div>

    <!-- 左上角标签（悬浮显示） -->
    <div class="absolute top-2 left-2 flex gap-1 opacity-0 group-hover:opacity-100 transition-opacity">
      <Badge variant="secondary" class="text-[10px] px-1.5 py-0 h-5 gap-1 bg-background/70 backdrop-blur-sm">
        <Cloud v-if="item.storage_type === 'cloud'" class="size-2.5" :stroke-width="1.5" />
        <HardDrive v-else class="size-2.5" :stroke-width="1.5" />
        {{ item.storage_type === 'cloud' ? '云' : '本地' }}
      </Badge>
      <Badge
        v-if="usageLoaded && !hasUsage"
        variant="secondary"
        class="text-[10px] px-1.5 py-0 h-5 gap-1 bg-amber-500/80 text-white backdrop-blur-sm border-0"
      >
        <Unlink class="size-2.5" :stroke-width="1.5" />
        未使用
      </Badge>
    </div>

    <!-- 底部毛玻璃信息栏 -->
    <CardFooter
      is-blurred
      class="justify-between py-1.5 rounded-2xl bottom-1 w-[calc(100%-8px)] shadow-small ml-1 border border-white/20"
    >
      <div class="min-w-0 flex-1">
        <p class="text-xs text-white/90 truncate">{{ item.file_name }}</p>
        <p class="text-[10px] text-white/60">{{ formatBytes(item.file_size) }}</p>
      </div>
      <button
        v-if="usageLoaded && hasUsage"
        class="flex items-center gap-1 text-[10px] text-white/80 hover:text-white shrink-0 ml-2"
        @click="emit('open-usage', item)"
      >
        <BookOpen class="size-3" :stroke-width="1.5" />
        {{ usageRefs?.length }}处
      </button>
    </CardFooter>
  </Card>
</template>
