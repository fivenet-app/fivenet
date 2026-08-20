<script lang="ts" setup>
import ExternalFileHolder from '~/components/partials/ExternalFileHolder.vue';

useHead({
    title: 'common.licenses',
});

definePageMeta({
    title: 'common.licenses',
    layout: 'landing',
    requiresAuth: false,
    redirectIfAuthed: false,
    showCookieOptions: true,
});

const licenses = [
    {
        label: 'FiveNet Project',
        path: '/licenses/LICENSE.txt',
    },
    {
        label: 'Frontend',
        path: '/licenses/frontend.txt',
    },
    {
        label: 'Backend',
        path: '/licenses/backend.txt',
    },
    {
        label: 'Images & Icons',
        path: '/licenses/images.txt',
    },
    {
        label: 'Fonts',
        path: '/licenses/fonts.txt',
    },
    {
        label: 'Sounds',
        path: '/licenses/sounds.txt',
    },
];
</script>

<template>
    <div class="relative min-h-dvh overflow-hidden">
        <div class="hero pointer-events-none absolute inset-0 z-[-1]" />

        <UPage :ui="{ root: 'pb-(--page-content-bottom-offset)' }">
            <UPageHero
                :title="$t('common.licenses')"
                :description="$t('pages.about.licenses.subtitle')"
                orientation="horizontal"
                :links="[
                    {
                        label: $t('common.about'),
                        icon: 'i-mdi-arrow-left',
                        size: 'lg',
                        color: 'gray',
                        variant: 'solid',
                        to: '/about',
                    },
                ]"
                :ui="{
                    container: 'py-24 sm:py-24 lg:py-24',
                    description: 'text-(--ui-text-highlighted)',
                    title: 'text-4xl sm:text-6xl',
                }"
            >
                <template #top>
                    <div
                        class="absolute left-1/2 size-60 -translate-x-1/2 -translate-y-80 transform rounded-full blur-[300px] sm:size-80 dark:bg-(--ui-primary)"
                    />
                </template>
            </UPageHero>

            <div class="relative">
                <div
                    class="pointer-events-none inset-x-0 h-20"
                    :style="{ background: 'linear-gradient(to bottom, transparent, var(--ui-bg))' }"
                />

                <UPageSection class="bg-default" :ui="{ container: 'py-8 sm:py-8 lg:py-8' }">
                    <UCard class="mx-auto w-full max-w-5xl bg-white/70 shadow-sm backdrop-blur-sm dark:bg-white/5">
                        <UAccordion :items="licenses" type="multiple" :ui="{ content: 'pb-0' }">
                            <template #content="{ item: license }">
                                <div class="space-y-5 border-t border-neutral/10 px-4 py-5 sm:px-6">
                                    <ExternalFileHolder :path="license.path" />

                                    <div class="flex justify-end border-t border-neutral/10 pt-4">
                                        <a
                                            class="inline-flex items-center gap-2 rounded-md border border-neutral/15 bg-neutral/5 px-3 py-2 text-sm font-medium text-highlighted transition hover:bg-neutral/10"
                                            :href="license.path"
                                            download
                                        >
                                            <UIcon name="i-mdi-download" class="size-4" />
                                            <span>{{ $t('common.download') }}</span>
                                        </a>
                                    </div>
                                </div>
                            </template>
                        </UAccordion>
                    </UCard>
                </UPageSection>
            </div>
        </UPage>
    </div>
</template>
