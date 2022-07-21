<script>

import axios from 'axios';



axios.defaults.baseURL = '/api/v1';
// axios.defaults.baseURL = 'http://localhost:3000';

export default {
  name: 'login-page',
  data() {
    return {
      username: '',
      password: '',
      loading: false,
    }
  },
  created() {
    /// user alrdy auth'd
    if (localStorage.getItem('auth_token')) {
      this.$notify({type: 'warning', title: 'Logged!', text: "You are already logged in."});
      this.$router.push('/');
    }
  },
  methods: {
    async login() {
      this.loading = true;

      let user = {
        username: this.username,
        password: this.password
      }

      await axios.get(`http://localhost:3000/${axios.defaults.baseURL}`) // it works

      const res = await axios.post('/auth/login', user)
          .then(response => {
            if (response.status === 201) {
              this.$notify({type: 'error', title: 'Error!', text: "Invalid credentials..."});
              return;
            }

            this.$notify({
              clean: true
            })
            this.$notify({type: 'success', title: 'Sucess!', text: 'User logged in.'});

            localStorage.setItem('auth_token', response.data.token);
            this.$store.authToken = response.data.token
            this.$router.push('/');
          }, err => {
            this.$notify({type: 'error', title: 'Error!', text: "Trouble logging in..."});
          });


      this.loading = false;
    },
    reset() {
      this.password = "";
    },
    validate() {
      return !(this.username === '' || this.password === '');
    },
    processUserInfo() {
      if (this.validate()) {
        this.login();
        this.reset();
        return;
      }
      this.$notify({type: 'error', title: 'Error!', text: "Please make sure to fill up the forms!"});
    }
  }
}
</script>

<style>
@import "./bg.scss";

.bg-anim {
  /*height: 0px;*/
  height: 0;
  overflow: hidden;
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
      Сказал и точка
    </h1>

    <a
        class="tracking-widest font-extrabold text-2xl uppercase rounded-lg text-t-secondary focus:outline-none focus:shadow-outline"
    >
      Лучшая социальная сеть? Да и точка.
    </a>

    <div class="form flex-col flex items-center space-y-5 text-t-secondary font-bold">

      <div class="space-y-3">
        <h4 class="block">
          Логин:
        </h4>
        <input
            type="username"
            name="username"
            placeholder="Введите логин"
            class="block rounded focus:outline-none text-secondary bg-t-accent p-4 placeholder:text-primary"
            v-model="username"
        />
      </div>

      <div class="space-y-3">
        <h4 class="block">
          Пароль:
        </h4>
        <input
            type="password"
            name="password"
            placeholder="Введите пароль"
            class="block rounded focus:outline-none text-secondary bg-t-accent p-4 placeholder:text-primary"
            v-model="password"
        />
      </div>

      <div>
        <button
            @click="processUserInfo"
            class="rounded bg-cyan-700 p-3 mt-5 shadow-xl hover:bg-cyan-700/75 text-gray-300"
        >
          Войти
        </button>
      </div>

      <div>
        <button
            class="rounded bg-indigo-500 hover:bg-indigo-500/75  p-3 mt-5 text-xs text-gray-300"
            @click="$router.push('/signup')"
        >
          Зарегистрироваться
        </button>
      </div>

      <span v-if="loading" class="text-red-500 opacity-75 !mt-12">
                <font-awesome-icon icon="circle-notch" size="5x" class="animate-spin"/>
            </span>
    </div>
  </div>
</template>
