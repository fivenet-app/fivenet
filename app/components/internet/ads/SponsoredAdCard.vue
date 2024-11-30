<script lang="ts" setup>
import type { Ad } from '~~/gen/ts/resources/internet/ads';

defineProps<{
    ad: Ad;
}>();
</script>

<template>
    <UPageCard :title="ad.title" :description="ad.description">
        <template #header>
            <UBadge :label="$t('components.internet.ads.sponsored_ad')" class="w-full" />
        </template>

        <template v-if="ad.image?.url" #icon>
            <img :src="ad.image.url" :alt="$t('common.logo')" class="mx-auto h-40" />
        </template>
    </UPageCard>
</template>
