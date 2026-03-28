import type { StreamEvent, ReActStep } from '@/types/agent-chat'

let _idCounter = 0
function uid(): string {
  return `step_${Date.now()}_${++_idCounter}`
}

function delay(ms: number): Promise<void> {
  return new Promise(r => setTimeout(r, ms))
}

function splitText(text: string, chunkSize = 2): string[] {
  const chunks: string[] = []
  for (let i = 0; i < text.length; i += chunkSize) {
    chunks.push(text.slice(i, i + chunkSize))
  }
  return chunks
}

// ============================================================
// Mock 预设回答
// ============================================================

interface MockReply {
  keywords: string[]
  steps: Array<{
    type: ReActStep['type']
    content: string
    toolCall?: { name: string; args: Record<string, unknown> }
  }>
  finalAnswer: string
}

const MOCK_REPLIES: MockReply[] = [
  {
    keywords: ['天气', '温度', '今天', '合适', '出钓'],
    steps: [
      { type: 'thought', content: '用户想知道今天是否适合出钓。我先给出天气，再结合平台监测服务返回的代表钓场环境参数给出建议。' },
      { type: 'tool_call', content: '调用天气查询工具，获取主要垂钓城市天气（哈尔滨 / 北京）...', toolCall: { name: 'query_weather', args: { cities: ['哈尔滨', '北京'], detail: true } } },
      { type: 'tool_result', content: '☀️ 今天天气概览：\n• 哈尔滨：晴转多云，3°C，湿度 63%，东北风 2级，气压 1014hPa\n• 北京：多云，15°C，湿度 55%，西北风 2-3级，气压 1016hPa\n降水概率整体较低（<15%），适合户外垂钓。' },
      { type: 'thought', content: '天气可出钓。我再核对监测服务中的代表性钓场（松花江北岸、镜泊湖大坝、密云水库大坝）的水温与溶氧，确认鱼口活跃度。' },
      { type: 'tool_call', content: '查询代表钓场实时环境参数...', toolCall: { name: 'query_environment', args: { spot_ids: [1, 7, 13], fields: ['water_temp', 'air_temp', 'humidity', 'dissolved_oxygen'] } } },
      { type: 'tool_result', content: '📊 关键钓场环境数据：\n• 松花江北岸钓场(ID:1)：水温 8.5°C，气温 3.2°C，湿度 65.3%，溶氧 8.2mg/L\n• 镜泊湖大坝钓场(ID:7)：水温 6.8°C，气温 1.5°C，湿度 70.2%，溶氧 9.0mg/L\n• 密云水库大坝钓场(ID:13)：水温 12.5°C，气温 15.3°C，湿度 55.2%，溶氧 8.6mg/L' },
      { type: 'thought', content: '三个钓场溶氧均高于 8.0mg/L，活性条件良好。若在华北，密云水库当前温度更友好；若在黑龙江，松花江北岸和镜泊湖均可出钓。' },
    ],
    finalAnswer: '## 🎣 今日出钓建议\n\n综合天气与钓场监测参数，**今天适合出钓**。\n\n### 🌤️ 天气结论\n- 哈尔滨与北京均为低降水窗口\n- 风力 2-3 级，整体适合抛投与台钓\n\n### 🏞️ 推荐钓场\n1. **密云水库大坝钓场（ID:13）**\n   - 水温 12.5°C，溶氧 8.6mg/L\n   - 对鲤鱼、翘嘴活性更友好\n2. **松花江北岸钓场（ID:1）**\n   - 水温 8.5°C，溶氧 8.2mg/L\n   - 适合鲫鱼慢钓，红虫/蚯蚓表现稳定\n3. **镜泊湖大坝钓场（ID:7）**\n   - 水温 6.8°C，溶氧 9.0mg/L\n   - 适合长竿守深水和大体型目标鱼\n\n### 💡 出钓提示\n- 上午 6:00-9:00 与傍晚 16:00-18:00 通常鱼口更稳\n- 低温水域优先活饵，温暖水域可加玉米粒打窝\n\n祝你今天顺利上鱼！'
  },
  {
    keywords: ['推荐', '去哪', '哪里', '钓点', '水域'],
    steps: [
      { type: 'thought', content: '用户需要钓点推荐。我先按后端返回的设备在线与钓场开放状态做实时热度排序。' },
      { type: 'tool_call', content: '读取热门钓场（基于实时设备 fishing_count + 容量）...', toolCall: { name: 'query_popular_spots', args: { source: 'realtime_service', limit: 5, include_capacity: true } } },
      { type: 'tool_result', content: '🏆 热门钓场：\n1. 密云水库大坝钓场(ID:13) - 当前 18 人 / 容量 100\n2. 兴凯湖北岸垂钓基地(ID:9) - 当前 15 人 / 容量 120\n3. 岗南水库坝下钓场(ID:16) - 当前 14 人 / 容量 90\n4. 镜泊湖大坝钓场(ID:7) - 当前 12 人 / 容量 100\n5. 伊犁河伊宁段钓场(ID:21) - 当前 11 人 / 容量 60' },
      { type: 'thought', content: '热度有了，我再补环境可钓性指标（温度、溶氧、pH），避免只按人气推荐。' },
      { type: 'tool_call', content: '获取 Top3 钓场环境参数...', toolCall: { name: 'query_spot_environment', args: { spot_ids: [13, 9, 16], metrics: ['water_temp', 'dissolved_oxygen', 'ph', 'status'] } } },
      { type: 'tool_result', content: '📡 Top3 环境参数：\n【密云水库大坝】水温 12.5°C | 溶氧 8.6mg/L | pH 7.3 | 状态 open ✅\n【兴凯湖北岸】水温 5.5°C | 溶氧 8.8mg/L | pH 7.2 | 状态 open ✅\n【岗南水库坝下】水温 14.2°C | 溶氧 8.3mg/L | pH 7.4 | 状态 open ✅' },
      { type: 'thought', content: '三处都可钓。密云和岗南温度更舒适，适合大众；兴凯湖偏冷但目标鱼质量高，适合进阶钓手。' },
    ],
    finalAnswer: '## 🎯 钓点推荐\n\n结合实时热度与环境参数，推荐如下：\n\n### 🥇 密云水库大坝钓场（ID:13）\n| 指标 | 数值 | 说明 |\n|------|------|------|\n| 当前人数 | 18 / 100 | 热度高但不拥挤 |\n| 水温 | 12.5°C | 适合春季鲤鱼活跃 |\n| 溶氧 | 8.6 mg/L | 鱼活性良好 |\n| 状态 | open | 可正常出钓 |\n\n### 🥈 岗南水库坝下钓场（ID:16）\n- 当前人数 14 / 90\n- 水温 14.2°C，溶氧 8.3mg/L\n- 适合草鱼、鲤鱼目标鱼种\n\n### 🥉 兴凯湖北岸垂钓基地（ID:9）\n- 当前人数 15 / 120，空间充足\n- 水温 5.5°C，建议长竿+活饵\n- 适合追求大体型目标鱼\n\n### 📌 注意\n- `龙凤湿地外围钓场(ID:6)` 当前为 `closed`\n- `扎龙湿地观鸟钓场(ID:10)` 当前为 `maintenance`，不建议前往'
  },
  {
    keywords: ['鱼种', '技巧', '钓法', '饵料', '装备', '怎么钓'],
    steps: [
      { type: 'thought', content: '用户问钓法。我按平台近期鱼情与建议文本给出更贴近实况的数据化技巧。' },
      { type: 'tool_call', content: '读取近期鱼情与建议关键词（fishing_suggestions + reminders）...', toolCall: { name: 'query_fish_activity', args: { source: 'realtime_service', include_bait_recommendation: true } } },
      { type: 'tool_result', content: '🐟 近期鱼情汇总：\n1. 鲤鱼（密云/岗南/镜泊湖）- 活跃度 高\n2. 鲫鱼（松花江/太阳岛/岗南苇塘）- 活跃度 中高\n3. 高白鲑 / 虹鳟（赛里木湖）- 活跃度 中高\n4. 大白鱼（兴凯湖）- 活跃度 中' },
      { type: 'thought', content: '继续提取不同水域下的实战配置，避免给通用模板。' },
      { type: 'tool_call', content: '生成分水域钓法建议...', toolCall: { name: 'query_fishing_guide', args: { spots: ['松花江北岸', '密云水库大坝', '赛里木湖南岸'], include_rig_setup: true } } },
      { type: 'tool_result', content: '📖 分场景钓法：\n【松花江北岸钓场(ID:1)】鲫鱼为主：红虫/蚯蚓，主线1.5+子线0.8，袖钩4号，晨钓更稳\n【密云水库大坝钓场(ID:13)】鲤鱼活跃：玉米粒打窝，主线2.0+子线1.0，伊势尼6号，4.5米手竿\n【赛里木湖南岸钓场(ID:19)】高白鲑：路亚银色勺形亮片5g，清晨窗口期更好' },
      { type: 'thought', content: '可以按新手/进阶两档整理，兼顾常见淡水钓和冷水路亚。' },
    ],
    finalAnswer: '## 🎣 垂钓技巧指南\n\n### 当前推荐目标鱼\n- **鲤鱼**：密云水库、岗南水库活跃\n- **鲫鱼**：松花江、太阳岛更稳定\n- **高白鲑/虹鳟**：赛里木湖冷水路亚优势明显\n\n### 🐟 新手优先：松花江北岸钓场（ID:1）\n| 项目 | 建议 |\n|------|------|\n| 饵料 | 红虫 > 蚯蚓 |\n| 线组 | 主线 1.5 + 子线 0.8 |\n| 钓钩 | 袖钩 4 号 |\n| 时段 | 早 6:00-8:00 |\n\n### 🐟 进阶：密云水库大坝钓场（ID:13）\n- 玉米粒打窝，鲤鱼口更稳定\n- 主线 2.0 + 子线 1.0，伊势尼 6 号\n- 4.5m 手竿或抛竿均可\n\n### 🐟 路亚：赛里木湖南岸钓场（ID:19）\n- 银色勺形亮片 5g（监测与提醒数据已验证）\n- 清晨窗口期优先，风小时更易中鱼\n\n### 💡 通用建议\n1. 先看钓场状态：open 才建议前往\n2. 低温水域优先活饵，温暖水域可加颗粒窝料\n3. 每 10-15 分钟检查饵团状态，保持雾化与附钩平衡'
  },
  {
    keywords: ['水质', '溶氧', 'pH', '环境', '检测'],
    steps: [
      { type: 'thought', content: '用户关心水质，我优先读取后端监测服务中有水下设备的钓场，并补上关闭/维护场地提示。' },
      { type: 'tool_call', content: '查询监测设备水质数据（water_quality_data + environment_data）...', toolCall: { name: 'query_water_quality', args: { spot_ids: [1, 7, 13, 16, 19], include_alerts: true } } },
      { type: 'tool_result', content: '💧 水质监测：\n【松花江北岸(ID:1)】pH 7.2 | 溶氧 8.2mg/L | 浊度 12NTU ✅\n【镜泊湖大坝(ID:7)】pH 7.4 | 溶氧 9.0mg/L | 浊度 6NTU ✅\n【密云水库大坝(ID:13)】pH 7.3 | 溶氧 8.6mg/L | 浊度 9NTU ✅\n【岗南水库坝下(ID:16)】pH 7.4 | 溶氧 8.3mg/L | 浊度 11NTU ✅\n【赛里木湖南岸(ID:19)】pH 7.5 | 溶氧 9.2mg/L | 浊度 5NTU ✅\n\n附加状态：龙凤湿地外围钓场(ID:6) closed；扎龙湿地观鸟钓场(ID:10) maintenance。' },
      { type: 'thought', content: '当前数据整体健康，我再给一个溶氧趋势样例，帮助判断短时波动是否影响鱼口。' },
      { type: 'tool_call', content: '查询松花江北岸近 24h 溶氧趋势...', toolCall: { name: 'query_trend', args: { spot_id: 1, metric: 'dissolved_oxygen', period: '24h' } } },
      { type: 'tool_result', content: '📈 松花江北岸溶氧趋势(24h)：\n• 06:00  8.7 mg/L\n• 09:00  8.4 mg/L\n• 12:00  8.1 mg/L\n• 15:00  8.0 mg/L\n• 现在    8.2 mg/L\n趋势：白天下降、傍晚回升，属于正常波动区间。' },
      { type: 'thought', content: '所有已监测点当前都在可钓范围内，主要风险来自场地状态而非水质阈值。' },
    ],
    finalAnswer: '## 💧 水质环境监测报告\n\n### 📊 重点钓场水质\n| 钓场 | pH | 溶氧(mg/L) | 浊度(NTU) | 结论 |\n|------|----|------------|-----------|------|\n| 松花江北岸(ID:1) | 7.2 | 8.2 | 12 | 🟢 可钓 |\n| 镜泊湖大坝(ID:7) | 7.4 | 9.0 | 6 | 🟢 可钓 |\n| 密云水库大坝(ID:13) | 7.3 | 8.6 | 9 | 🟢 可钓 |\n| 岗南水库坝下(ID:16) | 7.4 | 8.3 | 11 | 🟢 可钓 |\n| 赛里木湖南岸(ID:19) | 7.5 | 9.2 | 5 | 🟢 可钓 |\n\n### ⚠️ 状态提醒（非水质）\n- `龙凤湿地外围钓场(ID:6)`：`closed`\n- `扎龙湿地观鸟钓场(ID:10)`：`maintenance`\n\n### 🧭 建议\n- 当前更应优先看钓场状态，其次看溶氧与温度\n- 所有 open 且有监测数据的点位，现阶段都在安全区间'
  },
  {
    keywords: [],
    steps: [
      { type: 'thought', content: '用户提的是泛问题。我先从平台返回的数据范围里找可回答的信息。' },
      { type: 'tool_call', content: '检索平台知识与监测索引...', toolCall: { name: 'search_knowledge', args: { source: 'backend_knowledge', top_k: 3 } } },
      { type: 'tool_result', content: '📚 可用数据范围：\n1. 钓场与区域：共 21 个钓场，覆盖黑龙江/北京/石家庄/伊犁\n2. 实时环境：设备温度、湿度、气压、水质（部分点位）\n3. 状态与提醒：open/closed/maintenance 与安全提醒' },
      { type: 'thought', content: '可基于这些数据回答天气适钓、钓点推荐、水质分析、鱼种技巧等问题。' },
    ],
    finalAnswer: '## 💬 智能助手说明\n\n当前回答基于平台数据服务返回的信息进行分析。\n\n我可以帮你：\n- 🌤️ 查询适钓天气（示例城市：哈尔滨、北京）\n- 🎯 推荐钓点（基于钓场状态与设备数据）\n- 🐟 提供钓法建议（按钓场/目标鱼）\n- 💧 做水质解读（pH、溶氧、浊度）\n\n### 💡 可以直接问\n- "今天适合去密云水库吗？"\n- "松花江北岸现在水质怎么样？"\n- "赛里木湖路亚怎么配饵？"\n- "哪些钓场当前是关闭状态？"'
  },
]

