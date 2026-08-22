<script lang="ts" setup>
import type { Form, SplitterItem } from '@nuxt/ui';
import { z } from 'zod';
import DispatchList from '~/components/dispatch/dispatches/DispatchList.vue';
import DispatchLayer from '~/components/dispatch/livemap/DispatchLayer.vue';
import BaseMap from '~/components/livemap/BaseMap.vue';
import Pagination from '~/components/partials/Pagination.vue';
import SelectMenu from '~/components/partials/SelectMenu.vue';
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import DataNoDataBlock from '~/components/partials/data/DataNoDataBlock.vue';
import DataPendingBlock from '~/components/partials/data/DataPendingBlock.vue';
import DispatchCenterLayoutPopover from '~/components/dispatch/DispatchCenterLayoutPopover.vue';
import { useLivemapStore } from '~/stores/livemap';
import { useSettingsStore, type DispatchCenterOuterPane } from '~/stores/settings';
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

const settingsStore = useSettingsStore();
const { centrum } = storeToRefs(settingsStore);

const completorStore = useCompletorStore();

const centrumDispatchesClient = await getCentrumDispatchesClient();

const schema = z.object({
    postal: z.coerce.string().trim().max(12).default(''),
    id: z.coerce.number().max(16).optional(),
    creatorIds: z.coerce.number().array().max(5).default([]),
    page: pageNumberSchema,
});

type Schema = z.output<typeof schema>;

const query = useSearchForm('centrum_dispatches_archive', schema);

const formRef = useTemplateRef<Form<typeof schema>>('formRef');
const { validatedQuery, commitValidatedQuery } = useFormSearchValidation<typeof schema>(query, formRef);

const dispatchesKey = computed(() => `centrum-dispatches-${JSON.stringify(validatedQuery.value)}`);

const { data, status, refresh, error } = useLazyAsyncData(dispatchesKey, () => listDispatches(validatedQuery.value));

async function listDispatches(values: Schema): Promise<ListDispatchesResponse> {
    try {
        const req: ListDispatchesRequest = {
            pagination: {
                offset: calculateOffset(values.page, data.value?.pagination),
            },
            notStatus: [],
            status: [],
            ids: [],
            postal: values.postal.replaceAll('-', '').replace(/\D/g, ''),
            creatorIds: values.creatorIds,
        };

        if (values.id && values.id > 0) {
            req.ids.push(values.id);
        }

        const call = centrumDispatchesClient.listDispatches(req);
        const { response } = await call;

        return response;
    } catch (e) {
        handleGRPCError(e as RpcError);
        throw e;
    }
}

const baseMapRef = useTemplateRef('baseMapRef');
const mapResizeFn = () => baseMapRef.value?.mapResize();

const roundPaneSize = (size: number): number => Math.round(size * 100) / 100;

const archivePaneItems = {
    map: {
        slot: 'map',
        minSize: 25,
    },
    sidebar: {
        slot: 'details',
        minSize: 40,
    },
} satisfies Record<DispatchCenterOuterPane, Pick<SplitterItem, 'slot' | 'minSize'>>;

const dispatchesSplitterItems = computed<SplitterItem[]>(() =>
    centrum.value.dispatchCenterPaneLayout.outer.map((pane) => ({
        ...archivePaneItems[pane],
        defaultSize: centrum.value.dispatchCenterPaneSizes[pane],
    })),
);

const dispatchesSplitterKey = computed(() => `dispatches-splitter-${centrum.value.dispatchCenterPaneLayout.outer.join('-')}`);

const splitterUi = {
    panel: 'min-h-0 min-w-0 overflow-hidden',
    handle: 'data-[orientation=horizontal]:w-px data-[orientation=vertical]:h-px bg-border transition-colors data-[state=hover]:bg-primary data-[state=drag]:bg-primary',
};

function onDispatchesLayout(sizes: number[]): void {
    const panes = centrum.value.dispatchCenterPaneLayout.outer;
    if (sizes.length < panes.length) return;

    panes.forEach((pane, index) => {
        const size = sizes[index];
        if (size === undefined) return;

        centrum.value.dispatchCenterPaneSizes[pane] = roundPaneSize(size);
    });
}

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

const inputRef = useTemplateRef('inputRef');

defineShortcuts({
    '/': () => inputRef.value?.inputRef?.focus(),
});

const mount = ref<boolean>(false);
</script>

<template>
    <UDashboardPanel :ui="{ root: 'pb-(--page-content-bottom-offset)', body: 'p-0 sm:p-0 gap-0 sm:gap-0' }">
        <template #header>
            <UDashboardNavbar :title="$t('common.dispatches')">
                <template #leading>
                    <UDashboardSidebarCollapse />
                </template>

                <template #right>
                    <div class="flex items-center gap-2">
                        <DispatchCenterLayoutPopover hide-inner-panes />

                        <PartialsBackButton fallback-to="/dispatch" />
                    </div>
                </template>
            </UDashboardNavbar>
        </template>

        <template #body>
            <div
                class="max-h-[calc(100dvh-var(--ui-header-height)-var(--page-content-bottom-offset))] min-h-[calc(100dvh-var(--ui-header-height)-var(--page-content-bottom-offset))] overflow-hidden"
            >
                <USplitter
                    v-if="mount"
                    id="dispatches-splitter"
                    :key="dispatchesSplitterKey"
                    class="relative size-full"
                    :items="dispatchesSplitterItems"
                    :ui="splitterUi"
                    @layout="onDispatchesLayout"
                >
                    <template #map>
                        <ClientOnly>
                            <BaseMap ref="baseMapRef" :map-options="{ zoomControl: false }">
                                <template #default>
                                    <LazyLivemapTempMarker />

                                    <DispatchLayer show-all-dispatches :dispatch-list="data?.dispatches ?? []" />
                                </template>
                            </BaseMap>
                        </ClientOnly>
                    </template>

                    <template #details>
                        <div class="max-h-full overflow-y-auto">
                            <div class="mb-2 px-2">
                                <UForm
                                    ref="formRef"
                                    class="flex flex-row gap-2"
                                    :schema="schema"
                                    :state="query"
                                    @submit="commitValidatedQuery"
                                >
                                    <UFormField class="flex-1" name="postal" :label="$t('common.postal')">
                                        <UInput
                                            ref="inputRef"
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
                    </template>
                </USplitter>
            </div>
        </template>
    </UDashboardPanel>
</template>
