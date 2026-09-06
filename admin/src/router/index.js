import { createRouter, createWebHistory } from 'vue-router'
import AdminLayout from '@/layouts/AdminLayout.vue'
import LoginView from '@/views/LoginView.vue'
import CompanyView from '@/views/CompanyView.vue'
import PortfolioListView from '@/views/PortfolioListView.vue'
import PortfolioEditView from '@/views/PortfolioEditView.vue'
import PricingView from '@/views/PricingView.vue'
import { getCompany } from '@/api/client'

const routes = [
  {
    path: '/login',
    name: 'login',
    component: LoginView,
    meta: { public: true },
  },
  {
    path: '/',
    component: AdminLayout,
    children: [
      { path: '', redirect: '/company' },
      { path: 'company', name: 'company', component: CompanyView },
      { path: 'portfolios', name: 'portfolios', component: PortfolioListView },
      { path: 'portfolios/new', name: 'portfolio-new', component: PortfolioEditView },
      { path: 'portfolios/:id', name: 'portfolio-edit', component: PortfolioEditView },
      { path: 'pricing', name: 'pricing', component: PricingView },
    ],
  },
]

const router = createRouter({
  history: createWebHistory('/admin/'),
  routes,
})

let authChecked = false
let authenticated = false

export function setAuthState(value) {
  authenticated = value
  authChecked = true
}

router.beforeEach(async (to) => {
  if (to.meta.public) return true

  if (!authChecked) {
    try {
      await getCompany()
      authenticated = true
    } catch {
      authenticated = false
    }
    authChecked = true
  }

  if (!authenticated) {
    return { name: 'login', query: { redirect: to.fullPath } }
  }
  return true
})

export default router
