<script lang="ts" setup>
import PageHeader from '~/components/landing/PageHeader.vue';
import PageFooter from '~/components/partials/PageFooter.vue';
</script>

<template>
    <div class="hero">
        <PageHeader />

        <slot />

        <PageFooter />

        <ClientOnly>
            <LazyPartialsEventsLayer />
        </ClientOnly>
    </div>
</template>
