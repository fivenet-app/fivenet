<script lang="ts" setup>
import { UButton, UCheckbox, UTooltip } from '#components';
import type { TableColumn, TableRow } from '@nuxt/ui';
import { h } from 'vue';
import { type ClipboardDocument, useClipboardStore } from '~/stores/clipboard';
import type { ObjectSpecs } from '~~/gen/ts/resources/documents/templates/templates';
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

const { documents, activeStack } = storeToRefs(clipboardStore);

const rowSelection = ref<Record<string, boolean>>(
    props.specs ? Object.fromEntries(activeStack.value.documents.map((document) => [String(document.id), true])) : {},
);
const selected = computed(() =>
    documents.value.filter((document) => rowSelection.value[String(document.id)]).map((document) => document.id),
);
const getRowId = (document: ClipboardDocument) => String(document.id);

async function select(): Promise<void> {
    const selectedLength = selected.value.length;
    if (!props.specs) {
        emit('statisfied', true);
        return;
    }

    emit(
        'statisfied',
        (!props.specs.required || selectedLength > 0) &&
            selectedLength >= (props.specs.min ?? 0) &&
            (props.specs.max === undefined || props.specs.max <= 0 || selectedLength <= props.specs.max),
    );
}

watch(selected, () => select());

async function remove(item: number, notify: boolean): Promise<void> {
    rowSelection.value = Object.fromEntries(
        Object.entries(rowSelection.value).filter(([documentId]) => documentId !== String(item)),
    );

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
    rowSelection.value = {};

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
                      id: 'select',
                      header: ({ table }) =>
                          props.specs?.max === 1
                              ? h('span', { class: 'block h-8' })
                              : h(UCheckbox, {
                                    modelValue: table.getIsSomePageRowsSelected()
                                        ? 'indeterminate'
                                        : table.getIsAllPageRowsSelected(),
                                    'onUpdate:modelValue': (value: unknown) => table.toggleAllPageRowsSelected(!!value),
                                }),
                      cell: ({ row }) =>
                          h(UCheckbox, {
                              modelValue: row.getIsSelected(),
                              ui: { label: 'hidden' },
                              'onUpdate:modelValue': (value: unknown) => row.toggleSelected(!!value),
                          }),
                      meta: {
                          class: {
                              td: 'pe-0 px-4 py-1.5',
                          },
                      },
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
            !props.specs
                ? {
                      id: 'delete',
                      header: () =>
                          h(
                              'div',
                              { class: 'flex h-6 items-center justify-center' },
                              selected.value.length > 0
                                  ? h(
                                        UTooltip,
                                        { text: t('common.delete') },
                                        {
                                            default: () =>
                                                h(UButton, {
                                                    variant: 'link',
                                                    icon: 'i-mdi-delete',
                                                    color: 'error',
                                                    size: 'xs',
                                                    onClick: removeAll,
                                                }),
                                        },
                                    )
                                  : undefined,
                          ),
                      cell: ({ row }) =>
                          h(
                              'div',
                              { class: 'flex h-6 items-center justify-center' },
                              h(
                                  UTooltip,
                                  { text: t('common.delete') },
                                  {
                                      default: () =>
                                          h(UButton, {
                                              variant: 'link',
                                              icon: 'i-mdi-delete',
                                              color: 'error',
                                              size: 'xs',
                                              onClick: () => remove(row.original.id, true),
                                          }),
                                  },
                              ),
                          ),
                  }
                : undefined,
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
            rowSelection.value = { [getRowId(documents.value[0])]: true };
        }
    }
});

function onSelect(_event: Event, row: TableRow<ClipboardDocument>): void {
    if (!props.showSelect) return;
    if (props.specs?.max === 1 && !row.getIsSelected()) rowSelection.value = {};
    row.toggleSelected(!row.getIsSelected());
}
</script>

<template>
    <div>
        <h3 v-if="!hideHeader" class="flex items-center justify-between text-lg font-medium">
            <span>{{ $t('common.document', 2) }}</span>
            <slot name="header" />
        </h3>

        <UTable
            v-model:row-selection="rowSelection"
            :columns="columns"
            :data="documents"
            :get-row-id="getRowId"
            :row-selection-options="{ enableRowSelection: showSelect }"
            :empty="$t('common.not_found', [$t('common.citizen', 2)])"
            @select="onSelect"
        />
    </div>
</template>
