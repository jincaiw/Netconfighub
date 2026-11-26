<template>
  <span :class="containerCls">
    <span class="trend-icon-container">
      <svg v-if="type === 'down'" width="16" height="16" viewBox="0 0 16 16" fill="none" xmlns="http://www.w3.org/2000/svg">
        <path d="M11.5 8L8 11.5L4.5 8" stroke="currentColor" stroke-width="1.5" />
        <path d="M8 11L8 4" stroke="currentColor" stroke-width="1.5" />
      </svg>
      <svg v-else width="16" height="16" viewBox="0 0 16 16" fill="none" xmlns="http://www.w3.org/2000/svg">
        <path d="M4.5 8L8 4.5L11.5 8" stroke="currentColor" stroke-width="1.5" />
        <path d="M8 5V12" stroke="currentColor" stroke-width="1.5" />
      </svg>
    </span>
    <span>{{ describe }}</span>
  </span>
</template>

<script setup lang="ts">
import { computed } from 'vue'

// 忠实复刻 TDesign Starter `src/components/trend/index.vue`：SVG 箭头 + 三类语义色。
const props = withDefaults(
  defineProps<{
    type?: string
    describe?: string | number
    isReverseColor?: boolean
  }>(),
  { type: '', describe: undefined, isReverseColor: false },
)

const containerCls = computed(() => [
  'trend-container',
  {
    'trend-container__reverse': props.isReverseColor,
    'trend-container__up': !props.isReverseColor && props.type === 'up',
    'trend-container__down': !props.isReverseColor && props.type === 'down',
  },
])
</script>

<style scoped>
.trend-container {
  display: inline-flex;
  align-items: center;
  justify-content: center;
}

.trend-container__up {
  color: var(--td-error-color);

  .trend-icon-container {
    background: var(--td-error-color-2);
    margin-right: 8px;
  }
}

.trend-container__down {
  color: var(--td-success-color);

  .trend-icon-container {
    background: var(--td-success-color-2);
    margin-right: 8px;
  }
}

.trend-container__reverse {
  color: #fff;

  .trend-icon-container {
    background: var(--td-brand-color-5);
    margin-right: 8px;
  }
}

.trend-icon-container {
  border-radius: 50%;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 16px;
  height: 16px;
}
</style>
