<template>
  <div class="image-list">
    <div class="image-list-head">
      <h2>{{ title }}</h2>
      <button type="button" class="btn btn-sm" @click="add">添加</button>
    </div>

    <div v-if="!modelValue.length" class="muted empty-list">暂无图片</div>

    <div v-for="(item, i) in modelValue" :key="i" class="image-list-item">
      <DualImageField
        :model-value="item"
        :label="`#${i + 1}`"
        hint=""
        @update:model-value="(v) => update(i, v)"
      />
      <div class="image-list-ops">
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
  title: { type: String, required: true },
})

const emit = defineEmits(['update:modelValue'])

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
  margin-bottom: 0.75rem;
}

.image-list-head h2 {
  margin: 0;
  font-family: var(--font-display);
  font-size: 1.05rem;
  font-weight: 560;
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
