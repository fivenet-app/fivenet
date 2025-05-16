<script lang="ts" setup>
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import DataNoDataBlock from '~/components/partials/data/DataNoDataBlock.vue';
import DataPendingBlock from '~/components/partials/data/DataPendingBlock.vue';
import TiptapEditor from '~/components/partials/editor/TiptapEditor.vue';
import { useMailerStore } from '~/stores/mailer';
import { AccessLevel } from '~~/gen/ts/resources/mailer/access';
import type { ListTemplatesResponse } from '~~/gen/ts/services/mailer/mailer';
import { canAccess } from './helpers';
import TemplateEditForm from './TemplateEditForm.vue';

const { $grpc } = useNuxtApp();

const { isOpen } = useModal();

const mailerStore = useMailerStore();
const { selectedEmail } = storeToRefs(mailerStore);

const {
    data: templates,
    pending: loading,
    error,
    refresh,
} = useLazyAsyncData(`mailer-templates:${selectedEmail.value!.id}`, () => listTemplates());

async function listTemplates(): Promise<ListTemplatesResponse> {
    try {
        const call = $grpc.mailer.mailer.listTemplates({
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

const creating = ref(false);
const editing = ref(false);
</script>

<template>
    <UModal fullscreen>
        <UCard
            :ui="{
                ring: '',
                divide: 'divide-y divide-gray-100 dark:divide-gray-800',
                base: 'flex flex-1 flex-col',
                body: { base: 'flex flex-1 flex-col' },
            }"
        >
            <template #header>
                <div class="flex items-center justify-between">
                    <h3 class="text-2xl font-semibold leading-6">
                        {{ $t('common.template', 2) }}
                    </h3>

                    <UButton class="-my-1" color="gray" variant="ghost" icon="i-mdi-window-close" @click="isOpen = false" />
                </div>
            </template>

            <div class="mx-auto flex w-full max-w-screen-xl flex-col gap-2">
                <UButton
                    v-if="!creating && !editing && canManage"
                    :label="$t('common.create')"
                    trailing-icon="i-mdi-plus"
                    @click="creating = true"
                />

                <TemplateEditForm v-if="creating" @refresh="refresh" @close="creating = false" />
                <template v-else>
                    <DataPendingBlock v-if="loading" :message="$t('common.loading', [$t('common.template')])" />
                    <DataErrorBlock
                        v-else-if="error"
                        :title="$t('common.unable_to_load', [$t('common.template')])"
                        :error="error"
                        :retry="refresh"
                    />
                    <DataNoDataBlock
                        v-if="!templates?.templates || templates?.templates.length === 0"
                        :type="$t('common.template', 2)"
                        icon="i-mdi-file-outline"
                    />

                    <UAccordion v-else :items="accordionItems">
                        <template #item="{ index }">
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

            <template #footer>
                <UButtonGroup class="inline-flex w-full">
                    <UButton class="flex-1" block color="black" @click="isOpen = false">
                        {{ $t('common.close', 1) }}
                    </UButton>
                </UButtonGroup>
            </template>
        </UCard>
    </UModal>
</template>
