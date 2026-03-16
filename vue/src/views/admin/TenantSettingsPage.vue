<!--
  TenantSettingsPage.vue - 租户设置页面
  职责：提供存储设置（本地/S3）和 AI 设置的管理界面
  对外接口：无 props，直接访问路由 /admin/settings
-->
<script setup lang="ts">
import { onMounted, computed, ref, watch } from 'vue'
import { useTenantSettings } from '@/lib/tenant-settings'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Switch } from '@/components/ui/switch'
import { Tooltip, TooltipContent, TooltipTrigger } from '@/components/ui/tooltip'
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select'
import {
  Accordion,
  AccordionContent,
  AccordionItem,
  AccordionTrigger,
} from '@/components/ui/accordion'
import { Separator } from '@/components/ui/separator'
import { Badge } from '@/components/ui/badge'
import { Loader2, HardDrive, Cloud, TestTube2, Check, Zap, Info, ChevronDown, Shield } from 'lucide-vue-next'

const {
  loading,
  saving,
  testing,
  storageForm,
  aiForm,
  accessForm,
  usage,
  loadSettings,
  loadUsage,
  saveStorage,
  testS3Connection,
  saveAI,
  saveAccess,
  formatBytes,
} = useTenantSettings()

const PROVIDER_CFG = {
  aws:        { endpoint: '留空（自动使用 AWS 官方域名）',                regionHint: 'us-east-1',    pathStyleNote: 'AWS S3 使用 Virtual-Hosted Style，无需开启' },
  cloudflare: { endpoint: 'https://<AccountID>.r2.cloudflarestorage.com', regionHint: 'auto',          pathStyleNote: 'R2 建议开启 Path Style' },
  aliyun:     { endpoint: 'https://oss-cn-hangzhou.aliyuncs.com',          regionHint: 'cn-hangzhou',   pathStyleNote: '阿里云 OSS 不支持 Path Style，请关闭' },
  tencent:    { endpoint: 'https://cos.ap-guangzhou.myqcloud.com',          regionHint: 'ap-guangzhou',  pathStyleNote: '腾讯云 COS 不支持 Path Style，请关闭' },
  minio:      { endpoint: 'http://minio.example.com:9000',                  regionHint: 'us-east-1',     pathStyleNote: 'MinIO 必须开启 Path Style' },
  custom:     { endpoint: 'https://your-s3-endpoint.com',                   regionHint: 'your-region',   pathStyleNote: '根据服务商要求设置' },
} as const

const providerCfg = computed(() =>
  PROVIDER_CFG[storageForm.provider as keyof typeof PROVIDER_CFG] ?? PROVIDER_CFG.custom
)

