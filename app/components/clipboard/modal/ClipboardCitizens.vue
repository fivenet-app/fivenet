<script lang="ts" setup>
import DataNoDataBlock from '~/components/partials/data/DataNoDataBlock.vue';
import type { ClipboardUser } from '~/stores/clipboard';
import { useClipboardStore } from '~/stores/clipboard';
import { useNotificatorStore } from '~/stores/notificator';
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

const clipboardStore = useClipboardStore();
const notifications = useNotificatorStore();

const { users, activeStack } = storeToRefs(clipboardStore);

const selected = ref<ClipboardUser[]>([]);

async function select(item: ClipboardUser): Promise<void> {
    const idx = selected.value.indexOf(item);
    if (idx > -1) {
        selected.value.splice(idx, 1);
    } else {
        // If specs are defined and max is set, clear the selection if we are over the limit
        if (props.specs && props.specs.max !== undefined && props.specs.max > 0 && selected.value.length >= props.specs.max) {
            selected.value.splice(0, props.specs.max);
        }
        selected.value.push(item);
    }

    const selectedLength = selected.value.length;
    if (props.specs !== undefined) {
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
    <div>
        <h3 v-if="!hideHeader" class="flex items-center justify-between text-lg font-medium">
            <span>{{ $t('common.citizen', 2) }}</span>
            <slot name="header" />
        </h3>

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
        />
        <table v-else class="min-w-full divide-y divide-gray-700">
            <thead>
                <tr>
                    <th v-if="showSelect" class="py-3.5 pl-4 pr-3 text-left text-sm font-semibold sm:pl-1" scope="col">
                        {{ $t('common.select', 1) }}
                    </th>
                    <th class="py-3.5 pl-4 pr-3 text-left text-sm font-semibold sm:pl-1" scope="col">
                        {{ $t('common.name') }}
                    </th>
                    <th class="px-3 py-3.5 text-left text-sm font-semibold" scope="col">
                        {{ $t('common.job', 1) }}
                    </th>
                    <th class="relative py-3.5 pl-3 pr-4 sm:pr-0" scope="col">
                        {{ $t('common.action', 2) }}
                        <UTooltip v-if="selected.length > 0" :text="$t('common.delete')">
                            <UButton variant="link" icon="i-mdi-delete" color="error" :padded="false" @click="removeAll()" />
                        </UTooltip>
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
                            {{ !selected.includes(item) ? $t('common.select', 1) : $t('common.select', 2) }}
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

                    <td class="whitespace-nowrap px-2 py-2 text-sm sm:px-4">
                        {{ item.jobLabel }}
                    </td>

                    <td class="relative whitespace-nowrap py-4 pl-3 pr-4 text-right text-sm font-medium sm:pr-0">
                        <UTooltip :text="$t('common.delete')">
                            <UButton
                                variant="link"
                                icon="i-mdi-delete"
                                color="error"
                                :padded="false"
                                @click="remove(item, true)"
                            />
                        </UTooltip>
                    </td>
                </tr>
            </tbody>
        </table>
    </div>
</template>
