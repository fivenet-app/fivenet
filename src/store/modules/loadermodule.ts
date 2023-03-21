import { Module } from 'vuex';
import { RootState } from '../store';

export interface LoaderModuleState {
    loading: boolean,
}

const loaderModule: Module<LoaderModuleState, RootState> = {
    namespaced: true,
    state: {
        loading: false,
    } as LoaderModuleState,
    actions: {
        show({ commit }) {
            commit('show');
        },
        hide({ commit }) {
            commit('hide');
        },
    },
    mutations: {
        show(state: LoaderModuleState) {
            state.loading = true;
        },
        hide(state: LoaderModuleState) {
            state.loading = false;
        },
    },
};

export default loaderModule;
