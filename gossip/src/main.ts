import { PiniaColada } from "@pinia/colada";
import { createPinia } from "pinia";
import { createApp } from "vue";
import { createRouter, createWebHistory } from "vue-router";
import App from "./App.vue";
import "./style.css";

import LoginView from "@/views/auth/Login.vue";
import VerifyLinkView from "@/views/auth/VerifyLink.vue";

const routes = [
  {
    path: "/login",
    component: LoginView,
  },
  {
    path: "/verify-link",
    component: VerifyLinkView,
  },
];

const router = createRouter({
  history: createWebHistory(),
  routes,
});

const app = createApp(App);
app.use(createPinia());
app.use(PiniaColada, {
  queryOptions: {
    staleTime: 0,
  },
  mutationOptions: {},
  plugins: [],
});

app.use(router);
app.mount("#app");
