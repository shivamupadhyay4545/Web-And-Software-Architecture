<template>
    <div>
      <div>
      <router-link
        :to="{ name: 'HomeView', params: { username: $route.params.username } }"
      >
        Go to Home
      </router-link>
    </div>
    <div>
      <router-link :to="{ name: 'UserProfile', params: { username: $route.params.username } }">
      Go to User Profile
    </router-link>
  </div>
      <img :src="getImageUrl(photoBytes)" alt="User Photo" />
      <button v-if="isOwner" @click="deletePhoto">Delete Photo</button>
      <div v-for="comment in comments" :key="comment.DateTime">
        <p>{{ comment.CommentUser }}:</p>
        <p>{{ comment.Comment }}</p>
        <p>{{ formatTimestamp(comment.DateTime) }}</p>
        <button v-if="this.username === comment.CommentUser" @click="deleteComment(comment.Comment)">
        Delete Comment
      </button>
        <hr />
      </div>
    </div>
  </template>
  
  <script>
  // import axios from 'axios';
  import moment from 'moment';
  
  export default {
    data() {
      return {
        comments: [],
        photoBytes: '',
      };
    },
    props: {
      username: String,
      PhotoId: String,
    },
    computed: {
    isOwner() {
      // Check if the current user is the owner of the photo
      return this.username === this.PhotoId.split('_')[0];
    },
  },
    methods: {
        async deleteComment(Content) {
            const jsonPayload = { Content };
      try {
        const jsonPayload = { Content };

        const queryParamValue = this.PhotoId; // Replace with your actual query parameter value
        const url = `/user/${this.username}/photos/comment?Photoid=${queryParamValue}`;

        const response = await this.$axios.delete(url, { data: jsonPayload });


        if (response.status === 200) {
          console.log('Comment deleted successfully!');
          // Fetch comments data after deleting the comment
          this.fetchComments();
        } else {
          console.error('Failed to delete comment:', response.statusText);
        }
      } catch (error) {
        console.log(jsonPayload)
        console.error('Error deleting comment11:', error.message);
      }
    },

    async deletePhoto() {
      try {
        // Assuming this.username, this.PhotoId, and other necessary values are defined
        const url = `/user/${this.username}/deleted_photos?Photoid=${this.PhotoId}`;
        
        const response = await this.$axios.delete(url);

        if (response.status === 200) {
            window.history.back();
          console.log('Photo deleted successfully!');
          // Add additional logic if needed
        } else {
          console.error('Failed to delete photo:', response.statusText);
        }
      } catch (error) {
        console.error('Error deleting photo:', error.message);
      }
    },
  


      async fetchComments() {
        const parts = this.PhotoId.split('_');
        try {
            const parts = this.PhotoId.split('_');
          // Make a GET request to fetch comments data
          const response = await this.$axios.get(`/user/${parts[0]}/photos/${this.PhotoId}`);
  
          // Update data with the fetched comments and photoBytes
          this.comments = response.data.comments;
          this.photoBytes = response.data.photobytes;
        } catch (error) {
            console.log("bhenchod",parts[0],this.PhotoId)
          console.error('Error fetching comments:', this.PhotoId,this.username);
        }
      },
      getImageUrl(photobytes) {
        if (photobytes) {
          return `data:image/jpeg;base64,${photobytes}`;
        }
        return ''; // Return placeholder or empty string if photobytes is null or undefined
      },
      formatTimestamp(timestamp) {
        // Use moment to format the timestamp
        return moment(timestamp).format('YYYY-MM-DD HH:mm:ss');
      },
    },
    mounted() {
      // Fetch comments data when the component is mounted
      this.fetchComments();
    },
  };
  </script>
  
  <style scoped>
  /* Add your component styles here */
  </style>
  