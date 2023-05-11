<script lang="ts" setup>
import { useClipboardStore } from '~/store/clipboard';
import { computed, ref, watch } from 'vue';
import { TrashIcon } from '@heroicons/vue/24/solid';
import { DocumentTextIcon } from '@heroicons/vue/20/solid';
import { ClipboardDocument } from '~/store/clipboard';
import { useNotificationsStore } from '~/store/notifications';
import { ObjectSpecs } from '@fivenet/gen/resources/documents/templates_pb';

const clipboardStore = useClipboardStore();
const notifications = useNotificationsStore();

const { t } = useI18n();

const documents = computed(() => clipboardStore.$state.documents);

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
    specs: {
        required: false,
        type: ObjectSpecs,
    },
});

const selected = ref<ClipboardDocument[]>([]);

async function select(item: ClipboardDocument): Promise<void> {
    const idx = selected.value.indexOf(item);
    if (idx !== undefined && idx > -1) {
        selected.value.splice(idx, 1);
    } else {
        if (props.specs && props.specs.getMax()) {
            selected.value.splice(0, selected.value.length);
        }
        selected.value.push(item);
    }

    if (props.specs) {
        if (selected.value.length >= props.specs.getMin()) {
            emit('statisfied', true);
        } else if (selected.value.length === props.specs.getMax()) {
            emit('statisfied', true);
        } else {
            emit('statisfied', false);
        }
    }
}

async function remove(item: ClipboardDocument, notify: boolean): Promise<void> {
    const idx = selected.value.indexOf(item);
    if (idx !== undefined && idx > -1) {
        selected.value.splice(idx, 1);
    }

    clipboardStore.removeDocument(item.id);
    if (notify) {
        notifications.dispatchNotification({
            title: t('notifications.clipboard.document_removed.title'),
            content: t('notifications.clipboard.document_removed.content'),
            duration: 3500,
            type: 'info'
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
        title: t('notifications.clipboard.documents_removed.title'),
        content: t('notifications.clipboard.documents_removed.content'),
        duration: 3500,
        type: 'info'
    });
}

watch(props, async (newVal) => {
    if (newVal.submit) {
        if (clipboardStore.activeStack) {
            clipboardStore.activeStack.documents.length = 0;
            selected.value.forEach((v) => clipboardStore.activeStack.documents.push(v));
        } else if (documents.value && documents.value.length === 1) {
            selected.value.unshift(documents.value[0]);
        }
    }
});
</script>

<template>
    <h3 class="font-medium pt-1 pb-1">Documents</h3>
    <button v-if="documents?.length == 0" type="button"
        class="relative block w-full p-4 text-center border-2 border-dashed rounded-lg border-base-300 hover:border-base-400 focus:outline-none focus:ring-2 focus:ring-neutral focus:ring-offset-2"
        disabled>
        <DocumentTextIcon class="w-12 h-12 mx-auto text-neutral" />
        <span class="block mt-2 text-sm font-semibold text-gray-300">
            {{ $t('components.clipboard.clipboard_modal.no_data', [$t('common.document', 2)]) }}
        </span>
    </button>
    <table v-else class="min-w-full divide-y divide-gray-700">
        <thead>
            <tr>
                <th scope="col" class="py-3.5 pl-4 pr-3 text-left text-sm font-semibold text-white sm:pl-0"
                    v-if="showSelect">
                    {{ $t('common.select', 1) }}
                </th>
                <th scope="col" class="py-3.5 pl-4 pr-3 text-left text-sm font-semibold text-white sm:pl-0">
                    {{ $t('common.title') }}
                </th>
                <th scope="col" class="py-3.5 px-3 text-left text-sm font-semibold text-white">
                    {{ $t('common.creator') }}
                </th>
                <th scope="col" class="relative py-3.5 pl-3 pr-4 sm:pr-0">
                    <span class="sr-only">{{ $t('common.action', 2) }}</span>
                    <button v-if="selected.length > 0" @click="removeAll()">
                        <TrashIcon class="w-6 h-6 mx-auto text-neutral" />
                    </button>
                </th>
            </tr>
        </thead>
        <tbody class="divide-y divide-gray-800">
            <tr v-for="item in documents" :key="item.id">
                <td class="whitespace-nowrap py-4 pl-4 pr-3 text-sm font-medium text-white sm:pl-0" v-if="showSelect">
                    <div v-if="specs && specs.getMax() === 1">
                        <button @click="select(item)"
                            class="inline-flex w-full justify-center rounded-md bg-indigo-600 px-3 py-2 text-sm font-semibold text-white shadow-sm hover:bg-indigo-500 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-indigo-600 sm:col-start-2">
                            <span v-if="!selected.includes(item)">
                                {{ $t('common.select', 1).toUpperCase() }}
                            </span>
                            <span v-else>
                                {{ $t('common.select', 2).toUpperCase() }}
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
                    {{ item.title }}
                </td>
                <td class="whitespace-nowrap py-4 px-3 text-sm text-gray-300">
                    {{ item.creator.firstname }}, {{ item.creator.lastname }}
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
