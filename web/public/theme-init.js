// 同步应用首屏主题，避免 deferred 模块执行前闪白（FOUC）。
// 浅色写入空字符串，与 setting.ts / theme.spec.ts 的语义保持一致。
// 独立文件（非内联）以满足生产 CSP script-src 'self'。
;(function () {
  var mode = 'auto'
  try {
    var stored = window.localStorage.getItem('nch-theme-mode')
    if (stored === 'light' || stored === 'dark' || stored === 'auto') mode = stored
  } catch (e) {
    mode = 'auto'
  }
  var isDark = false
  try {
    isDark = mode === 'dark' || (mode === 'auto' && window.matchMedia('(prefers-color-scheme: dark)').matches)
  } catch (e) {
    isDark = false
  }
  document.documentElement.setAttribute('theme-mode', isDark ? 'dark' : '')
})()
