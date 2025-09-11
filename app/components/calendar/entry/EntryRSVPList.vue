<script lang="ts" setup>
import CitizenInfoPopover from '~/components/partials/citizens/CitizenInfoPopover.vue';
import ConfirmModal from '~/components/partials/ConfirmModal.vue';
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import DataNoDataBlock from '~/components/partials/data/DataNoDataBlock.vue';
import DataPendingBlock from '~/components/partials/data/DataPendingBlock.vue';
import { useAuthStore } from '~/stores/auth';
import { useCalendarStore } from '~/stores/calendar';
import { type CalendarEntryRSVP, RsvpResponses } from '~~/gen/ts/resources/calendar/calendar';
import type { ListCalendarEntryRSVPResponse, RSVPCalendarEntryResponse } from '~~/gen/ts/services/calendar/calendar';
import EntryShareForm from './EntryShareForm.vue';

const props = withDefaults(
    defineProps<{
        modelValue: CalendarEntryRSVP | undefined;
        entryId: number;
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

const emit = defineEmits<{
    (e: 'update:modelValue', entry: CalendarEntryRSVP | undefined): void;
}>();

const overlay = useOverlay();

const authStore = useAuthStore();
const { activeChar } = storeToRefs(authStore);

const calendarStore = useCalendarStore();

const ownEntry = useVModel(props, 'modelValue', emit);

const page = useRouteQuery('page', '1', { transform: Number });

const { data, status, refresh, error } = useLazyAsyncData(`calendar-entry:${props.entryId}-${page.value}`, () =>
    listCalendarEntryRSVP(),
);

async function listCalendarEntryRSVP(): Promise<ListCalendarEntryRSVPResponse> {
    try {
        const response = await calendarStore.listCalendarEntryRSVP({
            pagination: {
                offset: calculateOffset(page.value, data.value?.pagination),
            },
            entryId: props.entryId,
        });

        return response;
    } catch (e) {
        handleGRPCError(e as RpcError);
        throw e;
    }
}

async function rsvpCalendarEntry(
    rsvpResponse: RsvpResponses,
    remove?: boolean,
): Promise<undefined | RSVPCalendarEntryResponse> {
    if (ownEntry.value?.response === rsvpResponse) {
        return;
    }

    try {
        const response = await calendarStore.rsvpCalendarEntry({
            entry: {
                entryId: props.entryId,
                response: rsvpResponse,
                userId: activeChar.value!.userId!,
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

const confirmModal = overlay.create(ConfirmModal);
</script>

<template>
    <div>
        <div class="mt-2 flex gap-2">
            <UButtonGroup v-if="rsvpOpen" class="inline-flex w-full">
                <UButton
                    class="flex-1"
                    block
                    :disabled="!canSubmit || disabled"
                    :loading="!canSubmit"
                    color="green"
                    :variant="ownEntry?.response === RsvpResponses.YES ? 'soft' : 'solid'"
                    @click="onSubmitThrottle(RsvpResponses.YES)"
                >
                    {{ $t('common.yes') }}
                </UButton>

                <UButton
                    class="flex-1"
                    block
                    :disabled="!canSubmit || disabled"
                    :loading="!canSubmit"
                    color="warning"
                    :variant="ownEntry?.response === RsvpResponses.MAYBE ? 'soft' : 'solid'"
                    @click="onSubmitThrottle(RsvpResponses.MAYBE)"
                >
                    {{ $t('common.maybe') }}
                </UButton>

                <UButton
                    class="flex-1"
                    block
                    :disabled="!canSubmit || disabled"
                    :loading="!canSubmit"
                    color="error"
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
                    color="neutral"
                    @click="
                        confirmModal.open({
                            confirm: async () => rsvpCalendarEntry(RsvpResponses.NO, true),
                        })
                    "
                />

                <UButton v-if="canShare" :icon="!openShare ? 'i-mdi-invite' : 'i-mdi-close'" @click="openShare = !openShare" />
            </UButtonGroup>
        </div>

        <EntryShareForm v-if="canShare && openShare" :entry-id="entryId" @close="openShare = false" @refresh="refresh()" />

        <UCollapsible class="my-2">
            <UButton
                class="group flex flex-col gap-2"
                color="neutral"
                variant="ghost"
                icon="i-mdi-calendar-question"
                trailing-icon="i-mdi-chevron-down"
                :label="$t('common.rsvp')"
                :ui="{
                    trailingIcon: 'group-data-[state=open]:rotate-180 transition-transform duration-200',
                }"
                block
            />

            <template #content>
                <DataPendingBlock v-if="isRequestPending(status)" :message="$t('common.loading', [$t('common.entry', 1)])" />
                <DataErrorBlock
                    v-else-if="error"
                    :title="$t('common.unable_to_load', [$t('common.entry', 1)])"
                    :error="error"
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
                                    <CitizenInfoPopover v-for="entry in rsvp" :key="entry.userId" :user="entry.user" />
                                </div>
                            </div>
                        </template>
                    </template>
                </div>
            </template>
        </UCollapsible>
    </div>
</template>
