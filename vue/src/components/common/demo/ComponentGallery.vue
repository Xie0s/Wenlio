<!--
  ComponentGallery.vue - 组件库展示中心
  职责：集中展示 src/components/ui 目录下的所有原子组件及其用法示例
  对外接口：无
-->
<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import * as Icons from 'lucide-vue-next'
import { cn } from '@/utils'
import type { ButtonVariants } from '@/components/ui/button'

// UI 组件导入
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Textarea } from '@/components/ui/textarea'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { 
  Card as HeroCard, 
  CardHeader as HeroCardHeader, 
  CardBody as HeroCardBody, 
  CardFooter as HeroCardFooter 
} from '@/components/ui/a-heroui/card'
import { Badge } from '@/components/ui/badge'
import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar'
import { Tooltip, TooltipContent, TooltipProvider, TooltipTrigger } from '@/components/ui/tooltip'
import { Accordion, AccordionContent, AccordionItem, AccordionTrigger } from '@/components/ui/accordion'
import { Alert, AlertDescription, AlertTitle } from '@/components/ui/alert'
import { AspectRatio } from '@/components/ui/aspect-ratio'
import { 
  AlertDialog, AlertDialogAction, AlertDialogCancel, AlertDialogContent, 
  AlertDialogDescription, AlertDialogFooter, AlertDialogHeader, AlertDialogTitle, AlertDialogTrigger 
} from '@/components/ui/alert-dialog'
import { Breadcrumb, BreadcrumbItem, BreadcrumbLink, BreadcrumbList, BreadcrumbPage, BreadcrumbSeparator } from '@/components/ui/breadcrumb'
import { ButtonGroup, ButtonGroupSeparator, ButtonGroupText } from '@/components/ui/button-group'
import { Carousel, CarouselContent, CarouselItem, CarouselNext, CarouselPrevious } from '@/components/ui/carousel'
import { CalendarDatePickerDemo, CalendarRangeDemo } from '@/components/common/demo/Calendar'
import { DrawerBasicDemo, DrawerDirectionDemo, DrawerResponsiveDemo } from '@/components/common/demo/Drawer'
import { ChartBarDemo, ChartAreaDemo, ChartDonutDemo, ChartLineDemo, ChartTooltipDefaultDemo } from '@/components/common/demo/charts'
import { Checkbox } from '@/components/ui/checkbox'
import { Collapsible, CollapsibleContent, CollapsibleTrigger } from '@/components/ui/collapsible'
import {
  Combobox,
  ComboboxAnchor,
  ComboboxEmpty,
  ComboboxGroup,
  ComboboxInput,
  ComboboxItem,
  ComboboxItemIndicator,
  ComboboxList,
  ComboboxTrigger,
  ComboboxViewport,
} from '@/components/ui/combobox'
import {
  Command,
  CommandEmpty,
  CommandGroup,
  CommandInput,
  CommandItem,
  CommandList,
  CommandSeparator,
} from '@/components/ui/command'
import {
  ContextMenu,
  ContextMenuCheckboxItem,
  ContextMenuContent,
  ContextMenuItem,
  ContextMenuSeparator,
  ContextMenuTrigger,
} from '@/components/ui/context-menu'
import { 
  Dialog, DialogContent, DialogDescription, DialogFooter, DialogHeader, 
  DialogTitle, DialogTrigger, DialogClose 
} from '@/components/ui/dialog'
import { 
  DropdownMenu, DropdownMenuContent, DropdownMenuItem, DropdownMenuLabel, 
  DropdownMenuSeparator, DropdownMenuTrigger 
} from '@/components/ui/dropdown-menu'
import { Empty, EmptyContent, EmptyDescription, EmptyHeader, EmptyMedia, EmptyTitle } from '@/components/ui/empty'
import { Field, FieldContent, FieldDescription, FieldGroup, FieldLabel, FieldSet, FieldTitle } from '@/components/ui/field'
import { Form, FormControl, FormDescription, FormField, FormItem, FormLabel, FormMessage } from '@/components/ui/form'
import { HoverCard, HoverCardContent, HoverCardTrigger } from '@/components/ui/hover-card'
import { InputGroup, InputGroupAddon, InputGroupButton, InputGroupInput } from '@/components/ui/input-group'
import { InputOTP, InputOTPGroup, InputOTPSeparator, InputOTPSlot } from '@/components/ui/input-otp'
import {
  Item,
  ItemActions,
  ItemContent,
  ItemDescription,
  ItemGroup,
  ItemHeader,
  ItemMedia,
  ItemSeparator,
  ItemTitle,
} from '@/components/ui/item'
import { Kbd, KbdGroup } from '@/components/ui/kbd'
import { Label } from '@/components/ui/label'
import {
  Menubar,
  MenubarCheckboxItem,
  MenubarContent,
  MenubarItem,
  MenubarMenu,
  MenubarSeparator,
  MenubarShortcut,
  MenubarTrigger,
} from '@/components/ui/menubar'
import { NativeSelect, NativeSelectOptGroup, NativeSelectOption } from '@/components/ui/native-select'
import { Popover, PopoverContent, PopoverTrigger } from '@/components/ui/popover'
import { Select, SelectContent, SelectGroup, SelectItem, SelectLabel, SelectTrigger, SelectValue } from '@/components/ui/select'
import { Separator } from '@/components/ui/separator'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'
import { Table, TableBody, TableCaption, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table'
import { NoiseBackground, TooltipCard } from '@/components/ui/a-aceternity'
import ExpandableCardDemo from '@/components/common/demo/expandable-card/ExpandableCardDemo.vue'
import GlassmorphismDemo from '@/components/common/demo/GlassmorphismDemo.vue'
import TextHoverEffectDemo from '@/components/common/demo/text-hover-effect/TextHoverEffectDemo.vue'
import { AnimatedThemeToggler } from '@/components/ui/a-magicui'
import ThemeToggle from '@/components/common/ThemeToggle.vue'
import { toast } from 'vue-sonner'

// 定义组件列表
const componentsList = [
  { id: 'accordion', name: 'Accordion - 手风琴', icon: Icons.ListCollapse },
  { id: 'alert', name: 'Alert - 警告提示', icon: Icons.AlertCircle },
  { id: 'alert-dialog', name: 'AlertDialog - 确认弹窗', icon: Icons.MessageSquareWarning },
  { id: 'aspect-ratio', name: 'AspectRatio - 比例容器', icon: Icons.Square },
  { id: 'avatar', name: 'Avatar - 头像', icon: Icons.UserCircle },
  { id: 'badge', name: 'Badge - 徽章', icon: Icons.Tag },
  { id: 'breadcrumb', name: 'Breadcrumb - 面包屑', icon: Icons.ChevronRightSquare },
  { id: 'button', name: 'Button - 按钮', icon: Icons.MousePointer2 },
  { id: 'button-group', name: 'ButtonGroup - 按钮组', icon: Icons.Columns },
  { id: 'calendar-date-picker', name: 'Calendar - 日期选择器', icon: Icons.Calendar },
  { id: 'calendar-range', name: 'Range Calendar - 日期范围', icon: Icons.Calendar },
  { id: 'card', name: 'Card - 卡片', icon: Icons.CreditCard },
  { id: 'heroui-card', name: 'HeroUI Card - 高级卡片', icon: Icons.LayoutTemplate },
  { id: 'carousel', name: 'Carousel - 轮播图', icon: Icons.GalleryHorizontal },
  { id: 'chart', name: 'Chart - 图表', icon: Icons.BarChart3 },
  { id: 'checkbox', name: 'Checkbox - 复选框', icon: Icons.CheckSquare },
  { id: 'collapsible', name: 'Collapsible - 可折叠', icon: Icons.ChevronsUpDown },
  { id: 'combobox', name: 'Combobox - 组合框', icon: Icons.Search },
  { id: 'command', name: 'Command - 指令面板', icon: Icons.Terminal },
  { id: 'context-menu', name: 'ContextMenu - 右键菜单', icon: Icons.SquareMousePointer },
  { id: 'dialog', name: 'Dialog - 对话框', icon: Icons.Maximize2 },
  { id: 'drawer', name: 'Drawer - 抽屉', icon: Icons.PanelBottom },
  { id: 'dropdown-menu', name: 'Dropdown - 下拉菜单', icon: Icons.Menu },
  { id: 'empty', name: 'Empty - 空状态', icon: Icons.FileWarning },
  { id: 'field', name: 'Field - 字段', icon: Icons.CaseSensitive },
  { id: 'form', name: 'Form - 表单', icon: Icons.Layout },
  { id: 'glassmorphism', name: 'Glassmorphism - 玻璃质感', icon: Icons.GlassWater },
  { id: 'hover-card', name: 'HoverCard - 悬浮卡片', icon: Icons.IdCard },
  { id: 'input', name: 'Input - 输入框', icon: Icons.Type },
  { id: 'input-group', name: 'InputGroup - 输入框组', icon: Icons.TextSelect },
  { id: 'input-otp', name: 'InputOTP - 验证码', icon: Icons.Hash },
  { id: 'item', name: 'Item - 条目', icon: Icons.Layers },
  { id: 'kbd', name: 'Kbd - 快捷键', icon: Icons.Keyboard },
  { id: 'label', name: 'Label - 标签', icon: Icons.CaseUpper },
  { id: 'menubar', name: 'Menubar - 菜单栏', icon: Icons.MoreHorizontal },
  { id: 'native-select', name: 'NativeSelect - 原生选择', icon: Icons.ChevronDown },
  { id: 'popover', name: 'Popover - 气泡卡片', icon: Icons.MessageSquare },
  { id: 'select', name: 'Select - 选择器', icon: Icons.ListFilter },
  { id: 'separator', name: 'Separator - 分割线', icon: Icons.Minus },
  { id: 'sonner', name: 'Sonner - 通知', icon: Icons.Bell },
  { id: 'table', name: 'Table - 表格', icon: Icons.Table2 },
  { id: 'tabs', name: 'Tabs - 标签页', icon: Icons.Columns },
  { id: 'textarea', name: 'Textarea - 文本域', icon: Icons.AlignLeft },
  { id: 'tooltip', name: 'Tooltip - 文字提示', icon: Icons.Info },
  { id: 'tooltip-advanced', name: 'TooltipCard - 高级 Tooltip', icon: Icons.MousePointerClick },
  { id: 'noise-background', name: 'NoiseBackground - 噪声背景', icon: Icons.Eraser },
  { id: 'animated-theme-toggler', name: 'ThemeToggler - 主题切换器', icon: Icons.SunMoon },
  { id: 'expandable-card', name: 'ExpandableCard - 高级展开卡', icon: Icons.LayoutTemplate },
  { id: 'text-hover-effect', name: 'TextHoverEffect - 高级文本悬停效果', icon: Icons.Sparkles },
].sort((a, b) => a.name.localeCompare(b.name))

const SELECTED_COMPONENT_STORAGE_KEY = 'demo-gallery:selected-component-id'
const getDefaultSelectedId = () => componentsList[0]?.id || 'accordion'

function getInitialSelectedId() {
  const defaultId = getDefaultSelectedId()
  try {
    const cachedId = window.localStorage.getItem(SELECTED_COMPONENT_STORAGE_KEY)
    if (!cachedId) {
      return defaultId
    }

    const exists = componentsList.some(item => item.id === cachedId)
    return exists ? cachedId : defaultId
  } catch {
    return defaultId
  }
}

const selectedId = ref(getInitialSelectedId())
const searchQuery = ref('')

const filteredComponents = computed(() => {
  return componentsList.filter(c => 
    c.name.toLowerCase().includes(searchQuery.value.toLowerCase()) ||
    c.id.toLowerCase().includes(searchQuery.value.toLowerCase())
  )
})

watch(filteredComponents, (nextList) => {
  if (nextList.length === 0) {
    selectedId.value = ''
    return
  }

  const stillExists = nextList.some(item => item.id === selectedId.value)
  const firstItem = nextList[0]
  if (!stillExists && firstItem) {
    selectedId.value = firstItem.id
  }
}, { immediate: true })

watch(selectedId, (nextId) => {
  try {
    if (!nextId) {
      window.localStorage.removeItem(SELECTED_COMPONENT_STORAGE_KEY)
      return
    }
    window.localStorage.setItem(SELECTED_COMPONENT_STORAGE_KEY, nextId)
  } catch {
    // 忽略存储不可用场景（例如隐私模式限制）
  }
})

const selectedComponent = computed(() => {
  if (!selectedId.value) {
    return undefined
  }

  const component = componentsList.find(c => c.id === selectedId.value)
  return component
})

const styleSpecNotes: Record<string, { title: string, items: string[] }> = {
  'heroui-card': {
    title: '规范对齐（HeroUI Card）',
    items: [
      'HeroUI Card 提供了更丰富的交互状态：isHoverable, isPressable。',
      '支持 radius (none/sm/md/lg/2xl/full) 定制。',
      'CardFooter 支持 isBlurred 效果，常用于图片覆盖场景。',
    ],
  },
  card: {
    title: '规范对齐（Card）',
    items: [
      'Card 属于容器级组件，根容器统一使用 rounded-2xl。',
      '卡片用于内容分组，不承担模态交互职责。',
    ],
  },
  'dropdown-menu': {
    title: '规范对齐（DropdownMenu）',
    items: [
      'DropdownMenuContent 根容器使用 rounded-2xl。',
      '轻量交互场景建议保持 modal: false，避免页面滚动抖动。',
    ],
  },
  select: {
    title: '规范对齐（Select）',
    items: [
      'SelectContent 属于容器级浮层，使用 rounded-2xl。',
      '轻量下拉场景应避免额外滚动锁定，保持背景稳定。',
    ],
  },
  'context-menu': {
    title: '规范对齐（ContextMenu）',
    items: [
      'ContextMenuContent 根容器使用 rounded-2xl。',
      '右键菜单为轻量交互，保持 modal: false 以减少布局偏移。',
    ],
  },
  drawer: {
    title: '规范对齐（Drawer）',
    items: [
      'DrawerContent 属于容器级组件，顶部/底部圆角使用 rounded-2xl。',
      '移动端长内容与操作聚合优先考虑 Drawer。',
    ],
  },
  'hover-card': {
    title: '规范对齐（HoverCard）',
    items: [
      'HoverCardContent 根容器使用 rounded-2xl。',
      '适合轻量说明，不建议承载复杂表单。',
    ],
  },
  input: {
    title: '规范对齐（Input）',
    items: [
      'Input 属于基础输入控件，统一 rounded-2xl。',
      '作为字段控件时建议配合 Label / Form 体系。',
    ],
  },
  textarea: {
    title: '规范对齐（Textarea）',
    items: [
      'Textarea 属于基础输入控件，统一 rounded-2xl。',
      '多行文本输入建议明确最小高度与状态提示。',
    ],
  },
  dialog: {
    title: '规范对齐（Dialog）',
    items: [
      'DialogContent 属于重型模态浮层，使用 rounded-3xl。',
      '该演示场景关闭右上角按钮，仅保留底部“取消/确认”操作。',
    ],
  },
  'alert-dialog': {
    title: '规范对齐（AlertDialog）',
    items: [
      'AlertDialogContent 属于重型确认弹层，使用 rounded-3xl。',
      'Action/Cancel 使用自定义根组件默认图标化圆形按钮规范。',
    ],
  },
  popover: {
    title: '规范对齐（Popover）',
    items: [
      'PopoverContent 在项目中定义为 rounded-3xl 的信息浮层。',
      '适合短配置与补充说明，不替代完整 Dialog。',
    ],
  },
  'calendar-date-picker': {
    title: '规范对齐（Calendar 日期选择器）',
    items: [
      '触发按钮遵循圆形图标 Button + Tooltip 规范，文本按钮模式保持 rounded-2xl。',
      'PopoverContent 根容器沿用 rounded-2xl，避免过度阴影影响视线。',
      '选中日期后主动关闭浮层，提升效率。',
    ],
  },
  'calendar-range': {
    title: '规范对齐（Range Calendar 日期范围）',
    items: [
      '1440px 以上桌面使用双月并排，移动端单月展示，使用 useMediaQuery 控制。',
      'RangeCalendarCellTrigger 保持 8px 圆角与主题色 data-selected 状态，不使用 !important。',
      '页眉左右箭头使用圆形按钮搭配 lucide 图标，hover 态统一。',
    ],
  },
  glassmorphism: {
    title: '规范对齐（Glassmorphism 玻璃质感）',
    items: [
      '`.glass` 是唯一的玻璃质感工具类，导航栏、卡片、弹窗、标签等所有场景统一使用这一个 class。',
      '效果依赖 `backdrop-filter: blur + saturate`，元素必须层叠在图片、渐变或纯色背景之上才能呈现玻璃感。',
      '无需额外变体 — 背景本身的密度与色彩丰富度会自然地控制玻璃遮蔽感的强弱。',
    ],
  },
}

const currentStyleSpecNote = computed(() => {
  return styleSpecNotes[selectedId.value]
})

const demoButtonVariants: NonNullable<ButtonVariants['variant']>[] = ['default', 'secondary', 'outline', 'ghost', 'destructive']

const comboboxValue = ref('vue')
const commandValue = ref('')
const checkboxValue = ref(true)
const collapsibleOpen = ref(false)
const contextMenuPinned = ref(true)
const selectValue = ref('vue')
const popoverNote = ref('需要跟进组件演示文档')
const otpValue = ref('')
const nativeSelectValue = ref('vue')
const textareaValue = ref('这是一个 Textarea 演示内容，可用于多行输入。')
const menubarStatusBar = ref(true)

const carouselItems = ['第一屏内容', '第二屏内容', '第三屏内容']


function handleFormSubmit() {
  toast.success('表单提交成功')
}

function showToast() {
  toast.success('这是一条测试通知')
}

function clearSearch() {
  searchQuery.value = ''
}
</script>

<template>
  <div class="gallery-root flex h-screen overflow-hidden">
    <!-- 左侧侧边栏 -->
    <aside class="relative flex w-72 shrink-0 flex-col overflow-hidden border-r">
      <div class="sidebar-gradient absolute inset-0 pointer-events-none" />

      <!-- 侧边栏头部 -->
      <div class="relative z-10 space-y-3 border-b border-border/60 px-4 py-4">
        <div class="flex items-center gap-2.5">
          <div class="flex size-8 shrink-0 items-center justify-center rounded-xl bg-primary shadow-sm">
            <Icons.Layers class="size-4 text-primary-foreground" />
          </div>
          <div class="min-w-0 flex-1">
            <p class="text-sm font-semibold leading-none tracking-tight">组件库</p>
            <p class="mt-0.5 text-[11px] text-muted-foreground">{{ componentsList.length }} 个组件</p>
          </div>
          <ThemeToggle class="shrink-0" />
        </div>

        <div class="flex items-center gap-2">
          <div class="relative flex-1">
            <Icons.Search class="absolute left-2.5 top-1/2 size-3.5 -translate-y-1/2 text-muted-foreground" />
            <Input
              v-model="searchQuery"
              placeholder="搜索组件..."
              class="h-8 rounded-xl bg-background/60 pl-8 text-sm"
            />
          </div>
          <button
            v-if="searchQuery"
            type="button"
            class="inline-flex size-8 shrink-0 items-center justify-center rounded-full border border-rose-200 bg-rose-50/70 text-rose-400 transition-colors hover:bg-rose-100 hover:text-rose-500 dark:border-rose-900/60 dark:bg-rose-950/30 dark:text-rose-300 dark:hover:bg-rose-950/50"
            aria-label="清空搜索"
            @click="clearSearch"
          >
            <Icons.X class="size-3.5" />
          </button>
        </div>
      </div>

      <!-- 组件列表 -->
      <nav class="relative z-10 flex-1 overflow-y-auto px-2 py-2 custom-scrollbar">
        <button
          v-for="item in filteredComponents"
          :key="item.id"
          type="button"
          @click="selectedId = item.id"
          :class="cn(
            'w-full flex items-center gap-2.5 px-3 py-2 rounded-xl transition-all duration-150 text-left',
            selectedId === item.id
              ? 'bg-primary/10 text-primary'
              : 'text-muted-foreground hover:bg-accent hover:text-foreground'
          )"
        >
          <component :is="item.icon" class="size-4 shrink-0" />
          <span class="truncate text-sm">{{ item.name }}</span>
          <Icons.ChevronRight
            v-if="selectedId === item.id"
            class="ml-auto size-3.5 shrink-0 text-primary"
          />
        </button>

        <div v-if="filteredComponents.length === 0" class="px-3 py-8 text-center text-sm text-muted-foreground">
          未找到相关组件
        </div>
      </nav>
    </aside>

    <!-- 右侧内容区域 -->
    <main class="main-area min-w-0 flex-1 overflow-x-hidden overflow-y-auto overscroll-contain p-8 custom-scrollbar">
      <div class="mx-auto max-w-4xl space-y-8">

        <!-- 头部 -->
        <header class="flex items-start justify-between gap-4">
          <div class="flex items-center gap-4">
            <div
              v-if="selectedComponent"
              class="glass flex size-11 shrink-0 items-center justify-center rounded-2xl"
            >
              <component :is="selectedComponent.icon" class="size-5" />
            </div>
            <div>
              <h1 v-if="selectedComponent" class="text-2xl font-bold leading-none tracking-tight">
                {{ (selectedComponent.name || '').split(' - ')[0] }}
              </h1>
              <p v-if="selectedComponent" class="mt-1.5 text-sm text-muted-foreground">
                展示 {{ selectedId }} 组件的各种变体和使用场景
              </p>
            </div>
          </div>

          <div class="flex shrink-0 gap-1.5">
            <TooltipProvider>
              <Tooltip>
                <TooltipTrigger as-child>
                  <Button variant="outline" size="icon" class="size-8 rounded-full" @click="showToast">
                    <Icons.Bell class="size-3.5" />
                  </Button>
                </TooltipTrigger>
                <TooltipContent>预览通知</TooltipContent>
              </Tooltip>
            </TooltipProvider>

            <TooltipProvider>
              <Tooltip>
                <TooltipTrigger as-child>
                  <Button variant="outline" size="icon" class="size-8 rounded-full">
                    <Icons.Code class="size-3.5" />
                  </Button>
                </TooltipTrigger>
                <TooltipContent>查看代码</TooltipContent>
              </Tooltip>
            </TooltipProvider>
          </div>
        </header>

        <Separator class="bg-border/50" />

        <div
          v-if="currentStyleSpecNote"
          class="glass rounded-2xl p-4"
        >
          <div class="mb-2.5 flex items-center gap-2">
            <Icons.BookOpen class="size-4 text-primary" />
            <h3 class="text-sm font-semibold">{{ currentStyleSpecNote.title }}</h3>
          </div>
          <ul class="list-disc space-y-1.5 pl-5 text-sm text-muted-foreground">
            <li v-for="(item, index) in currentStyleSpecNote.items" :key="index">
              {{ item }}
            </li>
          </ul>
        </div>

        <!-- 组件内容展示（直接展示，不再用 Card 作为统一包裹） -->
        <div class="grid gap-6">
          <div v-if="selectedId === 'accordion'" class="space-y-4">
            <Accordion type="single" collapsible class="w-full rounded-2xl border bg-card/50 px-4">
              <AccordionItem value="item-1">
                <AccordionTrigger>是否支持自定义样式？</AccordionTrigger>
                <AccordionContent>是的，你可以通过 Tailwind CSS 快速定制样式。</AccordionContent>
              </AccordionItem>
              <AccordionItem value="item-2">
                <AccordionTrigger>是否支持动画？</AccordionTrigger>
                <AccordionContent>默认内置了平滑过渡动画。</AccordionContent>
              </AccordionItem>
            </Accordion>
          </div>

          <div v-else-if="selectedId === 'alert'" class="space-y-4">
            <Alert class="rounded-2xl">
              <Icons.Terminal class="size-4" />
              <AlertTitle>提示</AlertTitle>
              <AlertDescription>你可以通过命令行快速安装更多组件。</AlertDescription>
            </Alert>
            <Alert variant="destructive" class="rounded-2xl">
              <Icons.AlertCircle class="size-4" />
              <AlertTitle>错误</AlertTitle>
              <AlertDescription>会话已过期，请重新登录。</AlertDescription>
            </Alert>
          </div>

          <div v-else-if="selectedId === 'alert-dialog'" class="space-y-4">
            <div class="rounded-2xl border bg-card/50 p-4 space-y-3">
              <h3 class="text-sm font-semibold">使用说明（AlertDialog）</h3>
              <p class="text-sm text-muted-foreground">
                AlertDialog 用于“高风险且不可逆”的确认操作（例如删除、清空、下线），必须显式由用户二次确认。
              </p>
              <ul class="list-disc pl-5 space-y-2 text-sm text-muted-foreground">
                <li>
                  <span class="font-medium text-foreground">推荐结构：</span>
                  AlertDialog → AlertDialogTrigger → AlertDialogContent → Header / Footer。
                </li>
                <li>
                  <span class="font-medium text-foreground">触发按钮（项目规范）：</span>
                  使用我们自定义 Button，优先 <code>size="icon"</code> + <code>rounded-full</code>，并配合 Tooltip 给出语义（如“删除项目”）。
                </li>
                <li>
                  <span class="font-medium text-foreground">按钮规范：</span>
                  取消与确认按钮优先使用根组件 <code>AlertDialogCancel</code> / <code>AlertDialogAction</code>，以继承统一图标与圆形风格。
                </li>
                <li>
                  <span class="font-medium text-foreground">我们自定义的根组件行为：</span>
                  <code>AlertDialogCancel</code> 默认是 <code>outline + 圆形 + X 图标</code>；
                  <code>AlertDialogAction</code> 默认是 <code>primary + 圆形 + Check 图标</code>。
                  演示区直接使用这两个组件即可，不建议在业务里重复覆盖其核心样式。
                </li>
                <li>
                  <span class="font-medium text-foreground">文案建议：</span>
                  标题直接描述动作，描述文本明确“后果 + 不可撤销性 + 影响范围”。
                </li>
                <li>
                  <span class="font-medium text-foreground">交互建议：</span>
                  Trigger 建议使用 Tooltip 提示动作语义，降低误触成本。
                </li>
              </ul>
            </div>

            <AlertDialog>
              <TooltipProvider>
                <Tooltip>
                  <TooltipTrigger as-child>
                    <AlertDialogTrigger as-child>
                      <Button variant="destructive" size="icon" class="rounded-full">
                        <Icons.Trash2 class="size-4" />
                      </Button>
                    </AlertDialogTrigger>
                  </TooltipTrigger>
                  <TooltipContent>删除项目</TooltipContent>
                </Tooltip>
              </TooltipProvider>
              <AlertDialogContent>
                <AlertDialogHeader>
                  <AlertDialogTitle>确认删除该项目吗？</AlertDialogTitle>
                  <AlertDialogDescription>此操作不可撤销，请谨慎确认。</AlertDialogDescription>
                </AlertDialogHeader>
                <AlertDialogFooter>
                  <AlertDialogCancel />
                  <AlertDialogAction />
                </AlertDialogFooter>
              </AlertDialogContent>
            </AlertDialog>
          </div>

          <div v-else-if="selectedId === 'aspect-ratio'" class="max-w-md">
            <AspectRatio :ratio="16 / 9" class="overflow-hidden rounded-2xl border bg-muted">
              <img src="https://images.unsplash.com/photo-1523206489230-c012c64b2b48?w=1200" alt="aspect ratio demo" class="h-full w-full object-cover" />
            </AspectRatio>
          </div>

          <div v-else-if="selectedId === 'avatar'" class="flex items-center gap-4">
            <Avatar size="lg" class="border-2 border-primary/20">
              <AvatarImage src="https://github.com/shadcn.png" alt="@shadcn" />
              <AvatarFallback>CN</AvatarFallback>
            </Avatar>
            <Avatar size="default"><AvatarFallback>JD</AvatarFallback></Avatar>
            <Avatar size="sm"><AvatarFallback><Icons.User class="size-4" /></AvatarFallback></Avatar>
          </div>

          <div v-else-if="selectedId === 'badge'" class="flex flex-wrap gap-2">
            <Badge>Default</Badge>
            <Badge variant="secondary">Secondary</Badge>
            <Badge variant="outline">Outline</Badge>
            <Badge variant="destructive">Destructive</Badge>
          </div>

          <div v-else-if="selectedId === 'breadcrumb'">
            <Breadcrumb>
              <BreadcrumbList>
                <BreadcrumbItem><BreadcrumbLink href="/">首页</BreadcrumbLink></BreadcrumbItem>
                <BreadcrumbSeparator />
                <BreadcrumbItem><BreadcrumbLink href="/components">组件库</BreadcrumbLink></BreadcrumbItem>
                <BreadcrumbSeparator />
                <BreadcrumbItem><BreadcrumbPage>当前页</BreadcrumbPage></BreadcrumbItem>
              </BreadcrumbList>
            </Breadcrumb>
          </div>

          <div v-else-if="selectedId === 'button'" class="space-y-6">
            <div class="flex flex-wrap gap-4">
              <TooltipProvider>
                <Tooltip v-for="(variant, i) in demoButtonVariants" :key="i">
                  <TooltipTrigger as-child>
                    <Button :variant="variant" size="icon" class="rounded-full">
                      <Icons.Zap v-if="variant === 'default'" class="size-4" />
                      <Icons.Settings v-else-if="variant === 'secondary'" class="size-4" />
                      <Icons.Share2 v-else-if="variant === 'outline'" class="size-4" />
                      <Icons.Ghost v-else-if="variant === 'ghost'" class="size-4" />
                      <Icons.Trash2 v-else class="size-4" />
                    </Button>
                  </TooltipTrigger>
                  <TooltipContent>{{ variant }} 按钮</TooltipContent>
                </Tooltip>
              </TooltipProvider>
            </div>
          </div>

          <div v-else-if="selectedId === 'button-group'" class="space-y-5">
            <div class="grid gap-4 md:grid-cols-2">
              <div class="space-y-2 rounded-2xl border bg-card/40 p-4">
                <p class="text-sm font-medium">场景 1：基础图标分组</p>
                <ButtonGroup>
                  <Button size="icon" variant="outline" aria-label="加粗"><Icons.Bold class="size-4" /></Button>
                  <Button size="icon" variant="outline" aria-label="斜体"><Icons.Italic class="size-4" /></Button>
                  <ButtonGroupSeparator />
                  <Button size="icon" variant="outline" aria-label="下划线"><Icons.Underline class="size-4" /></Button>
                </ButtonGroup>
              </div>

              <div class="space-y-2 rounded-2xl border bg-card/40 p-4">
                <p class="text-sm font-medium">场景 2：编辑工具组</p>
                <ButtonGroup>
                  <Button size="icon" variant="outline" aria-label="左对齐"><Icons.AlignLeft class="size-4" /></Button>
                  <Button size="icon" variant="outline" aria-label="居中对齐"><Icons.AlignCenter class="size-4" /></Button>
                  <Button size="icon" variant="outline" aria-label="右对齐"><Icons.AlignRight class="size-4" /></Button>
                </ButtonGroup>
              </div>

              <div class="space-y-2 rounded-2xl border bg-card/40 p-4">
                <p class="text-sm font-medium">场景 3：纵向操作组</p>
                <ButtonGroup orientation="vertical">
                  <Button size="icon" variant="outline" aria-label="上移"><Icons.ChevronUp class="size-4" /></Button>
                  <Button size="icon" variant="outline" aria-label="下移"><Icons.ChevronDown class="size-4" /></Button>
                </ButtonGroup>
              </div>

              <div class="space-y-2 rounded-2xl border bg-card/40 p-4">
                <p class="text-sm font-medium">场景 4：状态说明条</p>
                <ButtonGroupText><Icons.Wand2 class="size-4" /> 快速样式组</ButtonGroupText>
                <ButtonGroupText><Icons.ShieldCheck class="size-4" /> 权限已生效</ButtonGroupText>
              </div>
            </div>
          </div>

          <div v-else-if="selectedId === 'calendar-date-picker'" class="max-w-2xl">
            <CalendarDatePickerDemo />
          </div>

          <div v-else-if="selectedId === 'calendar-range'" class="w-fit overflow-hidden">
            <CalendarRangeDemo />
          </div>

          <div v-else-if="selectedId === 'heroui-card'" class="grid gap-6">
            <div class="grid gap-6 md:grid-cols-2">
              <!-- 基础展示 -->
              <HeroCard class="max-w-[400px]">
                <HeroCardHeader class="flex gap-3">
                  <Avatar class="rounded-2xl">
                    <AvatarImage src="https://heroui.com/avatars/avatar-1.png" />
                    <AvatarFallback>HU</AvatarFallback>
                  </Avatar>
                  <div class="flex flex-col">
                    <p class="text-md font-bold">HeroUI</p>
                    <p class="text-small text-muted-foreground">heroui.com</p>
                  </div>
                </HeroCardHeader>
                <HeroCardBody>
                  <p>Make beautiful websites regardless of your design experience.</p>
                </HeroCardBody>
                <HeroCardFooter>
                  <a href="https://github.com/heroui-inc/heroui" target="_blank" class="text-primary text-sm hover:underline">
                    Visit source code on GitHub.
                  </a>
                </HeroCardFooter>
              </HeroCard>

              <!-- 可交互卡片 -->
              <HeroCard is-pressable is-hoverable class="max-w-[400px]" @click="showToast">
                <HeroCardBody class="p-0">
                  <AspectRatio :ratio="16 / 9">
                    <img src="https://heroui.com/images/hero-card-complete.jpeg" class="w-full h-full object-cover" alt="Card background" />
                  </AspectRatio>
                </HeroCardBody>
                <HeroCardFooter class="justify-between">
                  <b class="text-sm">Frontend Framework</b>
                  <p class="text-muted-foreground text-xs">Vue 3</p>
                </HeroCardFooter>
              </HeroCard>

              <!-- 模糊页脚展示 -->
              <HeroCard class="max-w-[400px] h-[300px] relative">
                <HeroCardHeader class="absolute z-10 top-1 flex-col !items-start">
                  <p class="text-tiny text-white/60 uppercase font-bold">Your Day at a Glance</p>
                  <h4 class="text-white font-medium text-large">High-performance apps</h4>
                </HeroCardHeader>
                <img src="https://heroui.com/images/card-example-6.jpeg" class="z-0 w-full h-full object-cover" alt="Card background" />
                <HeroCardFooter is-blurred class="absolute bg-white/30 bottom-0 border-t-1 border-zinc-100/50 z-10 justify-between">
                  <div>
                    <p class="text-black text-tiny">Available soon.</p>
                    <p class="text-black text-tiny">Get notified.</p>
                  </div>
                  <Button size="sm" class="rounded-full bg-black text-white px-3 py-1 text-xs">Notify Me</Button>
                </HeroCardFooter>
              </HeroCard>

              <!-- 不同圆角展示 -->
              <div class="grid gap-4">
                <HeroCard radius="lg" class="border border-border/40">
                  <HeroCardBody><p class="text-sm">Radius LG</p></HeroCardBody>
                </HeroCard>
                <HeroCard radius="2xl" class="bg-primary/10 border border-primary/20">
                  <HeroCardBody><p class="text-sm">Radius 2xl (Brand style)</p></HeroCardBody>
                </HeroCard>
                <HeroCard is-blurred radius="full" class="border border-white/20">
                  <HeroCardBody class="items-center"><p class="text-sm font-medium">Blurred Glassmorphism</p></HeroCardBody>
                </HeroCard>
              </div>

              <!-- 用户提供的 React 风格示例适配 -->
              <div class="col-span-full mt-4">
                <p class="text-sm font-medium mb-4">React 风格高级示例 (适配版)</p>
                <HeroCard class="w-[200px] h-[200px]" radius="2xl">
                  <img
                    alt="Woman listing to music"
                    class="object-cover w-full h-full"
                    src="https://heroui.com/images/hero-card.jpeg"
                  />
                  <HeroCardFooter class="justify-between before:bg-white/10 border-white/20 border px-3 py-1 absolute before:rounded-xl rounded-2xl bottom-1 w-[calc(100%_-_8px)] ml-1 z-10 overflow-hidden bg-black/20 backdrop-blur-md">
                    <p class="text-[10px] text-white/80">Available soon.</p>
                    <Button
                      class="text-[10px] text-white bg-black/40 h-7 px-2 rounded-lg"
                      variant="ghost"
                    >
                      Notify me
                    </Button>
                  </HeroCardFooter>
                </HeroCard>
              </div>
            </div>
          </div>

          <div v-else-if="selectedId === 'card'" class="max-w-md">
            <Card class="rounded-2xl">
              <CardHeader>
                <CardTitle>Card 组件示例</CardTitle>
                <CardDescription>用于信息分组与内容承载</CardDescription>
              </CardHeader>
              <CardContent>这是卡片正文内容。</CardContent>
            </Card>
          </div>

          <div v-else-if="selectedId === 'carousel'" class="mx-auto w-full max-w-xl px-12">
            <Carousel class="w-full">
              <CarouselContent>
                <CarouselItem v-for="(item, index) in carouselItems" :key="index">
                  <div class="flex h-44 items-center justify-center rounded-2xl border bg-card/50 text-lg font-medium">
                    {{ item }}
                  </div>
                </CarouselItem>
              </CarouselContent>
              <CarouselPrevious class="rounded-full" />
              <CarouselNext class="rounded-full" />
            </Carousel>
          </div>

          <div v-else-if="selectedId === 'chart'" class="space-y-4 overflow-hidden">
            <div class="grid gap-4 xl:grid-cols-2">
              <ChartBarDemo />
              <ChartDonutDemo />
            </div>
            <ChartAreaDemo />
            <div class="grid gap-4 xl:grid-cols-2">
              <ChartLineDemo />
              <ChartTooltipDefaultDemo />
            </div>
          </div>

          <div v-else-if="selectedId === 'checkbox'" class="flex items-center gap-3">
            <Checkbox id="agree" v-model="checkboxValue" />
            <Label for="agree">同意组件库使用协议</Label>
          </div>

          <div v-else-if="selectedId === 'collapsible'" class="max-w-md space-y-3 rounded-2xl border bg-card/50 p-4">
            <Collapsible v-model:open="collapsibleOpen" class="space-y-3">
              <div class="flex items-center justify-between">
                <span class="text-sm font-medium">展开高级配置</span>
                <CollapsibleTrigger as-child>
                  <Button size="icon" variant="outline" class="rounded-full">
                    <Icons.ChevronsUpDown class="size-4" />
                  </Button>
                </CollapsibleTrigger>
              </div>
              <CollapsibleContent class="space-y-2">
                <Input placeholder="API Endpoint" />
                <Input placeholder="Token Scope" />
              </CollapsibleContent>
            </Collapsible>
          </div>

          <div v-else-if="selectedId === 'combobox'" class="max-w-md">
            <Combobox v-model="comboboxValue">
              <ComboboxAnchor class="w-full">
                <div class="relative w-full items-center">
                  <ComboboxInput class="pl-9" placeholder="选择框架" />
                  <span class="pointer-events-none absolute inset-y-0 left-0 flex items-center pl-3">
                    <Icons.Search class="size-4 text-muted-foreground" />
                  </span>
                  <ComboboxTrigger class="absolute inset-y-0 right-3 flex items-center">
                    <Icons.ChevronsUpDown class="size-4 text-muted-foreground" />
                  </ComboboxTrigger>
                </div>
              </ComboboxAnchor>
              <ComboboxList>
                <ComboboxEmpty>未找到匹配项</ComboboxEmpty>
                <ComboboxViewport>
                  <ComboboxGroup>
                    <ComboboxItem value="vue">
                      Vue 3
                      <ComboboxItemIndicator><Icons.Check class="size-4" /></ComboboxItemIndicator>
                    </ComboboxItem>
                    <ComboboxItem value="react">
                      React
                      <ComboboxItemIndicator><Icons.Check class="size-4" /></ComboboxItemIndicator>
                    </ComboboxItem>
                    <ComboboxItem value="svelte">
                      Svelte
                      <ComboboxItemIndicator><Icons.Check class="size-4" /></ComboboxItemIndicator>
                    </ComboboxItem>
                  </ComboboxGroup>
                </ComboboxViewport>
              </ComboboxList>
            </Combobox>
          </div>

          <div v-else-if="selectedId === 'command'" class="max-w-md overflow-hidden rounded-2xl border">
            <Command v-model="commandValue">
              <CommandInput placeholder="输入命令或搜索..." />
              <CommandList>
                <CommandEmpty>没有找到对应命令。</CommandEmpty>
                <CommandGroup heading="建议操作">
                  <CommandItem value="new" @select="showToast"><Icons.Plus class="size-4" /> 新建页面</CommandItem>
                  <CommandItem value="deploy" @select="showToast"><Icons.Rocket class="size-4" /> 触发部署</CommandItem>
                </CommandGroup>
                <CommandSeparator />
                <CommandGroup heading="跳转">
                  <CommandItem value="dashboard"><Icons.LayoutDashboard class="size-4" /> 控制台</CommandItem>
                  <CommandItem value="settings"><Icons.Settings class="size-4" /> 设置</CommandItem>
                </CommandGroup>
              </CommandList>
            </Command>
          </div>

          <div v-else-if="selectedId === 'context-menu'" class="max-w-md">
            <ContextMenu>
              <ContextMenuTrigger>
                <div class="flex h-28 items-center justify-center rounded-2xl border border-dashed bg-card/50 text-sm text-muted-foreground">
                  右键此区域打开菜单
                </div>
              </ContextMenuTrigger>
              <ContextMenuContent class="w-52 rounded-2xl">
                <ContextMenuItem><Icons.Copy class="mr-2 size-4" /> 复制</ContextMenuItem>
                <ContextMenuItem><Icons.Scissors class="mr-2 size-4" /> 剪切</ContextMenuItem>
                <ContextMenuSeparator />
                <ContextMenuCheckboxItem v-model="contextMenuPinned">固定到侧边栏</ContextMenuCheckboxItem>
              </ContextMenuContent>
            </ContextMenu>
          </div>

          <div v-else-if="selectedId === 'dialog'" class="space-y-4">
            <div class="rounded-2xl border bg-card/50 p-4 mb-4 space-y-3">
              <h3 class="text-sm font-semibold">使用说明（Dialog）</h3>
              <p class="text-sm text-muted-foreground">
                Dialog 用于“可编辑、可中断、信息量较高”的模态任务（如新建、编辑、配置），与 AlertDialog 的“危险确认”职责要区分。
              </p>
              <ul class="list-disc pl-5 space-y-2 text-sm text-muted-foreground">
                <li>
                  <span class="font-medium text-foreground">推荐结构：</span>
                  Dialog → DialogTrigger → DialogContent → DialogHeader / DialogDescription / DialogFooter。
                </li>
                <li>
                  <span class="font-medium text-foreground">触发按钮（项目规范）：</span>
                  触发入口用自定义 Button 的图标模式（<code>size="icon" + rounded-full</code>），并通过 Tooltip 明确动作名称。
                </li>
                <li>
                  <span class="font-medium text-foreground">样式规范：</span>
                  Content 使用根组件默认圆角与关闭按钮能力（含 Tooltip）；底部动作按钮遵循“图标按钮 + Tooltip”的统一规则。
                </li>
                <li>
                  <span class="font-medium text-foreground">底部按钮建议：</span>
                  取消按钮使用 <code>variant="outline"</code>，确认按钮使用 <code>variant="default"</code> 或 <code>variant="destructive"</code>；
                  两者均建议 <code>size="icon"</code> 并由 <code>DialogClose as-child</code> 包裹以保持关闭行为一致。
                </li>
                <li>
                  <span class="font-medium text-foreground">状态处理：</span>
                  提交成功后主动关闭；若校验失败，保留弹窗并在字段附近显示错误信息，不要直接关闭。
                </li>
                <li>
                  <span class="font-medium text-foreground">适用边界：</span>
                  复杂表单可放 Dialog；超长内容或移动端底部交互优先考虑 Drawer。
                </li>
                <li>
                  <span class="font-medium text-foreground">演示策略：</span>
                  以下示例按场景拆分为“新建任务 / 编辑成员 / 长内容审批”，避免单个弹窗堆叠过多内容。
                </li>
              </ul>
            </div>

            <div class="grid gap-4 sm:grid-cols-2">
              <div class="rounded-2xl border bg-card/40 p-4 space-y-3">
                <p class="text-sm font-medium">场景 1：新建任务（轻表单）</p>
                <Dialog>
                  <TooltipProvider>
                    <Tooltip>
                      <TooltipTrigger as-child>
                        <DialogTrigger as-child>
                          <Button variant="outline" size="icon" class="rounded-full">
                            <Icons.Plus class="size-4" />
                          </Button>
                        </DialogTrigger>
                      </TooltipTrigger>
                      <TooltipContent>新建任务</TooltipContent>
                    </Tooltip>
                  </TooltipProvider>
                  <DialogContent :show-close-button="false" class="sm:max-w-[460px]">
                    <DialogHeader>
                      <DialogTitle>新建任务</DialogTitle>
                      <DialogDescription>创建一个新的工作项。</DialogDescription>
                    </DialogHeader>
                    <div class="space-y-3">
                      <Input placeholder="任务名称" />
                      <Input placeholder="负责人" />
                      <Textarea class="min-h-24" placeholder="任务说明" />
                    </div>
                    <DialogFooter>
                      <TooltipProvider>
                        <Tooltip>
                          <TooltipTrigger as-child>
                            <DialogClose as-child>
                              <Button variant="outline" size="icon" class="rounded-full"><Icons.X class="size-4" /></Button>
                            </DialogClose>
                          </TooltipTrigger>
                          <TooltipContent>取消</TooltipContent>
                        </Tooltip>
                        <Tooltip>
                          <TooltipTrigger as-child>
                            <DialogClose as-child>
                              <Button variant="default" size="icon" class="rounded-full"><Icons.Check class="size-4" /></Button>
                            </DialogClose>
                          </TooltipTrigger>
                          <TooltipContent>确认</TooltipContent>
                        </Tooltip>
                      </TooltipProvider>
                    </DialogFooter>
                  </DialogContent>
                </Dialog>
              </div>

              <div class="rounded-2xl border bg-card/40 p-4 space-y-3">
                <p class="text-sm font-medium">场景 2：编辑成员（资料更新）</p>
                <Dialog>
                  <TooltipProvider>
                    <Tooltip>
                      <TooltipTrigger as-child>
                        <DialogTrigger as-child>
                          <Button variant="outline" size="icon" class="rounded-full">
                            <Icons.UserCog class="size-4" />
                          </Button>
                        </DialogTrigger>
                      </TooltipTrigger>
                      <TooltipContent>编辑成员</TooltipContent>
                    </Tooltip>
                  </TooltipProvider>
                  <DialogContent :show-close-button="false" class="sm:max-w-[460px]">
                    <DialogHeader>
                      <DialogTitle>编辑成员信息</DialogTitle>
                      <DialogDescription>更新成员的基础资料与角色。</DialogDescription>
                    </DialogHeader>
                    <div class="space-y-3">
                      <Input placeholder="成员姓名" />
                      <Input type="email" placeholder="企业邮箱" />
                      <Select>
                        <SelectTrigger class="w-full rounded-2xl"><SelectValue placeholder="选择角色" /></SelectTrigger>
                        <SelectContent class="rounded-2xl">
                          <SelectItem value="admin">管理员</SelectItem>
                          <SelectItem value="editor">编辑</SelectItem>
                          <SelectItem value="viewer">访客</SelectItem>
                        </SelectContent>
                      </Select>
                    </div>
                    <DialogFooter>
                      <TooltipProvider>
                        <Tooltip>
                          <TooltipTrigger as-child>
                            <DialogClose as-child>
                              <Button variant="outline" size="icon" class="rounded-full"><Icons.X class="size-4" /></Button>
                            </DialogClose>
                          </TooltipTrigger>
                          <TooltipContent>取消</TooltipContent>
                        </Tooltip>
                        <Tooltip>
                          <TooltipTrigger as-child>
                            <DialogClose as-child>
                              <Button variant="default" size="icon" class="rounded-full"><Icons.Check class="size-4" /></Button>
                            </DialogClose>
                          </TooltipTrigger>
                          <TooltipContent>保存</TooltipContent>
                        </Tooltip>
                      </TooltipProvider>
                    </DialogFooter>
                  </DialogContent>
                </Dialog>
              </div>

              <div class="rounded-2xl border bg-card/40 p-4 space-y-3 sm:col-span-2">
                <p class="text-sm font-medium">场景 3：长内容审批（内部滚动）</p>
                <Dialog>
                  <TooltipProvider>
                    <Tooltip>
                      <TooltipTrigger as-child>
                        <DialogTrigger as-child>
                          <Button variant="outline" size="icon" class="rounded-full">
                            <Icons.ScrollText class="size-4" />
                          </Button>
                        </DialogTrigger>
                      </TooltipTrigger>
                      <TooltipContent>审批详情</TooltipContent>
                    </Tooltip>
                  </TooltipProvider>
                  <DialogContent :show-close-button="false" class="sm:max-w-[560px]">
                    <DialogHeader>
                      <DialogTitle>审批详情</DialogTitle>
                      <DialogDescription>长内容采用内部滚动，底部操作始终可见。</DialogDescription>
                    </DialogHeader>
                    <div class="max-h-80 space-y-3 overflow-y-auto pr-2">
                      <Textarea class="min-h-28" placeholder="项目背景" />
                      <Textarea class="min-h-28" placeholder="执行计划" />
                      <Textarea class="min-h-28" placeholder="风险评估" />
                      <Textarea class="min-h-28" placeholder="资源需求" />
                      <Textarea class="min-h-28" placeholder="验收标准" />
                    </div>
                    <DialogFooter>
                      <TooltipProvider>
                        <Tooltip>
                          <TooltipTrigger as-child>
                            <DialogClose as-child>
                              <Button variant="outline" size="icon" class="rounded-full"><Icons.X class="size-4" /></Button>
                            </DialogClose>
                          </TooltipTrigger>
                          <TooltipContent>取消</TooltipContent>
                        </Tooltip>
                        <Tooltip>
                          <TooltipTrigger as-child>
                            <DialogClose as-child>
                              <Button variant="default" size="icon" class="rounded-full"><Icons.Check class="size-4" /></Button>
                            </DialogClose>
                          </TooltipTrigger>
                          <TooltipContent>通过</TooltipContent>
                        </Tooltip>
                      </TooltipProvider>
                    </DialogFooter>
                  </DialogContent>
                </Dialog>
              </div>
            </div>
          </div>

          <div v-else-if="selectedId === 'drawer'" class="space-y-4">
            <div class="grid gap-4 md:grid-cols-3">
              <div class="space-y-2 rounded-2xl border bg-card/40 p-4">
                <p class="text-sm font-medium">场景 1：基础设置抽屉</p>
                <p class="text-xs text-muted-foreground">标准底部抽屉，包含表单与确认操作。</p>
                <DrawerBasicDemo />
              </div>

              <div class="space-y-2 rounded-2xl border bg-card/40 p-4">
                <p class="text-sm font-medium">场景 2：多方向抽屉</p>
                <p class="text-xs text-muted-foreground">展示 top / left / right 多方向弹出布局。</p>
                <DrawerDirectionDemo />
              </div>

              <div class="space-y-2 rounded-2xl border bg-card/40 p-4">
                <p class="text-sm font-medium">场景 3：复杂菜单抽屉</p>
                <p class="text-xs text-muted-foreground">承载多操作项与危险操作按钮。</p>
                <DrawerResponsiveDemo />
              </div>
            </div>
          </div>

          <div v-else-if="selectedId === 'dropdown-menu'">
            <DropdownMenu>
              <DropdownMenuTrigger as-child>
                <Button variant="outline" size="icon" class="rounded-full"><Icons.MoreVertical class="size-4" /></Button>
              </DropdownMenuTrigger>
              <DropdownMenuContent class="w-48 rounded-2xl" align="start">
                <DropdownMenuLabel>操作菜单</DropdownMenuLabel>
                <DropdownMenuSeparator />
                <DropdownMenuItem><Icons.User class="mr-2 size-4" /> 个人资料</DropdownMenuItem>
                <DropdownMenuItem><Icons.Settings class="mr-2 size-4" /> 设置</DropdownMenuItem>
              </DropdownMenuContent>
            </DropdownMenu>
          </div>

          <div v-else-if="selectedId === 'empty'" class="max-w-md">
            <Empty>
              <EmptyHeader>
                <EmptyMedia variant="icon"><Icons.Inbox class="size-6" /></EmptyMedia>
                <EmptyTitle>暂无数据</EmptyTitle>
                <EmptyDescription>当前筛选条件下暂无结果，请尝试调整筛选项。</EmptyDescription>
              </EmptyHeader>
              <EmptyContent>
                <Button size="icon" class="rounded-full"><Icons.RefreshCcw class="size-4" /></Button>
              </EmptyContent>
            </Empty>
          </div>

          <div v-else-if="selectedId === 'field'" class="max-w-md">
            <FieldGroup>
              <FieldSet>
                <FieldTitle>基础字段组</FieldTitle>
                <FieldDescription>Field 负责字段级布局和说明信息组织。</FieldDescription>
                <Field class="mt-4">
                  <FieldLabel>邮箱地址</FieldLabel>
                  <FieldContent>
                    <Input placeholder="name@company.com" />
                  </FieldContent>
                </Field>
              </FieldSet>
            </FieldGroup>
          </div>

          <div v-else-if="selectedId === 'form'" class="max-w-md rounded-2xl border bg-card/50 p-4">
            <Form v-slot="{ handleSubmit }" :initial-values="{ email: '' }">
              <form class="space-y-4" @submit="handleSubmit(handleFormSubmit)">
                <FormField v-slot="{ componentField }" name="email">
                  <FormItem>
                    <FormLabel>邮箱</FormLabel>
                    <FormControl>
                      <Input type="email" placeholder="name@company.com" v-bind="componentField" />
                    </FormControl>
                    <FormDescription>用于接收系统通知。</FormDescription>
                    <FormMessage />
                  </FormItem>
                </FormField>
                <Button type="submit" size="icon" class="rounded-full">
                  <Icons.Send class="size-4" />
                </Button>
              </form>
            </Form>
          </div>

          <div v-else-if="selectedId === 'glassmorphism'" class="space-y-6">
            <GlassmorphismDemo />
          </div>

          <div v-else-if="selectedId === 'hover-card'">
            <HoverCard>
              <HoverCardTrigger as-child>
                <Button variant="outline" size="icon" class="rounded-full"><Icons.UserRound class="size-4" /></Button>
              </HoverCardTrigger>
              <HoverCardContent class="rounded-2xl">
                <p class="text-sm">@microswift - 企业官网组件库维护中</p>
              </HoverCardContent>
            </HoverCard>
          </div>

          <div v-else-if="selectedId === 'input'" class="max-w-md space-y-3">
            <Input placeholder="请输入关键词" />
            <Input type="password" placeholder="请输入密码" />
          </div>

          <div v-else-if="selectedId === 'input-group'" class="max-w-md">
            <InputGroup>
              <InputGroupAddon>https://</InputGroupAddon>
              <InputGroupInput placeholder="your-domain.com" />
              <InputGroupAddon align="inline-end">
                <InputGroupButton size="icon-sm" variant="ghost">
                  <Icons.ArrowUpRight class="size-4" />
                </InputGroupButton>
              </InputGroupAddon>
            </InputGroup>
          </div>

          <div v-else-if="selectedId === 'input-otp'" class="space-y-4">
            <InputOTP v-model="otpValue" :maxlength="6">
              <InputOTPGroup>
                <InputOTPSlot v-for="index in 3" :key="index" :index="index - 1" />
              </InputOTPGroup>
              <InputOTPSeparator />
              <InputOTPGroup>
                <InputOTPSlot v-for="index in 3" :key="index + 10" :index="index + 2" />
              </InputOTPGroup>
            </InputOTP>
            <p class="text-sm text-muted-foreground">当前值：{{ otpValue || '未输入' }}</p>
          </div>

          <div v-else-if="selectedId === 'item'" class="max-w-xl space-y-4">
            <ItemGroup>
              <Item variant="outline">
                <ItemMedia variant="icon"><Icons.FileText class="size-4" /></ItemMedia>
                <ItemContent>
                  <ItemHeader>
                    <ItemTitle>组件文档</ItemTitle>
                    <ItemDescription>最后更新于 2 小时前</ItemDescription>
                  </ItemHeader>
                </ItemContent>
                <ItemActions>
                  <Button variant="ghost" size="icon" class="rounded-full"><Icons.MoreHorizontal class="size-4" /></Button>
                </ItemActions>
              </Item>
              <ItemSeparator />

              <Item variant="muted">
                <ItemMedia>
                  <Avatar size="sm">
                    <AvatarFallback>MS</AvatarFallback>
                  </Avatar>
                </ItemMedia>
                <ItemContent>
                  <ItemHeader>
                    <ItemTitle>Microswift 设计评审</ItemTitle>
                    <ItemDescription>今天 15:00 · 线上会议</ItemDescription>
                  </ItemHeader>
                </ItemContent>
                <ItemActions>
                  <Badge variant="secondary" class="rounded-full">进行中</Badge>
                </ItemActions>
              </Item>
            </ItemGroup>
          </div>

          <div v-else-if="selectedId === 'kbd'" class="space-y-4">
            <div class="flex items-center gap-2 text-sm">
              保存：
              <KbdGroup>
                <Kbd>Ctrl</Kbd>
                <Kbd>S</Kbd>
              </KbdGroup>
            </div>
          </div>

          <div v-else-if="selectedId === 'label'" class="max-w-md space-y-4">
            <div class="space-y-2">
              <Label for="label-demo">用户名</Label>
              <Input id="label-demo" placeholder="请输入用户名" />
            </div>
            <div class="space-y-2">
              <div class="flex items-center justify-between">
                <Label for="label-email">企业邮箱</Label>
                <span class="text-xs text-muted-foreground">必填</span>
              </div>
              <Input id="label-email" type="email" placeholder="name@company.com" />
            </div>
            <div class="flex items-center gap-2">
              <Checkbox id="label-agree" />
              <Label for="label-agree">我已阅读并同意服务条款</Label>
            </div>
          </div>

          <div v-else-if="selectedId === 'menubar'" class="space-y-4">
            <Menubar>
              <MenubarMenu>
                <MenubarTrigger>文件</MenubarTrigger>
                <MenubarContent>
                  <MenubarItem>新建 <MenubarShortcut>⌘N</MenubarShortcut></MenubarItem>
                  <MenubarItem>保存 <MenubarShortcut>⌘S</MenubarShortcut></MenubarItem>
                  <MenubarSeparator />
                  <MenubarCheckboxItem v-model="menubarStatusBar">显示状态栏</MenubarCheckboxItem>
                </MenubarContent>
              </MenubarMenu>

              <MenubarMenu>
                <MenubarTrigger>编辑</MenubarTrigger>
                <MenubarContent>
                  <MenubarItem>撤销 <MenubarShortcut>⌘Z</MenubarShortcut></MenubarItem>
                  <MenubarItem>重做 <MenubarShortcut>⇧⌘Z</MenubarShortcut></MenubarItem>
                  <MenubarSeparator />
                  <MenubarItem>查找 <MenubarShortcut>⌘F</MenubarShortcut></MenubarItem>
                </MenubarContent>
              </MenubarMenu>

              <MenubarMenu>
                <MenubarTrigger>帮助</MenubarTrigger>
                <MenubarContent>
                  <MenubarItem>
                    <Icons.BookOpen class="mr-2 size-4" />
                    文档中心
                  </MenubarItem>
                  <MenubarItem>
                    <Icons.MessageCircle class="mr-2 size-4" />
                    联系支持
                  </MenubarItem>
                </MenubarContent>
              </MenubarMenu>
            </Menubar>
          </div>

          <div v-else-if="selectedId === 'native-select'" class="max-w-sm">
            <NativeSelect v-model="nativeSelectValue" class="w-full">
              <NativeSelectOptGroup label="前端框架">
                <NativeSelectOption value="vue">Vue 3</NativeSelectOption>
                <NativeSelectOption value="react">React</NativeSelectOption>
              </NativeSelectOptGroup>
              <NativeSelectOptGroup label="后端语言">
                <NativeSelectOption value="go">Go</NativeSelectOption>
                <NativeSelectOption value="node">Node.js</NativeSelectOption>
              </NativeSelectOptGroup>
            </NativeSelect>
          </div>

          <div v-else-if="selectedId === 'popover'" class="space-y-4">
            <Popover>
              <PopoverTrigger as-child>
                <Button variant="outline" size="icon" class="rounded-full"><Icons.FileText class="size-4" /></Button>
              </PopoverTrigger>
              <PopoverContent class="rounded-2xl">
                <div class="space-y-2">
                  <Label for="popover-note">备注</Label>
                  <Textarea id="popover-note" v-model="popoverNote" class="min-h-24" />
                </div>
              </PopoverContent>
            </Popover>
          </div>

          <div v-else-if="selectedId === 'select'" class="max-w-sm">
            <Select v-model="selectValue">
              <SelectTrigger class="w-[220px] rounded-2xl"><SelectValue placeholder="选择语言" /></SelectTrigger>
              <SelectContent class="rounded-2xl">
                <SelectGroup>
                  <SelectLabel>语言</SelectLabel>
                  <SelectItem value="vue">Vue.js</SelectItem>
                  <SelectItem value="react">React</SelectItem>
                  <SelectItem value="go">Golang</SelectItem>
                </SelectGroup>
              </SelectContent>
            </Select>
          </div>

          <div v-else-if="selectedId === 'separator'" class="space-y-4 max-w-md">
            <div class="text-sm">上半部分内容</div>
            <Separator />
            <div class="text-sm">下半部分内容</div>
          </div>

          <div v-else-if="selectedId === 'sonner'" class="space-y-3">
            <Button size="icon" class="rounded-full" @click="showToast"><Icons.Bell class="size-4" /></Button>
            <p class="text-sm text-muted-foreground">点击按钮触发 Sonner 通知。</p>
          </div>

          <div v-else-if="selectedId === 'table'" class="overflow-hidden rounded-2xl border bg-card/50">
            <Table>
              <TableCaption>最近的交易记录</TableCaption>
              <TableHeader class="bg-muted/30">
                <TableRow>
                  <TableHead class="w-[100px]">状态</TableHead>
                  <TableHead>方法</TableHead>
                  <TableHead>金额</TableHead>
                  <TableHead class="text-right">日期</TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                <TableRow v-for="i in 3" :key="i">
                  <TableCell><Badge variant="outline" class="rounded-full">已完成</Badge></TableCell>
                  <TableCell>信用卡</TableCell>
                  <TableCell>¥ 250.00</TableCell>
                  <TableCell class="text-right text-muted-foreground">2024-02-{{ 10 + i }}</TableCell>
                </TableRow>
              </TableBody>
            </Table>
          </div>

          <div v-else-if="selectedId === 'tabs'" class="space-y-4">
            <Tabs default-value="account" class="space-y-4">
              <TabsList class="w-fit">
                <TabsTrigger value="account">账户资料</TabsTrigger>
                <TabsTrigger value="security">安全设置</TabsTrigger>
                <TabsTrigger value="notifications">通知偏好</TabsTrigger>
              </TabsList>

              <div class="w-[360px] max-w-full">
                <TabsContent value="account" class="space-y-3 rounded-2xl border bg-card/60 p-4">
                  <Input placeholder="姓名" />
                  <Input placeholder="企业邮箱" />
                  <div class="flex items-center justify-between text-xs text-muted-foreground">
                    <span>资料修改后将同步到成员目录</span>
                    <Badge variant="secondary" class="rounded-full">已同步</Badge>
                  </div>
                </TabsContent>

                <TabsContent value="security" class="space-y-3 rounded-2xl border bg-card/60 p-4">
                  <Input type="password" placeholder="当前密码" />
                  <Input type="password" placeholder="新密码" />
                  <Input type="password" placeholder="确认新密码" />
                </TabsContent>

                <TabsContent value="notifications" class="space-y-3 rounded-2xl border bg-card/60 p-4">
                  <div class="rounded-2xl border bg-background/70 p-3 text-sm">产品更新通知：已开启</div>
                  <div class="rounded-2xl border bg-background/70 p-3 text-sm">账单提醒：每周一推送</div>
                  <div class="rounded-2xl border bg-background/70 p-3 text-sm">安全告警：实时推送</div>
                </TabsContent>
              </div>
            </Tabs>
          </div>

          <div v-else-if="selectedId === 'textarea'" class="max-w-xl space-y-3">
            <Textarea v-model="textareaValue" class="min-h-32" />
            <p class="text-sm text-muted-foreground">字符数：{{ String(textareaValue).length }}</p>
          </div>

          <div v-else-if="selectedId === 'tooltip'" class="flex gap-4">
            <TooltipProvider>
              <Tooltip>
                <TooltipTrigger as-child>
                  <Button size="icon" variant="outline" class="rounded-full"><Icons.Info class="size-4" /></Button>
                </TooltipTrigger>
                <TooltipContent>这是 Tooltip 文本</TooltipContent>
              </Tooltip>
            </TooltipProvider>
          </div>

          <div v-else-if="selectedId === 'tooltip-advanced'" class="space-y-8">
            <p class="text-sm text-muted-foreground">鼠标悬浮在元素上，Tooltip 跟随光标出现并自动检测视口边界。<code>#content</code> slot 支持任意富文本内容。</p>

            <!-- 文字触发 -->
            <div class="space-y-2">
              <p class="text-xs font-medium text-muted-foreground uppercase tracking-wide">文字触发</p>
              <div class="flex flex-wrap gap-6 items-center">
                <TooltipCard>
                  <span class="cursor-default underline decoration-dotted decoration-muted-foreground text-sm">悬浮查看详情</span>
                  <template #content>
                    <p class="font-medium text-foreground">这是一条提示</p>
                    <p class="mt-1 text-xs text-muted-foreground">支持多行文本内容，自动撑高卡片。</p>
                  </template>
                </TooltipCard>

                <TooltipCard>
                  <span class="cursor-default underline decoration-dotted decoration-muted-foreground text-sm">Vue 3 技术栈</span>
                  <template #content>
                    <ul class="space-y-1 text-xs">
                      <li>⚡ Vite 8 构建</li>
                      <li>🎨 TailwindCSS v4</li>
                      <li>🧩 shadcn-vue 组件库</li>
                      <li>🗃️ Pinia 状态管理</li>
                    </ul>
                  </template>
                </TooltipCard>
              </div>
            </div>

            <!-- 按钮触发 -->
            <div class="space-y-2">
              <p class="text-xs font-medium text-muted-foreground uppercase tracking-wide">按钮触发</p>
              <div class="flex flex-wrap gap-4 items-center">
                <TooltipCard>
                  <Button size="icon" variant="outline" class="rounded-full">
                    <Icons.Info class="size-4" />
                  </Button>
                  <template #content>
                    <p class="font-medium text-foreground">关于此功能</p>
                    <p class="mt-1 text-xs text-muted-foreground">点击可查看完整说明文档。</p>
                  </template>
                </TooltipCard>

                <TooltipCard>
                  <Button size="icon" variant="outline" class="rounded-full">
                    <Icons.Settings class="size-4" />
                  </Button>
                  <template #content>
                    <p class="font-medium text-foreground">系统设置</p>
                    <p class="mt-1 text-xs text-muted-foreground">管理通知、主题、语言等偏好。</p>
                  </template>
                </TooltipCard>

                <TooltipCard>
                  <Button size="icon" variant="destructive" class="rounded-full">
                    <Icons.Trash2 class="size-4" />
                  </Button>
                  <template #content>
                    <p class="font-medium">删除操作</p>
                    <p class="mt-1 text-xs opacity-70">此操作不可撤销，请谨慎确认。</p>
                  </template>
                </TooltipCard>
              </div>
            </div>

            <!-- 卡片式富内容 -->
            <div class="space-y-2">
              <p class="text-xs font-medium text-muted-foreground uppercase tracking-wide">富内容卡片</p>
              <div class="flex flex-wrap gap-6 items-center">
                <TooltipCard>
                  <div class="flex items-center gap-2 cursor-default rounded-lg border px-3 py-2 text-sm hover:bg-accent transition-colors">
                    <Icons.User class="size-4" />
                    <span>用户资料</span>
                  </div>
                  <template #content>
                    <div class="flex items-center gap-3">
                      <div class="size-9 rounded-full bg-gradient-to-br from-violet-500 to-blue-500 flex items-center justify-center text-white text-sm font-bold">A</div>
                      <div>
                        <p class="font-medium text-foreground text-sm">Alice Chen</p>
                        <p class="text-xs text-muted-foreground">alice@example.com</p>
                      </div>
                    </div>
                    <div class="mt-2 pt-2 border-t flex gap-3 text-xs text-muted-foreground">
                      <span>📦 12 项目</span>
                      <span>⭐ 4.9 评分</span>
                    </div>
                  </template>
                </TooltipCard>
              </div>
            </div>
          </div>

          <div v-else-if="selectedId === 'noise-background'" class="space-y-6">
            <div class="space-y-3">
              <div class="flex justify-center">
                <NoiseBackground
                  container-class-name="mx-auto w-fit rounded-full p-2"
                  :gradient-colors="[
                    'rgb(255, 100, 150)',
                    'rgb(100, 150, 255)',
                    'rgb(255, 200, 100)',
                  ]"
                >
                  <button
                    type="button"
                    class="h-full w-full cursor-pointer rounded-full bg-linear-to-r from-neutral-100 via-neutral-100 to-white px-4 py-2 text-black shadow-[0px_2px_0px_0px_var(--color-neutral-50)_inset,0px_0.5px_1px_0px_var(--color-neutral-400)] transition-all duration-100 active:scale-98 dark:from-black dark:via-black dark:to-neutral-900 dark:text-white dark:shadow-[0px_1px_0px_0px_var(--color-neutral-950)_inset,0px_1px_0px_0px_var(--color-neutral-800)]"
                  >
                    Start publishing →
                  </button>
                </NoiseBackground>
              </div>
              <p class="text-center text-sm text-muted-foreground">官方风格噪点渐变容器按钮示例</p>
            </div>

            <div class="space-y-3">
              <div class="mx-auto max-w-sm">
                <NoiseBackground
                  container-class-name="rounded-3xl p-2"
                  :gradient-colors="[
                    'rgb(255, 100, 150)',
                    'rgb(100, 150, 255)',
                    'rgb(255, 200, 100)',
                  ]"
                >
                  <div :class="cn('flex h-full min-h-80 flex-col overflow-hidden rounded-3xl bg-white text-center dark:bg-neutral-800')">
                    <img
                      src="https://assets.aceternity.com/blog/how-to-create-a-bento-grid.png"
                      alt="Task Complete"
                      class="h-60 w-full rounded-3xl object-cover"
                    >
                    <div class="px-4 py-2">
                      <h3 class="text-left text-lg font-semibold text-balance text-neutral-800 dark:text-neutral-200">
                        How to create a bento grid with Tailwind
                      </h3>
                      <p class="mt-2 text-left text-sm text-neutral-600 dark:text-neutral-400">
                        Learn how to create a bento grid with Tailwind CSS, Next.js and Framer Motion.
                      </p>
                    </div>
                  </div>
                </NoiseBackground>
              </div>
              <p class="text-center text-sm text-muted-foreground">官方风格噪点渐变卡片示例</p>
            </div>
          </div>

          <div v-else-if="selectedId === 'animated-theme-toggler'" class="space-y-3">
            <AnimatedThemeToggler class="border rounded-full" />
            <p class="text-sm text-muted-foreground">点击图标切换亮暗主题并触发转场动画。</p>
          </div>

          <div v-else-if="selectedId === 'text-hover-effect'" class="space-y-6">
            <div class="rounded-2xl border bg-card/40 p-4 space-y-3">
              <h3 class="text-sm font-semibold">使用示例</h3>
              <pre class="overflow-x-auto rounded-xl bg-muted/60 p-4 text-xs leading-relaxed"><code>&lt;script setup&gt;
