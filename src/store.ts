import { createStore } from 'vuex';
import { RpcError } from 'grpc-web';
import { AccountServiceClient } from '@arpanet/gen/auth/AuthServiceClientPb';
import { LoginRequest, LogoutRequest } from '@arpanet/gen/auth/auth_pb';
import { version } from '../package.json';
import { Character } from '@arpanet/gen/common/character_pb';
import { clientAuthOptions } from './grpc';
import config from './config';

const store = createStore({
    state: {
        // Persisted to Local Storage
        version: '' as string,
        accessToken: null as null | string,
        activeChar: null as Character | null,
        activeCharID: '' as null | string,
        // Temporary
        loggingIn: false as boolean,
        loginError: null as null | string,
        permissions: [] as Array<String>,
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
        updateActiveCharID: (state, identifier: null | string) => {
            state.activeCharID = identifier;
        },
        updatePermissions: (state, permissions: string[]) => {
            state.permissions = permissions;
        },
    },
    actions: {
        async doLogin({ commit }, loginData: LoginRequest) {
            commit('loginStart');

            const client = new AccountServiceClient(config.apiProtoURL, null, null);
            return client
                .login(loginData, null)
                .then((response) => {
                    commit('loginStop', null);
                    commit('updateAccessToken', response.getToken());
                    commit('updateActiveChar', null);
                    commit('updateActiveCharID', null);
                    commit('updatePermissions', []);
                })
                .catch((err: RpcError) => {
                    commit('loginStop', err.message);
                    commit('updateAccessToken', null);
                    commit('updateActiveChar', null);
                    commit('updateActiveCharID', null);
                    commit('updatePermissions', []);
                });
        },
        async doLogout({ commit }) {
            commit('loginStart');
            commit('updateActiveChar', null);
            commit('updateActiveCharID', null);
            commit('updatePermissions', []);

            return new AccountServiceClient(config.apiProtoURL, null, clientAuthOptions)
                .logout(new LogoutRequest(), null)
                .then((response) => {
                    commit('loginStop', null);
                    commit('updateAccessToken', null);

                    if (response.getSuccess()) {
                        return;
                    }
                })
                .catch((err: RpcError) => {
                    commit('loginStop', err.message);
                    commit('updateAccessToken', null);
                    console.log('Error during logout process: ' + err);
                });
        },
        updateAccessToken({ commit }, token: string): void {
            commit('updateAccessToken', token);
        },
        updateActiveChar({ commit }, char: null | Character): void {
            commit('updateActiveChar', char);
        },
        updateActiveCharID({ commit }, identifier: null | string): void {
            commit('updateActiveCharID', identifier);
        },
        updatePermissions({ commit }, permissions: string[]): void {
            commit('updatePermissions', permissions);
        },
    },
});

export default store;

store.subscribe((mutation, state) => {
    const s = {
        version: state.version,
        accessToken: state.accessToken,
        activeCharID: state.activeCharID,
    };

    localStorage.setItem('store', JSON.stringify(s));
});
