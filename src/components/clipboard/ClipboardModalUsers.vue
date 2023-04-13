<script lang="ts" setup>
import { useClipboardStore, ClipboardUser } from '~/store/clipboard';
import { computed, ref, watch } from 'vue';
import { TrashIcon } from '@heroicons/vue/24/solid';
import { UsersIcon } from '@heroicons/vue/20/solid';
import { dispatchNotification } from '~/components/partials/notification';

const store = useClipboardStore();

const users = computed(() => store.$state.users);

const emit = defineEmits<{
    (e: 'statisfied', payload: boolean): void,
}>();

const props = defineProps({
    submit: {
        required: false,
        type: Boolean,
        default: false,
    },
    showSelect: {
        required: false,
        type: Boolean,
        default: false,
    },
    min: {
        required: false,
        type: Number,
        default: 0,
    },
    max: {
        required: false,
        type: Number,
        default: 0,
    },
});

const selected = ref<ClipboardUser[]>([]);

async function select(item: ClipboardUser): Promise<void> {
    const idx = selected.value.indexOf(item);
    if (idx !== undefined && idx > -1) {
        selected.value.splice(idx, 1);
    } else {
        if (props.max) {
            selected.value.splice(0, selected.value.length);
        }
        selected.value.push(item);
    }

    if (selected.value.length >= props.min) {
        emit('statisfied', true);
    } else if (selected.value.length === props.max) {
        emit('statisfied', true);
    } else {
        emit('statisfied', false);
    }
}

async function remove(item: ClipboardUser, notify: boolean): Promise<void> {
    const idx = selected.value.indexOf(item);
    if (idx !== undefined && idx > -1) {
        selected.value.splice(idx, 1);
    }

    await store.removeUser(item.id!);
    if (notify) {
        dispatchNotification({ title: 'Clipboard: Citizen removed', content: 'Selected citizen removed from clipboard', duration: 3500 });
    }
}

async function removeAll(): Promise<void> {
    while (selected.value.length > 0) {
        selected.value.forEach((v) => {
            remove(v, false);
        });
    }

    emit('statisfied', false);
    dispatchNotification({ title: 'Clipboard: Citizens removed', content: 'All citizens have been removed from your clipboard', duration: 3500 });
}

watch(props, async (newVal) => {
    if (newVal.submit) {
        if (store.activeStack) {
            store.activeStack.users.length = 0;
            selected.value.forEach((v) => store.activeStack.users.push(v));
        } else if (users.value && users.value.length === 1) {
            selected.value.unshift(users.value[0]);
        }
    }
});
</script>

<template>
    <h3 class="font-medium pt-1 pb-1">Users</h3>
    <button v-if="users?.length == 0" type="button"
        class="relative block w-full p-4 text-center border-2 border-dashed rounded-lg border-base-300 hover:border-base-400 focus:outline-none focus:ring-2 focus:ring-neutral focus:ring-offset-2"
        disabled>
        <UsersIcon class="w-12 h-12 mx-auto text-neutral" />
        <span class="block mt-2 text-sm font-semibold text-gray-300">
            No Users in Clipboard.
        </span>
    </button>
    <table v-else class="min-w-full divide-y divide-gray-700">
        <thead>
            <tr>
                <th scope="col" class="py-3.5 pl-4 pr-3 text-left text-sm font-semibold text-white sm:pl-0"
                    v-if="showSelect">
                    Select
                </th>
                <th scope="col" class="py-3.5 pl-4 pr-3 text-left text-sm font-semibold text-white sm:pl-0">
                    Name
                </th>
                <th scope="col" class="py-3.5 px-3 text-left text-sm font-semibold text-white">
                    Job
                </th>
                <th scope="col" class="relative py-3.5 pl-3 pr-4 sm:pr-0">
                    Actions
                    <button v-if="selected.length > 0" @click="removeAll()">
                        <TrashIcon class="w-6 h-6 mx-auto text-neutral" />
                    </button>
                </th>
            </tr>
        </thead>
        <tbody class="divide-y divide-gray-800">
            <tr v-for="item in users" :key="item.id">
                <td class="whitespace-nowrap py-4 pl-4 pr-3 text-sm font-medium text-white sm:pl-0" v-if="select">
                    <div v-if="max === 1">
                        <button @click="select(item)"
                            class="inline-flex w-full justify-center rounded-md bg-indigo-600 px-3 py-2 text-sm font-semibold text-white shadow-sm hover:bg-indigo-500 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-indigo-600 sm:col-start-2">
                            <span v-if="!selected.includes(item)">
                                SELECT
                            </span>
                            <span v-else>
                                SELECTED
                            </span>
                        </button>
                    </div>
                    <div v-else>
                        <input id="selected" name="selected" :key="item.id" :checked="selected.includes(item)" :value="item"
                            v-model="selected" type="checkbox"
                            class="h-4 w-4 rounded border-gray-300 text-indigo-600 focus:ring-indigo-600" />
                    </div>
                </td>
                <td class="whitespace-nowrap py-4 pl-4 pr-3 text-sm font-medium text-white sm:pl-0">
                    {{ item.firstname }}, {{ item.lastname }}
                </td>
                <td class="whitespace-nowrap py-4 px-3 text-sm text-gray-300">
                    {{ item.jobLabel }} (Rank: {{ item.jobGradeLabel }})
                </td>
                <td class="relative whitespace-nowrap py-4 pl-3 pr-4 text-right text-sm font-medium sm:pr-0">
                    <button @click="remove(item, true)">
                        <TrashIcon class="w-6 h-6 mx-auto text-neutral" />
                    </button>
                </td>
            </tr>
        </tbody>
    </table>
</template>
