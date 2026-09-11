<template>
  <div class="image-list">
    <div class="image-list-head">
      <h2>{{ title }}</h2>
      <div class="image-list-head-actions">
        <label class="btn btn-sm" :class="{ disabled: processing }">
          {{ processing ? '处理中…' : '批量选择' }}
          <input
            type="file"
            accept="image/*"
            multiple
            hidden
            :disabled="processing"
            @change="onBatchFiles"
          />
        </label>
        <button type="button" class="btn btn-sm btn-ghost" :disabled="processing" @click="add">
          添加一张
        </button>
      </div>
    </div>

    <p v-if="error" class="image-list-error">{{ error }}</p>
    <div v-if="!modelValue.length" class="muted empty-list">暂无图片，可批量选择上传</div>

    <div v-for="(item, i) in modelValue" :key="i" class="image-list-item">
      <DualImageField
        :model-value="item"
        :label="`#${i + 1}`"
        hint=""
        :disabled="processing"
        @update:model-value="(v) => update(i, v)"
      />
      <div class="image-list-ops">
        <button type="button" class="btn btn-sm btn-ghost" :disabled="processing || i === 0" @click="move(i, -1)">
          上移
        </button>
        <button
          type="button"
          class="btn btn-sm btn-ghost"
          :disabled="processing || i === modelValue.length - 1"
          @click="move(i, 1)"
        >
          下移
        </button>
        <button type="button" class="btn btn-sm btn-danger" :disabled="processing" @click="remove(i)">
          移除
        </button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import DualImageField from '@/components/DualImageField.vue'
import { processDualImage } from '@/utils/imageProcess'

const props = defineProps({
  modelValue: { type: Array, default: () => [] },
  title: { type: String, required: true },
})

const emit = defineEmits(['update:modelValue'])

const processing = ref(false)
const error = ref('')

function emptyDual() {
  return { thumb: '', original: '', pending: null }
}

function update(i, val) {
  const next = props.modelValue.slice()
  next[i] = val
  emit('update:modelValue', next)
}

function add() {
  emit('update:modelValue', [...props.modelValue, emptyDual()])
}

function remove(i) {
  const next = props.modelValue.slice()
  next.splice(i, 1)
  emit('update:modelValue', next)
}

function move(i, delta) {
  const j = i + delta
  if (j < 0 || j >= props.modelValue.length) return
  const next = props.modelValue.slice()
  const [row] = next.splice(i, 1)
  next.splice(j, 0, row)
  emit('update:modelValue', next)
}

async function onBatchFiles(e) {
  const files = Array.from(e.target.files || []).filter(
    (f) => f && typeof f.type === 'string' && f.type.startsWith('image/'),
  )
  e.target.value = ''
  if (!files.length) return

  processing.value = true
  error.value = ''
  const added = []
  try {
    for (const file of files) {
      const pair = await processDualImage(file)
      added.push({
        thumb: '',
        original: '',
        pending: pair,
      })
    }
    emit('update:modelValue', [...props.modelValue, ...added])
  } catch (err) {
    error.value = err.message || '批量处理失败'
    if (added.length) {
      emit('update:modelValue', [...props.modelValue, ...added])
    }
  } finally {
    processing.value = false
  }
}
</script>

<style scoped>
.image-list {
  margin-top: 0.5rem;
  padding-top: 1rem;
  border-top: 1px solid var(--border);
}

.image-list-head {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 0.75rem;
  flex-wrap: wrap;
  margin-bottom: 0.75rem;
}

.image-list-head h2 {
  margin: 0;
  font-family: var(--font-display);
  font-size: 1.05rem;
  font-weight: 560;
}

.image-list-head-actions {
  display: flex;
  gap: 0.5rem;
  flex-wrap: wrap;
  align-items: center;
}

.image-list-head-actions .btn.disabled,
.image-list-head-actions .btn:has(input:disabled) {
  opacity: 0.6;
  pointer-events: none;
}

.image-list-item {
  display: grid;
  grid-template-columns: 1fr auto;
  gap: 0.75rem;
  align-items: start;
  margin-bottom: 0.75rem;
}

.image-list-ops {
  display: flex;
  flex-direction: column;
  gap: 0.35rem;
}

.empty-list {
  margin: 0;
  font-size: 0.9rem;
}

.image-list-error {
  margin: 0 0 0.5rem;
  color: #b42318;
  font-size: 0.88rem;
}

@media (max-width: 700px) {
  .image-list-item {
    grid-template-columns: 1fr;
  }

  .image-list-ops {
    flex-direction: row;
    flex-wrap: wrap;
  }
}
</style>
