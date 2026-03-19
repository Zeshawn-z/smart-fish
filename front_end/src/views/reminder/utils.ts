/** 格式化时间为 MM/DD HH:mm 格式 */
export function formatTime(ts: string): string {
  const d = new Date(ts)
  return d.toLocaleString('zh-CN', { month: '2-digit', day: '2-digit', hour: '2-digit', minute: '2-digit' })
}

/** 根据分数返回样式类名 */
export function getScoreClass(score: number): string {
  if (score >= 80) return 'high'
  if (score >= 50) return 'mid'
  return 'low'
}
