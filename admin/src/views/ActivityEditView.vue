<template>
  <section>
    <div class="header">
      <div>
        <h1 class="page-title">{{ isNew ? '新建新闻' : '编辑新闻' }}</h1>
        <p class="page-sub">标题、封面图与正文（支持字号 / 加粗 / 颜色）</p>
      </div>
      <router-link class="btn btn-ghost" to="/activities">返回列表</router-link>
    </div>

    <div v-if="loadError" class="alert alert-error">{{ loadError }}</div>
    <div v-if="message" class="alert" :class="messageOk ? 'alert-ok' : 'alert-error'">{{ message }}</div>

    <form v-if="!loading" class="glass-panel form" @submit.prevent="onSubmit">
      <div class="field">
        <label for="title">标题（必填）</label>
        <input id="title" v-model.trim="form.title" required maxlength="255" placeholder="请输入新闻标题" />
      </div>

      <DualImageField v-model="form.image" label="封面图（必填）" />

      <div class="field">
        <label>正文</label>
        <RichTextEditor v-model="form.bodyHtml" />
      </div>

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
import RichTextEditor from '@/components/RichTextEditor.vue'
import { createActivity, getActivity, resolveDualImage, updateActivity } from '@/api/client'

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

const route = useRoute()
const router = useRouter()
const isNew = computed(() => route.name === 'activity-new')

const loading = ref(!isNew.value)
const saving = ref(false)
const loadError = ref('')
const message = ref('')
const messageOk = ref(false)

const form = reactive({
  title: '',
  bodyHtml: '',
  image: emptyDual(),
})

function apply(data) {
  form.title = data?.title || ''
  form.bodyHtml = data?.bodyHtml || ''
  form.image = emptyDual(data?.image)
}

onMounted(async () => {
  if (isNew.value) {
    loading.value = false
    return
  }
  try {
    const data = await getActivity(route.params.id)
    apply(data || {})
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
    if (!form.title.trim()) throw new Error('请填写新闻标题')
    if (isEmptyDual(form.image)) throw new Error('请上传封面图')
    const image = await resolveDualImage(form.image, 'activity')
    if (!image.thumb && !image.original) throw new Error('请上传封面图')
    const payload = {
      title: form.title.trim(),
      bodyHtml: form.bodyHtml || '',
      image,
    }
    if (isNew.value) {
      await createActivity(payload)
      messageOk.value = true
      message.value = '已创建'
    } else {
      await updateActivity(route.params.id, payload)
      messageOk.value = true
      message.value = '已保存'
    }
    await router.replace('/activities')
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
  max-width: 720px;
}
</style>
