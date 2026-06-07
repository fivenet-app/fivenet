<script lang="ts" setup>
import { Pane, Splitpanes } from 'splitpanes';
import 'splitpanes/dist/splitpanes.css';
import { z } from 'zod';
import DispatchList from '~/components/dispatch/dispatches/DispatchList.vue';
import DispatchLayer from '~/components/dispatch/livemap/DispatchLayer.vue';
import BaseMap from '~/components/livemap/BaseMap.vue';
import Pagination from '~/components/partials/Pagination.vue';
import SelectMenu from '~/components/partials/SelectMenu.vue';
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import DataNoDataBlock from '~/components/partials/data/DataNoDataBlock.vue';
import DataPendingBlock from '~/components/partials/data/DataPendingBlock.vue';
import { useLivemapStore } from '~/stores/livemap';
import { getCentrumDispatchesClient } from '~~/gen/ts/clients';
import type { UserShort } from '~~/gen/ts/resources/users/short/user';
import type { ListDispatchesRequest, ListDispatchesResponse } from '~~/gen/ts/services/centrum/dispatches';

useHead({
    title: 'common.dispatches',
});

definePageMeta({
    title: 'common.dispatches',
    requiresAuth: true,
    permission: 'centrum.CentrumService/TakeControl',
});

const livemapStore = useLivemapStore();
const { showLocationMarker } = storeToRefs(livemapStore);

const completorStore = useCompletorStore();

const centrumDispatchesClient = await getCentrumDispatchesClient();

const schema = z.object({
    postal: z.coerce.string().trim().max(12).default(''),
    id: z.coerce.number().max(16).optional(),
    creatorIds: z.coerce.number().array().max(5).default([]),
    page: pageNumberSchema,
});

const query = useSearchForm('centrum_dispatches_archive', schema);

const { data, status, refresh, error } = useLazyAsyncData(`centrum-dispatches-${query.page}`, () => listDispatches());

async function listDispatches(): Promise<ListDispatchesResponse> {
    try {
        const req: ListDispatchesRequest = {
            pagination: {
                offset: calculateOffset(query.page, data.value?.pagination),
            },
            notStatus: [],
            status: [],
            ids: [],
            postal: query.postal.replaceAll('-', '').replace(/\D/g, ''),
            creatorIds: query.creatorIds,
        };

        if (query.id && query.id > 0) {
            req.ids.push(query.id);
        }

        const call = centrumDispatchesClient.listDispatches(req);
        const { response } = await call;

        return response;
    } catch (e) {
        handleGRPCError(e as RpcError);
        throw e;
    }
}

useDebouncedRefresh(query, refresh, {
    debounce: 200,
    maxWait: 1250,
});

const baseMapRef = useTemplateRef('baseMapRef');
const mapResizeFn = () => baseMapRef.value?.mapResize();

onBeforeMount(() => (showLocationMarker.value = true));
onMounted(async () => {
    useTimeoutFn(() => (mount.value = true), 35);

    nuiEvents.on('openTablet', mapResizeFn);
});

onBeforeUnmount(() => {
    showLocationMarker.value = false;
});
onUnmounted(() => {
    nuiEvents.off('openTablet', mapResizeFn);
});

const input = useTemplateRef('input');

defineShortcuts({
    '/': () => input.value?.inputRef?.focus(),
});

const mount = ref<boolean>(false);
</script>

