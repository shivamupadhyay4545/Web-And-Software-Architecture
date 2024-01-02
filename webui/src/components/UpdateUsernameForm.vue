<template>
    <div class="register">
      <h2>Update Username</h2>
      <form @submit.prevent="updateUsername">
        <input type="text" v-model="name" placeholder="Enter current name" />
  
        <input type="text" v-model="newName" placeholder="Enter new name" />
  
        <button type="submit">Update Username</button>
      </form>
      <p v-if="successMessage" style="color: green;">{{ successMessage }}</p>
      <p v-if="errorMessage" style="color: red;">{{ errorMessage }}</p>
    </div>
  </template>
  
  <script>
  import axios from 'axios';
  
  export default {
    data() {
      return {
        name: '',
        newName: '',
        successMessage: '',
        errorMessage: '',
      };
    },
    methods: {
      async updateUsername() {
        try {
          const username = this.$route.params.username;
          
          const response = await axios.put(`http://localhost:8080/user/${username}`, {
            Name: this.name,
            Newname: this.newName,
          });
  
          if (response.status === 200) {
            console.log('Username updated successfully!');
            this.name = '';
            this.newName = '';
            this.successMessage = "Username Updated Successfully"
            // You can perform additional actions on successful update
          } else {
            console.error('Failed to update username:', response.statusText);
            this.errorMessage="Oh No!"
          }
        } catch (error) {
          console.error('Error during username update:', error.message);
          this.errorMessage =" Oh No!"
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
  