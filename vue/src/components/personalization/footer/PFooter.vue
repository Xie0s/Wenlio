<!-- PFooter.vue - 页脚展示组件
     职责：展示 Logo、标语、自定义链接、版权信息和备案信息
     对外接口：
       Props: config, tenantName, tenantLogo -->
<script setup lang="ts">
import { computed } from 'vue'
import type { FooterConfig } from '../types'

const props = defineProps<{
  config: FooterConfig
  tenantName: string
  tenantLogo: string
}>()

const displayLogoUrl = computed(() => props.config.logo_url || props.tenantLogo)
const displaySiteName = computed(() => props.config.site_name || props.tenantName)

const hasCopyrightOrFiling = computed(() => !!(
  props.config.copyright.owner ||
  props.config.site_name ||
  props.config.filing.icp_number ||
  props.config.filing.police_number
))
</script>

<template>
  <footer class="relative w-full border-t border-border dark:border-border/50 mt-20">
    <div 
      class="mx-auto px-6 pt-6 lg:pt-10 pb-6" 
      style="max-width: var(--hp-max-width, 1200px)"
    >
      <div class="flex flex-col items-center space-y-6 text-center">
        <!-- Logo 和标语 -->
        <div v-if="displayLogoUrl || config.slogan" class="flex flex-col items-center space-y-4">
          <div v-if="displayLogoUrl" class="flex justify-center">
            <img
              :src="displayLogoUrl"
              alt="Logo"
              class="w-auto h-auto max-w-full object-contain dark:brightness-0 dark:invert"
              style="max-height: 80px; width: auto;"
            />
          </div>
          <p v-if="config.slogan" class="text-xl font-light text-muted-foreground max-w-lg leading-relaxed">
            {{ config.slogan }}
          </p>
        </div>

        <!-- 自定义链接 -->
        <div v-if="config.custom_links.length > 0" class="flex flex-wrap justify-center items-center gap-x-6 gap-y-2">
          <component
            v-for="(link, i) in config.custom_links"
            :key="i"
            :is="link.url ? 'a' : 'span'"
            :href="link.url || undefined"
            target="_blank"
            rel="noopener noreferrer"
            class="text-base text-muted-foreground hover:text-primary transition-colors focus:text-foreground focus:outline-none focus:ring-2 focus:ring-primary focus:ring-offset-2 rounded-sm"
          >
            {{ link.name }}
          </component>
        </div>
      </div>

      <!-- 底部版权和备案信息 -->
      <div v-if="hasCopyrightOrFiling" class="mt-6 border-t border-border/60 dark:border-border/40 pt-4">
        <div class="flex flex-col lg:flex-row lg:justify-between lg:items-center gap-2 text-center lg:text-left">
          <!-- 版权信息 -->
          <p v-if="config.copyright.owner" class="text-base text-muted-foreground">
            <component
              :is="config.copyright.link ? 'a' : 'span'"
              :href="config.copyright.link || undefined"
              :target="config.copyright.link ? '_blank' : undefined"
              :rel="config.copyright.link ? 'noopener noreferrer' : undefined"
              class="text-muted-foreground hover:text-foreground focus:text-foreground focus:outline-none focus:ring-2 focus:ring-primary focus:ring-offset-2 rounded-sm transition-colors duration-200 truncate max-w-[320px] sm:max-w-none"
              :title="`${config.copyright.owner} 版权所有 ©${config.copyright.year}`"
            >
              {{ config.copyright.owner }} 版权所有 ©{{ config.copyright.year }}
            </component>
          </p>

          <!-- 网站名称 -->
          <p v-if="displaySiteName" class="text-base text-muted-foreground">
            {{ displaySiteName }}
          </p>

          <!-- 备案信息 -->
          <div v-if="config.filing.icp_number || config.filing.police_number" class="flex flex-col sm:flex-row sm:items-center gap-2 justify-center lg:justify-end">
            <component
              v-if="config.filing.icp_number"
              :is="config.filing.icp_link ? 'a' : 'span'"
              :href="config.filing.icp_link || undefined"
              :target="config.filing.icp_link ? '_blank' : undefined"
              :rel="config.filing.icp_link ? 'noopener noreferrer' : undefined"
              :title="config.filing.icp_number"
              class="text-base text-muted-foreground hover:text-foreground transition-colors truncate max-w-[260px] sm:max-w-none"
            >
              {{ config.filing.icp_number }}
            </component>
            <component
              v-if="config.filing.police_number"
              :is="config.filing.police_link ? 'a' : 'span'"
              :href="config.filing.police_link || undefined"
              :target="config.filing.police_link ? '_blank' : undefined"
              :rel="config.filing.police_link ? 'noopener noreferrer' : undefined"
              :title="config.filing.police_number"
              class="text-base text-muted-foreground hover:text-foreground transition-colors truncate max-w-[260px] sm:max-w-none"
            >
              {{ config.filing.police_number }}
            </component>
          </div>
        </div>
      </div>
    </div>
  </footer>
</template>