/**
 * 流式 Mock 回复生成器
 *
 * 节奏说明：
 * - 思考步骤：逐字 50ms/chunk，流式输出
 * - 工具调用描述：逐字 40ms/chunk
 * - 工具执行等待：2 秒模拟耗时
 * - 观察结果：逐字 35ms/chunk，流式输出
 * - 步骤间暂停：1.2 秒
 * - 最终回答：逐字 25ms/chunk
 */
export function createMockStreamReply(
  userInput: string,
  onEvent: (event: StreamEvent) => void,
): { cancel: () => void } {
  let cancelled = false

  const reply = MOCK_REPLIES.find(r =>
    r.keywords.length === 0 || r.keywords.some(k => userInput.includes(k))
  ) ?? MOCK_REPLIES[MOCK_REPLIES.length - 1]

  const thoughtCharDelay = 52
  const toolDescCharDelay = 48
  const toolResultCharDelay = 44
  const finalCharDelay = 34
  const stepPause = 500
  const toolCallDuration = 1650

  function nowMs(): number {
    return typeof performance !== 'undefined' ? performance.now() : Date.now()
  }

  async function waitDuration(ms: number): Promise<boolean> {
    const endAt = nowMs() + ms

    while (!cancelled) {
      const remaining = endAt - nowMs()
      if (remaining <= 0) return true
      await delay(Math.min(remaining, 50))
    }

    return false
  }

  async function streamTextWithTimeCatchUp(
    text: string,
    chunkSize: number,
    perChunkDelay: number,
    emitChunk: (chunk: string) => void,
  ): Promise<boolean> {
    const chunks = splitText(text, chunkSize)
    if (chunks.length === 0) return true

    let index = 0
    let nextAt = nowMs() + perChunkDelay

    while (index < chunks.length) {
      if (cancelled) return false

      const now = nowMs()
      if (now < nextAt) {
        await delay(Math.min(nextAt - now, 50))
        continue
      }

      // 切屏后定时器会被节流；恢复时按经过时间补发多个 chunk，避免“停住”错觉。
      const dueCount = Math.max(1, Math.floor((now - nextAt) / perChunkDelay) + 1)
      const batchCount = Math.min(dueCount, chunks.length - index, 24)

      for (let i = 0; i < batchCount; i++) {
        emitChunk(chunks[index])
        index += 1
      }

      nextAt += batchCount * perChunkDelay

      if (batchCount > 8) {
        await delay(0)
      }
    }

    return true
  }

  async function run() {
    for (const step of reply.steps) {
      if (cancelled) return

      const stepId = uid()

      if (step.type === 'thought') {
        onEvent({ type: 'step_start', stepId, stepType: 'thought', content: '' })
        if (!(await streamTextWithTimeCatchUp(step.content, 2, thoughtCharDelay, (chunk) => {
          onEvent({ type: 'step_delta', stepId, delta: chunk })
        }))) return
        onEvent({ type: 'step_end', stepId })
        if (!(await waitDuration(stepPause))) return
      } else if (step.type === 'tool_call') {
        // 工具调用：先显示调用状态，流式输出描述，再等待执行
        onEvent({ type: 'step_start', stepId, stepType: 'tool_call', content: '' })
        onEvent({
          type: 'tool_start',
          stepId,
          toolCall: {
            name: step.toolCall?.name ?? 'unknown',
            args: step.toolCall?.args ?? {},
            status: 'calling',
          },
        })
        if (!(await streamTextWithTimeCatchUp(step.content, 2, toolDescCharDelay, (chunk) => {
          onEvent({ type: 'step_delta', stepId, delta: chunk })
        }))) return
        // 模拟工具执行等待
        if (!(await waitDuration(toolCallDuration))) return
        onEvent({
          type: 'tool_result',
          stepId,
          toolCall: {
            name: step.toolCall?.name ?? 'unknown',
            args: step.toolCall?.args ?? {},
            status: 'success',
          },
        })
        onEvent({ type: 'step_end', stepId })
        if (!(await waitDuration(stepPause))) return
      } else if (step.type === 'tool_result') {
        onEvent({ type: 'step_start', stepId, stepType: 'tool_result', content: '' })
        if (!(await streamTextWithTimeCatchUp(step.content, 3, toolResultCharDelay, (chunk) => {
          onEvent({ type: 'step_delta', stepId, delta: chunk })
        }))) return
        onEvent({ type: 'step_end', stepId })
        if (!(await waitDuration(stepPause))) return
      }
    }

    if (cancelled) return

    // 最终回答
    onEvent({ type: 'final_start', content: '' })
    if (!(await streamTextWithTimeCatchUp(reply.finalAnswer, 2, finalCharDelay, (chunk) => {
      onEvent({ type: 'final_delta', delta: chunk })
    }))) return
    onEvent({ type: 'final_end' })
    onEvent({ type: 'done' })
  }

  run()

  return {
    cancel() {
      cancelled = true
    },
  }
}
