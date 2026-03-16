<!--
  ThemeToggle.vue - 主题切换控件（公开入口）
  职责：统一对外暴露主题切换入口，内部委托 AnimatedThemeToggler 实现带视觉转场动效的明暗切换
  对外暴露：透传 AnimatedThemeToggler 的所有 Props

  ┌─────────────────────────────────────────────────────────────────────┐
  │ Props（均可选，通过 v-bind="$attrs" 透传）                           │
  ├──────────────┬──────────┬──────────────────────────────────────────┤
  │ prop         │ 默认值   │ 说明                                     │
  ├──────────────┼──────────┼──────────────────────────────────────────┤
  │ button-size  │ size-9   │ 按钮尺寸（Tailwind size-* 类名）          │
  │ icon-size    │ size-5   │ 图标尺寸（Tailwind size-* / h-* w-* 类名）│
  │ duration     │ 500      │ View Transition 动画时长（ms）            │
  │ class        │ —        │ 附加样式（仅用于装饰：border、bg、hover 等）│
  └──────────────┴──────────┴──────────────────────────────────────────┘

  使用示例：

  1. 默认（36×36，图标 20px）：
     <ThemeToggle />

  2. 大按钮（44×44，图标 24px）—— 侧边栏工具栏：
     <ThemeToggle button-size="size-11" icon-size="size-6" />

  3. 超大按钮（48×48）—— AdminDock：
     <ThemeToggle button-size="size-12" class="text-foreground" icon-size="size-[22px]" />

  4. 紧凑按钮（32×32）—— 登录页：
     <ThemeToggle button-size="size-8" class="border border-border/60 bg-background/55" icon-size="h-4 w-4" />

  5. 带边框装饰（36×36，图标 18px）—— EditorToolbar：
     <ThemeToggle button-size="size-9" class="border border-input bg-background" icon-size="size-[18px]" />

  6. 导航栏（40×40）：
     <ThemeToggle button-size="size-10" class="border border-foreground/20" />

  注意：尺寸由 button-size / icon-size 控制，class 中不要再写 h-* w-* size-* 等尺寸类。
-->
<script setup lang="ts">
import AnimatedThemeToggler from '@/components/ui/a-magicui/AnimatedThemeToggler.vue'

defineOptions({ inheritAttrs: false })
</script>

<template>
  <AnimatedThemeToggler v-bind="$attrs" />
</template>
