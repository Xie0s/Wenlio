# HomepageEditorPage — 添加区块 Popover 实现总结

> 文件：`vue/src/views/admin/HomepageEditorPage.vue`

---

## 一、功能概述

顶部工具栏右侧有一个「添加区块」按钮（Plus 图标），点击后弹出一个 Popover，以 2 列图标网格的形式展示当前**可添加**的区块类型，点击任意区块即可将其追加到首页布局。

---

## 二、核心状态

```ts
const addPopoverOpen = ref(false)          // 控制 Popover 开/关

const SECTION_ICONS: Record<string, unknown> = {
  navbar: Navigation,
  hero: ImageIcon,
  introduction: LayoutGrid,
  theme_list: Library,
  cta: Megaphone,
  footer: PanelBottom,
}

const addableSections = computed(() =>
  (Object.keys(SECTION_TYPE_META) as SectionType[]).filter(t => editor.canAddSection(t)),
)
```

- `SECTION_ICONS`：区块类型 → Lucide 图标的静态映射表。
- `addableSections`：从 `SECTION_TYPE_META`（所有区块类型的元信息）中过滤掉已不可再添加的类型（某些区块如 navbar/footer 只允许存在一个）。

---

## 三、添加逻辑

```ts
function addSection(type: SectionType) {
  editor.addSection(type)       // 委托给 editor composable 修改 layout
  addPopoverOpen.value = false  // 添加后自动关闭 Popover
}
```

---

## 四、模板结构

```html
<Popover v-model:open="addPopoverOpen">

  <!-- 锚点绑定到触发按钮，PopoverContent 相对此元素定位 -->
  <PopoverAnchor as-child>

    <!-- Tooltip 嵌套在锚点内，与 Popover 共享同一个触发元素 -->
    <Tooltip>
      <TooltipTrigger as-child>
        <Button
          variant="ghost" size="icon" class="rounded-full"
          :disabled="!addableSections.length"
          @click="addPopoverOpen = !addPopoverOpen"
        >
          <!-- Plus 在 Popover 打开时旋转 45° → 视觉上变为 × -->
          <Plus class="h-4 w-4 transition-transform duration-200"
                :class="{ 'rotate-45': addPopoverOpen }" />
        </Button>
      </TooltipTrigger>
      <!-- Popover 打开时隐藏 Tooltip，避免两者同时出现 -->
      <TooltipContent v-if="!addPopoverOpen">添加区块</TooltipContent>
    </Tooltip>

  </PopoverAnchor>

  <!-- Popover 内容：2 列图标网格 -->
  <PopoverContent align="end" :side-offset="8" class="w-44 rounded-2xl p-2">
    <div class="grid grid-cols-2 gap-1.5">
      <button
        v-for="type in addableSections" :key="type"
        class="flex flex-col items-center justify-center gap-1.5 p-2.5
               rounded-xl hover:bg-muted transition-colors aspect-square"
        @click="addSection(type)"
      >
        <component :is="SECTION_ICONS[type]" class="h-5 w-5 text-muted-foreground" />
        <span class="text-[11px] leading-tight text-center">
          {{ SECTION_TYPE_META[type].label }}
        </span>
      </button>
    </div>
    <!-- 所有区块均已添加时的空态提示 -->
    <p v-if="!addableSections.length" class="py-2 text-xs text-muted-foreground text-center">
      所有区块已添加
    </p>
  </PopoverContent>

</Popover>
```

---

## 五、关键设计要点

| 要点 | 说明 |
|------|------|
| **`PopoverAnchor as-child`** | 将锚点绑定到按钮元素自身，`PopoverContent` 相对按钮定位，而非默认的 `PopoverTrigger` |
| **Tooltip 内嵌锚点** | Tooltip 与 Popover 共用同一个 DOM 节点（Button），通过 `as-child` 传递透明 |
| **`v-if="!addPopoverOpen"` on TooltipContent** | Popover 展开时屏蔽 Tooltip，防止两个浮层叠加显示 |
| **Plus 旋转动画** | `rotate-45` + `transition-transform` 让 + 变成 ×，无需额外图标，提供明确的关闭意图反馈 |
| **按钮禁用态** | `!addableSections.length` 为空时禁用按钮，阻止点击 |
| **空态文案** | `addableSections` 为空时在 Popover 内显示「所有区块已添加」提示，确保内容区不为空白 |

---

## 六、组件依赖

- `Popover / PopoverAnchor / PopoverContent` — `@/components/ui/popover`
- `Tooltip / TooltipContent / TooltipTrigger` — `@/components/ui/tooltip`
- `SECTION_TYPE_META / SectionType` — `@/components/personalization/types`
- `useHomepageEditor` — `@/composables/personalization`（提供 `canAddSection` / `addSection`）
