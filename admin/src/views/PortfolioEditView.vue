<template>
  <section>
    <div class="header">
      <div>
        <h1 class="page-title">{{ isNew ? '新建作品' : '编辑作品' }}</h1>
        <p class="page-sub">图片选择后自动上传；全部就绪后再保存</p>
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
          <label for="area">面积㎡</label>
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

      <DualImageField v-model="form.cover" label="封面" category="portfolio" />

      <ImageListEditor v-model="form.images" title="作品图" category="portfolio" />

      <div class="form-actions">
        <button class="btn btn-primary" type="submit" :disabled="saving || !!blockingReason">
          {{ savingLabel }}
        </button>
      </div>
      <p v-if="blockingReason" class="form-block-hint">{{ blockingReason }}</p>
    </form>

    <p v-else class="muted">加载中…</p>
  </section>
</template>

<script setup>
import { computed, onMounted, reactive, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import DualImageField from '@/components/DualImageField.vue'
import ImageListEditor from '@/components/ImageListEditor.vue'
import {
  createPortfolio,
  getPortfolio,
  putPortfolioImages,
  updatePortfolio,
} from '@/api/client'
import { PORTFOLIO_CATEGORIES } from '@/constants/portfolioCategories'
import { createEmptyDual, dualsBlockingReason, isEmptyDual } from '@/utils/dualSlot'

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
  cover: createEmptyDual(),
  images: [],
})

const blockingReason = computed(() =>
  dualsBlockingReason([form.cover, ...(form.images || [])]),
)

const savingLabel = computed(() => {
  if (!saving.value) return '保存'
  return '保存中…'
})

function toPayloadDual(d) {
  if (isEmptyDual(d)) return { thumb: '', original: '' }
  return {
    thumb: d.thumb || '',
    original: d.original || '',
  }
}

function applyPortfolio(data) {
  form.category = data?.category || route.params.category || 'residential'
  form.address = data?.address || ''
  form.area = data?.area || ''
  form.style = data?.style || ''
  form.heartFlow = data?.heartFlow || ''
  form.cover = createEmptyDual(data?.cover)
  form.images = (data?.images || []).map((img) => createEmptyDual(img))
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

async function onSubmit() {
  const block = blockingReason.value
  if (block) {
    messageOk.value = false
    message.value = block
    return
  }

  saving.value = true
  message.value = ''
  try {
    const cover = toPayloadDual(form.cover)
    const images = (form.images || [])
      .filter((item) => !isEmptyDual(item))
      .map((item) => toPayloadDual(item))

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
      const created = await createPortfolio({ ...base, images })
      const id = created?.id ?? created
      if (created && created.images === undefined) {
        await putPortfolioImages(id, { images })
      }
      messageOk.value = true
      message.value = '已创建'
      await router.replace(`/portfolios/${form.category}`)
    } else {
      const id = route.params.id
      await updatePortfolio(id, base)
      await putPortfolioImages(id, { images })
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

.form-block-hint {
  margin: 0;
  color: var(--danger);
  font-size: 0.85rem;
}
</style>
