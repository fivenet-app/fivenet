import { Module } from 'vuex';
import { RootState } from '../store';

export interface ClipboardModuleState {
    loading: boolean,
}

const clipboardModule: Module<ClipboardModuleState, RootState> = {
    namespaced: true,
    state: {
        loading: false,
    } as ClipboardModuleState,
    actions: {
        show({ commit }) {
            commit('show');
        },
    },
    mutations: {
        show(state: ClipboardModuleState) {
            state.loading = true;
        },
    },
};

export default clipboardModule;
