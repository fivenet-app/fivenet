<script lang="ts" setup>
import { RpcError } from '@protobuf-ts/runtime-rpc';
import { useThrottleFn } from '@vueuse/core';
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
    const response = await setMOTD(e).finally(() => setTimeout(() => (canSubmit.value = true), 400));

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
            <h4 v-if="data.motd.length > 0 || canEdit" class="flex-1 mt-2 text-base font-semibold leading-6 text-neutral">
                {{ $t('common.motd') }}
            </h4>

            <template v-if="canEdit">
                <button
                    v-if="!editing"
                    type="button"
                    class="text-primary-500 hover:text-primary-400"
                    @click="editing = !editing"
                >
                    <PencilIcon class="h-5 w-5" aria-hidden="true" />
                </button>
                <div v-else class="flex flex-row gap-1">
                    <button
                        type="button"
                        class="inline-flex flex-row text-primary-500 hover:text-primary-400"
                        @click="
                            onSubmitThrottle(data?.motd ?? '');
                            editing = !editing;
                        "
                    >
                        <ContentSaveIcon class="h-5 w-5" aria-hidden="true" />
                        <template v-if="!canSubmit">
                            <LoadingIcon class="mr-2 h-5 w-5 animate-spin" aria-hidden="true" />
                        </template>
                    </button>
                    <button type="button" class="text-primary-500 hover:text-primary-400" @click="editing = !editing">
                        <CancelIcon class="h-5 w-5" aria-hidden="true" />
                    </button>
                </div>
            </template>
        </div>

        <div class="flex">
            <template v-if="!editing">
                <div class="flex-1 w-full">
                    <p class="prose prose-invert">
                        {{ data.motd }}
                    </p>
                </div>
            </template>
            <template v-else>
                <textarea
                    v-model="data.motd"
                    rows="2"
                    name="content"
                    class="flex-1 w-full rounded-md border-0 bg-base-700 py-1.5 text-neutral placeholder:text-accent-200 focus:ring-2 focus:ring-inset border-2 border-base-200 focus:ring-base-300 sm:text-sm sm:leading-6"
                    @focusin="focusTablet(true)"
                    @focusout="focusTablet(false)"
                ></textarea>
            </template>
        </div>
    </div>
</template>
