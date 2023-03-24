<script lang="ts" setup>
const navigation = {
    logo: '/images/logo.png',
};
</script>

<template>
    <footer class="bg-base-850">
        <div class="py-6 sm:py-8 flex flex-row items-center mx-auto w-fit gap-5">
            <img class="w-16 h-auto" :src="navigation.logo" alt="CIP Logo" />
            <p class=" text-center text-s leading-5 text-neutral">
                Copyright &copy; 2023 Galexrt All rights reserved.
            </p>
        </div>
    </footer>
</template>
