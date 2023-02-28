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
                    <li><router-link active-class="active" to="/overview">Overview</router-link></li>
                    <li><router-link active-class="active" to="/citizens">Citizens</router-link></li>
                    <li tabindex="0">
                        <router-link active-class="active" to="/documents" class="justify-between">
                            Documents
                            <svg class="fill-current" xmlns="http://www.w3.org/2000/svg" width="24" height="24"
                                viewBox="0 0 24 24">
                                <path d="M8.59,16.58L13.17,12L8.59,7.41L10,6L16,12L10,18L8.59,16.58Z" />
                            </svg>
                        </router-link>
                        <ul class="p-2">
                            <li><router-link active-class="active" to="/">Submenu 1</router-link></li>
                        </ul>
                    </li>
                    <li><router-link active-class="active" to="/livemap">Livemap</router-link></li>
                </ul>
            </div>
            <router-link active-class="active" to="/" class="flex items-center">
                <img src="/images/logo.png" class="h-6 mr-3 sm:h-9" alt="aRPaNet Logo" />
                <span class="self-center text-xl font-semibold whitespace-nowrap">aRPaNet</span>
            </router-link>
        </div>
        <div v-if="accessToken" class="navbar-center hidden lg:flex">
            <ul class="menu menu-horizontal px-1">
                <li><router-link active-class="active" to="/overview">Overview</router-link></li>
                <li><router-link active-class="active" to="/citizens">Citizens</router-link></li>
                <li tabindex="0">
                    <router-link active-class="active" to="/documents">
                        Documents
                        <svg class="fill-current" xmlns="http://www.w3.org/2000/svg" width="20" height="20"
                            viewBox="0 0 24 24">
                            <path d="M7.41,8.58L12,13.17L16.59,8.58L18,10L12,16L6,10L7.41,8.58Z" />
                        </svg>
                    </router-link>
                    <ul class="p-2">
                        <li><router-link active-class="active" to="/">Submenu 1</router-link></li>
                    </ul>
                </li>
                <li><router-link active-class="active" to="/livemap">Livemap</router-link></li>
            </ul>
        </div>
        <div class="navbar-end">
            <ul class="menu menu-horizontal px-1">
                <li tabindex="0">
                    <button>
                        <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" fill="currentColor"
                            class="bi bi-person-circle" viewBox="0 0 16 16">
                            <path d="M11 6a3 3 0 1 1-6 0 3 3 0 0 1 6 0z" />
                            <path fill-rule="evenodd"
                                d="M0 8a8 8 0 1 1 16 0A8 8 0 0 1 0 8zm8-7a7 7 0 0 0-5.468 11.37C3.242 11.226 4.805 10 8 10s4.757 1.225 5.468 2.37A7 7 0 0 0 8 1z" />
                        </svg>
                        Account
                        <svg class="fill-current" xmlns="http://www.w3.org/2000/svg" width="20" height="20"
                            viewBox="0 0 24 24">
                            <path d="M7.41,8.58L12,13.17L16.59,8.58L18,10L12,16L6,10L7.41,8.58Z" />
                        </svg>
                    </button>
                    <ul class="p-2 bg-base-100">
                        <li><a><label for="modal-navbar-changetheme">Change Theme</label></a></li>
                        <li class="menu-title">
                            <span>Account</span>
                        </li>
                        <li><router-link to="/logout" active-class="active"
                                for="modal-navbar-changetheme">Logout</router-link></li>
                    </ul>
                </li>
            </ul>
        </div>
    </div>
    <input type="checkbox" id="modal-navbar-changetheme" class="modal-toggle" />
    <label for="modal-navbar-changetheme" class="modal cursor-pointer">
        <label class="modal-box w-11/12 max-w-5xl" for="">
            <h2 class="text-lg font-bold">Change Theme</h2>
            <div class="rounded-box grid grid-cols-3 gap-4 sm:grid-cols-3 md:grid-cols-3 lg:grid-cols-3">
                <div class="border-base-content/20 hover:border-base-content/40 outline-base-content overflow-hidden rounded-lg border outline-2 outline-offset-2"
                    data-set-theme="light" data-act-class="outline">
                    <div data-theme="light" class="bg-base-100 text-base-content w-full cursor-pointer font-sans">
                        <div class="grid grid-cols-5 grid-rows-3">
                            <div class="bg-base-200 col-start-1 row-span-2 row-start-1"></div>
                            <div class="bg-base-300 col-start-1 row-start-3"></div>
                            <div class="bg-base-100 col-span-4 col-start-2 row-span-3 row-start-1 flex flex-col gap-1 p-2">
                                <div class="font-bold">light</div>
                                <div class="flex flex-wrap gap-1">
                                    <div
                                        class="bg-primary flex aspect-square w-5 items-center justify-center rounded lg:w-6">
                                        <div class="text-primary-content text-sm font-bold">A</div>
                                    </div>
                                    <div
                                        class="bg-secondary flex aspect-square w-5 items-center justify-center rounded lg:w-6">
                                        <div class="text-secondary-content text-sm font-bold">A</div>
                                    </div>
                                    <div
                                        class="bg-accent flex aspect-square w-5 items-center justify-center rounded lg:w-6">
                                        <div class="text-accent-content text-sm font-bold">A</div>
                                    </div>
                                    <div
                                        class="bg-neutral flex aspect-square w-5 items-center justify-center rounded lg:w-6">
                                        <div class="text-neutral-content text-sm font-bold">A</div>
                                    </div>
                                </div>
                            </div>
                        </div>
                    </div>
                </div>
                <div class="border-base-content/20 hover:border-base-content/40 outline-base-content overflow-hidden rounded-lg border outline-2 outline-offset-2"
                    data-set-theme="dark" data-act-class="outline">
                    <div data-theme="dark" class="bg-base-100 text-base-content w-full cursor-pointer font-sans">
                        <div class="grid grid-cols-5 grid-rows-3">
                            <div class="bg-base-200 col-start-1 row-span-2 row-start-1"></div>
                            <div class="bg-base-300 col-start-1 row-start-3"></div>
                            <div class="bg-base-100 col-span-4 col-start-2 row-span-3 row-start-1 flex flex-col gap-1 p-2">
                                <div class="font-bold">dark</div>
                                <div class="flex flex-wrap gap-1">
                                    <div
                                        class="bg-primary flex aspect-square w-5 items-center justify-center rounded lg:w-6">
                                        <div class="text-primary-content text-sm font-bold">A</div>
                                    </div>
                                    <div
                                        class="bg-secondary flex aspect-square w-5 items-center justify-center rounded lg:w-6">
                                        <div class="text-secondary-content text-sm font-bold">A</div>
                                    </div>
                                    <div
                                        class="bg-accent flex aspect-square w-5 items-center justify-center rounded lg:w-6">
                                        <div class="text-accent-content text-sm font-bold">A</div>
                                    </div>
                                    <div
                                        class="bg-neutral flex aspect-square w-5 items-center justify-center rounded lg:w-6">
                                        <div class="text-neutral-content text-sm font-bold">A</div>
                                    </div>
                                </div>
                            </div>
                        </div>
                    </div>
                </div>
                <div class="border-base-content/20 hover:border-base-content/40 outline-base-content overflow-hidden rounded-lg border outline-2 outline-offset-2"
                    data-set-theme="corporate" data-act-class="outline">
                    <div data-theme="corporate" class="bg-base-100 text-base-content w-full cursor-pointer font-sans">
                        <div class="grid grid-cols-5 grid-rows-3">
                            <div class="bg-base-200 col-start-1 row-span-2 row-start-1"></div>
                            <div class="bg-base-300 col-start-1 row-start-3"></div>
                            <div class="bg-base-100 col-span-4 col-start-2 row-span-3 row-start-1 flex flex-col gap-1 p-2">
                                <div class="font-bold">corporate</div>
                                <div class="flex flex-wrap gap-1">
                                    <div
                                        class="bg-primary flex aspect-square w-5 items-center justify-center rounded lg:w-6">
                                        <div class="text-primary-content text-sm font-bold">A</div>
                                    </div>
                                    <div
                                        class="bg-secondary flex aspect-square w-5 items-center justify-center rounded lg:w-6">
                                        <div class="text-secondary-content text-sm font-bold">A</div>
                                    </div>
                                    <div
                                        class="bg-accent flex aspect-square w-5 items-center justify-center rounded lg:w-6">
                                        <div class="text-accent-content text-sm font-bold">A</div>
                                    </div>
                                    <div
                                        class="bg-neutral flex aspect-square w-5 items-center justify-center rounded lg:w-6">
                                        <div class="text-neutral-content text-sm font-bold">A</div>
                                    </div>
                                </div>
                            </div>
                        </div>
                    </div>
                </div>
                <div class="border-base-content/20 hover:border-base-content/40 outline-base-content overflow-hidden rounded-lg border outline-2 outline-offset-2"
                    data-set-theme="business" data-act-class="outline">
                    <div data-theme="business" class="bg-base-100 text-base-content w-full cursor-pointer font-sans">
                        <div class="grid grid-cols-5 grid-rows-3">
                            <div class="bg-base-200 col-start-1 row-span-2 row-start-1"></div>
                            <div class="bg-base-300 col-start-1 row-start-3"></div>
                            <div class="bg-base-100 col-span-4 col-start-2 row-span-3 row-start-1 flex flex-col gap-1 p-2">
                                <div class="font-bold">business</div>
                                <div class="flex flex-wrap gap-1">
                                    <div
                                        class="bg-primary flex aspect-square w-5 items-center justify-center rounded lg:w-6">
                                        <div class="text-primary-content text-sm font-bold">A</div>
                                    </div>
                                    <div
                                        class="bg-secondary flex aspect-square w-5 items-center justify-center rounded lg:w-6">
                                        <div class="text-secondary-content text-sm font-bold">A</div>
                                    </div>
                                    <div
                                        class="bg-accent flex aspect-square w-5 items-center justify-center rounded lg:w-6">
                                        <div class="text-accent-content text-sm font-bold">A</div>
                                    </div>
                                    <div
                                        class="bg-neutral flex aspect-square w-5 items-center justify-center rounded lg:w-6">
                                        <div class="text-neutral-content text-sm font-bold">A</div>
                                    </div>
                                </div>
                            </div>
                        </div>
                    </div>
                </div>
                <div class="border-base-content/20 hover:border-base-content/40 outline-base-content overflow-hidden rounded-lg border outline-2 outline-offset-2"
                    data-set-theme="acid" data-act-class="outline">
                    <div data-theme="acid" class="bg-base-100 text-base-content w-full cursor-pointer font-sans">
                        <div class="grid grid-cols-5 grid-rows-3">
                            <div class="bg-base-200 col-start-1 row-span-2 row-start-1"></div>
                            <div class="bg-base-300 col-start-1 row-start-3"></div>
                            <div class="bg-base-100 col-span-4 col-start-2 row-span-3 row-start-1 flex flex-col gap-1 p-2">
                                <div class="font-bold">acid</div>
                                <div class="flex flex-wrap gap-1">
                                    <div
                                        class="bg-primary flex aspect-square w-5 items-center justify-center rounded lg:w-6">
                                        <div class="text-primary-content text-sm font-bold">A</div>
                                    </div>
                                    <div
                                        class="bg-secondary flex aspect-square w-5 items-center justify-center rounded lg:w-6">
                                        <div class="text-secondary-content text-sm font-bold">A</div>
                                    </div>
                                    <div
                                        class="bg-accent flex aspect-square w-5 items-center justify-center rounded lg:w-6">
                                        <div class="text-accent-content text-sm font-bold">A</div>
                                    </div>
                                    <div
                                        class="bg-neutral flex aspect-square w-5 items-center justify-center rounded lg:w-6">
                                        <div class="text-neutral-content text-sm font-bold">A</div>
                                    </div>
                                </div>
                            </div>
                        </div>
                    </div>
                </div>
                <div class="border-base-content/20 hover:border-base-content/40 outline-base-content overflow-hidden rounded-lg border outline-2 outline-offset-2"
                    data-set-theme="night" data-act-class="outline">
                    <div data-theme="night" class="bg-base-100 text-base-content w-full cursor-pointer font-sans">
                        <div class="grid grid-cols-5 grid-rows-3">
                            <div class="bg-base-200 col-start-1 row-span-2 row-start-1"></div>
                            <div class="bg-base-300 col-start-1 row-start-3"></div>
                            <div class="bg-base-100 col-span-4 col-start-2 row-span-3 row-start-1 flex flex-col gap-1 p-2">
                                <div class="font-bold">night</div>
                                <div class="flex flex-wrap gap-1">
                                    <div
                                        class="bg-primary flex aspect-square w-5 items-center justify-center rounded lg:w-6">
                                        <div class="text-primary-content text-sm font-bold">A</div>
                                    </div>
                                    <div
                                        class="bg-secondary flex aspect-square w-5 items-center justify-center rounded lg:w-6">
                                        <div class="text-secondary-content text-sm font-bold">A</div>
                                    </div>
                                    <div
                                        class="bg-accent flex aspect-square w-5 items-center justify-center rounded lg:w-6">
                                        <div class="text-accent-content text-sm font-bold">A</div>
                                    </div>
                                    <div
                                        class="bg-neutral flex aspect-square w-5 items-center justify-center rounded lg:w-6">
                                        <div class="text-neutral-content text-sm font-bold">A</div>
                                    </div>
                                </div>
                            </div>
                        </div>
                    </div>
                </div>
            </div>
        </label>
    </label>
</template>
