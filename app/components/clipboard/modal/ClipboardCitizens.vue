<script lang="ts" setup>
import { UButton, UCheckbox, UTooltip } from '#components';
import type { TableColumn, TableRow } from '@nuxt/ui';
import { h } from 'vue';
import { type ClipboardUser, useClipboardStore } from '~/stores/clipboard';
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

const { users, activeStack } = storeToRefs(clipboardStore);

const rowSelection = ref<Record<string, boolean>>(
    props.specs ? Object.fromEntries(activeStack.value.users.map((user) => [String(user.userId), true])) : {},
);
const selected = computed(() =>
    users.value.filter((user) => rowSelection.value[String(user.userId)]).map((user) => user.userId),
);
const getRowId = (user: ClipboardUser) => String(user.userId);

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
    rowSelection.value = Object.fromEntries(Object.entries(rowSelection.value).filter(([userId]) => userId !== String(item)));

    clipboardStore.removeUser(item);
    if (notify) {
        notifications.add({
            title: { key: 'notifications.clipboard.citizen_removed.title', parameters: {} },
            description: { key: 'notifications.clipboard.citizen_removed.content', parameters: {} },
            duration: 3250,
            type: NotificationType.INFO,
        });
    }
}

async function removeAll(): Promise<void> {
    // Make a shallow copy to avoid mutation issues
    const toRemove = [...selected.value];
    toRemove.forEach((v) => remove(v, false));
    rowSelection.value = {};

    emit('statisfied', !(props.specs !== undefined));

    notifications.add({
        title: { key: 'notifications.clipboard.citizens_removed.title', parameters: {} },
        description: { key: 'notifications.clipboard.citizens_removed.content', parameters: {} },
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
                accesssorKey: 'name',
                header: t('common.name'),
                cell: ({ row }) => h('span', { class: 'text-highlighted' }, userToLabel(row.original)),
            },
            {
                accesssorKey: 'job',
                header: t('common.job'),
                cell: ({ row }) => h('span', row.original.jobLabel ?? row.original.job),
            },
            {
                accesssorKey: 'dateofbirth',
                header: t('common.date_of_birth'),
                cell: ({ row }) => h('span', `${row.original.dateofbirth}`),
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
                                  : h('span', { class: 'block w-6 h-6 px-2 py-1' }),
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
                                              onClick: () => remove(row.original.userId, true),
                                          }),
                                  },
                              ),
                          ),
                  }
                : undefined,
        ] as TableColumn<ClipboardUser>[]
    ).filter((c) => c !== undefined),
);

watch(props, async (newVal) => {
    if (newVal.submit) {
        if (activeStack.value) {
            activeStack.value.users.length = 0;
            selected.value.forEach((v) => activeStack.value.users.push(clipboardStore.users.find((u) => u.userId === v)!));
        } else if (users.value && users.value[0]) {
            rowSelection.value = { [getRowId(users.value[0])]: true };
        }
    }
});

function onSelect(_event: Event, row: TableRow<ClipboardUser>): void {
    if (!props.showSelect) return;
    if (props.specs?.max === 1 && !row.getIsSelected()) rowSelection.value = {};
    row.toggleSelected(!row.getIsSelected());
}
</script>

<template>
    <div>
        <h3 v-if="!hideHeader" class="flex items-center justify-between text-lg font-medium">
            <span>{{ $t('common.citizen', 2) }}</span>
            <slot name="header" />
        </h3>

        <UTable
            v-model:row-selection="rowSelection"
            :columns="columns"
            :data="users"
            :get-row-id="getRowId"
            :row-selection-options="{ enableRowSelection: showSelect }"
            :empty="$t('common.not_found', [$t('common.citizen', 2)])"
            @select="onSelect"
        />
    </div>
</template>
