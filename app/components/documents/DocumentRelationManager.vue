<script setup lang="ts">
import type { TabItem } from '#ui/types';
import CitizenInfoPopover from '~/components/partials/citizens/CitizenInfoPopover.vue';
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import { useAuthStore } from '~/stores/auth';
import { getUser, useClipboardStore } from '~/stores/clipboard';
import type { DocumentRelation } from '~~/gen/ts/resources/documents/documents';
import { DocRelation } from '~~/gen/ts/resources/documents/documents';
import type { User } from '~~/gen/ts/resources/users/users';

const props = defineProps<{
    open: boolean;
    documentId?: number;
    modelValue: Map<number, DocumentRelation>;
}>();

defineEmits<{
    (e: 'close'): void;
    (e: 'update:modelValue', payload: Map<number, DocumentRelation>): void;
}>();

const { $grpc } = useNuxtApp();

const { t } = useI18n();

const authStore = useAuthStore();

const clipboardStore = useClipboardStore();

const { activeChar } = storeToRefs(authStore);

const items = ref<TabItem[]>([
    {
        label: t('components.documents.document_managers.view_current'),
        icon: 'i-mdi-view-list-outline',
        slot: 'current',
    },
    {
        label: t('common.clipboard'),
        icon: 'i-mdi-clipboard-list',
        slot: 'clipboard',
    },
    {
        label: t('components.documents.document_managers.add_new'),
        icon: 'i-mdi-account-search',
        slot: 'new',
    },
]);

const queryCitizens = ref('');

const {
    data: citizens,
    pending: loading,
    refresh,
    error,
} = useLazyAsyncData(`document-${props.documentId?.toString()}-relations-citzens-${queryCitizens.value}`, () => listCitizens());

watchDebounced(queryCitizens, async () => await refresh(), {
    debounce: 200,
    maxWait: 1750,
});

async function listCitizens(): Promise<User[]> {
    try {
        const call = $grpc.citizens.citizens.listCitizens({
            pagination: {
                offset: 0,
                pageSize: 8,
            },
            search: queryCitizens.value,
        });
        const { response } = await call;

        return response.users.filter(
            (user) => !Array.from(props.modelValue.values()).find((r) => r.targetUserId === user.userId),
        );
    } catch (e) {
        handleGRPCError(e as RpcError);
        throw e;
    }
}

async function addRelation(user: User, relation: DocRelation): Promise<void> {
    const keys = Array.from(props.modelValue.keys());
    const key = !keys.length ? 1 : keys[keys.length - 1]! + 1;

    props.modelValue.set(key, {
        id: key,
        documentId: props.documentId ?? 0,
        sourceUserId: activeChar.value!.userId,
        sourceUser: activeChar.value!,
        targetUserId: user.userId,
        targetUser: user,
        relation,
    });

    await refresh();
}

async function removeRelation(id: number): Promise<void> {
    props.modelValue.delete(id);
    refresh();
}

const columnsCurrent = [
    {
        key: 'name',
        label: t('common.name'),
    },
    {
        key: 'creator',
        label: t('common.creator'),
    },
    {
        key: 'relation',
        label: t('common.relation', 1),
    },
    {
        key: 'actions',
        label: t('common.action', 2),
        sortable: false,
    },
];

const columnsClipboard = [
    {
        key: 'name',
        label: t('common.name'),
    },
    {
        key: 'job',
        label: t('common.job'),
    },
    {
        key: 'relations',
        label: t('components.documents.document_managers.add_relation'),
        sortable: false,
    },
];

const columnsNew = [
    {
        key: 'name',
        label: t('common.name'),
    },
    {
        key: 'job',
        label: t('common.job'),
    },
    {
        key: 'relations',
        label: t('components.documents.document_managers.add_relation'),
        sortable: false,
    },
];
</script>