import { TextHoverEffect } from '@/components/ui/a-aceternity'
&lt;/script&gt;

&lt;template&gt;
  &lt;!-- 即时跟随 --&gt;
  &lt;div class="h-40 w-full"&gt;
    &lt;TextHoverEffect text="Aceternity" /&gt;
  &lt;/div&gt;

  &lt;!-- 平滑跟随，duration 单位：秒 --&gt;
  &lt;div class="h-40 w-full"&gt;
    &lt;TextHoverEffect text="UI Magic" :duration="0.5" /&gt;
  &lt;/div&gt;
&lt;/template&gt;</code></pre>
            </div>
            <TextHoverEffectDemo />
          </div>

          <div v-else-if="selectedId === 'expandable-card'" class="space-y-6">
            <div class="rounded-2xl border bg-card/40 p-4 space-y-3">
              <h3 class="text-sm font-semibold">使用示例</h3>
              <pre class="overflow-x-auto rounded-xl bg-muted/60 p-4 text-xs leading-relaxed"><code>&lt;script setup&gt;
import { ExpandableCard } from '@/components/ui/a-aceternity'
import type { CardItem } from '@/components/ui/a-aceternity'

const cards: CardItem[] = [
  {
    title: 'Summertime Sadness',
    description: 'Lana Del Rey',
    src: 'https://assets.aceternity.com/demos/lana-del-rey.jpeg',
    ctaText: 'Play',
    ctaLink: 'https://ui.aceternity.com/templates',
    content: '&lt;p&gt;Lana Del Rey, an iconic American singer-songwriter...&lt;/p&gt;',
  },
]
&lt;/script&gt;

