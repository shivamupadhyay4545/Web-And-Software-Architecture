<!-- App.vue -->

<template>
  <div id="app">
    <router-view @login-successful="handleLoginSuccess"/>
    
  </div>
</template>

<script>
//import LoginForm from './components/LoginForm.vue';
//import Home from './views/HomeView.vue';
//import axios from 'axios';

export default {
  
  methods: {
    handleLoginSuccess(username) {
      // Use Vue Router to navigate to the user profile page
      this.$router.push({ path: `/${username}/home` });
    },
  },
};
</script>

<style>
/* Add global styles for your application */
body {
  font-family: 'Arial', sans-serif;
  padding: 0;
  display: flex;
  justify-content: center;
  height: 100vh;
  background-color: #f0f0f0;
  margin-top: 20px;
}
</style>
