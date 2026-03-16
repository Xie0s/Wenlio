<!--
  ProfilePage.vue - 管理后台个人中心页面

  职责：展示/编辑当前登录用户的个人资料（头像、背景、姓名、签名）并提供修改密码入口。
  功能边界：
  1) 头像 / 背景图通过 authStore.uploadAvatar / uploadProfileBg 上传；
  2) 姓名 / 签名通过 authStore.updateProfile 更新；
  3) 密码修改通过 authStore.changePassword 调用后端 /auth/me/password。

  对外暴露：
  - 路由页面组件 <ProfilePage />（由 /admin/profile 路由挂载）
-->
<script setup lang="ts">
import { computed, ref, nextTick } from 'vue'
import { useAuthStore } from '@/stores/auth'
import { Input } from '@/components/ui/input'
import { Textarea } from '@/components/ui/textarea'
import { Button } from '@/components/ui/button'
import { Tooltip, TooltipContent, TooltipTrigger } from '@/components/ui/tooltip'
import { Avatar, AvatarImage, AvatarFallback } from '@/components/ui/avatar'
import { Separator } from '@/components/ui/separator'
import { Check, X, Camera, ImagePlus, Pencil, ChevronDown } from 'lucide-vue-next'
import { toast } from 'vue-sonner'

const authStore = useAuthStore()

const displayName = computed(() => authStore.user?.name || authStore.user?.username || '-')
const roleLabel = computed(() => (authStore.isSuperAdmin ? '超级管理员' : '租户管理员'))
const avatarFallback = computed(() => displayName.value.charAt(0).toUpperCase())

// ── 头像上传 ──
const avatarInputRef = ref<HTMLInputElement | null>(null)

async function handleAvatarChange(e: Event) {
  const file = (e.target as HTMLInputElement).files?.[0]
  if (!file) return
  const result = await authStore.uploadAvatar(file)
  if (!result.success) { toast.error(result.message || '头像上传失败'); return }
  toast.success('头像已更新')
}

// ── 背景上传 ──
const bgInputRef = ref<HTMLInputElement | null>(null)

async function handleBgChange(e: Event) {
  const file = (e.target as HTMLInputElement).files?.[0]
  if (!file) return
  const result = await authStore.uploadProfileBg(file)
  if (!result.success) { toast.error(result.message || '背景上传失败'); return }
  toast.success('背景已更新')
}

// ── 姓名编辑 ──
const editingName = ref(false)
const nameInput = ref('')
const nameInputRef = ref<HTMLInputElement | null>(null)

function startEditName() {
  nameInput.value = authStore.user?.name || ''
  editingName.value = true
  nextTick(() => nameInputRef.value?.focus())
}

async function saveName() {
  const trimmed = nameInput.value.trim()
  if (!trimmed) { toast.error('姓名不能为空'); return }
  const result = await authStore.updateProfile({ name: trimmed })
  if (!result.success) { toast.error(result.message || '更新失败'); return }
  editingName.value = false
  toast.success('姓名已更新')
}

// ── 签名编辑 ──
const editingBio = ref(false)
const bioInput = ref('')
const bioInputRef = ref<InstanceType<typeof Textarea> | null>(null)

function startEditBio() {
  bioInput.value = authStore.user?.bio || ''
  editingBio.value = true
  nextTick(() => bioInputRef.value?.$el?.focus())
}

async function saveBio() {
  const result = await authStore.updateProfile({ bio: bioInput.value.trim() })
  if (!result.success) { toast.error(result.message || '更新失败'); return }
  editingBio.value = false
  toast.success('签名已更新')
}

// ── 修改密码 ──
const showPasswordForm = ref(false)
const oldPassword = ref('')
const newPassword = ref('')
const confirmPassword = ref('')
const submitting = ref(false)

async function handleChangePassword() {
  if (!oldPassword.value || !newPassword.value || !confirmPassword.value) {
    toast.error('请完整填写密码信息'); return
  }
  if (newPassword.value !== confirmPassword.value) {
    toast.error('两次输入的新密码不一致'); return
  }
  submitting.value = true
  const result = await authStore.changePassword(oldPassword.value, newPassword.value)
  submitting.value = false
  if (!result.success) { toast.error(result.message || '修改密码失败'); return }
  oldPassword.value = ''
  newPassword.value = ''
  confirmPassword.value = ''
  showPasswordForm.value = false
  toast.success('密码修改成功')
}
</script>

