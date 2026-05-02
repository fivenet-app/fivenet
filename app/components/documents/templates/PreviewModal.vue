<script lang="ts" setup>
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import DataNoDataBlock from '~/components/partials/data/DataNoDataBlock.vue';
import DataPendingBlock from '~/components/partials/data/DataPendingBlock.vue';
import CategoryBadge from '~/components/partials/documents/CategoryBadge.vue';
import { useAuthStore } from '~/stores/auth';
import { useClipboardStore } from '~/stores/clipboard';
import { getDocumentsDocumentsClient } from '~~/gen/ts/clients';
import type { Template } from '~~/gen/ts/resources/documents/templates/templates';

const props = defineProps<{
    templateId: number;
}>();

defineEmits<{
    (e: 'close', v: boolean): void;
}>();

const authStore = useAuthStore();
const clipboardStore = useClipboardStore();

const { activeChar } = storeToRefs(authStore);

const logger = useLogger('📃 Doc Templates');

const documentsDocumentsClient = await getDocumentsDocumentsClient();

const {
    data: template,
    status,
    refresh,
    error,
} = useLazyAsyncData(`documents-templates-${props.templateId}`, () => getTemplate());

async function getTemplate(): Promise<Template> {
    try {
        const data = clipboardStore.getTemplateData();
        data.activeChar = activeChar.value!;
        logger.debug('Documents: Editor - Clipboard Template Data', data);

        const call = documentsDocumentsClient.getTemplate({
            templateId: props.templateId,
            data: data,
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
        <!-- eslint-disable vue/no-v-html -->

        <template #body>
            <DataPendingBlock v-if="isRequestPending(status)" :message="$t('common.loading', [$t('common.template', 2)])" />
            <DataErrorBlock
                v-else-if="error"
                :title="$t('common.unable_to_load', [$t('common.template', 2)])"
                :error="error"
                :retry="refresh"
            />
            <DataNoDataBlock v-else-if="!template" :type="$t('common.template', 2)" />

            <div v-else class="mx-auto w-full max-w-(--breakpoint-xl)">
                <div class="mb-2">
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
                </div>

                <UFormField name="content" :label="$t('common.content')">
                    <div
                        class="mx-auto w-full max-w-(--breakpoint-xl) rounded-lg bg-neutral-100 p-4 break-words dark:bg-neutral-800"
                    >
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
                        ></div>
                    </div>
                </UFormField>
            </div>
        </template>

        <template #footer>
            <UButton class="flex-1" color="neutral" block :label="$t('common.close', 1)" @click="$emit('close', false)" />
        </template>
    </UModal>
</template>
