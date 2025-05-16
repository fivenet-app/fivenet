<script lang="ts" setup>
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

        <template v-if="ad.image?.url" #icon>
            <NuxtImg class="mx-auto h-40" :src="ad.image.url" :alt="$t('common.logo')" loading="lazy" />
        </template>
    </UPageCard>
</template>
