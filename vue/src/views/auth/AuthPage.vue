<!--
  AuthPage.vue - 统一认证页面

  职责：根据当前认证路由组装登录或注册布局与表单。
  对外暴露：路由页面组件。
-->
<script setup lang="ts">
import { computed } from 'vue'
import { useRoute } from 'vue-router'
import AuthLayout from '@/components/auth/AuthLayout.vue'
import LoginForm from '@/components/auth/LoginForm.vue'
import RegisterForm from '@/components/auth/RegisterForm.vue'

const route = useRoute()

const authView = computed(() => {
  if (route.path === '/admin/register') {
    return {
      title: '创建租户并注册管理员',
      description: '只需填写必要信息，完成安全校验后即可创建租户与首个管理员账户。',
      cardClass: 'max-w-[860px]',
      mode: 'register' as const,
      component: RegisterForm,
      footerText: '已有管理员账号？',
      footerActionText: '去登录',
      footerTo: '/admin/login',
    }
  }

  return {
    title: '登录文档管理平台',
    description: '输入管理员账户信息并完成安全检查后进入后台。',
    cardClass: 'max-w-[520px]',
    mode: 'login' as const,
    component: LoginForm,
    footerText: '还没有租户空间？',
    footerActionText: '去注册',
    footerTo: '/admin/register',
  }
})
</script>

<template>
  <AuthLayout
    :title="authView.title"
    :description="authView.description"
    :card-class="authView.cardClass"
    :mode="authView.mode"
  >
    <Transition name="auth-panel" mode="out-in">
      <component :is="authView.component" :key="authView.mode" />
    </Transition>
    <template #footer>
      <Transition name="auth-switch" mode="out-in">
        <div
          :key="authView.mode"
          class="flex flex-col gap-2 text-center text-sm text-muted-foreground sm:flex-row sm:items-center sm:justify-between sm:text-left"
        >
          <span>{{ authView.footerText }}</span>
          <RouterLink class="font-medium text-foreground underline underline-offset-4" :to="authView.footerTo">{{ authView.footerActionText }}</RouterLink>
        </div>
      </Transition>
    </template>
  </AuthLayout>
</template>

<style scoped>
.auth-panel-enter-active,
.auth-panel-leave-active {
  transition: opacity 0.32s ease, transform 0.32s ease, filter 0.32s ease;
}

.auth-panel-enter-from {
  opacity: 0;
  transform: translateY(16px) scale(0.985);
  filter: blur(6px);
}

.auth-panel-leave-to {
  opacity: 0;
  transform: translateY(-10px) scale(0.985);
  filter: blur(4px);
}

.auth-switch-enter-active,
.auth-switch-leave-active {
  transition: opacity 0.24s ease, transform 0.24s ease;
}

.auth-switch-enter-from {
  opacity: 0;
  transform: translateY(8px);
}

.auth-switch-leave-to {
  opacity: 0;
  transform: translateY(-6px);
}
</style>
