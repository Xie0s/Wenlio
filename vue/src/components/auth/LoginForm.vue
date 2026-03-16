<!--
  LoginForm.vue - 登录表单组件

  职责：收集用户名与密码，调用认证状态管理完成登录，并在成功后跳转至管理后台首页。
  对外暴露：默认组件，无额外属性或事件。
-->
<script setup lang="ts">
import { computed, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { Check, Loader2, LockKeyhole, User2 } from 'lucide-vue-next'
import { toast } from 'vue-sonner'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Tooltip, TooltipContent, TooltipProvider, TooltipTrigger } from '@/components/ui/tooltip'
import SecurityCheck from '@/components/auth/SecurityCheck.vue'
import { useAuthStore } from '@/stores/auth'

const route = useRoute()
const router = useRouter()
const authStore = useAuthStore()

const username = ref('')
const password = ref('')
const captchaToken = ref('')
const loading = ref(false)
const securityCheckRef = ref<InstanceType<typeof SecurityCheck> | null>(null)
const canSubmit = computed(() => Boolean(username.value.trim() && password.value && captchaToken.value))

function resetSecurityCheck() {
  securityCheckRef.value?.reset()
}

async function handleLogin() {
  if (loading.value) {
    return
  }

  if (!username.value || !password.value) {
    toast.error('请输入用户名和密码')
    return
  }

  if (!captchaToken.value) {
    toast.error('请先完成安全检查')
    return
  }

  loading.value = true
  try {
    const result = await authStore.login({
      username: username.value,
      password: password.value,
      captcha_token: captchaToken.value,
    })

    if (result.success) {
      toast.success('登录成功')
      const redirectTo = (route.query.redirect as string) || '/admin'
      await router.push(redirectTo)
      return
    }

    resetSecurityCheck()
    toast.error(result.message || '登录失败')
  } catch {
    resetSecurityCheck()
    toast.error('登录请求失败，请稍后重试')
  } finally {
    loading.value = false
  }
}
</script>

  <template>
  <form class="flex flex-col gap-3" @submit.prevent="handleLogin">
    <div class="rounded-3xl border border-border/50 bg-background/35 p-5 backdrop-blur-md">
      <div class="mb-4 flex items-start justify-between gap-4">
        <div>
          <p class="text-base font-medium text-foreground">账户验证</p>
          <p class="mt-1 text-[13px] leading-5 text-muted-foreground">请输入管理员用户名和密码，完成安全检查后进入后台。</p>
        </div>
        <div class="rounded-full border border-border/50 bg-background/60 px-3 py-1 text-xs text-muted-foreground">
          安全认证
        </div>
      </div>

      <div class="grid gap-3">
        <div class="flex flex-col gap-1.5">
          <Label for="username">用户名</Label>
          <div class="relative">
            <User2 class="pointer-events-none absolute left-4 top-1/2 h-4 w-4 -translate-y-1/2 text-muted-foreground/70" />
            <Input
              id="username"
              v-model="username"
              placeholder="请输入用户名"
              autocomplete="username"
              class="h-11 rounded-full bg-background/60 pl-11 text-sm"
            />
          </div>
        </div>

        <div class="flex flex-col gap-1.5">
          <div class="flex items-center justify-between gap-3">
            <Label for="password">密码</Label>
            <span class="text-xs text-muted-foreground">区分大小写</span>
          </div>
          <div class="relative">
            <LockKeyhole class="pointer-events-none absolute left-4 top-1/2 h-4 w-4 -translate-y-1/2 text-muted-foreground/70" />
            <Input
              id="password"
              v-model="password"
              type="password"
              placeholder="请输入密码"
              autocomplete="current-password"
              class="h-11 rounded-full bg-background/60 pl-11 text-sm"
            />
          </div>
        </div>
      </div>
    </div>
  
    <!-- 安全校验 + 提交：同一行，校验完成后按钮才可点击 -->
    <div class="mx-auto flex w-full max-w-[460px] items-center gap-3">
      <div class="flex-1">
        <SecurityCheck ref="securityCheckRef" v-model="captchaToken" scene="login" :disabled="loading" />
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
          <TooltipContent>登录</TooltipContent>
        </Tooltip>
      </TooltipProvider>
    </div>
  </form>
</template>
