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
  data() {
    return {
      isUserLoaded: false,
      user: {
        userData: {
          id: 'Loading...',
          username: 'Loading...',
          firstName: 'Loading...',
          lastName: 'Loading...',
          email: 'Loading...',
        },
        followers: [], // TODO
        follows: []
      }
    }
  },
  created() { // TODO это хуки
    if (!localStorage.getItem("refreshToken")) {
      this.$notify({
        type: 'error',
        title: 'No login!',
        text: "Please login first."
      });
      this.$router.push('/login');
    }
  },

  async mounted() { // TODO
    await axiosInstance.get('/users/me', {
          headers: {
            "Authorization": localStorage.getItem('accessToken'),
          }
        }
    ).then(
        res => {
          this.user = res.data

          if (this.user.follows === null) {
            this.user.follows = []
          }
          if (this.user.followers === null) {
            this.user.followers = []
          }

          this.isUserLoaded = true;
          this.$store.dispatch('saveUser', this.user.userData); //обращаемся к store модулю
        }, err => {
          this.$notify({type: 'error', title: 'Error!', text: "Trouble in getting user..."});
        }
    )
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

      if (localStorage.getItem("refreshToken")) {
        axiosInstance.post("/auth/logout", {"refreshToken": localStorage.getItem("refreshToken")})
      }
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