const PROVIDER_GUIDE = {
  aws: {
    title: 'AWS S3 配置指南',
    steps: [
      '在 AWS 控制台创建 S3 Bucket，记录 Bucket 名称和所在区域（Region）',
      '进入 IAM，创建用户并授予 s3:PutObject / GetObject / DeleteObject / ListBucket 权限，创建 Access Key',
      '在 Bucket 权限设置中关闭「阻止所有公有访问」，并添加允许公有读取的 Bucket Policy',
      'Endpoint 留空，Region 填写 Bucket 所在区域代码，如 us-east-1',
    ],
    note: '文件 URL 格式：https://bucket.s3.region.amazonaws.com/key',
  },
  cloudflare: {
    title: 'Cloudflare R2 配置指南',
    steps: [
      '在 Cloudflare 控制台进入 R2，创建存储桶',
      '在 R2 总览页找到账户 ID，Endpoint 格式为 https://<AccountID>.r2.cloudflarestorage.com',
      '进入「R2 → API 令牌」创建 Token，权限选择「对象读和写」，获得 Access Key ID 和 Secret Key',
      '建议在 R2 Bucket 设置中绑定自定义域名（CDN）作为文件访问地址',
    ],
    note: 'Region 填 auto；建议开启 Path Style',
  },
  aliyun: {
    title: '阿里云 OSS 配置指南',
    steps: [
      '在 OSS 控制台创建 Bucket，访问控制选择「公共读」，记录所在地域',
      '进入 RAM 访问控制，创建子账户并授予 AliyunOSSFullAccess（或自定义策略）',
      '在子账户下创建 AccessKey，获得 Access Key ID 和 Secret Key',
      'Endpoint 格式：https://oss-<地域>.aliyuncs.com，如 https://oss-cn-hangzhou.aliyuncs.com',
    ],
    note: '请关闭 Path Style；Endpoint 不含 Bucket 名称',
  },
  tencent: {
    title: '腾讯云 COS 配置指南',
    steps: [
      'Endpoint 必须填写不含 Bucket 的通用地址：https://cos.<地域>.myqcloud.com（如 https://cos.ap-chengdu.myqcloud.com）。注意：不要从控制台复制带 Bucket 名的 URL！',
      '在 COS 控制台创建存储桶，访问权限选择「公有读私有写」，Bucket 名称填入下方字段',
      '进入「访问管理 CAM → API 密钥」，新建密钥，SecretId 填入 Access Key ID，SecretKey 填入 Secret Access Key',
      '在 Bucket 设置 → 安全管理 → 跨域访问 CORS 中添加允许规则，以支持前端跨域上传',
    ],
    note: '请关闭 Path Style；自定义域名（CDN）可填写控制台绑定的自定义加速域名',
  },
  minio: {
    title: 'MinIO 配置指南',
    steps: [
      '确保 MinIO 服务正常运行并可被当前服务器网络访问',
      '通过 MinIO Console 或 mc 工具创建 Bucket，设置 Bucket 策略为 public（允许匿名读取）',
      'Endpoint 填写 MinIO 服务地址，格式：http://host:9000 或 https://minio.example.com',
      '使用 MinIO Root 账户或自建账户的 Access Key / Secret Key',
    ],
    note: '必须开启 Path Style；Region 填任意值，如 us-east-1',
  },
  custom: {
    title: '自定义 S3 兼容服务',
    steps: [
      '确认服务商支持 S3 兼容 API（兼容 AWS SDK v2 协议）',
      '从服务商控制台获取 Endpoint、Region、Access Key ID 和 Secret Access Key',
      '根据服务商文档决定是否开启 Path Style',
      '保存后使用「测试连接」按钮验证配置是否正确',
    ],
    note: '若文件 URL 访问异常，建议配置自定义域名（CDN）',
  },
} as const

const showGuide = ref(true)

const providerGuide = computed(() =>
  PROVIDER_GUIDE[storageForm.provider as keyof typeof PROVIDER_GUIDE] ?? PROVIDER_GUIDE.custom
)

watch(() => storageForm.provider, () => {
  showGuide.value = true
})

onMounted(() => {
  loadSettings()
  loadUsage()
})
</script>

