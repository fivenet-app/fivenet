import { createStore } from 'vuex';

import * as grpcWeb from 'grpc-web';
import { AccountServiceClient } from '@arpanet/gen/auth/AuthServiceClientPb';
import { Character, LoginRequest } from '@arpanet/gen/auth/auth_pb';

import { version } from '../package.json';

export const store = createStore({
	state: {
		version: '',
		accessToken: null as null | string,
		loggingIn: false,
		loginError: null,
		chars: [] as Array<Character>,
		activeChar: "" as string,
	},
	mutations: {
		initialiseStore(state) {
			// Check if the store exists
			if (localStorage.getItem('store')) {
				let store = JSON.parse(localStorage.getItem('store'));

				// Check the version stored against current. If different, don't
				// load the cached version
				if (store.version == version) {
					this.replaceState(Object.assign(state, store));
				} else {
					state.version = version;
				}
			}
		},
		loginStart: (state) => (state.loggingIn = true),
		loginStop: (state, errorMessage) => {
			state.loggingIn = false;
			state.loginError = errorMessage;
		},
		updateAccessToken: (state, accessToken) => {
			state.accessToken = accessToken;
		},
		updateChars: (state, chars) => {
			state.chars = chars;
		},
		updateActiveChar: (state, index: string) => {
			state.activeChar = index;
		},
	},
	actions: {
		async doLogin({ commit }, loginData: LoginRequest) {
			commit('loginStart');
			const client = new AccountServiceClient('https://localhost:8181', null, {});
			return client
				.login(loginData, null)
				.then((response) => {
					commit('loginStop', null);
					commit('updateAccessToken', response.getToken());
					commit('updateChars', response.getCharsList());
					commit('updateActiveChar', "");
				})
				.catch((err: grpcWeb.RpcError) => {
					commit('loginStop', err.message);
					commit('updateAccessToken', null);
					commit('updateChars', []);
					commit('updateActiveChar', "");
				});
		},
		doLogout({ commit }): void {
			commit('updateAccessToken', null);
			commit('updateChars', []);
			commit('updateActiveChar', "");
		},
		updateChars({ commit }, chars) {
			commit('updateChars', chars);
		},
		updateActiveChar({ commit }, identifier: string): void {
			commit('updateActiveChar', identifier);
		},
		updateAccessToken({ commit }, token: string): void {
			commit('updateAccessToken', token);
		},
	},
});

store.subscribe((mutation, state) => {
	let store = {
		version: state.version,
		accessToken: state.accessToken,
		activeChar: state.activeChar,
	};

	localStorage.setItem('store', JSON.stringify(store));
});
