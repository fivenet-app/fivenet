<script lang="ts" setup>
useHeadSafe({
    title: 'pages.dereferer.title',
    meta: [{ name: 'referrer', content: 'no-referrer' }],
});
definePageMeta({
    title: 'pages.dereferer.title',
    layout: 'blank',
    requiresAuth: false,
});

const route = useRoute();
const router = useRouter();

if (!route.query || !route.query.target) {
    await navigateTo('/');
} else {
    const url = route.query.target as string;
    useTimeoutFn(async () => {
        if (isNUIAvailable()) {
            openURLInWindow(url);
            router.back();
        } else {
            await navigateTo(url, { external: true });
        }
    }, 3250);
}

const target = route.query.target as string;
</script>

<template>
    <div class="hero flex h-dscreen flex-col">
        <div class="flex w-full flex-1 items-center justify-center bg-black/50">
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
                    <UButton
                        :label="$t('pages.dereferer.goto')"
                        trailing-icon="i-mdi-link-variant"
                        size="lg"
                        :to="target"
                        :external="true"
                        rel="noreferrer"
                    />
                </template>
            </UCard>
        </div>
    </div>
</template>
