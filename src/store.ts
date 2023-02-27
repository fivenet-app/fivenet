import { createStore } from "vuex";
import axios from "axios";

export const store = createStore({
  state: {
    accessToken: null as null | string,
    loggingIn: false,
    loginError: null,
  },
  mutations: {
    loginStart: (state) => (state.loggingIn = true),
    loginStop: (state, errorMessage) => {
      state.loggingIn = false;
      state.loginError = errorMessage;
    },
    updateAccessToken: (state, accessToken) => {
      state.accessToken = accessToken;
    },
  },
  actions: {
    async doLogin({ commit }, loginData) {
      commit("loginStart");
      axios
        .post("http://localhost:8080/auth/login", {
          ...loginData,
        })
        .then((response) => {
          localStorage.accessToken = response.data.token;
          commit("loginStop", null);
          commit("updateAccessToken", response.data.token);
        })
        .catch((error) => {
          commit("loginStop", error.response.data.error);
          commit('updateAccessToken', null);
        });
    },
    fetchAccessToken({ commit }): void {
      commit('updateAccessToken', localStorage.getItem('accessToken'));
    },
  },
});
