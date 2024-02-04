<template>
    <div class="register">
      <h2>Unban User</h2>
      <form @submit.prevent="unbanUser">
        <input type="text" v-model="unbanUsername" placeholder="Enter username" />
  
        <button type="submit">Unban</button>
      </form>
  
      <p v-if="successMessage" style="color: green;">{{ successMessage }}</p>
      <p v-if="errorMessage" style="color: red;">{{ errorMessage }}</p>
    </div>
  </template>
  
  <script>
  // import axios from 'axios';
  
  export default {
    data() {
      return {
        unbanUsername: '',
        successMessage: '',
        errorMessage: '',
      };
    },
    methods: {
      async unbanUser() {
        try {
          const username = this.$route.params.username;
  
          // Make a DELETE request to unban the user
          const response = await this.$axios.delete(`/user/${username}/ban_list`, {
            data: { banned: this.unbanUsername },
          });
  
          if (response.status === 200) {
            console.log('User unbanned successfully!');
            this.successMessage = 'User unbanned successfully!';
            this.unbanUsername=''
            // Clear the success message after 5 seconds
            setTimeout(() => {
              this.successMessage = '';
            }, 5000);
          } else {
            console.error('Failed to unban user:', response.statusText);
            this.errorMessage = 'Failed to unban user.';
  
            // Clear the error message after 5 seconds
            setTimeout(() => {
              this.errorMessage = '';
            }, 5000);
          }
        } catch (error) {
          console.error('Error during user unban:', error.message);
          this.errorMessage = 'An error occurred while unbanning the user.';
  
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
  