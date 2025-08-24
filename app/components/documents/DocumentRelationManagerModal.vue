<script setup lang="ts">
import CitizenInfoPopover from '~/components/partials/citizens/CitizenInfoPopover.vue';
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import { useAuthStore } from '~/stores/auth';
import { getUser, useClipboardStore } from '~/stores/clipboard';
import { getCitizensCitizensClient } from '~~/gen/ts/clients';
import { type DocumentRelation, DocRelation } from '~~/gen/ts/resources/documents/documents';
import type { User } from '~~/gen/ts/resources/users/users';

const props = defineProps<{
    documentId?: number;
}>();

const modelValue = defineModel<DocumentRelation[]>('relations', {
    type: Array,
    required: true,
});

const { isOpen } = useOverlay();

const { t } = useI18n();

const authStore = useAuthStore();
const { activeChar } = storeToRefs(authStore);

const clipboardStore = useClipboardStore();

const citizensCitizensClient = await getCitizensCitizensClient();

const items = ref<TabsItem[]>([
    {
        label: t('components.documents.document_managers.view_current'),
        icon: 'i-mdi-view-list-outline',
        slot: 'current' as const,
        value: 'current',
    },
    {
        label: t('common.clipboard'),
        icon: 'i-mdi-clipboard-list',
        slot: 'clipboard' as const,
        value: 'clipboard',
    },
    {
        label: t('components.documents.document_managers.add_new'),
        icon: 'i-mdi-account-search',
        slot: 'new' as const,
        value: 'new',
    },
]);

const queryCitizens = ref('');

const {
    data: citizens,
    status,
    refresh,
    error,
} = useLazyAsyncData(`document-${props.documentId?.toString()}-relations-citzens-${queryCitizens.value}`, () => listCitizens());

watchDebounced(queryCitizens, async () => await refresh(), {
    debounce: 200,
    maxWait: 1750,
});

async function listCitizens(): Promise<User[]> {
    try {
        const call = citizensCitizensClient.listCitizens({
            pagination: {
                offset: 0,
                pageSize: 8,
            },
            search: queryCitizens.value,
        });
        const { response } = await call;

        return response.users.filter((user) => !modelValue.value.find((r) => r.targetUserId === user.userId));
    } catch (e) {
        handleGRPCError(e as RpcError);
        throw e;
    }
}

let lastId = 0;

async function addRelation(user: User, relation: DocRelation): Promise<void> {
    modelValue.value.push({
        id: lastId--,
        documentId: props.documentId ?? 0,
        sourceUserId: activeChar.value!.userId,
        sourceUser: activeChar.value!,
        targetUserId: user.userId,
        targetUser: user,
        relation: relation,
    });

    await refresh();
}

async function removeRelation(id: number): Promise<void> {
    const idx = modelValue.value.findIndex((r) => r.id === id);
    if (idx > -1) {
        modelValue.value.splice(idx, 1);
    }
    refresh();
}

const columnsCurrent = [
    {
        accessorKey: 'name',
        label: t('common.name'),
    },
    {
        accessorKey: 'creator',
        label: t('common.creator'),
    },
    {
        accessorKey: 'relation',
        label: t('common.relation', 1),
    },
    {
        accessorKey: 'actions',
        label: t('common.action', 2),
        sortable: false,
    },
];

const columnsClipboard = [
    {
        accessorKey: 'name',
        label: t('common.name'),
    },
    {
        accessorKey: 'job',
        label: t('common.job'),
    },
    {
        accessorKey: 'relations',
        label: t('components.documents.document_managers.add_relation'),
        sortable: false,
    },
];

const columnsNew = [
    {
        accessorKey: 'name',
        label: t('common.name'),
    },
    {
        accessorKey: 'job',
        label: t('common.job'),
    },
    {
        accessorKey: 'relations',
        label: t('components.documents.document_managers.add_relation'),
        sortable: false,
    },
];
</script>

<template>
    <UModal :ui="{ width: 'w-full sm:max-w-5xl' }">
        <UCard>
            <template #header>
                <div class="flex items-center justify-between">
                    <h3 class="text-2xl leading-6 font-semibold">
                        {{ $t('common.citizen', 1) }}
                        {{ $t('common.relation', 2) }}
                    </h3>

                    <UButton class="-my-1" color="neutral" variant="ghost" icon="i-mdi-window-close" @click="isOpen = false" />
                </div>
            </template>

            <div>
                <UTabs :items="items">
                    <template #current>
                        <div>
                            <UTable
                                :columns="columnsCurrent"
                                :data="modelValue"
                                :empty-state="{ icon: 'i-mdi-file', label: $t('common.not_found', [$t('common.relation', 2)]) }"
                            >
                                <template #name-cell="{ row }">
                                    <CitizenInfoPopover
                                        :user="!row.targetUser?.userId ? undefined : row.targetUser"
                                        :user-id="row.targetUserId"
                                        show-birthdate
                                    />
                                </template>

                                <template #creator-cell="{ row }">
                                    <CitizenInfoPopover
                                        :user="!row.sourceUser?.userId ? undefined : row.sourceUser"
                                        :user-id="row.sourceUserId"
                                        :trailing="false"
                                    />
                                </template>

                                <template #relation-cell="{ row }">
                                    {{ $t(`enums.documents.DocRelation.${DocRelation[row.relation]}`) }}
                                </template>

                                <template #actions-cell="{ row }">
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
                                :data="clipboardStore.$state.users"
                                :empty-state="{ icon: 'i-mdi-file', label: $t('common.not_found', [$t('common.citizen', 2)]) }"
                            >
                                <template #name-cell="{ row }">
                                    <CitizenInfoPopover :user="row.targetUser" show-birthdate />
                                </template>

                                <template #job-cell="{ row }">
                                    {{ row.jobLabel }}
                                </template>

                                <template #relations-cell="{ row }">
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
                                                color="warning"
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
                        <UFormField class="mb-2" name="name" :label="$t('common.search')">
                            <UInput
                                v-model="queryCitizens"
                                type="text"
                                name="name"
                                :placeholder="`${$t('common.citizen', 1)} ${$t('common.name')}`"
                                leading-icon="i-mdi-search"
                            />
                        </UFormField>

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
                                :loading="isRequestPending(status)"
                                :data="citizens"
                                :empty-state="{ icon: 'i-mdi-file', label: $t('common.not_found', [$t('common.citizen', 2)]) }"
                            >
                                <template #name-cell="{ row }">
                                    <CitizenInfoPopover :user="row" show-birthdate />
                                </template>

                                <template #job-cell="{ row }">
                                    {{ row.jobLabel }}
                                </template>

                                <template #relations-cell="{ row }">
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
                                                color="warning"
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
                <UButton class="flex-1" block color="neutral" @click="isOpen = false">
                    {{ $t('common.close', 1) }}
                </UButton>
            </template>
        </UCard>
    </UModal>
</template>
