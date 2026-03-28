/**
 * 智能体对话类型定义
 * 基于 ReAct（Reasoning + Acting）范式的智能体消息模型
 */

/** ReAct 步骤类型 */
export type StepType = 'thought' | 'tool_call' | 'tool_result' | 'action' | 'final_answer'

/** 工具调用状态 */
export type ToolCallStatus = 'calling' | 'success' | 'error'

/** 工具调用信息 */
export interface ToolCall {
  /** 工具名称 */
  name: string
  /** 调用参数 */
  args: Record<string, unknown>
  /** 调用状态 */
  status: ToolCallStatus
}

/** ReAct 思维链中的一个步骤 */
export interface ReActStep {
  /** 步骤类型 */
  type: StepType
  /** 步骤唯一标识 */
  id: string
  /** 每步的内容 */
  content: string
  /** 工具调用信息（仅 tool_call 类型有） */
  toolCall?: ToolCall
  /** 步骤是否正在流式生成中 */
  streaming?: boolean
}

/** 消息角色 */
export type MessageRole = 'user' | 'assistant'

/** 智能体消息 */
export interface AgentMessage {
  /** 消息唯一标识 */
  id: string
  /** 消息角色 */
  role: MessageRole
  /** 最终回答文本 */
  content: string
  /** ReAct 思维链步骤 */
  steps: ReActStep[]
  /** 是否正在流式生成中 */
  streaming?: boolean
  /** 创建时间 */
  createdAt: Date
}

/** 对话会话 */
export interface ChatSession {
  id: string
  title: string
  messages: AgentMessage[]
  createdAt: Date
}

/** 流式事件类型（用于 SSE 模拟） */
export type StreamEventType =
  | 'step_start'       // 开始新步骤
  | 'step_delta'       // 步骤内容增量
  | 'step_end'         // 步骤结束
  | 'tool_start'       // 工具开始调用
  | 'tool_result'      // 工具返回结果
  | 'final_start'      // 开始最终回答
  | 'final_delta'      // 最终回答增量
  | 'final_end'        // 最终回答结束
  | 'done'             // 整条消息完成

/** 流式事件 */
export interface StreamEvent {
  type: StreamEventType
  /** 步骤 ID */
  stepId?: string
  /** 步骤类型 */
  stepType?: StepType
  /** 增量文本 */
  delta?: string
  /** 完整内容（用于某些事件） */
  content?: string
  /** 工具信息 */
  toolCall?: ToolCall
}
