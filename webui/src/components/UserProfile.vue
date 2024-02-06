<template>
    <div>
      <router-link
        :to="{ name: 'HomeView', params: { username: $route.params.username } }"
      >
        Go to Home
      </router-link>
      <div class="profile-container">
        <div class="user-details">
          <p class="detail">Photos: {{ numberOfPhotos }}</p>
          <p class="detail">Followers: {{ followersCount }}</p>
          <p class="detail">Followings: {{ followingsCount }}</p>
        </div>
  
        <div v-if="photos" class="photo-list">
          <h3 class="photo-heading">Photos</h3>
          <ul>
            <div v-for="photo in photos" :key="photo.photoId" class="photo-item">
              <button @click="showComments(photo)">
                <img :src="getImageUrl(photo.Photobytes)" alt="User Photo" />
              </button>
              <div class="photo-details">
                <p class="comment-count">Comments: {{ photo.NoComments }}</p>
                <input
                  v-model="commentInput[photo.photoId]"
                  class="comment-input"
                  placeholder="Add a comment..."
                />
                <button
                  @click="postComment(photo.PhotoId)"
                  class="post-comment-button"
                >
                  Post Comment
                </button>
                <p class="like-count">Likes: {{ photo.Likes }}</p>
                <button style="background-color: red;"
                  v-if="photo.Liked === 1"
                  @click="toggleLike(photo.PhotoId, true)"
                  class="dislike-button"
                >
                  Dislike
                </button>
                <button
                  v-else
                  @click="toggleLike(photo.PhotoId, false)"
                  class="like-button"
                >
                  Like
                </button>
                <p class="upload-time">
                  Uploaded At {{ formatTimestamp(photo.CreatedAt) }}
                </p>
              </div>
            </div>
          </ul>
        </div>
        <div v-else>
          <p>No photos available</p>
        </div>
      </div>
    </div>
  </template>
  
  <script>
  // import axios from 'axios';
  import moment from 'moment';
  
  export default {
    data() {
      return {
        numberOfPhotos: 0,
        followersCount: 0,
        followingsCount: 0,
        photos: [],
        commentInput: '',
      };
    },
    methods: {
      async fetchUserProfile() {
        try {
          const username = this.$route.params.username;
          const response = await this.$axios.get(`/user/${username}/profile`);
          const data = response.data;
          console.log(data)
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
      formatTimestamp(timestamp) {
        // Use moment to format the timestamp
        return moment(timestamp).format('YYYY-MM-DD HH:mm:ss');
      },
      async postComment(PhotoId) {
        try {
          const username = this.$route.params.username;
          const response = await this.$axios.post(
            `/user/${username}/photos/comment?Photoid=${PhotoId}`,
            {
              content: this.commentInput,
            }
          );
  
          if (response.status === 200) {
            // Successfully posted comment, fetch updated user data
            this.fetchUserProfile(); // Corrected method name
            // Clear the comment input
            this.commentInput = '';
          } else {
            console.error('Failed to post comment:', response.statusText);
          }
        } catch (error) {
          const username = this.$route.params.username;
          console.error('Error while posting comment:', error.message, 'username:', username, 'Photoid:', PhotoId);
        }
      },
      showComments(photo) {
        // Use Vue Router to navigate to the comments route
        this.$router.push({
          name: 'CommentPhotoForm', // Replace with the actual name of your Comments route
          params: { PhotoId: photo.PhotoId, username: this.$route.params.username },
        });
      },
      async toggleLike(photoId, liked) {
        try {
          const username = this.$route.params.username;
  
          if (liked) {
            // Send DELETE request to unlike the photo
            await this.$axios.delete(`/user/${username}/photos/likes?Photoid=${photoId}`);
          } else {
            // Send POST request to like the photo
            await this.$axios.post(`/user/${username}/photos/likes?Photoid=${photoId}`);
          }
          window.location.reload();
          // Update dislikeStatus after toggling the like state
        } catch (error) {
          console.error('Error while toggling like:', error.message);
        }
      },
    },
    created() {
      this.fetchUserProfile();
    },
  };
  </script>
  
  <style scoped>
  .profile-container {
    max-width: 800px;
    margin: 0 auto;
  }
  
  .user-details {
  margin-bottom: 20px;
  padding: 20px;
  background-color: #fff;
  border: 1px solid #ddd;
  border-radius: 5px;
  box-shadow: 0 0 10px rgba(0, 0, 0, 0.1);
  display: flex;
  justify-content: space-around;
  align-items: center;
}
  
  .detail {
    font-size: 16px;
    margin: 5px 0;
  }
  
  .photo-list {
    padding: 20px;
    background-color: #fff;
    border: 1px solid #ddd;
    border-radius: 5px;
  }
  
  .photo-heading {
    font-size: 18px;
    margin-bottom: 10px;
  }
  
  .photo-item {
    margin-bottom: 20px;
    display: flex;
    flex-direction: column;
    align-items: center;
  }
  
  .photo-details {
    width: 100%;
    margin-top: 10px;
    display: flex;
    flex-direction: column;
    align-items: center;
  }
  
  .comment-count,
  .like-count,
  .upload-time {
    font-size: 14px;
    margin: 5px 0;
  }
  
  .comment-input {
    width: 100%;
    margin: 10px 0;
    padding: 5px;
  }
  
  .post-comment-button,
  .like-button,
  .dislike-button {
    background-color: #3498db;
    color: #fff;
    padding: 5px 10px;
    border: none;
    border-radius: 3px;
    cursor: pointer;
    margin: 5px 0;
  }
  
  .post-comment-button:hover,
  .like-button:hover,
  .dislike-button:hover {
    background-color: #2980b9;
  }
  
  .upload-time {
    margin-top: 10px;
  }
  </style>
  



