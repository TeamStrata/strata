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
          path: 'queries',
          name: 'queries',
          component: () => import('../components/Queries.vue'),
          meta: {
            title: 'Queries',
            description: 'Build and manage your organizations queries.',
            configurable: false,
          },
        },
        {
          path: 'admin',
          name: 'admin',
          children: [
            {
              path: 'users',
              meta: {
                title: 'Members',
                description: 'Manage your organization\'s members.',
                configurable: false,
              },
              component: () => import('../components/AdminUsers.vue')
            },
            {
              path: 'roles',
              meta: {
                title: 'Roles',
                description: 'Manage your organization\'s roles.',
                configurable: false,
              },
              component: () => import('../components/AdminRoles.vue')
            },
            {
              path: 'settings',
              meta: {
                title: 'Settings',
                description: 'Manage your platform settings.',
                configurable: false,
              },
              component: () => import('../components/AdminSettings.vue')
            }
          ]
        },
        {
          path: '/dashboard/:id',
          name: 'dashboard',
          meta: {

          },
          component: () => import('../components/Dashboard.vue')
        },
        {
          path: '/charts',
          name: 'charts',
          children: [
            {
              path: "list",
              meta: {
                title: 'Charts',
                description: 'Manage your organization\'s charts.',
                configurable: false,
              },
              component: () => import("../components/Charts.vue")
            }
          ]
        }
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
