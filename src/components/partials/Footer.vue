<script lang="ts" setup>
const navigation = {
    logo: '/images/logo.png',
};
</script>

<template>
    <footer class="bg-gray-900">
        <div class="mx-auto max-w-7xl overflow-hidden py-14 px-6 sm:py-14 lg:px-8">
            <img class="h-14 mx-auto" :src="navigation.logo" alt="CIP Logo" />
            <p class="mt-10 text-center text-xs leading-5 text-gray-500">
                Copyright &copy; 2023 Galexrt All rights reserved.
            </p>
        </div>
    </footer>
</template>
