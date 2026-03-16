<!-- FooterEditor.vue - 页脚区块配置编辑器
     职责：编辑 FooterConfig 的 Logo、标语、自定义链接、版权信息和备案信息
     对外接口：
       Props: config
       Emits: update(config: FooterConfig) -->
<script setup lang="ts">
import { reactive, watch } from 'vue'
import type { FooterConfig } from '../types'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Textarea } from '@/components/ui/textarea'
import { Button } from '@/components/ui/button'
import { Tooltip, TooltipContent, TooltipTrigger } from '@/components/ui/tooltip'
import { Plus, Trash2 } from 'lucide-vue-next'

const props = defineProps<{ config: FooterConfig }>()
const emit = defineEmits<{ update: [config: FooterConfig] }>()

const form = reactive<FooterConfig>(JSON.parse(JSON.stringify(props.config)))

watch(form, () => emit('update', JSON.parse(JSON.stringify(form))), { deep: true })

function addCustomLink() {
  form.custom_links.push({ name: '', url: '' })
}

function removeCustomLink(i: number) {
  form.custom_links.splice(i, 1)
}
</script>

<template>
  <div class="space-y-5">

    <!-- Logo 和标语 -->
    <fieldset class="rounded-xl border p-3.5 space-y-3">
      <legend class="text-xs font-semibold px-1.5 text-muted-foreground uppercase tracking-wider">Logo 和标语</legend>
      <div class="space-y-1.5">
        <Label for="footer-logo-url" class="text-xs text-muted-foreground">Logo URL</Label>
        <Input id="footer-logo-url" v-model="form.logo_url" placeholder="留空使用租户 Logo" />
      </div>
      <div class="space-y-1.5">
        <Label for="footer-slogan" class="text-xs text-muted-foreground">宣传标语</Label>
        <Textarea id="footer-slogan" v-model="form.slogan" placeholder="在此输入标语" :rows="2" />
      </div>
    </fieldset>

    <!-- 自定义链接 -->
    <fieldset class="rounded-xl border p-3.5 space-y-2">
      <legend class="text-xs font-semibold px-1.5 text-muted-foreground uppercase tracking-wider">自定义链接</legend>

      <div v-for="(link, i) in form.custom_links" :key="i" class="flex items-center gap-1.5">
        <Input :name="`footer-link-name-${i}`" v-model="link.name" placeholder="链接名称" class="flex-1" />
        <Input :name="`footer-link-url-${i}`" v-model="link.url" placeholder="URL（可选）" class="flex-1" />
        <Tooltip>
          <TooltipTrigger as-child>
            <Button variant="ghost" size="icon" class="rounded-full h-7 w-7 shrink-0 text-destructive/70 hover:text-destructive" @click="removeCustomLink(i)">
              <Trash2 class="h-3.5 w-3.5" />
            </Button>
          </TooltipTrigger>
          <TooltipContent>删除链接</TooltipContent>
        </Tooltip>
      </div>

      <p v-if="!form.custom_links.length" class="py-2 text-xs text-center text-muted-foreground">暂无自定义链接</p>

      <button
        class="w-full rounded-lg border border-dashed border-border py-2 text-xs text-muted-foreground hover:border-primary/50 hover:text-primary transition-colors"
        @click="addCustomLink"
      >
        <Plus class="h-3 w-3 inline mr-1" />
        添加链接
      </button>
    </fieldset>

    <!-- 版权信息 -->
    <fieldset class="rounded-xl border p-3.5 space-y-3">
      <legend class="text-xs font-semibold px-1.5 text-muted-foreground uppercase tracking-wider">版权信息</legend>
      <div class="grid gap-2 sm:grid-cols-2">
        <div class="space-y-1.5">
          <Label for="footer-copyright-owner" class="text-xs text-muted-foreground">版权所有者</Label>
          <Input id="footer-copyright-owner" v-model="form.copyright.owner" placeholder="如：Microswift Core™" />
        </div>
        <div class="space-y-1.5">
          <Label for="footer-copyright-year" class="text-xs text-muted-foreground">版权年份</Label>
          <Input id="footer-copyright-year" v-model="form.copyright.year" placeholder="如：2020-2025" />
        </div>
      </div>
      <div class="space-y-1.5">
        <Label for="footer-copyright-link" class="text-xs text-muted-foreground">版权所有者链接（可选）</Label>
        <Input id="footer-copyright-link" v-model="form.copyright.link" placeholder="https://..." />
      </div>
    </fieldset>

    <!-- 网站名称 -->
    <fieldset class="rounded-xl border p-3.5 space-y-3">
      <legend class="text-xs font-semibold px-1.5 text-muted-foreground uppercase tracking-wider">网站名称</legend>
      <div class="space-y-1.5">
        <Label for="footer-site-name" class="text-xs text-muted-foreground">网站名称</Label>
        <Input id="footer-site-name" v-model="form.site_name" placeholder="留空使用租户名称" />
      </div>
    </fieldset>

    <!-- 备案信息 -->
    <fieldset class="rounded-xl border p-3.5 space-y-3">
      <legend class="text-xs font-semibold px-1.5 text-muted-foreground uppercase tracking-wider">备案信息</legend>
      <div class="grid gap-2 sm:grid-cols-2">
        <div class="space-y-1.5">
          <Label for="footer-icp-number" class="text-xs text-muted-foreground">ICP 备案号</Label>
          <Input id="footer-icp-number" v-model="form.filing.icp_number" placeholder="如：陇ICP备20002844号-1" />
        </div>
        <div class="space-y-1.5">
          <Label for="footer-icp-link" class="text-xs text-muted-foreground">ICP 备案链接</Label>
          <Input id="footer-icp-link" v-model="form.filing.icp_link" placeholder="https://..." />
        </div>
      </div>
      <div class="grid gap-2 sm:grid-cols-2">
        <div class="space-y-1.5">
          <Label for="footer-police-number" class="text-xs text-muted-foreground">公安备案号</Label>
          <Input id="footer-police-number" v-model="form.filing.police_number" placeholder="如：甘公网安备62090202000540号" />
        </div>
        <div class="space-y-1.5">
          <Label for="footer-police-link" class="text-xs text-muted-foreground">公安备案链接</Label>
          <Input id="footer-police-link" v-model="form.filing.police_link" placeholder="https://..." />
        </div>
      </div>
    </fieldset>

  </div>
</template>
