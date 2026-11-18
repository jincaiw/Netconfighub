<template>
  <div class="t-result" :class="`t-result--${theme}`">
    <div class="t-result__body">
      <div v-if="themeIcon" class="t-result__icon"><component :is="themeIcon" /></div>
      <div v-if="title" class="t-result__title">{{ title }}</div>
      <div v-if="description" class="t-result__description">{{ description }}</div>
      <div v-if="$slots.footer" class="t-result__footer"><slot name="footer" /></div>
    </div>
  </div>
</template>

<script setup lang="ts">
// TDesign Vue Next 1.20.11 does not ship a Result component (no ESM export and
// no es/result style directory), yet NotFound.vue uses <t-result>. Register a
// local equivalent so the 404 page renders its title/description/footer.
import { computed } from 'vue'
import {
  CheckCircleIcon,
  ErrorCircleIcon,
  ErrorTriangleIcon,
  InfoCircleIcon,
} from 'tdesign-icons-vue-next'

const props = withDefaults(
  defineProps<{
    theme?: string
    title?: string
    description?: string
  }>(),
  { theme: 'default', title: '', description: '' },
)

const themeIconMap: Record<string, typeof InfoCircleIcon> = {
  default: InfoCircleIcon,
  info: InfoCircleIcon,
  success: CheckCircleIcon,
  warning: ErrorTriangleIcon,
  error: ErrorCircleIcon,
}

const themeIcon = computed(() => themeIconMap[props.theme] || InfoCircleIcon)
</script>

<style scoped>
.t-result {
  display: flex;
  justify-content: center;
  padding: 48px 24px;
}

.t-result__icon {
  width: 56px;
  height: 56px;
  margin: 0 auto 20px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 50%;
  font-size: 30px;
}

/* 状态样式按 theme 类展开；项目未安装 sass，使用等价的原生 CSS 作用域规则，不新增依赖。 */
.t-result--default .t-result__icon,
.t-result--info .t-result__icon {
  color: var(--td-brand-color);
  background: var(--td-brand-color-light);
}

.t-result--success .t-result__icon {
  color: var(--td-success-color-5);
  background: var(--td-success-color-1);
}

.t-result--warning .t-result__icon {
  color: var(--td-warning-color-5);
  background: var(--td-warning-color-1);
}

.t-result--error .t-result__icon {
  color: var(--td-error-color-5);
  background: var(--td-error-color-1);
}

.t-result--success .t-result__title {
  color: var(--td-success-color-6);
}

.t-result--warning .t-result__title {
  color: var(--td-warning-color-6);
}

.t-result--error .t-result__title {
  color: var(--td-error-color-6);
}

.t-result__title {
  color: var(--td-text-color-primary);
  font: var(--td-font-title-large);
  text-align: center;
}

.t-result__description {
  margin-top: 8px;
  color: var(--td-text-color-secondary);
  font: var(--td-font-body-medium);
  text-align: center;
}

.t-result__footer {
  margin-top: 24px;
  text-align: center;
}
</style>