<template>
    <UModal :ui="{ width: 'w-full sm:max-w-5xl' }" :model-value="open" @update:model-value="$emit('close')">
        <UCard :ui="{ ring: '', divide: 'divide-y divide-gray-100 dark:divide-gray-800' }">
            <template #header>
                <div class="flex items-center justify-between">
                    <h3 class="text-2xl font-semibold leading-6">
                        {{ $t('common.citizen', 1) }}
                        {{ $t('common.relation', 2) }}
                    </h3>

                    <UButton class="-my-1" color="gray" variant="ghost" icon="i-mdi-window-close" @click="$emit('close')" />
                </div>
            </template>

            <div>
                <UTabs :items="items">
                    <template #current>
                        <div>
                            <UTable
                                :columns="columnsCurrent"
                                :rows="[...modelValue.values()]"
                                :empty-state="{ icon: 'i-mdi-file', label: $t('common.not_found', [$t('common.relation', 2)]) }"
                            >
                                <template #name-data="{ row }">
                                    <CitizenInfoPopover :user="row.targetUser" show-birthdate />
                                </template>

                                <template #creator-data="{ row }">
                                    <CitizenInfoPopover :user="row.sourceUser" :trailing="false" />
                                </template>

                                <template #relation-data="{ row }">
                                    {{ $t(`enums.documents.DocRelation.${DocRelation[row.relation]}`) }}
                                </template>

                                <template #actions-data="{ row }">
                                    <UButtonGroup>
                                        <UTooltip :text="$t('components.documents.document_managers.open_citizen')">
                                            <UButton
                                                :to="{
                                                    name: 'citizens-id',
                                                    params: {
                                                        id: row.targetUserId,
                                                    },
                                                }"
                                                target="_blank"
                                                variant="link"
                                                icon="i-mdi-open-in-new"
                                            />
                                        </UTooltip>

                                        <UTooltip :text="$t('components.documents.document_managers.remove_relation')">
                                            <UButton
                                                variant="link"
                                                icon="i-mdi-account-minus"
                                                color="error"
                                                @click="removeRelation(row.id!)"
                                            />
                                        </UTooltip>
                                    </UButtonGroup>
                                </template>
                            </UTable>
                        </div>
                    </template>

                    <template #clipboard>
                        <div>
                            <UTable
                                :columns="columnsClipboard"
                                :rows="clipboardStore.$state.users"
                                :empty-state="{ icon: 'i-mdi-file', label: $t('common.not_found', [$t('common.citizen', 2)]) }"
                            >
                                <template #name-data="{ row }">
                                    <CitizenInfoPopover :user="row.targetUser" show-birthdate />
                                </template>

                                <template #job-data="{ row }">
                                    {{ row.jobLabel }}
                                </template>

                                <template #relations-data="{ row }">
                                    <UButtonGroup>
                                        <UTooltip :text="$t('components.documents.document_managers.mentioned')">
                                            <UButton
                                                color="blue"
                                                icon="i-mdi-at"
                                                @click="addRelation(getUser(row), DocRelation.MENTIONED)"
                                            />
                                        </UTooltip>

                                        <UTooltip :text="$t('components.documents.document_managers.targets')">
                                            <UButton
                                                color="amber"
                                                icon="i-mdi-target"
                                                @click="addRelation(getUser(row), DocRelation.TARGETS)"
                                            />
                                        </UTooltip>

                                        <UTooltip :text="$t('components.documents.document_managers.caused')">
                                            <UButton
                                                color="error"
                                                icon="i-mdi-source-commit-start"
                                                @click="addRelation(getUser(row), DocRelation.CAUSED)"
                                            />
                                        </UTooltip>
                                    </UButtonGroup>
                                </template>
                            </UTable>
                        </div>
                    </template>

                    <template #new>
                        <UFormGroup class="mb-2" name="name" :label="$t('common.search')">
                            <UInput
                                v-model="queryCitizens"
                                type="text"
                                name="name"
                                :placeholder="`${$t('common.citizen', 1)} ${$t('common.name')}`"
                                leading-icon="i-mdi-search"
                            />
                        </UFormGroup>

                        <div>
                            <DataErrorBlock
                                v-if="error"
                                :title="$t('common.unable_to_load', [$t('common.citizen', 2)])"
                                :error="error"
                                :retry="refresh"
                            />

                            <UTable
                                v-else
                                :columns="columnsNew"
                                :loading="loading"
                                :rows="citizens"
                                :empty-state="{ icon: 'i-mdi-file', label: $t('common.not_found', [$t('common.citizen', 2)]) }"
                            >
                                <template #name-data="{ row }">
                                    <CitizenInfoPopover :user="row" show-birthdate />
                                </template>

                                <template #job-data="{ row }">
                                    {{ row.jobLabel }}
                                </template>

                                <template #relations-data="{ row }">
                                    <UButtonGroup>
                                        <UTooltip :text="$t('components.documents.document_managers.mentioned')">
                                            <UButton
                                                color="blue"
                                                icon="i-mdi-at"
                                                @click="addRelation(row, DocRelation.MENTIONED)"
                                            />
                                        </UTooltip>

                                        <UTooltip :text="$t('components.documents.document_managers.targets')">
                                            <UButton
                                                color="amber"
                                                icon="i-mdi-target"
                                                @click="addRelation(row, DocRelation.TARGETS)"
                                            />
                                        </UTooltip>

                                        <UTooltip :text="$t('components.documents.document_managers.caused')">
                                            <UButton
                                                color="error"
                                                icon="i-mdi-source-commit-start"
                                                @click="addRelation(row, DocRelation.CAUSED)"
                                            />
                                        </UTooltip>
                                    </UButtonGroup>
                                </template>
                            </UTable>
                        </div>
                    </template>
                </UTabs>
            </div>

            <template #footer>
                <UButton class="flex-1" block color="black" @click="$emit('close')">
                    {{ $t('common.close', 1) }}
                </UButton>
            </template>
        </UCard>
    </UModal>
</template>
