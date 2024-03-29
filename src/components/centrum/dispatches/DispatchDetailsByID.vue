<script lang="ts" setup>
import type { GetDispatchResponse } from '~~/gen/ts/services/centrum/centrum';
import DispatchDetails from '~/components/centrum/dispatches/DispatchDetails.vue';

const props = defineProps<{
    open: boolean;
    dispatchId: string;
}>();

const emit = defineEmits<{
    (e: 'close'): void;
    (e: 'goto', loc: Coordinate): void;
}>();

const { $grpc } = useNuxtApp();

const { data, refresh } = useLazyAsyncData(`centrum-dispatch-${props.dispatchId}`, () => getDispatch(props.dispatchId));

async function getDispatch(id: string): Promise<GetDispatchResponse> {
    try {
        const call = $grpc.getCentrumClient().getDispatch({ id });
        const { response } = await call;

        // TODO show notification when no dispatch

        return response;
    } catch (e) {
        $grpc.handleError(e as RpcError);
        emit('close');
        throw e;
    }
}

watch(props, () => refresh());
</script>

<template>
    <DispatchDetails
        v-if="data?.dispatch"
        :open="open"
        :dispatch="data.dispatch"
        @close="$emit('close')"
        @goto="$emit('goto', $event)"
    />
</template>
