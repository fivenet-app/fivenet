<script lang="ts" setup>
import type { JSONContent } from '@tiptap/core';
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import DataNoDataBlock from '~/components/partials/data/DataNoDataBlock.vue';
import DataPendingBlock from '~/components/partials/data/DataPendingBlock.vue';
import TiptapEditor from '~/components/partials/editor/TiptapEditor.vue';
import { useMailerStore } from '~/stores/mailer';
import { getMailerSettingsClient } from '~~/gen/ts/clients';
import { Struct } from '~~/gen/ts/google/protobuf/struct';
import { AccessLevel } from '~~/gen/ts/resources/mailer/access/access';
import type { ListTemplatesResponse } from '~~/gen/ts/services/mailer/settings';
import { canAccess } from './helpers';
import TemplateEditForm from './TemplateEditForm.vue';

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
} = useLazyAsyncData(`mailer-templates:${selectedEmail.value!.id}`, () => listTemplates());

async function listTemplates(): Promise<ListTemplatesResponse> {
    try {
        const call = mailerSettingsClient.listTemplates({
            emailId: selectedEmail.value!.id,
        });
        const { response } = await call;

        return response;
    } catch (e) {
        handleGRPCError(e as RpcError);
        throw e;
    }
}

const accordionItems = computed(() =>
    templates.value?.templates.map((t) => ({
        label: t.title,
    })),
);

const canManage = computed(() => canAccess(selectedEmail.value?.access, selectedEmail.value?.userId, AccessLevel.MANAGE));

const creating = ref<boolean>(false);
const editing = ref<boolean>(false);
const childDirty = ref<boolean>(false);

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
    editing.value = false;
    childDirty.value = false;
}
</script>

<template>
    <UModal :title="$t('common.template', 2)" :close="false" :dismissible="!hasUnsavedChanges" fullscreen>
        <template #header>
            <div class="flex w-full items-center justify-between gap-1.5">
                <h3 class="font-semibold text-highlighted">
                    {{ $t('common.template', 2) }}
                </h3>

                <UButton color="neutral" variant="ghost" icon="i-mdi-close" @click="closeModal" />
            </div>
        </template>

        <template #body>
            <div class="mx-auto flex w-full max-w-(--breakpoint-xl) flex-col gap-2">
                <UButton
                    v-if="!creating && !editing && canManage"
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

                    <UAccordion v-else :items="accordionItems">
                        <template #content="{ index }">
                            <template v-if="templates?.templates[index]">
                                <template v-if="!editing">
                                    <UFieldGroup v-if="canManage" class="mx-4 mb-2 flex">
                                        <UTooltip :text="$t('common.edit')">
                                            <UButton
                                                class="flex-1"
                                                icon="i-mdi-pencil"
                                                :label="$t('common.edit')"
                                                @click="editing = !editing"
                                            />
                                        </UTooltip>

                                        <UTooltip :text="$t('common.delete')">
                                            <UButton icon="i-mdi-delete" color="error" :label="$t('common.delete')" />
                                        </UTooltip>
                                    </UFieldGroup>

                                    <ClientOnly>
                                        <TiptapEditor
                                            :model-value="
                                                templates.templates[index].content?.tiptapJson
                                                    ? (Struct.toJson(
                                                          templates.templates[index].content.tiptapJson,
                                                      ) as JSONContent)
                                                    : (templates.templates[index].content?.rawHtml ?? '')
                                            "
                                            disabled
                                            hide-toolbar
                                            wrapper-class="min-h-44"
                                        />
                                    </ClientOnly>
                                </template>
                                <TemplateEditForm
                                    v-else
                                    :template="templates.templates[index]"
                                    @refresh="refresh"
                                    @dirty-change="childDirty = $event"
                                    @close="resetEditing"
                                />
                            </template>
                        </template>
                    </UAccordion>
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
