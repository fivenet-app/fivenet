<script lang="ts" setup>
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import DataNoDataBlock from '~/components/partials/data/DataNoDataBlock.vue';
import DataPendingBlock from '~/components/partials/data/DataPendingBlock.vue';
import TiptapEditor from '~/components/partials/editor/TiptapEditor.vue';
import { useMailerStore } from '~/stores/mailer';
import { getMailerMailerClient } from '~~/gen/ts/clients';
import { AccessLevel } from '~~/gen/ts/resources/mailer/access';
import type { ListTemplatesResponse } from '~~/gen/ts/services/mailer/mailer';
import { canAccess } from './helpers';
import TemplateEditForm from './TemplateEditForm.vue';

defineEmits<{
    (e: 'close', v: boolean): void;
}>();

const mailerStore = useMailerStore();
const { selectedEmail } = storeToRefs(mailerStore);

const mailerMailerClient = await getMailerMailerClient();

const {
    data: templates,
    status,
    error,
    refresh,
} = useLazyAsyncData(`mailer-templates:${selectedEmail.value!.id}`, () => listTemplates());

async function listTemplates(): Promise<ListTemplatesResponse> {
    try {
        const call = mailerMailerClient.listTemplates({
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
        ...t,
        label: t.title,
    })),
);

const canManage = computed(() => canAccess(selectedEmail.value?.access, selectedEmail.value?.userId, AccessLevel.MANAGE));

const creating = ref(false);
const editing = ref(false);
</script>

<template>
    <UModal :title="$t('common.template', 2)" fullscreen>
        <template #body>
            <div class="mx-auto flex w-full max-w-(--breakpoint-xl) flex-col gap-2">
                <UButton
                    v-if="!creating && !editing && canManage"
                    :label="$t('common.create')"
                    trailing-icon="i-mdi-plus"
                    @click="creating = true"
                />

                <TemplateEditForm v-if="creating" @refresh="refresh" @close="creating = false" />
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
                                    <UButtonGroup v-if="canManage" class="mx-4 mb-2 flex">
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
                                    </UButtonGroup>

                                    <ClientOnly>
                                        <TiptapEditor
                                            v-model="templates.templates[index].content"
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
                                    @close="editing = false"
                                />
                            </template>
                        </template>
                    </UAccordion>
                </template>
            </div>
        </template>

        <template #footer>
            <UButtonGroup class="inline-flex w-full">
                <UButton class="flex-1" block color="neutral" @click="$emit('close', false)">
                    {{ $t('common.close', 1) }}
                </UButton>
            </UButtonGroup>
        </template>
    </UModal>
</template>
