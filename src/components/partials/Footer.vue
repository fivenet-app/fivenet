<script lang="ts">
import { defineComponent } from 'vue'

export default defineComponent({
    data() {
        return {
            'navigation': {
                logo: '/images/logo.png',
            },
        };
    },
});
</script>

<template>
    <footer class="bg-gray-900">
        <div class="mx-auto max-w-7xl overflow-hidden py-14 px-6 sm:py-14 lg:px-8">
            <img class="h-14 mx-auto" :src="navigation.logo" alt="aRPaNet Logo" />
            <p class="mt-10 text-center text-xs leading-5 text-gray-500">
                Copyright &copy; 2023 aRPaNet All rights reserved.
            </p>
            <div class="mt-10 flex justify-center space-x-10">
                <a class="text-gray-400 hover:text-gray-500"><label for="modal-navbar-changetheme">Change Theme</label></a>
            </div>
        </div>
    </footer>
</template>
