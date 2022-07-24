<script>
import Sidebar from './sidebar/Sidebar.vue';
import Header from './header/Header.vue';

import {axiosInstance} from "../service/axiosService";

export default {
  name: 'main-view',
  components: {
    Header,
    Sidebar,
  },
  setup() {
  },
  data() {
    return {
      isUserLoaded: false,
      user: {
        id: 'Loading...', // поменял c _id
        username: 'Loading...',
        firstname: 'Loading...',
        lastname: 'Loading...',
        avatar: ''
      }
    }
  },
  created() {
    //user not auth'd
    if (!localStorage.getItem('auth_token')) {
      this.$notify({
        type: 'error',
        title: 'No login!',
        text: "Please log in first."
      });
      this.$router.push('/login');
    }
  },
  async mounted() {
    await axiosInstance.get('/auth/user', {
          headers: {
            token: localStorage.getItem('auth_token')
          }
        }
    ).then(
        res => {
          this.user = res.data.user;
          this.user.myFollowings = this.user.followings.length;
          this.isUserLoaded = true;
          this.$store.dispatch('saveUser', this.user);
        }, err => {
          this.$notify({type: 'error', title: 'Error!', text: "Trouble in getting user..."});
        })
  },
  methods: {
    logout() {
      this.$notify({clean: true});
      this.$notify({
        type: 'warning',
        title: 'Logged out',
        text: "Logged out successfully."
      });
      this.$router.push('/login');
      localStorage.clear();
    }
  }
}
</script>

<style>
@import "./bg.scss";

.bg-anim {
  overflow: hidden;
  width: 80vw;
}
</style>

<template>
  <div class="bg-anim flex">
    <div v-for="x in 100" v-bind:key="x" class="circle-container">
      <div class="circle"></div>
    </div>
  </div>

  <div class="md:flex">
    <Sidebar/>
    <div class="px-6 py-4 space-y-3.5 flex-row md:flex-col md:w-full">
      <Header :isUserLoaded="this.isUserLoaded" v-bind:user="this.user"/>

      <div class="page-container py-3">
        <router-view :isUserLoaded="this.isUserLoaded" v-bind:user="this.user"/>
      </div>
    </div>
  </div>
</template>
