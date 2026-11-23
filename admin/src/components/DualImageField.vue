<template>
  <div
    class="dual"
    :class="{ 'is-dragover': dragOver, disabled: disabled }"
    @dragenter.prevent="onDragEnter"
    @dragover.prevent="onDragOver"
    @dragleave.prevent="onDragLeave"
    @drop.prevent="onDrop"
  >
    <div class="dual-head">
      <span class="dual-label">{{ label }}</span>
      <span
        v-if="statusText"
        class="dual-status"
        :class="statusToneClass"
      >
        {{ statusText }}
      </span>
      <span v-else-if="dragOver" class="dual-status pending">松开以添加</span>
      <button
        v-if="canWatermark"
        type="button"
        class="btn btn-sm btn-ghost dual-watermark"
        :disabled="watermarking"
        title="给已上线或待保存图片补加水印，保存时上传替换（重复添加会叠加）"
        @click="addWatermark"
      >
        {{ watermarking ? '加水印中…' : '加水印' }}
      </button>
    </div>

    <div class="dual-body">
      <div
        class="preview"
        :class="{ empty: !previewUrl, clickable: true }"
        role="button"
        tabindex="0"
        @click="onPreviewClick"
        @keydown.enter.prevent="onPreviewClick"
      >
        <img v-if="previewUrl" :src="previewUrl" :alt="label" />
        <span v-else class="placeholder">拖拽图片到此处，或点击选择</span>
      </div>

      <div class="dual-actions">
        <label class="btn btn-sm">
          {{ previewUrl ? '重新选择' : '选择图片' }}
          <input
            ref="fileInput"
            type="file"
            accept="image/*"
            hidden
            :disabled="disabled"
            @change="onFile"
          />
        </label>
        <button
          v-if="modelValue?.status === 'error'"
          type="button"
          class="btn btn-sm"
          :disabled="disabled || isBusyStatus"
          @click="retry"
        >
          重试
        </button>
        <button
          v-if="previewUrl"
          type="button"
          class="btn btn-sm btn-ghost"
          :disabled="disabled"
          @click="clear"
        >
          清除
        </button>
      </div>
    </div>

    <p v-if="error" class="dual-error">{{ error }}</p>
    <p v-if="hintText" class="dual-hint">{{ hintText }}</p>

    <Teleport to="body">
      <div
        v-if="lightboxOpen && lightboxUrl"
        class="lightbox"
        role="dialog"
        aria-modal="true"
        @click.self="closeLightbox"
        @keydown.esc.prevent="closeLightbox"
      >
        <button type="button" class="lightbox-close btn btn-sm" @click="closeLightbox">关闭</button>
        <img class="lightbox-img" :src="lightboxUrl" :alt="label" @click.stop />
      </div>
    </Teleport>
  </div>
</template>

<script setup>
import { computed, onBeforeUnmount, onMounted, onUnmounted, ref, watch } from 'vue'
import { revokeObjectUrl, WATERMARK_TEXT } from '@/utils/imageProcess'
import { createPreviewObjectUrl } from '@/utils/previewImage'
import { blobAsImageFile, fetchRemoteImageFile } from '@/utils/remoteImage'
import { createEmptyDual, dualStatusLabel, dualStatusTone, SLOT_STATUS } from '@/utils/dualSlot'
import {
  cancelSlotUpload,
  enqueuePendingUpload,
  startSlotPipeline,
} from '@/utils/imageUploadQueue'

const props = defineProps({
  modelValue: {
    type: Object,
    default: () => createEmptyDual(),
  },
  label: { type: String, default: '图片' },
  hint: { type: String, default: '' },
  disabled: { type: Boolean, default: false },
  category: { type: String, default: 'general' },
})

const emit = defineEmits(['update:modelValue'])

const error = ref('')
const watermarking = ref(false)
const localPreview = ref('')
const localOriginalPreview = ref('')
const pendingPreview = ref('')
const pendingOriginalPreview = ref('')
const lightboxOpen = ref(false)
const dragOver = ref(false)
const dragDepth = ref(0)
const fileInput = ref(null)
let previewToken = 0

