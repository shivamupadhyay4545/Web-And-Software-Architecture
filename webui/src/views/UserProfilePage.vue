<!-- components/Profile.vue -->
<template>
    <div>
      <h2>User Profile</h2>
      
      <!-- Display user details: number of photos, followers, and followings -->
      <div>
        <p>Number of Photos: {{ numberOfPhotos }}</p>
        <p>Followers: {{ followersCount }}</p>
        <p>Followings: {{ followingsCount }}</p>
      </div>
  
      <!-- Display a list of user photos -->
      <div v-if="photos.length > 0">
        <h3>Photos</h3>
        <ul>
          <li v-for="photo in photos" :key="photo.photoId">
            <!-- Display each photo -->
            <img :src="getImageUrl(photo.photoBytes)" alt="User Photo" />
          </li>
        </ul>
      </div>
      <div v-else>
        <p>No photos available</p>
      </div>
    </div>
  </template>
  
  <script>
  import axios from 'axios';
  export default {
    data() {
      return {
        numberOfPhotos: 0,
        followersCount: 0,
        followingsCount: 0,
        photos: [],
      };
    },
    methods: {
      async fetchUserProfile() {
        try {
            const username = this.$route.params.username;
            const response = await axios.get(`http://localhost:8080/user/${username}/profile`);
          const data = await response.json();
  
          this.numberOfPhotos = data['my profile']['PhotoNo'];
          this.followersCount = data['my profile']['Followers'];
          this.followingsCount = data['my profile']['Following'];
          this.photos = data['my profile']['Photos'];
        } catch (error) {
          console.error('Error fetching user profile:', error);
        }
      },
      getImageUrl(photobytes) {
      if (photobytes) {
        return `data:image/jpeg;base64,${photobytes}`;
      }
  
        // If photobytes is null or undefined, return a placeholder or empty string
        return '';
      },
    },
    created() {
      this.fetchUserProfile();
    },
  };
  </script>
  
  <style scoped>
  /* Add your component styles here */
  </style>
  