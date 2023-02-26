<script lang="ts">
import { themeChange } from 'theme-change';
import { defineComponent } from 'vue';

export default defineComponent({
    mounted: function () {
        themeChange(false);
    },
});
</script>

<template>
    <div class="navbar bg-base-100">
        <div class="flex-1">
            <div class="w-10 rounded-full">
                <img alt="aRPaNet Logo" src="/images/logo.png" />
            </div>
            <a class="btn btn-ghost normal-case text-xl">aRPaNet</a>
        </div>
        <div class="flex-none">
            <ul class="menu menu-horizontal px-1">
                <li tabindex="0">
                    <a>
                        Change Theme
                        <svg class="fill-current" xmlns="http://www.w3.org/2000/svg" width="20" height="20"
                            viewBox="0 0 24 24">
                            <path d="M7.41,8.58L12,13.17L16.59,8.58L18,10L12,16L6,10L7.41,8.58Z" />
                        </svg>
                    </a>
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
