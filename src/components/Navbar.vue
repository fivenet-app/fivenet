<script lang="ts">
import { themeChange } from 'theme-change';
import { defineComponent } from 'vue';
import { mapState } from 'vuex';

export default defineComponent({
    computed: {
        ...mapState({
            accessToken: 'accessToken',
        }),
    },
    mounted: function () {
        themeChange(false);
    },
});
</script>

<template>
    <div class="navbar bg-base-100">
        <div class="navbar-start">
            <div v-if="accessToken" class="dropdown">
                <label tabindex="0" class="btn btn-ghost lg:hidden">
                    <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24"
                        stroke="currentColor">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6h16M4 12h8m-8 6h16" />
                    </svg>
                </label>
                <ul tabindex="0" class="menu menu-compact dropdown-content mt-3 p-2 shadow bg-base-100 rounded-box w-52">
                    <li><router-link to="/overview">Overview</router-link></li>
                    <li><router-link to="/citizens">Citizens</router-link></li>
                    <li tabindex="0">
                        <router-link to="/documents" class="justify-between">
                            Documents
                            <svg class="fill-current" xmlns="http://www.w3.org/2000/svg" width="24" height="24"
                                viewBox="0 0 24 24">
                                <path d="M8.59,16.58L13.17,12L8.59,7.41L10,6L16,12L10,18L8.59,16.58Z" />
                            </svg>
                        </router-link>
                        <ul class="p-2">
                            <li><router-link to="/">Submenu 1</router-link></li>
                        </ul>
                    </li>
                    <li><router-link to="/livemap">Livemap</router-link></li>
                </ul>
            </div>
            <router-link to="/" class="flex items-center">
                <img src="/images/logo.png" class="h-6 mr-3 sm:h-9" alt="aRPaNet Logo" />
                <span class="self-center text-xl font-semibold whitespace-nowrap">aRPaNet</span>
            </router-link>
        </div>
        <div v-if="accessToken" class="navbar-center hidden lg:flex">
            <ul class="menu menu-horizontal px-1">
                <li><router-link to="/overview">Overview</router-link></li>
                <li><router-link to="/citizens">Citizens</router-link></li>
                <li tabindex="0">
                    <router-link to="/documents">
                        Documents
                        <svg class="fill-current" xmlns="http://www.w3.org/2000/svg" width="20" height="20"
                            viewBox="0 0 24 24">
                            <path d="M7.41,8.58L12,13.17L16.59,8.58L18,10L12,16L6,10L7.41,8.58Z" />
                        </svg>
                    </router-link>
                    <ul class="p-2">
                        <li><router-link to="/">Submenu 1</router-link></li>
                    </ul>
                </li>
                <li><router-link to="/livemap">Livemap</router-link></li>
            </ul>
        </div>
        <div class="navbar-end">
            <ul class="menu menu-horizontal px-1">
                <li tabindex="0">
                    <button>
                        Change Theme
                        <svg class="fill-current" xmlns="http://www.w3.org/2000/svg" width="20" height="20"
                            viewBox="0 0 24 24">
                            <path d="M7.41,8.58L12,13.17L16.59,8.58L18,10L12,16L6,10L7.41,8.58Z" />
                        </svg>
                    </button>
                    <ul class="p-2 bg-base-100">
                        <li><button data-set-theme="dark" data-act-class="ACTIVECLASS">Dark</button></li>
                        <li><button data-set-theme="cyberpunk" data-act-class="ACTIVECLASS">Cyberpunk</button></li>
                        <li><button data-set-theme="light" data-act-class="ACTIVECLASS">Light</button></li>
                    </ul>
                </li>
            </ul>
        </div>
    </div>
</template>
