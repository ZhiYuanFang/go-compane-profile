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
      <span v-if="statusText" class="dual-status" :class="{ pending: isBusyStatus, error: modelValue?.status === 'error' }">
        {{ statusText }}
      </span>
      <span v-else-if="dragOver" class="dual-status pending">松开以添加</span>
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
          :disabled="disabled || isBusyStatus"
          @click="clear"
        >
          清除
        </button>
      </div>
    </div>

    <p v-if="error" class="dual-error">{{ error }}</p>
    <p v-if="hint" class="dual-hint">{{ hint }}</p>

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
import { processDualImage, revokeObjectUrl } from '@/utils/imageProcess'
import { createEmptyDual, dualStatusLabel, SLOT_STATUS } from '@/utils/dualSlot'
import { cancelSlotUpload, enqueuePendingUpload } from '@/utils/imageUploadQueue'

const props = defineProps({
  modelValue: {
    type: Object,
    default: () => createEmptyDual(),
  },
  label: { type: String, default: '图片' },
  hint: {
    type: String,
    default: '可拖拽图片到此处；选择后自动压缩并上传（原图≤3MB，缩略≤0.5MB）',
  },
  disabled: { type: Boolean, default: false },
  category: { type: String, default: 'general' },
  eagerUpload: { type: Boolean, default: true },
})

const emit = defineEmits(['update:modelValue'])

const error = ref('')
const localPreview = ref('')
const localOriginalPreview = ref('')
const pendingPreview = ref('')
const pendingOriginalPreview = ref('')
const lightboxOpen = ref(false)
const dragOver = ref(false)
const dragDepth = ref(0)
const fileInput = ref(null)

const isBusyStatus = computed(() => {
  const s = props.modelValue?.status
  return s === SLOT_STATUS.COMPRESSING || s === SLOT_STATUS.QUEUED || s === SLOT_STATUS.UPLOADING
})

const statusText = computed(() => dualStatusLabel(props.modelValue))

const previewUrl = computed(() => {
  if (localPreview.value) return localPreview.value
  if (pendingPreview.value) return pendingPreview.value
  const v = props.modelValue
  if (!v) return ''
  return v.thumb || v.original || ''
})

const lightboxUrl = computed(() => {
  if (localOriginalPreview.value) return localOriginalPreview.value
  if (pendingOriginalPreview.value) return pendingOriginalPreview.value
  const v = props.modelValue
  if (!v) return previewUrl.value
  return v.original || v.thumb || previewUrl.value
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
    // External pending (e.g. batch import) needs blob previews; local applyFile already sets localPreview.
    if (!pending) {
      revokePendingPreviews()
      return
    }
    if (localPreview.value) return
    syncPendingPreviews(pending)
  },
  { immediate: true },
)

watch(
  () => props.modelValue,
  (v) => {
    if (!v?.pending && localPreview.value) {
      revokeObjectUrl(localPreview.value)
      localPreview.value = ''
      revokeObjectUrl(localOriginalPreview.value)
      localOriginalPreview.value = ''
    }
  },
  { deep: true },
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
  const id = props.modelValue?.slotId
  if (id) cancelSlotUpload(id)
  revokeObjectUrl(localPreview.value)
  revokeObjectUrl(localOriginalPreview.value)
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
  emitValue({ ...props.modelValue, ...partial })
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

async function applyFile(file) {
  if (!file) return
  if (props.disabled) return

  error.value = ''
  const slotId = ensureSlotId()
  cancelSlotUpload(slotId)

  if (props.eagerUpload) {
    emitValue({
      slotId,
      status: SLOT_STATUS.COMPRESSING,
      error: '',
    })
    try {
      const pair = await processDualImage(file)
      revokeObjectUrl(localPreview.value)
      revokeObjectUrl(localOriginalPreview.value)
      localPreview.value = URL.createObjectURL(pair.thumb)
      localOriginalPreview.value = URL.createObjectURL(pair.original)
      emitValue({
        slotId,
        thumb: props.modelValue?.thumb || '',
        original: props.modelValue?.original || '',
        pending: pair,
        status: SLOT_STATUS.QUEUED,
        error: '',
      })
      enqueuePendingUpload({
        slotId,
        pending: pair,
        category: props.category,
        patch: (partial) => patchSelf(partial),
      })
    } catch (err) {
      error.value = err.message || '图片处理失败'
      emitValue({
        slotId,
        status: SLOT_STATUS.ERROR,
        error: err.message || '图片处理失败',
      })
    }
    return
  }

  try {
    const pair = await processDualImage(file)
    revokeObjectUrl(localPreview.value)
    revokeObjectUrl(localOriginalPreview.value)
    localPreview.value = URL.createObjectURL(pair.thumb)
    localOriginalPreview.value = URL.createObjectURL(pair.original)
    emitValue({
      slotId,
      thumb: props.modelValue?.thumb || '',
      original: props.modelValue?.original || '',
      pending: pair,
      status: SLOT_STATUS.QUEUED,
      error: '',
    })
  } catch (err) {
    error.value = err.message || '图片处理失败'
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
    })
    return
  }
  openFilePicker()
}

async function onFile(e) {
  const file = e.target.files?.[0]
  e.target.value = ''
  if (!file) return
  await applyFile(file)
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

  await applyFile(file)
}

function clear() {
  const id = props.modelValue?.slotId
  if (id) cancelSlotUpload(id)
  revokeObjectUrl(localPreview.value)
  revokeObjectUrl(localOriginalPreview.value)
  localPreview.value = ''
  localOriginalPreview.value = ''
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

.dual-status.error {
  color: var(--danger);
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
