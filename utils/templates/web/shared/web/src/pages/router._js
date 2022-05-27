import {createRouter, createWebHistory} from 'vue-router'

const routes = [

    {
        path: '/',
        component: () => import('./PageWithSidebar.vue'),
        children: [
            {
                path: '/payment',
                redirect: '/payment/waiting',
                component: () => import('./yourpage/ViewTab.vue'),
                children: [
                    {
                        path: '/payment/waiting',
                        component: () => import('../usecase/sample/Page1.vue'),
                    },
                    {
                        path: '/payment/processing',
                        component: () => import('../usecase/sample/Page2.vue'),
                    },
                ],
            },
            {
                path: '/order',
                component: () => import('../usecase/sample/Page3.vue'),
            },
        ],
    },

]

const router = createRouter({
    history: createWebHistory(),
    routes
})

export default router