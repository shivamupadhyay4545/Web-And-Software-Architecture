// router.js


import { createRouter, createWebHistory } from 'vue-router';
import Login from './components/LoginForm.vue';
import HomeView from './views/HomeView.vue';
import UserProfile from './components/UserProfile.vue';
import CommentPhotoForm from './components/CommentPhotoForm.vue'



const router = createRouter({
  history: createWebHistory(),
  routes: [
    { path: '/', component: Login },
    { path: '/home', component: HomeView },
    { path: '/:username/home', name : 'HomeView' ,component: HomeView }, // Dynamic route for username
    { path: '/:username/profile', name: 'UserProfile'  , component: UserProfile },
    { path: '/:username/profile/:PhotoId', name: 'CommentPhotoForm'  , component: CommentPhotoForm, props: true, },
  ],
});

export default router;

