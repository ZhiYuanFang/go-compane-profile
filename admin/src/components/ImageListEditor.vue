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
      <div
        v-for="(item, i) in modelValue"
        :key="slotKey(item, i)"
        class="gallery-cell"
        :class="{ selected: selected.has(i), busy: slotBusy === i }"
      >
        <label class="gallery-check" @click.stop>
          <input type="checkbox" :checked="selected.has(i)" :disabled="processing" @change="toggleSelect(i, $event)" />
        </label>
        <span class="gallery-index">#{{ i + 1 }}</span>

        <div class="gallery-thumb-wrap">
          <button
            type="button"
            class="gallery-thumb"
            :class="{ empty: !cellPreview(item) }"
            :disabled="processing || slotBusy === i"
            @click="onThumbClick(i, item)"
          >
            <img v-if="cellPreview(item)" :src="cellPreview(item)" :alt="`#${i + 1}`" />
            <span v-else class="gallery-placeholder">选择图片</span>
          </button>

          <div v-if="cellPreview(item)" class="gallery-ops" @click.stop>
            <label class="op-btn">
              换
              <input
                type="file"
                accept="image/*"
                hidden
                :disabled="processing || slotBusy !== null"
                @change="onReplaceFile(i, $event)"
              />
            </label>
            <button type="button" class="op-btn" :disabled="processing || i === 0" @click="move(i, -1)">上</button>
            <button
              type="button"
              class="op-btn"
              :disabled="processing || i === modelValue.length - 1"
              @click="move(i, 1)"
            >
              下
            </button>
            <button type="button" class="op-btn danger" :disabled="processing" @click="removeOne(i)">删</button>
          </div>
          <div v-else class="gallery-ops always" @click.stop>
            <label class="op-btn">
              选图
              <input
                type="file"
                accept="image/*"
                hidden
                :disabled="processing || slotBusy !== null"
                @change="onReplaceFile(i, $event)"
              />
            </label>
            <button type="button" class="op-btn danger" :disabled="processing" @click="removeOne(i)">删</button>
          </div>
        </div>

        <span v-if="slotBusy === i" class="gallery-badge">处理中…</span>
        <span v-else-if="item?.pending" class="gallery-badge">待上传</span>
      </div>
    </div>

    <input
      ref="emptyPickInput"
      type="file"
      accept="image/*"
      hidden
      :disabled="processing || slotBusy !== null"
      @change="onEmptyPickFile"
    />

    <Teleport to="body">
      <div
        v-if="lightboxOpen && lightboxUrl"
        class="lightbox"
        role="dialog"
        aria-modal="true"
        @click.self="closeLightbox"
      >
        <button type="button" class="lightbox-close btn btn-sm" @click="closeLightbox">关闭</button>
        <img class="lightbox-img" :src="lightboxUrl" alt="" @click.stop />
      </div>
    </Teleport>
  </div>
</template>

<script setup>
import { onBeforeUnmount, onMounted, onUnmounted, ref } from 'vue'
import { processDualImage, revokeObjectUrl } from '@/utils/imageProcess'

const props = defineProps({
  modelValue: { type: Array, default: () => [] },
  title: { type: String, required: true },
})

const emit = defineEmits(['update:modelValue'])

const processing = ref(false)
const progressDone = ref(0)
const progressTotal = ref(0)
const slotBusy = ref(null)
const error = ref('')
const selected = ref(new Set())
const blobUrlCache = new Map()
const lightboxOpen = ref(false)
const lightboxUrl = ref('')
const emptyPickInput = ref(null)
const emptyPickIndex = ref(null)

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

function cellLightbox(item) {
  if (!item) return ''
  if (item.pending?.original) return getBlobUrl(item.pending.original)
  if (item.pending?.thumb) return getBlobUrl(item.pending.thumb)
  return item.original || item.thumb || ''
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
}

function removeOne(i) {
  if (!confirm(`确定移除 #${i + 1}？`)) return
  remove(i)
}

function move(i, delta) {
  const j = i + delta
  if (j < 0 || j >= props.modelValue.length) return
  const next = props.modelValue.slice()
  const [row] = next.splice(i, 1)
  next.splice(j, 0, row)
  emit('update:modelValue', next)
  clearSelection()
}

function onThumbClick(i, item) {
  if (processing.value || slotBusy.value !== null) return
  if (!cellPreview(item)) {
    emptyPickIndex.value = i
    emptyPickInput.value?.click()
    return
  }
  lightboxUrl.value = cellLightbox(item)
  lightboxOpen.value = true
}

