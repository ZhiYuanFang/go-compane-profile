import { createRouter, createWebHistory } from 'vue-router'
import AdminLayout from '@/layouts/AdminLayout.vue'
import LoginView from '@/views/LoginView.vue'
import CompanyView from '@/views/CompanyView.vue'
import PortfolioListView from '@/views/PortfolioListView.vue'
import PortfolioEditView from '@/views/PortfolioEditView.vue'
import ActivityListView from '@/views/ActivityListView.vue'
import ActivityEditView from '@/views/ActivityEditView.vue'
import { getCompany } from '@/api/client'
import { isValidCategory } from '@/constants/portfolioCategories'

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
      { path: 'portfolios', redirect: '/portfolios/residential' },
      {
        path: 'portfolios/:category',
        name: 'portfolios',
        component: PortfolioListView,
        beforeEnter: (to) => {
          if (!isValidCategory(to.params.category)) {
            return { path: '/portfolios/residential' }
          }
          return true
        },
      },
      {
        path: 'portfolios/:category/new',
        name: 'portfolio-new',
        component: PortfolioEditView,
        beforeEnter: (to) => {
          if (!isValidCategory(to.params.category)) {
            return { path: '/portfolios/residential' }
          }
          return true
        },
      },
      {
        path: 'portfolios/:category/:id',
        name: 'portfolio-edit',
        component: PortfolioEditView,
        beforeEnter: (to) => {
          if (!isValidCategory(to.params.category)) {
            return { path: '/portfolios/residential' }
          }
          return true
        },
      },
      { path: 'activity', redirect: '/activities' },
      { path: 'activities', name: 'activities', component: ActivityListView },
      { path: 'activities/new', name: 'activity-new', component: ActivityEditView },
      { path: 'activities/:id', name: 'activity-edit', component: ActivityEditView },
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
