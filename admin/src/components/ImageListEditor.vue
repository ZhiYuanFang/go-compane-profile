<template>
  <div class="image-list">
    <div class="image-list-head">
      <h2>{{ title }}</h2>
      <div class="image-list-head-actions">
        <label class="btn btn-sm" :class="{ disabled: locked }">
          {{ picking ? '添加中…' : '批量选择' }}
          <input
            type="file"
            accept="image/*"
            multiple
            hidden
            :disabled="locked"
            @change="onBatchFiles"
          />
        </label>
        <button type="button" class="btn btn-sm btn-ghost" :disabled="locked" @click="add">
          添加一张
        </button>
        <button
          type="button"
          class="btn btn-sm btn-ghost"
          :disabled="locked || watermarking || !selectedWatermarkTargets.length"
          title="给勾选的已上线图片补加水印，保存时上传替换（重复添加会叠加）"
          @click="watermarkSelected"
        >
          {{ watermarking ? '处理中…' : '加水印' }}
        </button>
        <button type="button" class="btn btn-sm btn-ghost" :disabled="locked || !listAcc.length" @click="selectAll">
          全选
        </button>
        <button type="button" class="btn btn-sm btn-ghost" :disabled="locked || !selected.size" @click="clearSelection">
          取消全选
        </button>
        <button
          type="button"
          class="btn btn-sm btn-danger"
          :disabled="locked || !selected.size"
          @click="removeSelected"
        >
          删除所选{{ selected.size ? ` (${selected.size})` : '' }}
        </button>
      </div>
    </div>

    <p v-if="error" class="image-list-error">{{ error }}</p>
    <div v-if="!listAcc.length" class="muted empty-list">暂无图片，可批量选择（选后本地处理并加水印，保存时上传）</div>

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
            :disabled="locked"
            @change="toggleSelect(item.slotId || i, $event)"
          />
        </label>
        <span class="gallery-index">#{{ i + 1 }}</span>

        <div class="gallery-thumb-wrap">
          <button
            type="button"
            class="gallery-thumb"
            :class="{ empty: !cellPreview(item) }"
            :disabled="locked"
            @click="onThumbClick(i, item)"
          >
            <img
              v-if="cellPreview(item)"
              :src="cellPreview(item)"
              :alt="`#${i + 1}`"
              loading="lazy"
              decoding="async"
            />
            <span v-else class="gallery-placeholder">{{ isBusy(item) ? '处理中…' : '选择图片' }}</span>
          </button>

          <div v-if="cellPreview(item)" class="gallery-ops" @click.stop>
            <label class="op-btn">
              换
              <input
                type="file"
                accept="image/*"
                hidden
                :disabled="locked"
                @change="onReplaceFile(item, $event)"
              />
            </label>
            <button type="button" class="op-btn" :disabled="locked || i === 0" @click="move(i, -1)">上</button>
            <button
              type="button"
              class="op-btn"
              :disabled="locked || i === listAcc.length - 1"
              @click="move(i, 1)"
            >
              下
            </button>
            <button
              v-if="item.status === 'error'"
              type="button"
              class="op-btn"
              :disabled="locked"
              @click="retry(item)"
            >
              重试
            </button>
            <button type="button" class="op-btn danger" :disabled="locked" @click="removeOne(item)">删</button>
          </div>
          <div v-else class="gallery-ops always" @click.stop>
            <label class="op-btn">
              选图
              <input
                type="file"
                accept="image/*"
                hidden
                :disabled="locked"
                @change="onReplaceFile(item, $event)"
              />
            </label>
            <button type="button" class="op-btn danger" :disabled="locked" @click="removeOne(item)">删</button>
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
      :disabled="locked"
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
import { computed, nextTick, onBeforeUnmount, onMounted, onUnmounted, ref, watch } from 'vue'
import { revokeObjectUrl } from '@/utils/imageProcess'
import { mapPool } from '@/utils/mapPool'
import { createPreviewObjectUrl } from '@/utils/previewImage'
import { fetchRemoteImageFile } from '@/utils/remoteImage'
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
  disabled: { type: Boolean, default: false },
})

const emit = defineEmits(['update:modelValue'])

const picking = ref(false)
const watermarking = ref(false)
const locked = computed(() => props.disabled || picking.value)
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

function revokePendingCache(pending) {
  if (!pending) return
  for (const blob of [pending.thumb, pending.original]) {
    if (!blob) continue
    const url = blobUrlCache.get(blob)
    if (url) {
      revokeObjectUrl(url)
      blobUrlCache.delete(blob)
    }
  }
}

function cellPreview(item) {
  if (!item) return ''
  // Pending blobs are watermarked: prefer them over the raw local preview.
  if (item.pending?.thumb) return getBlobUrl(item.pending.thumb)
  if (item.localPreview) return item.localPreview
  if (item.thumb) return item.thumb
  if (item.original) return item.original
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
    item?.status === SLOT_STATUS.DOWNLOADING ||
    item?.status === SLOT_STATUS.COMPRESSING ||
    item?.status === SLOT_STATUS.WATERMARKING ||
    item?.status === SLOT_STATUS.QUEUED ||
    item?.status === SLOT_STATUS.UPLOADING
  )
}

/** Selected images that are still untouched remote uploads. */
const selectedWatermarkTargets = computed(() =>
  listAcc.value.filter(
    (item) =>
      selected.value.has(item.slotId) &&
      !item.pending &&
      !isBusy(item) &&
      !!(item.original || item.thumb),
  ),
)

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
  if (
    Object.prototype.hasOwnProperty.call(partial, 'pending') &&
    prev.pending &&
    partial.pending !== prev.pending
  ) {
    revokePendingCache(prev.pending)
  }
  list[i] = { ...prev, ...partial }
  setList(list)
}