<template>
  <div class="max-w-2xl mx-auto">
    <div class="rounded-3xl bg-card ring-1 ring-border/50 overflow-hidden">

      <!-- ═══ 背景封面（头像+姓名+用户名全部在背景内） ═══ -->
      <div
        class="group/bg relative m-3 rounded-2xl overflow-hidden bg-gradient-to-br from-primary/30 via-primary/15 to-muted/50 bg-cover bg-center"
        :style="authStore.user?.profile_bg_url ? { backgroundImage: `url(${authStore.user.profile_bg_url})` } : {}">
        <!-- 底部渐变（保证白色文字在任何背景上可读） -->
        <div class="absolute inset-0 bg-gradient-to-t from-black/75 via-black/30 to-transparent pointer-events-none" />

        <!-- 更换封面 -->
        <button
          class="absolute top-2.5 right-2.5 z-20 rounded-full bg-black/30 hover:bg-black/50 p-1.5 opacity-0 group-hover/bg:opacity-100 transition-opacity"
          @click="bgInputRef?.click()">
          <ImagePlus class="size-3.5 text-white/80" :stroke-width="1.5" />
        </button>
        <input ref="bgInputRef" type="file" accept="image/*" class="hidden" @change="handleBgChange" />

        <!-- 占位撑高 + 底部内容 -->
        <div class="pt-28" />
        <div class="relative z-10 px-5 pb-4 flex items-end gap-3.5">
          <!-- 头像 -->
          <div class="relative shrink-0 group cursor-pointer" @click="avatarInputRef?.click()">
            <Avatar class="size-[4.5rem] border-[3px] border-white/25 shadow-xl">
              <AvatarImage v-if="authStore.user?.avatar_url" :src="authStore.user.avatar_url" />
              <AvatarFallback class="text-xl font-semibold bg-white/15 text-white backdrop-blur-sm">
                {{ avatarFallback }}
              </AvatarFallback>
            </Avatar>
            <div
              class="absolute inset-0 rounded-full flex items-center justify-center opacity-0 group-hover:opacity-100 transition-opacity bg-black/40">
              <Camera class="size-4 text-white" :stroke-width="1.5" />
            </div>
            <input ref="avatarInputRef" type="file" accept="image/*" class="hidden" @change="handleAvatarChange" />
          </div>

          <!-- 姓名 + 用户名 -->
          <div class="min-w-0 pb-0.5">
            <div class="flex items-center gap-1.5">
              <template v-if="editingName">
                <Input ref="nameInputRef" v-model="nameInput"
                  class="h-8 w-44 text-sm font-semibold bg-black/25 border-white/20 text-white placeholder:text-white/40 backdrop-blur-sm"
                  @keyup.enter="saveName" @keyup.escape="editingName = false" />
                <button
                  class="rounded-full size-6 flex items-center justify-center bg-white/20 hover:bg-white/30 text-white transition-colors"
                  @click="saveName">
                  <Check class="size-3" />
                </button>
                <button
                  class="rounded-full size-6 flex items-center justify-center bg-white/10 hover:bg-white/20 text-white/70 transition-colors"
                  @click="editingName = false">
                  <X class="size-3" />
                </button>
              </template>
              <template v-else>
                <span class="text-2xl font-bold text-white drop-shadow-md truncate">{{ displayName }}</span>
                <button
                  class="shrink-0 rounded-full size-5 flex items-center justify-center text-white/40 hover:text-white hover:bg-white/15 transition-colors"
                  @click="startEditName">
                  <Pencil class="size-2.5" />
                </button>
              </template>
            </div>
            <span class="text-sm text-white/60 font-light">@{{ authStore.user?.username }}</span>
          </div>
        </div>
      </div>

      <!-- ═══ 个性签名 ═══ -->
      <div class="px-6 pt-4 pb-3">
        <div class="flex items-center gap-2 mb-2">
          <span class="text-xs font-medium text-muted-foreground uppercase tracking-wider">个性签名</span>
          <button v-if="!editingBio"
            class="rounded-full size-5 flex items-center justify-center text-muted-foreground/50 hover:text-foreground hover:bg-muted transition-colors"
            @click="startEditBio">
            <Pencil class="size-2.5" />
          </button>
        </div>
        <template v-if="editingBio">
          <Textarea ref="bioInputRef" v-model="bioInput" class="w-full text-sm min-h-[3.5rem] resize-none"
            placeholder="写一句话介绍自己..." @keydown.escape="editingBio = false" />
          <div class="flex justify-end gap-1 mt-2">
            <Tooltip>
              <TooltipTrigger as-child>
                <Button variant="ghost" size="icon" class="rounded-full size-7" @click="editingBio = false">
                  <X class="size-3.5" />
                </Button>
              </TooltipTrigger>
              <TooltipContent>取消</TooltipContent>
            </Tooltip>
            <Tooltip>
              <TooltipTrigger as-child>
                <Button size="icon" class="rounded-full size-7" @click="saveBio">
                  <Check class="size-3.5" />
                </Button>
              </TooltipTrigger>
              <TooltipContent>保存</TooltipContent>
            </Tooltip>
          </div>
        </template>
        <p v-else class="text-sm cursor-pointer"
          :class="authStore.user?.bio ? 'text-foreground' : 'text-muted-foreground/50 italic'" @click="startEditBio">
          {{ authStore.user?.bio || '点击添加个性签名…' }}
        </p>
      </div>

      <Separator class="mx-6" />

      <!-- ═══ 账号信息 ═══ -->
      <div class="px-6 py-4 space-y-3">
        <span class="text-xs font-medium text-muted-foreground uppercase tracking-wider">账号信息</span>
        <div class="space-y-2.5">
          <div class="flex items-center justify-between">
            <span class="text-sm text-muted-foreground">用户名</span>
            <span class="text-sm font-medium">{{ authStore.user?.username || '-' }}</span>
          </div>
          <div class="flex items-center justify-between">
            <span class="text-sm text-muted-foreground">角色</span>
            <span class="text-xs px-2.5 py-0.5 rounded-full bg-primary/10 text-primary font-medium">{{ roleLabel
              }}</span>
          </div>
          <div v-if="authStore.user?.tenant_name" class="flex items-center justify-between">
            <span class="text-sm text-muted-foreground">所属空间</span>
            <span class="text-sm font-medium">{{ authStore.user.tenant_name }}</span>
          </div>
        </div>
      </div>

      <Separator class="mx-6" />

      <!-- ═══ 安全设置 ═══ -->
      <div class="px-6 py-4">
        <button class="w-full flex items-center justify-between" @click="showPasswordForm = !showPasswordForm">
          <span class="text-xs font-medium text-muted-foreground uppercase tracking-wider">安全设置</span>
          <ChevronDown class="size-3.5 text-muted-foreground transition-transform duration-200"
            :class="{ 'rotate-180': showPasswordForm }" :stroke-width="1.5" />
        </button>

        <Transition name="collapse">
          <div v-if="showPasswordForm" class="mt-4 space-y-3">
            <div class="grid gap-3 sm:grid-cols-2">
              <Input v-model="oldPassword" type="password" placeholder="当前密码" class="sm:col-span-2 h-9" />
              <Input v-model="newPassword" type="password" placeholder="新密码" class="h-9" />
              <Input v-model="confirmPassword" type="password" placeholder="确认新密码" class="h-9" />
            </div>
            <div class="flex justify-end">
              <Tooltip>
                <TooltipTrigger as-child>
                  <Button size="icon" class="rounded-full" :disabled="submitting" @click="handleChangePassword">
                    <Check class="size-4" />
                  </Button>
                </TooltipTrigger>
                <TooltipContent>保存新密码</TooltipContent>
              </Tooltip>
            </div>
          </div>
        </Transition>
      </div>

    </div>
  </div>
</template>

<style scoped>
.collapse-enter-active,
.collapse-leave-active {
  overflow: hidden;
  transition: max-height 0.3s ease, opacity 0.25s ease;
}

.collapse-enter-from,
.collapse-leave-to {
  max-height: 0;
  opacity: 0;
}

.collapse-enter-to,
.collapse-leave-from {
  max-height: 300px;
  opacity: 1;
}
</style>
