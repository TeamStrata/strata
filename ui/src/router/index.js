import { createRouter, createWebHistory } from 'vue-router'
import HomeView from '../views/HomeView.vue'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
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
              component: () => import('../components/Queries.vue')
            },
            {
              path: 'run/:id',
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
              component: () => import('../components/AdminUsers.vue')
            },
            {
              path: 'roles',
              component: () => import('../components/AdminRoles.vue')
            }
          ]
        },
      ]
    }
  ],
})

export default router
