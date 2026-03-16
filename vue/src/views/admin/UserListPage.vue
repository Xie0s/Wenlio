<!-- UserListPage.vue - 全平台用户管理（超管专属）
     职责：纯 UI 渲染与用户交互，业务逻辑由 lib/user-list.ts 提供 -->
<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { useUserList, ROLE_LABEL, formatDate, formatDateTime } from '@/lib/user-list'
import { useAuthStore } from '@/stores/auth'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table'
import { Badge } from '@/components/ui/badge'
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
  DialogFooter,
} from '@/components/ui/dialog'
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu'
import { Tooltip, TooltipContent, TooltipTrigger } from '@/components/ui/tooltip'
import { Plus, KeyRound, Check, Pencil, Ban, CheckCircle, Trash2, Search, X, MoreHorizontal } from 'lucide-vue-next'
import UserHoverCard from '@/components/common/UserHoverCard.vue'
import ConfirmDialog from '@/components/common/ConfirmDialog.vue'

const {
  users, loading, keyword,
  showCreate, createForm,
  showEdit, editForm,
  showResetPwd, resetUserName, newPassword,
  loadUsers, createSuperAdmin, openEditUser, updateUser, toggleUserStatus, deleteUser, openResetPassword, resetPassword,
} = useUserList()
const authStore = useAuthStore()
const currentUserId = authStore.user?.id || ''
const appliedKeyword = ref('')
const isSearchApplied = computed(() => !!appliedKeyword.value && appliedKeyword.value === keyword.value)
const openActionUserId = ref('')

const banConfirmOpen = ref(false)
const pendingBanUser = ref<(typeof users.value)[number] | null>(null)
const deleteConfirmOpen = ref(false)
const pendingDeleteUser = ref<(typeof users.value)[number] | null>(null)

function requestToggleStatus(u: (typeof users.value)[number]) {
  if (u.status === 'active') {
    pendingBanUser.value = u
    banConfirmOpen.value = true
  } else {
    toggleUserStatus(u)
  }
}

function confirmBan() {
  if (pendingBanUser.value) toggleUserStatus(pendingBanUser.value)
  pendingBanUser.value = null
}

function requestDeleteUser(u: (typeof users.value)[number]) {
  pendingDeleteUser.value = u
  deleteConfirmOpen.value = true
}

function confirmDeleteUser() {
  if (pendingDeleteUser.value) deleteUser(pendingDeleteUser.value)
  pendingDeleteUser.value = null
}

function handleActionMenuOpenChange(userId: string, open: boolean) {
  openActionUserId.value = open ? userId : ''
}

function executeSearch() {
  keyword.value = keyword.value.trim()
  appliedKeyword.value = keyword.value
  loadUsers()
}

function handleSearchAction() {
  if (isSearchApplied.value) {
    keyword.value = ''
    appliedKeyword.value = ''
    loadUsers()
    return
  }
  executeSearch()
}

onMounted(loadUsers)
</script>

