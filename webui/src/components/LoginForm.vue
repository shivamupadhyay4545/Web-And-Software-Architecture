<template>
  <form @submit.prevent="login">
    <img class="logo" src="../assets/WasaPhotoLogo.png" />
    <h1>Login</h1>
    <div class="register">
      <!-- <input type="text" v-model="id" placeholder="Enter USERID" /> -->
      <input type="text" v-model="name" placeholder="Enter USERNAME" />
      <button type="submit">Login</button>
    </div>
  </form>
</template>

<script>
// import axios from 'axios';
import { setAuthToken } from '../services/axios';

export default {
  name: 'LoginForm',
  data() {
    return {
      token: '',
      name: '',
    };
  },
  methods: {
    async login() {
      try {
        const response = await this.$axios.post('/session', {
          name: this.name,
        });

        if (response.status === 200) {
          this.token = response.data["authtoken"];

          // Set the Authorization header for subsequent requests
          this.$axios.defaults.headers.common['Authorization'] = this.token;

          setAuthToken(this.token);

          console.log('Login successful!');
          const username = this.name;

          // Use Vue Router to navigate to the user profile page
          this.$router.push({ path: `/${username}/home` });
          // You can redirect the user to another page or perform other actions on successful login
        } else {
          console.error('Login failed:', response.statusText);
        }
      } catch (error) {
        console.error('Error during login:', error.message, this.name);
      }
    },
  },
};
</script>

<style>
.logo {
  width: 150px;
  margin-left: auto;
  margin-right: auto;
  display: block;
}
.register input {
  width: 300px;
  height: 40px;
  padding-left: 20px;
  display: block;
  margin-bottom: 30px;
  margin-right: auto;
  margin-left: auto;
  border: 1px cornflowerblue;
}
.register button {
  width: 320px;
  height: 40px;
  border: 1px crimson;
  background: skyblue;
  color: aliceblue;
  cursor: pointer;
}
</style>



