import { createRouter, createWebHistory } from 'vue-router'
import HomeView from '../views/HomeView.vue'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'home',
      component: HomeView,
    },
    {
      path: '/login',
      name: 'login',
      component: () => import('../views/Login.vue')
    },
    {
      path: '/admin',
      name: 'admin',
      component: () => import('../views/Admin.vue'),
      children: [
        {
          path: "users",
          component: () => import("../components/AdminUsers.vue")
        },
        {
          path: "roles",
          component: () => import("../components/AdminRoles.vue")
        },
      ]
    },
    {
      path: '/query',
      name: 'query',
      component: HomeView,
      children: [
        {
          path: "list",
          component: () => import("../components/Queries.vue")
        },
        {
          path: "run/:id",
          component: () => import("../components/QueryExecute.vue")
        }
      ]
    }
  ],
})

export default router
