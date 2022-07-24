import {createRouter, createWebHistory} from "vue-router";


import Main from '/src/components/Main.vue'
import Home from '/src/components/pages/home/Home.vue'
import Login from '/src/components/Login.vue'
import Signup from '/src/components/Signup.vue'

import NProgress from 'nprogress'

const routes = [
    {
        path: '/',
        name: 'Main',
        component: Main,
        children: [
            {
                path: '',
                name: 'Лента',
                component: Home

            },
            {
                path: '/chat',
                name: 'Чаты',
                component: () => import('/src/components/pages/chat/Chatest.vue')

            },
            {
                path: '/profile',
                name: 'Профиль',
                component: () => import('/src/components/pages/profile/Profile.vue')
            },
            {
                path: '/settings',
                name: 'Настройки',
                component: () => import('/src/components/pages/settings/Settings.vue')
            },
            {
                path: '/user/:username',
                name: 'Пользователь',
                component: () => import('/src/components/pages/userprofile/UserProfile.vue')
            }
        ]
    },
    {
        path: '/login',
        name: 'Login',
        component: Login,
    },
    {
        path: '/signup',
        name: 'Signup',
        component: Signup,
    }
]


const router = createRouter({
    history: createWebHistory(),
    routes,
})

router.beforeResolve((to, from, next) => {
    // If this isn't an initial page load.
    if (to.name) {
        NProgress.start()
    }
    next()
})


router.afterEach((to, from) => {
    // Complete the animation of the route progress bar.
    NProgress.done()
})


//page setup (Перед каждой загрузкой страницы)
router.beforeEach((to, from) => {
    //check which page
    switch (to.name) {
        case 'Login':
            document.title = 'S&D - Вход';
            break;
        case 'Signup':
            document.title = 'S&D - Регистрация';
            break;
        case 'Home':
            document.title = 'S&D - Home';
            break;
        case 'Profile':
            document.title = 'S&D - Profile';
            break;
        case 'Settings':
            document.title = 'S&D - Настройки';
            break;
        case 'User Profile':
            document.title = 'S&D - Профиль пользователя';
            break;
        default:
            document.title = 'S&D';
            break;
    }
})

export default router