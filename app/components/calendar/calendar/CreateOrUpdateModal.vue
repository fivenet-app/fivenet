<script lang="ts" setup>
import type { JSONContent } from '@tiptap/core';
import type { FormSubmitEvent } from '@nuxt/ui';
import { VueDraggable } from 'vue-draggable-plus';
import { z } from 'zod';
import ColorPickerTW from '~/components/partials/ColorPickerTW.vue';
import AccessManager from '~/components/partials/access/AccessManager.vue';
import { enumToAccessLevelEnums } from '~/components/partials/access/helpers';
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import DataNoDataBlock from '~/components/partials/data/DataNoDataBlock.vue';
import DataPendingBlock from '~/components/partials/data/DataPendingBlock.vue';
import DraggableHandle from '~/components/partials/DraggableHandle.vue';
import ReorderButtons from '~/components/partials/ReorderButtons.vue';
import SelectMenu from '~/components/partials/SelectMenu.vue';
import { useCalendarStore } from '~/stores/calendar';
import { isSystemManagedCalendar } from '~/components/calendar/helpers';
import { getSettingsSettingsClient } from '~~/gen/ts/clients';
import { AccessLevel } from '~~/gen/ts/resources/calendar/access/access';
import type { Channel } from '~~/gen/ts/resources/discord/discord';
import { NotificationType } from '~~/gen/ts/resources/notifications/notifications';
import type { CreateCalendarResponse, UpdateCalendarResponse } from '~~/gen/ts/services/calendar/calendar';
import TiptapEditor from '../../partials/editor/TiptapEditor.vue';
import ColorPicker from '~/components/partials/ColorPicker.vue';

const props = defineProps<{
    calendarId?: number;
    systemManaged?: boolean;
}>();

const emit = defineEmits<{
    (e: 'close', v: boolean): void;
}>();

const { t } = useI18n();

const { attr, activeChar, can } = useAuth();

const calendarStore = useCalendarStore();
const { hasPrivateCalendar } = storeToRefs(calendarStore);
const settingsSettingsClient = await getSettingsSettingsClient();
const appConfig = useAppConfig();

const notifications = useNotificationsStore();

const { maxAccessEntries } = useAppConfig();

const canDo = computed(() => ({
    privateCalendar: attr('calendar.CalendarService/CreateCalendar', 'Fields', 'Job').value,
    publicCalendar: attr('calendar.CalendarService/CreateCalendar', 'Fields', 'Public').value,
}));

const isSystemManaged = computed(() => props.systemManaged || isSystemManagedCalendar(data.value?.calendar));
const canEditDiscordReminderSettings = can('settings.SettingsService/SetJobProps');
const showDiscordReminderSection = computed(
    () => appConfig.discord.botEnabled && !isSystemManaged.value && !state.private && !!activeChar.value?.job,
);
const canConfigureDiscordReminders = computed(() => showDiscordReminderSection.value && canEditDiscordReminderSettings.value);

function emptyDiscordReminderStep() {
    return {
        atMinute: 15,
        message: '',
        embed: {
            title: '',
            description: '',
            color: '',
        },
    };
}

function emptyDiscordSettings() {
    return {
        enabled: false,
        channelId: '',
        reminderSteps: [] as ReturnType<typeof emptyDiscordReminderStep>[],
    };
}

