<script lang="ts" setup>
import { Combobox, ComboboxInput, ComboboxOption, ComboboxOptions } from '@headlessui/vue';
import { watchDebounced } from '@vueuse/core';
import { useLivemapStore } from '~/store/livemap';
import { useNotificatorStore } from '~/store/notificator';

const notifications = useNotificatorStore();

const livemapStore = useLivemapStore();
const { location } = storeToRefs(livemapStore);

type Postal = {
    x: number;
    y: number;
    code: string;
};

const selectedPostal = ref<Postal | undefined>();
const postalQuery = ref('');
let postalsLoaded = false;
const postals = ref<Postal[]>([]);
const filteredPostals = ref<Postal[]>([]);

async function loadPostals(): Promise<void> {
    if (postalsLoaded) {
        return;
    }
    postalsLoaded = true;

    try {
        const response = await fetch('/data/postals.json');
        postals.value.push(...((await response.json()) as Postal[]));
    } catch (_) {
        notifications.dispatchNotification({
            title: { key: 'notifications.livemap.failed_loading_postals.title', parameters: {} },
            content: { key: 'notifications.livemap.failed_loading_postals.content', parameters: {} },
            type: 'error',
        });
        postalsLoaded = false;
    }
}

async function findPostal(): Promise<void> {
    if (postalQuery.value === '') {
        return;
    }

    let results = 0;
    filteredPostals.value.length = 0;
    filteredPostals.value = postals.value.filter((p) => {
        if (results >= 10) {
            return false;
        }
        const result = p.code.startsWith(postalQuery.value!);
        if (result) results++;
        return result;
    });
}

watch(selectedPostal, () => {
    if (!selectedPostal.value) {
        return;
    }

    location.value = selectedPostal.value;
});

watchDebounced(postalQuery, () => findPostal(), {
    debounce: 250,
    maxWait: 750,
});
</script>

<template>
    <Combobox v-model="selectedPostal" as="div" class="w-full max-w-[11rem]" nullable>
        <ComboboxInput
            autocomplete="off"
            class="w-full p-0.5 px-1 bg-clip-padding rounded-md border-2 border-black/20"
            :display-value="(postal: any) => (postal ? postal?.code : '')"
            :placeholder="`${$t('common.postal')} ${$t('common.search')}`"
            @change="postalQuery = $event.target.value"
            @click="loadPostals"
            @focusin="focusTablet(true)"
            @focusout="focusTablet(false)"
        />
        <ComboboxOptions class="z-10 w-full py-1 mt-1 overflow-auto bg-neutral">
            <ComboboxOption v-for="postal in filteredPostals" :key="postal.code" v-slot="{ active }" :value="postal">
                <li
                    :class="[
                        'relative cursor-default select-none py-2 pl-8 pr-4',
                        active ? 'bg-primary-500 text-neutral' : 'text-gray-600',
                    ]"
                >
                    {{ postal.code }}
                </li>
            </ComboboxOption>
        </ComboboxOptions>
    </Combobox>
</template>