function closeLightbox() {
  lightboxOpen.value = false
  lightboxUrl.value = ''
}

function onKeydown(e) {
  if (e.key === 'Escape' && lightboxOpen.value) closeLightbox()
}

onMounted(() => window.addEventListener('keydown', onKeydown))
onUnmounted(() => window.removeEventListener('keydown', onKeydown))

async function applyFileToSlot(i, file) {
  if (!file || !file.type?.startsWith('image/')) return
  slotBusy.value = i
  error.value = ''
  try {
    const pair = await processDualImage(file)
    const cur = props.modelValue[i] || emptyDual()
    update(i, {
      thumb: cur.thumb || '',
      original: cur.original || '',
      pending: pair,
    })
  } catch (err) {
    error.value = err.message || '图片处理失败'
  } finally {
    slotBusy.value = null
  }
}

async function onReplaceFile(i, e) {
  const file = e.target.files?.[0]
  e.target.value = ''
  if (!file) return
  await applyFileToSlot(i, file)
}

async function onEmptyPickFile(e) {
  const file = e.target.files?.[0]
  e.target.value = ''
  const i = emptyPickIndex.value
  emptyPickIndex.value = null
  if (i == null || !file) return
  await applyFileToSlot(i, file)
}

function removeSelected() {
  const idxs = [...selected.value].sort((a, b) => b - a)
  if (!idxs.length) return
  if (!confirm(`确定删除所选 ${idxs.length} 张？`)) return
  const next = props.modelValue.slice()
  for (const i of idxs) next.splice(i, 1)
  emit('update:modelValue', next)
  clearSelection()
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
  padding: 0.45rem;
  display: flex;
  flex-direction: column;
  gap: 0.3rem;
}

.gallery-cell.selected {
  border-color: var(--accent, #6ea8fe);
}

.gallery-check {
  position: absolute;
  top: 0.4rem;
  left: 0.4rem;
  z-index: 3;
  display: flex;
  padding: 0.15rem;
  background: rgba(0, 0, 0, 0.45);
  border-radius: 4px;
}

.gallery-index {
  align-self: flex-end;
  font-size: 0.72rem;
  color: var(--text-muted);
  line-height: 1;
}

.gallery-thumb-wrap {
  position: relative;
  width: 100%;
}

.gallery-thumb {
  width: 100%;
  aspect-ratio: 4 / 3;
  border: 1px dashed var(--border);
  border-radius: 8px;
  overflow: hidden;
  padding: 0;
  background: rgba(255, 255, 255, 0.03);
  cursor: zoom-in;
  display: flex;
  align-items: center;
  justify-content: center;
}

.gallery-thumb.empty {
  cursor: pointer;
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

.gallery-ops {
  position: absolute;
  left: 0.35rem;
  right: 0.35rem;
  bottom: 0.35rem;
  display: flex;
  gap: 0.25rem;
  justify-content: center;
  flex-wrap: wrap;
  opacity: 0;
  pointer-events: none;
  transition: opacity 0.12s ease;
  z-index: 2;
}

.gallery-cell:hover .gallery-ops,
.gallery-cell:focus-within .gallery-ops,
.gallery-ops.always {
  opacity: 1;
  pointer-events: auto;
}

.op-btn {
  box-sizing: border-box;
  min-width: 1.7rem;
  padding: 0.2rem 0.35rem;
  border-radius: 5px;
  border: 1px solid rgba(255, 255, 255, 0.25);
  background: rgba(0, 0, 0, 0.72);
  color: #fff;
  font-size: 0.7rem;
  line-height: 1.2;
  cursor: pointer;
  text-align: center;
}

.op-btn:disabled {
  opacity: 0.45;
  cursor: not-allowed;
}

.op-btn.danger {
  border-color: rgba(180, 35, 24, 0.7);
  color: #ffb4a8;
}

.gallery-badge {
  font-size: 0.7rem;
  color: var(--accent, #6ea8fe);
}

.lightbox {
  position: fixed;
  inset: 0;
  z-index: 10000;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 2rem;
  background: rgba(0, 0, 0, 0.82);
  box-sizing: border-box;
}

.lightbox-close {
  position: absolute;
  top: 1rem;
  right: 1rem;
}

.lightbox-img {
  max-width: min(96vw, 1200px);
  max-height: 90vh;
  object-fit: contain;
  border-radius: 4px;
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
