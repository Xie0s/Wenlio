/**
 * components/personalization/types.ts - 租户个性化首页类型定义
 *
 * 职责：定义首页配置的完整类型系统，包括各区块配置、布局选项、动画效果等。
 *       作为前端个性化模块的类型唯一来源，后端存储结构与此对齐。
 * 对外暴露：所有个性化相关的 interface / type / 常量
 */

// ============================================================
// 顶层配置结构
// ============================================================

/** 租户首页完整配置（后端返回 / 前端提交） */
export interface TenantHomepageConfig {
  /** 已发布配置（读者可见） */
  published: HomepageLayout | null
  /** 草稿配置（管理员编辑中） */
  draft: HomepageLayout | null
  updated_at: string
}

/** 首页布局（一套完整的页面描述） */
export interface HomepageLayout {
  global: HomepageGlobal
  sections: HomepageSection[]
}

/** 全局样式设定 */
export interface HomepageGlobal {
  /** 页面背景色（亮色模式） */
  background_color: string
  /** 页面背景色（暗色模式），空值时继承系统默认 */
  dark_background_color: string
  /** 浏览器标签标题，空值时读者端回退到租户名称 */
  browser_title: string
  /** 浏览器标签图标 URL，空值时回退到默认 favicon */
  browser_icon_url: string
  /** 字体族，空值时继承系统默认 */
  font_family: string
  /** 内容区最大宽度，如 '1200px' */
  max_width: string
  /** 区块间距（上下 padding），单位 px，如 96 */
  section_spacing: number
}

/** 标题对齐方式 */
export type TitleAlign = 'left' | 'center' | 'right'

// ============================================================
// 区块通用结构
// ============================================================

/** 区块类型枚举 */
export type SectionType = 'navbar' | 'hero' | 'introduction' | 'theme_list' | 'cta' | 'footer'

/** 区块配置联合类型 */
export type SectionConfig =
  | NavbarConfig
  | HeroConfig
  | IntroductionConfig
  | ThemeListConfig
  | CtaConfig
  | FooterConfig

/** 单个区块描述 */
export interface HomepageSection {
  /** 唯一标识（UUID） */
  id: string
  /** 区块类型 */
  type: SectionType
  /** 是否对读者可见 */
  visible: boolean
  /** 区块具体配置，类型由 type 字段决定 */
  config: SectionConfig
}

// ============================================================
// Navbar — 导航栏
// ============================================================

/** 导航栏样式 */
export type NavbarStyle = 'transparent' | 'solid' | 'blur'

/** 导航链接 */
export interface NavLink {
  label: string
  url: string
  /** 二级链接 */
  children?: NavLink[]
  /** 是否在新标签页打开 */
  external?: boolean
}

/** 品牌区域显示模式 */
export type NavbarBrandMode = 'both' | 'logo_only' | 'text_only'

/** 导航栏配置 */
export interface NavbarConfig {
  /** 覆盖租户 Logo（空值时使用租户默认 Logo） */
  logo_url: string
  /** 品牌文字（空值时使用租户名称） */
  brand_text: string
  /** 品牌区域显示模式：图文都显示 / 仅 Logo / 仅文字 */
  brand_mode: NavbarBrandMode
  /** 导航链接列表 */
  links: NavLink[]
  /** 导航栏视觉风格 */
  style: NavbarStyle
  /** 是否吸顶 */
  sticky: boolean
  /** 右侧 CTA 按钮（可选） */
  cta_button: ButtonConfig | null
  /** 是否在 CTA 按钮后显示主题切换开关 */
  show_theme_toggle?: boolean
}

// ============================================================
// Hero — 首屏主视觉
// ============================================================

/** Hero 布局模式 */
export type HeroLayout = 'centered' | 'left-right' | 'right-left'

/** Hero 入场动画 */
export type HeroAnimation = 'none' | 'fade-up' | 'fade-in' | 'typewriter'

/** 按钮配置（Hero / CTA / Navbar 共用） */
export interface ButtonConfig {
  text: string
  url: string
  /** noise = NoiseBackground 渐变噪点风格按钮, dark = 深色药丸按钮, plain = 无样式纯文字链接 */
  variant: 'primary' | 'secondary' | 'outline' | 'noise' | 'dark' | 'plain'
  /** 是否在按钮右侧显示箭头 → */
  show_arrow: boolean
}

/** Hero 区块配置 */
export interface HeroConfig {
  /** 布局模式 */
  layout: HeroLayout
  /** 标题对齐方式 */
  title_align: TitleAlign
  /** 主标题 */
  title: string
  /** 副标题 */
  subtitle: string
  /** 描述文本 */
  description: string
  /** 主按钮 */
  primary_button: ButtonConfig | null
  /** 次按钮 */
  secondary_button: ButtonConfig | null
  /** 入场动画 */
  animation: HeroAnimation
  /** 配图（左右分栏布局时展示） */
  image_url: string
  /** 艺术画风背景图（居中布局时覆盖整个区域） */
  background_image_url: string
  /** 背景图遮罩颜色（亮色模式），如 'rgba(245,240,232,0.82)' */
  background_overlay: string
  /** 背景图遮罩颜色（暗色模式） */
  background_overlay_dark: string
  /** 标题中要以灰色显示的文字片段（Cursor 混合粗细风格） */
  highlight_text: string
}

// ============================================================
// Introduction — 特性介绍
// ============================================================

/** 介绍区布局 */
export type IntroductionLayout = 'grid-3' | 'grid-2' | 'list'

/** 介绍区卡片视觉风格 */
export type IntroductionCardStyle = 'elevated' | 'flat' | 'bordered' | 'glass'

