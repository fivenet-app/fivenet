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

onMounted(async () => {
    if (!route.query || !route.query.target) {
        await navigateTo('/');
    } else {
        const url = route.query.target as string;
        useTimeoutFn(async () => {
            if (nuiEnabled.value) {
                await openURLInWindow(url);
                router.back();
                return;
            }

            await navigateTo(url, { external: true });
        }, 5000);
    }
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
                        <span>{{ $t('pages.dereferer.title') }}</span> - <span>{{ $t('pages.dereferer.subtitle') }}</span>
                    </h3>
                </template>

                <p>{{ $t('pages.dereferer.description') }}</p>

                <div class="mt-4">
                    <pre>{{ target }}</pre>
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
                        color="error"
                        leading-icon="i-mdi-arrow-back"
                        block
                        :label="$t('common.back')"
                        @click="router.back()"
                    />
                </template>
            </UCard>
        </div>
    </div>
</template>
