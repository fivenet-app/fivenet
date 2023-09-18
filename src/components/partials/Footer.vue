<script lang="ts" setup>
import FiveNetLogo from './logos/FiveNetLogo.vue';
</script>

<template>
    <footer class="bg-base-800 w-full">
        <div class="py-6 sm:py-8 flex flex-row items-center mx-auto w-fit gap-5">
            <FiveNetLogo class="w-16 h-auto" />
            <p class="text-center text-sm leading-5 text-neutral">
                {{ $t('copyright') }}
            </p>
        </div>
    </footer>
</template>
