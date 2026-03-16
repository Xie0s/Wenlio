<!--
  ChartDonutDemo.vue - 饼图演示组件
  职责：使用 Unovis VisDonut 展示饼图，含中心标签（总数）、segment tooltip、Card 布局
  对外接口：作为独立 demo 组件直接使用，无 props
-->
<script setup lang="ts">
import { computed } from 'vue'
import type { ChartConfig } from '@/components/ui/chart'
import { Donut } from '@unovis/ts'
import { VisDonut, VisSingleContainer } from '@unovis/vue'
import { TrendingUp } from 'lucide-vue-next'
import {
  ChartContainer,
  ChartTooltip,
  ChartTooltipContent,
  componentToString,
} from '@/components/ui/chart'
import {
  Card,
  CardContent,
  CardDescription,
  CardFooter,
  CardHeader,
  CardTitle,
} from '@/components/ui/card'

const chartData = [
  { browser: 'chrome', visitors: 275, fill: 'var(--color-chrome)' },
  { browser: 'safari', visitors: 200, fill: 'var(--color-safari)' },
  { browser: 'firefox', visitors: 287, fill: 'var(--color-firefox)' },
  { browser: 'edge', visitors: 173, fill: 'var(--color-edge)' },
  { browser: 'other', visitors: 190, fill: 'var(--color-other)' },
]
type Data = (typeof chartData)[number]

const chartConfig = {
  visitors: {
    label: '访客数',
    color: undefined,
  },
  chrome: {
    label: 'Chrome',
    color: 'var(--chart-1)',
  },
  safari: {
    label: 'Safari',
    color: 'var(--chart-2)',
  },
  firefox: {
    label: 'Firefox',
    color: 'var(--chart-3)',
  },
  edge: {
    label: 'Edge',
    color: 'var(--chart-4)',
  },
  other: {
    label: '其他',
    color: 'var(--chart-5)',
  },
} satisfies ChartConfig

const totalVisitors = computed(() => chartData.reduce((acc, curr) => acc + curr.visitors, 0))
</script>

<template>
  <Card class="flex flex-col">
    <CardHeader class="items-center pb-0">
      <CardTitle>饼图</CardTitle>
      <CardDescription>2024 年 1 月 - 6 月</CardDescription>
    </CardHeader>
    <CardContent class="flex-1 pb-0">
      <ChartContainer
        :config="chartConfig"
        class="mx-auto aspect-square max-h-[250px]"
        :style="{
          '--vis-donut-central-label-font-size': 'var(--text-3xl)',
          '--vis-donut-central-label-font-weight': 'var(--font-weight-bold)',
          '--vis-donut-central-label-text-color': 'var(--foreground)',
          '--vis-donut-central-sub-label-text-color': 'var(--muted-foreground)',
        }"
      >
        <VisSingleContainer
          :data="chartData"
          :margin="{ top: 30, bottom: 30 }"
        >
          <VisDonut
            :value="(d: Data) => d.visitors"
            :color="(d: Data) => chartConfig[d.browser as keyof typeof chartConfig].color"
            :arc-width="30"
            :central-label-offset-y="10"
            :central-label="totalVisitors.toLocaleString()"
            central-sub-label="访客"
          />
          <ChartTooltip
            :triggers="{
              [Donut.selectors.segment]: componentToString(chartConfig, ChartTooltipContent, { hideLabel: true })!,
            }"
          />
        </VisSingleContainer>
      </ChartContainer>
    </CardContent>
    <CardFooter class="flex-col gap-2 text-sm">
      <div class="flex items-center gap-2 font-medium leading-none">
        本月增长 5.2% <TrendingUp class="h-4 w-4" />
      </div>
      <div class="leading-none text-muted-foreground">
        展示近 6 个月的访客总量
      </div>
    </CardFooter>
  </Card>
</template>
