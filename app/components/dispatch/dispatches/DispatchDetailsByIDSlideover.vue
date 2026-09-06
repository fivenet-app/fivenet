<script lang="ts" setup>
import DispatchDetailsSlideover from '~/components/dispatch/dispatches/DispatchDetailsSlideover.vue';
import { useCentrumStore } from '~/stores/centrum';
import { getCentrumDispatchesClient } from '~~/gen/ts/clients';
import type { GetDispatchResponse } from '~~/gen/ts/services/centrum/dispatches';

const props = defineProps<{
    dispatchId: number;
}>();

const emit = defineEmits<{
    (e: 'close', v: boolean): void;
}>();

const centrumStore = useCentrumStore();
const { dispatches } = storeToRefs(centrumStore);

const centrumDispatchesClient = await getCentrumDispatchesClient();

const { data, refresh } = useAuthedLazyAsyncData('userState', `centrum-dispatch-${props.dispatchId}`, ({ signal }) =>
    getDispatch(props.dispatchId, signal),
);

async function getDispatch(id: number, signal: AbortSignal): Promise<GetDispatchResponse> {
    if (dispatches.value.has(id)) {
        return {
            dispatch: dispatches.value.get(id),
        };
    }

    try {
        const call = centrumDispatchesClient.getDispatch({ id }, { abort: signal });
        const { response } = await call;

        return response;
    } catch (e) {
        handleGRPCError(e as RpcError);
        emit('close', false);
        throw e;
    }
}

watch(props, async () => refresh());
</script>

<template>
    <DispatchDetailsSlideover
        v-if="data?.dispatch"
        :dispatch-id="dispatchId"
        :dispatch="data.dispatch"
        @close="$emit('close', false)"
    />
</template>
