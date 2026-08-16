<script lang="ts" setup>
import PageHeader from '~/components/landing/PageHeader.vue';
import BannerMessage from '~/components/partials/BannerMessage.vue';
import PageFooter from '~/components/partials/PageFooter.vue';

const { system } = useAppConfig();
</script>

<template>
    <div>
        <PageHeader />

        <div>
            <slot />
        </div>

        <BannerMessage v-if="system.bannerMessageEnabled && system.bannerMessage" :message="system.bannerMessage" />

        <ClientOnly>
            <LazyPartialsEventsLayer />
        </ClientOnly>

        <PageFooter />
    </div>
</template>
