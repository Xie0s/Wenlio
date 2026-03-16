<script setup lang="ts">
/**
 * TextHoverEffectDemo.vue - TextHoverEffect 演示组件
 * 职责：展示 TextHoverEffect 的各类用法变体
 * 对外暴露：无（演示专用）
 */
import { TextHoverEffect } from '@/components/ui/a-aceternity'
</script>

<template>
  <div class="space-y-12">
    <!-- 即时跟随（默认）-->
    <section class="space-y-4">
      <div class="rounded-2xl border bg-card/40 p-4">
        <h3 class="text-sm font-semibold">即时跟随（duration: 0）</h3>
        <p class="text-sm text-muted-foreground mt-1">光效随鼠标即时移动，无平滑过渡。</p>
      </div>
      <div class="h-40 w-full">
        <TextHoverEffect text="Aceternity" />
      </div>
    </section>

    <div class="border-t border-dashed" />

    <!-- 平滑跟随 -->
    <section class="space-y-4">
      <div class="rounded-2xl border bg-card/40 p-4">
        <h3 class="text-sm font-semibold">平滑跟随（duration: 0.5）</h3>
        <p class="text-sm text-muted-foreground mt-1">光效以指数平滑方式跟随鼠标，过渡更自然。</p>
      </div>
      <div class="h-40 w-full">
        <TextHoverEffect text="UI Magic" :duration="0.5" />
      </div>
    </section>

    <div class="border-t border-dashed" />

    <!-- 长文字 -->
    <section class="space-y-4">
      <div class="rounded-2xl border bg-card/40 p-4">
        <h3 class="text-sm font-semibold">较长文字</h3>
        <p class="text-sm text-muted-foreground mt-1">SVG viewBox 固定 300×100，文字会自适应缩放。</p>
      </div>
      <div class="h-40 w-full">
        <TextHoverEffect text="Hello World" :duration="0.3" />
      </div>
    </section>
  </div>
</template>
