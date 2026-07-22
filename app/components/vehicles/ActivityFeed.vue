<script lang="ts" setup>
import type { Form } from '@nuxt/ui';
import { z } from 'zod';
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import DataNoDataBlock from '~/components/partials/data/DataNoDataBlock.vue';
import DataPendingBlock from '~/components/partials/data/DataPendingBlock.vue';
import Pagination from '~/components/partials/Pagination.vue';
import SortButton from '~/components/partials/SortButton.vue';
import { getVehiclesVehiclesClient } from '~~/gen/ts/clients';
import type { SortByColumn } from '~~/gen/ts/resources/common/database/database';
import { VehicleActivityType } from '~~/gen/ts/resources/vehicles/activity/activity';
import type { ListVehicleActivityResponse } from '~~/gen/ts/services/vehicles/vehicles';
import ActivityFeedEntry from './ActivityFeedEntry.vue';
import { vehicleActivityTypeBGColor, vehicleActivityTypeIcon } from './helpers';

const props = defineProps<{
    plate: string;
    ownerId?: number;
}>();

const { t } = useI18n();
const { attr, activeChar } = useAuth();

const vehiclesVehiclesClientPromise = getVehiclesVehiclesClient();

const activityTypes = Object.keys(VehicleActivityType)
    .map((t) => VehicleActivityType[t as keyof typeof VehicleActivityType])
    .filter((at) => typeof at === 'number' && at > VehicleActivityType.UNSPECIFIED);

const options = activityTypes.map((at) => ({
    label: t(`enums.vehicles.VehicleActivityType.${VehicleActivityType[at]}`),
    icon: vehicleActivityTypeIcon(at),
    value: at,
    ui: {
        itemLeadingIcon: vehicleActivityTypeBGColor(at),
    },
}));

const schema = z.object({
    types: z.enum(VehicleActivityType).array().max(activityTypes.length).default(activityTypes),
    sorting: z
        .object({
            columns: z
                .custom<SortByColumn>()
                .array()
                .max(3)
                .default([
                    {
                        id: 'createdAt',
                        desc: true,
                    },
                ]),
        })
        .default({ columns: [{ id: 'createdAt', desc: true }] }),
    page: pageNumberSchema,
});

type Schema = z.output<typeof schema>;

const query = useSearchForm(`vehicle_activity_${props.plate}`, schema);
const formRef = useTemplateRef<Form<typeof schema>>('formRef');
const { validatedQuery, commitValidatedQuery } = useFormSearchValidation<typeof schema>(query, formRef);

const activityKey = computed(() => `vehicle-activity-${props.plate}-${JSON.stringify(validatedQuery.value)}`);
const { data, status, refresh, error } = useLazyAsyncData(activityKey, () => listVehicleActivity(validatedQuery.value), {
    immediate: false,
});

defineExpose({
    refresh,
});

const denyView = computed(
    () =>
        props.ownerId !== undefined &&
        props.ownerId === activeChar.value?.userId &&
        !attr('vehicles.VehiclesService/ListVehicleActivity', 'Fields', 'Own').value,
);

watch(
    denyView,
    (value) => {
        if (!value) {
            refresh();
        }
    },
    { immediate: true },
);

async function listVehicleActivity(values: Schema): Promise<ListVehicleActivityResponse> {
    try {
        const vehiclesVehiclesClient = await vehiclesVehiclesClientPromise;
        const call = vehiclesVehiclesClient.listVehicleActivity({
            pagination: {
                offset: calculateOffset(values.page, data.value?.pagination),
            },
            sort: values.sorting,
            plate: props.plate,
            types: values.types,
        });
        const { response } = await call;

        return response;
    } catch (e) {
        handleGRPCError(e as RpcError);
        throw e;
    }
}
</script>

<template>
    <UDashboardPanel :ui="{ root: 'min-h-0 h-full', body: 'p-0 sm:p-0 gap-0 sm:gap-0' }">
        <template v-if="!denyView" #header>
            <UDashboardToolbar>
                <template #default>
                    <UForm
                        ref="formRef"
                        class="my-2 flex w-full flex-row gap-2"
                        :schema="schema"
                        :state="query"
                        @submit="commitValidatedQuery"
                    >
                        <UFormField class="flex-1 grow" name="types" :label="$t('common.type', 2)">
                            <ClientOnly>
                                <USelectMenu
                                    v-model="query.types"
                                    class="w-full min-w-40 flex-1"
                                    multiple
                                    :items="options"
                                    value-key="value"
                                    :search-input="{ placeholder: $t('common.type', 2) }"
                                >
                                    <template #default>
                                        {{ $t('common.selected', query.types.length) }}
                                    </template>

                                    <template #empty> {{ $t('common.not_found', [$t('common.type', 2)]) }} </template>
                                </USelectMenu>
                            </ClientOnly>
                        </UFormField>

                        <UFormField label="&nbsp;">
                            <SortButton
                                v-model="query.sorting"
                                :fields="[{ label: $t('common.created_at'), value: 'createdAt' }]"
                            />
                        </UFormField>
                    </UForm>
                </template>
            </UDashboardToolbar>
        </template>

        <template #body>
            <UContainer v-if="denyView" class="my-2">
                <UAlert
                    variant="subtle"
                    color="error"
                    icon="i-mdi-denied"
                    :title="$t('components.vehicles.VehicleActivityFeed.own.title')"
                    :description="$t('components.vehicles.VehicleActivityFeed.own.message')"
                />
            </UContainer>

            <DataPendingBlock
                v-else-if="isRequestPending(status)"
                :message="$t('common.loading', [`${$t('common.vehicle', 1)} ${$t('common.activity')}`])"
            />
            <DataErrorBlock
                v-else-if="error"
                :title="$t('common.not_found', [`${$t('common.vehicle', 1)} ${$t('common.activity')}`])"
                :error="error"
                :retry="refresh"
            />
            <DataNoDataBlock
                v-else-if="!data || data.activity.length === 0"
                :type="`${$t('common.vehicle', 1)} ${$t('common.activity')}`"
                icon="i-mdi-pulse"
            />

            <div v-else class="relative flex-1">
                <ul class="min-w-full divide-y divide-default" role="list">
                    <ActivityFeedEntry v-for="activity in data.activity" :key="activity.id" :activity="activity" />
                </ul>
            </div>
        </template>

        <template #footer>
            <Pagination v-model="query.page" :pagination="data?.pagination" :status="status" :refresh="refresh" />
        </template>
    </UDashboardPanel>
</template>
