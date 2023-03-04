<script lang="ts">
import { defineComponent } from 'vue';

export default defineComponent({
    props: {
        title: {
            required: true,
            type: String,
        },
    },
});
</script>

<template>
    <header class="bg-gray-700 shadow-sm">
      <div class="mx-auto max-w-7xl py-4 px-4 sm:px-6 lg:px-8">
        <h1 class="text-lg font-semibold leading-6 text-white">{{ title }}</h1>
        <slot />
      </div>
    </header>
</template>
