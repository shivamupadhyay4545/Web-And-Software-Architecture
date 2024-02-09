<template>
    <div class="register">
      <h2>Unfollow User</h2>
      <form @submit.prevent="unfollowUser">
        <input type="text" v-model="unfollowUsername" placeholder="Enter username" />
  
        <button type="submit">Unfollow</button>
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
        unfollowUsername: '',
        errorMessage: '',
        successMessage: '',
      };
    },
    methods: {
      async unfollowUser() {
        try {
          const username = this.$route.params.username;
  
          // Make a POST request to unfollow the user
          const response = await this.$axios.delete(`/user/${username}/follow_list`, {
            data: {
                     following: this.unfollowUsername,
                 },
            });

  
          if (response.status === 200) {
            console.log('User unfollowed successfully!');
            window.location.reload();
            this.successMessage= "User Unfollowed Successfully"
            console.log('User unfollowed successfully!');

            setTimeout(() => {
              this.successMessage = '';
            }, 5000);
            this.unfollowUsername=''
            // You can perform additional actions on successful unfollow
          } else {
            console.error('Failed to unfollow user:', response.statusText);
            this.errorMessage = response.status === 400 ? response.data.message : 'Failed to unfollow user.';
  
            // Clear the error message after 5 seconds
            setTimeout(() => {
              this.errorMessage = '';
            }, 5000);
          }
        } catch (error) {
          console.error('Error during user unfollow:', error.message);
          this.unfollowUsername=''
          this.errorMessage="Oh No!"
  
          if (error.response && error.response.status === 404) {
            // Access the error message from the backend
            this.errorMessage = error.response.data.message || 'User not found in the follow list.';
          } else {
            this.errorMessage = 'An error occurred while unfollowing the user.';
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
}  </style>
  