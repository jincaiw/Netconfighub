<template>
  <div class="not-found">
    <t-result theme="default" title="404" description="您访问的页面不存在">
      <template #footer>
        <t-button theme="primary" @click="router.push('/dashboard')">返回首页</t-button>
      </template>
    </t-result>
  </div>
</template>

<script setup lang="ts">
import { useRouter } from 'vue-router'

const router = useRouter()
</script>

<style scoped>
.not-found {
  display: flex;
  justify-content: center;
  align-items: center;
  height: 100vh;
  height: 100dvh;
}
</style>
