<template>
  <div class="layout">
    <aside class="sidebar glass-panel">
      <div class="brand">
        <div class="brand-mark">M</div>
        <div>
          <div class="brand-name">温州目后空间设计有限公司</div>
          <div class="brand-sub">Admin CMS</div>
        </div>
      </div>

      <nav class="nav">
        <router-link to="/company">公司资料</router-link>
        <router-link to="/portfolios/residential">住宅作品集</router-link>
        <router-link to="/portfolios/commercial">商业作品集</router-link>
        <router-link to="/portfolios/office">办公作品集</router-link>
        <router-link to="/portfolios/installation">装置作品集</router-link>
        <router-link to="/activities">活动</router-link>
      </nav>

      <button type="button" class="btn btn-ghost logout" :disabled="loggingOut" @click="onLogout">
        {{ loggingOut ? '退出中…' : '退出登录' }}
      </button>
    </aside>

    <main class="main">
      <router-view />
    </main>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { logout } from '@/api/client'
import { setAuthState } from '@/router'

const router = useRouter()
const loggingOut = ref(false)

async function onLogout() {
  loggingOut.value = true
  try {
    await logout()
  } catch {
    // still clear local auth gate
  } finally {
    setAuthState(false)
    loggingOut.value = false
    router.replace({ name: 'login' })
  }
}
</script>

<style scoped>
.layout {
  display: grid;
  grid-template-columns: 240px 1fr;
  min-height: 100vh;
  gap: 1.25rem;
  padding: 1.25rem;
}

.sidebar {
  display: flex;
  flex-direction: column;
  padding: 1.25rem;
  position: sticky;
  top: 1.25rem;
  height: calc(100vh - 2.5rem);
}

.brand {
  display: flex;
  gap: 0.75rem;
  align-items: center;
  margin-bottom: 2rem;
}

.brand-mark {
  width: 42px;
  height: 42px;
  border-radius: 12px;
  display: grid;
  place-items: center;
  font-family: var(--font-display);
  font-weight: 700;
  background: linear-gradient(145deg, rgba(255, 255, 255, 0.16), rgba(255, 255, 255, 0.04));
  border: 1px solid var(--border-strong);
  color: var(--accent);
}

.brand-name {
  font-family: var(--font-display);
  font-weight: 600;
  letter-spacing: 0.04em;
}

.brand-sub {
  font-size: 0.75rem;
  color: var(--text-muted);
  letter-spacing: 0.08em;
  text-transform: uppercase;
}

.nav {
  display: flex;
  flex-direction: column;
  gap: 0.35rem;
  flex: 1;
}

.nav a {
  padding: 0.7rem 0.9rem;
  border-radius: 10px;
  color: var(--text-muted);
  border: 1px solid transparent;
  transition: background 0.2s ease, color 0.2s ease, border-color 0.2s ease;
}

.nav a:hover {
  color: var(--text);
  background: rgba(255, 255, 255, 0.04);
}

.nav a.router-link-active {
  color: var(--text);
  background: var(--glass-strong);
  border-color: var(--border);
}

.logout {
  margin-top: auto;
  width: 100%;
}

.main {
  min-width: 0;
  padding: 0.25rem 0.5rem 1.5rem;
}

@media (max-width: 860px) {
  .layout {
    grid-template-columns: 1fr;
  }

  .sidebar {
    position: static;
    height: auto;
  }

  .nav {
    flex-direction: row;
    flex-wrap: wrap;
  }

  .logout {
    margin-top: 1rem;
    width: auto;
  }
}
</style>
