<!--
文档：TooltipCard 使用规范
职责：说明 Aceternity UI 风格鼠标跟随 Tooltip 组件的接入方式、Props/Slots 定义、
      视觉规范与生产环境已知坑点，避免重复踩坑。
主要接口：
  - Props：containerClassName
  - Slots：default（触发元素）、content（Tooltip 富内容）
  - 视觉规范：glass + ring-inset 组合边框
  - 已知坑：backdrop-filter 在 Lightning CSS 编译时的顺序问题
-->

# TooltipCard 使用规范

## 背景与目标

`TooltipCard` 移植自 [Aceternity UI Tooltip](https://ui.aceternity.com/components/tooltip-card)，
以 Vue 3 Composition API 重写，提供跟随鼠标的富内容 Tooltip，支持视口边界检测与弹簧高度动画。

组件位于 `src/components/ui/a-aceternity/TooltipCard.vue`，通过
`src/components/ui/a-aceternity/index.ts` 统一导出。

## Props

| Prop | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| `containerClassName` | `string` | `undefined` | 追加到外层容器的 Tailwind 类名 |

## Slots

| Slot | 说明 |
|------|------|
| `default` | 触发 Tooltip 的元素（任意 HTML / 组件） |
| `#content` | Tooltip 浮层内容，支持任意富文本 / 组件结构 |

## 示例

### 文字触发

```vue
<TooltipCard>
  <span class="underline decoration-dotted cursor-default">悬浮查看详情</span>
  <template #content>
    <p class="font-medium">标题</p>
    <p class="mt-1 text-xs text-muted-foreground">支持多行文本。</p>
  </template>
</TooltipCard>
```

### 按钮触发

```vue
<TooltipCard>
  <Button size="icon" variant="outline" class="rounded-full">
    <Icons.Info class="size-4" />
  </Button>
  <template #content>
    <p class="font-medium">关于此功能</p>
    <p class="mt-1 text-xs text-muted-foreground">点击查看完整说明。</p>
  </template>
</TooltipCard>
```

### 富内容卡片

```vue
<TooltipCard>
  <div class="flex items-center gap-2 cursor-default rounded-lg border px-3 py-2 text-sm">
    <Icons.User class="size-4" /><span>用户资料</span>
  </div>
  <template #content>
    <div class="flex items-center gap-3">
      <div class="size-9 rounded-full bg-gradient-to-br from-violet-500 to-blue-500
                  flex items-center justify-center text-white font-bold">A</div>
      <div>
        <p class="font-medium text-sm">Alice Chen</p>
        <p class="text-xs text-muted-foreground">alice@example.com</p>
      </div>
    </div>
  </template>
</TooltipCard>
```

## 视觉规范

- **浮层背景**：使用全局 `.glass` 工具类（`src/style.css`），提供扁平化毛玻璃效果。
- **轮廓线**：`ring-[0.5px] ring-inset ring-black/[0.07] dark:ring-white/[0.08]`。
  - 使用 `ring-inset`（`box-shadow` 实现）而非 `border-color`，原因：
    `.glass` 使用 `border` 简写属性，会覆盖后续的 `border-color`，
    `ring-inset` 与 `border` 正交，无级联冲突。
- **圆角**：`rounded-3xl`（浮层），外层容器默认 `inline-block`。
- **最小宽度**：`min-w-[15rem]`（240px），与原 React 版本保持一致。

### 扁平化补充约束

- Tooltip 浮层禁止增加 `shadow-*` 外投影类。
- Tooltip 浮层禁止叠加渐变高光背景（`bg-gradient-*`）。
- 若视觉对比不足，优先增加 `ring` 对比，而非新增 `glass-*` 变体。

## 动画实现

使用 Vue 内置 `<Transition>` + JavaScript hooks（`onBeforeEnter` / `onEnter` / `onLeave`）
实现弹簧高度动画，**不依赖任何外部动画库**：

| 阶段 | 曲线 | 时长 |
|------|------|------|
| 进入（展开） | `cubic-bezier(0.34, 1.56, 0.64, 1)`（spring 感） | 0.45s |
| 离开（收起） | `cubic-bezier(0.25, 0.46, 0.45, 0.94)`（ease-out） | 0.3s |

## 已知坑点

### backdrop-filter 在生产构建中丢失

**现象**：开发环境正常，`npm run build` 后毛玻璃变为纯透明。

**根因**：Vite 生产构建使用 Lightning CSS 进行 CSS 优化。当 CSS 中 `backdrop-filter`
写在 `-webkit-backdrop-filter` **之前**时，Lightning CSS 会将标准属性优化掉，
仅保留 `-webkit-` 版本，导致非 WebKit 浏览器失去毛玻璃效果，
部分浏览器渲染为透明。

**正确写法**（`-webkit-` 在前作为降级，标准属性在后覆盖）：

```css
.glass {
  -webkit-backdrop-filter: blur(16px) saturate(160%); /* 降级 */
  backdrop-filter: blur(16px) saturate(160%);         /* 标准，在后 */
}
```

**已修复**：`src/style.css` 内 `.glass` 已统一为此顺序（2026-02）。

## 维护建议

- 新增 Tooltip 触发场景无需修改组件，直接通过 `#content` slot 扩展即可。
- 若需改变浮层背景，统一继续使用 `.glass`，勿新增 `glass-*` 变体，也勿直接写内联 `backdrop-filter`。
- 组件不依赖 `motion` 包，可安全卸载（`npm uninstall motion`）。
- 禁止在浮层外层容器上叠加不透明 `bg-*` 类，会覆盖 `.glass` 的半透明背景。
