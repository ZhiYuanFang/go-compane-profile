<template>
  <div class="image-list">
    <div class="image-list-head">
      <h2>{{ title }}</h2>
      <div class="image-list-head-actions">
        <label class="btn btn-sm" :class="{ disabled: picking }">
          {{ picking ? '添加中…' : '批量选择' }}
          <input
            type="file"
            accept="image/*"
            multiple
            hidden
            :disabled="picking"
            @change="onBatchFiles"
          />
        </label>
        <button type="button" class="btn btn-sm btn-ghost" :disabled="picking" @click="add">
          添加一张
        </button>
        <button type="button" class="btn btn-sm btn-ghost" :disabled="picking || !listAcc.length" @click="selectAll">
          全选
        </button>
        <button type="button" class="btn btn-sm btn-ghost" :disabled="picking || !selected.size" @click="clearSelection">
          取消全选
        </button>
        <button
          type="button"
          class="btn btn-sm btn-danger"
          :disabled="picking || !selected.size"
          @click="removeSelected"
        >
          删除所选{{ selected.size ? ` (${selected.size})` : '' }}
        </button>
      </div>
    </div>

    <p v-if="error" class="image-list-error">{{ error }}</p>
    <div v-if="!listAcc.length" class="muted empty-list">暂无图片，可批量选择（选后自动上传）</div>

    <div v-else class="gallery-grid">
      <div
        v-for="(item, i) in listAcc"
        :key="item.slotId || i"
        class="gallery-cell"
        :class="{ selected: selected.has(item.slotId || i), busy: isBusy(item) }"
      >
        <label class="gallery-check" @click.stop>
          <input
            type="checkbox"
            :checked="selected.has(item.slotId || i)"
            :disabled="picking"
            @change="toggleSelect(item.slotId || i, $event)"
          />
        </label>
        <span class="gallery-index">#{{ i + 1 }}</span>

        <div class="gallery-thumb-wrap">
          <button
            type="button"
            class="gallery-thumb"
            :class="{ empty: !cellPreview(item) }"
            :disabled="picking"
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
                :disabled="picking"
                @change="onReplaceFile(item, $event)"
              />
            </label>
            <button type="button" class="op-btn" :disabled="picking || i === 0" @click="move(i, -1)">上</button>
            <button
              type="button"
              class="op-btn"
              :disabled="picking || i === listAcc.length - 1"
              @click="move(i, 1)"
            >
              下
            </button>
            <button
              v-if="item.status === 'error'"
              type="button"
              class="op-btn"
              :disabled="picking"
              @click="retry(item)"
            >
              重试
            </button>
            <button type="button" class="op-btn danger" :disabled="picking" @click="removeOne(item)">删</button>
          </div>
          <div v-else class="gallery-ops always" @click.stop>
            <label class="op-btn">
              选图
              <input
                type="file"
                accept="image/*"
                hidden
                :disabled="picking"
                @change="onReplaceFile(item, $event)"
              />
            </label>
            <button type="button" class="op-btn danger" :disabled="picking" @click="removeOne(item)">删</button>
          </div>
        </div>

        <span
          v-if="dualStatusLabel(item)"
          class="gallery-badge"
          :class="dualStatusTone(item) ? `tone-${dualStatusTone(item)}` : ''"
        >
          {{ dualStatusLabel(item) }}
        </span>
      </div>
    </div>

    <input
      ref="emptyPickInput"
      type="file"
      accept="image/*"
      hidden
      :disabled="picking"
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
import { onBeforeUnmount, onMounted, onUnmounted, ref, watch } from 'vue'
import { revokeObjectUrl } from '@/utils/imageProcess'
import {
  createEmptyDual,
  dualStatusLabel,
  dualStatusTone,
  SLOT_STATUS,
} from '@/utils/dualSlot'
import {
  cancelSlotUpload,
  enqueueFileUpload,
  enqueuePendingUpload,
  getSharedCompressQueue,
  getSharedUploadQueue,
} from '@/utils/imageUploadQueue'

const props = defineProps({
  modelValue: { type: Array, default: () => [] },
  title: { type: String, required: true },
  category: { type: String, default: 'portfolio' },
})

const emit = defineEmits(['update:modelValue'])

const picking = ref(false)
const error = ref('')
const selected = ref(new Set())
const blobUrlCache = new Map()
const lightboxOpen = ref(false)
const lightboxUrl = ref('')
const emptyPickInput = ref(null)
const emptyPickSlotId = ref(null)

/** Latest gallery list; patches always apply here to avoid stale-props races. */
const listAcc = ref([])
let syncingFromProps = false

watch(
  () => props.modelValue,
  (v) => {
    if (syncingFromProps) return
    listAcc.value = Array.isArray(v) ? v.slice() : []
  },
  { immediate: true },
)

function getBlobUrl(blob) {
  if (!blob) return ''
  let url = blobUrlCache.get(blob)
  if (!url) {
    url = URL.createObjectURL(blob)
    blobUrlCache.set(blob, url)
  }
  return url
}

function revokeLocalPreview(url) {
  if (url) revokeObjectUrl(url)
}

function cellPreview(item) {
  if (!item) return ''
  if (item.thumb) return item.thumb
  if (item.original) return item.original
  if (item.pending?.thumb) return getBlobUrl(item.pending.thumb)
  if (item.localPreview) return item.localPreview
  return ''
}

function cellLightbox(item) {
  if (!item) return ''
  if (item.pending?.original) return getBlobUrl(item.pending.original)
  if (item.pending?.thumb) return getBlobUrl(item.pending.thumb)
  if (item.localPreview) return item.localPreview
  return item.original || item.thumb || ''
}

