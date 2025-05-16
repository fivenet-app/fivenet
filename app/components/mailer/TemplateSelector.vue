<script lang="ts" setup>
import { useMailerStore } from '~/stores/mailer';
import type { Template } from '~~/gen/ts/resources/mailer/template';
import type { ListTemplatesResponse } from '~~/gen/ts/services/mailer/mailer';

const props = defineProps<{
    modelValue: string;
}>();

const emit = defineEmits<{
    (e: 'update:modelValue', value: string): void;
}>();

const content = useVModel(props, 'modelValue', emit);

defineOptions({
    inheritAttrs: false,
});

const { $grpc } = useNuxtApp();

const mailerStore = useMailerStore();
const { selectedEmail } = storeToRefs(mailerStore);

const { data: templates } = useLazyAsyncData(`mailer-templates:${selectedEmail.value!.id}`, () => listTemplates());

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

const selectedTemplate = ref<Template | undefined>(undefined);
</script>

<template>
    <ClientOnly v-if="templates?.templates && templates?.templates.length > 0">
        <USelectMenu
            v-if="!selectedTemplate"
            v-model="selectedTemplate"
            class="min-w-48"
            :options="templates?.templates"
            option-attribute="title"
            by="id"
            searchable
            :placeholder="$t('common.template')"
            v-bind="$attrs"
        >
            <template #option-empty="{ query: search }">
                <q>{{ search }}</q> {{ $t('common.query_not_found') }}
            </template>

            <template #empty>
                {{ $t('common.not_found', [$t('common.template', 2)]) }}
            </template>
        </USelectMenu>

        <UButtonGroup v-else v-bind="$attrs">
            <UButton
                :label="$t('common.confirm')"
                icon="i-mdi-check"
                @click="
                    selectedTemplate && (content = selectedTemplate.content + (selectedEmail?.settings?.signature ?? ''));
                    selectedTemplate = undefined;
                "
            />

            <UButton :label="$t('common.cancel')" color="error" icon="i-mdi-cancel" @click="selectedTemplate = undefined" />
        </UButtonGroup>
    </ClientOnly>
</template>
