<!-- TenantListPage.vue - 租户管理页面（超管专属）
     职责：纯 UI 渲染与用户交互，业务逻辑由 lib/tenant-list.ts 提供 -->
<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { useTenantList } from '@/lib/tenant-list'
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
import { Plus, Ban, CheckCircle, Search, Check, Pencil, Trash2, X, MoreHorizontal } from 'lucide-vue-next'
import UserHoverCard from '@/components/common/UserHoverCard.vue'
import ConfirmDialog from '@/components/common/ConfirmDialog.vue'

function formatDate(dateStr: string) {
  if (!dateStr || dateStr.startsWith('0001')) return '-'
  return new Date(dateStr).toLocaleDateString('zh-CN', { year: 'numeric', month: '2-digit', day: '2-digit' })
}

const {
  tenants, keyword, loading,
  showCreate, form,
  showEdit, editForm,
  loadTenants, createTenant, openEditTenant, updateTenant, toggleStatus, deleteTenant,
} = useTenantList()
const appliedKeyword = ref('')
const isSearchApplied = computed(() => !!appliedKeyword.value && appliedKeyword.value === keyword.value)
const openActionTenantId = ref('')

const banConfirmOpen = ref(false)
const pendingBanTenant = ref<(typeof tenants.value)[number] | null>(null)
const deleteConfirmOpen = ref(false)
const pendingDeleteTenant = ref<(typeof tenants.value)[number] | null>(null)

function requestToggleStatus(t: (typeof tenants.value)[number]) {
  if (t.status === 'active') {
    pendingBanTenant.value = t
    banConfirmOpen.value = true
  } else {
    toggleStatus(t)
  }
}

function confirmBan() {
  if (pendingBanTenant.value) toggleStatus(pendingBanTenant.value)
  pendingBanTenant.value = null
}

function requestDeleteTenant(t: (typeof tenants.value)[number]) {
  pendingDeleteTenant.value = t
  deleteConfirmOpen.value = true
}

function confirmDeleteTenant() {
  if (pendingDeleteTenant.value) deleteTenant(pendingDeleteTenant.value)
  pendingDeleteTenant.value = null
}

function handleActionMenuOpenChange(tenantId: string, open: boolean) {
  openActionTenantId.value = open ? tenantId : ''
}

function executeSearch() {
  keyword.value = keyword.value.trim()
  appliedKeyword.value = keyword.value
  loadTenants()
}

function handleSearchAction() {
  if (isSearchApplied.value) {
    keyword.value = ''
    appliedKeyword.value = ''
    loadTenants()
    return
  }
  executeSearch()
}

onMounted(loadTenants)
</script>

