<!-- PHero.vue - Hero 区块入口
     职责：根据 HeroConfig.layout 分发到对应的布局子组件
     对外接口：
       Props: config
       Emits: navigate(url: string) -->
<script setup lang="ts">
import type { HeroConfig } from '../types'
import HeroCentered from './HeroCentered.vue'
import HeroSplit from './HeroSplit.vue'

defineProps<{
  config: HeroConfig
}>()

defineEmits<{
  navigate: [url: string]
}>()
</script>

<template>
  <HeroCentered v-if="config.layout === 'centered'" :config="config" @navigate="$emit('navigate', $event)" />
  <HeroSplit v-else :config="config" @navigate="$emit('navigate', $event)" />
</template>
