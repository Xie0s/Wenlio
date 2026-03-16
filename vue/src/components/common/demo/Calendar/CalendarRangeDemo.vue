<!--
  CalendarRangeDemo.vue - Range Calendar 日期范围演示根组件
  职责：演示 RangeCalendar 原子组件的范围选择功能，对齐用户提供的简洁参考实现
  对外接口：无
-->
<script setup lang="ts">
import type { DateRange } from 'reka-ui'
import { getLocalTimeZone, today } from '@internationalized/date'
import { RangeCalendar } from '@/components/ui/range-calendar'
import { ref, type Ref } from 'vue'

const start = today(getLocalTimeZone())
const end = start.add({ days: 7 })

const dateRange = ref({
  start,
  end,
}) as Ref<DateRange>
</script>

<template>
  <RangeCalendar
    v-model="dateRange"
    class="rounded-2xl border shadow-sm"
    :number-of-months="2"
    disable-days-outside-current-view
  />
</template>


