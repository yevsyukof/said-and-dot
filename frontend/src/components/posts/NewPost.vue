<script>
import {axiosInstance} from "../../service/axiosService";
import {isJwtExpired} from 'jwt-check-expiration';
import moment from 'moment';
// import {refreshAccessToken} from "../../../service/tokenService"

import {StatusCodes} from 'http-status-codes';

export default {
  name: 'new-post-item',
  emits: ['closePost', 'newPost'],
  props: {
    user: {
      required: true
    },
    isOpen: {
      required: true
    },
  },
  data() {
    return {
      postContent: '',
      isEmpty: false,
      isLoading: false,
    }
  },
  methods: {
    formatDate(date) {
      return moment(date, 'YYYY-MM-DDThh:mm:ss.sssZ').format('DD MMM YYYY')
    },
    async submitPost() {
      this.isEmpty = false;
      if (this.postContent === '') {
        this.isEmpty = true;
      } else {
        this.isLoading = true;

        // if (isJwtExpired(localStorage.getItem('accessToken'))) {
        //   if (isJwtExpired(localStorage.getItem('refreshToken'))) {
        //     // TODO выкидывать чувака на окно Login
        //   }
        //
        //   refreshData = await refreshAccessToken(isJwtExpired(localStorage.getItem('refreshToken')))
        //   localStorage.removeItem('refreshToken')
        //   localStorage.removeItem('accessToken')
        //
        //   console.assert(refreshData)
        // }

        // //create new post
        await axiosInstance.post('/tweets/',
            {
              tweet: this.postContent
            },
            {
              headers: {
                Authorization: localStorage.getItem('accessToken')
              }
            }
        ).then(
            res => {
              if (res.status === StatusCodes.CREATED) {
                return this.$notify({
                      type: 'success',
                      title: 'Sucess!',
                      text: 'Ваше высказывание услышано!'
                    }
                );
              }
            }, err => {
              this.$notify({type: 'error', title: 'Error!', text: "Trouble posting..."});
            }
        );
        //end create post

        this.isLoading = false;

        this.$emit('newPost', this.postContent);
        this.postContent = "";
        this.$emit('closePost');
      }
    }
  }
}
</script>

<style scoped>
.fadeHeight-enter-active,
.fadeHeight-leave-active {
  transition: all 0.2s;
  max-height: 230px;
}

.fadeHeight-enter-from,
.fadeHeight-leave-to {
  opacity: 0;
  max-height: 0px;
}
</style>

<template>
  <!-- Менюшка написания поста -->
  <transition appear name="fadeHeight" mode="in-out">
    <div
        v-if="isOpen"
        class="bg-accent shadow-lg rounded-lg md:mx-auto my-5 md:max-w-2xl"
    >
      <div class="flex items-start px-4 py-6">
        <img
            class="w-12 h-12 rounded-full object-cover mr-4 shadow float-left"
            v-bind:src="user.avatar"
            alt="avatar"
        />

        <div class="w-full">
          <div class="flex items-center justify-between">
            <h2
                class="text-lg font-semibold text-t-primary -mt-1"
            >
              {{ user.firstName + ' ' + user.lastName }}
            </h2>
            <small class="text-sm text-t-accent">Just now</small>
          </div>

          <p class="text-t-accent text-xs">Joined {{ formatDate(user.createdAt) }}.</p>

          <textarea
              v-model="this.postContent"
              :class="{ 'border-red-500 border-2': this.isEmpty }"
              class="mt-3 text-t-secondary rounded w-full h-32 sm:h-32 p-3 outline-none resize-none text-sm bg-primary"
              :disabled="this.isLoading"
          ></textarea>

          <div class="mt-2 flex items-center justify-between">
            <div class="flex mr-2 text-t-accent text-xs">
              <svg
                  fill="none"
                  viewBox="0 0 24 24"
                  class="w-4 h-4 mr-1"
                  stroke="currentColor"
              >
                <path
                    stroke-linecap="round"
                    stroke-linejoin="round"
                    stroke-width="2"
                    d="M4.318 6.318a4.5 4.5 0 000 6.364L12 20.364l7.682-7.682a4.5 4.5 0 00-6.364-6.364L12 7.636l-1.318-1.318a4.5 4.5 0 00-6.364 0z"
                />
              </svg>
            </div>
            <span v-if="this.isEmpty" class="text-red-500">Не оставляйте это поле пустым</span>

            <button
                @click="this.submitPost"
                class="bg-green-700/50 hover:bg-green-700 text-sm font-semibold text-gray-200 px-4 py-2 rounded"
            >
              Высказать
            </button>
          </div>
        </div>
      </div>
    </div>
  </transition>
  <div v-if="this.isLoading" class="text-red-500 opacity-75 flex justify-center my-5">
    <font-awesome-icon icon="circle-notch" size="5x" class="animate-spin"/>
  </div>
</template>
