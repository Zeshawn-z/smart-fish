/**
 * ECharts 按需导入统一入口
 *
 * 全量 import * as echarts from 'echarts' 大约 1.1MB（gzip ~365KB）
 * 按需导入后只引入实际用到的图表类型和组件，大幅减小体积
 */

import * as echarts from 'echarts/core'

// 渲染器
import { CanvasRenderer } from 'echarts/renderers'

// 图表类型
import {
  LineChart,
  PieChart,
  BarChart,
  RadarChart,
  ScatterChart,
  EffectScatterChart,
  MapChart
} from 'echarts/charts'

// 组件
import {
  TooltipComponent,
  LegendComponent,
  GridComponent,
  GeoComponent,
  RadarComponent,
  GraphicComponent,
  VisualMapComponent
} from 'echarts/components'

// 注册
echarts.use([
  CanvasRenderer,
  // Charts
  LineChart,
  PieChart,
  BarChart,
  RadarChart,
  ScatterChart,
  EffectScatterChart,
  MapChart,
  // Components
  TooltipComponent,
  LegendComponent,
  GridComponent,
  GeoComponent,
  RadarComponent,
  GraphicComponent,
  VisualMapComponent
])

export default echarts
export type { ECharts } from 'echarts/core'
export type { EChartsOption } from 'echarts/types/dist/shared'
