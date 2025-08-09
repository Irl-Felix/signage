import { createRouter, createWebHistory } from 'vue-router'
import HomeView from '../views/HomeView.vue'
import UsersView from '@/views/UsersView.vue'
import RolesPermissionsView from '@/views/RolesPermissionsView.vue'
import TenantsView from '@/views/TenantsView.vue'
import LogsView from '@/views/LogsView.vue'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'home',
      component: HomeView,
    },
    {
      path: '/users',
      name: 'users',
      component: UsersView,
    },
    {
      path: '/roleperm',
      name: 'roles-permissions',
      component: RolesPermissionsView,
    },
    {
      path: '/tenants',
      name: 'tenants',
      component: TenantsView,
    },
    {
      path: '/logs',
      name: 'logs',
      component: LogsView,
    },
  ],
})

export default router
