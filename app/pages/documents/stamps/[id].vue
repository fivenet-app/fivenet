<script lang="ts" setup>
import type { FormSubmitEvent } from '@nuxt/ui';
import type { TypedRouteFromName } from '@typed-router';
import { z } from 'zod';
import EditorToolbar from '~/components/fabriceditor/EditorToolbar.vue';
import EditorWrapper from '~/components/fabriceditor/EditorWrapper.vue';
import AccessManager from '~/components/partials/access/AccessManager.vue';
import { enumToAccessLevelEnums } from '~/components/partials/access/helpers';
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import DataNoDataBlock from '~/components/partials/data/DataNoDataBlock.vue';
import DataPendingBlock from '~/components/partials/data/DataPendingBlock.vue';
import { getDocumentsStampsClient } from '~~/gen/ts/clients';
import { StampAccessLevel, type StampJobAccess } from '~~/gen/ts/resources/documents/stamps/stamp';

useHead({
    title: 'pages.documents.stamps.update',
});

definePageMeta({
    title: 'pages.documents.stamps.update',
    requiresAuth: true,
    permission: 'documents.StampsService/UpsertStamp',
    validate: async (route) => {
        route = route as TypedRouteFromName<'documents-stamps-id'>;
        // Check if the id is made up of digits
        if (typeof route.params.id !== 'string') {
            return false;
        }
        return !!(route.params.id && !isNaN(Number(route.params.id))) && Number(route.params.id) > -1;
    },
});

const { can } = useAuth();

const schema = z.object({
    name: z.string().min(1).max(120),
    svgData: z.string().max(99999),
    access: z.object({
        jobs: z.custom<StampJobAccess>().array().max(5),
    }),
});

type Schema = z.output<typeof schema>;

const state = reactive<Schema>({
    name: '',
    svgData: '',
    access: {
        jobs: [],
    },
});

const stampsClient = await getDocumentsStampsClient();

async function createOrUpsertStamp(values: Schema) {
    try {
        const call = stampsClient.upsertStamp({
            stamp: {
                id: 0,
                job: '',
                name: values.name,
                svgTemplate: values.svgData,
                access: values.access,
            },
        });
        const { response } = await call;

        return response;
    } catch (e) {
        handleGRPCError(e);
        throw e;
    }
}

const route = useRoute('documents-stamps-id');

const {
    data: stamp,
    status,
    error,
    refresh,
} = useLazyAsyncData('stamp', async () => {
    const stampsClient = await getDocumentsStampsClient();
    const { response } = await stampsClient.getStamp({
        id: Number.parseInt(route.params.id as string),
    });
    return response.stamp;
});

function setFromProps(): void {
    if (!stamp.value) return;

    state.name = stamp.value.name || '';
    state.svgData = stamp.value.svgTemplate || '';
}

setFromProps();
watch(stamp, () => setFromProps());

const canSubmit = ref<boolean>(true);
const onSubmitThrottle = useThrottleFn(async (event: FormSubmitEvent<Schema>) => {
    canSubmit.value = false;
    await createOrUpsertStamp(event.data).finally(() => useTimeoutFn(() => (canSubmit.value = true), 400));
}, 1000);
</script>

<template>
    <UForm class="flex w-full flex-1" :state="state" :schema="schema" @submit="onSubmitThrottle">
        <UDashboardPanel :ui="{ body: 'p-0 sm:p-0 gap-0 sm:gap-0' }">
            <template #header>
                <UDashboardNavbar :title="$t('pages.documents.stamps.update')">
                    <template #leading>
                        <UDashboardSidebarCollapse />
                    </template>

                    <template #right>
                        <PartialsBackButton fallback-to="/documents" />

                        <UTooltip v-if="can('documents.StampsService/UpsertStamp').value" :text="$t('common.save', 1)">
                            <UButton
                                type="submit"
                                trailing-icon="i-mdi-content-save"
                                color="neutral"
                                variant="outline"
                                truncate
                            >
                                <span class="hidden truncate sm:block">
                                    {{ $t('common.save', 1) }}
                                </span>
                            </UButton>
                        </UTooltip>
                    </template>
                </UDashboardNavbar>

                <UDashboardToolbar class="p-1 print:hidden">
                    <template #default>
                        <EditorToolbar />
                    </template>
                </UDashboardToolbar>
            </template>

            <template #body>
                <DataPendingBlock v-if="isRequestPending(status)" :message="$t('common.loading', [$t('common.stamp', 1)])" />
                <DataErrorBlock
                    v-else-if="error"
                    :title="$t('common.unable_to_load', [$t('common.stamp', 1)])"
                    :error="error"
                    :retry="refresh"
                />
                <DataNoDataBlock v-else-if="!stamp" :type="$t('common.stamp', 1)" icon="i-mdi-stamper" />

                <EditorWrapper v-else v-model="state.svgData" :max-width="900" :max-height="350" background-color="#ffffff">
                    <template #sidebar-top>
                        <UCard>
                            <template #header>
                                {{ $t('pages.documents.stamps.update') }}
                            </template>

                            <UFormField name="name" :label="$t('common.name')" required>
                                <UInput v-model="state.name" class="w-full" type="text" />
                            </UFormField>

                            <UFormField name="access" :label="$t('common.access')">
                                <AccessManager
                                    v-model:jobs="state.access.jobs"
                                    :target-id="0"
                                    name="jobs"
                                    :access-types="[{ label: $t('common.job', 2), value: 'job' }]"
                                    :access-roles="enumToAccessLevelEnums(StampAccessLevel, 'enums.documents.StampAccessLevel')"
                                />
                            </UFormField>
                        </UCard>
                    </template>
                </EditorWrapper>
            </template>
        </UDashboardPanel>
    </UForm>
</template>
