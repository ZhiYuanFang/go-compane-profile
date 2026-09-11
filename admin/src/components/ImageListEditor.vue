<template>
  <div class="image-list">
    <div class="image-list-head">
      <h2>{{ title }}</h2>
      <div class="image-list-head-actions">
        <label class="btn btn-sm" :class="{ disabled: processing }">
          {{ processing ? `处理中 ${progressDone}/${progressTotal}` : '批量选择' }}
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
        <button type="button" class="btn btn-sm btn-ghost" :disabled="processing || !modelValue.length" @click="selectAll">
          全选
        </button>
        <button type="button" class="btn btn-sm btn-ghost" :disabled="processing || !selected.size" @click="clearSelection">
          取消全选
        </button>
        <button
          type="button"
          class="btn btn-sm btn-danger"
          :disabled="processing || !selected.size"
          @click="removeSelected"
        >
          删除所选{{ selected.size ? ` (${selected.size})` : '' }}
        </button>
      </div>
    </div>

    <p v-if="error" class="image-list-error">{{ error }}</p>
    <div v-if="!modelValue.length" class="muted empty-list">暂无图片，可批量选择上传</div>

    <div v-else class="gallery-grid">
      <div v-for="(item, i) in modelValue" :key="slotKey(item, i)" class="gallery-cell" :class="{ selected: selected.has(i) }">
        <label class="gallery-check" @click.stop>
          <input type="checkbox" :checked="selected.has(i)" :disabled="processing" @change="toggleSelect(i, $event)" />
        </label>
        <span class="gallery-index">#{{ i + 1 }}</span>
        <button
          type="button"
          class="gallery-thumb"
          :class="{ empty: !cellPreview(item) }"
          :disabled="processing"
          @click="openEditor(i)"
        >
          <img v-if="cellPreview(item)" :src="cellPreview(item)" :alt="`#${i + 1}`" />
          <span v-else class="gallery-placeholder">点击编辑</span>
        </button>
        <span v-if="item?.pending" class="gallery-badge">待上传</span>
      </div>
    </div>

    <Teleport to="body">
      <div
        v-if="editIndex !== null"
        class="edit-dialog"
        role="dialog"
        aria-modal="true"
        @click.self="closeEditor"
        @keydown.esc.prevent="closeEditor"
      >
        <div class="edit-panel glass-panel" @click.stop>
          <div class="edit-panel-head">
            <h3>编辑 #{{ editIndex + 1 }}</h3>
            <button type="button" class="btn btn-sm btn-ghost" @click="closeEditor">关闭</button>
          </div>
          <DualImageField
            v-if="editIndex !== null && modelValue[editIndex]"
            :model-value="modelValue[editIndex]"
            :label="`#${editIndex + 1}`"
            hint=""
            :disabled="processing"
            @update:model-value="(v) => update(editIndex, v)"
          />
          <div class="edit-panel-ops">
            <button type="button" class="btn btn-sm btn-ghost" :disabled="processing || editIndex === 0" @click="moveInEditor(-1)">
              上移
            </button>
            <button
              type="button"
              class="btn btn-sm btn-ghost"
              :disabled="processing || editIndex >= modelValue.length - 1"
              @click="moveInEditor(1)"
            >
              下移
            </button>
            <button type="button" class="btn btn-sm btn-danger" :disabled="processing" @click="removeInEditor">
              移除
            </button>
          </div>
        </div>
      </div>
    </Teleport>
  </div>
</template>

<script setup>
import { onBeforeUnmount, ref } from 'vue'
import DualImageField from '@/components/DualImageField.vue'
import { processDualImage, revokeObjectUrl } from '@/utils/imageProcess'

const props = defineProps({
  modelValue: { type: Array, default: () => [] },
  title: { type: String, required: true },
})

const emit = defineEmits(['update:modelValue'])

const processing = ref(false)
const progressDone = ref(0)
const progressTotal = ref(0)
const error = ref('')
const selected = ref(new Set())
const editIndex = ref(null)
const blobUrlCache = new Map()

function emptyDual() {
  return { thumb: '', original: '', pending: null }
}

function slotKey(item, i) {
  if (item?.pending?.thumb) return `p-${i}-${item.pending.thumb.name || i}`
  return `s-${i}-${item?.thumb || item?.original || 'empty'}`
}

function getBlobUrl(blob) {
  if (!blob) return ''
  let url = blobUrlCache.get(blob)
  if (!url) {
    url = URL.createObjectURL(blob)
    blobUrlCache.set(blob, url)
  }
  return url
}

function cellPreview(item) {
  if (!item) return ''
  if (item.thumb) return item.thumb
  if (item.original) return item.original
  if (item.pending?.thumb) return getBlobUrl(item.pending.thumb)
  return ''
}

function clearSelection() {
  selected.value = new Set()
}

function selectAll() {
  selected.value = new Set(props.modelValue.map((_, i) => i))
}

