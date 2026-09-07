<template>
  <section>
    <h1 class="page-title">公司资料</h1>
    <p class="page-sub">编辑简介文案与品牌 Logo</p>

    <div v-if="loadError" class="alert alert-error">{{ loadError }}</div>
    <div v-if="message" class="alert" :class="messageOk ? 'alert-ok' : 'alert-error'">{{ message }}</div>

    <div v-if="!loading" class="glass-panel about-views">
      查阅关于我们次数为 {{ aboutViewCount }}
    </div>

    <form v-if="!loading" class="glass-panel form" @submit.prevent="onSubmit">
      <div class="field">
        <label for="introTitle">简介标题</label>
        <input id="introTitle" v-model="form.introTitle" />
      </div>

      <div class="field">
        <label for="introBody">简介正文</label>
        <textarea id="introBody" v-model="form.introBody" rows="8" />
      </div>

      <div class="field">
        <label for="values">设计价值观</label>
        <textarea id="values" v-model="form.values" rows="4" />
      </div>

      <div class="grid-2">
        <div class="field">
          <label for="yearsLabel">年份标识</label>
          <input id="yearsLabel" v-model="form.yearsLabel" />
        </div>
        <div class="field">
          <label for="address">地址</label>
          <input id="address" v-model="form.address" />
        </div>
      </div>

      <div class="field">
        <label for="awards">奖项描述</label>
        <textarea id="awards" v-model="form.awards" rows="6" placeholder="支持换行，小程序将按行展示" />
      </div>

      <div class="grid-2">
        <div class="field">
          <label for="phone">电话</label>
          <input id="phone" v-model="form.phone" />
        </div>
        <div class="field">
          <label for="wechat">微信号</label>
          <input id="wechat" v-model="form.wechat" />
        </div>
      </div>

      <div class="grid-2">
        <DualImageField v-model="form.logo" label="Logo（方）" />
        <DualImageField v-model="form.logoHor" label="Logo（横）" />
      </div>

      <div class="form-actions">
        <button class="btn btn-primary" type="submit" :disabled="saving">
          {{ saving ? '保存中…' : '保存' }}
        </button>
      </div>
    </form>

    <p v-else class="muted">加载中…</p>
  </section>
</template>

<script setup>
import { onMounted, reactive, ref } from 'vue'
import DualImageField from '@/components/DualImageField.vue'
import { getCompany, putCompany, resolveDualImage } from '@/api/client'

const loading = ref(true)
const saving = ref(false)
const loadError = ref('')
const message = ref('')
const messageOk = ref(false)
const aboutViewCount = ref(0)

const form = reactive({
  introTitle: '',
  introBody: '',
  values: '',
  yearsLabel: '',
  address: '',
  awards: '',
  phone: '',
  wechat: '',
  logo: emptyDual(),
  logoHor: emptyDual(),
})

function emptyDual(src) {
  return {
    thumb: src?.thumb || '',
    original: src?.original || '',
    pending: null,
  }
}

function applyCompany(data) {
  form.introTitle = data?.introTitle || ''
  form.introBody = data?.introBody || ''
  form.values = data?.values || ''
  form.yearsLabel = data?.yearsLabel || ''
  form.address = data?.address || ''
  form.awards = data?.awards || ''
  form.phone = data?.phone || ''
  form.wechat = data?.wechat || ''
  form.logo = emptyDual(data?.logo)
  form.logoHor = emptyDual(data?.logoHor)
  if (typeof data?.aboutViewCount === 'number') {
    aboutViewCount.value = data.aboutViewCount
  }
}

onMounted(async () => {
  try {
    const data = await getCompany()
    applyCompany(data || {})
  } catch (err) {
    loadError.value = err.message || '加载失败'
  } finally {
    loading.value = false
  }
})

async function onSubmit() {
  saving.value = true
  message.value = ''
  try {
    const logo = await resolveDualImage(form.logo, 'logo')
    const logoHor = await resolveDualImage(form.logoHor, 'logo')
    const payload = {
      introTitle: form.introTitle,
      introBody: form.introBody,
      values: form.values,
      yearsLabel: form.yearsLabel,
      address: form.address,
      awards: form.awards,
      phone: form.phone,
      wechat: form.wechat,
      logo,
      logoHor,
    }
    const saved = await putCompany(payload)
    applyCompany(saved || payload)
    messageOk.value = true
    message.value = '已保存'
  } catch (err) {
    messageOk.value = false
    message.value = err.message || '保存失败'
  } finally {
    saving.value = false
  }
}
</script>

<style scoped>
.form {
  padding: 1.35rem;
}

.about-views {
  padding: 0.95rem 1.2rem;
  margin-bottom: 1rem;
  font-size: 0.95rem;
}
</style>
