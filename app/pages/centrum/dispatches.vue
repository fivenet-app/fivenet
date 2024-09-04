<script lang="ts" setup>
import { Pane, Splitpanes } from 'splitpanes';
import 'splitpanes/dist/splitpanes.css';
import { z } from 'zod';
import DispatchList from '~/components/centrum/dispatches/DispatchList.vue';
import Pagination from '~/components/partials/Pagination.vue';
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import DataNoDataBlock from '~/components/partials/data/DataNoDataBlock.vue';
import DataPendingBlock from '~/components/partials/data/DataPendingBlock.vue';
import { useLivemapStore } from '~/store/livemap';
import type { ListDispatchesRequest, ListDispatchesResponse } from '~~/gen/ts/services/centrum/centrum';

useHead({
    title: 'common.dispatches',
});
definePageMeta({
    title: 'common.dispatches',
    requiresAuth: true,
    permission: 'CentrumService.TakeControl',
});

const livemapStore = useLivemapStore();
const { showLocationMarker } = storeToRefs(livemapStore);

const schema = z.object({
    postal: z.string().trim().max(12),
    id: z.string().trim().max(16),
});

type Schema = z.output<typeof schema>;

const query = reactive<Schema>({
    postal: '',
    id: '',
});

const page = ref(1);
const offset = computed(() => (data.value?.pagination?.pageSize ? data.value?.pagination?.pageSize * (page.value - 1) : 0));

const { data, pending: loading, refresh, error } = useLazyAsyncData(`centrum-dispatches-${page.value}`, () => listDispatches());

async function listDispatches(): Promise<ListDispatchesResponse> {
    try {
        const req: ListDispatchesRequest = {
            pagination: {
                offset: offset.value,
            },
            notStatus: [],
            status: [],
            ids: [],
            postal: query.postal.replaceAll('-', '').replace(/\D/g, ''),
        };

        if (query.id) {
            const id = query.id.replaceAll('-', '').replace(/\D/g, '');
            if (id.length > 0) {
                req.ids.push(id);
            }
        }

        const call = getGRPCCentrumClient().listDispatches(req);
        const { response } = await call;

        return response;
    } catch (e) {
        handleGRPCError(e as RpcError);
        throw e;
    }
}

watch(offset, async () => refresh());

watchDebounced(query, async () => refresh(), {
    debounce: 200,
    maxWait: 1250,
});

onMounted(() => {
    showLocationMarker.value = true;
});

onBeforeUnmount(() => {
    showLocationMarker.value = false;
});

const input = ref<{ input: HTMLInputElement }>();

defineShortcuts({
    '/': () => {
        input.value?.input?.focus();
    },
});

const mount = ref(false);
onMounted(async () => useTimeoutFn(() => (mount.value = true), 35));
</script>

<template>
    <UDashboardPage>
        <UDashboardPanel grow>
            <UDashboardNavbar :title="$t('common.dispatches')">
                <template #right>
                    <UButton color="black" icon="i-mdi-arrow-back" to="/centrum">
                        {{ $t('common.back') }}
                    </UButton>
                </template>
            </UDashboardNavbar>

            <div class="max-h-[calc(100dvh-var(--header-height))] min-h-[calc(100dvh-var(--header-height))] overflow-hidden">
                <Splitpanes v-if="mount" class="relative">
                    <Pane :min-size="25">
                        <ClientOnly>
                            <LazyLivemapBaseMap :map-options="{ zoomControl: false }">
                                <template #default>
                                    <LazyLivemapMapTempMarker />
                                </template>
                            </LazyLivemapBaseMap>
                        </ClientOnly>
                    </Pane>

                    <Pane :size="65">
                        <div class="max-h-full overflow-y-auto">
                            <div class="mb-2 px-2">
                                <UForm :schema="schema" :state="query" class="flex flex-row gap-2" @submit="refresh()">
                                    <UFormGroup name="postal" :label="$t('common.postal')" class="flex-1">
                                        <UInput
                                            ref="input"
                                            v-model="query.postal"
                                            type="text"
                                            name="postal"
                                            :placeholder="$t('common.postal')"
                                        >
                                            <template #trailing>
                                                <UKbd value="/" />
                                            </template>
                                        </UInput>
                                    </UFormGroup>
                                    <UFormGroup name="id" :label="$t('common.id')" class="flex-1">
                                        <UInput
                                            v-model="query.id"
                                            type="text"
                                            name="id"
                                            :min="1"
                                            :max="99999999999"
                                            :placeholder="$t('common.id')"
                                        />
                                    </UFormGroup>
                                </UForm>
                            </div>

                            <DataPendingBlock v-if="loading" :message="$t('common.loading', [$t('common.dispatches')])" />
                            <DataErrorBlock
                                v-else-if="error"
                                :title="$t('common.unable_to_load', [$t('common.dispatches')])"
                                :retry="refresh"
                            />
                            <DataNoDataBlock v-else-if="data?.dispatches.length === 0" :type="$t('common.dispatches')" />

                            <div v-else class="relative overflow-x-auto">
                                <DispatchList
                                    :show-button="false"
                                    :hide-actions="true"
                                    :always-show-day="true"
                                    :dispatches="data?.dispatches"
                                />
                            </div>

                            <Pagination v-model="page" :pagination="data?.pagination" :loading="loading" :refresh="refresh" />
                        </div>
                    </Pane>
                </Splitpanes>
            </div>
        </UDashboardPanel>
    </UDashboardPage>
</template>

<style scoped>
.splitpanes--vertical > .splitpanes__splitter {
    min-width: 2px;
    background-color: rgb(var(--color-gray-800));
}

.splitpanes--horizontal > .splitpanes__splitter {
    min-height: 2px;
    background-color: rgb(var(--color-gray-800));
}
</style>