function toggleSelect(i, e) {
  const next = new Set(selected.value)
  if (e.target.checked) next.add(i)
  else next.delete(i)
  selected.value = next
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
  clearSelection()
  if (editIndex.value === i) editIndex.value = null
  else if (editIndex.value !== null && editIndex.value > i) editIndex.value -= 1
}

function move(i, delta) {
  const j = i + delta
  if (j < 0 || j >= props.modelValue.length) return
  const next = props.modelValue.slice()
  const [row] = next.splice(i, 1)
  next.splice(j, 0, row)
  emit('update:modelValue', next)
  clearSelection()
  if (editIndex.value === i) editIndex.value = j
  else if (editIndex.value === j) editIndex.value = i
}

function openEditor(i) {
  if (processing.value) return
  editIndex.value = i
}

function closeEditor() {
  editIndex.value = null
}

function moveInEditor(delta) {
  if (editIndex.value === null) return
  move(editIndex.value, delta)
}

function removeInEditor() {
  if (editIndex.value === null) return
  if (!confirm(`确定移除 #${editIndex.value + 1}？`)) return
  remove(editIndex.value)
}

function removeSelected() {
  const idxs = [...selected.value].sort((a, b) => b - a)
  if (!idxs.length) return
  if (!confirm(`确定删除所选 ${idxs.length} 张？`)) return
  const next = props.modelValue.slice()
  for (const i of idxs) next.splice(i, 1)
  emit('update:modelValue', next)
  clearSelection()
  editIndex.value = null
}

async function onBatchFiles(e) {
  const files = Array.from(e.target.files || []).filter(
    (f) => f && typeof f.type === 'string' && f.type.startsWith('image/'),
  )
  e.target.value = ''
  if (!files.length) return

  processing.value = true
  error.value = ''
  progressDone.value = 0
  progressTotal.value = files.length
  let list = props.modelValue.slice()
  try {
    for (let i = 0; i < files.length; i++) {
      const file = files[i]
      try {
        const pair = await processDualImage(file)
        list = [...list, { thumb: '', original: '', pending: pair }]
        emit('update:modelValue', list)
      } catch (err) {
        error.value = err.message || `第 ${i + 1} 张处理失败`
        // keep already-appended items; stop further files
        break
      }
      progressDone.value = i + 1
    }
  } finally {
    processing.value = false
    progressDone.value = 0
    progressTotal.value = 0
  }
}

onBeforeUnmount(() => {
  for (const url of blobUrlCache.values()) revokeObjectUrl(url)
  blobUrlCache.clear()
})
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

.empty-list {
  margin: 0;
  font-size: 0.9rem;
}

.image-list-error {
  margin: 0 0 0.5rem;
  color: #b42318;
  font-size: 0.88rem;
}

.gallery-grid {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 0.75rem;
}

.gallery-cell {
  position: relative;
  border: 1px solid var(--border);
  border-radius: 10px;
  background: rgba(0, 0, 0, 0.22);
  padding: 0.55rem;
  display: flex;
  flex-direction: column;
  gap: 0.35rem;
}

.gallery-cell.selected {
  border-color: var(--accent, #6ea8fe);
}

.gallery-check {
  position: absolute;
  top: 0.45rem;
  left: 0.45rem;
  z-index: 2;
  display: flex;
  padding: 0.15rem;
  background: rgba(0, 0, 0, 0.45);
  border-radius: 4px;
}

.gallery-index {
  align-self: flex-end;
  font-size: 0.75rem;
  color: var(--text-muted);
}

.gallery-thumb {
  width: 100%;
  aspect-ratio: 4 / 3;
  border: 1px dashed var(--border);
  border-radius: 8px;
  overflow: hidden;
  padding: 0;
  background: rgba(255, 255, 255, 0.03);
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
}

.gallery-thumb:not(.empty) {
  border-style: solid;
}

.gallery-thumb:disabled {
  opacity: 0.65;
  cursor: not-allowed;
}

.gallery-thumb img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.gallery-placeholder {
  font-size: 0.78rem;
  color: var(--text-muted);
}

.gallery-badge {
  font-size: 0.72rem;
  color: var(--accent, #6ea8fe);
}

.edit-dialog {
  position: fixed;
  inset: 0;
  z-index: 10000;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 1.25rem;
  background: rgba(0, 0, 0, 0.72);
  box-sizing: border-box;
}

.edit-panel {
  width: min(520px, 100%);
  max-height: 90vh;
  overflow: auto;
  padding: 1rem 1.1rem 1.15rem;
}

.edit-panel-head {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 0.75rem;
  margin-bottom: 0.85rem;
}

.edit-panel-head h3 {
  margin: 0;
  font-size: 1rem;
  font-weight: 560;
}

.edit-panel-ops {
  display: flex;
  flex-wrap: wrap;
  gap: 0.45rem;
  margin-top: 0.85rem;
}

@media (max-width: 900px) {
  .gallery-grid {
    grid-template-columns: repeat(2, 1fr);
  }
}

@media (max-width: 600px) {
  .gallery-grid {
    grid-template-columns: 1fr;
  }
}
</style>
