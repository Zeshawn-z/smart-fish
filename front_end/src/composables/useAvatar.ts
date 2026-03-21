/**
 * 头像服务 composable
 *
 * 统一处理用户头像逻辑：
 * - 有真实头像 URL → 直接使用
 * - 无头像 → 根据用户名/ID 生成确定性的彩色背景 + 首字母
 *
 * 不同用户生成的背景颜色不同（基于用户名 hash）
 * 可在 Navbar、个人中心、社区帖子、评论等处统一复用
 */

/** 预设的渐变色方案 — 选了 16 组好看的 */
const AVATAR_GRADIENTS = [
  ['#667eea', '#764ba2'],
  ['#f093fb', '#f5576c'],
  ['#4facfe', '#00f2fe'],
  ['#43e97b', '#38f9d7'],
  ['#fa709a', '#fee140'],
  ['#a18cd1', '#fbc2eb'],
  ['#fccb90', '#d57eeb'],
  ['#e0c3fc', '#8ec5fc'],
  ['#f5576c', '#ff9068'],
  ['#48c6ef', '#6f86d6'],
  ['#feada6', '#f5efef'],
  ['#a1c4fd', '#c2e9fb'],
  ['#d4fc79', '#96e6a1'],
  ['#84fab0', '#8fd3f4'],
  ['#cfd9df', '#e2ebf0'],
  ['#ff9a9e', '#fecfef'],
]

/** 简单字符串 hash（确定性） */
function simpleHash(str: string): number {
  let hash = 0
  for (let i = 0; i < str.length; i++) {
    hash = ((hash << 5) - hash + str.charCodeAt(i)) | 0
  }
  return Math.abs(hash)
}

/** 根据用户名获取渐变色组 */
export function getAvatarGradient(name: string): [string, string] {
  const idx = simpleHash(name) % AVATAR_GRADIENTS.length
  return AVATAR_GRADIENTS[idx] as [string, string]
}

/** 根据用户名获取首字母（支持中英文） */
export function getAvatarLetter(name: string): string {
  if (!name) return '?'
  // 取第一个非空白字符
  const ch = name.trim().charAt(0)
  // 英文字母转大写
  if (/[a-zA-Z]/.test(ch)) return ch.toUpperCase()
  return ch
}

/** 生成 CSS gradient 背景字符串 */
export function getAvatarBackground(name: string): string {
  const [c1, c2] = getAvatarGradient(name)
  return `linear-gradient(135deg, ${c1}, ${c2})`
}

/** 生成内联 style 对象（用于无头像的 el-avatar / div） */
export function getAvatarStyle(name: string): Record<string, string> {
  return {
    background: getAvatarBackground(name),
    color: '#fff',
    fontWeight: '600',
  }
}

/**
 * 综合头像信息
 * @param avatarUrl 用户的真实头像 URL（可能为空）
 * @param username 用户名
 * @returns { src, letter, style, hasAvatar }
 */
export function useAvatar(avatarUrl?: string | null, username?: string) {
  const name = username || '用户'
  const hasAvatar = !!avatarUrl

  return {
    /** 真实头像 URL（可能 undefined） */
    src: avatarUrl || undefined,
    /** 无头像时显示的首字母 */
    letter: getAvatarLetter(name),
    /** 无头像时的背景样式 */
    style: getAvatarStyle(name),
    /** 是否有真实头像 */
    hasAvatar,
  }
}

/**
 * 生成 Data URL 头像图片（Canvas 方式，用于需要真实图片 URL 的场景）
 * @param name 用户名
 * @param size 图片大小 (px)
 */
export function generateAvatarDataURL(name: string, size = 80): string {
  if (typeof document === 'undefined') return ''

  const canvas = document.createElement('canvas')
  canvas.width = size
  canvas.height = size
  const ctx = canvas.getContext('2d')
  if (!ctx) return ''

  // 画渐变背景
  const [c1, c2] = getAvatarGradient(name)
  const gradient = ctx.createLinearGradient(0, 0, size, size)
  gradient.addColorStop(0, c1)
  gradient.addColorStop(1, c2)
  ctx.fillStyle = gradient
  ctx.fillRect(0, 0, size, size)

  // 画首字母
  const letter = getAvatarLetter(name)
  ctx.fillStyle = '#fff'
  ctx.font = `bold ${size * 0.45}px -apple-system, "Helvetica Neue", sans-serif`
  ctx.textAlign = 'center'
  ctx.textBaseline = 'middle'
  ctx.fillText(letter, size / 2, size / 2 + size * 0.03)

  return canvas.toDataURL('image/png')
}
