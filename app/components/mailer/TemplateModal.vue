<script lang="ts" setup>
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import DataNoDataBlock from '~/components/partials/data/DataNoDataBlock.vue';
import DataPendingBlock from '~/components/partials/data/DataPendingBlock.vue';
import TiptapEditor from '~/components/partials/editor/TiptapEditor.vue';
import { useMailerStore } from '~/stores/mailer';
import { getMailerSettingsClient } from '~~/gen/ts/clients';
import { contentToTiptapValue } from '~/utils/content';
import { AccessLevel } from '~~/gen/ts/resources/mailer/access/access';
import type { ListTemplatesResponse } from '~~/gen/ts/services/mailer/settings';
import { NotificationType } from '~~/gen/ts/resources/notifications/notifications';
import { canAccess } from './helpers';
import TemplateEditForm from './TemplateEditForm.vue';
import ConfirmModal from '~/components/partials/ConfirmModal.vue';

const emit = defineEmits<{
    (e: 'close', v: boolean): void;
}>();

const mailerStore = useMailerStore();
const { selectedEmail } = storeToRefs(mailerStore);

const mailerSettingsClient = await getMailerSettingsClient();

const {
    data: templates,
    status,
    error,
    refresh,
} = useAuthedLazyAsyncData('userState', `mailer-templates:${selectedEmail.value!.id}`, ({ signal }) => listTemplates(signal));

async function listTemplates(signal: AbortSignal): Promise<ListTemplatesResponse> {
    try {
        const call = mailerSettingsClient.listTemplates(
            {
                emailId: selectedEmail.value!.id,
            },
            { abort: signal },
        );
        const { response } = await call;

        return response;
    } catch (e) {
        handleGRPCError(e as RpcError);
        throw e;
    }
}

const templateSearch = ref('');
const visibleTemplates = computed(() => {
    const search = templateSearch.value.trim().toLowerCase();
    if (!search) return templates.value?.templates ?? [];

    return (templates.value?.templates ?? []).filter((template) => template.title.toLowerCase().includes(search));
});

const accordionItems = computed(() =>
    visibleTemplates.value.map((t) => ({
        label: t.title,
    })),
);

const canManage = computed(() => canAccess(selectedEmail.value?.access, selectedEmail.value?.userId, AccessLevel.MANAGE));

const creating = ref<boolean>(false);
const editingTemplateId = ref<number | undefined>(undefined);
const childDirty = ref<boolean>(false);
const overlay = useOverlay();
const notifications = useNotificationsStore();

const { hasUnsavedChanges, confirmLeave } = useUnsavedChanges({
    dirty: childDirty,
});

async function closeModal(): Promise<void> {
    if (hasUnsavedChanges.value && !(await confirmLeave())) return;

    emit('close', false);
}

function resetCreating(): void {
    creating.value = false;
    childDirty.value = false;
}

function resetEditing(): void {
    editingTemplateId.value = undefined;
    childDirty.value = false;
}

async function deleteTemplate(templateId: number): Promise<void> {
    if (!selectedEmail.value?.id) return;

    try {
        await mailerSettingsClient.deleteTemplate({
            emailId: selectedEmail.value.id,
            id: templateId,
        });

        notifications.add({
            title: { key: 'notifications.action_successful.title', parameters: {} },
            description: { key: 'notifications.action_successful.content', parameters: {} },
            type: NotificationType.SUCCESS,
        });

        await refresh();
    } catch (e) {
        handleGRPCError(e as RpcError);
        throw e;
    }
}

const confirmModal = overlay.create(ConfirmModal);
</script>

<template>
    <UModal :title="$t('common.template', 2)" :close="false" :dismissible="!hasUnsavedChanges" fullscreen>
        <template #header>
            <div class="flex w-full items-center justify-between gap-1.5">
                <h3 class="font-semibold text-highlighted">
                    {{ $t('common.template', 2) }}
                </h3>

                <UButton
                    color="neutral"
                    variant="ghost"
                    icon="i-mdi-close"
                    :aria-label="$t('common.close', 1)"
                    @click="closeModal"
                />
            </div>
        </template>

        <template #body>
            <div class="mx-auto flex w-full max-w-(--breakpoint-xl) flex-col gap-2">
                <UButton
                    v-if="!creating && editingTemplateId === undefined && canManage"
                    :label="$t('common.create')"
                    trailing-icon="i-mdi-plus"
                    @click="creating = true"
                />

                <TemplateEditForm
                    v-if="creating"
                    @refresh="refresh"
                    @dirty-change="childDirty = $event"
                    @close="resetCreating"
                />
                <template v-else>
                    <DataPendingBlock
                        v-if="isRequestPending(status)"
                        :message="$t('common.loading', [$t('common.template')])"
                    />
                    <DataErrorBlock
                        v-else-if="error"
                        :title="$t('common.unable_to_load', [$t('common.template')])"
                        :error="error"
                        :retry="refresh"
                    />
                    <DataNoDataBlock
                        v-else-if="!templates?.templates || templates?.templates.length === 0"
                        :type="$t('common.template', 2)"
                        icon="i-mdi-file-outline"
                    />

                    <template v-else-if="visibleTemplates.length === 0">
                        <UInput v-model="templateSearch" :placeholder="$t('common.search_field')" />
                        <DataNoDataBlock
                            :message="$t('common.not_found', [$t('common.template', 2)])"
                            :type="$t('common.template', 2)"
                            icon="i-mdi-file-search-outline"
                        />
                    </template>

                    <template v-else>
                        <UInput v-model="templateSearch" :placeholder="$t('common.search_field')" />
                        <UAccordion :items="accordionItems">
                            <template #content="{ index }">
                                <template v-if="visibleTemplates[index]">
                                    <template v-if="editingTemplateId !== visibleTemplates[index].id">
                                        <UCard variant="subtle" :ui="{ body: 'p-4 sm:p-4' }">
                                            <UFieldGroup v-if="canManage" class="mb-3 flex">
                                                <UTooltip :text="$t('common.edit')">
                                                    <UButton
                                                        class="flex-1"
                                                        icon="i-mdi-pencil"
                                                        :label="$t('common.edit')"
                                                        @click="editingTemplateId = visibleTemplates[index].id"
                                                    />
                                                </UTooltip>

                                                <UTooltip :text="$t('common.delete')">
                                                    <UButton
                                                        icon="i-mdi-delete"
                                                        color="error"
                                                        :label="$t('common.delete')"
                                                        @click="
                                                            confirmModal.open({
                                                                confirm: async () =>
                                                                    deleteTemplate(visibleTemplates[index]!.id),
                                                            })
                                                        "
                                                    />
                                                </UTooltip>
                                            </UFieldGroup>

                                            <ClientOnly>
                                                <TiptapEditor
                                                    :model-value="contentToTiptapValue(visibleTemplates[index].content)"
                                                    disabled
                                                    hide-toolbar
                                                    wrapper-class="min-h-44"
                                                />
                                            </ClientOnly>
                                        </UCard>
                                    </template>
                                    <TemplateEditForm
                                        v-else
                                        :template="visibleTemplates[index]"
                                        @refresh="refresh"
                                        @dirty-change="childDirty = $event"
                                        @close="resetEditing"
                                    />
                                </template>
                            </template>
                        </UAccordion>
                    </template>
                </template>
            </div>
        </template>

        <template #footer>
            <UFieldGroup class="inline-flex w-full">
                <UButton class="flex-1" block color="neutral" :label="$t('common.close', 1)" @click="closeModal" />
            </UFieldGroup>
        </template>
    </UModal>
</template>
