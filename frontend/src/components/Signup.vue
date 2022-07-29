<script>
import {axiosInstance} from "../service/axiosService";

import {validateRegister} from './validator';

import {StatusCodes} from 'http-status-codes';

export default {
  name: 'signup-page',
  data() {
    return {
      signupData: {
        firstName: "",
        lastName: "",
        username: "",
        password: "",
        email: ""
      },
      loading: false,
    }
  },
  created() {
    //check user alrdy auth'd
    if (localStorage.getItem('refreshToken')) {
      this.$notify({type: 'warning', title: 'Logged!', text: "You are already logged in."});
      this.$router.push('/');
    }
  },
  methods: {
    async signup() {
      this.loading = true;

      await axiosInstance.post('/auth/signup', this.signupData)
          .then(
              res => {
                if (res.status === StatusCodes.CONFLICT) {
                  this.$notify({type: 'error', title: 'Error!', text: "Username/email already used"});
                  return;
                }
                this.$notify({clean: true});
                this.$notify({type: 'success', title: 'Sucess!', text: 'User registered successfully'});
                this.$router.push('/login');
              }, err => {
                this.$notify({type: 'error', title: 'Error!', text: "Trouble registering: " + err.toString()});
              }
          );

      this.loading = false;

    },
    reset() {
      this.signupData.firstName = "";
      this.signupData.lastName = "";
      this.signupData.username = "";
      this.signupData.password = "";
      this.signupData.email = "";
    },
    validate() {
      return validateRegister(this.signupData)
    },
    processUserInfo() {
      if (this.validate()) {
        this.signup();
        this.reset();
      }
    }
  }

}
</script>

<style>
@import "./bg.scss";

.log {
  padding: 10vh 5vw;
}
</style>

<template>
  <div class="bg-anim">
    <div v-for="x in 100" v-bind:key="x" class="circle-container">
      <div class="circle"></div>
    </div>
  </div>
  <div class="space-y-10 flex flex-col items-center log">
    <h1 class="font-extrabold text-red-500 text-3xl">
      Сказал — и точка
    </h1>

    <a
        class="tracking-widest font-extrabold text-2xl uppercase rounded-lg text-t-secondary focus:outline-none focus:shadow-outline"
    >
      Регистрация
    </a>


    <div class="form flex-col flex items-center space-y-5 text-t-secondary font-bold">
      <div class="space-y-3">
        <h4 class="block">Имя:</h4>
        <input
            type="firstname"
            name="firstname"
            placeholder="Введите имя"
            class="block rounded focus:outline-none text-secondary bg-t-accent p-4 placeholder:text-primary"
            autocomplete="off"
            v-model="signupData.firstName"
        />
      </div>

      <div class="space-y-3">
        <h4 class="block">Фамилия:</h4>
        <input
            type="lastname"
            name="lastname"
            placeholder="Введите фамилию"
            class="block rounded focus:outline-none text-secondary bg-t-accent p-4 placeholder:text-primary"
            autocomplete="off"
            v-model="signupData.lastName"
        />
      </div>

      <div class="space-y-3">
        <h4 class="block">Логин:</h4>
        <input
            type="username"
            name="username"
            placeholder="Придумайте логин"
            class="block rounded focus:outline-none text-secondary bg-t-accent p-4 placeholder:text-primary"
            autocomplete="off"
            v-model="signupData.username"
        />
      </div>

      <div class="space-y-3">
        <h4 class="block">Пароль:</h4>
        <input
            type="password"
            name="password"
            placeholder="Введите пароль"
            class="block rounded focus:outline-none text-secondary bg-t-accent p-4 placeholder:text-primary"
            autocomplete="off"
            v-model="signupData.password"
        />
      </div>

      <div class="space-y-3">
        <h4 class="block">Email:</h4>
        <input
            type="email"
            name="email"
            placeholder="Введите email"
            class="block rounded focus:outline-none text-secondary bg-t-accent p-4 placeholder:text-primary"
            autocomplete="off"
            v-model="signupData.email"
        />
      </div>

      <div>
        <button
            class="rounded bg-indigo-500 hover:bg-indigo-500/75 shadow-xl  p-3 mt-5 text-gray-300"
            @click="processUserInfo"
        >
          Зарегистрироваться
        </button>
      </div>

      <div>
        <button
            class="rounded bg-cyan-700 hover:bg-cyan-700/75 p-3 mt-5 text-xs text-gray-300"
            @click="$router.push('/login')"
        >
          Назад
        </button>
      </div>

      <span v-if="loading" class="text-red-500 opacity-75 !mt-12">
                <font-awesome-icon icon="circle-notch" size="5x" class="animate-spin"/>
            </span>
    </div>
  </div>
</template>
