/**
 * ECharts 配色：语义对齐 TDesign starter `src/config/color.ts`（light/dark 两套），
 * 不包含 tvision-color 等 starter 基础设施。seriesColors 供首页趋势图使用。
 */
export interface ChartColors {
  textColor: string
  placeholderColor: string
  borderColor: string
  containerColor: string
  seriesColors: string[]
}

export const LIGHT_CHART_COLORS: ChartColors = {
  textColor: 'rgba(0, 0, 0, 0.9)',
  placeholderColor: 'rgba(0, 0, 0, 0.35)',
  borderColor: '#dcdcdc',
  containerColor: '#fff',
  seriesColors: ['#0052D9', '#0594FA', '#00A870', '#EBB105', '#ED7B2F', '#E34D59'],
}

export const DARK_CHART_COLORS: ChartColors = {
  textColor: 'rgba(255, 255, 255, 0.9)',
  placeholderColor: 'rgba(255, 255, 255, 0.35)',
  borderColor: '#5e5e5e',
  containerColor: '#242424',
  seriesColors: ['#4582e6', '#59B8FF', '#3CC9A4', '#F2C34E', '#F0954F', '#EE6B73'],
}

export const DEFAULT_CHART_COLORS: ChartColors = LIGHT_CHART_COLORS
