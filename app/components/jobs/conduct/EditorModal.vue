<script lang="ts" setup>
import type { FormSubmitEvent } from '@nuxt/ui';
import type { JSONContent } from '@tiptap/core';
import { z } from 'zod';
import ConfirmModal from '~/components/partials/ConfirmModal.vue';
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import DataNoDataBlock from '~/components/partials/data/DataNoDataBlock.vue';
import DataPendingBlock from '~/components/partials/data/DataPendingBlock.vue';
import TiptapEditor from '~/components/partials/editor/TiptapEditor.vue';
import InputDatePicker from '~/components/partials/InputDatePicker.vue';
import SelectMenu from '~/components/partials/SelectMenu.vue';
import { useAuthStore } from '~/stores/auth';
import { useCompletorStore } from '~/stores/completor';
import { getJobsConductClient } from '~~/gen/ts/clients';
import { Struct } from '~~/gen/ts/google/protobuf/struct';
import { ContentType } from '~~/gen/ts/resources/common/content/content';
import type { File } from '~~/gen/ts/resources/file/file';
import { type ConductEntry, ConductType } from '~~/gen/ts/resources/jobs/conduct';
import { NotificationType } from '~~/gen/ts/resources/notifications/notifications';
import ColleagueName from '../colleagues/ColleagueName.vue';
import { conductTypesToBadgeColor } from './helpers';

const props = defineProps<{
    entryId?: number;
    userId?: number;
}>();

const emit = defineEmits<{
    (e: 'close', v: boolean): void;
    (e: 'created', entry: ConductEntry): void;
    (e: 'updated', entry: ConductEntry): void;
}>();

const overlay = useOverlay();

const authStore = useAuthStore();
const { activeChar } = storeToRefs(authStore);

const completorStore = useCompletorStore();

const notifications = useNotificationsStore();

const jobsConductClient = await getJobsConductClient();

const cTypes = ref<{ status: ConductType }[]>([
    { status: ConductType.NOTE },
    { status: ConductType.NEUTRAL },
    { status: ConductType.POSITIVE },
    { status: ConductType.NEGATIVE },
    { status: ConductType.WARNING },
    { status: ConductType.SUSPENSION },
]);

const schema = z.object({
    targetUser: z.coerce.number().positive(),
    type: z.enum(ConductType),
    draft: z.coerce.boolean().default(false),
    message: z.custom<JSONContent | string>().optional(),
    files: z.custom<File>().array().max(5).default([]),
    expiresAt: z.date().optional(),
});

type Schema = z.output<typeof schema>;

const state = reactive<Schema>({
    targetUser: 0,
    type: ConductType.NOTE,
    draft: false,
    message: undefined,
    files: [],
    expiresAt: undefined,
});

const {
    data: entry,
    error,
    status,
    refresh,
} = await useLazyAsyncData(`conduct-entry-${props.entryId ?? 'new'}`, async () => {
    if (!props.entryId) {
        const call = jobsConductClient.createConductEntry({
            entry: {
                id: 0,
                targetUserId: props.userId ?? 0,
                job: '',
                creatorId: activeChar.value?.userId ?? 0,
                type: ConductType.NOTE,
                draft: true,
                message: {
                    contentType: ContentType.TIPTAP_JSON,
                    version: '',
                    tiptapJson: Struct.fromJsonString(JSON.stringify({ type: 'doc', content: [] })),
                },
                files: [],
            },
        });
        const { response } = await call;

        emit('created', response.entry!);

        return response.entry;
    }

    const call = jobsConductClient.getConductEntry({ id: props.entryId });
    const { response } = await call;

    return response.entry;
});

watch(props, async () => refresh());

