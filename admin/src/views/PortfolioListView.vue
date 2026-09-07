<template>
  <section>
    <div class="header">
      <div>
        <h1 class="page-title">{{ title }}作品集</h1>
        <p class="page-sub">管理本类作品顺序与内容</p>
      </div>
      <router-link class="btn btn-primary" :to="`/portfolios/${category}/new`">新建作品</router-link>
    </div>

    <div v-if="error" class="alert alert-error">{{ error }}</div>
    <div v-if="message" class="alert alert-ok">{{ message }}</div>

    <div v-if="loading" class="muted">加载中…</div>

    <div v-else-if="!items.length" class="glass-panel empty">暂无作品，点击「新建作品」开始录入。</div>

    <ul v-else class="list">
      <li v-for="(item, index) in items" :key="item.id" class="glass-panel row">
        <div class="cover">
          <img v-if="coverOf(item)" :src="coverOf(item)" alt="" />
          <span v-else class="no-cover">无封面</span>
        </div>

        <div class="meta">
          <div class="title">{{ item.address || item.slug || `作品 #${item.id}` }}</div>
          <div class="sub muted">
            <span>热度 {{ item.viewCount ?? 0 }}</span>
            <span v-if="item.slug"> · {{ item.slug }}</span>
            <span v-if="item.area"> · {{ item.area }}</span>
            <span v-if="item.style"> · {{ item.style }}</span>
          </div>
        </div>

        <div class="order">
          <input
            class="order-input"
            type="number"
            :value="index + 1"
            min="1"
            :max="items.length"
            @change="onOrderInput(index, $event)"
          />
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
          <router-link class="btn btn-sm" :to="`/portfolios/${category}/${item.id}`">编辑</router-link>
          <button type="button" class="btn btn-sm btn-danger" @click="onDelete(item)">删除</button>
        </div>
      </li>
    </ul>
  </section>
</template>

<script setup>
import { computed, onMounted, ref, watch } from 'vue'
import { useRoute } from 'vue-router'
import { deletePortfolio, listPortfolios, reorderPortfolios } from '@/api/client'
import { categoryLabel } from '@/constants/portfolioCategories'

const route = useRoute()
const category = computed(() => route.params.category)
const title = computed(() => categoryLabel(category.value))

const items = ref([])
const loading = ref(true)
const error = ref('')
const message = ref('')
const savingOrder = ref(false)

function coverOf(item) {
  const c = item.cover
  if (!c) return ''
  if (typeof c === 'string') return c
  return c.thumb || c.original || ''
}

async function load() {
  loading.value = true
  error.value = ''
  try {
    const data = await listPortfolios(category.value)
    items.value = Array.isArray(data) ? data : data?.list || data?.items || []
  } catch (err) {
    error.value = err.message || '加载失败'
  } finally {
    loading.value = false
  }
}

async function persistOrder() {
  savingOrder.value = true
  message.value = ''
  try {
    await reorderPortfolios({
      category: category.value,
      ids: items.value.map((x) => x.id),
    })
    message.value = '顺序已更新'
  } catch (err) {
    error.value = err.message || '排序保存失败'
    await load()
  } finally {
    savingOrder.value = false
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

async function onOrderInput(index, e) {
  let target = Number(e.target.value)
  if (!Number.isFinite(target)) {
    e.target.value = index + 1
    return
  }
  target = Math.min(Math.max(1, Math.round(target)), items.value.length)
  if (target - 1 === index) return
  const copy = items.value.slice()
  const [row] = copy.splice(index, 1)
  copy.splice(target - 1, 0, row)
  items.value = copy
  e.target.value = target
  await persistOrder()
}

async function onDelete(item) {
  if (!confirm(`确定删除「${item.address || item.slug || item.id}」？`)) return
  try {
    await deletePortfolio(item.id)
    message.value = '已删除'
    await load()
  } catch (err) {
    error.value = err.message || '删除失败'
  }
}

watch(category, load)
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

.sub {
  font-size: 0.85rem;
  margin-top: 0.2rem;
}

.order,
.actions {
  display: flex;
  gap: 0.4rem;
  align-items: center;
  flex-wrap: wrap;
}

.order-input {
  width: 3.2rem;
  padding: 0.35rem 0.4rem;
  border-radius: 8px;
  border: 1px solid var(--border);
  background: rgba(0, 0, 0, 0.28);
  text-align: center;
}

.empty {
  padding: 2rem;
  text-align: center;
  color: var(--text-muted);
}

@media (max-width: 900px) {
  .row {
    grid-template-columns: 72px 1fr;
  }

  .order,
  .actions {
    grid-column: 1 / -1;
  }
}
</style>
