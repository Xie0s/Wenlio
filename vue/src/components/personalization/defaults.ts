/**
 * components/personalization/defaults.ts - 首页区块默认配置
 *
 * 职责：为每种区块类型提供出厂默认值，供编辑器"新增区块"和首次创建首页时使用。
 * 对外暴露：各区块默认配置工厂函数、完整默认布局生成函数
 */

import type {
  HomepageLayout,
  HomepageGlobal,
  NavbarConfig,
  HeroConfig,
  IntroductionConfig,
  ThemeListConfig,
  CtaConfig,
  FooterConfig,
  HomepageSection,
  SectionType,
  SectionConfig,
} from './types'

let _counter = 0
function uid(): string {
  return `sec_${Date.now()}_${++_counter}`
}

// ── 全局默认 ─────────────────────────────────────────────────

export function defaultGlobal(): HomepageGlobal {
  return {
    background_color: '#f5f0e8',
    dark_background_color: '#1a1a1a',
    browser_title: '',
    browser_icon_url: '',
    font_family: '',
    max_width: '1200px',
    section_spacing: 96,
  }
}

// ── 各区块默认配置 ──────────────────────────────────────────

export function defaultNavbarConfig(): NavbarConfig {
  return {
    logo_url: '',
    brand_text: 'Wenlio 文流',
    brand_mode: 'both',
    links: [
      { label: '核心能力', url: '#features' },
      { label: '组织结构', url: '#structure' },
      { label: '发布流程', url: '#workflow' },
      { label: '功能总览', url: '#more-features' },
      { label: '个性化组件', url: '#personalization' },
    ],
    style: 'blur',
    sticky: true,
    cta_button: { text: '登录', url: '/admin/login', variant: 'dark', show_arrow: false },
    show_theme_toggle: true,
  }
}

export function defaultHeroConfig(): HeroConfig {
  return {
    layout: 'centered',
    title_align: 'center',
    title: '多租户文档管理与阅读\n一个应用即刻上线',
    subtitle: '文档管理平台 · v2',
    description: '从主题与版本管理、Markdown 编辑与渲染，到全文搜索与评论审核，形成完整闭环。发布仅变更状态，无需构建等待，读者刷新即见最新内容。',
    primary_button: { text: '查看核心能力与流程', url: '#features', variant: 'outline', show_arrow: false },
    secondary_button: null,
    animation: 'fade-up',
    image_url: '',
    background_image_url: '',
    background_overlay: 'rgba(245,240,232,0.82)',
    background_overlay_dark: 'rgba(26,26,26,0.85)',
    highlight_text: '一个应用即刻上线',
  }
}

export function defaultIntroductionConfig(): IntroductionConfig {
  return {
    title_align: 'center',
    title: '围绕需求闭环的四项核心能力',
    description: '从创作、协作、发布到审核，统一在一个界面里完成，降低团队交付与治理成本。',
    features: [
      {
        icon: 'file-text',
        title: '富文本与 Markdown 双模式编辑器',
        description: '支持 WYSIWYG 与 Markdown 源码双模式切换，内容互转且渲染一致，支持代码块、表格、自定义容器、图片上传与 .md 导入。',
        image_url: '',
        link_text: '',
        link_url: '',
      },
      {
        icon: 'users',
        title: '路径、数据、权限四维完全隔离',
        description: '所有文档路径以 /{tenant_id}/ 开头，数据与搜索强制按 tenant_id 过滤，接口层阻断跨租户访问。',
        image_url: '',
        link_text: '',
        link_url: '',
      },
      {
        icon: 'zap',
        title: '单页或整版发布，改状态即生效',
        description: '发布仅更新状态，无需构建等待。支持单篇发布/下线与整版本批量发布，刷新即可获取最新已发布内容。',
        image_url: '',
        link_text: '',
        link_url: '',
      },
      {
        icon: 'lock',
        title: '三类角色协同，评论先审后显',
        description: '超级管理员、租户管理员、读者权限边界明确；评论默认 pending，审核通过后公开展示。',
        image_url: '',
        link_text: '',
        link_url: '',
      },
    ],
    layout: 'grid-3',
    card_style: 'elevated',
  }
}

export function defaultThemeListConfig(): ThemeListConfig {
  return {
    title_align: 'center',
    title: '把核心能力落到可管理的层级结构',
    description: '四层路径从租户到文档页逐级展开。选择主题后可按默认版本继续阅读。',
    card_style: 'list-simple',
    show_description: true,
    show_slug: false,
    show_date: false,
  }
}

export function defaultCtaConfig(): CtaConfig {
  return {
    title_align: 'center',
    title: '按需求链路完成导览，立即开始管理文档。',
    description: '从主题与版本管理，到文档发布、搜索与评论审核，后台能力已完整就绪。登录后即可按租户开始配置与发布。',
    show_description: true,
    show_button: true,
    background_color: '',
    button: { text: '进入后台开始使用', url: '/admin/login', variant: 'dark', show_arrow: true },
    background_image_url: '',
    layout: 'simple',
    card_image_url: '',
    card_title: '',
    card_description: '',
    card_button_text: '',
    card_button_icon: '',
    card_footer_text_visible: false,
  }
}

export function defaultFooterConfig(): FooterConfig {
  return {
    logo_url: '',
    slogan: '支持多租户隔离、版本化管理、即时发布与评论审核的一体化文档平台。',
    custom_links: [
      { name: '关于我们', url: 'https://microswift.cn' },
      { name: '后台登录', url: '/admin/login' },
      { name: '服务条款', url: '#' },
      { name: '隐私政策', url: '#' },
    ],
    copyright: {
      owner: 'Microswift Core™',
      link: 'https://microswift.cn',
      year: '2020-2026',
    },
    site_name: 'Wenlio 文流',
    filing: {
      icp_number: '陇ICP备20002844号-1',
      icp_link: 'https://beian.miit.gov.cn/',
      police_number: '甘公网安备62090202000540号',
      police_link: 'https://www.beian.gov.cn/',
    },
  }
}

// ── 区块配置工厂映射 ────────────────────────────────────────

const CONFIG_FACTORY: Record<SectionType, () => SectionConfig> = {
  navbar: defaultNavbarConfig,
  hero: defaultHeroConfig,
  introduction: defaultIntroductionConfig,
  theme_list: defaultThemeListConfig,
  cta: defaultCtaConfig,
  footer: defaultFooterConfig,
}

/** 根据区块类型创建带默认配置的 HomepageSection */
export function createSection(type: SectionType): HomepageSection {
  return {
    id: uid(),
    type,
    visible: true,
    config: CONFIG_FACTORY[type](),
  }
}

// ── 完整默认布局 ─────────────────────────────────────────────

/** 生成一套包含所有区块的默认首页布局 */
export function defaultHomepageLayout(): HomepageLayout {
  return {
    global: defaultGlobal(),
    sections: [
      createSection('navbar'),
      createSection('hero'),
      createSection('introduction'),
      createSection('theme_list'),
      createSection('cta'),
      createSection('footer'),
    ],
  }
}
