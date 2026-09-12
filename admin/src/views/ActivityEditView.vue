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

      <DualImageField v-model="form.image" label="封面图（必填）" category="activity" />

      <div class="field">
        <label>正文</label>
        <RichTextEditor v-model="form.bodyHtml" />
      </div>

      <div class="form-actions">
        <button class="btn btn-primary" type="submit" :disabled="saving || !!blockingReason">
          {{ saving ? '保存中…' : '保存' }}
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
import RichTextEditor from '@/components/RichTextEditor.vue'
import { createActivity, getActivity, updateActivity } from '@/api/client'
import { createEmptyDual, dualsBlockingReason, isEmptyDual } from '@/utils/dualSlot'

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
  image: createEmptyDual(),
})

const blockingReason = computed(() => dualsBlockingReason([form.image]))

function apply(data) {
  form.title = data?.title || ''
  form.bodyHtml = data?.bodyHtml || ''
  form.image = createEmptyDual(data?.image)
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
    if (blockingReason.value) throw new Error(blockingReason.value)
    if (isEmptyDual(form.image)) throw new Error('请上传封面图')
    const image = {
      thumb: form.image.thumb || '',
      original: form.image.original || '',
    }
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

.form-block-hint {
  margin: 0;
  color: var(--danger);
  font-size: 0.85rem;
}
</style>
