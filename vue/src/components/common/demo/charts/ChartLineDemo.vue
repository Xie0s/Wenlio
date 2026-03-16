<!--
  ChartLineDemo.vue - 交互式折线图演示组件
  职责：使用 Unovis VisLine 展示可切换 desktop/mobile 的折线图，含坐标轴、Crosshair Tooltip
  对外接口：作为独立 demo 组件直接使用，无 props
-->
<script setup lang="ts">
import { computed, ref } from 'vue'
import type { ChartConfig } from '@/components/ui/chart'
import { VisAxis, VisLine, VisXYContainer } from '@unovis/vue'
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from '@/components/ui/card'
import {
  ChartContainer,
  ChartCrosshair,
  ChartTooltip,
  ChartTooltipContent,
  componentToString,
} from '@/components/ui/chart'

const chartData = [
  { date: new Date('2024-06-24'), desktop: 132, mobile: 180 },
  { date: new Date('2024-06-25'), desktop: 141, mobile: 190 },
  { date: new Date('2024-06-26'), desktop: 434, mobile: 380 },
  { date: new Date('2024-06-27'), desktop: 448, mobile: 490 },
  { date: new Date('2024-06-28'), desktop: 149, mobile: 200 },
  { date: new Date('2024-06-29'), desktop: 103, mobile: 160 },
  { date: new Date('2024-06-30'), desktop: 446, mobile: 400 },
]

type Data = (typeof chartData)[number]

const chartConfig = {
  views: {
    label: '浏览量',
    color: undefined,
  },
  desktop: {
    label: '桌面端',
    color: 'var(--chart-1)',
  },
  mobile: {
    label: '移动端',
    color: 'var(--chart-2)',
  },
} satisfies ChartConfig

const activeChart = ref<'desktop' | 'mobile'>('desktop')
const total = computed(() => ({
  desktop: chartData.reduce((acc, curr) => acc + curr.desktop, 0),
  mobile: chartData.reduce((acc, curr) => acc + curr.mobile, 0),
}))
</script>

<template>
  <Card class="py-4 sm:py-0">
    <CardHeader class="flex flex-col items-stretch border-b !p-0 sm:flex-row">
      <div class="flex flex-1 flex-col justify-center gap-1 px-6 pb-3 sm:pb-0">
        <CardTitle>折线图 - 交互式</CardTitle>
        <CardDescription>最近一周桌面端与移动端趋势</CardDescription>
      </div>
      <div class="flex">
        <button
          v-for="chart in (['desktop', 'mobile'] as const)"
          :key="chart"
          :data-active="activeChart === chart"
          class="data-[active=true]:bg-muted/50 flex flex-1 flex-col justify-center gap-1 border-t px-6 py-4 text-left even:border-l sm:border-t-0 sm:border-l sm:px-8 sm:py-6"
          @click="activeChart = chart"
        >
          <span class="text-muted-foreground text-xs">{{ chartConfig[chart].label }}</span>
          <span class="text-lg leading-none font-bold sm:text-3xl">{{ total[chart].toLocaleString() }}</span>
        </button>
      </div>
    </CardHeader>
    <CardContent class="px-2 sm:p-6">
      <ChartContainer :config="chartConfig" class="aspect-auto h-[250px] w-full" cursor>
        <VisXYContainer :data="chartData" :margin="{ left: -24 }" :y-domain="[0, undefined]">
          <VisLine
            :x="(d: Data) => d.date"
            :y="(d: Data) => d[activeChart]"
            :color="chartConfig[activeChart].color"
          />
          <VisAxis
            type="x"
            :x="(d: Data) => d.date"
            :tick-line="false"
            :domain-line="false"
            :grid-line="false"
            :tick-format="(d: number) => {
              const date = new Date(d)
              return date.toLocaleDateString('zh-CN', { month: 'short', day: 'numeric' })
            }"
          />
          <VisAxis type="y" :num-ticks="3" :tick-line="false" :domain-line="false" />
          <ChartTooltip />
          <ChartCrosshair
            :template="componentToString(chartConfig, ChartTooltipContent, {
              labelFormatter(d: number | Date) {
                return new Date(d).toLocaleDateString('zh-CN', {
                  month: 'short',
                  day: 'numeric',
                  year: 'numeric',
                })
              },
            })"
            :color="chartConfig.desktop.color"
          />
        </VisXYContainer>
      </ChartContainer>
    </CardContent>
  </Card>
</template>
