<script lang="ts" setup>
import { CancelIcon, ContentSaveIcon, LoadingIcon, PencilIcon } from 'mdi-vue3';
import type { GetMOTDResponse, SetMOTDResponse } from '~~/gen/ts/services/jobs/jobs';

const { $grpc } = useNuxtApp();

const { data, refresh } = useLazyAsyncData('jobs-motd', () => getMOTD());

async function getMOTD(): Promise<GetMOTDResponse> {
    try {
        const call = $grpc.getJobsClient().getMOTD({});
        const { response } = await call;

        return response;
    } catch (e) {
        $grpc.handleError(e as RpcError);
        throw e;
    }
}

async function setMOTD(motd: string): Promise<SetMOTDResponse> {
    try {
        const call = $grpc.getJobsClient().setMOTD({
            motd,
        });
        const { response } = await call;

        return response;
    } catch (e) {
        $grpc.handleError(e as RpcError);
        throw e;
    }
}

const canSubmit = ref(true);
const onSubmit = async (e: string): Promise<SetMOTDResponse> => {
    const response = await setMOTD(e).finally(() => useTimeoutFn(() => (canSubmit.value = true), 400));

    if (data.value !== null) {
        data.value.motd = response.motd;
    }

    return response;
};
const onSubmitThrottle = useThrottleFn(async (e: string) => {
    canSubmit.value = false;
    await onSubmit(e);
}, 1000);

const canEdit = can('JobsService.SetMOTD');

const editing = ref(false);

watch(editing, () => {
    if (!editing.value) {
        refresh();
    }
});
</script>

<template>
    <div v-if="data !== null" class="w-full flex-col">
        <div class="flex items-center">
            <h4 v-if="data.motd.length > 0 || canEdit" class="flex-1 text-base font-semibold leading-6">
                {{ $t('common.motd') }}
            </h4>

            <template v-if="canEdit">
                <UButton v-if="!editing" variant="link" icon="i-mdi-pencil" :loading="!canSubmit" @click="editing = !editing" />
                <div v-else class="flex flex-row gap-1">
                    <UButton
                        variant="link"
                        icon="i-mdi-content-save"
                        :loading="!canSubmit"
                        @click="
                            onSubmitThrottle(data?.motd ?? '');
                            editing = !editing;
                        "
                    />
                    <UButton variant="link" icon="i-mdi-cancel" :loading="!canSubmit" @click="editing = !editing" />
                </div>
            </template>
        </div>

        <div class="flex">
            <template v-if="!editing">
                <div class="w-full flex-1">
                    <p class="prose prose-invert">
                        {{ data.motd }}
                    </p>
                </div>
            </template>
            <template v-else>
                <UTextarea
                    v-model="data.motd"
                    rows="2"
                    name="content"
                    @focusin="focusTablet(true)"
                    @focusout="focusTablet(false)"
                />
            </template>
        </div>
    </div>
</template>