<template>
  <div class="max-w-3xl mx-auto px-4 py-8">
    <div v-if="loading" class="flex justify-center py-16">
      <Loader2 class="size-6 animate-spin text-muted-foreground" />
    </div>

    <Accordion v-else type="multiple" :default-value="['storage', 'ai', 'access']">
      <!-- ══════════ 存储设置 ══════════ -->
      <AccordionItem value="storage">
        <AccordionTrigger class="hover:no-underline">
          <div class="flex items-center gap-3 text-left w-full">
            <HardDrive class="size-5 text-muted-foreground shrink-0" :stroke-width="1.5" />
            <div class="flex-1">
              <div class="text-lg font-semibold">存储设置</div>
              <div class="text-sm text-muted-foreground font-normal mt-1">
                配置文件上传的存储策略，本地存储上限 100MB
              </div>
            </div>
            <Badge variant="outline" class="font-mono text-xs shrink-0">
              {{ usage.used_mb.toFixed(1) }} / {{ usage.limit_mb.toFixed(0) }} MB
            </Badge>
          </div>
        </AccordionTrigger>

        <AccordionContent class="space-y-5 pt-4">
          <!-- 本地用量进度条 -->
          <div class="space-y-1.5">
            <div class="h-1.5 w-full rounded-full bg-muted overflow-hidden">
              <div
                class="h-full rounded-full transition-all"
                :class="usage.percent > 85 ? 'bg-destructive' : 'bg-primary'"
                :style="{ width: `${Math.min(usage.percent, 100)}%` }"
              />
            </div>
            <div class="flex justify-between text-xs text-muted-foreground">
              <span>本地已用 {{ formatBytes(usage.used_bytes) }}</span>
              <span>{{ usage.percent.toFixed(1) }}% / 上限 {{ usage.limit_mb.toFixed(0) }} MB</span>
            </div>
          </div>

          <!-- 上传目标 -->
          <div class="flex items-center justify-between">
            <div>
              <Label class="text-sm font-medium">默认上传目标</Label>
              <p class="text-xs text-muted-foreground mt-0.5">本地满时自动切换到云存储</p>
            </div>
            <Select v-model="storageForm.default_target">
              <SelectTrigger class="w-28 rounded-xl">
                <SelectValue />
              </SelectTrigger>
              <SelectContent class="rounded-2xl">
                <SelectItem value="local">本地</SelectItem>
                <SelectItem value="cloud">云存储</SelectItem>
              </SelectContent>
            </Select>
          </div>

          <div class="flex justify-center">
            <div class="w-96 h-px bg-border" />
          </div>

          <!-- 云存储开关 -->
          <div class="flex items-center justify-between">
            <div class="flex items-center gap-2">
              <Cloud class="size-4 text-muted-foreground" :stroke-width="1.5" />
              <div>
                <Label class="text-sm font-medium">启用 S3 兼容云存储（上传限制：50M）</Label>
                <p class="text-xs text-muted-foreground mt-0.5">支持 AWS S3、MinIO、R2 等</p>
              </div>
            </div>
            <Switch v-model="storageForm.enabled" />
          </div>

          <template v-if="storageForm.enabled">
            <div class="flex justify-center">
              <Separator class="w-1/2" />
            </div>

            <!-- 服务商 -->
            <div class="grid gap-1.5">
              <Label>服务商</Label>
              <Select v-model="storageForm.provider">
                <SelectTrigger class="rounded-xl">
                  <SelectValue />
                </SelectTrigger>
                <SelectContent class="rounded-2xl">
                  <SelectItem value="aws">AWS S3</SelectItem>
                  <SelectItem value="cloudflare">Cloudflare R2</SelectItem>
                  <SelectItem value="aliyun">阿里云 OSS</SelectItem>
                  <SelectItem value="tencent">腾讯云 COS</SelectItem>
                  <SelectItem value="minio">MinIO</SelectItem>
                  <SelectItem value="custom">自定义</SelectItem>
                </SelectContent>
              </Select>
            </div>

            <!-- 配置指南 -->
            <div class="rounded-2xl border bg-muted/30 overflow-hidden">
              <button
                class="w-full flex items-center justify-between px-4 py-3 text-sm hover:bg-muted/50 transition-colors"
                @click="showGuide = !showGuide"
              >
                <div class="flex items-center gap-2 font-medium">
                  <Info class="size-3.5 text-muted-foreground shrink-0" :stroke-width="1.5" />
                  {{ providerGuide.title }}
                </div>
                <ChevronDown
                  class="size-3.5 text-muted-foreground transition-transform shrink-0"
                  :class="showGuide ? 'rotate-180' : ''"
                  :stroke-width="1.5"
                />
              </button>
              <div v-if="showGuide" class="px-4 pb-4 space-y-3">
                <ol class="space-y-2">
                  <li
                    v-for="(step, i) in providerGuide.steps"
                    :key="i"
                    class="flex gap-2.5 text-xs text-muted-foreground"
                  >
                    <span class="shrink-0 size-4 rounded-full bg-primary/10 text-primary text-[10px] flex items-center justify-center font-semibold mt-0.5">{{ i + 1 }}</span>
                    <span>{{ step }}</span>
                  </li>
                </ol>
                <p class="text-xs text-muted-foreground border-t pt-2.5 italic">
                  {{ providerGuide.note }}
                </p>
              </div>
            </div>

            <!-- 端点 / 区域 -->
            <div class="grid grid-cols-2 gap-4">
              <div class="grid gap-1.5">
                <Label>自定义端点（Endpoint）</Label>
                <Input v-model="storageForm.endpoint" :placeholder="providerCfg.endpoint" class="rounded-xl" />
              </div>
              <div class="grid gap-1.5">
                <Label>区域（Region）</Label>
                <Input v-model="storageForm.region" :placeholder="providerCfg.regionHint" class="rounded-xl" />
              </div>
            </div>

            <!-- Bucket -->
            <div class="grid gap-1.5">
              <Label>Bucket 名称</Label>
              <Input v-model="storageForm.bucket" placeholder="my-docs-bucket" class="rounded-xl" />
            </div>

            <!-- Access Key -->
            <div class="grid grid-cols-2 gap-4">
              <div class="grid gap-1.5">
                <Label>Access Key ID</Label>
                <Input v-model="storageForm.access_key_id" class="rounded-xl font-mono text-sm" />
              </div>
              <div class="grid gap-1.5">
                <Label>Secret Access Key</Label>
                <Input
                  v-model="storageForm.secret_access_key"
                  type="password"
                  placeholder="留空保持不变"
                  class="rounded-xl font-mono text-sm"
                />
              </div>
            </div>

            <!-- 自定义域名 -->
            <div class="grid gap-1.5">
              <Label>自定义域名（CDN，可选）</Label>
              <Input v-model="storageForm.custom_domain" placeholder="https://cdn.example.com" class="rounded-xl" />
            </div>

            <!-- Path Style -->
            <div class="flex items-center justify-between">
              <div>
                <Label class="text-sm font-medium">强制 Path Style</Label>
                <p class="text-xs text-muted-foreground mt-0.5">{{ providerCfg.pathStyleNote }}</p>
              </div>
              <Switch v-model="storageForm.force_path_style" />
            </div>
          </template>

          <!-- 操作按钮 -->
          <div class="flex justify-end gap-2 pt-2">
            <Tooltip v-if="storageForm.enabled">
              <TooltipTrigger as-child>
                <Button
                  variant="outline"
                  size="icon"
                  class="rounded-full"
                  :disabled="testing"
                  @click="testS3Connection"
                >
                  <Loader2 v-if="testing" class="size-4 animate-spin" />
                  <TestTube2 v-else class="size-4" :stroke-width="1.5" />
                </Button>
              </TooltipTrigger>
              <TooltipContent>测试链接</TooltipContent>
            </Tooltip>
            <Tooltip>
              <TooltipTrigger as-child>
                <Button
                  size="icon"
                  class="rounded-full"
                  :disabled="saving"
                  @click="saveStorage"
                >
                  <Loader2 v-if="saving" class="size-4 animate-spin" />
                  <Check v-else class="size-4" :stroke-width="2" />
                </Button>
              </TooltipTrigger>
              <TooltipContent>保存存储设置</TooltipContent>
            </Tooltip>
          </div>
        </AccordionContent>
      </AccordionItem>

      <!-- ══════════ AI 设置 ══════════ -->
      <AccordionItem value="ai">
        <AccordionTrigger class="hover:no-underline">
          <div class="flex items-center gap-3 text-left w-full">
            <Zap class="size-5 text-muted-foreground shrink-0" :stroke-width="1.5" />
            <div class="flex-1">
              <div class="text-lg font-semibold flex items-center gap-2">
                AI 设置
                <Badge variant="secondary" class="text-xs">预留</Badge>
              </div>
              <div class="text-sm text-muted-foreground font-normal mt-1">
                配置 AI 服务端点，用于文档摘要、智能问答等功能（即将推出）
              </div>
            </div>
          </div>
        </AccordionTrigger>

        <AccordionContent class="space-y-5 pt-4">
          <!-- 启用开关 -->
          <div class="flex items-center justify-between">
            <Label class="text-sm font-medium">启用 AI 服务</Label>
            <Switch v-model="aiForm.enabled" />
          </div>

          <template v-if="aiForm.enabled">
            <div class="flex justify-center">
              <Separator class="w-1/2" />
            </div>

            <!-- Chat -->
            <div class="space-y-3">
              <Label class="text-xs uppercase tracking-wider text-muted-foreground">对话模型（Chat）</Label>
              <div class="grid grid-cols-2 gap-4">
                <div class="grid gap-1.5">
                  <Label class="text-xs">Base URL</Label>
                  <Input v-model="aiForm.chat!.base_url" placeholder="https://api.openai.com/v1" class="rounded-xl" />
                </div>
                <div class="grid gap-1.5">
                  <Label class="text-xs">Model ID</Label>
                  <Input v-model="aiForm.chat!.model_id" placeholder="gpt-4o" class="rounded-xl" />
                </div>
              </div>
            </div>
          </template>

          <div class="flex justify-end pt-2">
            <Tooltip>
              <TooltipTrigger as-child>
                <Button
                  size="icon"
                  class="rounded-full"
                  :disabled="saving"
                  @click="saveAI"
                >
                  <Loader2 v-if="saving" class="size-4 animate-spin" />
                  <Check v-else class="size-4" :stroke-width="2" />
                </Button>
              </TooltipTrigger>
              <TooltipContent>保存 AI 设置</TooltipContent>
            </Tooltip>
          </div>
        </AccordionContent>
      </AccordionItem>
      <!-- ══════════ 访问控制设置 ══════════ -->
      <AccordionItem value="access">
        <AccordionTrigger class="hover:no-underline">
          <div class="flex items-center gap-3 text-left w-full">
            <Shield class="size-5 text-muted-foreground shrink-0" :stroke-width="1.5" />
            <div class="flex-1">
              <div class="text-lg font-semibold">访问控制</div>
              <div class="text-sm text-muted-foreground font-normal mt-1">
                维护模式、画廊登录可见等站点级访问策略
              </div>
            </div>
          </div>
        </AccordionTrigger>

        <AccordionContent class="space-y-5 pt-4">
          <!-- 维护模式 -->
          <div class="flex items-center justify-between">
            <div>
              <Label class="text-sm font-medium">维护模式</Label>
              <p class="text-xs text-muted-foreground mt-0.5">开启后读者端整站仅登录用户可见，未登录用户看到维护页面</p>
            </div>
            <Switch v-model="accessForm.maintenance_mode" />
          </div>

          <div class="flex justify-center">
            <Separator class="w-1/2" />
          </div>

          <!-- 画廊登录可见 -->
          <div class="flex items-center justify-between">
            <div>
              <Label class="text-sm font-medium">画廊登录可见</Label>
              <p class="text-xs text-muted-foreground mt-0.5">开启后主题列表页需登录后才能访问</p>
            </div>
            <Switch v-model="accessForm.gallery_login_required" />
          </div>

          <div class="flex justify-end pt-2">
            <Tooltip>
              <TooltipTrigger as-child>
                <Button
                  size="icon"
                  class="rounded-full"
                  :disabled="saving"
                  @click="saveAccess"
                >
                  <Loader2 v-if="saving" class="size-4 animate-spin" />
                  <Check v-else class="size-4" :stroke-width="2" />
                </Button>
              </TooltipTrigger>
              <TooltipContent>保存访问控制设置</TooltipContent>
            </Tooltip>
          </div>
        </AccordionContent>
      </AccordionItem>
    </Accordion>
  </div>
</template>
