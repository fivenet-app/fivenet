<script lang="ts" setup>
import type { FormSubmitEvent } from '@nuxt/ui';
import { z } from 'zod';
import EditorToolbar from '~/components/fabriceditor/EditorToolbar.vue';
import EditorWrapper from '~/components/fabriceditor/EditorWrapper.vue';
import { getDocumentsStampsClient } from '~~/gen/ts/clients';

useHead({
    title: 'pages.documents.stamps.create',
});

definePageMeta({
    title: 'pages.documents.stamps.create',
    requiresAuth: true,
    permission: 'documents.StampsService/UpsertStamp',
});

const { can } = useAuth();

const schema = z.object({
    name: z.string().min(1).max(120),
    svgData: z.string().min(1).max(99999),
});

type Schema = z.output<typeof schema>;

const state = reactive<Schema>({
    name: '',
    svgData: '',
});

const stampsClient = await getDocumentsStampsClient();

async function createOrUpsertStamp(values: Schema) {
    try {
        const call = stampsClient.upsertStamp({
            stamp: {
                id: 0,
                job: '',
                name: '',
                svgTemplate: values.svgData,
                access: {
                    jobs: [],
                },
            },
        });
        const { response } = await call;

        return response;
    } catch (e) {
        handleGRPCError(e);
        throw e;
    }
}

const canSubmit = ref(true);
const onSubmitThrottle = useThrottleFn(async (event: FormSubmitEvent<Schema>) => {
    canSubmit.value = false;
    await createOrUpsertStamp(event.data).finally(() => useTimeoutFn(() => (canSubmit.value = true), 400));
}, 1000);
</script>

<template>
    <UForm :state="state" :schema="schema" class="flex w-full flex-1" @submit="onSubmitThrottle">
        <UDashboardPanel :ui="{ body: 'p-0 sm:p-0 gap-0 sm:gap-0' }">
            <template #header>
                <UDashboardNavbar :title="$t('pages.documents.stamps.create')">
                    <template #leading>
                        <UDashboardSidebarCollapse />
                    </template>

                    <template #right>
                        <PartialsBackButton fallback-to="/documents" />

                        <UTooltip v-if="can('documents.StampsService/UpsertStamp').value" :text="$t('common.save', 1)">
                            <UButton trailing-icon="i-mdi-content-save" color="neutral" variant="outline" truncate>
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
                <EditorWrapper v-model="state.svgData" :max-width="900" :max-height="350" background-color="#ffffff">
                    <template #sidebar-top>
                        <UCard :ui="{ header: 'font-bold' }">
                            <template #header>
                                {{ $t('pages.documents.stamps.create') }}
                            </template>

                            <UFormField name="name" :label="$t('common.name')">
                                <UInput v-model="state.name" type="text" class="w-full" />
                            </UFormField>
                        </UCard>
                    </template>
                </EditorWrapper>
            </template>
        </UDashboardPanel>
    </UForm>
</template>
