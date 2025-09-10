<script lang="ts" setup>
import type { TableColumn } from '@nuxt/ui';
import { type ClipboardDocument, useClipboardStore } from '~/stores/clipboard';
import type { ObjectSpecs } from '~~/gen/ts/resources/documents/templates';
import { NotificationType } from '~~/gen/ts/resources/notifications/notifications';

const props = withDefaults(
    defineProps<{
        submit?: boolean;
        showSelect?: boolean;
        specs?: ObjectSpecs;
        hideHeader?: boolean;
    }>(),
    {
        submit: undefined,
        showSelect: true,
        specs: undefined,
        hideHeader: false,
    },
);

const emit = defineEmits<{
    (e: 'statisfied', payload: boolean): void;
    (e: 'close'): void;
}>();

const { t } = useI18n();

const clipboardStore = useClipboardStore();
const notifications = useNotificationsStore();

const { documents } = storeToRefs(clipboardStore);

const selected = ref<number[]>([]);

async function select(): Promise<void> {
    const selectedLength = selected.value.length;
    if (props.specs) {
        if (props.specs.min !== undefined && selectedLength >= props.specs.min) {
            emit('statisfied', true);
        } else if (props.specs.max !== undefined && selectedLength === props.specs.max) {
            emit('statisfied', true);
        } else {
            emit('statisfied', false);
        }
    } else {
        emit('statisfied', true);
    }
}

watch(selected, () => select());

async function remove(item: number, notify: boolean): Promise<void> {
    const idx = selected.value.indexOf(item);
    if (idx !== undefined && idx > -1) {
        selected.value.splice(idx, 1);
    }

    clipboardStore.removeDocument(item);
    if (notify) {
        notifications.add({
            title: { key: 'notifications.clipboard.document_removed.title', parameters: {} },
            description: { key: 'notifications.clipboard.document_removed.content', parameters: {} },
            duration: 3250,
            type: NotificationType.INFO,
        });
    }
}

async function removeAll(): Promise<void> {
    // Make a shallow copy to avoid mutation issues
    const toRemove = [...selected.value];
    toRemove.forEach((v) => {
        remove(v, false);
    });
    selected.value = [];

    if (props.specs !== undefined) {
        emit('statisfied', false);
    } else {
        emit('statisfied', true);
    }

    notifications.add({
        title: { key: 'notifications.clipboard.documents_removed.title', parameters: {} },
        description: { key: 'notifications.clipboard.documents_removed.content', parameters: {} },
        duration: 3250,
        type: NotificationType.INFO,
    });
}

const columns = computed(() =>
    (
        [
            props.showSelect
                ? {
                      id: 'actions',
                  }
                : undefined,
            {
                accesssorKey: 'title',
                header: t('common.title'),
                cell: ({ row }) => h('span', { class: 'text-highlighted' }, row.original.title),
            },
            {
                accesssorKey: 'creator',
                header: t('common.creator'),
                cell: ({ row }) => h('span', {}, `${row.original.creator?.firstname} ${row.original.creator?.lastname}`),
            },
            {
                id: 'delete',
            },
        ] as TableColumn<ClipboardDocument>[]
    ).filter((c) => c !== undefined),
);

watch(props, async (newVal) => {
    if (newVal.submit) {
        if (clipboardStore.activeStack) {
            clipboardStore.activeStack.documents.length = 0;
            selected.value.forEach((v) =>
                clipboardStore.activeStack.documents.push(clipboardStore.documents.find((d) => d.id === v)!),
            );
        } else if (documents.value && documents.value[0]) {
            selected.value.unshift(documents.value[0].id);
        }
    }
});
</script>

<template>
    <div>
        <h3 v-if="!hideHeader" class="flex items-center justify-between text-lg font-medium">
            <span>{{ $t('common.document', 2) }}</span>
            <slot name="header" />
        </h3>

        <UTable :columns="columns" :data="documents" :empty="$t('common.not_found', [$t('common.citizen', 2)])">
            <template #actions-cell="{ row }">
                <URadioGroup
                    v-if="specs && specs.max && specs.max === 1"
                    :model-value="selected[0]"
                    name="selected"
                    :items="[row.original.id]"
                    value-key="id"
                    :ui="{ label: 'hidden' }"
                    @update:model-value="(v) => (selected = [parseInt(v)])"
                />
                <UCheckboxGroup
                    v-else
                    :key="row.original.id"
                    v-model="selected"
                    name="selected"
                    :items="[row.original.id!]"
                    value-key="id"
                    :ui="{ label: 'hidden' }"
                />
            </template>

            <template v-if="selected.length > 0" #actions-header>
                <UTooltip :text="$t('common.delete')">
                    <UButton variant="link" icon="i-mdi-delete" color="error" size="xs" @click="removeAll()" />
                </UTooltip>
            </template>

            <template #delete-cell="{ row }">
                <UTooltip :text="$t('common.delete')">
                    <UButton variant="link" icon="i-mdi-delete" color="error" @click="remove(row.original.id!, true)" />
                </UTooltip>
            </template>
        </UTable>
    </div>
</template>
