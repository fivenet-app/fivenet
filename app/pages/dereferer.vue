<script lang="ts" setup>
import { useSettingsStore } from '~/store/settings';

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
                openURLInWindow(url);
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
    <div class="flex h-dscreen flex-col">
        <div class="hero absolute inset-0 z-[-1] [mask-image:radial-gradient(100%_100%_at_top,white,transparent)]" />

        <div class="flex w-full flex-1 items-center justify-center">
            <UCard class="w-full max-w-lg bg-white/75 backdrop-blur dark:bg-white/5">
                <template #header>
                    <div class="flex items-center justify-between">
                        <h3 class="text-2xl font-semibold leading-6">
                            {{ $t('pages.dereferer.title') }} - {{ $t('pages.dereferer.subtitle') }}
                        </h3>
                    </div>
                </template>

                <p>{{ $t('pages.dereferer.description') }}</p>

                <div class="mt-4">
                    <pre>{{ target }}</pre>
                </div>

                <template #footer>
                    <div class="inline-flex flex-col gap-2">
                        <UButton
                            :label="$t('pages.dereferer.goto')"
                            trailing-icon="i-mdi-link-variant"
                            size="lg"
                            :to="nuiEnabled ? undefined : target"
                            :external="true"
                            rel="noreferrer"
                            @click="nuiEnabled ? openURLInWindow(target) : undefined"
                        />

                        <UButton color="red" icon="i-mdi-arrow-back" :label="$t('common.back')" @click="router.back()" />
                    </div>
                </template>
            </UCard>
        </div>
    </div>
</template>
