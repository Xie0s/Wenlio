<!-- components/DocPage/DocComments.vue
     职责：文档评论区，包含评论列表展示和评论提交表单
     对外暴露事件：submit-comment({ name, email, content })
     对外暴露方法：resetForm()（提交成功后由父组件调用清空表单） -->

<script setup lang="ts">
import { ref } from 'vue'
import { MessageSquare, ArrowUpRight, Loader2 } from 'lucide-vue-next'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Textarea } from '@/components/ui/textarea'
import { Tooltip, TooltipContent, TooltipTrigger } from '@/components/ui/tooltip'
import type { Comment } from '@/utils/types'

defineProps<{
  comments: Comment[]
  submitting: boolean
}>()

const emit = defineEmits<{
  'submit-comment': [data: { name: string; email: string; content: string }]
}>()

const name = ref('')
const email = ref('')
const content = ref('')

function handleSubmit() {
  emit('submit-comment', {
    name: name.value,
    email: email.value,
    content: content.value,
  })
}

function resetForm() {
  name.value = ''
  email.value = ''
  content.value = ''
}

defineExpose({ resetForm })
</script>

<template>
  <section class="mt-16 border-t border-border pt-10">
    <!-- 标题 -->
    <h2 class="mb-7 flex items-center gap-2 text-lg font-semibold text-foreground">
      <MessageSquare class="h-5 w-5 text-muted-foreground" />
      评论
      <span v-if="comments.length > 0" class="text-sm font-normal text-muted-foreground">
        {{ comments.length }} 条
      </span>
    </h2>

    <!-- 评论列表 -->
    <div v-if="comments.length > 0" class="mb-10 divide-y divide-border/50">
      <div
        v-for="c in comments"
        :key="c.id"
        class="py-4"
        :class="{ 'pl-8': c.parent_id }"
      >
        <div class="flex items-baseline gap-2 mb-1">
          <span class="text-xs font-medium text-foreground">{{ c.author.name }}</span>
          <span class="text-[11px] text-muted-foreground">
            {{ new Date(c.created_at).toLocaleDateString('zh-CN', { year: 'numeric', month: 'long', day: 'numeric' }) }}
          </span>
        </div>
        <p class="text-base leading-relaxed text-foreground/80 whitespace-pre-wrap">{{ c.content }}</p>
      </div>
    </div>

    <!-- 空状态 -->
    <div v-else class="mb-10 flex flex-col items-center justify-center py-10 text-center">
      <MessageSquare class="mb-3 h-10 w-10 text-muted-foreground/25" />
      <p class="text-sm text-muted-foreground">暂无评论，来做第一个评论者吧</p>
    </div>

    <!-- 提交表单 -->
    <form class="space-y-3 max-w-lg mx-auto" @submit.prevent="handleSubmit">
      <div class="flex gap-3">
        <Input v-model="name" placeholder="昵称（必填）" class="flex-1" />
        <Input v-model="email" placeholder="邮箱（选填）" class="flex-1" type="email" />
      </div>
      <div class="relative">
        <Textarea v-model="content" placeholder="写下你的评论..." :rows="4" class="resize-none pb-12" />
        <Tooltip>
          <TooltipTrigger as-child>
            <Button
              type="submit"
              :disabled="submitting"
              class="absolute bottom-3 right-3 h-8 w-8 rounded-full bg-primary p-0 text-primary-foreground hover:bg-primary/90 disabled:opacity-50"
            >
              <Loader2 v-if="submitting" class="h-4 w-4 animate-spin" />
              <ArrowUpRight v-else class="h-4 w-4" />
            </Button>
          </TooltipTrigger>
          <TooltipContent>提交评论</TooltipContent>
        </Tooltip>
      </div>
    </form>
  </section>
</template>
