import { createStore } from "vuex";
import axios from "axios";

export const store = createStore({
  state: {
    loggingIn: false,
    loginError: null,
    loginSuccessful: false,
  },
  mutations: {
    loginStart: (state) => (state.loggingIn = true),
    loginStop: (state, errorMessage) => {
      state.loggingIn = false;
      state.loginError = errorMessage;
      state.loginSuccessful = !errorMessage;
    },
  },
  actions: {
    doLogin({ commit }, loginData) {
      commit("loginStart");
      axios
        .post("http://localhost:8080/auth/login", {
          ...loginData,
        })
        .then(() => {
          commit("loginStop", null);
        })
        .catch((error) => {
          commit("loginStop", error.response.data.error);
        });
    },
    doLogout({ commit }) {
      // TODO
    },
  },
});
