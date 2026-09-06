<template>
  <div class="login-page">
    <form class="login-card glass-panel" @submit.prevent="onSubmit">
      <div class="brand">目后空间</div>
      <h1>管理后台</h1>
      <p class="muted">使用管理员密码登录</p>

      <div v-if="error" class="alert alert-error">{{ error }}</div>

      <div class="field">
        <label for="password">密码</label>
        <input
          id="password"
          v-model="password"
          type="password"
          autocomplete="current-password"
          required
          autofocus
        />
      </div>

      <button class="btn btn-primary" type="submit" :disabled="loading">
        {{ loading ? '登录中…' : '登录' }}
      </button>
    </form>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { login } from '@/api/client'
import { setAuthState } from '@/router'

const password = ref('')
const loading = ref(false)
const error = ref('')
const router = useRouter()
const route = useRoute()

async function onSubmit() {
  loading.value = true
  error.value = ''
  try {
    await login(password.value)
    setAuthState(true)
    const redirect = typeof route.query.redirect === 'string' ? route.query.redirect : '/company'
    router.replace(redirect)
  } catch (err) {
    error.value = err.message || '登录失败'
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.login-page {
  min-height: 100vh;
  display: grid;
  place-items: center;
  padding: 1.5rem;
}

.login-card {
  width: min(420px, 100%);
  padding: 2rem 1.75rem;
}

.brand {
  font-family: var(--font-display);
  font-size: 1.35rem;
  font-weight: 600;
  letter-spacing: 0.08em;
  color: var(--accent);
  margin-bottom: 0.35rem;
}

h1 {
  margin: 0 0 0.35rem;
  font-family: var(--font-display);
  font-size: 1.75rem;
  font-weight: 600;
}

.login-card .btn {
  width: 100%;
  margin-top: 0.5rem;
}
</style>
