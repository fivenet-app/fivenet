<script lang="ts" setup>
import type { GetDispatchResponse } from '~~/gen/ts/services/centrum/centrum';
import DispatchDetailsSlideover from '~/components/centrum/dispatches/DispatchDetailsSlideover.vue';

const props = defineProps<{
    dispatchId: string;
}>();

defineEmits<{
    (e: 'goto', loc: Coordinate): void;
}>();

const { $grpc } = useNuxtApp();

const { isOpen } = useSlideover();

const { data, refresh } = useLazyAsyncData(`centrum-dispatch-${props.dispatchId}`, () => getDispatch(props.dispatchId));

async function getDispatch(id: string): Promise<GetDispatchResponse> {
    try {
        const call = $grpc.getCentrumClient().getDispatch({ id });
        const { response } = await call;

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
    <DispatchDetailsSlideover v-if="data?.dispatch" :dispatch="data.dispatch" @goto="$emit('goto', $event)" />
</template>
