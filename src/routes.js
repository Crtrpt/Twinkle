import Index from './page/Index.vue'
import About from './page/About.vue'
import Login from "./page/Login.vue"

const routes = [
    { path: '/', component: Index },
    { path: '/login', component: Login,name:"login" },
    { path: '/about', component: About },
]
export default routes
