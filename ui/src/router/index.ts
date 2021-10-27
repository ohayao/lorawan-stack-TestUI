import { createRouter, createWebHashHistory, createWebHistory, RouteRecordRaw, RouterOptions } from 'vue-router';
const routes: Array<RouteRecordRaw> = [
  {
    path: '/',
    redirect: '/ui'
  },
  {
    path: '/ui',
    component: () => import('../App.vue'),
    redirect: '/ui/index',
    children: [
      {
        path: 'index',
        component: () => import('../views/index.vue'),
        meta: {
          title: 'Bridge Test'
        }
      },
      {
        path: 'lora',
        component: () => import('../views/lora.vue'),
        meta: {
          title: 'Lora-ttn Test'
        }
      }
    ]
  }
];
const router = createRouter(<RouterOptions>{
  history: createWebHashHistory('#'),
  routes: routes
});
router.beforeEach((to, from, next) => {
  if (to.meta.title) {
    document.title = to.meta.title as string;
  }
  next();
});
export default router;
