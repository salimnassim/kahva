import { createRouter, createWebHistory } from 'vue-router'
import TorrentView from '../views/TorrentView.vue'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'torrents',
      component: TorrentView
    }
  ]
})

export default router
