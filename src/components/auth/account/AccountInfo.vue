<script lang="ts" setup>
import { Switch, SwitchGroup, SwitchLabel } from '@headlessui/vue';
import { RpcError } from '@protobuf-ts/runtime-rpc';
import { AccountIcon } from 'mdi-vue3';
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import DataNoDataBlock from '~/components/partials/data/DataNoDataBlock.vue';
import DataPendingBlock from '~/components/partials/data/DataPendingBlock.vue';
import { useAuthStore } from '~/store/auth';
import { useSettingsStore } from '~/store/settings';
import { GetAccountInfoResponse } from '~~/gen/ts/services/auth/auth';
import ChangePasswordModal from '~/components/auth/account/ChangePasswordModal.vue';
import DebugInfo from '~/components/auth/account/DebugInfo.vue';
import OAuth2Connections from '~/components/auth/account/OAuth2Connections.vue';

const { $grpc } = useNuxtApp();

const authStore = useAuthStore();
const { activeChar } = storeToRefs(authStore);

const settings = useSettingsStore();
const { startpage, documents } = storeToRefs(settings);

const { data: account, pending, refresh, error } = useLazyAsyncData(`accountinfo`, () => getAccountInfo());

async function getAccountInfo(): Promise<GetAccountInfoResponse | undefined> {
    try {
        const call = $grpc.getAuthClient().getAccountInfo({});

        return call.response;
    } catch (e) {
        $grpc.handleError(e as RpcError);
        throw e;
    }
}

const changePasswordModal = ref(false);

async function removeOAuth2Connection(provider: string): Promise<void> {
    const idx = account.value?.oauth2Connections.findIndex((v) => v.providerName === provider);
    if (idx !== undefined && idx > -1) {
        account.value?.oauth2Connections.splice(idx, 1);

        await refresh();
    }
}

const homepages: { name: string; path: string; permission?: string }[] = [
    { name: 'common.home', path: '/overview' },
    { name: 'pages.citizens.title', path: '/citizens', permission: 'CitizenStoreService.ListCitizens' },
    { name: 'pages.vehicles.title', path: '/vehicles', permission: 'DMVService.ListVehicles' },
    { name: 'pages.documents.title', path: '/documents', permission: 'DocStoreService.ListDocuments' },
    { name: 'pages.jobs.overview.title', path: '/jobs/overview', permission: 'JobsService.ColleaguesList' },
    { name: 'common.livemap', path: '/livemap', permission: 'LivemapperService.Stream' },
    { name: 'common.dispatch_center', path: '/centrum', permission: 'CentrumService.TakeControl' },
];

const selectedHomepage = ref<(typeof homepages)[0]>();
watch(selectedHomepage, () => (startpage.value = selectedHomepage.value?.path ?? '/overview'));

onBeforeMount(async () => {
    selectedHomepage.value = homepages.find((h) => h.path === startpage.value);
});

const darkModeActive = ref(documents.value.editorTheme === 'dark');

watch(darkModeActive, async () => {
    if (darkModeActive.value) {
        documents.value.editorTheme = 'dark';
    } else {
        documents.value.editorTheme = 'default';
    }
});
</script>

