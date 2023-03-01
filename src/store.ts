import { createStore } from 'vuex';

import { RpcError } from 'grpc-web';
import { AccountServiceClient } from '@arpanet/gen/auth/AuthServiceClientPb';
import { LoginRequest, LogoutRequest } from '@arpanet/gen/auth/auth_pb';

import { version } from '../package.json';
import { Character } from '@arpanet/gen/common/character_pb';
import authInterceptor from './grpcauth';

const store = createStore({
	state: {
		// Persisted to Local Storage
		version: '' as string,
		accessToken: null as null | string,
		activeChar: null as Character | null,
		activeCharIdentifier: '' as null | string,
		// Temporary
		loggingIn: false as boolean,
		loginError: null as null | string,
	},
	mutations: {
		initialiseStore(state) {
			// Check if the store exists
			if (localStorage.getItem('store')) {
				let store = JSON.parse(localStorage.getItem('store') as string);

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
		updateActiveChar: (state, char) => {
			state.activeChar = char;
		},
		updateActiveCharIdentifier: (state, identifier: null | string) => {
			state.activeCharIdentifier = identifier;
		},
	},
	actions: {
		async doLogin({ commit }, loginData: LoginRequest) {
			commit('loginStart');

			const client = new AccountServiceClient('https://localhost:8181', null, null);
			return client
				.login(loginData, null)
				.then((response) => {
					commit('loginStop', null);
					commit('updateAccessToken', response.getToken());
					commit('updateActiveChar', null);
					commit('updateActiveCharIdentifier', null);
				})
				.catch((err: RpcError) => {
					commit('loginStop', err.message);
					commit('updateAccessToken', null);
					commit('updateActiveChar', null);
					commit('updateActiveCharIdentifier', null);
				});
		},
		async doLogout({ commit }) {
			commit('loginStart');

			const client = new AccountServiceClient('https://localhost:8181', null, {
				unaryInterceptors: [authInterceptor],
				streamInterceptors: [authInterceptor],
			});
			return client
				.logout(new LogoutRequest(), null)
				.then((response) => {
					commit('loginStop', null);
					commit('updateAccessToken', null);
					commit('updateActiveChar', null);
					commit('updateActiveCharIdentifier', null);
					if (response.getSuccess()) {
						return;
					}
				})
				.catch((err: RpcError) => {
					commit('loginStop', err.message);
					console.log('Error during logout process: ' + err);
				});
		},
		updateActiveChar({ commit }, char: null | Character): void {
			commit('updateActiveChar', char);
		},
		updateActiveCharIdentifier({ commit }, identifier: null | string): void {
			commit('updateActiveCharIdentifier', identifier);
		},
		updateAccessToken({ commit }, token: string): void {
			commit('updateAccessToken', token);
		},
	},
});

export default store;

store.subscribe((mutation, state) => {
	const s = {
		version: state.version,
		accessToken: state.accessToken,
		activeCharIdentifier: state.activeCharIdentifier,
	};

	localStorage.setItem('store', JSON.stringify(s));
});
