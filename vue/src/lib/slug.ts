/**
 * lib/slug.ts - 标题转 URL 友好 slug 工具
 * 职责：将中/英文标题转为紧凑 slug，优先英文映射，否则拼音转写，限制最大长度
 * 对外暴露：titleToSlug(title, maxLen?)
 */
import { pinyin } from 'pinyin-pro'

/** 自动生成 slug 的默认最大长度 */
const MAX_SLUG_LENGTH = 10

/** 常见中文 → 英文 slug 映射（值长度须 ≤ MAX_SLUG_LENGTH） */
const WORD_MAP: Record<string, string> = {
  '快速开始': 'quickstart',
  '快速入门': 'quickstart',
  '入门指南': 'getstarted',
  '入门': 'start',
  '安装': 'install',
  '安装指南': 'install',
  '配置': 'config',
  '介绍': 'intro',
  '简介': 'intro',
  '概述': 'overview',
  '概览': 'overview',
  '指南': 'guide',
  '教程': 'tutorial',
  '参考': 'ref',
  '参考手册': 'reference',
  '常见问题': 'faq',
  '更新日志': 'changelog',
  '贡献': 'contrib',
  '部署': 'deploy',
  '开发': 'dev',
  '测试': 'test',
  '文档': 'docs',
  '首页': 'home',
  '关于': 'about',
  '联系': 'contact',
  '帮助': 'help',
  '设置': 'settings',
  '搜索': 'search',
  '登录': 'login',
  '注册': 'signup',
  '用户': 'user',
  '管理': 'admin',
  '接口': 'api',
  '权限': 'perm',
  '认证': 'auth',
  '授权': 'auth',
  '主题': 'theme',
  '版本': 'version',
  '章节': 'section',
  '页面': 'page',
  '标签': 'tag',
  '分类': 'category',
  '目录': 'catalog',
  '数据': 'data',
  '数据库': 'database',
  '服务': 'service',
  '工具': 'utils',
  '示例': 'example',
  '样式': 'style',
  '模板': 'template',
  '组件': 'component',
  '路由': 'router',
  '中间件': 'middleware',
  '发布': 'release',
  '迁移': 'migrate',
  '备份': 'backup',
  '监控': 'monitor',
  '日志': 'log',
  '错误': 'error',
  '通知': 'notify',
  '消息': 'message',
  '上传': 'upload',
  '下载': 'download',
}

/**
 * 将标题转为 URL 友好的 slug
 * 1. 优先匹配中文→英文映射表
 * 2. 否则使用 pinyin-pro 转拼音，按音节边界截断
 * 3. 结果仅含 a-z0-9，不超过 maxLen 字符
 */
export function titleToSlug(title: string, maxLen = MAX_SLUG_LENGTH): string {
  const source = title.trim()
  if (!source) return ''

  // 完整词映射
  const mapped = WORD_MAP[source]
  if (mapped) return mapped.slice(0, maxLen)

  // 中文转拼音（无声调，空格分隔各字拼音）
  const pinyinStr = pinyin(source, { toneType: 'none' })

  // 拆分为音节段
  const segments = pinyinStr
    .toLowerCase()
    .replace(/['']/g, '')
    .replace(/[^a-z0-9\s]/g, '')
    .trim()
    .split(/\s+/)
    .filter(Boolean)

  if (!segments.length) {
    return source.toLowerCase().replace(/[^a-z0-9]/g, '').slice(0, maxLen)
  }

  // 逐段拼接，不超过 maxLen
  let result = ''
  for (const seg of segments) {
    if (result.length + seg.length > maxLen) break
    result += seg
  }

  return result || segments[0]!.slice(0, maxLen)
}