async function conductCreateOrUpdateEntry(values: Schema, id?: number): Promise<void> {
    try {
        const req = {
            entry: {
                id: id ?? 0,
                job: '',
                creatorId: activeChar.value?.userId ?? 0,
                type: values.type,
                draft: values.draft,
                message: {
                    contentType: ContentType.TIPTAP_JSON,
                    version: '',
                    tiptapJson: Struct.fromJsonString(JSON.stringify(values.message)),
                },
                files: values.files,
                targetUserId: values.targetUser,
                expiresAt: values.expiresAt ? toTimestamp(values.expiresAt) : undefined,
            },
        };

        const call = jobsConductClient.updateConductEntry(req);
        const { response } = await call;

        emit('updated', response.entry!);

        notifications.add({
            title: { key: 'notifications.action_successful.title', parameters: {} },
            description: { key: 'notifications.action_successful.content', parameters: {} },
            type: NotificationType.SUCCESS,
        });

        emit('close', false);
    } catch (e) {
        handleGRPCError(e as RpcError);
        throw e;
    }
}

async function setFromData(): Promise<void> {
    if (!entry.value) return;

    state.targetUser = entry.value.targetUserId ?? 0;
    state.type = entry.value.type ?? ConductType.NOTE;
    state.message = entry.value.message?.tiptapJson
        ? (Struct.toJson(entry.value.message.tiptapJson) as JSONContent)
        : (entry.value.message?.rawHtml ?? '');
    state.expiresAt = entry.value.expiresAt ? toDate(entry.value.expiresAt) : undefined;
}

watch(entry, () => setFromData());

onBeforeMount(() => setFromData());

const canSubmit = ref(true);
const onSubmitThrottle = useThrottleFn(async (event: FormSubmitEvent<Schema>) => {
    canSubmit.value = false;
    await conductCreateOrUpdateEntry(event.data, entry.value?.id).finally(() =>
        useTimeoutFn(() => (canSubmit.value = true), 400),
    );
}, 1000);

const formRef = useTemplateRef('formRef');

const confirmModal = overlay.create(ConfirmModal);
</script>