function isBusy(item) {
  return (
    item?.status === SLOT_STATUS.COMPRESSING ||
    item?.status === SLOT_STATUS.QUEUED ||
    item?.status === SLOT_STATUS.UPLOADING
  )
}

function clearSelection() {
  selected.value = new Set()
}

function selectAll() {
  selected.value = new Set(listAcc.value.map((x, i) => x.slotId || i))
}

function toggleSelect(id, e) {
  const next = new Set(selected.value)
  if (e.target.checked) next.add(id)
  else next.delete(id)
  selected.value = next
}

function setList(next) {
  listAcc.value = next
  syncingFromProps = true
  emit('update:modelValue', next)
  queueMicrotask(() => {
    syncingFromProps = false
  })
}

function patchSlot(slotId, partial) {
  const list = listAcc.value.slice()
  const i = list.findIndex((x) => x.slotId === slotId)
  if (i < 0) return
  const prev = list[i]
  if (
    Object.prototype.hasOwnProperty.call(partial, 'localPreview') &&
    prev.localPreview &&
    partial.localPreview !== prev.localPreview
  ) {
    revokeLocalPreview(prev.localPreview)
  }
  list[i] = { ...prev, ...partial }
  setList(list)
}

function add() {
  setList([...listAcc.value, createEmptyDual()])
}

function removeBySlotId(slotId) {
  const item = listAcc.value.find((x) => x.slotId === slotId)
  cancelSlotUpload(slotId)
  if (item?.localPreview) revokeLocalPreview(item.localPreview)
  setList(listAcc.value.filter((x) => x.slotId !== slotId))
  const next = new Set(selected.value)
  next.delete(slotId)
  selected.value = next
}

function removeOne(item) {
  if (!item?.slotId) return
  if (!confirm('确定移除该图片？')) return
  removeBySlotId(item.slotId)
}

function removeSelected() {
  const ids = [...selected.value]
  if (!ids.length) return
  if (!confirm(`确定删除所选 ${ids.length} 张？`)) return
  for (const id of ids) {
    const item = listAcc.value.find((x) => x.slotId === id)
    cancelSlotUpload(id)
    if (item?.localPreview) revokeLocalPreview(item.localPreview)
  }
  setList(listAcc.value.filter((x) => !ids.includes(x.slotId)))
  clearSelection()
}

function move(i, delta) {
  const j = i + delta
  if (j < 0 || j >= listAcc.value.length) return
  const next = listAcc.value.slice()
  const [row] = next.splice(i, 1)
  next.splice(j, 0, row)
  setList(next)
  clearSelection()
}

function startUpload(item, file) {
  if (!item?.slotId || !file) return
  enqueueFileUpload({
    slotId: item.slotId,
    file,
    category: props.category,
    prevLocalPreview: item.localPreview || '',
    patch: (partial) => patchSlot(item.slotId, partial),
  })
}

function retry(item) {
  if (!item?.slotId) return
  if (item.pending?.original && item.pending?.thumb) {
    enqueuePendingUpload({
      slotId: item.slotId,
      pending: item.pending,
      category: props.category,
      patch: (partial) => patchSlot(item.slotId, partial),
    })
    return
  }
  error.value = '请重新选择图片'
}

function onReplaceFile(item, e) {
  const file = e.target.files?.[0]
  e.target.value = ''
  if (!file || !item) return
  startUpload(item, file)
}

function onThumbClick(i, item) {
  if (picking.value) return
  if (!cellPreview(item)) {
    if (isBusy(item)) return
    emptyPickSlotId.value = item.slotId
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

function onEmptyPickFile(e) {
  const file = e.target.files?.[0]
  e.target.value = ''
  const slotId = emptyPickSlotId.value
  emptyPickSlotId.value = null
  if (!slotId || !file) return
  const item = listAcc.value.find((x) => x.slotId === slotId)
  if (!item) return
  startUpload(item, file)
}

function onBatchFiles(e) {
  const files = Array.from(e.target.files || []).filter(
    (f) => f && typeof f.type === 'string' && f.type.startsWith('image/'),
  )
  e.target.value = ''
  if (!files.length) return

  picking.value = true
  error.value = ''
  try {
    const added = files.map((file) => {
      const dual = createEmptyDual()
      dual.localPreview = URL.createObjectURL(file)
      dual.status = SLOT_STATUS.QUEUED
      dual.error = ''
      return dual
    })
    setList([...listAcc.value, ...added])

    for (let i = 0; i < files.length; i++) {
      enqueueFileUpload({
        slotId: added[i].slotId,
        file: files[i],
        category: props.category,
        skipLocalPreviewSetup: true,
        patch: (partial) => patchSlot(added[i].slotId, partial),
      })
    }
  } catch (err) {
    error.value = err.message || '批量添加失败'
  } finally {
    picking.value = false
  }
}

onBeforeUnmount(() => {
  for (const item of listAcc.value || []) {
    if (item?.slotId) cancelSlotUpload(item.slotId)
    if (item?.localPreview) revokeLocalPreview(item.localPreview)
  }
  for (const url of blobUrlCache.values()) revokeObjectUrl(url)
  blobUrlCache.clear()
})

getSharedCompressQueue()
getSharedUploadQueue()
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
  color: var(--text-muted);
  line-height: 1.3;
  word-break: break-all;
}

.gallery-badge.tone-wait {
  color: #8b9bb4;
}

.gallery-badge.tone-busy {
  color: var(--accent, #6ea8fe);
}

.gallery-badge.tone-done {
  color: #3dd68c;
}

.gallery-badge.tone-error {
  color: var(--danger, #ff8f7a);
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
