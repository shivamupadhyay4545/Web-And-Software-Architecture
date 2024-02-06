<template>
    <div class="register">
      <h2>Update Username</h2>
      <form @submit.prevent="updateUsername">

  
        <input type="text" v-model="newName" placeholder="Enter new username" />
  
        <button type="submit">Update Username</button>
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
        newName: '',
        successMessage: '',
        errorMessage: '',
      };
    },
    methods: {
      async updateUsername() {
        try {
          const username = this.$route.params.username;
          const response = await this.$axios.put(`/user/${username}/`, {
            Newname: this.newName,
          });
  
          if (response.status === 200) {
            console.log('Username updated successfully!');
            
            this.successMessage = "Username Updated Successfully"
            setTimeout(() => {
              this.successMessage = '';
            }, 5000);
            this.$router.push({ path: `/${this.newName}/home` });
            this.newName = '';
            // You can perform additional actions on successful update
          } else {
            console.error('Failed to update username:', response.statusText);
            this.errorMessage="Oh No!"
            setTimeout(() => {
              this.errorMessage = '';
            }, 5000);
          }
        } catch (error) {
          console.error('Error during username update:',this.newName);
          this.errorMessage =" Oh No!"
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
  