<script lang="ts" setup>
import CitizenInfoPopover from '~/components/partials/citizens/CitizenInfoPopover.vue';
import ConfirmModal from '~/components/partials/ConfirmModal.vue';
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import DataNoDataBlock from '~/components/partials/data/DataNoDataBlock.vue';
import DataPendingBlock from '~/components/partials/data/DataPendingBlock.vue';
import { useAuthStore } from '~/stores/auth';
import { useCalendarStore } from '~/stores/calendar';
import { type CalendarEntryRSVP, RsvpResponses } from '~~/gen/ts/resources/calendar/entries/entries';
import type { ListCalendarEntryRSVPResponse, RSVPCalendarEntryResponse } from '~~/gen/ts/services/calendar/entries';
import EntryRsvpScopeModal from './EntryRsvpScopeModal.vue';
import EntryShareForm from './EntryShareForm.vue';

const props = withDefaults(
    defineProps<{
        entryId: number;
        occurrenceKey?: string;
        rsvpOpen?: boolean;
        disabled?: boolean;
        showRemove?: boolean;
        canShare?: boolean;
    }>(),
    {
        occurrenceKey: undefined,
        showRemove: true,
        canShare: false,
    },
);

const overlay = useOverlay();

const authStore = useAuthStore();
const { activeChar } = storeToRefs(authStore);

const calendarStore = useCalendarStore();

const ownEntry = defineModel<CalendarEntryRSVP | undefined>();

const page = ref<number>(1);
type RsvpScope = 'series' | 'occurrence';

const pendingRsvpAction = ref<{ response: RsvpResponses; remove: boolean } | null>(null);
const showRsvpScopeModal = ref<boolean>(false);

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
    scope: RsvpScope = 'series',
): Promise<undefined | RSVPCalendarEntryResponse> {
    if (!props.occurrenceKey && ownEntry.value?.response === rsvpResponse) return;

    try {
        const response = await calendarStore.rsvpCalendarEntry({
            entry: {
                entryId: props.entryId,
                response: rsvpResponse,
                userId: activeChar.value!.userId!,
                occurrenceKey: scope === 'occurrence' ? props.occurrenceKey : undefined,
            },
            subscribe: true,
            remove: remove,
            occurrenceKey: scope === 'occurrence' ? props.occurrenceKey : undefined,
        });

        if (response.entry && !response.entry.occurrenceKey) {
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

const submitRsvp = useThrottleFn(async (rsvpResponse: RsvpResponses, remove?: boolean, scope: RsvpScope = 'series') => {
    canSubmit.value = false;
    await rsvpCalendarEntry(rsvpResponse, remove, scope).finally(() => useTimeoutFn(() => (canSubmit.value = true), 400));
}, 1000);

function requestRsvp(rsvpResponse: RsvpResponses, remove = false): void {
    if (props.occurrenceKey) {
        pendingRsvpAction.value = { response: rsvpResponse, remove };
        showRsvpScopeModal.value = true;
        return;
    }

    void submitRsvp(rsvpResponse, remove, 'series');
}

async function chooseRsvpScope(scope: RsvpScope): Promise<void> {
    if (!pendingRsvpAction.value) return;

    const action = pendingRsvpAction.value;
    pendingRsvpAction.value = null;
    showRsvpScopeModal.value = false;

    await submitRsvp(action.response, action.remove, scope);
}

watch(showRsvpScopeModal, (open) => {
    if (!open) {
        pendingRsvpAction.value = null;
    }
});

const groupedEntries = computed(() => ({
    yes: data.value?.entries.filter((e) => e.response === RsvpResponses.YES),
    maybe: data.value?.entries.filter((e) => e.response === RsvpResponses.MAYBE),
    no: data.value?.entries.filter((e) => e.response === RsvpResponses.NO),
    invited: data.value?.entries.filter((e) => e.response === RsvpResponses.INVITED),
}));

const openShare = ref<boolean>(false);

const canSubmit = ref<boolean>(true);

const confirmModal = overlay.create(ConfirmModal);
</script>

<template>
    <div class="flex w-full flex-col gap-2">
        <div class="flex flex-1 flex-row gap-2">
            <UFieldGroup v-if="rsvpOpen" class="inline-flex w-full">
                <UButton
                    class="flex-1"
                    block
                    :disabled="!canSubmit || disabled"
                    :loading="!canSubmit"
                    color="success"
                    :variant="ownEntry?.response === RsvpResponses.YES ? 'soft' : 'solid'"
                    :label="$t('common.yes')"
                    @click="requestRsvp(RsvpResponses.YES)"
                />

                <UButton
                    class="flex-1"
                    block
                    :disabled="!canSubmit || disabled"
                    :loading="!canSubmit"
                    color="warning"
                    :variant="ownEntry?.response === RsvpResponses.MAYBE ? 'soft' : 'solid'"
                    :label="$t('common.maybe')"
                    @click="requestRsvp(RsvpResponses.MAYBE)"
                />

                <UButton
                    class="flex-1"
                    block
                    :disabled="!canSubmit || disabled"
                    :loading="!canSubmit"
                    color="error"
                    :variant="ownEntry?.response === RsvpResponses.NO ? 'soft' : 'solid'"
                    :label="$t('common.no')"
                    @click="requestRsvp(RsvpResponses.NO)"
                />
            </UFieldGroup>

            <UFieldGroup class="inline-flex">
                <UButton
                    v-if="ownEntry && showRemove"
                    icon="i-mdi-calendar-remove"
                    color="neutral"
                    @click="
                        props.occurrenceKey
                            ? requestRsvp(RsvpResponses.NO, true)
                            : confirmModal.open({
                                  confirm: async () => submitRsvp(RsvpResponses.NO, true, 'series'),
                              })
                    "
                />

                <UTooltip v-if="canShare" :text="$t('common.invite')">
                    <UButton :icon="!openShare ? 'i-mdi-invite' : 'i-mdi-close'" />
                </UTooltip>
            </UFieldGroup>
        </div>

        <EntryShareForm
            v-if="canShare && openShare"
            :entry-id="entryId"
            @refresh="refresh()"
            @close="() => (openShare = !openShare)"
        />

        <p v-if="ownEntry?.occurrenceKey" class="mt-1 text-xs text-toned">
            {{ $t('components.calendar.rsvp_scope.occurrence_active') }}
        </p>

        <EntryRsvpScopeModal
            v-model:open="showRsvpScopeModal"
            :remove="pendingRsvpAction?.remove ?? false"
            @confirm="chooseRsvpScope"
        />

        <UCollapsible :ui="{ content: 'p-1' }">
            <UButton
                class="group"
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
