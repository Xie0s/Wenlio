<!--
  common/UserHoverCard.vue - 用户悬停信息卡组件
  职责：封装用户头像 + TooltipCard 悬停卡片的完整样式，使用方只需传入用户数据即可
  对外暴露：
    Props:
      - avatarUrl?:     string   头像地址
      - name:           string   显示名称
      - username:       string   用户名（@xxx）
      - bio?:           string   个性签名
      - profileBgUrl?:  string   个人背景图
      - roleLabel?:     string   角色显示文字
      - tenantName?:    string   所属租户名称
    Slots:
      - default: 触发区域（表格中的头像+姓名行等）
-->
<script setup lang="ts">
import { Avatar, AvatarImage, AvatarFallback } from '@/components/ui/avatar'
import TooltipCard from '@/components/ui/a-aceternity/TooltipCard.vue'

const props = withDefaults(defineProps<{
  avatarUrl?: string
  name: string
  username: string
  bio?: string
  profileBgUrl?: string
  roleLabel?: string
  tenantName?: string
}>(), {
  avatarUrl: '',
  bio: '',
  profileBgUrl: '',
  roleLabel: '',
  tenantName: '',
})

const fallback = (props.name || props.username || '?').charAt(0).toUpperCase()
</script>

<template>
  <TooltipCard>
    <!-- 触发区域：使用方通过默认 slot 自定义 -->
    <slot>
      <div class="flex items-center gap-2">
        <Avatar class="size-7 shrink-0">
          <AvatarImage v-if="props.avatarUrl" :src="props.avatarUrl" />
          <AvatarFallback class="text-xs font-medium bg-primary/15 text-primary">
            {{ fallback }}
          </AvatarFallback>
        </Avatar>
        <span class="truncate">{{ props.name || props.username }}</span>
      </div>
    </slot>

    <template #content>
      <!-- Banner 短背景条 -->
      <div
        class="-mx-2 -mt-2 md:-mx-4 md:-mt-4 h-20 rounded-t-3xl bg-gradient-to-br from-primary/35 via-primary/15 to-muted/40 bg-cover bg-center"
        :style="props.profileBgUrl ? { backgroundImage: `url(${props.profileBgUrl})` } : {}" />

      <!-- 头像：上移与 banner 交叠 -->
      <div class="px-1 -mt-7">
        <Avatar class="size-14 ring-[3px] ring-background shadow-md">
          <AvatarImage v-if="props.avatarUrl" :src="props.avatarUrl" />
          <AvatarFallback class="text-base font-semibold bg-primary/15 text-primary">
            {{ fallback }}
          </AvatarFallback>
        </Avatar>
      </div>

      <!-- 用户信息 -->
      <div class="px-1 pt-2 pb-0.5 space-y-1.5">
        <div>
          <div class="flex items-center gap-2">
            <p class="text-sm font-bold truncate leading-tight">{{ props.name || props.username }}</p>
            <span v-if="props.roleLabel"
              class="shrink-0 px-1.5 py-px rounded-full bg-primary/10 text-primary text-[10px] font-medium">
              {{ props.roleLabel }}
            </span>
          </div>
          <p class="text-xs text-muted-foreground">@{{ props.username }}</p>
        </div>

        <p v-if="props.bio" class="text-xs leading-relaxed text-muted-foreground/80 line-clamp-2">{{ props.bio }}</p>

        <p v-if="props.tenantName" class="text-[11px] text-muted-foreground/70">{{ props.tenantName }}</p>
      </div>
    </template>
  </TooltipCard>
</template>
