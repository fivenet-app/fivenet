<script lang="ts" setup>
import type { Postal } from '~/types/livemap';
import { NotificationType } from '~~/gen/ts/resources/notifications/notifications';

defineOptions({
    inheritAttrs: false,
});

const props = withDefaults(
    defineProps<{
        modelValue?: Postal;
        selectedCode?: string;
        disabled?: boolean;
    }>(),
    {
        modelValue: undefined,
        selectedCode: '',
        disabled: false,
    },
);

const emit = defineEmits<{
    (e: 'update:modelValue', value: Postal | undefined): void;
}>();

const notifications = useNotificationsStore();

const selectedPostal = ref<Postal | undefined>(props.modelValue);
const postalQuery = ref('');

const {
    data: postals,
    status,
    execute,
} = useLazyAsyncData(
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

watch(
    () => props.modelValue,
    (value) => {
        if ((value?.code ?? '') === (selectedPostal.value?.code ?? '')) return;
        selectedPostal.value = value;
    },
    { immediate: true },
);

watch(
    () => props.selectedCode,
    async (code) => {
        if (!code) return;
        if (!postals.value) await execute();

        if (!postals.value) return;

        const match = postals.value.find((postal) => postal.code === code);
        if (match) {
            if ((props.modelValue?.code ?? '') === match.code) return;
            selectedPostal.value = match;
        }
    },
    { immediate: true },
);

watch(selectedPostal, () => {
    if ((selectedPostal.value?.code ?? '') === (props.modelValue?.code ?? '')) return;
    emit('update:modelValue', selectedPostal.value);
});

watchOnce(postalQuery, () => onOpen());
</script>

<template>
    <ClientOnly>
        <UInputMenu
            v-model="selectedPostal"
            v-model:search-term.trim="postalQuery"
            class="w-full"
            :items="postals ?? []"
            label-key="code"
            nullable
            :loading="status === 'pending'"
            :placeholder="`${$t('common.postal')} ${$t('common.search')}`"
            :search-input="{ placeholder: $t('common.search_field') }"
            :disabled="disabled"
            leading-icon="i-mdi-postage-stamp"
            virtualize
            v-bind="$attrs"
            @update:open="onOpen"
        >
            <template #empty> {{ $t('common.not_found', [$t('common.postal', 2)]) }} </template>
        </UInputMenu>
    </ClientOnly>
</template>
