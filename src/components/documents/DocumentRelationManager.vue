<script setup lang="ts">
import {
    Dialog,
    DialogPanel,
    DialogTitle,
    Tab,
    TabGroup,
    TabList,
    TabPanel,
    TabPanels,
    TransitionChild,
    TransitionRoot,
} from '@headlessui/vue';
import { watchDebounced } from '@vueuse/core';
import {
    AccountMinusIcon,
    AccountMultipleIcon,
    AccountSearchIcon,
    AtIcon,
    ClipboardListIcon,
    CloseIcon,
    OpenInNewIcon,
    SourceCommitStartIcon,
    TargetIcon,
    ViewListOutlineIcon,
} from 'mdi-vue3';
import { type DefineComponent } from 'vue';
import CitizenInfoPopover from '~/components/partials/citizens/CitizenInfoPopover.vue';
import DataNoDataBlock from '~/components/partials/data/DataNoDataBlock.vue';
import { useAuthStore } from '~/store/auth';
import { getUser, useClipboardStore } from '~/store/clipboard';
import { DocRelation, DocumentRelation } from '~~/gen/ts/resources/documents/documents';
import { User } from '~~/gen/ts/resources/users/users';
import DataPendingBlock from '~/components/partials/data/DataPendingBlock.vue';
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';

const { $grpc } = useNuxtApp();
const authStore = useAuthStore();
const clipboardStore = useClipboardStore();

const { activeChar } = storeToRefs(authStore);

const { t } = useI18n();

const props = defineProps<{
    open: boolean;
    document?: string;
    modelValue: Map<string, DocumentRelation>;
}>();

const emit = defineEmits<{
    (e: 'close'): void;
    (e: 'update:modelValue', payload: Map<string, DocumentRelation>): void;
}>();

const tabs = ref<{ name: string; icon: DefineComponent }[]>([
    {
        name: t('components.documents.document_managers.view_current'),
        icon: markRaw(ViewListOutlineIcon),
    },
    { name: t('common.clipboard'), icon: markRaw(ClipboardListIcon) },
    {
        name: t('components.documents.document_managers.add_new'),
        icon: markRaw(AccountSearchIcon),
    },
]);

const queryCitizens = ref('');

const {
    data: citizens,
    pending,
    refresh,
    error,
} = useLazyAsyncData(`document-${props.document?.toString()}-relations-citzens-${queryCitizens.value}`, () => listCitizens());

watchDebounced(queryCitizens, async () => await refresh(), {
    debounce: 600,
    maxWait: 1750,
});

async function listCitizens(): Promise<User[]> {
    try {
        const call = $grpc.getCitizenStoreClient().listCitizens({
            pagination: {
                offset: 0n,
            },
            searchName: queryCitizens.value,
        });
        const { response } = await call;

        return response.users.filter(
            (user) => !Array.from(props.modelValue.values()).find((r) => r.targetUserId === user.userId),
        );
    } catch (e) {
        $grpc.handleError(e as RpcError);
        throw e;
    }
}

