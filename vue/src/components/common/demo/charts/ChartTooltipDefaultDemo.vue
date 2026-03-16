<!--
  ChartTooltipDefaultDemo.vue - 默认 Tooltip 演示组件
  职责：使用 Unovis VisStackedBar 演示 ChartTooltip + ChartCrosshair + ChartTooltipContent 的默认用法
  对外接口：作为独立 demo 组件直接使用，无 props
-->
<script setup lang="ts">
import type { ChartConfig } from '@/components/ui/chart'
import { VisAxis, VisStackedBar, VisXYContainer } from '@unovis/vue'
import { TrendingUp } from 'lucide-vue-next'
import {
  Card,
  CardContent,
  CardDescription,
  CardFooter,
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
  { date: new Date('2024-07-15'), running: 450, swimming: 300 },
  { date: new Date('2024-07-16'), running: 380, swimming: 420 },
  { date: new Date('2024-07-17'), running: 520, swimming: 120 },
  { date: new Date('2024-07-18'), running: 140, swimming: 550 },
  { date: new Date('2024-07-19'), running: 600, swimming: 350 },
  { date: new Date('2024-07-20'), running: 480, swimming: 400 },
]

type Data = (typeof chartData)[number]

const chartConfig = {
  running: {
    label: '跑步',
    color: 'var(--chart-1)',
  },
  swimming: {
    label: '游泳',
    color: 'var(--chart-2)',
  },
} satisfies ChartConfig
</script>

<template>
  <Card>
    <CardHeader>
      <CardTitle>Tooltip - 默认样式</CardTitle>
      <CardDescription>基于 ChartTooltipContent 的默认 Tooltip 演示。</CardDescription>
    </CardHeader>
    <CardContent>
      <ChartContainer :config="chartConfig">
        <VisXYContainer :data="chartData" :padding="{ top: 10, bottom: 10, left: 10, right: 10 }">
          <VisStackedBar
            :x="(d: Data) => d.date"
            :y="[(d: Data) => d.running, (d: Data) => d.swimming]"
            :color="[chartConfig.running.color, chartConfig.swimming.color]"
            :rounded-corners="4"
            :bar-padding="0.1"
          />
          <VisAxis
            type="x"
            :x="(d: Data) => d.date"
            :tick-line="false"
            :domain-line="false"
            :grid-line="false"
            :num-ticks="6"
            :tick-format="(d: number) => {
              const date = new Date(d)
              return date.toLocaleDateString('zh-CN', { weekday: 'short' })
            }"
            :tick-values="chartData.map(d => d.date)"
          />
          <ChartTooltip />
          <ChartCrosshair
            :template="componentToString(chartConfig, ChartTooltipContent, {
              labelFormatter(d: number | Date) {
                return new Date(d).toLocaleDateString('zh-CN')
              },
            })"
            color="#0000"
          />
        </VisXYContainer>
      </ChartContainer>
    </CardContent>
    <CardFooter class="flex-col items-start gap-2 text-sm">
      <div class="flex gap-2 font-medium leading-none">
        本周训练强度提升 5.2% <TrendingUp class="h-4 w-4" />
      </div>
      <div class="leading-none text-muted-foreground">
        展示最近 6 天运动数据
      </div>
    </CardFooter>
  </Card>
</template>
