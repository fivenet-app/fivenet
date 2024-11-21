<script lang="ts" setup>
import { useMailerStore } from '~/store/mailer';
import type { Template } from '~~/gen/ts/resources/mailer/template';
import type { ListTemplatesResponse } from '~~/gen/ts/services/mailer/mailer';

const props = defineProps<{
    modelValue: string;
}>();

const emit = defineEmits<{
    (e: 'update:modelValue', value: string): void;
}>();

defineOptions({
    inheritAttrs: false,
});

const content = useVModel(props, 'modelValue', emit);

const mailerStore = useMailerStore();
const { selectedEmail } = storeToRefs(mailerStore);

const { data: templates } = useLazyAsyncData(`mailer-templates:${selectedEmail.value!.id}`, () => listTemplates());

async function listTemplates(): Promise<ListTemplatesResponse> {
    try {
        const call = getGRPCMailerClient().listTemplates({
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
            v-bind="$attrs"
            :options="templates?.templates"
            option-attribute="title"
            by="id"
            searchable
            class="min-w-48"
            :placeholder="$t('common.template')"
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

            <UButton :label="$t('common.cancel')" color="red" icon="i-mdi-cancel" @click="selectedTemplate = undefined" />
        </UButtonGroup>
    </ClientOnly>
</template>
