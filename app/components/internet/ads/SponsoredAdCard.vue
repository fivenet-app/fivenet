<script lang="ts" setup>
import GenericImg from '~/components/partials/elements/GenericImg.vue';
import type { Ad } from '~~/gen/ts/resources/internet/ads';

defineProps<{
    ad: Ad;
}>();
</script>

<template>
    <UPageCard :title="ad.title" :description="ad.description">
        <template #header>
            <UBadge class="w-full" :label="$t('components.internet.ads.sponsored_ad')" />
        </template>

        <template v-if="ad.image?.filePath" #icon>
            <GenericImg class="mx-auto h-40" :src="ad.image.filePath" :alt="$t('common.logo')" />
        </template>
    </UPageCard>
</template>