const isBusyStatus = computed(() => {
  const s = props.modelValue?.status
  return (
    s === SLOT_STATUS.DOWNLOADING ||
    s === SLOT_STATUS.COMPRESSING ||
    s === SLOT_STATUS.WATERMARKING ||
    s === SLOT_STATUS.QUEUED ||
    s === SLOT_STATUS.UPLOADING
  )
})

/** Already-uploaded or clean-pending content images can take an explicit watermark. */
const canWatermark = computed(() => {
  if (props.disabled || props.category === 'logo' || isBusyStatus.value) return false
  const v = props.modelValue
  if (!v) return false
  if (v.pending?.original || v.pending?.thumb) return true
  return !v.pending && !!(v.original || v.thumb)
})

const statusText = computed(() => dualStatusLabel(props.modelValue))

const hintText = computed(() => {
  if (props.hint) return props.hint
  const limits = '（原图≤3MB，缩略≤0.5MB）'
  if (props.category === 'logo') {
    return `可拖拽图片到此处；选择后本地压缩，保存时上传${limits}`
  }
  return `可拖拽图片到此处；选择后本地压缩，保存时上传；需要时点「加水印」添加 ${WATERMARK_TEXT}${limits}`
})

const statusToneClass = computed(() => {
  const tone = dualStatusTone(props.modelValue)
  return tone ? `tone-${tone}` : ''
})

const previewUrl = computed(() => {
  if (pendingPreview.value) return pendingPreview.value
  if (localPreview.value) return localPreview.value
  const v = props.modelValue
  if (!v) return ''
  return v.thumb || v.original || v.localPreview || ''
})

const lightboxUrl = computed(() => {
  if (pendingOriginalPreview.value) return pendingOriginalPreview.value
  if (pendingPreview.value) return pendingPreview.value
  if (localOriginalPreview.value) return localOriginalPreview.value
  const v = props.modelValue
  if (!v) return previewUrl.value
  return v.original || v.thumb || v.localPreview || previewUrl.value
})

function revokePendingPreviews() {
  revokeObjectUrl(pendingPreview.value)
  revokeObjectUrl(pendingOriginalPreview.value)
  pendingPreview.value = ''
  pendingOriginalPreview.value = ''
}

function syncPendingPreviews(pending) {
  revokePendingPreviews()
  if (!pending?.thumb) return
  pendingPreview.value = URL.createObjectURL(pending.thumb)
  if (pending.original) {
    pendingOriginalPreview.value = URL.createObjectURL(pending.original)
  }
}

watch(
  () => props.modelValue?.pending,
  (pending) => {
    if (!pending) {
      revokePendingPreviews()
      return
    }
    // Prefer pending (compressed, optionally watermarked) over the raw local preview.
    syncPendingPreviews(pending)
  },
  { immediate: true },
)

watch(
  () => props.modelValue?.status,
  (status) => {
    if (status === SLOT_STATUS.DONE && (props.modelValue?.thumb || props.modelValue?.original)) {
      revokeObjectUrl(localPreview.value)
      localPreview.value = ''
      revokeObjectUrl(localOriginalPreview.value)
      localOriginalPreview.value = ''
      revokePendingPreviews()
    }
  },
)

function onKeydown(e) {
  if (e.key === 'Escape' && lightboxOpen.value) {
    closeLightbox()
  }
}

onMounted(() => {
  window.addEventListener('keydown', onKeydown)
})

onUnmounted(() => {
  window.removeEventListener('keydown', onKeydown)
})

onBeforeUnmount(() => {
  previewToken += 1
  const id = props.modelValue?.slotId
  if (id) cancelSlotUpload(id)
  revokeObjectUrl(localPreview.value)
  revokeObjectUrl(localOriginalPreview.value)
  if (props.modelValue?.localPreview) revokeObjectUrl(props.modelValue.localPreview)
  revokePendingPreviews()
})

