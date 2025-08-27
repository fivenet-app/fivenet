<script lang="ts" setup>
import ClipboardCitizens from '~/components/clipboard/modal/ClipboardCitizens.vue';
import ClipboardDocuments from '~/components/clipboard/modal/ClipboardDocuments.vue';
import ClipboardVehicles from '~/components/clipboard/modal/ClipboardVehicles.vue';
import TemplateList from '~/components/documents/templates/TemplateList.vue';
import TemplateRequirementsList from '~/components/documents/templates/TemplateRequirementsList.vue';
import { useClipboardStore } from '~/stores/clipboard';
import type { TemplateRequirements, TemplateShort } from '~~/gen/ts/resources/documents/templates';

const clipboardStore = useClipboardStore();

const emit = defineEmits<{
    (e: 'close', v: boolean): void;
}>();

const template = ref<undefined | TemplateShort>();
const reqs = ref<undefined | TemplateRequirements>();

const steps = ref<{ selectTemplate: boolean; selectClipboard: boolean }>({
    selectTemplate: true,
    selectClipboard: false,
});

const requirementTypes = ['citizens', 'documents', 'vehicles'] as const;
type RequirementType = (typeof requirementTypes)[number];

const reqStatus = ref<Record<RequirementType, boolean>>({
    citizens: false,
    documents: false,
    vehicles: false,
});

const readyToCreate = ref(false);

watch(reqStatus.value, () => {
    readyToCreate.value = requirementTypes.every((type) => reqStatus.value[type]);
});

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

async function templateSelected(t: TemplateShort | undefined): Promise<void> {
    if (t) {
        template.value = t;
        if (t.schema) {
            reqs.value = t.schema?.requirements;
            clipboardStore.clearActiveStack();
            requirementTypes.forEach((type) => {
                const required = reqs.value?.[type === 'citizens' ? 'users' : type];
                let status = true;
                if (required) {
                    clipboardStore.promoteToActiveStack(type);
                    status = clipboardStore.checkRequirements(required, type);
                }
                reqStatus.value[type] = status;
            });
            steps.value.selectTemplate = false;
            steps.value.selectClipboard = true;
        } else {
            await documentsDocuments.createDocument(template.value.id);
            emit('close', false);
        }
    } else {
        requirementTypes.forEach((type) => {
            reqStatus.value[type] = false;
        });
        template.value = undefined;
        reqs.value = undefined;
    }
}

function goBackDialog(): void {
    steps.value.selectTemplate = true;
    steps.value.selectClipboard = false;
}

const submit = ref(false);

async function clipboardDialog(): Promise<void> {
    submit.value = true;
    await documentsDocuments.createDocument(template.value?.id);

    emit('close', false);
}

const filteredRequirementTypes = computed(() => {
    if (!reqs.value) return [];
    return requirementTypes.filter((type) => reqs.value && reqs.value[type === 'citizens' ? 'users' : type]);
});
</script>

<template>
    <UModal :title="`${$t('common.template', 2)}${template ? ` - ${template?.title}` : ''}`">
        <template #body>
            <template v-if="steps.selectTemplate">
                <UButton block @click="clipboardDialog()">
                    {{ $t('components.documents.templates.templates_modal.no_template') }}
                </UButton>

                <div class="pt-6">
                    <TemplateList @selected="templateSelected($event)" />
                </div>
            </template>
            <div v-else-if="template !== undefined && reqs !== undefined && steps.selectClipboard">
                <div>
                    <div v-for="type in filteredRequirementTypes" :key="type">
                        <component
                            :is="clipboardComponent(type)"
                            v-model:submit="submit"
                            :specs="reqs[type === 'citizens' ? 'users' : type]!"
                            @statisfied="(v: boolean) => (reqStatus[type] = v)"
                            @close="$emit('close', false)"
                        >
                            <template #header>
                                <span class="text-sm">
                                    <TemplateRequirementsList
                                        :name="$t('common.' + type.slice(0, -1), 2)"
                                        :plural="$t('common.' + type.slice(0, -1), 2)"
                                        :specs="reqs[type === 'citizens' ? 'users' : type]!"
                                    />
                                </span>
                            </template>
                        </component>
                    </div>
                </div>
            </div>
        </template>

        <template #footer>
            <UButtonGroup
                v-if="template !== undefined && reqs !== undefined && steps.selectClipboard"
                class="inline-flex w-full"
            >
                <UButton class="flex-1" color="neutral" block @click="goBackDialog">
                    {{ $t('common.go_back') }}
                </UButton>

                <UButton class="flex-1" block :disabled="!readyToCreate" @click="clipboardDialog()">
                    {{ $t('common.create') }}
                </UButton>
            </UButtonGroup>

            <UButton v-else class="flex-1" color="neutral" block @click="$emit('close', false)">
                {{ $t('common.close', 1) }}
            </UButton>
        </template>
    </UModal>
</template>
