<template>
  <div class="login-container">
    <section class="brand-panel">
      <div class="brand-mark">N</div>
      <div class="brand-copy">
        <h1>NetConfigHub</h1>
        <p>让网络设备配置备份、版本追踪与基线审计保持清晰可控。</p>
      </div>
      <div class="brand-features">
        <span>自动备份与失败重试</span>
        <span>Git 版本历史与差异追踪</span>
        <span>基线偏差与安全审计</span>
      </div>
    </section>
    <div class="login-content">
      <div class="login-header">
        <h1 class="login-title">登录管理控制台</h1>
        <p class="login-subtitle">使用管理员账户继续</p>
      </div>
      <t-card class="login-card">
        <t-form ref="formRef" :data="formData" :rules="rules" label-align="top">
          <t-form-item label="用户名" name="username">
            <t-input
              v-model="formData.username"
              placeholder="请输入用户名"
              @enter="handleLogin"
            >
              <template #prefixIcon><UserIcon /></template>
            </t-input>
          </t-form-item>
          <t-form-item label="密码" name="password">
            <t-input
              v-model="formData.password"
              type="password"
              placeholder="请输入密码"
              @enter="handleLogin"
            >
              <template #prefixIcon><LockOnIcon /></template>
            </t-input>
          </t-form-item>
          <t-button theme="primary" block size="large" :loading="loading" @click="handleLogin">
            登录
          </t-button>
        </t-form>
      </t-card>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { MessagePlugin, type FormInstanceFunctions, type FormRules } from 'tdesign-vue-next'
import { UserIcon, LockOnIcon } from 'tdesign-icons-vue-next'
import { useAuthStore } from '@/stores/auth'

const router = useRouter()
const route = useRoute()
const authStore = useAuthStore()

const formRef = ref<FormInstanceFunctions | null>(null)
const loading = ref(false)

const formData = reactive({
  username: '',
  password: '',
})

const rules: FormRules = {
  username: [{ required: true, message: '请输入用户名', trigger: 'blur' }],
  password: [{ required: true, message: '请输入密码', trigger: 'blur' }],
}

async function handleLogin() {
  const valid = await formRef.value?.validate()
  if (valid !== true) return

  loading.value = true
  try {
    await authStore.login(formData)
    MessagePlugin.success('登录成功')
    const redirect = (route.query.redirect as string) || '/dashboard'
    router.push(redirect)
  } catch (err: any) {
    const msg = err?.response?.data?.message || '登录失败，请检查用户名和密码'
    MessagePlugin.error(msg)
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.login-container {
  /* 登录页品牌深蓝面板为有意设计，独立于 --td-brand 的局部设计变量；集中在此管理。 */
  --login-brand-deep: #071a33;
  --login-brand-deep-2: #0b2749;
  --login-accent: #34d399;
  --login-accent-glow: rgba(16, 185, 129, 0.22);
  --login-accent-wash: rgba(16, 185, 129, 0.1);
  --login-accent-border: rgba(52, 211, 153, 0.45);
  --login-ring: rgba(148, 163, 184, 0.16);
  --login-ring-soft: rgba(148, 163, 184, 0.05);
  --login-ring-faint: rgba(148, 163, 184, 0.03);
  --login-copy: #b9c8dc;
  --login-feature: #dbe6f3;

  width: 100%;
  min-height: 100vh;
  min-height: 100dvh;
  display: grid;
  grid-template-columns: minmax(420px, 46%) 1fr;
  background: var(--td-bg-color-page);
}

.brand-panel {
  position: relative;
  display: flex;
  flex-direction: column;
  justify-content: center;
  padding: 9vw;
  overflow: hidden;
  color: var(--td-text-color-anti);
  background:
    radial-gradient(circle at 18% 20%, var(--login-accent-glow), transparent 28%),
    linear-gradient(145deg, var(--login-brand-deep) 0%, var(--login-brand-deep-2) 100%);
}

.brand-panel::after {
  content: '';
  position: absolute;
  right: -140px;
  bottom: -140px;
  width: 420px;
  height: 420px;
  border: 1px solid var(--login-ring);
  border-radius: 50%;
  box-shadow:
    0 0 0 70px var(--login-ring-soft),
    0 0 0 140px var(--login-ring-faint);
}

.brand-mark {
  width: 52px;
  height: 52px;
  align-items: center;
  justify-content: center;
  display: flex;
  border: 1px solid var(--login-accent-border);
  border-radius: 14px;
  color: var(--login-accent);
  background: var(--login-accent-wash);
  font-size: 24px;
  font-weight: 800;
}

.brand-copy {
  position: relative;
  z-index: 1;
  margin-top: 34px;
  max-width: 560px;
}

.brand-copy h1 {
  font-size: clamp(36px, 4vw, 58px);
  line-height: 1.1;
  letter-spacing: -0.04em;
}

.brand-copy p {
  margin-top: 20px;
  color: var(--login-copy);
  font-size: 17px;
  line-height: 1.8;
}

.brand-features {
  position: relative;
  z-index: 1;
  display: grid;
  gap: 14px;
  margin-top: 48px;
  color: var(--login-feature);
  font-size: 14px;
}

.brand-features span::before {
  content: '';
  display: inline-block;
  width: 7px;
  height: 7px;
  margin-right: 12px;
  border-radius: 50%;
  background: var(--login-accent);
}

.login-content {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 48px 24px;
}

.login-header {
  width: 420px;
  margin-bottom: 24px;
}

.login-title {
  font-size: 28px;
  font-weight: 700;
  color: var(--td-text-color-primary);
  letter-spacing: -0.02em;
  margin-bottom: 8px;
}

.login-subtitle {
  font-size: 14px;
  color: var(--td-text-color-secondary);
}

.login-card {
  width: 420px;
  border-radius: var(--td-radius-large);
}

.login-card :deep(.t-card__body) {
  padding: 30px;
}

/* 断点统一到 768px（见 global.css 断点契约），合并原 900/520 两档。 */
@media (max-width: 768px) {
  .login-container {
    grid-template-columns: 1fr;
  }

  .brand-panel {
    min-height: 220px;
    padding: 36px 28px;
  }

  .brand-copy {
    margin-top: 18px;
  }

  .brand-copy h1 {
    font-size: 32px;
  }

  .brand-copy p {
    margin-top: 10px;
    font-size: 14px;
  }

  .brand-features {
    display: none;
  }

  .login-header,
  .login-card {
    width: 100%;
  }
}
</style>
