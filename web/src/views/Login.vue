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
      <n-card class="login-card">
        <n-form ref="formRef" :model="formData" :rules="rules" label-placement="top">
          <n-form-item path="username" label="用户名">
            <n-input
              v-model:value="formData.username"
              placeholder="请输入用户名"
              @keydown.enter="handleLogin"
            >
              <template #prefix>
                <n-icon :component="PersonOutline" />
              </template>
            </n-input>
          </n-form-item>
          <n-form-item path="password" label="密码">
            <n-input
              v-model:value="formData.password"
              type="password"
              show-password-on="click"
              placeholder="请输入密码"
              @keydown.enter="handleLogin"
            >
              <template #prefix>
                <n-icon :component="LockClosedOutline" />
              </template>
            </n-input>
          </n-form-item>
          <n-button
            type="primary"
            block
            size="large"
            :loading="loading"
            @click="handleLogin"
          >
            登录
          </n-button>
        </n-form>
      </n-card>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { NCard, NForm, NFormItem, NInput, NButton, NIcon, useMessage } from 'naive-ui'
import type { FormInst, FormRules } from 'naive-ui'
import { PersonOutline, LockClosedOutline } from '@vicons/ionicons5'
import { useAuthStore } from '@/stores/auth'

const router = useRouter()
const route = useRoute()
const message = useMessage()
const authStore = useAuthStore()

const formRef = ref<FormInst | null>(null)
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
  try {
    await formRef.value?.validate()
  } catch {
    return
  }

  loading.value = true
  try {
    await authStore.login(formData)
    message.success('登录成功')
    const redirect = (route.query.redirect as string) || '/dashboard'
    router.push(redirect)
  } catch (err: any) {
    const msg = err?.response?.data?.message || '登录失败，请检查用户名和密码'
    message.error(msg)
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.login-container {
  width: 100%;
  min-height: 100vh;
  display: grid;
  grid-template-columns: minmax(420px, 46%) 1fr;
  background: #f4f7fb;
}

.brand-panel {
  position: relative;
  display: flex;
  flex-direction: column;
  justify-content: center;
  padding: 9vw;
  overflow: hidden;
  color: #f8fafc;
  background:
    radial-gradient(circle at 18% 20%, rgba(16, 185, 129, 0.22), transparent 28%),
    linear-gradient(145deg, #071a33 0%, #0b2749 100%);
}

.brand-panel::after {
  content: '';
  position: absolute;
  right: -140px;
  bottom: -140px;
  width: 420px;
  height: 420px;
  border: 1px solid rgba(148, 163, 184, 0.16);
  border-radius: 50%;
  box-shadow:
    0 0 0 70px rgba(148, 163, 184, 0.05),
    0 0 0 140px rgba(148, 163, 184, 0.03);
}

.brand-mark {
  width: 52px;
  height: 52px;
  align-items: center;
  justify-content: center;
  display: flex;
  border: 1px solid rgba(52, 211, 153, 0.45);
  border-radius: 14px;
  color: #34d399;
  background: rgba(16, 185, 129, 0.1);
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
  color: #b9c8dc;
  font-size: 17px;
  line-height: 1.8;
}

.brand-features {
  position: relative;
  z-index: 1;
  display: grid;
  gap: 14px;
  margin-top: 48px;
  color: #dbe6f3;
  font-size: 14px;
}

.brand-features span::before {
  content: '';
  display: inline-block;
  width: 7px;
  height: 7px;
  margin-right: 12px;
  border-radius: 50%;
  background: #34d399;
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
  color: #0f172a;
  letter-spacing: -0.02em;
  margin-bottom: 8px;
}

.login-subtitle {
  font-size: 14px;
  color: #64748b;
}

.login-card {
  width: 420px;
  border-color: #dfe7f1;
  border-radius: 12px;
  box-shadow: 0 18px 45px rgba(15, 23, 42, 0.08);
}

.login-card :deep(.n-card__content) {
  padding: 30px;
}

@media (max-width: 900px) {
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
}

@media (max-width: 520px) {
  .brand-panel {
    min-height: 190px;
  }

  .login-header,
  .login-card {
    width: 100%;
  }
}
</style>
