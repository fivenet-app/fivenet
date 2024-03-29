<script lang="ts" setup>
import ContentCenterWrapper from '~/components/partials/ContentCenterWrapper.vue';
import ContentHeroFull from '~/components/partials/ContentHeroFull.vue';
import PageFooter from '~/components/partials/PageFooter.vue';
</script>

<template>
    <div class="flex h-full flex-col">
        <ContentHeroFull>
            <ContentCenterWrapper class="mx-auto max-w-2xl text-center">
                <div class="px-5 sm:px-0">
                    <slot name="default" />
                </div>
            </ContentCenterWrapper>
        </ContentHeroFull>

        <PageFooter />
    </div>
</template>
