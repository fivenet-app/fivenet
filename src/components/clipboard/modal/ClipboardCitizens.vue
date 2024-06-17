<script lang="ts" setup>
import DataNoDataBlock from '~/components/partials/data/DataNoDataBlock.vue';
import { ClipboardUser, useClipboardStore } from '~/store/clipboard';
import { useNotificatorStore } from '~/store/notificator';
import { ObjectSpecs } from '~~/gen/ts/resources/documents/templates';
import { NotificationType } from '~~/gen/ts/resources/notifications/notifications';

const props = withDefaults(
    defineProps<{
        submit?: boolean;
        showSelect?: boolean;
        specs?: ObjectSpecs;
    }>(),
    {
        submit: undefined,
        showSelect: true,
        specs: undefined,
    },
);

const emits = defineEmits<{
    (e: 'statisfied', payload: boolean): void;
    (e: 'close'): void;
}>();

const clipboardStore = useClipboardStore();
const notifications = useNotificatorStore();

const { users, activeStack } = storeToRefs(clipboardStore);

const selected = ref<ClipboardUser[]>([]);

async function select(item: ClipboardUser): Promise<void> {
    const idx = selected.value.indexOf(item);
    if (idx !== undefined && idx > -1) {
        selected.value.splice(idx, 1);
    } else {
        if (props.specs && props.specs.max) {
            selected.value.splice(0, selected.value.length);
        }
        selected.value.push(item);
    }

    const selectedLength = selected.value.length;
    if (props.specs !== undefined) {
        if (props.specs.min !== undefined && selectedLength >= props.specs.min) {
            emits('statisfied', true);
        } else if (props.specs.max !== undefined && selectedLength === props.specs.max) {
            emits('statisfied', true);
        } else {
            emits('statisfied', false);
        }
    } else {
        emits('statisfied', true);
    }
}

async function remove(item: ClipboardUser, notify: boolean): Promise<void> {
    const idx = selected.value.indexOf(item);
    if (idx !== undefined && idx > -1) {
        selected.value.splice(idx, 1);
    }

    clipboardStore.removeUser(item.userId!);
    if (notify) {
        notifications.add({
            title: { key: 'notifications.clipboard.citizen_removed.title', parameters: {} },
            description: { key: 'notifications.clipboard.citizen_removed.content', parameters: {} },
            timeout: 3250,
            type: NotificationType.INFO,
        });
    }
}

async function removeAll(): Promise<void> {
    while (selected.value.length > 0) {
        selected.value.forEach((v) => {
            remove(v, false);
        });
    }

    if (props.specs !== undefined) {
        emits('statisfied', false);
    } else {
        emits('statisfied', true);
    }

    notifications.add({
        title: { key: 'notifications.clipboard.citizens_removed.title', parameters: {} },
        description: { key: 'notifications.clipboard.citizens_removed.content', parameters: {} },
        timeout: 3250,
        type: NotificationType.INFO,
    });
}

watch(props, async (newVal) => {
    if (newVal.submit) {
        if (activeStack.value) {
            activeStack.value.users.length = 0;
            selected.value.forEach((v) => activeStack.value.users.push(v));
        } else if (users.value && users.value[0]) {
            selected.value.unshift(users.value[0]);
        }
    }
});
</script>

<template>
    <h3 class="py-1 font-medium">{{ $t('common.citizen', 2) }}</h3>

    <DataNoDataBlock
        v-if="users?.length === 0"
        icon="i-mdi-account-multiple"
        :message="$t('components.clipboard.clipboard_modal.no_data', [$t('common.citizen', 2)])"
        :focus="
            async () => {
                navigateTo({ name: 'citizens' });
                $emit('close');
            }
        "
        padding="p-2"
    />
    <table v-else class="min-w-full divide-y divide-gray-700">
        <thead>
            <tr>
                <th v-if="showSelect" scope="col" class="py-3.5 pl-4 pr-3 text-left text-sm font-semibold sm:pl-1">
                    {{ $t('common.select', 1) }}
                </th>
                <th scope="col" class="py-3.5 pl-4 pr-3 text-left text-sm font-semibold sm:pl-1">
                    {{ $t('common.name') }}
                </th>
                <th scope="col" class="px-3 py-3.5 text-left text-sm font-semibold">
                    {{ $t('common.job', 1) }}
                </th>
                <th scope="col" class="relative py-3.5 pl-3 pr-4 sm:pr-0">
                    {{ $t('common.action', 2) }}
                    <UButton v-if="selected.length > 0" variant="link" icon="i-mdi-trash-can" @click="removeAll()" />
                </th>
            </tr>
        </thead>
        <tbody class="divide-y divide-gray-800">
            <tr v-for="item in users" :key="item.userId">
                <td v-if="showSelect" class="whitespace-nowrap py-4 pl-4 pr-3 text-sm font-medium sm:pl-1">
                    <UButton
                        v-if="specs && specs.max === 1"
                        block
                        :color="selected.includes(item) ? 'gray' : 'primary'"
                        @click="select(item)"
                    >
                        {{
                            !selected.includes(item)
                                ? $t('common.select', 1).toUpperCase()
                                : $t('common.select', 2).toUpperCase()
                        }}
                    </UButton>
                    <UCheckbox
                        v-else
                        :key="item.userId"
                        v-model="selected"
                        name="selected"
                        :checked="selected.includes(item)"
                        :value="item"
                        @click="select(item)"
                    />
                </td>
                <td class="whitespace-nowrap py-4 pl-4 pr-3 text-sm font-medium sm:pl-1">
                    {{ item.firstname }} {{ item.lastname }}
                </td>
                <td class="whitespace-nowrap px-3 py-4 text-sm">
                    {{ item.jobLabel }}
                </td>
                <td class="relative whitespace-nowrap py-4 pl-3 pr-4 text-right text-sm font-medium sm:pr-0">
                    <UButton variant="link" icon="i-mdi-trash-can" @click="remove(item, true)" />
                </td>
            </tr>
        </tbody>
    </table>
</template>
