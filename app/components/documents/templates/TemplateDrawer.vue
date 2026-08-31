<script lang="ts" setup>
import { z } from 'zod';
import ClipboardCitizens from '~/components/clipboard/modal/ClipboardCitizens.vue';
import ClipboardDocuments from '~/components/clipboard/modal/ClipboardDocuments.vue';
import ClipboardVehicles from '~/components/clipboard/modal/ClipboardVehicles.vue';
import List from '~/components/documents/templates/List.vue';
import RequirementsList from '~/components/documents/templates/RequirementsList.vue';
import { useClipboardStore } from '~/stores/clipboard';
import type { ObjectSpecs, TemplateRequirements, TemplateShort } from '~~/gen/ts/resources/documents/templates/templates';

const clipboardStore = useClipboardStore();

const emits = defineEmits<{
    (e: 'close', v: boolean): void;
}>();

const template = ref<undefined | TemplateShort>();
const reqs = ref<undefined | TemplateRequirements>();

const steps = ref<{ selectTemplate: boolean; selectClipboard: boolean }>({
    selectTemplate: true,
    selectClipboard: false,
});

const requirementDefinitions = [
    { type: 'citizens', key: 'users', name: 'citizen' },
    { type: 'documents', key: 'documents', name: 'document' },
    { type: 'vehicles', key: 'vehicles', name: 'vehicle' },
] as const;

type RequirementType = (typeof requirementDefinitions)[number]['type'];

const reqStatus = ref<Record<RequirementType, boolean>>({
    citizens: false,
    documents: false,
    vehicles: false,
});

const readyToCreate = computed(() =>
    requirementDefinitions.every(({ type, key }) => !hasRequirement(reqs.value?.[key]) || reqStatus.value[type]),
);

const documentsDocuments = await useDocumentsDocuments();

function clipboardComponent(type: RequirementType) {
    switch (type) {
        case 'citizens':
            return ClipboardCitizens;
        case 'vehicles':
            return ClipboardVehicles;
        case 'documents':
            return ClipboardDocuments;
    }
}

function hasRequirement(specs?: ObjectSpecs): boolean {
    return !!specs && (specs.required === true || (specs.min ?? 0) > 0 || (specs.max ?? 0) > 0);
}

async function selectTemplate(t?: TemplateShort | undefined): Promise<void> {
    if (t) {
        template.value = t;
        const requirements = t.schema?.requirements;

        if (!requirements) {
            await documentsDocuments.createDocument(template.value.id);
            emits('close', false);
            return;
        }

        reqs.value = requirements;
        reqStatus.value = {
            citizens: false,
            documents: false,
            vehicles: false,
        };
        clipboardStore.clearActiveStack();
        requirementDefinitions.forEach(({ type, key }) => {
            const specs = requirements[key];
            if (hasRequirement(specs)) {
                reqStatus.value[type] = clipboardStore.checkRequirements(specs!, type);
                clipboardStore.promoteToActiveStack(type);
            } else {
                reqStatus.value[type] = true;
            }
        });
        steps.value.selectTemplate = false;
        steps.value.selectClipboard = true;
    } else {
        requirementDefinitions.forEach(({ type }) => {
            reqStatus.value[type] = false;
        });
        template.value = undefined;
        reqs.value = undefined;
    }
}

function goBackDialog(): void {
    steps.value.selectTemplate = true;
    steps.value.selectClipboard = false;

    nextTick(() => selectTemplate());
}

const submit = ref<boolean>(false);

async function clipboardDialog(): Promise<void> {
    submit.value = true;
    await documentsDocuments.createDocument(template.value?.id);

    emits('close', false);
}

const requirementTypes = computed(() => {
    if (!reqs.value) return [];
    return requirementDefinitions.filter(({ key }) => reqs.value?.[key] !== undefined);
});

const schema = z.object({
    title: z.string().max(128).optional(),
});

const query = useSearchForm('documents-templates', schema);
</script>

<template>
    <UDrawer
        :close="{ onClick: () => $emit('close', false) }"
        :ui="{ container: 'flex-1', content: 'min-h-[70%]', title: 'flex flex-row gap-2 justify-between', body: 'h-full' }"
    >
        <template #title>
            <div class="relative flex flex-1 items-center gap-2">
                <div class="flex items-center justify-center gap-2">
                    <UIcon
                        v-if="template"
                        class="shrink-0"
                        :class="`text-${template.color ?? 'primary'}`"
                        :name="template.icon ? convertComponentIconNameToDynamic(template.icon) : 'i-mdi-file-outline'"
                    />
                    <span>{{ template?.title ?? $t('common.template', 2) }}</span>
                </div>
                <span class="pointer-events-none absolute inset-x-0 text-center text-xs text-muted">
                    {{ $t(template ? 'common.select_clipboard_items' : 'common.select_template_or_blank') }}
                </span>
            </div>
        </template>

        <template #body>
            <div class="mx-auto w-full max-w-[80%] min-w-3/4">
                <template v-if="steps.selectTemplate">
                    <UButton
                        block
                        icon="i-mdi-plus"
                        :label="$t('components.documents.templates.templates_modal.no_template')"
                        @click="clipboardDialog()"
                    />

                    <USeparator class="my-4" />

                    <List :search-title="query.title" @selected="selectTemplate($event)" />
                </template>

                <div v-else-if="template !== undefined && reqs !== undefined && steps.selectClipboard">
                    <div>
                        <template v-for="(requirement, index) in requirementTypes" :key="requirement.type">
                            <component
                                :is="clipboardComponent(requirement.type)"
                                v-model:submit="submit"
                                :specs="reqs[requirement.key]!"
                                @statisfied="(v: boolean) => (reqStatus[requirement.type] = v)"
                                @close="$emit('close', false)"
                            >
                                <template #header>
                                    <span class="text-sm">
                                        <RequirementsList
                                            :name="$t('common.' + requirement.name, 2)"
                                            :plural="$t('common.' + requirement.name, 2)"
                                            :specs="reqs[requirement.key]!"
                                            :fulfilled="reqStatus[requirement.type]"
                                        />
                                    </span>
                                </template>
                            </component>

                            <USeparator v-if="index < requirementTypes.length - 1" class="my-2" />
                        </template>
                    </div>
                </div>
            </div>
        </template>

        <template #footer>
            <UForm
                v-if="steps.selectTemplate"
                class="mx-auto my-2 flex w-full max-w-[80%] min-w-3/4 flex-1 flex-col gap-2"
                :schema="schema"
                :state="query"
            >
                <UFormField class="flex-1" name="title">
                    <UInput
                        ref="inputRef"
                        v-model="query.title"
                        class="w-full"
                        type="text"
                        name="title"
                        :placeholder="$t('common.template')"
                        leading-icon="i-mdi-search"
                    />
                </UFormField>
            </UForm>

            <UFieldGroup
                v-if="template !== undefined && reqs !== undefined && steps.selectClipboard"
                class="inline-flex w-full"
            >
                <UButton class="flex-1" color="neutral" block :label="$t('common.go_back')" @click="goBackDialog()" />

                <UButton
                    class="flex-1"
                    block
                    :disabled="!readyToCreate"
                    :label="$t('common.create')"
                    @click="clipboardDialog()"
                />
            </UFieldGroup>

            <UButton
                v-else
                class="flex-1"
                color="neutral"
                block
                :label="$t('common.close', 1)"
                @click="$emit('close', false)"
            />
        </template>
    </UDrawer>
</template>
