<script lang="ts" setup>
import { useLivemapStore } from '~/stores/livemap';
import type { Postal } from '~/types/livemap';
import { NotificationType } from '~~/gen/ts/resources/notifications/notifications';

defineOptions({
    inheritAttrs: false,
});

const notifications = useNotificationsStore();

const livemapStore = useLivemapStore();
const { location } = storeToRefs(livemapStore);

const filteredPostals = ref<Postal[]>([]);

const selectedPostal = ref<Postal | undefined>();
const postalQuery = ref('');

const {
    data: postals,
    status,
    execute,
} = await useLazyAsyncData(
    'postals',
    () =>
        $fetch<Postal[]>('/data/postals.json').catch(() =>
            notifications.add({
                title: { key: 'notifications.livemap.failed_loading_postals.title', parameters: {} },
                description: { key: 'notifications.livemap.failed_loading_postals.content', parameters: {} },
                type: NotificationType.ERROR,
            }),
        ),
    {
        immediate: false,
    },
);

function onOpen(): void {
    if (postals.value && postals.value.length > 0) return;

    execute();
}

async function findPostal(q: string): Promise<void> {
    if (q === '' || !postals.value) return;

    let count = 0;
    filteredPostals.value = postals.value.filter((p) => {
        if (count >= 10) return false;

        const result = p.code.startsWith(postalQuery.value!);
        if (result) count++;

        return result;
    });
}

watch(selectedPostal, () => {
    if (!selectedPostal.value) {
        return;
    }

    location.value = selectedPostal.value;
});

watchOnce(postalQuery, () => onOpen());
watchDebounced(postalQuery, (q) => findPostal(q), {
    debounce: 250,
    maxWait: 750,
});
</script>

<template>
    <ClientOnly>
        <UInputMenu
            v-model="selectedPostal"
            v-model:search-term.trim="postalQuery"
            :items="filteredPostals"
            label-key="code"
            class="w-full max-w-40"
            nullable
            :loading="status === 'pending'"
            :placeholder="`${$t('common.postal')} ${$t('common.search')}`"
            :search-input="{ placeholder: $t('common.search_field') }"
            size="xs"
            leading-icon="i-mdi-postage-stamp"
            v-bind="$attrs"
            @update:open="onOpen"
        >
            <template #empty> {{ $t('common.not_found', [$t('common.postal', 2)]) }} </template>
        </UInputMenu>
    </ClientOnly>
</template>