<template>
    <UDashboardPanel :ui="{ body: 'p-0 sm:p-0 gap-0 sm:gap-0' }">
        <template #header>
            <UDashboardNavbar :title="$t('common.dispatches')">
                <template #leading>
                    <UDashboardSidebarCollapse />
                </template>

                <template #right>
                    <PartialsBackButton fallback-to="/dispatch" />
                </template>
            </UDashboardNavbar>
        </template>

        <template #body>
            <div
                class="max-h-[calc(100dvh-var(--ui-header-height))] min-h-[calc(100dvh-var(--ui-header-height))] overflow-hidden"
            >
                <Splitpanes v-if="mount" class="relative">
                    <Pane :min-size="25">
                        <ClientOnly>
                            <BaseMap ref="baseMapRef" :map-options="{ zoomControl: false }">
                                <template #default>
                                    <LazyLivemapTempMarker />

                                    <DispatchLayer show-all-dispatches :dispatch-list="data?.dispatches ?? []" />
                                </template>
                            </BaseMap>
                        </ClientOnly>
                    </Pane>

                    <Pane :size="65">
                        <div class="max-h-full overflow-y-auto">
                            <div class="mb-2 px-2">
                                <UForm class="flex flex-row gap-2" :schema="schema" :state="query" @submit="refresh()">
                                    <UFormField class="flex-1" name="postal" :label="$t('common.postal')">
                                        <UInput
                                            ref="input"
                                            v-model="query.postal"
                                            class="w-full"
                                            type="text"
                                            name="postal"
                                            :placeholder="$t('common.postal')"
                                        >
                                            <template #trailing>
                                                <UKbd value="/" />
                                            </template>
                                        </UInput>
                                    </UFormField>

                                    <UFormField class="flex-1" name="id" :label="$t('common.id')">
                                        <UInput
                                            v-model="query.id"
                                            class="w-full"
                                            type="text"
                                            name="id"
                                            :min="1"
                                            :max="99999999999"
                                            :placeholder="$t('common.id')"
                                        />
                                    </UFormField>

                                    <UFormField class="flex-1" name="creator" :label="$t('common.creator')">
                                        <SelectMenu
                                            v-model="query.creatorIds"
                                            class="w-full"
                                            multiple
                                            nullable
                                            :searchable="
                                                async (q: string): Promise<UserShort[]> =>
                                                    await completorStore.completeCitizens({
                                                        search: q,
                                                        userIds: query.creatorIds,
                                                    })
                                            "
                                            searchable-key="completor-citizens"
                                            :search-input="{ placeholder: $t('common.search_field') }"
                                            :filter-fields="['firstname', 'lastname']"
                                            :placeholder="$t('common.creator')"
                                            trailing
                                            value-key="userId"
                                        >
                                            <template #default="{ modelValue }">
                                                {{ $t('common.selected', modelValue?.length ?? 0) }}
                                            </template>

                                            <template #item-label="{ item }">
                                                {{ userToLabel(item) }}
                                            </template>

                                            <template #empty>
                                                {{ $t('common.not_found', [$t('common.creator', 2)]) }}
                                            </template>
                                        </SelectMenu>
                                    </UFormField>
                                </UForm>
                            </div>

                            <DataPendingBlock
                                v-if="isRequestPending(status)"
                                :message="$t('common.loading', [$t('common.dispatches')])"
                            />
                            <DataErrorBlock
                                v-else-if="error"
                                :title="$t('common.unable_to_load', [$t('common.dispatches')])"
                                :error="error"
                                :retry="refresh"
                            />
                            <DataNoDataBlock
                                v-else-if="data?.dispatches.length === 0"
                                icon="i-mdi-car-emergency"
                                :type="$t('common.dispatches')"
                            />

                            <div v-else class="relative overflow-x-auto">
                                <DispatchList
                                    :show-button="false"
                                    hide-actions
                                    always-show-day
                                    :dispatches="data?.dispatches"
                                />
                            </div>

                            <Pagination
                                v-model="query.page"
                                :pagination="data?.pagination"
                                :status="status"
                                :refresh="refresh"
                            />
                        </div>
                    </Pane>
                </Splitpanes>
            </div>
        </template>
    </UDashboardPanel>
</template>

<style scoped>
.splitpanes--vertical > .splitpanes__splitter {
    min-width: 2px;
    background-color: var(--color-gray-800);
}

.splitpanes--horizontal > .splitpanes__splitter {
    min-height: 2px;
    background-color: var(--color-gray-800);
}
</style>
