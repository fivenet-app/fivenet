<script lang="ts" setup>
import DataNoDataBlock from '~/components/partials/data/DataNoDataBlock.vue';
import { ClipboardDocument, useClipboardStore } from '~/store/clipboard';
import { useNotificatorStore } from '~/store/notificator';
import { ObjectSpecs } from '~~/gen/ts/resources/documents/templates';
import { NotificationType } from '~~/gen/ts/resources/notifications/notifications';

const clipboardStore = useClipboardStore();
const notifications = useNotificatorStore();

const { documents } = storeToRefs(clipboardStore);

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
        showSelect: true,
        specs: undefined,
    },
);

const selected = ref<ClipboardDocument[]>([]);

async function select(item: ClipboardDocument): Promise<void> {
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

async function remove(item: ClipboardDocument, notify: boolean): Promise<void> {
    const idx = selected.value.indexOf(item);
    if (idx !== undefined && idx > -1) {
        selected.value.splice(idx, 1);
    }

    clipboardStore.removeDocument(item.id);
    if (notify) {
        notifications.add({
            title: { key: 'notifications.clipboard.document_removed.title', parameters: {} },
            description: { key: 'notifications.clipboard.document_removed.content', parameters: {} },
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
        emit('statisfied', false);
    } else {
        emit('statisfied', true);
    }

    notifications.add({
        title: { key: 'notifications.clipboard.documents_removed.title', parameters: {} },
        description: { key: 'notifications.clipboard.documents_removed.content', parameters: {} },
        timeout: 3250,
        type: NotificationType.INFO,
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
    <h3 class="py-1 font-medium">{{ $t('common.document', 2) }}</h3>
    <DataNoDataBlock
        v-if="documents?.length === 0"
        icon="i-mdi-file-document-multiple"
        :message="$t('components.clipboard.clipboard_modal.no_data', [$t('common.document', 2)])"
    />
    <table v-else class="min-w-full divide-y divide-gray-700">
        <thead>
            <tr>
                <th v-if="showSelect" scope="col" class="py-3.5 pl-4 pr-3 text-left text-sm font-semibold sm:pl-1">
                    {{ $t('common.select', 1) }}
                </th>
                <th scope="col" class="py-3.5 pl-4 pr-3 text-left text-sm font-semibold sm:pl-1">
                    {{ $t('common.title') }}
                </th>
                <th scope="col" class="px-3 py-3.5 text-left text-sm font-semibold">
                    {{ $t('common.creator') }}
                </th>
                <th scope="col" class="relative py-3.5 pl-3 pr-4 sm:pr-0">
                    <span class="sr-only">{{ $t('common.action', 2) }}</span>
                    <UButton v-if="selected.length > 0" variant="link" icon="i-mdi-trash-can" @click="removeAll()" />
                </th>
            </tr>
        </thead>
        <tbody class="divide-y divide-gray-800">
            <tr v-for="item in documents" :key="item.id">
                <td v-if="showSelect" class="whitespace-nowrap py-4 pl-4 pr-3 text-sm font-medium sm:pl-1">
                    <template v-if="specs && specs.max === 1">
                        <UButton block :color="selected.includes(item) ? 'gray' : 'primary'" @click="select(item)">
                            {{
                                !selected.includes(item)
                                    ? $t('common.select', 1).toUpperCase()
                                    : $t('common.select', 2).toUpperCase()
                            }}
                        </UButton>
                    </template>
                    <template v-else>
                        <UCheckbox
                            :key="item.id"
                            v-model="selected"
                            name="selected"
                            :checked="selected.includes(item)"
                            :value="item"
                            @click="select(item)"
                        />
                    </template>
                </td>
                <td class="whitespace-nowrap py-4 pl-4 pr-3 text-sm font-medium sm:pl-1">
                    {{ item.title }}
                </td>
                <td class="whitespace-nowrap px-3 py-4 text-sm">{{ item.creator.firstname }} {{ item.creator.lastname }}</td>
                <td class="relative whitespace-nowrap py-4 pl-3 pr-4 text-right text-sm font-medium sm:pr-0">
                    <UButton variant="link" icon="i-mdi-trash-can" @click="remove(item, true)" />
                </td>
            </tr>
        </tbody>
    </table>
</template>
