<template>
    <div>
      <h2>Upload Photo</h2>
      <form @submit.prevent="uploadPhoto">
        <label for="image">Select Photo:</label>
        <input type="file" ref="fileInput" @change="handleFileChange" />
  
        <button type="submit" :disabled="!photoSelected">Upload Photo</button>
      </form>
    </div>
  </template>
  
  <script>
  import axios from 'axios';
  
  export default {
    data() {
      return {
        photoSelected: false,
        photoFile: null,
      };
    },
    methods: {
      handleFileChange(event) {
        const fileInput = event.target;
        const file = fileInput.files[0];  
        if (file) {
          this.photoFile = file;
          this.photoSelected = true;
        }
      },
  
      async uploadPhoto() {
        try {
          const username = this.$route.params.username;
  
          if (!this.photoFile) {
            console.error('No photo selected');
            return;
          }
  
          const formData = new FormData();
          formData.append('image', this.photoFile);
  
          const response = await axios.post(`http://localhost:8080/user/${username}`, formData, {
            headers: {
              'Content-Type': 'multipart/form-data',
            },
          });
  
          if (response.status === 201) {
            console.log('Photo uploaded successfully!');
            this.$refs.fileInput.value = '';
            this.photoSelected = false;
            window.location.reload();
          } else {
            console.error('Failed to upload photo:', response.statusText);
          }
        } catch (error) {
          console.error('Error during photo upload:', error.message);
        }
      },
    },
  };
  </script>
  
  <style>
  /* Add styles as needed */
  </style>
  