const schema = z.object({
    name: z.coerce.string().min(3).max(255),
    description: z.custom<JSONContent | string>().optional(),
    private: z.coerce.boolean(),
    public: z.coerce.boolean(),
    closed: z.coerce.boolean(),
    color: z.coerce.string().max(12),
    access: z.object({
        jobs: jobsAccessEntries(t).max(maxAccessEntries).default([]),
        users: userAccessEntries(t).max(maxAccessEntries).default([]),
    }),
    discordSettings: z.object({
        enabled: z.coerce.boolean(),
        channelId: z.coerce.string().max(64),
        reminderSteps: z
            .object({
                atMinute: z.coerce.number().int().min(0).max(10080),
                message: z.coerce.string().max(2000),
                embed: z.object({
                    title: z.coerce.string().max(256),
                    description: z.coerce.string().max(4096),
                    color: z.union([z.literal(''), z.coerce.string().regex(/^#[A-Fa-f0-9]{6}$/)]),
                }),
            })
            .array()
            .max(2)
            .default([]),
    }),
});

type Schema = z.output<typeof schema>;

const state = reactive<Schema>({
    name: '',
    description: '',
    private: !hasPrivateCalendar.value,
    public: false,
    closed: false,
    color: 'blue',
    access: {
        jobs: [],
        users: [],
    },
    discordSettings: emptyDiscordSettings(),
});

const { hasUnsavedChanges, confirmLeave, syncSnapshot } = useSnapshotChanges(state);
const { moveUp, moveDown } = useListReorder(toRef(() => state.discordSettings.reminderSteps));

const {
    data: data,
    status,
    refresh,
    error,
} = useLazyAsyncData(
    `calendar-calendar:${props.calendarId}`,
    () => calendarStore.getCalendar({ calendarId: props.calendarId! }),
    {
        immediate: !!props.calendarId,
    },
);

async function createOrUpdateCalendar(values: Schema): Promise<CreateCalendarResponse | UpdateCalendarResponse> {
    values.access.users.forEach((user) => {
        if (user.id < 0) user.id = 0;
        user.user = undefined; // Clear user object to avoid sending unnecessary data
    });
    values.access.jobs.forEach((job) => job.id < 0 && (job.id = 0));

    try {
        const response = await calendarStore.createOrUpdateCalendar({
            id: data.value?.calendar?.id ?? 0,
            job: isSystemManaged.value
                ? (data.value?.calendar?.job ?? activeChar.value?.job)
                : values.private
                  ? undefined
                  : activeChar.value?.job,
            name: values.name,
            description: tiptapToContent(values.description),
            public: values.public,
            closed: values.closed,
            color: values.color,
            access: values.access,
            creatorJob: '',
            discordSettings: values.private
                ? undefined
                : canEditDiscordReminderSettings.value
                  ? toProtoDiscordSettings(values.discordSettings)
                  : data.value?.calendar?.discordSettings,
        });

        notifications.add({
            title: { key: 'notifications.action_successful.title', parameters: {} },
            description: { key: 'notifications.action_successful.content', parameters: {} },
            type: NotificationType.SUCCESS,
        });

        emit('close', true);
        emit('refresh');

        return response;
    } catch (e) {
        handleGRPCError(e as RpcError);
        throw e;
    }
}

function setFromProps(): void {
    if (!data.value?.calendar) return;

    const calendar = data.value?.calendar;
    state.name = calendar.name;
    state.description = contentToTiptapValue(calendar.description);
    state.private = calendar.job === undefined;
    state.public = calendar.public;
    state.closed = calendar.closed;
    state.color = calendar.color ?? 'primary';
    if (calendar.access) {
        state.access = calendar.access;
    }
    state.discordSettings = fromProtoDiscordSettings(calendar.discordSettings);

    syncSnapshot();
}

function fromProtoDiscordSettings(
    settings:
        | {
              enabled: boolean;
              channelId: string;
              reminderSteps: {
                  atMinute: number;
                  message?: string;
                  embed?: {
                      title?: string;
                      description?: string;
                      color?: string;
                  };
              }[];
          }
        | undefined,
) {
    if (!settings) return emptyDiscordSettings();

    return {
        enabled: settings.enabled ?? false,
        channelId: settings.channelId ?? '',
        reminderSteps:
            settings.reminderSteps?.map((step) => ({
                atMinute: step.atMinute ?? 0,
                message: step.message ?? '',
                embed: {
                    title: step.embed?.title ?? '',
                    description: step.embed?.description ?? '',
                    color: step.embed?.color ?? '',
                },
            })) ?? [],
    };
}

function toProtoDiscordSettings(values: Schema['discordSettings']) {
    const channelId = values.channelId.trim();
    const reminderSteps = values.reminderSteps.map((step) => {
        const message = step.message.trim();
        const title = step.embed.title.trim();
        const description = step.embed.description.trim();
        const color = step.embed.color.trim();

        const embed =
            title.length > 0 || description.length > 0
                ? {
                      title: title.length > 0 ? title : undefined,
                      description: description.length > 0 ? description : undefined,
                      color: color.length > 0 ? color : undefined,
                  }
                : undefined;

        return {
            atMinute: step.atMinute,
            message: message.length > 0 ? message : undefined,
            embed,
        };
    });

    if (!values.enabled && channelId.length === 0 && reminderSteps.length === 0) {
        return undefined;
    }

    return {
        enabled: values.enabled,
        channelId: channelId.length > 0 ? channelId : '',
        reminderSteps,
    };
}

async function searchChannels() {
    if (!canEditDiscordReminderSettings.value) {
        return [] as Channel[];
    }

    try {
        const call = settingsSettingsClient.listDiscordChannels({});
        const { response } = await call;

        return response.channels.sort((a, b) => a.position - b.position);
    } catch (e) {
        handleGRPCError(e as RpcError);
        return [] as Channel[];
    }
}

watch(data, () => setFromProps());
watch(props, async () => refresh());
watch(
    () => state.private,
    (isPrivate) => {
        if (isPrivate) {
            state.discordSettings = emptyDiscordSettings();
        }
    },
);

const canSubmit = ref<boolean>(true);
const onSubmitThrottle = useThrottleFn(async (event: FormSubmitEvent<Schema>) => {
    canSubmit.value = false;
    await createOrUpdateCalendar(event.data).finally(() => useTimeoutFn(() => (canSubmit.value = true), 400));
}, 1000);

const formRef = useTemplateRef('formRef');

async function closeModal(): Promise<void> {
    if (!canSubmit.value) return;

    if (hasUnsavedChanges.value && !(await confirmLeave())) return;

    emit('close', false);
}
</script>

<template>
    <UModal
        :title="
            calendarId
                ? $t('components.calendar.calendar.CreateOrUpdateModal.update.title')
                : $t('components.calendar.calendar.CreateOrUpdateModal.create.title')
        "
        :close="false"
        :dismissible="!hasUnsavedChanges && canSubmit"
        fullscreen
    >
        <template #header>
            <div class="flex w-full items-center justify-between gap-2">
                <h3 class="font-semibold text-highlighted">
                    {{
                        calendarId
                            ? $t('components.calendar.calendar.CreateOrUpdateModal.update.title')
                            : $t('components.calendar.calendar.CreateOrUpdateModal.create.title')
                    }}
                </h3>

                <UButton
                    color="neutral"
                    variant="ghost"
                    icon="i-mdi-close"
                    :disabled="!canSubmit"
                    :aria-label="$t('common.close', 1)"
                    @click="closeModal"
                />
            </div>
        </template>

        <template #body>
            <UForm ref="formRef" class="flex flex-col gap-2" :schema="schema" :state="state" @submit="onSubmitThrottle">
                <DataPendingBlock
                    v-if="props.calendarId && isRequestPending(status)"
                    :message="$t('common.loading', [$t('common.calendar')])"
                />
                <DataErrorBlock
                    v-else-if="props.calendarId && error"
                    :title="$t('common.unable_to_load', [$t('common.calendar')])"
                    :error="error"
                    :retry="refresh"
                />
                <DataNoDataBlock
                    v-else-if="props.calendarId && (!data || !data.calendar)"
                    :type="$t('common.calendar')"
                    icon="i-mdi-calendar"
                />

                <template v-else>
                    <p v-if="isSystemManaged" class="text-sm text-neutral-500 dark:text-neutral-400">
                        {{ $t('common.read_only') }}
                    </p>

                    <UFormField v-else class="flex-1" name="title" :label="$t('common.name')" required>
                        <UInput v-model="state.name" class="w-full" name="name" type="text" :placeholder="$t('common.name')" />
                    </UFormField>

                    <UFormField class="flex-1" name="color" :label="$t('common.color')">
                        <ColorPickerTW v-model="state.color" class="w-full" />
                    </UFormField>

                    <template v-if="!isSystemManaged">
                        <UFormField class="flex-1" name="description" :label="$t('common.description')">
                            <TiptapEditor
                                v-model="state.description"
                                class="w-full"
                                name="content"
                                wrapper-class="min-h-80"
                                :limit="1_000"
                            />
                        </UFormField>

                        <div class="grid grid-cols-1 gap-4 sm:grid-cols-3">
                            <UFormField
                                class="flex-1"
                                name="private"
                                :label="$t('components.calendar.calendar.CreateOrUpdateModal.private')"
                            >
                                <USwitch
                                    v-model="state.private"
                                    :disabled="
                                        !canDo.privateCalendar ||
                                        calendarId !== undefined ||
                                        (!props.calendarId && hasPrivateCalendar)
                                    "
                                />
                            </UFormField>

                            <UFormField v-if="canDo.publicCalendar" class="flex-1" name="public" :label="$t('common.public')">
                                <USwitch v-model="state.public" />
                            </UFormField>

                            <UFormField class="flex-1" name="closed" :label="`${$t('common.close', 2)}?`">
                                <USwitch v-model="state.closed" />
                            </UFormField>
                        </div>

                        <UPageCard
                            v-if="showDiscordReminderSection"
                            :title="$t('components.calendar.calendar.CreateOrUpdateModal.discord_settings.title')"
                            :description="$t('components.calendar.calendar.CreateOrUpdateModal.discord_settings.description')"
                        >
                            <UFormField
                                class="grid grid-cols-2 items-center gap-2"
                                name="discordSettings.enabled"
                                :label="$t('common.enabled')"
                            >
                                <USwitch v-model="state.discordSettings.enabled" :disabled="!canEditDiscordReminderSettings" />
                            </UFormField>

                            <UFormField
                                class="grid grid-cols-1 items-center gap-2"
                                name="discordSettings.channelId"
                                :label="$t('common.channel')"
                                required
                            >
                                <SelectMenu
                                    v-if="canEditDiscordReminderSettings"
                                    v-model="state.discordSettings.channelId"
                                    class="w-full"
                                    name="discordSettings.channelId"
                                    :disabled="!canConfigureDiscordReminders"
                                    :searchable="
                                        () =>
                                            searchChannels().then((channels) =>
                                                channels.map((channel) => ({
                                                    id: channel.id,
                                                    type: 'item',
                                                    label: `${channel.name} (${channel.id})`,
                                                    item: channel,
                                                })),
                                            )
                                    "
                                    searchable-key="calendar-discord-reminder-channels"
                                    :filter-fields="['name']"
                                    :search-input="{ placeholder: $t('common.search_field') }"
                                    value-key="id"
                                    :placeholder="$t('common.channel')"
                                >
                                    <template #empty>
                                        {{ $t('common.not_found', [$t('common.channel', 1)]) }}
                                    </template>
                                </SelectMenu>

                                <UInput
                                    v-else
                                    class="w-full"
                                    :model-value="state.discordSettings.channelId"
                                    :placeholder="$t('common.channel')"
                                    disabled
                                />
                            </UFormField>

                            <UFormField
                                name="discordSettings.reminderSteps"
                                :label="
                                    $t('components.calendar.calendar.CreateOrUpdateModal.discord_settings.reminder_steps.label')
                                "
                                :description="
                                    $t(
                                        'components.calendar.calendar.CreateOrUpdateModal.discord_settings.reminder_steps.description',
                                    )
                                "
                            >
                                <VueDraggable
                                    v-model="state.discordSettings.reminderSteps"
                                    class="flex flex-col gap-4"
                                    :disabled="!canConfigureDiscordReminders"
                                    handle=".handle-choice"
                                >
                                    <div
                                        v-for="(_, idx) in state.discordSettings.reminderSteps"
                                        :key="idx"
                                        class="rounded-lg border border-default p-3"
                                    >
                                        <div class="mb-3 flex items-center gap-1">
                                            <DraggableHandle
                                                handle-class="handle-choice"
                                                :disabled="!canConfigureDiscordReminders"
                                            />
                                            <ReorderButtons
                                                :idx="idx"
                                                :move-up="moveUp"
                                                :move-down="moveDown"
                                                :button="{ disabled: !canConfigureDiscordReminders }"
                                            />

                                            <div class="ml-auto">
                                                <UTooltip :text="$t('common.remove')">
                                                    <UButton
                                                        icon="i-mdi-close"
                                                        color="error"
                                                        :disabled="!canConfigureDiscordReminders"
                                                        @click="state.discordSettings.reminderSteps.splice(idx, 1)"
                                                    />
                                                </UTooltip>
                                            </div>
                                        </div>

                                        <div class="flex flex-col gap-4">
                                            <UFormField
                                                class="w-full"
                                                :name="`discordSettings.reminderSteps.${idx}.atMinute`"
                                                :label="
                                                    $t(
                                                        'components.calendar.calendar.CreateOrUpdateModal.discord_settings.reminder_steps.at_minute.label',
                                                    )
                                                "
                                                :description="
                                                    $t(
                                                        'components.calendar.calendar.CreateOrUpdateModal.discord_settings.reminder_steps.at_minute.description',
                                                    )
                                                "
                                            >
                                                <UInputNumber
                                                    v-model="state.discordSettings.reminderSteps[idx]!.atMinute"
                                                    class="w-full"
                                                    :min="0"
                                                    :max="10080"
                                                    :step="1"
                                                    :disabled="!canConfigureDiscordReminders"
                                                />
                                            </UFormField>

                                            <UFormField
                                                class="w-full"
                                                :name="`discordSettings.reminderSteps.${idx}.message`"
                                                :label="
                                                    $t(
                                                        'components.calendar.calendar.CreateOrUpdateModal.discord_settings.reminder_steps.message',
                                                    )
                                                "
                                                :description="
                                                    $t(
                                                        'components.calendar.calendar.CreateOrUpdateModal.discord_settings.reminder_steps.message_description',
                                                    )
                                                "
                                            >
                                                <UTextarea
                                                    v-model="state.discordSettings.reminderSteps[idx]!.message"
                                                    class="w-full"
                                                    :disabled="!canConfigureDiscordReminders"
                                                    :rows="3"
                                                    :maxrows="6"
                                                    :placeholder="
                                                        $t(
                                                            'components.calendar.calendar.CreateOrUpdateModal.discord_settings.reminder_steps.message_placeholder',
                                                        )
                                                    "
                                                />
                                            </UFormField>

                                            <USeparator />

                                            <div class="flex flex-row gap-4">
                                                <UFormField
                                                    :name="`discordSettings.reminderSteps.${idx}.embed.color`"
                                                    :label="
                                                        $t(
                                                            'components.calendar.calendar.CreateOrUpdateModal.discord_settings.reminder_steps.embed_color',
                                                        )
                                                    "
                                                >
                                                    <ColorPicker
                                                        v-model="state.discordSettings.reminderSteps[idx]!.embed.color"
                                                        :disabled="!canConfigureDiscordReminders"
                                                    />
                                                </UFormField>

                                                <UFormField
                                                    class="flex-1"
                                                    :name="`discordSettings.reminderSteps.${idx}.embed.title`"
                                                    :label="
                                                        $t(
                                                            'components.calendar.calendar.CreateOrUpdateModal.discord_settings.reminder_steps.embed_title',
                                                        )
                                                    "
                                                >
                                                    <UInput
                                                        v-model="state.discordSettings.reminderSteps[idx]!.embed.title"
                                                        class="w-full"
                                                        :disabled="!canConfigureDiscordReminders"
                                                        :placeholder="
                                                            $t(
                                                                'components.calendar.calendar.CreateOrUpdateModal.discord_settings.reminder_steps.embed_title_placeholder',
                                                            )
                                                        "
                                                    />
                                                </UFormField>
                                            </div>

                                            <UFormField
                                                class="w-full"
                                                :name="`discordSettings.reminderSteps.${idx}.embed.description`"
                                                :label="
                                                    $t(
                                                        'components.calendar.calendar.CreateOrUpdateModal.discord_settings.reminder_steps.embed_description',
                                                    )
                                                "
                                            >
                                                <UTextarea
                                                    v-model="state.discordSettings.reminderSteps[idx]!.embed.description"
                                                    class="w-full"
                                                    :disabled="!canConfigureDiscordReminders"
                                                    :rows="3"
                                                    :maxrows="6"
                                                    :placeholder="
                                                        $t(
                                                            'components.calendar.calendar.CreateOrUpdateModal.discord_settings.reminder_steps.embed_description_placeholder',
                                                        )
                                                    "
                                                />
                                            </UFormField>
                                        </div>
                                    </div>
                                </VueDraggable>

                                <UButton
                                    class="mt-3"
                                    icon="i-mdi-plus"
                                    :disabled="state.discordSettings.reminderSteps.length >= 2 || !canConfigureDiscordReminders"
                                    @click="state.discordSettings.reminderSteps.push(emptyDiscordReminderStep())"
                                />
                            </UFormField>
                        </UPageCard>

                        <UFormField class="flex-1" name="access" :label="$t('common.access')">
                            <AccessManager
                                v-model:jobs="state.access.jobs"
                                v-model:users="state.access.users"
                                :target-id="calendarId ?? 0"
                                :access-roles="enumToAccessLevelEnums(AccessLevel, 'enums.calendar.AccessLevel')"
                            />
                        </UFormField>
                    </template>
                </template>
            </UForm>
        </template>

        <template #footer>
            <UFieldGroup class="inline-flex w-full">
                <UButton
                    class="flex-1"
                    color="neutral"
                    block
                    :disabled="!canSubmit"
                    :label="$t('common.close', 1)"
                    @click="closeModal"
                />

                <UButton
                    class="flex-1"
                    block
                    :disabled="!canSubmit"
                    :loading="!canSubmit"
                    :label="data ? $t('common.save') : $t('common.create')"
                    @click="formRef?.submit()"
                />
            </UFieldGroup>
        </template>
    </UModal>
</template>