function emitValue(next) {
  const base = createEmptyDual(props.modelValue)
  emit('update:modelValue', { ...base, ...next })
}

function ensureSlotId() {
  if (props.modelValue?.slotId) return props.modelValue.slotId
  const dual = createEmptyDual(props.modelValue)
  emit('update:modelValue', dual)
  return dual.slotId
}

function patchSelf(partial) {
  const prev = props.modelValue
  if (
    Object.prototype.hasOwnProperty.call(partial, 'localPreview') &&
    prev?.localPreview &&
    partial.localPreview !== prev.localPreview
  ) {
    revokeObjectUrl(prev.localPreview)
  }
  emitValue({ ...prev, ...partial })
}

function openLightbox() {
  if (!previewUrl.value) return
  lightboxOpen.value = true
}

function closeLightbox() {
  lightboxOpen.value = false
}

function openFilePicker() {
  if (props.disabled) return
  fileInput.value?.click()
}

function onPreviewClick() {
  if (props.disabled) return
  if (!previewUrl.value) {
    openFilePicker()
    return
  }
  openLightbox()
}

function firstImageFile(fileList) {
  if (!fileList?.length) return null
  for (const file of fileList) {
    if (file && typeof file.type === 'string' && file.type.startsWith('image/')) {
      return file
    }
  }
  return null
}

function applyFile(file, { watermark = false } = {}) {
  if (!file) return
  if (props.disabled) return

  error.value = ''
  const slotId = ensureSlotId()

  revokeObjectUrl(localPreview.value)
  revokeObjectUrl(localOriginalPreview.value)
  localPreview.value = ''
  localOriginalPreview.value = ''

  // Small preview only: a raw full-size photo can be tens of MiB in memory.
  const token = ++previewToken
  createPreviewObjectUrl(file)
    .then((url) => {
      if (token !== previewToken) {
        revokeObjectUrl(url)
        return
      }
      revokeObjectUrl(localPreview.value)
      revokeObjectUrl(localOriginalPreview.value)
      localPreview.value = url
      localOriginalPreview.value = url
    })
    .catch(() => {
      /* pending thumb will show when compression finishes */
    })

  startSlotPipeline({
    slotId,
    file,
    category: props.category,
    watermark,
    skipLocalPreviewSetup: true,
    patch: (partial) => {
      patchSelf(partial)
      if (partial.error) error.value = partial.error
      if (partial.status === SLOT_STATUS.DONE) error.value = ''
    },
  })
}

async function addWatermark() {
  if (!canWatermark.value || watermarking.value) return
  watermarking.value = true
  error.value = ''

  const pending = props.modelValue?.pending
  if (pending?.original || pending?.thumb) {
    try {
      const file = blobAsImageFile(pending.original || pending.thumb, 'pending')
      if (props.disabled) return
      applyFile(file, { watermark: true })
    } catch (err) {
      error.value = err.message || '加水印失败'
    } finally {
      watermarking.value = false
    }
    return
  }

  const sourceUrl = props.modelValue.original || props.modelValue.thumb
  if (!sourceUrl) {
    watermarking.value = false
    return
  }

  patchSelf({ status: SLOT_STATUS.DOWNLOADING, error: '' })
  try {
    const file = await fetchRemoteImageFile(sourceUrl)
    if ((props.modelValue?.original || props.modelValue?.thumb) !== sourceUrl) return
    if (props.disabled) {
      patchSelf({ status: SLOT_STATUS.DONE, error: '' })
      return
    }
    applyFile(file, { watermark: true })
  } catch (err) {
    if ((props.modelValue?.original || props.modelValue?.thumb) === sourceUrl) {
      patchSelf({ status: SLOT_STATUS.DONE, error: '' })
    }
    error.value = err.message || '加水印失败'
  } finally {
    watermarking.value = false
  }
}

