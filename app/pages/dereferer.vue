<script lang="ts" setup>
import { useSettingsStore } from '~/stores/settings';

useHeadSafe({
    title: 'pages.dereferer.title',
    meta: [{ name: 'referrer', content: 'no-referrer' }],
});

definePageMeta({
    title: 'pages.dereferer.title',
    layout: 'landing',
    requiresAuth: false,
    redirectIfAuthed: false,
    showCookieOptions: true,
});

const settingsStore = useSettingsStore();
const { nuiEnabled } = storeToRefs(settingsStore);

const route = useRoute();
const router = useRouter();

const { remaining, start } = useCountdown(6, {
    onComplete: async () => {
        if (nuiEnabled.value) {
            openURLInWindow(route.query.target as string);
            router.back();
        } else {
            await navigateTo(route.query.target as string, { external: true });
        }
    },
});

onMounted(async () => {
    if (!route.query || !route.query.target) {
        await navigateTo('/');
        return;
    }

    start();
});

const target = route.query.target as string;
</script>

<template>
    <div class="flex h-dvh flex-col">
        <div class="hero absolute inset-0 z-[-1] mask-[radial-gradient(100%_100%_at_top,white,transparent)]" />

        <div class="flex w-full flex-1 items-center justify-center">
            <UCard class="w-full max-w-xl bg-white/75 backdrop-blur-sm dark:bg-white/5">
                <template #header>
                    <h3 class="inline-flex gap-2 text-2xl leading-6 font-semibold">
                        <span>{{ $t('pages.dereferer.title') }}</span>
                        <span>-</span>
                        <span>{{ $t('pages.dereferer.subtitle') }}</span>
                    </h3>
                </template>

                <div class="space-y-4">
                    <p>{{ $t('pages.dereferer.description') }}</p>

                    <p>{{ $t('pages.dereferer.countdown', [remaining]) }}</p>

                    <div class="min-w-0 text-sm">
                        <pre class="line-clamp-4">{{ target }}</pre>
                    </div>
                </div>

                <template #footer>
                    <UButton
                        trailing-icon="i-mdi-link-variant"
                        size="lg"
                        :to="nuiEnabled ? undefined : target"
                        external
                        rel="noreferrer"
                        class="mb-2"
                        block
                        @click="nuiEnabled ? openURLInWindow(target) : undefined"
                    >
                        <span>{{ $t('pages.dereferer.goto') }}</span>
                    </UButton>

                    <UButton
                        block
                        color="gray"
                        variant="ghost"
                        leading-icon="i-mdi-arrow-back"
                        :label="$t('common.back')"
                        @click="router.back()"
                    />
                </template>
            </UCard>
        </div>
    </div>
</template>
