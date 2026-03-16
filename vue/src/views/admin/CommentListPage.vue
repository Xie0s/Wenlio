<!-- CommentListPage.vue - 评论管理页面
     职责：纯 UI 渲染与用户交互，业务逻辑由 lib/comment.ts 提供 -->
<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { COMMENT_STATUS_LABEL, COMMENT_STATUS_COLOR } from '@/utils/types'
import type { Comment } from '@/utils/types'
import { useCommentList } from '@/lib/comment'
import { Button } from '@/components/ui/button'
import { Tabs, TabsList, TabsTrigger } from '@/components/ui/tabs'
import { Input } from '@/components/ui/input'
import { Badge } from '@/components/ui/badge'
import { Dialog, DialogContent, DialogDescription, DialogHeader, DialogTitle, DialogFooter } from '@/components/ui/dialog'
import { DropdownMenu, DropdownMenuContent, DropdownMenuItem, DropdownMenuSeparator, DropdownMenuTrigger } from '@/components/ui/dropdown-menu'
import { Tooltip, TooltipContent, TooltipTrigger } from '@/components/ui/tooltip'
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table'
import { Check, X, Trash2, Reply, ChevronLeft, ChevronRight, Eye, MoreHorizontal } from 'lucide-vue-next'
import ConfirmDialog from '@/components/common/ConfirmDialog.vue'

const {
  comments, page, statusFilter, loading, totalPages,
  showReply, replyTargetName, replyContent,
  loadComments, setFilter, approve, reject, deleteComment,
  openReply, submitReply, prevPage, nextPage, formatDate,
} = useCommentList()

const deleteConfirmOpen = ref(false)
const pendingDeleteId = ref<string | null>(null)

const showView = ref(false)
const viewTarget = ref<Comment | null>(null)
const openActionCommentId = ref('')

function openView(c: Comment) {
  viewTarget.value = c
  showView.value = true
}

function requestDelete(id: string) {
  pendingDeleteId.value = id
  deleteConfirmOpen.value = true
}

function confirmDelete() {
  if (pendingDeleteId.value) deleteComment(pendingDeleteId.value)
  pendingDeleteId.value = null
}

function handleActionMenuOpenChange(commentId: string, open: boolean) {
  openActionCommentId.value = open ? commentId : ''
}

onMounted(loadComments)
</script>

