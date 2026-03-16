/**
 * comment.ts - 评论管理业务逻辑层
 *
 * 职责：封装评论列表的状态管理与 API 操作（加载、筛选、分页、审核、回复、删除）
 * 对外暴露：useCommentList() composable 函数
 */
import { ref, computed } from 'vue'
import { http } from '@/utils/http'
import type { PageData } from '@/utils/http'
import type { Comment } from '@/utils/types'
import { toast } from 'vue-sonner'

const PAGE_SIZE = 20

export function useCommentList() {
  const comments = ref<Comment[]>([])
  const total = ref(0)
  const page = ref(1)
  const statusFilter = ref('')
  const loading = ref(false)

  // 回复弹窗状态
  const showReply = ref(false)
  const replyTargetId = ref('')
  const replyTargetName = ref('')
  const replyContent = ref('')

  const totalPages = computed(() => Math.ceil(total.value / PAGE_SIZE) || 1)

  async function loadComments() {
    loading.value = true
    const res = await http.get<PageData<Comment>>('/tenant/comments', {
      page: page.value,
      page_size: PAGE_SIZE,
      status: statusFilter.value,
    })
    loading.value = false
    if (res.code === 0 && res.data) {
      comments.value = res.data.list || []
      total.value = res.data.pagination.total
    }
  }

  function setFilter(status: string) {
    statusFilter.value = status
    page.value = 1
    loadComments()
  }

  async function approve(id: string) {
    const res = await http.post(`/tenant/comments/${id}/approve`)
    res.code === 0 ? (toast.success('已批准'), loadComments()) : toast.error(res.message)
  }

  async function reject(id: string) {
    const res = await http.post(`/tenant/comments/${id}/reject`)
    res.code === 0 ? (toast.success('已拒绝'), loadComments()) : toast.error(res.message)
  }

  async function deleteComment(id: string) {
    const res = await http.delete(`/tenant/comments/${id}`)
    res.code === 0 ? (toast.success('已删除'), loadComments()) : toast.error(res.message)
  }

  function openReply(c: Comment) {
    replyTargetId.value = c.id
    replyTargetName.value = c.author.name
    replyContent.value = ''
    showReply.value = true
  }

  async function submitReply() {
    if (!replyContent.value.trim()) {
      toast.error('请输入回复内容')
      return
    }
    const res = await http.post(`/tenant/comments/${replyTargetId.value}/reply`, {
      content: replyContent.value,
    })
    if (res.code === 0) {
      toast.success('回复成功')
      showReply.value = false
      loadComments()
    } else {
      toast.error(res.message)
    }
  }

  function prevPage() {
    if (page.value > 1) { page.value--; loadComments() }
  }

  function nextPage() {
    if (page.value < totalPages.value) { page.value++; loadComments() }
  }

  function formatDate(dateStr: string) {
    if (!dateStr) return '-'
    return new Date(dateStr).toLocaleDateString('zh-CN', { year: 'numeric', month: '2-digit', day: '2-digit' })
  }

  return {
    // 状态
    comments,
    total,
    page,
    statusFilter,
    loading,
    totalPages,
    // 回复弹窗
    showReply,
    replyTargetName,
    replyContent,
    // 操作
    loadComments,
    setFilter,
    approve,
    reject,
    deleteComment,
    openReply,
    submitReply,
    prevPage,
    nextPage,
    formatDate,
  }
}
