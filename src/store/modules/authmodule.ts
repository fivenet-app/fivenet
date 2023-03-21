import { RpcError } from 'grpc-web';
import { AuthServiceClient } from '@arpanet/gen/services/auth/AuthServiceClientPb';
import { LoginRequest, LogoutRequest } from '@arpanet/gen/services/auth/auth_pb';
import { User } from '@arpanet/gen/resources/users/users_pb';
import { getAuthClient } from '../../grpc';
import config from '../../config';
import { dispatchNotification } from '../../components/notification';
import { RootState } from '../store';
import { Module } from 'vuex';

export interface AuthModuleState {
    accessToken: null | string;
    lastCharID: number;
    activeChar: null | User;
    loggingIn: boolean;
    loginError: null | string;
    permissions: Array<String>;
}

const authModule: Module<AuthModuleState, RootState> = {
    namespaced: true,
    state: {
        // Persisted to Local Storage
        accessToken: null as null | string,
        lastCharID: 0 as number,
        // Temporary
        activeChar: null as null | User,
        loggingIn: false as boolean,
        loginError: null as null | string,
        permissions: [] as Array<String>,
    } as AuthModuleState,
    mutations: {
        loginStart: (state: AuthModuleState) => (state.loggingIn = true),
        loginStop: (state: AuthModuleState, errorMessage: string) => {
            state.loggingIn = false;
            state.loginError = errorMessage;
        },
        updateAccessToken: (state: AuthModuleState, accessToken: string) => {
            state.accessToken = accessToken;
        },
        updateActiveChar: (state: AuthModuleState, char: null | User) => {
            state.activeChar = char;
            state.lastCharID = char ? char.getUserId() : state.lastCharID;
        },
        updatePermissions: (state: AuthModuleState, permissions: string[]) => {
            state.permissions = permissions;
        },
    },
    actions: {
        async doLogin({ commit }, loginData: LoginRequest) {
            commit('loginStart');

            const client = new AuthServiceClient(config.apiProtoURL, null, null);
            return client
                .login(loginData, null)
                .then((response) => {
                    commit('loginStop', null);
                    commit('updateAccessToken', response.getToken());
                    commit('updateActiveChar', null);
                    commit('updatePermissions', []);
                })
                .catch((err: RpcError) => {
                    commit('loginStop', err.message);
                    commit('updateAccessToken', null);
                    commit('updateActiveChar', null);
                    commit('updatePermissions', []);
                });
        },
        async doLogout({ commit }) {
            commit('loginStart');
            commit('updateActiveChar', null);
            commit('updatePermissions', []);

            return getAuthClient()
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
                    dispatchNotification({ title: 'Error during logout!', content: err.message, type: 'error' });
                });
        },
        updateAccessToken({ commit }, token: string): void {
            commit('updateAccessToken', token);
        },
        updateActiveChar({ commit }, char: null | User): void {
            commit('updateActiveChar', char);
        },
        updatePermissions({ commit }, permissions: string[]): void {
            commit('updatePermissions', permissions);
        },
    },
};

export default authModule;
