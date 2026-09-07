<template>
  <section>
    <div class="header">
      <div>
        <h1 class="page-title">活动</h1>
        <p class="page-sub">管理活动顺序与内容</p>
      </div>
      <router-link class="btn btn-primary" to="/activities/new">新建活动</router-link>
    </div>

    <div v-if="error" class="alert alert-error">{{ error }}</div>
    <div v-if="message" class="alert alert-ok">{{ message }}</div>
    <div v-if="loading" class="muted">加载中…</div>
    <div v-else-if="!items.length" class="glass-panel empty">暂无活动，点击「新建活动」开始录入。</div>

    <ul v-else class="list">
      <li v-for="(item, index) in items" :key="item.id" class="glass-panel row">
        <div class="cover">
          <img v-if="coverOf(item)" :src="coverOf(item)" alt="" />
          <span v-else class="no-cover">无图</span>
        </div>
        <div class="meta">
          <div class="title">{{ item.title || `活动 #${item.id}` }}</div>
        </div>
        <div class="order">
          <button type="button" class="btn btn-sm btn-ghost" :disabled="index === 0" @click="move(index, -1)">
            上移
          </button>
          <button
            type="button"
            class="btn btn-sm btn-ghost"
            :disabled="index === items.length - 1"
            @click="move(index, 1)"
          >
            下移
          </button>
        </div>
        <div class="actions">
          <router-link class="btn btn-sm" :to="`/activities/${item.id}`">编辑</router-link>
          <button type="button" class="btn btn-sm btn-danger" @click="onDelete(item)">删除</button>
        </div>
      </li>
    </ul>
  </section>
</template>

<script setup>
import { onMounted, ref } from 'vue'
import { deleteActivity, listActivities, reorderActivities } from '@/api/client'

const items = ref([])
const loading = ref(true)
const error = ref('')
const message = ref('')

function coverOf(item) {
  const c = item.image
  if (!c) return ''
  return c.thumb || c.original || ''
}

async function load() {
  loading.value = true
  error.value = ''
  try {
    const data = await listActivities()
    items.value = Array.isArray(data) ? data : data?.list || []
  } catch (err) {
    error.value = err.message || '加载失败'
  } finally {
    loading.value = false
  }
}

async function persistOrder() {
  message.value = ''
  try {
    await reorderActivities({ ids: items.value.map((x) => x.id) })
    message.value = '顺序已更新'
  } catch (err) {
    error.value = err.message || '排序保存失败'
    await load()
  }
}

async function move(index, delta) {
  const next = index + delta
  if (next < 0 || next >= items.value.length) return
  const copy = items.value.slice()
  const [row] = copy.splice(index, 1)
  copy.splice(next, 0, row)
  items.value = copy
  await persistOrder()
}

async function onDelete(item) {
  if (!confirm(`确定删除「${item.title || item.id}」？`)) return
  try {
    await deleteActivity(item.id)
    message.value = '已删除'
    await load()
  } catch (err) {
    error.value = err.message || '删除失败'
  }
}

onMounted(load)
</script>

<style scoped>
.header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  gap: 1rem;
  flex-wrap: wrap;
}

.list {
  list-style: none;
  margin: 0;
  padding: 0;
  display: flex;
  flex-direction: column;
  gap: 0.75rem;
}

.row {
  display: grid;
  grid-template-columns: 72px 1fr auto auto;
  gap: 1rem;
  align-items: center;
  padding: 0.85rem 1rem;
}

.cover {
  width: 72px;
  height: 56px;
  border-radius: 8px;
  overflow: hidden;
  border: 1px solid var(--border);
  background: rgba(0, 0, 0, 0.25);
  display: grid;
  place-items: center;
}

.cover img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.no-cover {
  font-size: 0.7rem;
  color: var(--text-muted);
}

.title {
  font-weight: 560;
}

.order,
.actions {
  display: flex;
  gap: 0.4rem;
  align-items: center;
  flex-wrap: wrap;
}

.empty {
  padding: 2rem;
  text-align: center;
  color: var(--text-muted);
}
</style>
