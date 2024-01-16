import { createRouter, createWebHistory } from 'vue-router'
import IndexView from '../views/IndexView.vue'
import TorrentView from '../views/TorrentView.vue'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'index',
      component: IndexView
    },
    {
      path: '/',
      name: 'torrents',
      component: TorrentView
    }
  ]
})

export default router
