<script setup lang="ts">
import CitizenInfoPopover from '~/components/partials/citizens/CitizenInfoPopover.vue';
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import DataNoDataBlock from '~/components/partials/data/DataNoDataBlock.vue';
import DataPendingBlock from '~/components/partials/data/DataPendingBlock.vue';
import { useAuthStore } from '~/store/auth';
import { getUser, useClipboardStore } from '~/store/clipboard';
import type { DocumentRelation } from '~~/gen/ts/resources/documents/documents';
import { DocRelation } from '~~/gen/ts/resources/documents/documents';
import type { User } from '~~/gen/ts/resources/users/users';

const props = defineProps<{
    open: boolean;
    documentId?: string;
    modelValue: Map<string, DocumentRelation>;
}>();

defineEmits<{
    (e: 'close'): void;
    (e: 'update:modelValue', payload: Map<string, DocumentRelation>): void;
}>();

const { t } = useI18n();

const authStore = useAuthStore();

const clipboardStore = useClipboardStore();

const { activeChar } = storeToRefs(authStore);

const tabs = ref<{ key: string; label: string; icon: string }[]>([
    {
        key: 'current',
        label: t('components.documents.document_managers.view_current'),
        icon: 'i-mdi-view-list-outline',
    },
    {
        key: 'clipboard',
        label: t('common.clipboard'),
        icon: 'i-mdi-clipboard-list',
    },
    {
        key: 'new',
        label: t('components.documents.document_managers.add_new'),
        icon: 'i-mdi-account-search',
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
        const call = getGRPCCitizenStoreClient().listCitizens({
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

function addRelation(user: User, relation: DocRelation): void {
    const keys = Array.from(props.modelValue.keys());
    const key = !keys.length ? '1' : (parseInt(keys[keys.length - 1]!) + 1).toString();

    props.modelValue.set(key, {
        id: key,
        documentId: props.documentId ?? '0',
        sourceUserId: activeChar.value!.userId,
        sourceUser: activeChar.value!,
        targetUserId: user.userId,
        targetUser: user,
        relation,
    });
    refresh();
}

function removeRelation(id: string): void {
    props.modelValue.delete(id);
    refresh();
}
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

                    <UButton color="gray" variant="ghost" icon="i-mdi-window-close" class="-my-1" @click="$emit('close')" />
                </div>
            </template>

            <div>
                <UTabs :items="tabs">
                    <template #item="{ item }">
                        <template v-if="item.key === 'current'">
                            <div class="flow-root">
                                <div class="-my-2 mx-0 overflow-x-auto">
                                    <div class="inline-block min-w-full py-2 align-middle">
                                        <table class="min-w-full divide-y divide-base-200">
                                            <thead>
                                                <tr>
                                                    <th
                                                        scope="col"
                                                        class="py-3.5 pl-4 pr-3 text-left text-sm font-semibold sm:pl-6 lg:pl-8"
                                                    >
                                                        {{ $t('common.name') }}
                                                    </th>
                                                    <th scope="col" class="px-3 py-3.5 text-left text-sm font-semibold">
                                                        {{ $t('common.creator') }}
                                                    </th>
                                                    <th scope="col" class="px-3 py-3.5 text-left text-sm font-semibold">
                                                        {{ $t('common.relation', 1) }}
                                                    </th>
                                                    <th scope="col" class="px-3 py-3.5 text-left text-sm font-semibold">
                                                        {{ $t('common.action', 2) }}
                                                    </th>
                                                </tr>
                                            </thead>
                                            <tbody class="divide-y divide-base-500">
                                                <tr v-for="[key, relation] in modelValue" :key="key.toString()">
                                                    <td
                                                        class="truncate whitespace-nowrap py-4 pl-4 pr-3 text-sm font-medium sm:pl-6 lg:pl-8"
                                                    >
                                                        <span class="inline-flex items-center gap-1">
                                                            <CitizenInfoPopover :user="relation.targetUser" :trailing="false" />
                                                            ({{ relation.targetUser?.dateofbirth }})
                                                        </span>
                                                    </td>
                                                    <td class="whitespace-nowrap px-2 py-2 text-sm sm:px-4">
                                                        <CitizenInfoPopover :user="relation.sourceUser" :trailing="false" />
                                                    </td>
                                                    <td class="whitespace-nowrap px-2 py-2 text-sm sm:px-4">
                                                        {{ $t(`enums.docstore.DocRelation.${DocRelation[relation.relation]}`) }}
                                                    </td>
                                                    <td class="whitespace-nowrap px-2 py-2 text-sm sm:px-4">
                                                        <div class="flex flex-row gap-2">
                                                            <UButton
                                                                :to="{
                                                                    name: 'citizens-id',
                                                                    params: {
                                                                        id: relation.targetUserId,
                                                                    },
                                                                }"
                                                                target="_blank"
                                                                :title="
                                                                    $t('components.documents.document_managers.open_citizen')
                                                                "
                                                                variant="link"
                                                                icon="i-mdi-open-in-new"
                                                            />
                                                            <UButton
                                                                :title="
                                                                    $t('components.documents.document_managers.remove_relation')
                                                                "
                                                                variant="link"
                                                                icon="i-mdi-account-minus"
                                                                color="red"
                                                                @click="removeRelation(relation.id!)"
                                                            />
                                                        </div>
                                                    </td>
                                                </tr>
                                            </tbody>
                                        </table>
                                    </div>
                                </div>
                            </div>
                        </template>

                        <template v-else-if="item.key === 'clipboard'">
                            <div class="mt-2 flow-root">
                                <div class="-my-2 mx-0 overflow-x-auto">
                                    <div class="inline-block min-w-full py-2 align-middle">
                                        <DataNoDataBlock
                                            v-if="clipboardStore.$state.users.length === 0"
                                            :type="$t('common.citizen', 2)"
                                            icon="i-mdi-account-multiple"
                                        />
                                        <table v-else class="min-w-full divide-y divide-base-200">
                                            <thead>
                                                <tr>
                                                    <th
                                                        scope="col"
                                                        class="py-3.5 pl-4 pr-3 text-left text-sm font-semibold sm:pl-6 lg:pl-8"
                                                    >
                                                        {{ $t('common.name') }}
                                                    </th>
                                                    <th scope="col" class="px-3 py-3.5 text-left text-sm font-semibold">
                                                        {{ $t('common.job', 1) }}
                                                    </th>
                                                    <th scope="col" class="px-3 py-3.5 text-left text-sm font-semibold">
                                                        {{ $t('components.documents.document_managers.add_relation') }}
                                                    </th>
                                                </tr>
                                            </thead>
                                            <tbody class="divide-y divide-base-500">
                                                <tr v-for="user in clipboardStore.$state.users" :key="user.userId">
                                                    <td
                                                        class="truncate whitespace-nowrap py-4 pl-4 pr-3 text-sm font-medium sm:pl-6 lg:pl-8"
                                                    >
                                                        <span class="inline-flex items-center gap-1">
                                                            <CitizenInfoPopover :user="user" :trailing="false" />
                                                            ({{ user.dateofbirth }})
                                                        </span>
                                                    </td>
                                                    <td class="whitespace-nowrap px-2 py-2 text-sm sm:px-4">
                                                        {{ user.jobLabel }}
                                                    </td>
                                                    <td class="whitespace-nowrap px-2 py-2 text-sm sm:px-4">
                                                        <div class="flex flex-row gap-2">
                                                            <UButton
                                                                :title="$t('components.documents.document_managers.mentioned')"
                                                                color="blue"
                                                                icon="i-mdi-at"
                                                                @click="addRelation(getUser(user), DocRelation.MENTIONED)"
                                                            />

                                                            <UButton
                                                                :title="$t('components.documents.document_managers.targets')"
                                                                color="amber"
                                                                icon="i-mdi-target"
                                                                @click="addRelation(getUser(user), DocRelation.TARGETS)"
                                                            />

                                                            <UButton
                                                                :title="$t('components.documents.document_managers.caused')"
                                                                color="red"
                                                                icon="i-mdi-source-commit-start"
                                                                @click="addRelation(getUser(user), DocRelation.CAUSED)"
                                                            />
                                                        </div>
                                                    </td>
                                                </tr>
                                            </tbody>
                                        </table>
                                    </div>
                                </div>
                            </div>
                        </template>

                        <template v-else-if="item.key === 'new'">
                            <UFormGroup name="name" :label="$t('common.search')">
                                <UInput
                                    v-model="queryCitizens"
                                    type="text"
                                    name="name"
                                    :placeholder="`${$t('common.citizen', 1)} ${$t('common.name')}`"
                                    leading-icon="i-mdi-search"
                                />
                            </UFormGroup>

                            <div class="mt-2 flow-root">
                                <div class="-my-2 mx-0 overflow-x-auto">
                                    <div class="inline-block min-w-full py-2 align-middle">
                                        <DataPendingBlock
                                            v-if="loading"
                                            :message="$t('common.loading', [$t('common.citizen', 2)])"
                                        />
                                        <DataErrorBlock
                                            v-else-if="error"
                                            :title="$t('common.unable_to_load', [$t('common.citizen', 2)])"
                                            :retry="refresh"
                                        />
                                        <DataNoDataBlock
                                            v-else-if="!citizens || citizens.length === 0"
                                            :message="$t('components.citizens.CitizensList.no_citizens')"
                                        />

                                        <table v-else class="min-w-full divide-y divide-base-200">
                                            <thead>
                                                <tr>
                                                    <th
                                                        scope="col"
                                                        class="py-3.5 pl-4 pr-3 text-left text-sm font-semibold sm:pl-6 lg:pl-8"
                                                    >
                                                        {{ $t('common.name') }}
                                                    </th>
                                                    <th scope="col" class="px-3 py-3.5 text-left text-sm font-semibold">
                                                        {{ $t('common.job', 1) }}
                                                    </th>
                                                    <th scope="col" class="px-3 py-3.5 text-left text-sm font-semibold">
                                                        {{ $t('components.documents.document_managers.add_relation') }}
                                                    </th>
                                                </tr>
                                            </thead>
                                            <tbody class="divide-y divide-base-500">
                                                <tr v-for="user in citizens.slice(0, 8)" :key="user.userId">
                                                    <td
                                                        class="truncate whitespace-nowrap py-4 pl-4 pr-3 text-sm font-medium sm:pl-6 lg:pl-8"
                                                    >
                                                        <span class="inline-flex items-center gap-1">
                                                            <CitizenInfoPopover :user="user" :trailing="false" />
                                                            ({{ user.dateofbirth }})
                                                        </span>
                                                    </td>
                                                    <td class="whitespace-nowrap px-2 py-2 text-sm sm:px-4">
                                                        {{ user.jobLabel }}
                                                    </td>
                                                    <td class="whitespace-nowrap px-2 py-2 text-sm sm:px-4">
                                                        <div class="flex flex-row gap-2">
                                                            <UButton
                                                                :title="$t('components.documents.document_managers.mentioned')"
                                                                color="blue"
                                                                icon="i-mdi-at"
                                                                @click="addRelation(user, DocRelation.MENTIONED)"
                                                            />

                                                            <UButton
                                                                :title="$t('components.documents.document_managers.targets')"
                                                                color="amber"
                                                                icon="i-mdi-target"
                                                                @click="addRelation(user, DocRelation.TARGETS)"
                                                            />

                                                            <UButton
                                                                :title="$t('components.documents.document_managers.caused')"
                                                                color="red"
                                                                icon="i-mdi-source-commit-start"
                                                                @click="addRelation(user, DocRelation.CAUSED)"
                                                            />
                                                        </div>
                                                    </td>
                                                </tr>
                                            </tbody>
                                        </table>
                                    </div>
                                </div>
                            </div>
                        </template>
                    </template>
                </UTabs>
            </div>

            <template #footer>
                <UButton block class="flex-1" color="black" @click="$emit('close')">
                    {{ $t('common.close', 1) }}
                </UButton>
            </template>
        </UCard>
    </UModal>
</template>
