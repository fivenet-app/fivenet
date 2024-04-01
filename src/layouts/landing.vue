<script lang="ts" setup>
import PageHeader from '~/components/landing/PageHeader.vue';
import PageFooter from '~/components/partials/PageFooter.vue';
</script>

<template>
    <div>
        <PageHeader />

        <div>
            <slot />
        </div>

        <PageFooter />
    </div>
</template>
