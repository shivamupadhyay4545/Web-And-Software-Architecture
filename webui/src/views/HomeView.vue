<template>
  <div>
    <UserHeader />
    <h1 class="main-title">Hello {{ this.$route.params.username }}, Welcome to your feed</h1>
    <div>
    <router-link :to="{ name: 'UserProfile', params: { username: $route.params.username } }">
      Go to User Profile
    </router-link>
  </div>
  <div>
    <router-link :to="{ path: '/', component: Login }">
      LOGOUT
    </router-link>
  </div>
    <div v-if="userData && userData.photos" class="photo-container">
      <!-- Display user data when 'photos' is not null -->
      <div class="photo-card" v-for="photo in userData.photos" :key="photo.PhotoId">
        <div class="user-info">
          <router-link :to="{ name: 'OtherUserProfile', params: { username:  $route.params.username, otherusername: photo.Username } }">
            <p class="username">{{ photo.Username }}</p>
          </router-link>
          
        </div>
        <button class="photo-button" @click="showComments(photo)">
          <img :src="getImageUrl(photo.Photobytes)" alt="User Photo" />
        </button>
        <div class="like-buttons">
          <button style="background-color:red;" v-if="photo.Liked === 1" @click="toggleLike(photo.PhotoId, true)">Dislike</button>
          <button v-else @click="toggleLike(photo.PhotoId, false)">Like</button>
          <p class="likes-count">Likes: {{ photo.Likes }}</p>
        </div>
        <div class="comments-section">
          <img class="comment-logo" src="../assets/comment-logo.png" />
          <p class="comments-count">Comments: {{ photo.NoComments }}</p>
          <input v-model="commentInputs[photo.PhotoId]" class="comment-input" placeholder="Add a comment..." />
          <button :disabled="isCommentInputEmpty(photo.PhotoId)" @click="postComment(photo.PhotoId)" class="post-comment-button">Post Comment</button>
        </div>
        <p class="upload-time">Uploaded At {{ formatTimestamp(photo.CreatedAt) }}</p>
      </div>
    </div>
    <div v-else>
      <p class="welcome-message">Welcome Our new user! Follow someone to see their photos.</p>
    </div>
  </div>
</template>

<script>
// import axios from 'axios';
import UserHeader from '../components/UserHeader.vue';
import moment from 'moment';
// import UploadPhotoForm from '../components/UploadPhotoForm.vue';
// import FollowUser from '../components/FollowUser.vue';
// import UnFollowUser from '../components/UnFollowUser.vue';
// import BanUser from '../components/BanUser.vue';
// import UnBanUser from '../components/UnBanUser.vue';
// import UpdateUsernameForm from '../components/UpdateUsernameForm.vue';

export default {
  components: {
    UserHeader
    // UploadPhotoForm,
    // FollowUser,
    // UnFollowUser,
    // BanUser,
    // UnBanUser,
    // UpdateUsernameForm
},
  name: 'HomeView',
  data() {
    return {
      userData: null,
      commentInputs: {},
      userName: ''
    };
  },
  mounted() {
    // Fetch user data after component is mounted
    this.fetchUserData();
  },
  methods: {
    async toggleLike(photoId, liked) {
      try {
        const username = this.$route.params.username;
        const authToken = this.$axios.defaults.headers.common['Authorization'];

        const config = {
          headers: {
            'Authorization': authToken,
          },
        };

        if (liked) {
          // Send DELETE request to unlike the photo
          await this.$axios.delete(`/user/${username}/photos/likes?Photoid=${photoId}`, config);
        } else {
          // Send POST request to like the photo
          await this.$axios.post(`/user/${username}/photos/likes?Photoid=${photoId}`, config);
        }

        window.location.reload()
        // Update dislikeStatus after toggling the like state
        
      } catch (error) {
        console.error('Error while toggling like:', error.message,photoId);
      }
    },
    async postComment(PhotoId) {
      try {
        const authToken = this.$axios.defaults.headers.common['Authorization'];

        const username = this.$route.params.username;
        const response = await this.$axios.post(
          `/user/${username}/photos/comment?Photoid=${PhotoId}`,
          {
            content: this.commentInputs[PhotoId],
          },
          {
      headers: {
        'Authorization': authToken,
      },
    }
        );

        if (response.status === 200) {
          // Successfully posted comment, fetch updated user data
          this.fetchUserData();
          // Clear the comment input
          this.commentInputs[PhotoId] = '';
        } else {
          console.error('Failed to post comment:', response.statusText);
        }
      } catch (error) {
        const username = this.$route.params.username;
        console.error('Error while posting comment:', error.message, 'username:', username, 'Photoid:', PhotoId);
      }
    },
    isCommentInputEmpty(photoId) {
      return !this.commentInputs[photoId] || !this.commentInputs[photoId].trim();
    },

    async fetchUserData() {
      try {
        const username = this.$route.params.username;
        const authToken = this.$axios.defaults.headers.common['Authorization'];
        const config = {
          headers: {
            'Authorization': authToken,
          },
        };
        const response = await this.$axios.get(`/user/${username}`, config);

        if (response.status === 200) {
          this.userData = response.data;
          
        } else {
          console.error('Failed to fetch user data:', response.statusText);
        }
      } catch (error) {
        console.error('Error during user data fetch:', error.message, this.$axios.headers.Authorization);
      }
    },
    formatTimestamp(timestamp) {
      // Use moment to format the timestamp
      return moment(timestamp).format('YYYY-MM-DD HH:mm:ss');
    },
    showComments(photo) {
      const username = this.$route.params.username;
      // Use Vue Router to navigate to the comments route
      this.$router.push({
        name: 'CommentPhotoForm', // Replace with the actual name of your Comments route
        params: { PhotoId: photo.PhotoId, username: username },
      });
    },
    getImageUrl(photobytes) {
      if (photobytes) {
        return `data:image/jpeg;base64,${photobytes}`;
      }

      // If photobytes is null or undefined, return a placeholder or empty string
      return '';
    },
  },
};
</script>


<style scoped>
.main-title {
  font-size: 24px;
  margin-bottom: 20px;
}

.photo-container {
  display: flex;
  flex-wrap: wrap;
  justify-content: space-around;
}

.photo-card {
  width: 300px;
  margin: 20px;
  padding: 15px;
  border: 1px solid #ddd;
  border-radius: 5px;
  box-shadow: 0 0 10px rgba(0, 0, 0, 0.1);
}

.user-info {
  font-weight: bold;
}

.photo-button {
  border: none;
  background: none;
  cursor: pointer;
}

.like-buttons {
  display: flex;
  align-items: center;
  margin-top: 10px;
}

.like-buttons button {
  margin-right: 10px;
  background-color: #3498db;
  color: #fff;
  padding: 5px 10px;
  border: none;
  border-radius: 3px;
  cursor: pointer;
}

.like-buttons button:hover {
  background-color: #2980b9;
}

.likes-count {
  margin-left: 10px;
}

.comments-section {
  display: flex;
  align-items: center;
  margin-top: 10px;
}

.comment-logo {
  width: 20px;
  margin-right: 5px;
}

.comments-count {
  margin-left: 5px;
}

.comment-input {
  margin-right: 10px;
  padding: 5px;
  width: 100%;
}

.post-comment-button {
  padding: 5px 10px;
  background-color: #3498db;
  color: #fff;
  cursor: pointer;
  border: none;
  border-radius: 3px;
}

.post-comment-button:hover {
  background-color: #2980b9;
}

.upload-time {
  margin-top: 10px;
  font-size: 12px;
  color: #888;
}

.welcome-message {
  font-size: 18px;
  color: #888;
}


</style>







˝