<template>
    <div class="py-2 mt-5 max-w-5xl mx-auto">
        <ChangePasswordModal :open="changePasswordModal" @close="changePasswordModal = false" />
        <DataPendingBlock v-if="pending" :message="$t('common.loading', [`${$t('common.account')} ${$t('common.info')}`])" />
        <DataErrorBlock
            v-else-if="error"
            :title="$t('common.unable_to_load', [`${$t('common.account')} ${$t('common.info')}`])"
            :retry="refresh"
        />
        <DataNoDataBlock v-else-if="!account" :type="`${$t('common.account')} ${$t('common.data')}`" :icon="AccountIcon" />
        <div v-else>
            <div class="overflow-hidden bg-base-800 shadow sm:rounded-lg text-neutral">
                <div class="px-4 py-5 sm:px-6">
                    <h3 class="text-base font-semibold leading-6">
                        {{ $t('components.auth.account_info.title') }}
                    </h3>
                    <p class="mt-1 max-w-2xl text-sm">
                        {{ $t('components.auth.account_info.subtitle') }}
                    </p>
                </div>
                <div class="border-t border-base-400 px-4 py-5 sm:p-0">
                    <dl class="sm:divide-y sm:divide-base-400">
                        <div class="py-4 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-6 sm:py-5">
                            <dt class="text-sm font-medium">
                                {{ $t('common.username') }}
                            </dt>
                            <dd class="mt-1 text-sm sm:col-span-2 sm:mt-0">
                                {{ account.account?.username }}
                            </dd>
                        </div>
                        <div class="py-4 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-6 sm:py-5">
                            <dt class="text-sm font-medium">
                                {{ $t('components.auth.account_info.license') }}
                            </dt>
                            <dd class="mt-1 text-sm sm:col-span-2 sm:mt-0">
                                {{ account.account?.license }}
                            </dd>
                        </div>
                        <div class="py-4 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-6 sm:py-5">
                            <dt class="text-sm font-medium">
                                {{ $t('components.auth.account_info.change_password') }}
                            </dt>
                            <dd class="mt-1 text-sm sm:col-span-2 sm:mt-0">
                                <button
                                    type="button"
                                    class="rounded-md bg-base-500 py-2.5 px-3.5 text-sm font-semibold text-neutral hover:bg-base-400"
                                    @click="changePasswordModal = true"
                                >
                                    {{ $t('components.auth.account_info.change_password_button') }}
                                </button>
                            </dd>
                        </div>
                        <div class="py-4 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-6 sm:py-5">
                            <dt class="text-sm font-medium">
                                {{ $t('components.auth.account_info.set_startpage.title') }}
                            </dt>
                            <dd class="mt-1 text-sm sm:col-span-2 sm:mt-0">
                                <select
                                    v-if="activeChar"
                                    v-model="selectedHomepage"
                                    class="block w-full rounded-md border-0 py-1.5 bg-base-700 text-neutral placeholder:text-base-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6"
                                >
                                    <option
                                        v-for="page in homepages"
                                        :key="page.path"
                                        :value="page"
                                        :disabled="!(page.permission === undefined || can(page.permission))"
                                    >
                                        {{ $t(page.name ?? 'common.page') }}
                                    </option>
                                </select>
                                <p v-else class="text-neutral">
                                    {{ $t('components.auth.account_info.set_startpage.no_char_selected') }}
                                </p>
                            </dd>
                        </div>
                        <div class="py-4 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-6 sm:py-5">
                            <dt class="text-sm font-medium">
                                {{ $t('components.auth.account_info.editor_theme.title') }}
                            </dt>
                            <dd class="mt-1 text-sm sm:col-span-2 sm:mt-0">
                                <SwitchGroup as="div" class="flex items-center">
                                    <Switch
                                        v-model="darkModeActive"
                                        :class="[
                                            documents.editorTheme === 'dark' ? 'bg-indigo-600' : 'bg-gray-200',
                                            'relative inline-flex h-6 w-11 flex-shrink-0 cursor-pointer rounded-full border-2 border-transparent transition-colors duration-200 ease-in-out focus:outline-none focus:ring-2 focus:ring-indigo-600 focus:ring-offset-2',
                                        ]"
                                    >
                                        <span
                                            aria-hidden="true"
                                            :class="[
                                                documents.editorTheme === 'dark' ? 'translate-x-5' : 'translate-x-0',
                                                'pointer-events-none inline-block h-5 w-5 transform rounded-full bg-white shadow ring-0 transition duration-200 ease-in-out',
                                            ]"
                                        />
                                    </Switch>
                                    <SwitchLabel as="span" class="ml-3 text-sm">
                                        <span class="font-medium text-gray-300">{{
                                            $t('components.auth.account_info.editor_theme.dark_mode')
                                        }}</span>
                                    </SwitchLabel>
                                </SwitchGroup>
                            </dd>
                        </div>
                    </dl>
                </div>
            </div>

            <OAuth2Connections
                v-if="account"
                :providers="account.oauth2Providers"
                :connections="account.oauth2Connections"
                @click="removeOAuth2Connection($event)"
            />

            <DebugInfo />
        </div>
    </div>
</template>
