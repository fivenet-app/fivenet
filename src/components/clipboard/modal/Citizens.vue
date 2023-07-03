<script lang="ts" setup>
import SvgIcon from '@jamescoyle/vue-icon';
import { mdiAccountMultiple, mdiTrashCan } from '@mdi/js';
import DataNoDataBlock from '~/components/partials/data/DataNoDataBlock.vue';
import { ClipboardUser, useClipboardStore } from '~/store/clipboard';
import { useNotificationsStore } from '~/store/notifications';
import { ObjectSpecs } from '~~/gen/ts/resources/documents/templates';

const clipboardStore = useClipboardStore();
const notifications = useNotificationsStore();

const { users, activeStack } = storeToRefs(clipboardStore);

const emit = defineEmits<{
    (e: 'statisfied', payload: boolean): void;
}>();

const props = withDefaults(
    defineProps<{
        submit?: boolean;
        showSelect?: boolean;
        specs?: ObjectSpecs;
    }>(),
    {
        submit: false,
        showSelect: false,
    }
);

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

    const selectedLength = BigInt(selected.value.length);
    if (props.specs) {
        if (props.specs.min && selectedLength >= props.specs.min) {
            emit('statisfied', true);
        } else if (props.specs.max && selectedLength === props.specs.max) {
            emit('statisfied', true);
        } else {
            emit('statisfied', false);
        }
    }
}

async function remove(item: ClipboardUser, notify: boolean): Promise<void> {
    const idx = selected.value.indexOf(item);
    if (idx !== undefined && idx > -1) {
        selected.value.splice(idx, 1);
    }

    clipboardStore.removeUser(item.userId!);
    if (notify) {
        notifications.dispatchNotification({
            title: { key: 'notifications.clipboard.citizen_removed.title', parameters: [] },
            content: { key: 'notifications.clipboard.citizen_removed.content', parameters: [] },
            duration: 3500,
            type: 'info',
        });
    }
}

async function removeAll(): Promise<void> {
    while (selected.value.length > 0) {
        selected.value.forEach((v) => {
            remove(v, false);
        });
    }

    emit('statisfied', false);
    notifications.dispatchNotification({
        title: { key: 'notifications.clipboard.citizens_removed.title', parameters: [] },
        content: { key: 'notifications.clipboard.citizens_removed.content', parameters: [] },
        duration: 3500,
        type: 'info',
    });
}

watch(props, async (newVal) => {
    if (newVal.submit) {
        if (activeStack.value) {
            activeStack.value.users.length = 0;
            selected.value.forEach((v) => activeStack.value.users.push(v));
        } else if (users.value && users.value.length === 1) {
            selected.value.unshift(users.value[0]);
        }
    }
});
</script>

<template>
    <h3 class="font-medium pt-1 pb-1">Users</h3>
    <DataNoDataBlock
        v-if="users?.length === 0"
        :icon="mdiAccountMultiple"
        :message="$t('components.clipboard.clipboard_modal.no_data', [$t('common.citizen', 2)])"
    />
    <table v-else class="min-w-full divide-y divide-gray-700">
        <thead>
            <tr>
                <th scope="col" class="py-3.5 pl-4 pr-3 text-left text-sm font-semibold text-white sm:pl-0" v-if="showSelect">
                    {{ $t('common.select', 1) }}
                </th>
                <th scope="col" class="py-3.5 pl-4 pr-3 text-left text-sm font-semibold text-white sm:pl-0">
                    {{ $t('common.name') }}
                </th>
                <th scope="col" class="py-3.5 px-3 text-left text-sm font-semibold text-white">
                    {{ $t('common.job', 1) }}
                </th>
                <th scope="col" class="relative py-3.5 pl-3 pr-4 sm:pr-0">
                    {{ $t('common.action', 2) }}
                    <button v-if="selected.length > 0" @click="removeAll()">
                        <SvgIcon class="w-6 h-6 mx-auto text-neutral" type="mdi" :path="mdiTrashCan" />
                    </button>
                </th>
            </tr>
        </thead>
        <tbody class="divide-y divide-gray-800">
            <tr v-for="item in users" :key="item.userId?.toString()">
                <td class="whitespace-nowrap py-4 pl-4 pr-3 text-sm font-medium text-white sm:pl-0" v-if="select">
                    <div v-if="specs && specs.max === 1n">
                        <button
                            @click="select(item)"
                            class="inline-flex w-full justify-center rounded-md bg-primary-600 px-3 py-2 text-sm font-semibold text-white shadow-sm hover:bg-primary-500 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-primary-600 sm:col-start-2"
                        >
                            <span v-if="!selected.includes(item)">
                                {{ $t('common.select', 1).toUpperCase() }}
                            </span>
                            <span v-else>
                                {{ $t('common.select', 2).toUpperCase() }}
                            </span>
                        </button>
                    </div>
                    <div v-else>
                        <input
                            name="selected"
                            :key="item.userId?.toString()"
                            :checked="selected.includes(item)"
                            :value="item"
                            v-model="selected"
                            type="checkbox"
                            class="h-4 w-4 rounded border-gray-300 text-primary-600 focus:ring-primary-600"
                        />
                    </div>
                </td>
                <td class="whitespace-nowrap py-4 pl-4 pr-3 text-sm font-medium text-white sm:pl-0">
                    {{ item.firstname }}, {{ item.lastname }}
                </td>
                <td class="whitespace-nowrap py-4 px-3 text-sm text-gray-300">
                    {{ item.jobLabel }}
                </td>
                <td class="relative whitespace-nowrap py-4 pl-3 pr-4 text-right text-sm font-medium sm:pr-0">
                    <button @click="remove(item, true)">
                        <SvgIcon class="w-6 h-6 mx-auto text-neutral" type="mdi" :path="mdiTrashCan" />
                    </button>
                </td>
            </tr>
        </tbody>
    </table>
</template>