<template>
  <div class="space-y-6 text-base">
    <!-- 搜索 -->
    <div class="flex gap-2">
      <Input v-model="keyword" placeholder="搜索租户名称或租户 ID..." class="max-w-xs" @keyup.enter="executeSearch" />
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
        <TooltipContent>创建租户</TooltipContent>
      </Tooltip>
    </div>

    <!-- 数据表 -->
    <div class="max-h-[420px] overflow-auto rounded-xl border">
      <Table class="w-full table-fixed text-lg">
        <TableHeader>
          <TableRow>
            <TableHead class="h-11 w-[16%] text-lg font-light">租户名称</TableHead>
            <TableHead class="h-11 w-[14%] text-lg font-light">租户 ID</TableHead>
            <TableHead class="h-11 w-[14%] text-lg font-light">管理员用户名</TableHead>
            <TableHead class="h-11 w-[12%] text-lg font-light">管理员姓名</TableHead>
            <TableHead class="h-11 w-[8%] text-lg font-light">用户数</TableHead>
            <TableHead class="h-11 w-[14%] text-lg font-light">创建时间</TableHead>
            <TableHead class="h-11 w-[10%] text-lg font-light">状态</TableHead>
            <TableHead class="h-11 w-[8%] text-right text-lg font-light">操作</TableHead>
          </TableRow>
        </TableHeader>
        <TableBody>
          <TableRow v-for="t in tenants" :key="t.id" class="h-12">
            <TableCell class="truncate py-3 text-lg font-light">{{ t.name }}</TableCell>
            <TableCell class="truncate py-3 text-lg font-light text-muted-foreground">{{ t.id }}</TableCell>
            <TableCell class="truncate py-3 text-lg font-light text-muted-foreground">{{ t.admin_username || '-' }}
            </TableCell>
            <TableCell class="py-3 text-lg font-light text-muted-foreground">
              <UserHoverCard v-if="t.admin_name" :avatar-url="t.admin_avatar_url" :name="t.admin_name"
                :username="t.admin_username || ''" role-label="租户管理员" :tenant-name="t.name" />
              <span v-else>-</span>
            </TableCell>
            <TableCell class="py-3 text-lg font-light text-muted-foreground">{{ t.user_count ?? 0 }}</TableCell>
            <TableCell class="py-3 text-lg font-light text-muted-foreground">{{ formatDate(t.created_at) }}</TableCell>
            <TableCell class="py-3">
              <Badge
                :variant="t.status === 'active' ? 'default' : t.status === 'deleting' ? 'secondary' : 'destructive'"
                class="px-2.5 py-1 text-sm font-light">
                {{ t.status === 'active' ? '活跃' : t.status === 'deleting' ? '删除中' : '已封禁' }}
              </Badge>
            </TableCell>
            <TableCell class="py-3 text-right">
              <DropdownMenu :modal="false" @update:open="(open: boolean) => handleActionMenuOpenChange(t.id, open)">
                <DropdownMenuTrigger as-child>
                  <Button variant="ghost" size="icon" class="h-8 w-8 rounded-full hover:bg-foreground/10"
                    :class="{ 'bg-foreground/10': openActionTenantId === t.id }">
                    <MoreHorizontal class="h-4 w-4" />
                  </Button>
                </DropdownMenuTrigger>
                <DropdownMenuContent align="end" class="rounded-3xl">
                  <DropdownMenuItem :disabled="t.status === 'deleting'" @click="openEditTenant(t)">
                    <Pencil class="mr-2 h-4 w-4" />
                    {{ t.status === 'deleting' ? '删除中不可编辑' : '编辑' }}
                  </DropdownMenuItem>
                  <DropdownMenuItem :disabled="t.status === 'deleting'" @click="requestToggleStatus(t)">
                    <Ban v-if="t.status === 'active'" class="mr-2 h-4 w-4" />
                    <CheckCircle v-else class="mr-2 h-4 w-4" />
                    {{ t.status === 'deleting' ? '删除中不可变更状态' : t.status === 'active' ? '封禁' : '解封' }}
                  </DropdownMenuItem>
                  <DropdownMenuSeparator />
                  <DropdownMenuItem variant="destructive" :disabled="t.status === 'deleting'"
                    @click="requestDeleteTenant(t)">
                    <Trash2 class="mr-2 h-4 w-4" />
                    {{ t.status === 'deleting' ? '删除任务执行中' : '删除' }}
                  </DropdownMenuItem>
                </DropdownMenuContent>
              </DropdownMenu>
            </TableCell>
          </TableRow>
        </TableBody>
      </Table>
      <p v-if="!loading && tenants.length === 0" class="py-6 text-center text-base text-muted-foreground">暂无租户</p>
    </div>

    <ConfirmDialog v-model:open="banConfirmOpen" type="ban"
      :description="`封禁后租户「${pendingBanTenant?.name}」将无法登录后台、无法访问读者端页面，可随时在此页面解封。`" @confirm="confirmBan" />
    <ConfirmDialog v-model:open="deleteConfirmOpen" type="delete"
      :description="`将为租户「${pendingDeleteTenant?.name}」提交异步删除任务，系统会立即阻断访问，并在后台清理其用户、主题、版本、文档、评论与媒体文件。`"
      @confirm="confirmDeleteTenant" />

    <!-- 创建对话框 -->
    <Dialog v-model:open="showCreate">
      <DialogContent class="rounded-3xl">
        <DialogHeader>
          <DialogTitle>创建租户</DialogTitle>
          <DialogDescription>
            填写租户基础信息并创建初始管理员账号。
          </DialogDescription>
        </DialogHeader>
        <form class="flex flex-col gap-3" @submit.prevent="createTenant">
          <div class="flex flex-col gap-1">
            <Label>租户 ID（URL 路径段，创建后不可修改）</Label>
            <Input v-model="form.id" placeholder="如 acme" />
          </div>
          <div class="flex flex-col gap-1">
            <Label>租户名称</Label>
            <Input v-model="form.name" placeholder="如 Acme Corp" />
          </div>
          <div class="flex flex-col gap-1">
            <Label>初始管理员用户名</Label>
            <Input v-model="form.admin_username" placeholder="admin" />
          </div>
          <div class="flex flex-col gap-1">
            <Label>初始管理员密码</Label>
            <Input v-model="form.admin_password" type="password" placeholder="至少8位" />
          </div>
          <div class="flex flex-col gap-1">
            <Label>初始管理员姓名</Label>
            <Input v-model="form.admin_name" placeholder="管理员" />
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
          <DialogTitle>编辑租户</DialogTitle>
          <DialogDescription>
            更新租户名称与 Logo 地址，租户 ID 不支持修改。
          </DialogDescription>
        </DialogHeader>
        <form class="flex flex-col gap-3" @submit.prevent="updateTenant">
          <div class="flex flex-col gap-1">
            <Label>租户名称</Label>
            <Input v-model="editForm.name" placeholder="如 Acme Corp" />
          </div>
          <div class="flex flex-col gap-1">
            <Label>Logo 地址（选填）</Label>
            <Input v-model="editForm.logo_url" placeholder="https://example.com/logo.png" />
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
  </div>
</template>