function add() {
  if (locked.value) return
  setList([...listAcc.value, createEmptyDual()])
}

function removeBySlotId(slotId) {
  const item = listAcc.value.find((x) => x.slotId === slotId)
  cancelSlotUpload(slotId)
  if (item?.localPreview) revokeLocalPreview(item.localPreview)
  revokePendingCache(item?.pending)
  setList(listAcc.value.filter((x) => x.slotId !== slotId))
  const next = new Set(selected.value)
  next.delete(slotId)
  selected.value = next
}

function removeOne(item) {
  if (locked.value || !item?.slotId) return
  if (!confirm('确定移除该图片？')) return
  removeBySlotId(item.slotId)
}

function removeSelected() {
  if (locked.value) return
  const ids = [...selected.value]
  if (!ids.length) return
  if (!confirm(`确定删除所选 ${ids.length} 张？`)) return
  for (const id of ids) {
    const item = listAcc.value.find((x) => x.slotId === id)
    cancelSlotUpload(id)
    if (item?.localPreview) revokeLocalPreview(item.localPreview)
    revokePendingCache(item?.pending)
  }
  setList(listAcc.value.filter((x) => !ids.includes(x.slotId)))
  clearSelection()
}

function move(i, delta) {
  if (locked.value) return
  const j = i + delta
  if (j < 0 || j >= listAcc.value.length) return
  const next = listAcc.value.slice()
  const [row] = next.splice(i, 1)
  next.splice(j, 0, row)
  setList(next)
  clearSelection()
}

function startUpload(item, file) {
  if (locked.value || !item?.slotId || !file) return
  enqueueFileUpload({
    slotId: item.slotId,
    file,
    category: props.category,
    prevLocalPreview: item.localPreview || '',
    patch: (partial) => patchSlot(item.slotId, partial),
  })
}

const WATERMARK_FETCH_CONCURRENCY = 4

async function watermarkSelected() {
  if (locked.value || watermarking.value) return
  const targets = selectedWatermarkTargets.value.map((item) => ({
    item,
    index: listAcc.value.indexOf(item),
  }))
  if (!targets.length) return

  watermarking.value = true
  error.value = ''
  // The remote download can be slow; show it on every selected slot immediately.
  for (const { item } of targets) {
    const current = listAcc.value.find((x) => x.slotId === item.slotId)
    if (current) patchSlot(current.slotId, { status: SLOT_STATUS.DOWNLOADING, error: '' })
  }

  let failures = []
  try {
    failures = (
      await mapPool(targets, WATERMARK_FETCH_CONCURRENCY, async ({ item, index }) => {
        const before = listAcc.value.find((x) => x.slotId === item.slotId)
        if (!before || before.status !== SLOT_STATUS.DOWNLOADING) return ''
        try {
          const file = await fetchRemoteImageFile(before.original || before.thumb)
          // The slot may have been removed or replaced while its download was in flight.
          const current = listAcc.value.find((x) => x.slotId === item.slotId)
          if (!current || current.status !== SLOT_STATUS.DOWNLOADING || current.pending) return ''
          if (!current.original && !current.thumb) return ''
          startUpload(current, file)
          return ''
        } catch (err) {
          const current = listAcc.value.find((x) => x.slotId === item.slotId)
          if (!current) return ''
          if (current.status === SLOT_STATUS.DOWNLOADING) {
            // The remote image is still valid; return the slot to 已上传.
            patchSlot(current.slotId, { status: SLOT_STATUS.DONE, error: '' })
          }
          return `#${index + 1} ${err.message || '处理失败'}`
        }
      })
    ).filter(Boolean)
  } finally {
    watermarking.value = false
  }
  if (failures.length) error.value = failures.join('；')
}

function retry(item) {
  if (locked.value || !item?.slotId) return
  if (item.pending?.original && item.pending?.thumb) {
    enqueuePendingUpload({
      slotId: item.slotId,
      pending: item.pending,
      category: props.category,
      patch: (partial) => patchSlot(item.slotId, partial),
    }).catch(() => {})
    return
  }
  error.value = '请重新选择图片'
}

function onReplaceFile(item, e) {
  if (locked.value) return
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
  if (locked.value) return
  const file = e.target.files?.[0]
  e.target.value = ''
  const slotId = emptyPickSlotId.value
  emptyPickSlotId.value = null
  if (!slotId || !file) return
  const item = listAcc.value.find((x) => x.slotId === slotId)
  if (!item) return
  startUpload(item, file)
}

const BATCH_PREVIEW_CONCURRENCY = 4

async function onBatchFiles(e) {
  if (locked.value) return
  const files = Array.from(e.target.files || []).filter(
    (f) => f && typeof f.type === 'string' && f.type.startsWith('image/'),
  )
  e.target.value = ''
  if (!files.length) return

  picking.value = true
  error.value = ''
  try {
    const added = files.map(() => {
      const dual = createEmptyDual()
      dual.status = SLOT_STATUS.QUEUED
      dual.error = ''
      return dual
    })
    setList([...listAcc.value, ...added])
    await nextTick()

    // Build small previews first: full-size previews of a large batch exhaust
    // memory and stall the tab. Slots show「排队中」until their preview is ready.
    await mapPool(files, BATCH_PREVIEW_CONCURRENCY, async (file, i) => {
      try {
        const url = await createPreviewObjectUrl(file)
        if (!listAcc.value.some((x) => x.slotId === added[i].slotId)) {
          revokeObjectUrl(url)
          return
        }
        patchSlot(added[i].slotId, { localPreview: url })
      } catch {
        /* preview is optional; pending thumb will show when compression finishes */
      }
    })

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
  grid-template-columns: repeat(8, minmax(0, 1fr));
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
