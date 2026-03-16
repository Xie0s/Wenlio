<!--
  RegisterForm.vue - 注册表单组件

  职责：收集租户与管理员信息，完成基础校验后调用认证状态管理进行注册，并在成功后进入管理后台。
  对外暴露：默认组件，无额外属性或事件。
-->
<script setup lang="ts">
import { computed, ref } from 'vue'
import { useRouter } from 'vue-router'
import { Building2, Check, Loader2, UserRound } from 'lucide-vue-next'
import { toast } from 'vue-sonner'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Tooltip, TooltipContent, TooltipProvider, TooltipTrigger } from '@/components/ui/tooltip'
import SecurityCheck from '@/components/auth/SecurityCheck.vue'
import { useAuthStore } from '@/stores/auth'
import type { RegisterTenantReq } from '@/utils/types'

const router = useRouter()
const authStore = useAuthStore()

const form = ref<RegisterTenantReq>({
  tenant_id: '',
  tenant_name: '',
  admin_username: '',
  admin_password: '',
  admin_name: '',
  captcha_token: '',
})
const confirmPassword = ref('')
const loading = ref(false)
const securityCheckRef = ref<InstanceType<typeof SecurityCheck> | null>(null)

const tenantIdHint = computed(() => form.value.tenant_id.trim().toLowerCase())
const normalizedForm = computed(() => ({
  ...form.value,
  tenant_id: form.value.tenant_id.trim().toLowerCase(),
  tenant_name: form.value.tenant_name.trim(),
  admin_name: form.value.admin_name.trim(),
  admin_username: form.value.admin_username.trim(),
}))
const validationMessage = computed(() => {
  const { tenant_name, tenant_id, admin_name, admin_username, admin_password, captcha_token } = normalizedForm.value

  if (!tenant_name || !tenant_id || !admin_name || !admin_username || !admin_password) {
    return '请先填写完整的必填信息'
  }

  if (!isValidTenantId(tenant_id)) {
    return '租户 ID 需为 3-32 位小写字母、数字或中划线，且不能以中划线开头或结尾'
  }

  if (admin_password !== confirmPassword.value) {
    return '两次输入的密码不一致'
  }

  if (!isStrongPassword(admin_password)) {
    return '密码至少 8 位，且必须包含字母和数字'
  }

  if (!captcha_token) {
    return '请先完成安全检查'
  }

  return ''
})
const canSubmit = computed(() => !validationMessage.value)

function isValidTenantId(value: string) {
  return /^[a-z0-9][a-z0-9-]{1,30}[a-z0-9]$/.test(value)
}

function isStrongPassword(value: string) {
  return /^(?=.*[A-Za-z])(?=.*\d).{8,}$/.test(value)
}

function resetSecurityCheck() {
  securityCheckRef.value?.reset()
  form.value.captcha_token = ''
}

