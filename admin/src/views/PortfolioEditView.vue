<template>
  <section>
    <div class="header">
      <div>
        <h1 class="page-title">{{ isNew ? '新建作品' : '编辑作品' }}</h1>
        <p class="page-sub">封面与效果图 / 实景图在提交时上传</p>
      </div>
      <router-link class="btn btn-ghost" :to="`/portfolios/${listCategory}`">返回列表</router-link>
    </div>

    <div v-if="loadError" class="alert alert-error">{{ loadError }}</div>
    <div v-if="message" class="alert" :class="messageOk ? 'alert-ok' : 'alert-error'">{{ message }}</div>

    <form v-if="!loading" class="glass-panel form" @submit.prevent="onSubmit">
      <div class="field">
        <label for="category">类别</label>
        <select id="category" v-model="form.category">
          <option v-for="c in categories" :key="c.value" :value="c.value">{{ c.label }}</option>
        </select>
      </div>

      <div class="field">
        <label for="address">地址 / 项目名</label>
        <input id="address" v-model="form.address" />
      </div>

      <div class="grid-2">
        <div class="field">
          <label for="area">面积</label>
          <input id="area" v-model="form.area" />
        </div>
        <div class="field">
          <label for="style">风格</label>
          <input id="style" v-model="form.style" />
        </div>
      </div>

      <div class="field">
        <label for="heartFlow">心流 / 文案</label>
        <textarea id="heartFlow" v-model="form.heartFlow" rows="5" />
      </div>

      <DualImageField v-model="form.cover" label="封面" />

      <PairImageEditor v-model="form.pairs" title="效果图 / 实景图对比" />

      <div class="form-actions">
        <button class="btn btn-primary" type="submit" :disabled="saving">
          {{ saving ? '保存中…' : '保存' }}
        </button>
      </div>
    </form>

    <p v-else class="muted">加载中…</p>
  </section>
</template>

<script setup>
import { computed, onMounted, reactive, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import DualImageField from '@/components/DualImageField.vue'
import PairImageEditor from '@/components/PairImageEditor.vue'
import {
  createPortfolio,
  getPortfolio,
  putPortfolioImages,
  resolveDualImage,
  updatePortfolio,
} from '@/api/client'
import { PORTFOLIO_CATEGORIES } from '@/constants/portfolioCategories'

function emptyDual(src) {
  return {
    thumb: src?.thumb || '',
    original: src?.original || '',
    pending: null,
  }
}

function isEmptyDual(d) {
  if (!d) return true
  if (d.pending) return false
  return !d.thumb && !d.original
}

function zipPairs(renders, reals) {
  const r = renders || []
  const a = reals || []
  const n = Math.max(r.length, a.length)
  const pairs = []
  for (let i = 0; i < n; i++) {
    pairs.push({
      render: emptyDual(r[i]),
      real: emptyDual(a[i]),
    })
  }
  return pairs
}

const route = useRoute()
const router = useRouter()
const isNew = computed(() => route.name === 'portfolio-new')
const listCategory = computed(() => route.params.category)
const categories = PORTFOLIO_CATEGORIES

const loading = ref(!isNew.value)
const saving = ref(false)
const loadError = ref('')
const message = ref('')
const messageOk = ref(false)

const form = reactive({
  category: route.params.category || 'residential',
  address: '',
  area: '',
  style: '',
  heartFlow: '',
  cover: emptyDual(),
  pairs: [],
})

function applyPortfolio(data) {
  form.category = data?.category || route.params.category || 'residential'
  form.address = data?.address || ''
  form.area = data?.area || ''
  form.style = data?.style || ''
  form.heartFlow = data?.heartFlow || ''
  form.cover = emptyDual(data?.cover)
  form.pairs = zipPairs(data?.renders, data?.reals)
}

onMounted(async () => {
  if (isNew.value) {
    form.category = route.params.category || 'residential'
    loading.value = false
    return
  }
  try {
    const data = await getPortfolio(route.params.id)
    applyPortfolio(data || {})
  } catch (err) {
    loadError.value = err.message || '加载失败'
  } finally {
    loading.value = false
  }
})

async function resolvePairs(pairs) {
  const renders = []
  const reals = []
  for (const pair of pairs) {
    if (isEmptyDual(pair.render) && isEmptyDual(pair.real)) continue
    renders.push(await resolveDualImage(pair.render, 'portfolio'))
    reals.push(await resolveDualImage(pair.real, 'portfolio'))
  }
  return { renders, reals }
}

async function onSubmit() {
  saving.value = true
  message.value = ''
  try {
    const cover = await resolveDualImage(form.cover, 'portfolio')
    const { renders, reals } = await resolvePairs(form.pairs)

    const base = {
      slug: '',
      category: form.category,
      address: form.address,
      area: form.area,
      style: form.style,
      heartFlow: form.heartFlow,
      cover,
    }

    if (isNew.value) {
      const created = await createPortfolio({ ...base, renders, reals })
      const id = created?.id ?? created
      if (created && created.renders === undefined) {
        await putPortfolioImages(id, { renders, reals })
      }
      messageOk.value = true
      message.value = '已创建'
      await router.replace(`/portfolios/${form.category}`)
    } else {
      const id = route.params.id
      await updatePortfolio(id, base)
      await putPortfolioImages(id, { renders, reals })
      messageOk.value = true
      message.value = '已保存'
      await router.replace(`/portfolios/${form.category}`)
    }
  } catch (err) {
    messageOk.value = false
    message.value = err.message || '保存失败'
  } finally {
    saving.value = false
  }
}
</script>

<style scoped>
.header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  gap: 1rem;
  flex-wrap: wrap;
}

.form {
  padding: 1.35rem;
  display: flex;
  flex-direction: column;
  gap: 1rem;
}
</style>
