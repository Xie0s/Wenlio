<!--
文档：Sonner通知使用规范
职责：说明项目中 Toast 通知的接入方式、API 用法、视觉规范与样式维护边界。
主要接口：
  - 全局注册说明（无需在页面级导入组件）
  - toast() API 用法与类型清单
  - 全局配置参数
  - 视觉规范（glass 毛玻璃、圆角、轮廓线）
  - 样式覆盖位置与维护建议
-->

# Sonner 通知使用规范

## 背景与目标

项目使用 [vue-sonner](https://vue-sonner.vercel.app/) 作为全局 Toast 通知方案。`<Toaster>` 组件已在 `App.vue` 根组件中全局注册，**业务页面与组件无需再次导入或挂载 `<Toaster>`，直接调用 `toast()` API 即可。**

---

## 全局注册说明

`<Toaster>` 已在 `src/App.vue` 中统一挂载：

```vue
<!-- src/App.vue -->
<Toaster position="bottom-right" :duration="3000" />
```

| 参数 | 值 | 说明 |
|------|-----|------|
| `position` | `bottom-right` | 通知出现位置 |
| `duration` | `3000` | 默认显示时长（ms） |

---

## 使用方式

### 1. 导入 `toast`

在任意业务文件中只需导入 `toast` 函数，**不需要导入或注册组件**：

```ts
import { toast } from 'vue-sonner'
```

### 2. 调用示例

```ts
// 默认
toast('操作完成')

// 成功
toast.success('保存成功')

// 失败
toast.error('请求失败，请重试')

// 警告
toast.warning('操作存在风险')

// 提示
toast.info('版本已是最新')

// 加载中（返回 toastId，可用于后续 dismiss/update）
const id = toast.loading('正在提交...')

// Promise（自动处理 loading → success/error 状态）
toast.promise(fetchData(), {
  loading: '加载中...',
  success: '加载完成',
  error: '加载失败',
})

// 手动关闭
toast.dismiss(id)
```

### 3. 自定义时长

```ts
toast.success('已同步', { duration: 1500 })
```

---

## 视觉规范

Toast 通知采用与全局 `.glass` 工具类一致的毛玻璃效果，通过 `src/style.css` 的高优先级选择器统一覆盖：

| 属性 | 亮色模式 | 暗色模式 |
|------|---------|---------|
| 背景 | `oklch(1 0 0 / 60%)` | `oklch(0.2 0.006 286 / 50%)` |
| 模糊 | `blur(16px) saturate(180%)` | 同左 |
| 圆角 | `1rem`（`rounded-2xl`） | 同左 |
| 轮廓线 | `1px solid oklch(0 0 0 / 8%)` | `1px solid oklch(1 0 0 / 8%)` |
| 阴影 | 无 | 无 |

---

## 样式覆盖位置

所有 Sonner 视觉覆盖集中在 `src/style.css` 底部的 **"Sonner Toast 毛玻璃样式"** 区块：

```css
/* src/style.css */
:root [data-sonner-toaster] [data-sonner-toast] { ... }
:root.dark [data-sonner-toaster] [data-sonner-toast] { ... }
```

> **注意**：选择器前缀必须保留 `:root`（或 `:root.dark`）以提升特异性，确保覆盖 vue-sonner 运行时动态注入的样式。

---

## 维护建议

- 不在页面或组件中单独覆盖 Toast 样式，统一在 `src/style.css` 的专属区块维护。
- 升级 `vue-sonner` 版本后，需回归验证毛玻璃效果与圆角是否仍正常显示。
- 不修改 `src/components/ui/sonner/Sonner.vue` 中的 `--border-radius` 内联变量，该值由组件负责注入，CSS 层独立覆盖。