async function handleRegister() {
  if (loading.value) {
    return
  }

  if (validationMessage.value) {
    toast.error(validationMessage.value)
    return
  }

  loading.value = true
  try {
    const result = await authStore.register(normalizedForm.value)
    if (result.success) {
      toast.success('注册成功')
      await router.push('/admin')
      return
    }

    resetSecurityCheck()
    toast.error(result.message || '注册失败')
  } catch {
    resetSecurityCheck()
    toast.error('注册请求失败，请稍后重试')
  } finally {
    loading.value = false
  }
}
</script>

  <template>
  <form class="flex flex-col gap-3" @submit.prevent="handleRegister">

    <div class="grid gap-4 lg:grid-cols-2">

      <div class="rounded-3xl border border-border/50 bg-background/35 p-5 backdrop-blur-md">
        <div class="mb-4 flex items-center gap-3.5">
          <div class="flex h-8 w-8 shrink-0 items-center justify-center rounded-2xl bg-primary/10 text-primary">
            <Building2 class="h-4 w-4" />
          </div>
          <div class="min-w-0 flex-1">
            <p class="text-base font-medium text-foreground">租户信息</p>
            <p class="truncate text-[13px] text-muted-foreground">名称与路径标识符</p>
          </div>
          <span class="shrink-0 rounded-full border border-border/50 bg-background/60 px-3 py-1 text-xs text-muted-foreground">租户注册</span>
        </div>
        <div class="flex flex-col gap-3">
          <div class="flex flex-col gap-1.5">
            <Label for="tenant_name">租户名称 <span class="text-destructive">*</span></Label>
            <Input id="tenant_name" v-model="form.tenant_name" placeholder="如 极客团队" class="h-11 rounded-full bg-background/60 text-sm" />
          </div>
          <div class="flex flex-col gap-1.5">
            <Label for="tenant_id">租户 ID <span class="text-destructive">*</span></Label>
            <Input id="tenant_id" v-model="form.tenant_id" placeholder="如 acme-docs" autocomplete="off" maxlength="32" class="h-11 rounded-full bg-background/60 text-sm" />
          </div>
          <div>
            <p class="text-[13px] text-muted-foreground">路径预览：{{ tenantIdHint ? `/${tenantIdHint}` : '/your-tenant-id' }}</p>
          </div>
        </div>
      </div>

      <div class="rounded-3xl border border-border/50 bg-background/35 p-5 backdrop-blur-md">
        <div class="mb-4 flex items-center gap-3.5">
          <div class="flex h-8 w-8 shrink-0 items-center justify-center rounded-2xl bg-primary/10 text-primary">
            <UserRound class="h-4 w-4" />
          </div>
          <div class="min-w-0 flex-1">
            <p class="text-base font-medium text-foreground">管理员账户</p>
            <p class="text-[13px] text-muted-foreground">注册成功后自动登录</p>
          </div>
        </div>
        <div class="grid gap-3 lg:grid-cols-2">
          <div class="flex flex-col gap-1.5">
            <Label for="admin_name">姓名 <span class="text-destructive">*</span></Label>
            <Input id="admin_name" v-model="form.admin_name" placeholder="管理员姓名" class="h-11 rounded-full bg-background/60 text-sm" />
          </div>
          <div class="flex flex-col gap-1.5">
            <Label for="admin_username">用户名 <span class="text-destructive">*</span></Label>
            <Input id="admin_username" v-model="form.admin_username" placeholder="登录用户名" autocomplete="username" maxlength="64" class="h-11 rounded-full bg-background/60 text-sm" />
          </div>
          <div class="flex flex-col gap-1.5">
            <Label for="admin_password">密码 <span class="text-destructive">*</span></Label>
            <Input id="admin_password" v-model="form.admin_password" type="password" placeholder="至少 8 位含字母数字" autocomplete="new-password" class="h-11 rounded-full bg-background/60 text-sm" />
          </div>
          <div class="flex flex-col gap-1.5">
            <Label for="confirm_password">确认密码 <span class="text-destructive">*</span></Label>
            <Input id="confirm_password" v-model="confirmPassword" type="password" placeholder="再次输入密码" autocomplete="new-password" class="h-11 rounded-full bg-background/60 text-sm" />
          </div>
        </div>
      </div>

    </div>

    <!-- 安全校验 + 提交：同一行，校验完成后按钮才可点击 -->
    <div class="mx-auto flex w-full max-w-[460px] items-center gap-3">
      <div class="flex-1">
        <SecurityCheck ref="securityCheckRef" v-model="form.captcha_token" scene="register" :disabled="loading" />
      </div>
      <TooltipProvider>
        <Tooltip>
          <TooltipTrigger as-child>
            <button
              type="submit"
              class="inline-flex h-[52px] w-[52px] shrink-0 items-center justify-center rounded-full bg-primary text-primary-foreground transition-all hover:bg-primary/90 disabled:cursor-not-allowed disabled:bg-muted disabled:text-muted-foreground disabled:opacity-60"
              :disabled="loading || !canSubmit"
            >
              <Loader2 v-if="loading" class="h-[18px] w-[18px] animate-spin" />
              <Check v-else class="h-[18px] w-[18px]" />
            </button>
          </TooltipTrigger>
          <TooltipContent>注册</TooltipContent>
        </Tooltip>
      </TooltipProvider>
    </div>

  </form>
</template>
