<script lang="ts" setup>
import CitizenInfoPopover from '~/components/partials/citizens/CitizenInfoPopover.vue';
import ConfirmModal from '~/components/partials/ConfirmModal.vue';
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import DataNoDataBlock from '~/components/partials/data/DataNoDataBlock.vue';
import DataPendingBlock from '~/components/partials/data/DataPendingBlock.vue';
import { useAuthStore } from '~/store/auth';
import { useCalendarStore } from '~/store/calendar';
import { CalendarEntryRSVP, RsvpResponses } from '~~/gen/ts/resources/calendar/calendar';
import type { ListCalendarEntryRSVPResponse, RSVPCalendarEntryResponse } from '~~/gen/ts/services/calendar/calendar';
import EntryShareForm from './EntryShareForm.vue';

const props = withDefaults(
    defineProps<{
        modelValue: CalendarEntryRSVP | undefined;
        entryId: string;
        rsvpOpen?: boolean;
        disabled?: boolean;
        showRemove?: boolean;
        canShare?: boolean;
    }>(),
    {
        showRemove: true,
        canShare: false,
    },
);

const emits = defineEmits<{
    (e: 'update:modelValue', entry: CalendarEntryRSVP | undefined): void;
}>();

const modal = useModal();

const authStore = useAuthStore();
const { activeChar } = storeToRefs(authStore);

const calendarStore = useCalendarStore();

const ownEntry = useVModel(props, 'modelValue', emits);

const page = ref(1);
const offset = computed(() => (data.value?.pagination?.pageSize ? data.value?.pagination?.pageSize * (page.value - 1) : 0));

const {
    data,
    pending: loading,
    refresh,
    error,
} = useLazyAsyncData(`calendar-entry:${props.entryId}-${page.value}`, () => listCalendarEntryRSVP());

async function listCalendarEntryRSVP(): Promise<ListCalendarEntryRSVPResponse> {
    try {
        const response = await calendarStore.listCalendarEntryRSVP({
            pagination: {
                offset: offset.value,
            },
            entryId: props.entryId,
        });

        return response;
    } catch (e) {
        handleGRPCError(e as RpcError);
        throw e;
    }
}

async function rsvpCalendarEntry(rsvpResponse: RsvpResponses, remove?: boolean): Promise<void | RSVPCalendarEntryResponse> {
    if (ownEntry.value?.response === rsvpResponse) {
        return;
    }

    try {
        const response = await calendarStore.rsvpCalendarEntry({
            entry: {
                entryId: props.entryId,
                response: rsvpResponse,
                userId: activeChar.value?.userId!,
            },
            subscribe: true,
            remove: remove,
        });

        if (response.entry) {
            const idx = data.value!.entries.findIndex(
                (e) => e.entryId === response.entry?.entryId && e.userId === response.entry?.userId,
            );
            if (idx > -1) {
                data.value!.entries[idx] = response.entry;
            } else {
                data.value!.entries.push(response.entry);
            }
        }
        ownEntry.value = response.entry;

        return response;
    } catch (e) {
        handleGRPCError(e as RpcError);
        throw e;
    }
}

const groupedEntries = computed(() => ({
    yes: data.value?.entries.filter((e) => e.response === RsvpResponses.YES),
    maybe: data.value?.entries.filter((e) => e.response === RsvpResponses.MAYBE),
    no: data.value?.entries.filter((e) => e.response === RsvpResponses.NO),
    invited: data.value?.entries.filter((e) => e.response === RsvpResponses.INVITED),
}));

const openShare = ref(false);

const canSubmit = ref(true);
const onSubmitThrottle = useThrottleFn(async (rsvpResponse: RsvpResponses) => {
    canSubmit.value = false;
    await rsvpCalendarEntry(rsvpResponse).finally(() => useTimeoutFn(() => (canSubmit.value = true), 400));
}, 1000);
</script>

<template>
    <div>
        <div class="mt-2 flex gap-2">
            <UButtonGroup v-if="rsvpOpen" class="inline-flex w-full">
                <UButton
                    block
                    class="flex-1"
                    :disabled="!canSubmit || disabled"
                    :loading="!canSubmit"
                    color="green"
                    :variant="ownEntry?.response === RsvpResponses.YES ? 'soft' : 'solid'"
                    @click="onSubmitThrottle(RsvpResponses.YES)"
                >
                    {{ $t('common.yes') }}
                </UButton>

                <UButton
                    block
                    class="flex-1"
                    :disabled="!canSubmit || disabled"
                    :loading="!canSubmit"
                    color="amber"
                    :variant="ownEntry?.response === RsvpResponses.MAYBE ? 'soft' : 'solid'"
                    @click="onSubmitThrottle(RsvpResponses.MAYBE)"
                >
                    {{ $t('common.maybe') }}
                </UButton>

                <UButton
                    block
                    class="flex-1"
                    :disabled="!canSubmit || disabled"
                    :loading="!canSubmit"
                    color="red"
                    :variant="ownEntry?.response === RsvpResponses.NO ? 'soft' : 'solid'"
                    @click="onSubmitThrottle(RsvpResponses.NO)"
                >
                    {{ $t('common.no') }}
                </UButton>
            </UButtonGroup>

            <UButtonGroup class="inline-flex">
                <UButton
                    v-if="ownEntry && showRemove"
                    icon="i-mdi-calendar-remove"
                    color="white"
                    @click="
                        modal.open(ConfirmModal, {
                            confirm: async () => rsvpCalendarEntry(RsvpResponses.NO, true),
                        })
                    "
                />

                <UButton v-if="canShare" :icon="!openShare ? 'i-mdi-invite' : 'i-mdi-close'" @click="openShare = !openShare" />
            </UButtonGroup>
        </div>

        <EntryShareForm v-if="canShare && openShare" :entry-id="entryId" @close="openShare = false" @refresh="refresh()" />

        <UAccordion
            variant="ghost"
            :items="[{ slot: 'rsvp', label: $t('common.rsvp'), icon: 'i-mdi-calendar-question' }]"
            class="mt-2 flex flex-col"
        >
            <template #rsvp>
                <UContainer>
                    <DataPendingBlock v-if="loading" :message="$t('common.loading', [$t('common.entry', 1)])" />
                    <DataErrorBlock
                        v-else-if="error"
                        :title="$t('common.unable_to_load', [$t('common.entry', 1)])"
                        :retry="refresh"
                    />
                    <DataNoDataBlock v-else-if="!data" :type="$t('common.entry', 1)" icon="i-mdi-calendar" />

                    <div v-else class="flex flex-col gap-2">
                        <template v-if="data.entries.length === 0">
                            <p>{{ $t('common.none', [$t('common.response', 2)]) }}</p>
                        </template>

                        <template v-else>
                            <template v-for="(rsvp, key) in groupedEntries" :key="key">
                                <div v-if="!rsvp || rsvp?.length > 0">
                                    <h3 class="font-bold text-black dark:text-white">{{ $t(`common.${key}`) }}</h3>
                                    <div class="grid grid-cols-2 gap-2 lg:grid-cols-4">
                                        <CitizenInfoPopover v-for="entry in rsvp" :user="entry.user" />
                                    </div>
                                </div>
                            </template>
                        </template>
                    </div>
                </UContainer>
            </template>
        </UAccordion>
    </div>
</template>