<template>
  <Tabs :model-value="statusFilter" class="space-y-6 text-base" @update:model-value="(v) => setFilter(v as string)">
    <div class="flex items-center justify-between">
      <TabsList>
        <TabsTrigger value="">全部</TabsTrigger>
        <TabsTrigger value="pending">待审核</TabsTrigger>
        <TabsTrigger value="approved">已批准</TabsTrigger>
        <TabsTrigger value="rejected">已拒绝</TabsTrigger>
      </TabsList>
    </div>

    <div class="max-h-[420px] overflow-auto rounded-xl border">
      <Table class="text-lg">
        <TableHeader>
          <TableRow>
            <TableHead class="h-11 w-28 text-lg font-light">作者</TableHead>
            <TableHead class="h-11 text-lg font-light">评论内容</TableHead>
            <TableHead class="h-11 w-20 text-lg font-light">状态</TableHead>
            <TableHead class="h-11 w-20 text-lg font-light">类型</TableHead>
            <TableHead class="h-11 w-[120px] text-lg font-light">时间</TableHead>
            <TableHead class="h-11 w-[80px] text-right text-lg font-light">操作</TableHead>
          </TableRow>
        </TableHeader>
        <TableBody>
          <TableRow v-for="c in comments" :key="c.id" class="h-12">
            <TableCell class="py-3">
              <span class="text-lg font-light">{{ c.author.name }}</span>
              <span v-if="c.author.email" class="block text-xs text-muted-foreground">{{ c.author.email }}</span>
            </TableCell>
            <TableCell class="py-3 max-w-0">
              <p class="truncate text-lg font-light">{{ c.content }}</p>
            </TableCell>
            <TableCell class="py-3">
              <Badge :class="COMMENT_STATUS_COLOR[c.status]" class="px-2.5 py-1 text-sm font-light">
                {{ COMMENT_STATUS_LABEL[c.status] }}
              </Badge>
            </TableCell>
            <TableCell class="py-3 text-lg font-light text-muted-foreground">
              {{ c.parent_id ? '回复' : '评论' }}
            </TableCell>
            <TableCell class="py-3 text-lg font-light text-muted-foreground">
              {{ formatDate(c.created_at) }}
            </TableCell>
            <TableCell class="py-3 text-right">
              <DropdownMenu :modal="false" @update:open="(open: boolean) => handleActionMenuOpenChange(c.id, open)">
                <DropdownMenuTrigger as-child>
                  <Button variant="ghost" size="icon" class="h-8 w-8 rounded-full hover:bg-foreground/10"
                    :class="{ 'bg-foreground/10': openActionCommentId === c.id }">
                    <MoreHorizontal class="h-4 w-4" />
                  </Button>
                </DropdownMenuTrigger>
                <DropdownMenuContent align="end" class="rounded-3xl">
                  <DropdownMenuItem v-if="c.status === 'pending'" @click="approve(c.id)">
                    <Check class="mr-2 h-4 w-4 text-green-600" />
                    批准
                  </DropdownMenuItem>
                  <DropdownMenuItem v-if="c.status === 'pending'" @click="reject(c.id)">
                    <X class="mr-2 h-4 w-4 text-red-600" />
                    拒绝
                  </DropdownMenuItem>
                  <DropdownMenuItem v-if="!c.parent_id" @click="openReply(c)">
                    <Reply class="mr-2 h-4 w-4" />
                    回复
                  </DropdownMenuItem>
                  <DropdownMenuItem @click="openView(c)">
                    <Eye class="mr-2 h-4 w-4" />
                    查看详情
                  </DropdownMenuItem>
                  <DropdownMenuSeparator />
                  <DropdownMenuItem variant="destructive" @click="requestDelete(c.id)">
                    <Trash2 class="mr-2 h-4 w-4" />
                    删除
                  </DropdownMenuItem>
                </DropdownMenuContent>
              </DropdownMenu>
            </TableCell>
          </TableRow>
        </TableBody>
      </Table>
      <p v-if="!loading && comments.length === 0" class="py-6 text-center text-base text-muted-foreground">暂无评论</p>
    </div>

    <!-- 分页 -->
    <div v-if="totalPages > 1" class="flex items-center justify-center gap-3">
      <Tooltip>
        <TooltipTrigger as-child>
          <Button variant="outline" size="icon" class="rounded-full" :disabled="page <= 1" @click="prevPage">
            <ChevronLeft class="h-4 w-4" />
          </Button>
        </TooltipTrigger>
        <TooltipContent>上一页</TooltipContent>
      </Tooltip>
      <span class="text-sm text-muted-foreground">{{ page }} / {{ totalPages }}</span>
      <Tooltip>
        <TooltipTrigger as-child>
          <Button variant="outline" size="icon" class="rounded-full" :disabled="page >= totalPages" @click="nextPage">
            <ChevronRight class="h-4 w-4" />
          </Button>
        </TooltipTrigger>
        <TooltipContent>下一页</TooltipContent>
      </Tooltip>
    </div>

    <ConfirmDialog
      v-model:open="deleteConfirmOpen"
      type="delete"
      description="此操作将永久删除该评论，无法撤销。"
      @confirm="confirmDelete"
    />

    <!-- 查看详情对话框 -->
    <Dialog v-model:open="showView">
      <DialogContent class="rounded-3xl max-w-md">
        <DialogHeader>
          <DialogTitle>评论详情</DialogTitle>
          <DialogDescription>该条评论的完整信息</DialogDescription>
        </DialogHeader>
        <template v-if="viewTarget">
          <dl class="space-y-3 text-sm">
            <div class="flex gap-2">
              <dt class="w-16 shrink-0 text-muted-foreground">作者</dt>
              <dd class="font-medium">{{ viewTarget.author.name }}</dd>
            </div>
            <div v-if="viewTarget.author.email" class="flex gap-2">
              <dt class="w-16 shrink-0 text-muted-foreground">邮箱</dt>
              <dd class="break-all">{{ viewTarget.author.email }}</dd>
            </div>
            <div class="flex gap-2">
              <dt class="w-16 shrink-0 text-muted-foreground">类型</dt>
              <dd>{{ viewTarget.parent_id ? '回复' : '评论' }}</dd>
            </div>
            <div class="flex gap-2">
              <dt class="w-16 shrink-0 text-muted-foreground">状态</dt>
              <dd>
                <Badge :class="COMMENT_STATUS_COLOR[viewTarget.status]" class="text-xs">
                  {{ COMMENT_STATUS_LABEL[viewTarget.status] }}
                </Badge>
              </dd>
            </div>
            <div class="flex gap-2">
              <dt class="w-16 shrink-0 text-muted-foreground">时间</dt>
              <dd class="text-muted-foreground">{{ formatDate(viewTarget.created_at) }}</dd>
            </div>
            <div class="flex gap-2">
              <dt class="w-16 shrink-0 text-muted-foreground">内容</dt>
              <dd class="whitespace-pre-wrap leading-relaxed">{{ viewTarget.content }}</dd>
            </div>
          </dl>
        </template>
      </DialogContent>
    </Dialog>

    <!-- 回复对话框 -->
    <Dialog v-model:open="showReply">
      <DialogContent class="rounded-3xl">
        <DialogHeader><DialogTitle>回复评论</DialogTitle><DialogDescription>输入回复内容以回应该评论。</DialogDescription></DialogHeader>
        <form class="flex flex-col gap-3" @submit.prevent="submitReply">
          <p class="text-sm text-muted-foreground">
            回复 <span class="font-medium text-foreground">{{ replyTargetName }}</span> 的评论
          </p>
          <Input v-model="replyContent" placeholder="输入回复内容..." />
          <DialogFooter>
            <Tooltip>
              <TooltipTrigger as-child>
                <Button type="submit" size="icon" class="rounded-full">
                  <Check class="size-4" />
                </Button>
              </TooltipTrigger>
              <TooltipContent>回复</TooltipContent>
            </Tooltip>
          </DialogFooter>
        </form>
      </DialogContent>
    </Dialog>
  </Tabs>
</template>