function retry() {
  const slotId = props.modelValue?.slotId || ensureSlotId()
  const pending = props.modelValue?.pending
  if (pending?.original && pending?.thumb) {
    enqueuePendingUpload({
      slotId,
      pending,
      category: props.category,
      patch: (partial) => patchSelf(partial),
    }).catch(() => {})
    return
  }
  openFilePicker()
}

async function onFile(e) {
  const file = e.target.files?.[0]
  e.target.value = ''
  if (!file) return
  applyFile(file)
}

function onDragEnter() {
  if (props.disabled) return
  dragDepth.value += 1
  dragOver.value = true
}

function onDragOver() {
  if (props.disabled) return
  dragOver.value = true
}

function onDragLeave() {
  dragDepth.value = Math.max(0, dragDepth.value - 1)
  if (dragDepth.value === 0) {
    dragOver.value = false
  }
}

async function onDrop(e) {
  dragDepth.value = 0
  dragOver.value = false
  if (props.disabled) return

  const file = firstImageFile(e.dataTransfer?.files)
  if (!file) {
    error.value = '请拖入图片文件'
    return
  }

  if (previewUrl.value) {
    const ok = window.confirm('已有图片，确定要替换吗？')
    if (!ok) return
  }

  applyFile(file)
}

function clear() {
  const id = props.modelValue?.slotId
  if (id) cancelSlotUpload(id)
  revokeObjectUrl(localPreview.value)
  revokeObjectUrl(localOriginalPreview.value)
  localPreview.value = ''
  localOriginalPreview.value = ''
  if (props.modelValue?.localPreview) revokeObjectUrl(props.modelValue.localPreview)
  revokePendingPreviews()
  error.value = ''
  lightboxOpen.value = false
  emitValue(createEmptyDual())
}
</script>

<style scoped>
.dual {
  padding: 0.9rem;
  border-radius: var(--radius-sm);
  border: 1px solid var(--border);
  background: rgba(0, 0, 0, 0.2);
  transition: border-color 0.15s ease, background 0.15s ease;
}

.dual.is-dragover {
  border-color: var(--accent);
  background: rgba(255, 255, 255, 0.06);
}

.dual.disabled {
  opacity: 0.7;
}

.dual-head {
  display: flex;
  align-items: center;
  gap: 0.6rem;
  margin-bottom: 0.75rem;
}

.dual-label {
  font-size: 0.82rem;
  letter-spacing: 0.04em;
  text-transform: uppercase;
  color: var(--text-muted);
}

.dual-status {
  font-size: 0.75rem;
  color: var(--text-muted);
}

.dual-status.pending {
  color: var(--accent);
}

.dual-watermark {
  margin-left: auto;
}

.dual-status.tone-wait {
  color: #8b9bb4;
}

.dual-status.tone-busy {
  color: var(--accent, #6ea8fe);
}

.dual-status.tone-done {
  color: #3dd68c;
}

.dual-status.tone-error {
  color: var(--danger, #ff8f7a);
}

.dual-body {
  display: flex;
  gap: 1rem;
  align-items: flex-end;
  flex-wrap: wrap;
}

.preview {
  width: 140px;
  height: 100px;
  border-radius: 8px;
  overflow: hidden;
  border: 1px solid var(--border);
  background: rgba(255, 255, 255, 0.03);
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: zoom-in;
}

.preview.empty {
  border-style: dashed;
  cursor: pointer;
}

.preview img {
  width: 100%;
  height: 100%;
  object-fit: contain;
}

.placeholder {
  font-size: 0.72rem;
  line-height: 1.35;
  padding: 0 0.5rem;
  text-align: center;
  color: var(--text-muted);
}

.dual-actions {
  display: flex;
  gap: 0.5rem;
  flex-wrap: wrap;
}

.dual-error {
  margin: 0.6rem 0 0;
  color: var(--danger);
  font-size: 0.85rem;
}

.dual-hint {
  margin: 0.55rem 0 0;
  color: var(--text-muted);
  font-size: 0.78rem;
}

.lightbox {
  position: fixed;
  inset: 0;
  z-index: 9999;
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
</style>
