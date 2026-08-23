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
import { contentToTiptapValue, tiptapToContent } from '~/utils/content';
import { useAuthStore } from '~/stores/auth';
import { useCompletorStore } from '~/stores/completor';
import { getJobsConductClient } from '~~/gen/ts/clients';
import type { File } from '~~/gen/ts/resources/file/file';
import { type ConductEntry, ConductType } from '~~/gen/ts/resources/jobs/conduct/conduct';
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

const { t } = useI18n();

const overlay = useOverlay();

const authStore = useAuthStore();
const { activeChar } = storeToRefs(authStore);

const completorStore = useCompletorStore();

const notifications = useNotificationsStore();

const { maxContentLength } = useAppConfig();

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
    targetUserId: z.coerce.number({ error: t('zod.custom.user') }).positive(),
    type: z.enum(ConductType),
    draft: z.coerce.boolean().default(false),
    message: z.custom<JSONContent | string>().optional(),
    files: z.custom<File>().array().max(5).default([]),
    expiresAt: z.date().optional(),
});

type Schema = z.output<typeof schema>;

const state = reactive<Schema>({
    targetUserId: 0,
    type: ConductType.NOTE,
    draft: false,
    message: undefined,
    files: [],
    expiresAt: undefined,
});

const formSnapshot = computed(() => ({
    targetUserId: state.targetUserId,
    type: state.type,
    draft: state.draft,
    message: state.message,
    files: state.files.map((file) => ({
        id: file.id,
        parentId: file.parentId ?? null,
        filePath: file.filePath,
        byteSize: file.byteSize,
        contentType: file.contentType,
        isDir: file.isDir,
    })),
    expiresAt: state.expiresAt?.toISOString() ?? null,
}));

const { hasUnsavedChanges, confirmLeave, syncSnapshot } = useSnapshotChanges(formSnapshot);

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
                message: tiptapToContent(),
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
                message: tiptapToContent(values.message),
                files: values.files,
                targetUserId: values.targetUserId,
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
        syncSnapshot();
    } catch (e) {
        handleGRPCError(e as RpcError);
        throw e;
    }
}

async function setFromProps(): Promise<void> {
    if (!entry.value) return;

    state.draft = entry.value.draft;
    state.targetUserId = entry.value.targetUserId;
    state.type = entry.value.type;
    state.message = contentToTiptapValue(entry.value.message);
    state.files = entry.value.files ?? [];
    state.expiresAt = entry.value.expiresAt ? toDate(entry.value.expiresAt) : undefined;
    syncSnapshot();
}

setFromProps();
watch(entry, () => setFromProps());
onBeforeMount(() => setFromProps());

const canSubmit = ref<boolean>(true);
const onSubmitThrottle = useThrottleFn(async (event: FormSubmitEvent<Schema>) => {
    canSubmit.value = false;
    await conductCreateOrUpdateEntry(event.data, entry.value?.id).finally(() =>
        useTimeoutFn(() => (canSubmit.value = true), 400),
    );
}, 1000);

const formRef = useTemplateRef('formRef');

const confirmModal = overlay.create(ConfirmModal);

async function closeModal(): Promise<void> {
    if (!canSubmit.value || isRequestPending(status.value)) return;

    if (hasUnsavedChanges.value && !(await confirmLeave())) return;

    emit('close', false);
}
</script>

<template>
    <UModal
        :title="`${$t('components.jobs.conduct.CreateOrUpdateModal.update.title')}: #${entry?.id ?? ''}`"
        :close="false"
        :dismissible="!hasUnsavedChanges && canSubmit && !isRequestPending(status)"
        fullscreen
    >
        <template #header>
            <div class="flex w-full items-center justify-between gap-1.5">
                <h3 class="font-semibold text-highlighted">
                    {{ `${$t('components.jobs.conduct.CreateOrUpdateModal.update.title')}: #${entry?.id ?? ''}` }}
                </h3>

                <UButton
                    color="neutral"
                    variant="ghost"
                    icon="i-mdi-close"
                    :disabled="!canSubmit || isRequestPending(status)"
                    @click="closeModal"
                />
            </div>
        </template>

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
                                        class="w-full"
                                        :items="cTypes"
                                        value-key="status"
                                        :search-input="{ placeholder: $t('common.search_field') }"
                                    >
                                        <template #default>
                                            <UBadge
                                                :color="conductTypesToBadgeColor(state.type)"
                                                truncate
                                                :label="$t(`enums.jobs.ConductType.${ConductType[state.type ?? 0]}`)"
                                            />
                                        </template>

                                        <template #item-label="{ item }">
                                            <UBadge
                                                :color="conductTypesToBadgeColor(item.status)"
                                                truncate
                                                :label="$t(`enums.jobs.ConductType.${ConductType[item.status ?? 0]}`)"
                                            />
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
                            <label class="block text-sm leading-6 font-medium" for="targetUserId">
                                {{ $t('common.target') }}
                            </label>
                        </dt>

                        <dd class="mt-1 text-sm leading-6 sm:col-span-2 sm:mt-0">
                            <UFormField name="targetUserId">
                                <SelectMenu
                                    v-model="state.targetUserId"
                                    class="w-full"
                                    :searchable="
                                        async (q: string) =>
                                            await completorStore.completeColleagues(
                                                q,
                                                state.targetUserId > 0 ? [state.targetUserId] : [],
                                            )
                                    "
                                    searchable-key="completor-colleagues"
                                    :search-input="{ placeholder: $t('common.search_field') }"
                                    :filter-fields="['firstname', 'lastname']"
                                    block
                                    :placeholder="$t('common.colleague')"
                                    trailing
                                    value-key="userId"
                                >
                                    <template #default="{ items }">
                                        <ColleagueName
                                            v-if="items?.find((c) => c.userId === state.targetUserId)"
                                            class="truncate"
                                            :colleague="items.find((c) => c.userId === state.targetUserId)!"
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
                                        class="min-h-120 w-full"
                                        name="message"
                                        :target-id="entry?.id"
                                        history-type="jobs-conduct"
                                        :limit="maxContentLength"
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
                <UButton
                    class="flex-1"
                    color="neutral"
                    block
                    :disabled="!canSubmit || isRequestPending(status)"
                    :label="$t('common.close', 1)"
                    @click="closeModal"
                />

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
