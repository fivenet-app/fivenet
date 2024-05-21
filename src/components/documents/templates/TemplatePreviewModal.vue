<script lang="ts" setup>
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import DataNoDataBlock from '~/components/partials/data/DataNoDataBlock.vue';
import DataPendingBlock from '~/components/partials/data/DataPendingBlock.vue';
import { useAuthStore } from '~/store/auth';
import { useClipboardStore } from '~/store/clipboard';
import { Template } from '~~/gen/ts/resources/documents/templates';

const props = defineProps<{
    templateId: string;
}>();

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
        console.debug('Documents: Editor - Clipboard Template Data', data);

        const call = getGRPCDocStoreClient().getTemplate({
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
                    :retry="refresh"
                />
                <DataNoDataBlock v-else-if="template === null" :type="$t('common.template', 2)" />

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
                        <!-- eslint-disable-next-line vue/no-v-html -->
                        <p v-html="template?.content"></p>
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
