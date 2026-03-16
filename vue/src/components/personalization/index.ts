/**
 * components/personalization/index.ts - 个性化模块统一导出
 *
 * 职责：汇聚导出所有个性化组件和类型，供外部消费方一站式引入
 * 对外暴露：所有区块组件、类型、默认配置工厂
 */

// 类型
export type * from './types'
export {
  isNavbarConfig,
  isHeroConfig,
  isIntroductionConfig,
  isThemeListConfig,
  isCtaConfig,
  isFooterConfig,
  SECTION_TYPE_META,
} from './types'

// 默认配置
export {
  defaultGlobal,
  defaultNavbarConfig,
  defaultHeroConfig,
  defaultIntroductionConfig,
  defaultThemeListConfig,
  defaultCtaConfig,
  defaultFooterConfig,
  createSection,
  defaultHomepageLayout,
} from './defaults'

// 区块组件（懒加载，由 HomepageRenderer 按需引入）
export { default as HomepageRenderer } from './renderer/HomepageRenderer.vue'
