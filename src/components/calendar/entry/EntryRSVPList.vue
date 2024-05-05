<script lang="ts" setup>
import CitizenInfoPopover from '~/components/partials/citizens/CitizenInfoPopover.vue';
import { useAuthStore } from '~/store/auth';
import { RsvpResponses } from '~~/gen/ts/resources/calendar/calendar';
import type { ListCalendarEntryRSVPResponse, RSVPCalendarEntryResponse } from '~~/gen/ts/services/calendar/calendar';

const props = defineProps<{
    entryId: string;
    rsvpOpen?: boolean;
}>();

const { $grpc } = useNuxtApp();

const authStore = useAuthStore();
const { activeChar } = storeToRefs(authStore);

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

async function rsvpCalendarEntry(rsvpResponse: RsvpResponses): Promise<void | RSVPCalendarEntryResponse> {
    if (data.value?.ownEntry?.response === rsvpResponse) {
        return;
    }

    try {
        const call = $grpc.getCalendarClient().rSVPCalendarEntry({
            entry: {
                entryId: props.entryId,
                response: rsvpResponse,
                userId: activeChar.value?.userId!,
            },
            subscribe: true,
        });
        const { response } = await call;

        data.value!.ownEntry = response.entry;

        return response;
    } catch (e) {
        $grpc.handleError(e as RpcError);
        throw e;
    }
}

const groupedEntries = computed(() => ({
    yes: data.value?.entries.filter((e) => e.response === RsvpResponses.YES),
    maybe: data.value?.entries.filter((e) => e.response === RsvpResponses.MAYBE),
    no: data.value?.entries.filter((e) => e.response === RsvpResponses.NO),
}));

const canSubmit = ref(true);
const onSubmitThrottle = useThrottleFn(async (rsvpResponse: RsvpResponses) => {
    canSubmit.value = false;
    await rsvpCalendarEntry(rsvpResponse).finally(() => useTimeoutFn(() => (canSubmit.value = true), 400));
}, 1000);
</script>

<template>
    <div>
        <div class="flex flex-col">
            <template v-if="!data?.entries || !data?.entries.length">
                {{ $t('common.not_found', [$t('common.rsvp')]) }}
            </template>
            <template v-else>
                <div class="mb-2 inline-flex items-center gap-2">
                    <UAvatarGroup size="sm" :max="3">
                        <UAvatar
                            v-for="rsvp in data?.entries.slice(0, 3)"
                            :src="rsvp.user?.avatar?.url"
                            :alt="`${rsvp.user?.firstname} ${rsvp.user?.lastname}`"
                        />
                    </UAvatarGroup>
                    <p v-if="data?.entries.length > 3">...</p>
                </div>

                <UAccordion
                    variant="ghost"
                    :items="[{ slot: 'rsvp', label: $t('common.rsvp'), icon: 'i-mdi-calendar-question' }]"
                >
                    <template #rsvp>
                        <UContainer>
                            <div class="flex flex-col gap-2">
                                <div v-for="(rsvp, key) in groupedEntries" :key="key">
                                    <h3 class="font-bold text-black dark:text-white">{{ $t(`common.${key}`) }}</h3>
                                    <div class="grid grid-cols-2 gap-2 lg:grid-cols-4">
                                        <template v-if="!rsvp?.length">
                                            {{ $t('common.none') }}
                                        </template>
                                        <template v-else>
                                            <CitizenInfoPopover v-for="entry in rsvp" :user="entry.user" />
                                        </template>
                                    </div>
                                </div>
                            </div>
                        </UContainer>
                    </template>
                </UAccordion>
            </template>
        </div>

        <template v-if="rsvpOpen">
            <UDivider />

            <UButtonGroup class="inline-flex w-full">
                <UButton
                    block
                    class="flex-1"
                    :disabled="!canSubmit"
                    :loading="!canSubmit"
                    color="green"
                    :variant="data?.ownEntry?.response === RsvpResponses.YES ? 'soft' : 'solid'"
                    @click="onSubmitThrottle(RsvpResponses.YES)"
                >
                    {{ $t('common.yes') }}
                </UButton>
                <UButton
                    block
                    class="flex-1"
                    :disabled="!canSubmit"
                    :loading="!canSubmit"
                    color="amber"
                    :variant="data?.ownEntry?.response === RsvpResponses.MAYBE ? 'soft' : 'solid'"
                    @click="onSubmitThrottle(RsvpResponses.MAYBE)"
                >
                    {{ $t('common.maybe') }}
                </UButton>
                <UButton
                    block
                    class="flex-1"
                    :disabled="!canSubmit"
                    :loading="!canSubmit"
                    color="red"
                    :variant="data?.ownEntry?.response === RsvpResponses.NO ? 'soft' : 'solid'"
                    @click="onSubmitThrottle(RsvpResponses.NO)"
                >
                    {{ $t('common.no') }}
                </UButton>
            </UButtonGroup>
        </template>
    </div>
</template>
