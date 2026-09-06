<template>
  <section>
    <h1 class="page-title">资费</h1>
    <p class="page-sub">上传资费长图（原图 + 缩略图）</p>

    <div v-if="loadError" class="alert alert-error">{{ loadError }}</div>
    <div v-if="message" class="alert" :class="messageOk ? 'alert-ok' : 'alert-error'">{{ message }}</div>

    <form v-if="!loading" class="glass-panel form" @submit.prevent="onSubmit">
      <DualImageField v-model="image" label="资费图" />

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
import { onMounted, ref } from 'vue'
import DualImageField from '@/components/DualImageField.vue'
import { getPricing, putPricing, resolveDualImage } from '@/api/client'

const loading = ref(true)
const saving = ref(false)
const loadError = ref('')
const message = ref('')
const messageOk = ref(false)

const image = ref({ thumb: '', original: '', pending: null })

onMounted(async () => {
  try {
    const data = await getPricing()
    image.value = {
      thumb: data?.thumb || '',
      original: data?.original || '',
      pending: null,
    }
  } catch (err) {
    loadError.value = err.message || '加载失败'
  } finally {
    loading.value = false
  }
})

async function onSubmit() {
  saving.value = true
  message.value = ''
  try {
    const pair = await resolveDualImage(image.value, 'pricing')
    const saved = await putPricing(pair)
    image.value = {
      thumb: saved?.thumb || pair.thumb,
      original: saved?.original || pair.original,
      pending: null,
    }
    messageOk.value = true
    message.value = '已保存'
  } catch (err) {
    messageOk.value = false
    message.value = err.message || '保存失败'
  } finally {
    saving.value = false
  }
}
</script>

<style scoped>
.form {
  padding: 1.35rem;
  max-width: 640px;
}
</style>
