<script lang="ts" setup>
import type { FormSubmitEvent } from '@nuxt/ui';
import { z } from 'zod';
import { getJobsJobsClient } from '~~/gen/ts/clients';
import type { GetMOTDResponse, SetMOTDResponse } from '~~/gen/ts/services/jobs/jobs';

const { can } = useAuth();

const jobsJobsClient = await getJobsJobsClient();

const { data, status, refresh } = useLazyAsyncData('jobs-motd', () => getMOTD());

async function getMOTD(): Promise<GetMOTDResponse> {
    try {
        const call = jobsJobsClient.getMOTD({});
        const { response } = await call;

        return response;
    } catch (e) {
        handleGRPCError(e as RpcError);
        throw e;
    }
}

const schema = z.object({
    motd: z.string().min(0).max(1024),
});

type Schema = z.output<typeof schema>;

const state = reactive<Schema>({
    motd: data.value?.motd ?? '',
});

watch(data, () => {
    if (!data.value) {
        return;
    }

    state.motd = data.value.motd;
});

async function setMOTD(values: Schema): Promise<SetMOTDResponse> {
    try {
        const call = jobsJobsClient.setMOTD({
            motd: values.motd,
        });
        const { response } = await call;

        return response;
    } catch (e) {
        handleGRPCError(e as RpcError);
        throw e;
    }
}

const canDo = can('jobs.JobsService/SetMOTD');

const editing = ref(false);

watch(editing, async () => {
    if (!editing.value) {
        refresh();
    }
});

const canSubmit = ref(true);
const onSubmitThrottle = useThrottleFn(async (event: FormSubmitEvent<Schema>) => {
    canSubmit.value = false;
    await setMOTD(event.data).finally(() => useTimeoutFn(() => (canSubmit.value = true), 400));
    editing.value = !editing.value;
}, 1000);
</script>

<template>
    <UForm class="w-full flex-col" :schema="schema" :state="state" @submit="onSubmitThrottle">
        <div class="flex items-center">
            <h4 v-if="data && (data.motd.length > 0 || canDo)" class="flex-1 text-base leading-6 font-semibold">
                {{ $t('common.motd') }}
            </h4>

            <template v-if="canDo">
                <UTooltip v-if="!editing" :text="$t('common.edit')">
                    <UButton variant="link" icon="i-mdi-pencil" :loading="!canSubmit" @click="editing = !editing" />
                </UTooltip>
                <div v-else class="flex flex-row gap-1">
                    <UTooltip :text="$t('common.save')">
                        <UButton type="submit" variant="link" icon="i-mdi-content-save" :loading="!canSubmit" />
                    </UTooltip>
                    <UTooltip :text="$t('common.cancel')">
                        <UButton variant="link" icon="i-mdi-cancel" :loading="!canSubmit" @click="editing = !editing" />
                    </UTooltip>
                </div>
            </template>
        </div>

        <div class="flex">
            <template v-if="!editing">
                <USkeleton v-if="isRequestPending(status)" class="h-7 w-full" />
                <div v-else class="w-full flex-1">
                    <p class="prose line-clamp-5 max-w-full whitespace-pre-wrap dark:prose-invert">
                        {{ state.motd }}
                    </p>
                </div>
            </template>
            <template v-else>
                <UFormField class="w-full" name="motd">
                    <UTextarea v-model="state.motd" name="motd" :rows="4" :maxrows="8" resize />
                </UFormField>
            </template>
        </div>
    </UForm>
</template>
