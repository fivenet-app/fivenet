<script lang="ts" setup>
import { logger } from '~/components/documents/helpers';
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import DataNoDataBlock from '~/components/partials/data/DataNoDataBlock.vue';
import DataPendingBlock from '~/components/partials/data/DataPendingBlock.vue';
import { useAuthStore } from '~/store/auth';
import { useClipboardStore } from '~/store/clipboard';
import type { Template } from '~~/gen/ts/resources/documents/templates';

const props = defineProps<{
    templateId: number;
}>();

const { $grpc } = useNuxtApp();

const { isOpen } = useModal();

const authStore = useAuthStore();
const clipboardStore = useClipboardStore();

const { activeChar } = storeToRefs(authStore);

const {
    data: template,
    pending: loading,
    refresh,
    error,
} = useLazyAsyncData(`documents-templates-${props.templateId}`, () => getTemplate());

async function getTemplate(): Promise<Template> {
    try {
        const data = clipboardStore.getTemplateData();
        data.activeChar = activeChar.value!;
        logger.debug('Documents: Editor - Clipboard Template Data', data);

        const call = $grpc.docstore.docStore.getTemplate({
            templateId: props.templateId,
            data,
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
    <UModal :ui="{ width: 'w-full sm:max-w-5xl' }">
        <!-- eslint-disable vue/no-v-html -->
        <UCard :ui="{ ring: '', divide: 'divide-y divide-gray-100 dark:divide-gray-800' }">
            <template #header>
                <div class="flex items-center justify-between">
                    <h3 class="text-2xl font-semibold leading-6">
                        {{ $t('common.document', 1) }}
                        {{ $t('common.preview') }}
                    </h3>

                    <UButton color="gray" variant="ghost" icon="i-mdi-window-close" class="-my-1" @click="isOpen = false" />
                </div>
            </template>

            <div>
                <DataPendingBlock v-if="loading" :message="$t('common.loading', [$t('common.template', 2)])" />
                <DataErrorBlock
                    v-else-if="error"
                    :title="$t('common.unable_to_load', [$t('common.template', 2)])"
                    :error="error"
                    :retry="refresh"
                />
                <DataNoDataBlock v-else-if="!template" :type="$t('common.template', 2)" />

                <template v-else>
                    <div>
                        <label class="mb-2 block text-sm font-medium leading-6">
                            {{ $t('common.title') }}
                        </label>
                        <h1 class="mt-4 break-words rounded-lg p-2 text-2xl font-bold">
                            {{ template?.title }}
                        </h1>
                    </div>
                    <div>
                        <label class="mb-2 block text-sm font-medium leading-6">
                            {{ $t('common.state') }}
                        </label>
                        <p class="mt-4 break-words rounded-lg p-2 text-base font-bold">
                            {{ template?.state }}
                        </p>
                    </div>

                    <label class="mb-2 block text-sm font-medium leading-6">
                        {{ $t('common.content') }}
                    </label>
                    <div class="mt-4 break-words rounded-lg p-2">
                        <div
                            class="tiptap ProseMirror prose prose-sm sm:prose-base lg:prose-lg dark:prose-invert min-w-full max-w-full break-words"
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
                </template>
            </div>

            <template #footer>
                <UButton color="black" block class="flex-1" @click="isOpen = false">
                    {{ $t('common.close', 1) }}
                </UButton>
            </template>
        </UCard>
    </UModal>
</template>
