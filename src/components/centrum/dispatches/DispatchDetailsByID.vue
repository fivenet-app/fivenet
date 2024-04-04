<script lang="ts" setup>
import type { GetDispatchResponse } from '~~/gen/ts/services/centrum/centrum';
import DispatchDetails from '~/components/centrum/dispatches/DispatchDetails.vue';

const props = defineProps<{
    dispatchId: string;
}>();

const { $grpc } = useNuxtApp();

const { isOpen } = useModal();

const { data, refresh } = useLazyAsyncData(`centrum-dispatch-${props.dispatchId}`, () => getDispatch(props.dispatchId));

async function getDispatch(id: string): Promise<GetDispatchResponse> {
    try {
        const call = $grpc.getCentrumClient().getDispatch({ id });
        const { response } = await call;

        // TODO show notification when no dispatch

        return response;
    } catch (e) {
        $grpc.handleError(e as RpcError);
        isOpen.value = false;
        throw e;
    }
}

watch(props, () => refresh());
</script>

<template>
    <DispatchDetails v-if="data?.dispatch" :dispatch="data.dispatch" @goto="$emit('goto', $event)" />
</template>