&lt;template&gt;
  &lt;!-- 列表布局（默认）--&gt;
  &lt;ExpandableCard :cards="cards" /&gt;

  &lt;!-- 网格布局 --&gt;
  &lt;ExpandableCard :cards="cards" layout="grid" /&gt;
&lt;/template&gt;</code></pre>
            </div>
            <ExpandableCardDemo />
          </div>
        </div>
      </div>
    </main>
  </div>
</template>

<style scoped>
.gallery-root {
  background: var(--background);
}

.sidebar-gradient {
  background: linear-gradient(
    170deg,
    oklch(0.97 0.012 264 / 60%) 0%,
    oklch(1 0 0 / 0%) 65%
  );
}

.dark .sidebar-gradient {
  background: linear-gradient(
    170deg,
    oklch(0.19 0.018 265 / 70%) 0%,
    oklch(0 0 0 / 0%) 65%
  );
}

.main-area {
  background:
    radial-gradient(ellipse 70% 40% at 80% 10%, oklch(0.72 0.08 264 / 5%), transparent),
    var(--background);
}

.dark .main-area {
  background:
    radial-gradient(ellipse 70% 40% at 80% 10%, oklch(0.42 0.1 264 / 8%), transparent),
    var(--background);
}

.custom-scrollbar::-webkit-scrollbar {
  width: 5px;
}
.custom-scrollbar::-webkit-scrollbar-track {
  background: transparent;
}
.custom-scrollbar::-webkit-scrollbar-thumb {
  background-color: rgba(156, 163, 175, 0.15);
  border-radius: 9999px;
}
.custom-scrollbar::-webkit-scrollbar-thumb:hover {
  background-color: rgba(156, 163, 175, 0.28);
}
</style>
