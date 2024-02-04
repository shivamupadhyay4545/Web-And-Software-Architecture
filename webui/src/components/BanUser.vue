<template>
    <div class="register">
      <h2>Ban User</h2>
      <form @submit.prevent="banUser">
        <input type="text" v-model="banUsername" placeholder="Enter username" />
  
        <button type="submit">Ban</button>
      </form>
  
      <p v-if="errorMessage" style="color: red;">{{ errorMessage }}</p>
      <p v-if="successMessage" style="color: green;">{{ successMessage }}</p>
    </div>
  </template>
  
  <script>
  // import axios from 'axios';
  
  export default {
    data() {
      return {
        banUsername: '',
        errorMessage: '',
        successMessage: '',
      };
    },
    methods: {
      async banUser() {
        try {
          const currentUsername = this.$route.params.username;
  
          // Make a POST request to ban the user
          const response = await this.$axios.post(`/user/${currentUsername}/ban_list`, {
            banned: this.banUsername,
          });
  
          if (response.status === 200) {
            console.log('User banned successfully!');
            this.banUsername=''
            window.location.reload();
            
            // You can perform additional actions on successful ban
          } else {
            console.error('Failed to ban user:', response.statusText);
            this.errorMessage = response.status === 409 ? response.data.message : 'Failed to ban user.';
  
            // Clear the error message after 5 seconds
            setTimeout(() => {
              this.errorMessage = '';
            }, 5000);
          }
        } catch (error) {
          console.error('Error during user ban:', error.message);
  
          if (error.response && error.response.status === 409) {
            // Access the error message from the backend
            this.errorMessage = error.response.data.message || 'User does not exist.';
          } else {
            this.errorMessage = 'An error occurred while banning the user.';
          }
  
          // Clear the error message after 5 seconds
          setTimeout(() => {
            this.errorMessage = '';
          }, 5000);
        }
      },
    },
  };
  </script>
  
  <style>
.register input {
  width: 200px;
  height: 40px;
  padding-left: 20px;
  display: block;
  margin-bottom: 30px;
  margin-right: auto;
  margin-left: auto;
  border: 1px cornflowerblue;
}
.register button {
  width: 300px;
  height: 40px;
  border: 1px crimson;
  background: skyblue;
  color: aliceblue;
  cursor: pointer;
}
  </style>
  