<template>
    <div class="register">
      <h2>Follow User</h2>
      <form @submit.prevent="followUser">
        <input type="text" v-model="followUsername" placeholder="Enter username" />
  
        <button type="submit">Follow</button>
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
        followUsername: '',
        errorMessage: '',
        successMessage: '',
        token: '',
      };
    },
    methods: {
      async followUser() {
        try {
          // Retrieve the authorization token from Axios defaults
          const authToken = this.$axios.defaults.headers.common['Authorization'];
          console.log('Authorization Token:', authToken);

          const username = this.$route.params.username;
  
          // Make a POST request to follow the user
          const response = await this.$axios.post(`/user/${username}/follow_list`, {
      following: this.followUsername,
    }, {
      headers: {
        'Authorization': authToken,
      },
    });
  
          if (response.status === 200) {
            console.log('User followed successfully!');
            this.followUsername =''
            this.successMessage ="User Followed Successfully"
            window.location.reload();
            console.log('User followed successfully!');
            // You can perform additional actions on successful follow
          } else {
            this.followUsername= ''
            console.error('Failed to follow user:', response.message);
            this.errorMessage = response.status === 409 ? response.message : 'Failed to follow user.';
            setTimeout(() => {
          this.errorMessage = '';
        }, 5000);
          }
        } catch (error) {
            this.followUsername= ''
          console.error('Error during user follow:', error.message);
  
          if (error.response && error.response.status === 409) {
            // Access the error message from the backend
            this.errorMessage = error.response.data.message || 'User does not exist.';
          } else {
            this.errorMessage = 'An error occurred while following the user.';
          }
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
  