
<template>
    <div class="nav">
      <div class="dropdown">
        <a @click="toggleDropdown('updateUsername')">Update Username</a>
        <div v-show="activeDropdown === 'updateUsername'" class="dropdown-content">
          <UpdateUsernameForm />
        </div>
      </div>
  
      <div class="dropdown">
        <a @click="toggleDropdown('follow')">Follow</a>
        <div v-show="activeDropdown === 'follow'" class="dropdown-content">
          <FollowUser />
        </div>
      </div>
  
      <div class="dropdown">
        <a @click="toggleDropdown('unfollow')">Unfollow</a>
        <div v-show="activeDropdown === 'unfollow'" class="dropdown-content">
          <UnFollowUser />
        </div>
      </div>
  
      <div class="dropdown">
        <a @click="toggleDropdown('ban')">Ban</a>
        <div v-show="activeDropdown === 'ban'" class="dropdown-content">
          <BanUser />
        </div>
      </div>
      <div class="dropdown">
        <a @click="toggleDropdown('unban')">UnBan</a>
        <div v-show="activeDropdown === 'unban'" class="dropdown-content">
          <UnBanUser />
        </div>
      </div>
      <div class="dropdown">
        <a @click="toggleDropdown('upload')">Upload Photo</a>
        <div v-show="activeDropdown === 'upload'" class="dropdown-content">
          <UploadPhotoForm />
        </div>
      </div>
    </div>
  </template>
  
  <script>
  import UpdateUsernameForm from "../components/UpdateUsernameForm.vue";
  import FollowUser from "../components/FollowUser.vue";
  import UnFollowUser from "../components/UnFollowUser.vue";
  import BanUser from "../components/BanUser.vue";
  import UnBanUser from "../components/UnBanUser.vue";
  import UploadPhotoForm from "../components/UploadPhotoForm.vue";
  
  export default {
    name: "UserHeader",
    components: {
      UpdateUsernameForm,
      FollowUser,
      UnFollowUser,
      BanUser,
      UnBanUser,
      UploadPhotoForm
    },
    data() {
      return {
        activeDropdown: null,
      };
    },
    methods: {
      toggleDropdown(name) {
        // Toggle visibility of the dropdown content
        this.activeDropdown = this.activeDropdown === name ? null : name;
      },
    },
  };
  </script>



  
  <style scoped>
      .nav a:hover {
    background-color: #3B5998; /* Facebook Dark Blue */
  }
  .nav {
    background-color: aqua;
    overflow: hidden;
    margin-top: 0%;
    height: 50px;
    font-size: large;
  }
  
  .dropdown {
    float: left;
  }
  
  .dropdown-content {
    display: none;
    position: absolute;
    background-color: #f9f9f9;
    min-width: 160px;
    box-shadow: 0 8px 16px 0 rgba(0, 0, 0, 0.2);
    z-index: 1;
  }
  
  .dropdown:hover .dropdown-content {
    display: block;
  }
  
  .nav a {
    color: black;
    padding: 14px 16px;
    text-align: center;
    font-size: 17px;
    text-decoration: none;
    margin-right: 5px;
    cursor: pointer;
  }
  </style>
  