<!--
  CalendarDatePickerDemo.vue - Calendar 日期选择器演示根组件
  职责：封装“Calendar + Popover”单日期选择器演示，展示月/年快速跳转与选择后自动关闭
  对外接口：无（仅作为 demo 展示组件）
-->
<script setup lang="ts">
import type { DateValue } from '@internationalized/date'
import { CalendarDate, getLocalTimeZone, today } from '@internationalized/date'
import { CalendarIcon, ChevronDownIcon } from 'lucide-vue-next'
import type { Ref } from 'vue'
import { ref } from 'vue'

import { Button } from '@/components/ui/button'
import { Calendar } from '@/components/ui/calendar'
import { Label } from '@/components/ui/label'
import {
  Popover,
  PopoverContent,
  PopoverTrigger,
} from '@/components/ui/popover'
import { cn } from '@/utils'

const localTimeZone = getLocalTimeZone()

// 场景 1：Popover Calendar（生日选择）
const birthDate = ref(today(localTimeZone)) as Ref<DateValue>

// 场景 2：Popover Calendar（预约日期，初始为空）
const bookingPlaceholder = today(localTimeZone)
const bookingDate = ref<DateValue>()

// 场景 3：Popover Calendar（排期日期，带可选范围）
const scheduleDate = ref(today(localTimeZone)) as Ref<DateValue>

function formatDate(date?: DateValue) {
  if (!date) {
    return '请选择日期'
  }
  return date.toDate(localTimeZone).toLocaleDateString('zh-CN')
}
</script>

<template>
  <div class="space-y-4">
    <div class="rounded-2xl border bg-card/30 p-4 space-y-3">
      <h4 class="text-sm font-semibold">场景 1：Popover Calendar（生日选择）</h4>
      <div class="flex flex-col gap-3">
        <Label for="calendar-date-of-birth" class="px-1 text-xs text-muted-foreground">
          出生日期
        </Label>

        <Popover v-slot="{ close }">
          <PopoverTrigger as-child>
            <Button
              id="calendar-date-of-birth"
              variant="outline"
              :class="cn('w-64 justify-between font-normal', !birthDate && 'text-muted-foreground')"
            >
              <span class="inline-flex items-center gap-2">
                <CalendarIcon class="h-4 w-4" />
                {{ formatDate(birthDate) }}
              </span>
              <ChevronDownIcon class="h-4 w-4" />
            </Button>
          </PopoverTrigger>

          <PopoverContent class="w-auto overflow-hidden p-0" align="start">
            <Calendar
              :model-value="birthDate"
              locale="zh-CN"
              layout="month-and-year"
              @update:model-value="(value) => {
                if (value) {
                  birthDate = value
                  close()
                }
              }"
            />
          </PopoverContent>
        </Popover>
      </div>
    </div>

    <div class="rounded-2xl border bg-card/30 p-4 space-y-3">
      <h4 class="text-sm font-semibold">场景 2：Popover Calendar（预约日期，初始空值）</h4>
      <div class="flex flex-col gap-3">
        <Label for="calendar-booking-date" class="px-1 text-xs text-muted-foreground">
          预约日期
        </Label>

        <Popover v-slot="{ close }">
          <PopoverTrigger as-child>
            <Button
              id="calendar-booking-date"
              variant="outline"
              :class="cn('w-64 justify-between font-normal', !bookingDate && 'text-muted-foreground')"
            >
              <span class="inline-flex items-center gap-2">
                <CalendarIcon class="h-4 w-4" />
                {{ formatDate(bookingDate) }}
              </span>
              <ChevronDownIcon class="h-4 w-4" />
            </Button>
          </PopoverTrigger>

          <PopoverContent class="w-auto overflow-hidden p-0" align="start">
            <Calendar
              :model-value="bookingDate"
              locale="zh-CN"
              :default-placeholder="bookingPlaceholder"
              layout="month-and-year"
              @update:model-value="(value) => {
                if (value) {
                  bookingDate = value
                  close()
                }
              }"
            />
          </PopoverContent>
        </Popover>
      </div>
    </div>

    <div class="rounded-2xl border bg-card/30 p-4 space-y-3">
      <h4 class="text-sm font-semibold">场景 3：Popover Calendar（排期日期 + 可选范围）</h4>
      <div class="flex flex-col gap-3">
        <Label for="calendar-schedule-date" class="px-1 text-xs text-muted-foreground">
          排期日期
        </Label>

        <Popover v-slot="{ close }">
          <PopoverTrigger as-child>
            <Button
              id="calendar-schedule-date"
              variant="outline"
              :class="cn('w-64 justify-between font-normal', !scheduleDate && 'text-muted-foreground')"
            >
              <span class="inline-flex items-center gap-2">
                <CalendarIcon class="h-4 w-4" />
                {{ formatDate(scheduleDate) }}
              </span>
              <ChevronDownIcon class="h-4 w-4" />
            </Button>
          </PopoverTrigger>

          <PopoverContent class="w-auto overflow-hidden p-0" align="start">
            <Calendar
              :model-value="scheduleDate"
              locale="zh-CN"
              layout="month-and-year"
              :min-value="new CalendarDate(2020, 1, 1)"
              :max-value="new CalendarDate(2035, 12, 31)"
              @update:model-value="(value) => {
                if (value) {
                  scheduleDate = value
                  close()
                }
              }"
            />
          </PopoverContent>
        </Popover>
      </div>
    </div>
  </div>
</template>
