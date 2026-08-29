<script lang="ts" setup>
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import DataNoDataBlock from '~/components/partials/data/DataNoDataBlock.vue';
import DataPendingBlock from '~/components/partials/data/DataPendingBlock.vue';
import CategoryBadge from '~/components/partials/documents/CategoryBadge.vue';
import RefreshButton from '~/components/partials/RefreshButton.vue';
import { useClipboardStore } from '~/stores/clipboard';
import { getDocumentsTemplatesClient } from '~~/gen/ts/clients';
import type { Template } from '~~/gen/ts/resources/documents/templates/templates';

const props = defineProps<{
    templateId: number;
}>();

defineEmits<{
    (e: 'close', v: boolean): void;
}>();

const clipboardStore = useClipboardStore();

const logger = useLogger('📃 Doc Templates');

const documentsTemplatesClient = await getDocumentsTemplatesClient();

const {
    data: template,
    status,
    refresh,
    error,
} = useLazyAsyncData(`documents-templates-${props.templateId}`, () => getTemplate());

const loading = computed(() => isRequestPending(status.value));

async function getTemplate(): Promise<Template> {
    try {
        const selection = clipboardStore.getTemplateSelection(false);
        logger.debug('Documents: Editor - Clipboard Template Selection', selection);

        const call = documentsTemplatesClient.getTemplate({
            templateId: props.templateId,
            selection: selection,
            render: true,
        });
        const { response } = await call;

        return response.template!;
    } catch (e) {
        handleGRPCError(e as RpcError);
        throw e;
    }
}
</script>

<template>
    <UModal :title="`${$t('common.template', 1)} ${$t('common.preview')}`" fullscreen>
        <template #body>
            <DataPendingBlock v-if="loading" :message="$t('common.loading', [$t('common.template', 2)])" />
            <DataErrorBlock
                v-else-if="error"
                :title="$t('common.unable_to_load', [$t('common.template', 2)])"
                :error="error"
                :retry="refresh"
            />
            <DataNoDataBlock v-else-if="!template" :type="$t('common.template', 2)" />

            <div v-else class="mx-auto flex w-full max-w-(--breakpoint-xl) flex-col gap-2">
                <UFormField name="title" :label="$t('common.title')">
                    <UInput class="w-full" :model-value="template?.title" type="text" size="xl" disabled />
                </UFormField>

                <div class="flex flex-row gap-2">
                    <UFormField class="flex-1" name="category" :label="$t('common.category', 1)">
                        <CategoryBadge v-if="template?.category" :category="template.category" />
                        <span v-else>{{ $t('common.categories', 0) }}</span>
                    </UFormField>

                    <UFormField class="flex-1" name="state" :label="$t('common.state')">
                        <UInput class="w-full" :model-value="template?.state" type="text" disabled />
                    </UFormField>
                </div>

                <UFormField name="content" :label="$t('common.content')">
                    <div
                        class="mx-auto w-full max-w-(--breakpoint-xl) rounded-lg bg-neutral-100 p-4 break-words dark:bg-neutral-800"
                    >
                        <!-- eslint-disable vue/no-v-html -->
                        <div
                            class="tiptap prose prose-sm max-w-full min-w-full break-words sm:prose-base lg:prose-lg dark:prose-invert"
                            :class="[
                                'hover:prose-a:text-blue-500',
                                'dark:hover:prose-a:text-blue-300',
                                'prose-headings:mt-0.5',
                                'prose-lead:mt-0.5',
                                'prose-h1:mt-0.5',
                                'prose-h2:mt-0.5',
                                'prose-h3:mt-0.5',
                                'prose-h4:mt-0.5',
                                'prose-p:mt-0.5',
                                'prose-a:mt-0.5',
                                'prose-blockquote:mt-0.5',
                                'prose-figure:mt-0.5',
                                'prose-figcaption:mt-0.5',
                                'prose-strong:mt-0.5',
                                'prose-em:mt-0.5',
                                'prose-kbd:mt-0.5',
                                'prose-code:mt-0.5',
                                'prose-pre:mt-0.5',
                                'prose-ol:mt-0.5',
                                'prose-ul:mt-0.5',
                                'prose-li:mt-0.5',
                                'prose-table:mt-0.5',
                                'prose-thead:mt-0.5',
                                'prose-tr:mt-0.5',
                                'prose-th:mt-0.5',
                                'prose-td:mt-0.5',
                                'prose-img:mt-0.5',
                                'prose-video:mt-0.5',
                                'prose-hr:mt-0.5',
                            ]"
                            v-html="template?.content"
                        />
                    </div>
                </UFormField>
            </div>
        </template>

        <template #footer>
            <UFieldGroup class="inline-flex w-full">
                <UButton class="flex-1" color="neutral" block :label="$t('common.close', 1)" @click="$emit('close', false)" />
                <RefreshButton variant="solid" :loading="loading" :disabled="loading" @click="() => refresh()" />
            </UFieldGroup>
        </template>
    </UModal>
</template>