<template>
  <div class="space-y-6 text-base">
    <div class="flex gap-2">
      <Input v-model="keyword" placeholder="搜索用户名、姓名、邮箱或租户 ID..." class="max-w-xs" @keyup.enter="executeSearch" />
      <Tooltip>
        <TooltipTrigger as-child>
          <Button variant="outline" size="icon" class="rounded-full" @click="handleSearchAction">
            <X v-if="isSearchApplied" class="h-4 w-4" />
            <Search v-else class="h-4 w-4" />
          </Button>
        </TooltipTrigger>
        <TooltipContent>{{ isSearchApplied ? '清除搜索' : '搜索' }}</TooltipContent>
      </Tooltip>
      <Tooltip>
        <TooltipTrigger as-child>
          <Button size="icon" class="rounded-full" @click="showCreate = true">
            <Plus class="h-5 w-5" />
          </Button>
        </TooltipTrigger>
        <TooltipContent>创建超级管理员</TooltipContent>
      </Tooltip>
    </div>

    <div class="max-h-[420px] overflow-auto rounded-xl border">
      <Table class="text-lg">
        <TableHeader>
          <TableRow>
            <TableHead class="h-11 w-[130px] text-lg font-light">姓名</TableHead>
            <TableHead class="h-11 text-lg font-light">用户名</TableHead>
            <TableHead class="h-11 text-lg font-light">邮箱</TableHead>
            <TableHead class="h-11 w-[120px] text-lg font-light">所属租户</TableHead>
            <TableHead class="h-11 w-[110px] text-lg font-light">角色</TableHead>
            <TableHead class="h-11 w-[80px] text-lg font-light">状态</TableHead>
            <TableHead class="h-11 w-[110px] text-lg font-light">创建时间</TableHead>
            <TableHead class="h-11 w-[110px] text-lg font-light">最近登录</TableHead>
            <TableHead class="h-11 w-[80px] text-right text-lg font-light">操作</TableHead>
          </TableRow>
        </TableHeader>
        <TableBody>
          <TableRow v-for="u in users" :key="u.id" class="h-12">
            <TableCell class="py-3 text-lg font-light">
              <UserHoverCard :avatar-url="u.avatar_url" :name="u.name" :username="u.username" :bio="u.bio"
                :profile-bg-url="u.profile_bg_url" :role-label="ROLE_LABEL[u.role] || u.role" />
            </TableCell>
            <TableCell class="py-3 text-lg font-light text-muted-foreground">{{ u.username }}</TableCell>
            <TableCell class="py-3 text-lg font-light text-muted-foreground">{{ u.email || '-' }}</TableCell>
            <TableCell class="py-3 text-lg font-light text-muted-foreground">{{ u.tenant_id || '-' }}</TableCell>
            <TableCell class="py-3">
              <Badge variant="outline" class="px-2.5 py-1 text-sm font-light">{{ ROLE_LABEL[u.role] || u.role }}</Badge>
            </TableCell>
            <TableCell class="py-3">
              <Badge :variant="u.status === 'active' ? 'default' : 'destructive'"
                class="px-2.5 py-1 text-sm font-light">
                {{ u.status === 'active' ? '活跃' : '已禁用' }}
              </Badge>
            </TableCell>
            <TableCell class="py-3 text-lg font-light text-muted-foreground">{{ formatDate(u.created_at) }}</TableCell>
            <TableCell class="py-3 text-lg font-light text-muted-foreground">{{ formatDateTime(u.last_login_at) }}
            </TableCell>
            <TableCell class="py-3 text-right">
              <DropdownMenu :modal="false" @update:open="(open: boolean) => handleActionMenuOpenChange(u.id, open)">
                <DropdownMenuTrigger as-child>
                  <Button variant="ghost" size="icon" class="h-8 w-8 rounded-full hover:bg-foreground/10"
                    :class="{ 'bg-foreground/10': openActionUserId === u.id }">
                    <MoreHorizontal class="h-4 w-4" />
                  </Button>
                </DropdownMenuTrigger>
                <DropdownMenuContent align="end" class="rounded-3xl">
                  <DropdownMenuItem @click="openResetPassword(u)">
                    <KeyRound class="mr-2 h-4 w-4" />
                    重置密码
                  </DropdownMenuItem>
                  <DropdownMenuItem @click="openEditUser(u)">
                    <Pencil class="mr-2 h-4 w-4" />
                    编辑
                  </DropdownMenuItem>
                  <template v-if="u.id !== currentUserId">
                    <DropdownMenuSeparator />
                    <DropdownMenuItem @click="requestToggleStatus(u)">
                      <Ban v-if="u.status === 'active'" class="mr-2 h-4 w-4" />
                      <CheckCircle v-else class="mr-2 h-4 w-4" />
                      {{ u.status === 'active' ? '禁用' : '启用' }}
                    </DropdownMenuItem>
                    <DropdownMenuItem variant="destructive" @click="requestDeleteUser(u)">
                      <Trash2 class="mr-2 h-4 w-4" />
                      删除
                    </DropdownMenuItem>
                  </template>
                </DropdownMenuContent>
              </DropdownMenu>
            </TableCell>
          </TableRow>
        </TableBody>
      </Table>
      <p v-if="!loading && users.length === 0" class="py-6 text-center text-base text-muted-foreground">暂无用户，点击右上角创建</p>
    </div>

    <ConfirmDialog v-model:open="banConfirmOpen" type="ban"
      :description="`禁用后用户「${pendingBanUser?.name || pendingBanUser?.username}」将无法登录，可随时在此页面启用。`"
      @confirm="confirmBan" />
    <ConfirmDialog v-model:open="deleteConfirmOpen" type="delete"
      :description="`将永久删除用户「${pendingDeleteUser?.name || pendingDeleteUser?.username}」，该操作不可恢复。`"
      @confirm="confirmDeleteUser" />

    <!-- 创建超管对话框 -->
    <Dialog v-model:open="showCreate">
      <DialogContent class="rounded-3xl">
        <DialogHeader>
          <DialogTitle>创建超级管理员</DialogTitle>
          <DialogDescription>
            填写用户名、密码和姓名，创建新的超级管理员账号。
          </DialogDescription>
        </DialogHeader>
        <form class="flex flex-col gap-3" @submit.prevent="createSuperAdmin">
          <div class="flex flex-col gap-1">
            <Label>用户名</Label>
            <Input v-model="createForm.username" placeholder="admin" />
          </div>
          <div class="flex flex-col gap-1">
            <Label>密码</Label>
            <Input v-model="createForm.password" type="password" placeholder="至少8位" />
          </div>
          <div class="flex flex-col gap-1">
            <Label>姓名</Label>
            <Input v-model="createForm.name" placeholder="管理员" />
          </div>
          <DialogFooter>
            <Tooltip>
              <TooltipTrigger as-child>
                <Button type="submit" size="icon" class="rounded-full">
                  <Check class="size-4" />
                </Button>
              </TooltipTrigger>
              <TooltipContent>创建</TooltipContent>
            </Tooltip>
          </DialogFooter>
        </form>
      </DialogContent>
    </Dialog>

    <Dialog v-model:open="showEdit">
      <DialogContent class="rounded-3xl">
        <DialogHeader>
          <DialogTitle>编辑用户</DialogTitle>
          <DialogDescription>
            更新用户姓名与邮箱信息。
          </DialogDescription>
        </DialogHeader>
        <form class="flex flex-col gap-3" @submit.prevent="updateUser">
          <div class="flex flex-col gap-1">
            <Label>姓名</Label>
            <Input v-model="editForm.name" placeholder="管理员" />
          </div>
          <div class="flex flex-col gap-1">
            <Label>邮箱（选填）</Label>
            <Input v-model="editForm.email" placeholder="email@example.com" />
          </div>
          <DialogFooter>
            <Tooltip>
              <TooltipTrigger as-child>
                <Button type="submit" size="icon" class="rounded-full">
                  <Check class="size-4" />
                </Button>
              </TooltipTrigger>
              <TooltipContent>保存</TooltipContent>
            </Tooltip>
          </DialogFooter>
        </form>
      </DialogContent>
    </Dialog>

    <!-- 重置密码对话框 -->
    <Dialog v-model:open="showResetPwd">
      <DialogContent class="rounded-3xl">
        <DialogHeader>
          <DialogTitle>重置密码</DialogTitle>
          <DialogDescription>
            为指定用户设置新的登录密码，请确保符合密码安全要求。
          </DialogDescription>
        </DialogHeader>
        <form class="flex flex-col gap-3" @submit.prevent="resetPassword">
          <p class="text-sm text-muted-foreground">为用户 <span class="font-medium text-foreground">{{ resetUserName
          }}</span> 设置新密码</p>
          <div class="flex flex-col gap-1">
            <Label>新密码</Label>
            <Input v-model="newPassword" type="password" placeholder="至少8位" />
          </div>
          <DialogFooter>
            <Tooltip>
              <TooltipTrigger as-child>
                <Button type="submit" size="icon" class="rounded-full">
                  <Check class="size-4" />
                </Button>
              </TooltipTrigger>
              <TooltipContent>确认重置</TooltipContent>
            </Tooltip>
          </DialogFooter>
        </form>
      </DialogContent>
    </Dialog>
  </div>
</template>
