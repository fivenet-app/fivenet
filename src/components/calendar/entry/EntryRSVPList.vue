<script lang="ts" setup>
import type { ListCalendarEntryRSVPResponse } from '~~/gen/ts/services/calendar/calendar';

const props = defineProps<{
    entryId: string;
}>();

const { $grpc } = useNuxtApp();

const page = ref(1);
const offset = computed(() => (data.value?.pagination?.pageSize ? data.value?.pagination?.pageSize * (page.value - 1) : 0));

const { data, pending: loading, refresh, error } = useLazyAsyncData(`vehicles-${page.value}`, () => listCalendarEntryRSVP());

async function listCalendarEntryRSVP(): Promise<ListCalendarEntryRSVPResponse> {
    try {
        const call = $grpc.getCalendarClient().listCalendarEntryRSVP({
            pagination: {
                offset: offset.value,
            },
            entryId: props.entryId,
        });
        const { response } = await call;

        return response;
    } catch (e) {
        $grpc.handleError(e as RpcError);
        throw e;
    }
}
// TODO
</script>

<template>
    <div>
        <!-- TODO -->
    </div>
</template>
