<template>
  <div class="rte">
    <div class="rte-toolbar">
      <div class="rte-row">
        <button type="button" class="btn btn-sm" @mousedown.prevent="exec('bold')">加粗</button>
        <label class="rte-color" @mousedown="saveSelection">
          颜色
          <input type="color" :value="color" @input="onColor" />
        </label>
      </div>
      <div class="rte-row rte-size-row">
        <span class="rte-hint">需先选中文字再调节</span>
        <button
          v-for="opt in sizeOptions"
          :key="opt.value"
          type="button"
          class="rte-size-btn"
          :class="{ active: fontSize === opt.value }"
          @mousedown.prevent="applyFontSize(opt.value)"
        >
          {{ opt.label }}
        </button>
      </div>
    </div>
    <div
      ref="editor"
      class="rte-body"
      contenteditable="true"
      @input="onInput"
      @blur="onInput"
      @mouseup="saveSelection"
      @keyup="saveSelection"
    ></div>
  </div>
</template>

<script setup>
import { onMounted, ref, watch } from 'vue'

const props = defineProps({
  modelValue: { type: String, default: '' },
})
const emit = defineEmits(['update:modelValue'])

const sizeOptions = [
  { label: '小', value: '12px' },
  { label: '默认', value: '14px' },
  { label: '中', value: '16px' },
  { label: '大', value: '18px' },
  { label: '更大', value: '22px' },
]

const editor = ref(null)
const color = ref('#1a1a1a')
const fontSize = ref('14px')
let applying = false
let savedRange = null

onMounted(() => {
  if (editor.value) {
    applying = true
    editor.value.innerHTML = props.modelValue || ''
    applying = false
  }
})

watch(
  () => props.modelValue,
  (v) => {
    if (!editor.value || applying) return
    if (editor.value.innerHTML === (v || '')) return
    applying = true
    editor.value.innerHTML = v || ''
    applying = false
  },
)

function onInput() {
  if (!editor.value || applying) return
  emit('update:modelValue', editor.value.innerHTML)
}

function saveSelection() {
  const sel = window.getSelection()
  if (!sel || sel.rangeCount === 0) return
  const range = sel.getRangeAt(0)
  if (!editor.value?.contains(range.commonAncestorContainer)) return
  savedRange = range.cloneRange()
}

function restoreSelection() {
  if (!savedRange || !editor.value) return false
  const sel = window.getSelection()
  if (!sel) return false
  sel.removeAllRanges()
  sel.addRange(savedRange)
  return true
}

function exec(cmd) {
  restoreSelection()
  document.execCommand(cmd, false)
  saveSelection()
  onInput()
}

function onColor(e) {
  color.value = e.target.value
  restoreSelection()
  document.execCommand('foreColor', false, color.value)
  // Convert <font color> to span style so sanitizer keeps color
  editor.value?.querySelectorAll('font[color]').forEach((el) => {
    const span = document.createElement('span')
    span.style.color = el.getAttribute('color') || color.value
    while (el.firstChild) span.appendChild(el.firstChild)
    el.replaceWith(span)
  })
  saveSelection()
  onInput()
}

function applyFontSize(size) {
  fontSize.value = size
  restoreSelection()
  const sel = window.getSelection()
  if (!sel || sel.rangeCount === 0 || sel.isCollapsed) return
  if (!editor.value?.contains(sel.getRangeAt(0).commonAncestorContainer)) return

  const range = sel.getRangeAt(0)
  const span = document.createElement('span')
  span.style.fontSize = size
  try {
    range.surroundContents(span)
  } catch {
    const frag = range.extractContents()
    span.appendChild(frag)
    range.insertNode(span)
  }
  sel.removeAllRanges()
  const next = document.createRange()
  next.selectNodeContents(span)
  sel.addRange(next)
  saveSelection()
  onInput()
}
</script>

<style scoped>
.rte {
  border: 1px solid var(--border);
  border-radius: 10px;
  overflow: hidden;
  background: rgba(0, 0, 0, 0.2);
}

.rte-toolbar {
  display: flex;
  flex-direction: column;
  gap: 0.45rem;
  padding: 0.55rem 0.75rem;
  border-bottom: 1px solid var(--border);
}

.rte-row {
  display: flex;
  flex-wrap: wrap;
  gap: 0.45rem;
  align-items: center;
}

.rte-color {
  display: inline-flex;
  align-items: center;
  gap: 0.35rem;
  font-size: 0.85rem;
  color: var(--text-muted);
}

.rte-color input {
  width: 2rem;
  height: 1.5rem;
  padding: 0;
  border: none;
  background: transparent;
}

.rte-hint {
  font-size: 0.78rem;
  color: var(--text-muted);
  margin-right: 0.15rem;
  white-space: nowrap;
}

.rte-size-btn {
  box-sizing: border-box;
  width: 2.6em;
  padding: 0.28rem 0;
  border-radius: 6px;
  border: 1px solid var(--border);
  background: rgba(0, 0, 0, 0.28);
  color: var(--text);
  font-size: 0.8rem;
  line-height: 1.2;
  text-align: center;
  cursor: pointer;
}

.rte-size-btn.active {
  border-color: var(--accent, #6ea8fe);
  color: var(--accent, #6ea8fe);
}

.rte-body {
  min-height: 160px;
  padding: 0.85rem 1rem;
  outline: none;
  line-height: 1.6;
}
</style>
