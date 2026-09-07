<template>
  <div class="pair-list">
    <div class="pair-list-head">
      <h2>{{ title }}</h2>
      <button type="button" class="btn btn-sm" @click="add">添加一组</button>
    </div>

    <div v-if="!modelValue.length" class="muted empty-list">暂无对比组</div>

    <div v-for="(pair, i) in modelValue" :key="i" class="pair-item">
      <div class="pair-fields">
        <DualImageField
          :model-value="pair.render"
          :label="`效果图 #${i + 1}（可选）`"
          hint="至少上传效果图或实景图一侧"
          @update:model-value="(v) => updateField(i, 'render', v)"
        />
        <DualImageField
          :model-value="pair.real"
          :label="`实景图 #${i + 1}（可选）`"
          hint=""
          @update:model-value="(v) => updateField(i, 'real', v)"
        />
      </div>
      <div class="pair-ops">
        <button type="button" class="btn btn-sm btn-ghost" :disabled="i === 0" @click="move(i, -1)">
          上移
        </button>
        <button
          type="button"
          class="btn btn-sm btn-ghost"
          :disabled="i === modelValue.length - 1"
          @click="move(i, 1)"
        >
          下移
        </button>
        <button type="button" class="btn btn-sm btn-danger" @click="remove(i)">移除</button>
      </div>
    </div>
  </div>
</template>

<script setup>
import DualImageField from '@/components/DualImageField.vue'

const props = defineProps({
  modelValue: { type: Array, default: () => [] },
  title: { type: String, default: '效果图 / 实景图' },
})

const emit = defineEmits(['update:modelValue'])

function emptyDual() {
  return { thumb: '', original: '', pending: null }
}

function emptyPair() {
  return { render: emptyDual(), real: emptyDual() }
}

function updateField(i, key, val) {
  const next = props.modelValue.slice()
  next[i] = { ...next[i], [key]: val }
  emit('update:modelValue', next)
}

function add() {
  emit('update:modelValue', [...props.modelValue, emptyPair()])
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
</script>

<style scoped>
.pair-list {
  margin-top: 0.5rem;
  padding-top: 1rem;
  border-top: 1px solid var(--border);
}

.pair-list-head {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 0.75rem;
}

.pair-list-head h2 {
  margin: 0;
  font-family: var(--font-display);
  font-size: 1.05rem;
  font-weight: 560;
}

.pair-item {
  display: grid;
  grid-template-columns: 1fr auto;
  gap: 0.75rem;
  align-items: start;
  margin-bottom: 0.9rem;
  padding-bottom: 0.9rem;
  border-bottom: 1px solid var(--border);
}

.pair-item:last-child {
  border-bottom: none;
  padding-bottom: 0;
}

.pair-fields {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 0.75rem;
}

.pair-ops {
  display: flex;
  flex-direction: column;
  gap: 0.35rem;
}

.empty-list {
  margin: 0;
  font-size: 0.9rem;
}

@media (max-width: 900px) {
  .pair-fields {
    grid-template-columns: 1fr;
  }
}

@media (max-width: 700px) {
  .pair-item {
    grid-template-columns: 1fr;
  }

  .pair-ops {
    flex-direction: row;
    flex-wrap: wrap;
  }
}
</style>
