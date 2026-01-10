// Map global CDN variables if they use different names
const ElementUI = window.ELEMENT;
const VueRouter = window.VueRouter;
const VueApexCharts = window.VueApexCharts;
const axios = window.axios;

if (typeof Vue !== 'undefined') {
    if (typeof VueRouter !== 'undefined') Vue.use(VueRouter);
    if (typeof ElementUI !== 'undefined') Vue.use(ElementUI);
    if (typeof VueApexCharts !== 'undefined') Vue.component('apexchart', VueApexCharts);
}

// Components
import Dashboard from './components/salev2/Dashboard.vue';
import ProductVault from './components/salev2/ProductVault.vue';
import OrderManager from './components/salev2/OrderManager.vue';
import SmartUpload from './components/salev2/SmartUpload.vue';
import CustomerManager from './components/salev2/CustomerManager.vue';
import OrderCreate from './components/salev2/OrderCreate.vue';
import ProductCreate from './components/salev2/ProductCreate.vue';
const ProductEdit = require('./components/salev2/ProductEdit.vue').default;
const OrderEdit = require('./components/salev2/OrderEdit.vue').default;
const ExpenseManager = require('./components/salev2/ExpenseManager.vue').default;
const ProductBatchUpload = require('./components/salev2/ProductBatchUpload.vue').default;

const routes = [
    { path: '/salev2', component: Dashboard },
    { path: '/salev2/products', component: ProductVault, name: 'product-vault' },
    { path: '/salev2/orders', component: OrderManager, name: 'order-manager' },
    { path: '/salev2/upload', component: ProductBatchUpload, name: 'smart-upload' },
    { path: '/salev2/customers', component: CustomerManager, name: 'customer-manager' },
    { path: '/salev2/order/create', name: 'order-create', component: OrderCreate },
    { path: '/salev2/order/:id/edit', name: 'order-edit', component: OrderEdit },
    { path: '/salev2/product/create', name: 'product-create', component: ProductCreate },
    { path: '/salev2/product/:id/edit', name: 'product-edit', component: ProductEdit },
    { path: '/salev2/expenses', component: ExpenseManager, name: 'expense-manager' },
    { path: '*', redirect: '/salev2' }
];

const router = new VueRouter({
    mode: 'history',
    routes
});

// CSRF for axios
axios.defaults.headers.common['X-Requested-With'] = 'XMLHttpRequest';
let token = document.head.querySelector('meta[name="csrf-token"]');
if (token) {
    axios.defaults.headers.common['X-CSRF-TOKEN'] = token.content;
}

const app = new Vue({
    el: '#app',
    router
});
