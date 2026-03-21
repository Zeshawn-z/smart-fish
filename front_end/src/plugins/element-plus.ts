import type { App } from 'vue'
import zhCn from 'element-plus/es/locale/lang/zh-cn'
import { ElConfigProvider } from 'element-plus'

// Element Plus 按需导入说明：
// - 组件 + 样式由 vite.config.ts 中的 unplugin-vue-components + ElementPlusResolver 自动按需导入
// - 不再导入全量 CSS（element-plus/dist/index.css），大幅减小打包体积
// - 图标由各组件自行 import，不再全局注册全部图标
// - 此处仅注册 ElConfigProvider 用于全局 locale/size 配置

export const elLocale = zhCn

export function setupElementPlus(app: App) {
  app.component('ElConfigProvider', ElConfigProvider)
}
