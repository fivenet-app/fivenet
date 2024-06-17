<script lang="ts" setup>
import type { GetDispatchResponse } from '~~/gen/ts/services/centrum/centrum';
import DispatchDetailsSlideover from '~/components/centrum/dispatches/DispatchDetailsSlideover.vue';
import { useCentrumStore } from '~/store/centrum';

const props = defineProps<{
    dispatchId: string;
}>();

const centrumStore = useCentrumStore();
const { dispatches } = storeToRefs(centrumStore);

const { isOpen } = useSlideover();

const { data, refresh } = useLazyAsyncData(`centrum-dispatch-${props.dispatchId}`, () => getDispatch(props.dispatchId));

async function getDispatch(id: string): Promise<GetDispatchResponse> {
    if (dispatches.value.has(id)) {
        return {
            dispatch: dispatches.value.get(id),
        };
    }

    try {
        const call = getGRPCCentrumClient().getDispatch({ id });
        const { response } = await call;

        return response;
    } catch (e) {
        handleGRPCError(e as RpcError);
        isOpen.value = false;
        throw e;
    }
}

watch(props, async () => refresh());
</script>

<template>
    <DispatchDetailsSlideover v-if="data?.dispatch" :dispatch-id="dispatchId" :dispatch="data.dispatch" />
</template>
