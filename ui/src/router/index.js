import { createRouter, createWebHistory } from 'vue-router'
import HomeView from '../views/HomeView.vue'

import { useUserStore } from '@/stores/user'





const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/login',
      name: 'login',
      component: () => import('../views/Login.vue')
    },
    {
      path: '',
      component: HomeView,
      children: [
        {
          path: 'query',
          name: 'query',
          children: [
            {
              path: 'list',
              meta: {
                title: 'Queries',
                description: 'Build and manage SQL queries.',
                configurable: false,
              },
              component: () => import('../components/Queries.vue')
            },
            {
              path: 'run/:id',
              meta: {
                title: 'Queries',
                description: 'Build and manage SQL queries.',
                configurable: false,
              },
              component: () => import('../components/QueryExecute.vue')
            }
          ]
        },
        {
          path: 'admin',
          name: 'admin',
          children: [
            {
              path: 'users',
              meta: {
                title: 'Members',
                description: 'Manage your organizations members.',
                configurable: false,
              },
              component: () => import('../components/AdminUsers.vue')
            },
            {
              path: 'roles',
              meta: {
                title: 'Roles',
                description: 'Manage your organizations roles.',
                configurable: false,
              },
              component: () => import('../components/AdminRoles.vue')
            }
          ]
        },
      ]
    }
  ],
})

router.beforeEach(async (to, from) => {
  const store = useUserStore()
  if (!store.isLoggedIn && to.name !== 'login') {
    // redirect the user to the login page
    return { name: 'login' }
  }
})


export default router
