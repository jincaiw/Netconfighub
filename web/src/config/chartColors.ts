/**
 * ECharts 配色：复刻 shadcn-admin `styles/theme.css` 的 --chart-1..5（zinc 主题）。
 * canvas 不支持 CSS 变量/oklch，故用等价 hex；亮/暗两套随 displayMode 切换。
 * primaryColor 对应 shadcn --primary（近黑/近白），用于主色柱状等。
 */
export interface ChartColors {
  textColor: string
  placeholderColor: string
  borderColor: string
  containerColor: string
  primaryColor: string
  seriesColors: string[]
}

// shadcn light: chart-1 orange / chart-2 teal / chart-3 blue / chart-4 yellow / chart-5 red
export const LIGHT_CHART_COLORS: ChartColors = {
  textColor: '#18181b',
  placeholderColor: '#71717a',
  borderColor: '#e4e4e7',
  containerColor: '#ffffff',
  primaryColor: '#27272a',
  seriesColors: ['#f97316', '#0d9488', '#2563eb', '#ca8a04', '#dc2626'],
}

// shadcn dark: chart-1 indigo / chart-2 emerald / chart-3 amber / chart-4 purple / chart-5 rose
export const DARK_CHART_COLORS: ChartColors = {
  textColor: '#fafafa',
  placeholderColor: '#a1a1aa',
  borderColor: 'rgba(255,255,255,0.1)',
  containerColor: '#1f1f22',
  primaryColor: '#e4e4e7',
  seriesColors: ['#818cf8', '#34d399', '#fbbf24', '#c084fc', '#fb7185'],
}

export const DEFAULT_CHART_COLORS: ChartColors = LIGHT_CHART_COLORS