function addRelation(user: User, relation: DocRelation): void {
    const keys = Array.from(props.modelValue.keys());
    const key = !keys.length ? '1' : (parseInt(keys[keys.length - 1]) + 1).toString();

    props.modelValue.set(key, {
        id: key,
        documentId: props.document ?? '0',
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
    <TransitionRoot as="template" :show="open">
        <Dialog as="div" class="relative z-30" @close="emit('close')">
            <TransitionChild
                as="template"
                enter="ease-out duration-300"
                enter-from="opacity-0"
                enter-to="opacity-100"
                leave="ease-in duration-200"
                leave-from="opacity-100"
                leave-to="opacity-0"
            >
                <div class="fixed inset-0 bg-base-900/75 transition-opacity" />
            </TransitionChild>

            <div class="fixed inset-0 z-30 overflow-y-auto">
                <div class="flex min-h-full items-end justify-center p-4 text-center sm:items-center sm:p-0">
                    <TransitionChild
                        as="template"
                        enter="ease-out duration-300"
                        enter-from="opacity-0 translate-y-4 sm:translate-y-0 sm:scale-95"
                        enter-to="opacity-100 translate-y-0 sm:scale-100"
                        leave="ease-in duration-200"
                        leave-from="opacity-100 translate-y-0 sm:scale-100"
                        leave-to="opacity-0 translate-y-4 sm:translate-y-0 sm:scale-95"
                    >
                        <DialogPanel
                            class="relative my-auto w-full overflow-hidden rounded-lg bg-base-800 px-4 pb-4 pt-5 text-left text-neutral transition-all sm:my-8 sm:max-w-6xl sm:p-6"
                        >
                            <div class="absolute right-0 top-0 hidden pr-4 pt-4 sm:block">
                                <button
                                    type="button"
                                    class="rounded-md transition-colors hover:text-base-300 focus:outline-none focus:ring-2 focus:ring-primary-500 focus:ring-offset-2"
                                    @click="emit('close')"
                                >
                                    <span class="sr-only">{{ $t('common.close') }}</span>
                                    <CloseIcon class="size-5" aria-hidden="true" />
                                </button>
                            </div>
                            <DialogTitle as="h3" class="text-base font-semibold leading-6">
                                {{ $t('common.citizen', 1) }}
                                {{ $t('common.relation', 2) }}
                            </DialogTitle>
                            <TabGroup>
                                <TabList class="mb-4 flex flex-row">
                                    <Tab v-for="tab in tabs" :key="tab.name" v-slot="{ selected }" class="w-full flex-initial">
                                        <button
                                            :class="[
                                                selected
                                                    ? 'border-primary-500 text-primary-500'
                                                    : 'border-transparent text-base-300 hover:border-base-300 hover:text-accent-200',
                                                'group inline-flex w-full items-center justify-center border-b-2 px-1 py-4 text-sm font-medium transition-colors',
                                            ]"
                                            :aria-current="selected ? 'page' : undefined"
                                        >
                                            <component
                                                :is="tab.icon"
                                                :class="[
                                                    selected ? 'text-primary-500' : 'text-base-300 group-hover:text-accent-200',
                                                    '-ml-0.5 mr-2 size-5 transition-colors',
                                                ]"
                                                aria-hidden="true"
                                            />
                                            <span>{{ tab.name }}</span>
                                        </button>
                                    </Tab>
                                </TabList>
                                <TabPanels>
                                    <div class="px-4 sm:flex sm:items-start sm:px-6 lg:px-8">
                                        <TabPanel class="w-full">
                                            <div class="flow-root">
                                                <div class="-my-2 mx-0 overflow-x-auto">
                                                    <div class="inline-block min-w-full py-2 align-middle">
                                                        <table class="min-w-full divide-y divide-base-200 text-neutral">
                                                            <thead>
                                                                <tr>
                                                                    <th
                                                                        scope="col"
                                                                        class="py-3.5 pl-4 pr-3 text-left text-sm font-semibold sm:pl-6 lg:pl-8"
                                                                    >
                                                                        {{ $t('common.name') }}
                                                                    </th>
                                                                    <th
                                                                        scope="col"
                                                                        class="px-3 py-3.5 text-left text-sm font-semibold"
                                                                    >
                                                                        {{ $t('common.creator') }}
                                                                    </th>
                                                                    <th
                                                                        scope="col"
                                                                        class="px-3 py-3.5 text-left text-sm font-semibold"
                                                                    >
                                                                        {{ $t('common.relation', 1) }}
                                                                    </th>
                                                                    <th
                                                                        scope="col"
                                                                        class="px-3 py-3.5 text-left text-sm font-semibold"
                                                                    >
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
                                                                            <CitizenInfoPopover :user="relation.targetUser" />
                                                                            ({{ relation.targetUser?.dateofbirth }})
                                                                        </span>
                                                                    </td>
                                                                    <td class="whitespace-nowrap px-3 py-4 text-sm">
                                                                        <CitizenInfoPopover :user="relation.sourceUser" />
                                                                    </td>
                                                                    <td class="whitespace-nowrap px-3 py-4 text-sm">
                                                                        {{
                                                                            $t(
                                                                                `enums.docstore.DocRelation.${
                                                                                    DocRelation[relation.relation]
                                                                                }`,
                                                                            )
                                                                        }}
                                                                    </td>
                                                                    <td class="whitespace-nowrap px-3 py-4 text-sm">
                                                                        <div class="flex flex-row gap-2">
                                                                            <div class="flex">
                                                                                <NuxtLink
                                                                                    :to="{
                                                                                        name: 'citizens-id',
                                                                                        params: {
                                                                                            id: relation.targetUserId,
                                                                                        },
                                                                                    }"
                                                                                    target="_blank"
                                                                                    data-te-toggle="tooltip"
                                                                                    :title="
                                                                                        $t(
                                                                                            'components.documents.document_managers.open_citizen',
                                                                                        )
                                                                                    "
                                                                                >
                                                                                    <OpenInNewIcon
                                                                                        class="h-auto w-5 text-primary-500 hover:text-primary-300"
                                                                                        aria-hidden="true"
                                                                                    />
                                                                                </NuxtLink>
                                                                            </div>
                                                                            <div class="flex">
                                                                                <button
                                                                                    role="button"
                                                                                    data-te-toggle="tooltip"
                                                                                    :title="
                                                                                        $t(
                                                                                            'components.documents.document_managers.remove_relation',
                                                                                        )
                                                                                    "
                                                                                    @click="removeRelation(relation.id!)"
                                                                                >
                                                                                    <AccountMinusIcon
                                                                                        class="h-auto w-5 text-error-400 hover:text-error-200"
                                                                                        aria-hidden="true"
                                                                                    />
                                                                                </button>
                                                                            </div>
                                                                        </div>
                                                                    </td>
                                                                </tr>
                                                            </tbody>
                                                        </table>
                                                    </div>
                                                </div>
                                            </div>
                                        </TabPanel>
                                        <TabPanel class="w-full">
                                            <div class="mt-2 flow-root">
                                                <div class="-my-2 mx-0 overflow-x-auto">
                                                    <div class="inline-block min-w-full py-2 align-middle">
                                                        <DataNoDataBlock
                                                            v-if="clipboardStore.$state.users.length === 0"
                                                            :type="$t('common.citizen', 2)"
                                                            :icon="AccountMultipleIcon"
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
                                                                    <th
                                                                        scope="col"
                                                                        class="px-3 py-3.5 text-left text-sm font-semibold"
                                                                    >
                                                                        {{ $t('common.job', 1) }}
                                                                    </th>
                                                                    <th
                                                                        scope="col"
                                                                        class="px-3 py-3.5 text-left text-sm font-semibold"
                                                                    >
                                                                        {{
                                                                            $t(
                                                                                'components.documents.document_managers.add_relation',
                                                                            )
                                                                        }}
                                                                    </th>
                                                                </tr>
                                                            </thead>
                                                            <tbody class="divide-y divide-base-500">
                                                                <tr
                                                                    v-for="user in clipboardStore.$state.users"
                                                                    :key="user.userId"
                                                                >
                                                                    <td
                                                                        class="truncate whitespace-nowrap py-4 pl-4 pr-3 text-sm font-medium sm:pl-6 lg:pl-8"
                                                                    >
                                                                        <span class="inline-flex items-center gap-1">
                                                                            <CitizenInfoPopover :user="user" />
                                                                            ({{ user.dateofbirth }})
                                                                        </span>
                                                                    </td>
                                                                    <td class="whitespace-nowrap px-3 py-4 text-sm">
                                                                        {{ user.jobLabel }}
                                                                    </td>
                                                                    <td class="whitespace-nowrap px-3 py-4 text-sm">
                                                                        <div class="flex flex-row gap-2">
                                                                            <div class="flex">
                                                                                <button
                                                                                    role="button"
                                                                                    data-te-toggle="tooltip"
                                                                                    :title="
                                                                                        $t(
                                                                                            'components.documents.document_managers.mentioned',
                                                                                        )
                                                                                    "
                                                                                    @click="
                                                                                        addRelation(
                                                                                            getUser(user),
                                                                                            DocRelation.MENTIONED,
                                                                                        )
                                                                                    "
                                                                                >
                                                                                    <AtIcon
                                                                                        class="h-auto w-5 text-success-500 hover:text-success-300"
                                                                                        aria-hidden="true"
                                                                                    />
                                                                                </button>
                                                                            </div>
                                                                            <div class="flex">
                                                                                <button
                                                                                    role="button"
                                                                                    data-te-toggle="tooltip"
                                                                                    :title="
                                                                                        $t(
                                                                                            'components.documents.document_managers.targets',
                                                                                        )
                                                                                    "
                                                                                    @click="
                                                                                        addRelation(
                                                                                            getUser(user),
                                                                                            DocRelation.TARGETS,
                                                                                        )
                                                                                    "
                                                                                >
                                                                                    <TargetIcon
                                                                                        class="h-auto w-5 text-warn-400 hover:text-warn-200"
                                                                                        aria-hidden="true"
                                                                                    />
                                                                                </button>
                                                                            </div>
                                                                            <div class="flex">
                                                                                <button
                                                                                    role="button"
                                                                                    data-te-toggle="tooltip"
                                                                                    :title="
                                                                                        $t(
                                                                                            'components.documents.document_managers.caused',
                                                                                        )
                                                                                    "
                                                                                    @click="
                                                                                        addRelation(
                                                                                            getUser(user),
                                                                                            DocRelation.CAUSED,
                                                                                        )
                                                                                    "
                                                                                >
                                                                                    <SourceCommitStartIcon
                                                                                        class="h-auto w-5 text-error-400 hover:text-error-200"
                                                                                        aria-hidden="true"
                                                                                    />
                                                                                </button>
                                                                            </div>
                                                                        </div>
                                                                    </td>
                                                                </tr>
                                                            </tbody>
                                                        </table>
                                                    </div>
                                                </div>
                                            </div>
                                        </TabPanel>
                                        <TabPanel class="w-full">
                                            <div>
                                                <label for="name" class="sr-only">Name</label>
                                                <input
                                                    v-model="queryCitizens"
                                                    type="text"
                                                    name="name"
                                                    class="block w-full rounded-md border-0 bg-base-700 py-1.5 text-neutral placeholder:text-accent-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6"
                                                    :placeholder="`${$t('common.citizen', 1)} ${$t('common.name')}`"
                                                    @focusin="focusTablet(true)"
                                                    @focusout="focusTablet(false)"
                                                />
                                            </div>
                                            <div class="mt-2 flow-root">
                                                <div class="-my-2 mx-0 overflow-x-auto">
                                                    <div class="inline-block min-w-full py-2 align-middle">
                                                        <DataPendingBlock
                                                            v-if="pending"
                                                            :message="$t('common.loading', [$t('common.citizen', 2)])"
                                                        />
                                                        <DataErrorBlock
                                                            v-else-if="error"
                                                            :title="$t('common.unable_to_load', [$t('common.citizen', 2)])"
                                                            :retry="refresh"
                                                        />
                                                        <DataNoDataBlock
                                                            v-else-if="citizens === null || citizens.length === 0"
                                                            :message="$t('components.citizens.citizens_list.no_citizens')"
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
                                                                    <th
                                                                        scope="col"
                                                                        class="px-3 py-3.5 text-left text-sm font-semibold"
                                                                    >
                                                                        {{ $t('common.job', 1) }}
                                                                    </th>
                                                                    <th
                                                                        scope="col"
                                                                        class="px-3 py-3.5 text-left text-sm font-semibold"
                                                                    >
                                                                        {{
                                                                            $t(
                                                                                'components.documents.document_managers.add_relation',
                                                                            )
                                                                        }}
                                                                    </th>
                                                                </tr>
                                                            </thead>
                                                            <tbody class="divide-y divide-base-500">
                                                                <tr v-for="user in citizens.slice(0, 8)" :key="user.userId">
                                                                    <td
                                                                        class="truncate whitespace-nowrap py-4 pl-4 pr-3 text-sm font-medium sm:pl-6 lg:pl-8"
                                                                    >
                                                                        <span class="inline-flex items-center gap-1">
                                                                            <CitizenInfoPopover :user="user" />
                                                                            ({{ user.dateofbirth }})
                                                                        </span>
                                                                    </td>
                                                                    <td class="whitespace-nowrap px-3 py-4 text-sm">
                                                                        {{ user.jobLabel }}
                                                                    </td>
                                                                    <td class="whitespace-nowrap px-3 py-4 text-sm">
                                                                        <div class="flex flex-row gap-2">
                                                                            <div class="flex">
                                                                                <button
                                                                                    role="button"
                                                                                    data-te-toggle="tooltip"
                                                                                    :title="
                                                                                        $t(
                                                                                            'components.documents.document_managers.mentioned',
                                                                                        )
                                                                                    "
                                                                                    @click="
                                                                                        addRelation(user, DocRelation.MENTIONED)
                                                                                    "
                                                                                >
                                                                                    <AtIcon
                                                                                        class="h-auto w-5 text-success-500 hover:text-success-300"
                                                                                        aria-hidden="true"
                                                                                    />
                                                                                </button>
                                                                            </div>
                                                                            <div class="flex">
                                                                                <button
                                                                                    role="button"
                                                                                    data-te-toggle="tooltip"
                                                                                    :title="
                                                                                        $t(
                                                                                            'components.documents.document_managers.targets',
                                                                                        )
                                                                                    "
                                                                                    @click="
                                                                                        addRelation(user, DocRelation.TARGETS)
                                                                                    "
                                                                                >
                                                                                    <TargetIcon
                                                                                        class="h-auto w-5 text-warn-400 hover:text-warn-200"
                                                                                        aria-hidden="true"
                                                                                    />
                                                                                </button>
                                                                            </div>
                                                                            <div class="flex">
                                                                                <button
                                                                                    role="button"
                                                                                    data-te-toggle="tooltip"
                                                                                    :title="
                                                                                        $t(
                                                                                            'components.documents.document_managers.caused',
                                                                                        )
                                                                                    "
                                                                                    @click="
                                                                                        addRelation(user, DocRelation.CAUSED)
                                                                                    "
                                                                                >
                                                                                    <SourceCommitStartIcon
                                                                                        class="h-auto w-5 text-error-400 hover:text-error-200"
                                                                                        aria-hidden="true"
                                                                                    />
                                                                                </button>
                                                                            </div>
                                                                        </div>
                                                                    </td>
                                                                </tr>
                                                            </tbody>
                                                        </table>
                                                    </div>
                                                </div>
                                            </div>
                                        </TabPanel>
                                    </div>
                                </TabPanels>
                            </TabGroup>
                            <div class="mt-5 gap-2 sm:mt-4 sm:flex sm:flex-row-reverse">
                                <button
                                    type="button"
                                    class="rounded-md bg-base-500 px-3.5 py-2.5 text-sm font-semibold text-neutral hover:bg-base-400"
                                    @click="emit('close')"
                                >
                                    {{ $t('common.close') }}
                                </button>
                            </div>
                        </DialogPanel>
                    </TransitionChild>
                </div>
            </div>
        </Dialog>
    </TransitionRoot>
</template>
