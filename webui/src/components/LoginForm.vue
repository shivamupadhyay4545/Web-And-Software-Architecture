<template>
  <form @submit.prevent="login">
    <img class="logo" src="../assets/WasaPhotoLogo.png" />
    <h1>Login</h1>
    <div class="register">
      <input type="text" v-model="id" placeholder="Enter ID" />
      <input type="text" v-model="name" placeholder="Enter Name" />
      <button type="submit">Login</button>
    </div>
  </form>
</template>

<script>
import axios from 'axios';

export default {
  name: 'LoginForm',
  data() {
    return {
      id: '',
      name: '',
    };
  },
  methods: {
    async login() {
      try {
        const response = await axios.post('http://localhost:8080/session', {
          id: this.id,
          name: this.name,
        });

        if (response.status === 200) {
          console.log('Login successful!');
          const username = this.id;

          // Use Vue Router to navigate to the user profile page
          this.$router.push({ path: `/${username}/home` });
          // You can redirect the user to another page or perform other actions on successful login
        } else {
          console.error('Login failed:', response.statusText);
        }
      } catch (error) {
        console.error('Error during login:', error.message);
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






























<!-- components/LoginForm.vue

<template>
  <div class ="register">
    <h2>Login</h2>
    <form @submit.prevent="login">
      
      <label for="id">ID:</label>
      <input type="text" id="id" v-model="id" />

      <label for="name">Name:</label>
      <input type="text" id="name" v-model="name" />

      <buttontype="submit">Login</button>
    </form>
  </div>
</template>

<script>
export default {
  data() {
    return {
      id: '',
      name: '',
    };
  },
  methods: {
    async login() {
      try {
        const response = await fetch('http://localhost:8080/session', {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
          },
          body: JSON.stringify({
            id: this.id,
            name: this.name,
          }),
        });

        if (response.ok) {
          console.log('Login successful!');
          const username = this.id; // Use the id or any other property as the username

          // Use Vue Router to navigate to the user profile page
          this.$router.push({ path: `http://localhost:8080/user/${username}` });
          // You can redirect the user to another page or perform other actions on successful login
        } else {
          console.error('Login failed:', response.statusText);
        }
      } catch (error) {
        console.error('Error during login:', error.message);
      }
    },
  },
};
</script>

<style >
.click button{
  width: 500px;
}
</style> -->