/** 单个特性项 */
export interface FeatureItem {
  /** Lucide 图标名 */
  icon: string
  title: string
  description: string
  /** 卡片内可选截图/图片 URL */
  image_url: string
  /** 卡片底部可选链接文字，如"了解更多 ↗" */
  link_text?: string
  /** 卡片底部可选链接地址 */
  link_url?: string
}

/** 介绍区配置 */
export interface IntroductionConfig {
  /** 标题对齐方式 */
  title_align: TitleAlign
  title: string
  description: string
  /** 特性条目 */
  features: FeatureItem[]
  /** 布局模式 */
  layout: IntroductionLayout
  /** 卡片视觉风格 */
  card_style: IntroductionCardStyle
}

// ============================================================
// ThemeList — 文档主题列表
// ============================================================

/** 主题列表行风格（全部为节省空间的列表形式） */
export type ThemeCardStyle =
  | 'list-simple'    // 简洁行：图标+名称+描述+箭头，带统一边框容器与分割线
  | 'list-compact'   // 密集行：名称与描述在同一行，行高极小
  | 'list-numbered'  // 编号行：序号+名称+描述，适合有序主题
  | 'list-two-col'   // 双列网格：每项独立圆角边框，两列排列
  | 'list-tag'       // 标签云：药丸式 flex-wrap，空间利用率最高
  | 'list-headline'  // 大标题行：无容器边框，字体突出，底部分割线

/** 主题列表配置 */
export interface ThemeListConfig {
  /** 标题对齐方式 */
  title_align: TitleAlign
  title: string
  description: string
  /** 列表行风格 */
  card_style: ThemeCardStyle
  /** 是否显示主题描述 */
  show_description: boolean
  /** 是否显示主题 slug（URL 路径段） */
  show_slug: boolean
  /** 是否显示创建时间 */
  show_date: boolean
}

// ============================================================
// CTA — 行动号召
// ============================================================

/** CTA 布局模式 */
export type CtaLayout = 'simple' | 'card'

/** CTA 区块配置 */
export interface CtaConfig {
  /** 标题对齐方式 */
  title_align: TitleAlign
  title: string
  description: string
  /** 是否显示描述文字 */
  show_description: boolean
  /** 是否显示行动按钮 */
  show_button: boolean
  /** 区块背景色 */
  background_color: string
  /** 行动按钮 */
  button: ButtonConfig
  /** 艺术画风背景图 */
  background_image_url: string
  /** 布局模式：simple 居中 / card 带浮动卡片 */
  layout: CtaLayout
  /** 浮动卡片内截图 URL（layout=card 时使用） */
  card_image_url: string
  /** 浮动卡片标题 */
  card_title: string
  /** 浮动卡片描述 */
  card_description: string
  /** 浮动卡片底栏按钮文字（空值时不显示按钮） */
  card_button_text: string
  /** 浮动卡片底栏按钮图标（Lucide 图标 id，如 'arrow-right'） */
  card_button_icon: string
  /** 是否在毛玻璃底栏显示左侧标题/描述文字区 */
  card_footer_text_visible: boolean
}

// ============================================================
// Footer — 页脚
// ============================================================

/** 页脚自定义链接 */
export interface FooterCustomLink {
  name: string
  url: string
}

/** 页脚版权信息 */
export interface FooterCopyright {
  /** 版权所有者 */
  owner: string
  /** 版权所有者链接 */
  link: string
  /** 版权年份，如 '2020-2025' */
  year: string
}

/** 页脚备案信息 */
export interface FooterFiling {
  /** ICP 备案号 */
  icp_number: string
  /** ICP 备案链接 */
  icp_link: string
  /** 公安备案号 */
  police_number: string
  /** 公安备案链接 */
  police_link: string
}

/** 页脚配置 */
export interface FooterConfig {
  /** Logo URL（留空使用租户 Logo） */
  logo_url: string
  /** 宣传标语 */
  slogan: string
  /** 自定义链接列表 */
  custom_links: FooterCustomLink[]
  /** 版权信息 */
  copyright: FooterCopyright
  /** 网站名称（留空使用租户名称） */
  site_name: string
  /** 备案信息 */
  filing: FooterFiling
}

// ============================================================
// 类型守卫辅助（按 SectionType 安全取 config）
// ============================================================

export function isNavbarConfig(s: HomepageSection): s is HomepageSection & { config: NavbarConfig } {
  return s.type === 'navbar'
}
export function isHeroConfig(s: HomepageSection): s is HomepageSection & { config: HeroConfig } {
  return s.type === 'hero'
}
export function isIntroductionConfig(s: HomepageSection): s is HomepageSection & { config: IntroductionConfig } {
  return s.type === 'introduction'
}
export function isThemeListConfig(s: HomepageSection): s is HomepageSection & { config: ThemeListConfig } {
  return s.type === 'theme_list'
}
export function isCtaConfig(s: HomepageSection): s is HomepageSection & { config: CtaConfig } {
  return s.type === 'cta'
}
export function isFooterConfig(s: HomepageSection): s is HomepageSection & { config: FooterConfig } {
  return s.type === 'footer'
}

// ============================================================
// 常量：区块类型元数据（用于编辑器 UI）
// ============================================================

export const SECTION_TYPE_META: Record<SectionType, { label: string; icon: string; singleton: boolean }> = {
  navbar:        { label: '导航栏', icon: 'navigation', singleton: true },
  hero:          { label: '首屏主视觉', icon: 'image', singleton: true },
  introduction:  { label: '特性介绍', icon: 'layout-grid', singleton: false },
  theme_list:    { label: '文档主题', icon: 'library', singleton: true },
  cta:           { label: '行动号召', icon: 'megaphone', singleton: false },
  footer:        { label: '页脚', icon: 'panel-bottom', singleton: true },
}