<template>
    <UModal :title="`${$t('components.jobs.conduct.CreateOrUpdateModal.update.title')}: #${entry?.id ?? ''}`" fullscreen>
        <template #body>
            <DataPendingBlock v-if="isRequestPending(status)" :message="$t('common.loading', [$t('common.entry', 1)])" />
            <DataErrorBlock
                v-else-if="error"
                :title="$t('common.unable_to_load', [$t('common.conduct_register')])"
                :error="error"
                :retry="refresh"
            />
            <DataNoDataBlock v-else-if="!entry" class="w-full" icon="i-mdi-pulse" :type="$t('common.entry', 1)" />

            <UForm v-else ref="formRef" :schema="schema" :state="state" @submit="onSubmitThrottle">
                <dl class="divide-y divide-default">
                    <div class="px-4 py-3 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-0">
                        <dt class="text-sm leading-6 font-medium">
                            <label class="block text-sm leading-6 font-medium" for="type">
                                {{ $t('common.type') }}
                            </label>
                        </dt>
                        <dd class="mt-1 text-sm leading-6 sm:col-span-2 sm:mt-0">
                            <UFormField name="type">
                                <ClientOnly>
                                    <USelectMenu
                                        v-model="state.type"
                                        :items="cTypes"
                                        value-key="status"
                                        :search-input="{ placeholder: $t('common.search_field') }"
                                        class="w-full"
                                    >
                                        <template #default>
                                            <UBadge :color="conductTypesToBadgeColor(state.type)" truncate>
                                                {{ $t(`enums.jobs.ConductType.${ConductType[state.type ?? 0]}`) }}
                                            </UBadge>
                                        </template>

                                        <template #item-label="{ item }">
                                            <UBadge :color="conductTypesToBadgeColor(item.status)" truncate>
                                                {{ $t(`enums.jobs.ConductType.${ConductType[item.status ?? 0]}`) }}
                                            </UBadge>
                                        </template>

                                        <template #empty>
                                            {{ $t('common.not_found', [$t('common.type', 2)]) }}
                                        </template>
                                    </USelectMenu>
                                </ClientOnly>
                            </UFormField>
                        </dd>
                    </div>

                    <div class="px-4 py-3 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-0">
                        <dt class="text-sm leading-6 font-medium">
                            <label class="block text-sm leading-6 font-medium" for="targetUser">
                                {{ $t('common.target') }}
                            </label>
                        </dt>
                        <dd class="mt-1 text-sm leading-6 sm:col-span-2 sm:mt-0">
                            <UFormField name="targetUserId">
                                <SelectMenu
                                    v-model="state.targetUser"
                                    :searchable="async (q: string) => await completorStore.completeColleagues(q)"
                                    searchable-key="completor-colleagues"
                                    :search-input="{ placeholder: $t('common.search_field') }"
                                    :filter-fields="['firstname', 'lastname']"
                                    block
                                    :placeholder="$t('common.colleague')"
                                    trailing
                                    class="w-full"
                                    value-key="userId"
                                >
                                    <template #default="{ items }">
                                        <ColleagueName
                                            v-if="items?.find((c) => c.userId === state.targetUser)"
                                            class="truncate"
                                            :colleague="items.find((c) => c.userId === state.targetUser)!"
                                            birthday
                                        />
                                        <span v-else>&nbsp;</span>
                                    </template>

                                    <template #item-label="{ item }">
                                        <ColleagueName class="truncate" :colleague="item" birthday />
                                    </template>

                                    <template #empty>
                                        {{ $t('common.not_found', [$t('common.creator', 2)]) }}
                                    </template>
                                </SelectMenu>
                            </UFormField>
                        </dd>
                    </div>

                    <div class="px-4 py-3 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-0">
                        <dt class="text-sm leading-6 font-medium">
                            <label class="block text-sm leading-6 font-medium" for="message">
                                {{ $t('common.content') }}
                            </label>
                        </dt>
                        <dd class="mt-1 text-sm leading-6 sm:col-span-2 sm:mt-0">
                            <UFormField name="message" :ui="{ error: 'hidden' }">
                                <ClientOnly>
                                    <TiptapEditor
                                        v-model="state.message"
                                        v-model:files="state.files"
                                        name="message"
                                        class="min-h-120 w-full"
                                        :target-id="entry?.id"
                                        filestore-namespace="jobs-conduct"
                                        :filestore-service="(opts) => jobsConductClient.uploadFile(opts)"
                                    />
                                </ClientOnly>
                            </UFormField>
                        </dd>
                    </div>

                    <div class="px-4 py-3 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-0">
                        <dt class="text-sm leading-6 font-medium">
                            <label class="block text-sm leading-6 font-medium" for="expiresAt">
                                {{ $t('common.expires_at') }}?
                            </label>
                        </dt>
                        <dd class="mt-1 text-sm leading-6 sm:col-span-2 sm:mt-0">
                            <UFormField name="expiresAt">
                                <InputDatePicker v-model="state.expiresAt" clearable time />
                            </UFormField>
                        </dd>
                    </div>
                </dl>
            </UForm>
        </template>

        <template #footer>
            <UFieldGroup class="inline-flex w-full">
                <UButton class="flex-1" color="neutral" block :label="$t('common.close', 1)" @click="$emit('close', false)" />

                <UButton
                    class="flex-1"
                    block
                    trailing-icon="i-mdi-content-save"
                    :disabled="!canSubmit || isRequestPending(status)"
                    :loading="!canSubmit || isRequestPending(status)"
                    :label="$t('common.save')"
                    @click="() => formRef?.submit()"
                />

                <UButton
                    v-if="entry?.draft"
                    class="flex-1"
                    block
                    color="info"
                    trailing-icon="i-mdi-publish"
                    :disabled="!canSubmit || isRequestPending(status)"
                    :loading="!canSubmit || isRequestPending(status)"
                    :label="$t('common.publish')"
                    @click="
                        confirmModal.open({
                            title: $t('common.publish_confirm.title', { type: $t('common.entry', 1) }),
                            description: $t('common.publish_confirm.description'),
                            color: 'info',
                            iconClass: 'text-info-500 dark:text-info-400',
                            icon: 'i-mdi-publish',
                            confirm: () => {
                                state.draft = false;
                                formRef?.submit();
                            },
                        })
                    "
                />
            </UFieldGroup>
        </template>
    </UModal>
</template>
