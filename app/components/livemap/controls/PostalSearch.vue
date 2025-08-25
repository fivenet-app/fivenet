<script lang="ts" setup>
import { useLivemapStore } from '~/stores/livemap';
import { NotificationType } from '~~/gen/ts/resources/notifications/notifications';

defineOptions({
    inheritAttrs: false,
});

const notifications = useNotificationsStore();

const livemapStore = useLivemapStore();
const { location } = storeToRefs(livemapStore);

type Postal = {
    x: number;
    y: number;
    code: string;
};

const postalsLoaded = ref(false);
const postals: Postal[] = [];
const filteredPostals = ref<Postal[]>([]);

const selectedPostal = ref<Postal | undefined>();
const postalQuery = ref('');

async function loadPostals(): Promise<void> {
    if (postalsLoaded.value) {
        return;
    }

    try {
        const response = await fetch('/data/postals.json');
        postals.push(...((await response.json()) as Postal[]));
        postalsLoaded.value = true;
    } catch (_) {
        notifications.add({
            title: { key: 'notifications.livemap.failed_loading_postals.title', parameters: {} },
            description: { key: 'notifications.livemap.failed_loading_postals.content', parameters: {} },
            type: NotificationType.ERROR,
        });
        postalsLoaded.value = false;
    }
}

async function findPostal(): Promise<void> {
    if (postalQuery.value === '') {
        return;
    }

    let results = 0;
    filteredPostals.value.length = 0;
    filteredPostals.value = postals.filter((p) => {
        if (results >= 10) {
            return false;
        }

        const result = p.code.startsWith(postalQuery.value!);
        if (result) {
            results++;
        }

        return result;
    });
}

watch(selectedPostal, () => {
    if (!selectedPostal.value) {
        return;
    }

    location.value = selectedPostal.value;
});

watchOnce(postalQuery, async () => loadPostals());
watchDebounced(postalQuery, () => findPostal(), {
    debounce: 250,
    maxWait: 750,
});
</script>

<template>
    <ClientOnly>
        <UInputMenu
            v-model="selectedPostal"
            v-model:query="postalQuery"
            class="w-full max-w-40"
            :items="filteredPostals"
            nullable
            :placeholder="`${$t('common.postal')} ${$t('common.search')}`"
            option-attribute="code"
            :searchable-placeholder="$t('common.search_field')"
            size="xs"
            leading-icon="i-mdi-postage-stamp"
            v-bind="$attrs"
        >
            <template #empty> {{ $t('common.not_found', [$t('common.postal', 2)]) }} </template>
        </UInputMenu>
    </ClientOnly>
</template>
