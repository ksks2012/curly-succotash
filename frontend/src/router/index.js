import { createRouter, createWebHistory } from 'vue-router'
import GameForm from '../components/GameForm.vue'
import About from '../components/About.vue'

const routes = [
  {
    path: '/',
    name: 'Home',
    component: GameForm,
  },
  {
    path: '/about',
    name: 'About',
    component: About,
  },
]

const router = createRouter({
  history: createWebHistory(),
  routes,
})

export default router