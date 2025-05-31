<script lang="ts" setup>
import ClipboardCitizens from '~/components/clipboard/modal/ClipboardCitizens.vue';
import ClipboardDocuments from '~/components/clipboard/modal/ClipboardDocuments.vue';
import ClipboardVehicles from '~/components/clipboard/modal/ClipboardVehicles.vue';
import TemplateRequirementsList from '~/components/documents/templates/TemplateRequirementsList.vue';
import TemplatesList from '~/components/documents/templates/TemplatesList.vue';
import { useClipboardStore } from '~/stores/clipboard';
import type { TemplateRequirements, TemplateShort } from '~~/gen/ts/resources/documents/templates';

const clipboardStore = useClipboardStore();

const props = withDefaults(
    defineProps<{
        autoFill?: boolean;
    }>(),
    {
        autoFill: false,
    },
);

const { isOpen } = useModal();

const { can } = useAuth();

const template = ref<undefined | TemplateShort>();
const reqs = ref<undefined | TemplateRequirements>();

const steps = ref<{ selectTemplate: boolean; selectClipboard: boolean }>({
    selectTemplate: true,
    selectClipboard: false,
});

const reqStatus = ref<{
    documents: boolean;
    users: boolean;
    vehicles: boolean;
}>({ documents: false, users: false, vehicles: false });

const readyToCreate = ref(false);

watch(reqStatus.value, () => {
    readyToCreate.value = reqStatus.value.documents && reqStatus.value.users && reqStatus.value.vehicles;
    // Auto redirect users when the requirements are matched
    if (readyToCreate.value && props.autoFill) {
        clipboardDialog();
    }
});

const documentsDocuments = useDocumentsDocuments();

function closeDialog(): void {
    isOpen.value = false;
}

async function templateSelected(t: TemplateShort | undefined): Promise<void> {
    if (t) {
        template.value = t;
        if (t.schema) {
            reqs.value = t.schema?.requirements;
            let reqDocuments = false;
            let reqUsers = false;
            let reqVehicles = false;

            clipboardStore.clearActiveStack();
            if (reqs.value) {
                if (reqs.value.documents) {
                    reqDocuments = clipboardStore.checkRequirements(reqs.value.documents, 'documents');
                    if (reqDocuments) {
                        clipboardStore.promoteToActiveStack('documents');
                    }
                } else {
                    reqDocuments = true;
                }
                if (reqs.value.users) {
                    reqUsers = clipboardStore.checkRequirements(reqs.value.users, 'users');
                    if (reqUsers) {
                        clipboardStore.promoteToActiveStack('users');
                    }
                } else {
                    reqUsers = true;
                }
                if (reqs.value.vehicles) {
                    reqVehicles = clipboardStore.checkRequirements(reqs.value.vehicles, 'vehicles');
                    if (reqVehicles) {
                        clipboardStore.promoteToActiveStack('vehicles');
                    }
                } else {
                    reqVehicles = true;
                }
            } else {
                reqDocuments = true;
                reqUsers = true;
                reqVehicles = true;
            }

            reqStatus.value.documents = reqDocuments;
            reqStatus.value.users = reqUsers;
            reqStatus.value.vehicles = reqVehicles;

            steps.value.selectTemplate = false;
            steps.value.selectClipboard = true;
        } else {
            await documentsDocuments.createDocument(template.value.id);
            isOpen.value = false;
        }
    } else {
        reqStatus.value.documents = false;
        reqStatus.value.users = false;
        reqStatus.value.vehicles = false;

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

    isOpen.value = false;
}

onBeforeMount(async () => {
    if (!can('documents.DocumentsService.CreateDocument').value) {
        await documentsDocuments.createDocument();
    }
});
</script>

<template>
    <UModal :ui="{ width: 'w-full sm:max-w-5xl' }">
        <UCard :ui="{ ring: '', divide: 'divide-y divide-gray-100 dark:divide-gray-800' }">
            <template #header>
                <div class="flex items-center justify-between">
                    <h3 class="text-2xl font-semibold leading-6">
                        {{ $t('common.template', 2) }}
                        <template v-if="template">- {{ template.title }} </template>
                    </h3>

                    <UButton class="-my-1" color="gray" variant="ghost" icon="i-mdi-window-close" @click="closeDialog()" />
                </div>
            </template>

            <template v-if="steps.selectTemplate">
                <UButton block @click="clipboardDialog()">
                    {{ $t('components.documents.templates.templates_modal.no_template') }}
                </UButton>

                <div class="pt-6">
                    <TemplatesList @selected="templateSelected($event)" />
                </div>
            </template>
            <div v-else-if="template !== undefined && reqs !== undefined && steps.selectClipboard">
                <div>
                    <div v-if="reqs.users">
                        <ClipboardCitizens
                            v-model:submit="submit"
                            :specs="reqs.users!"
                            @statisfied="(v: boolean) => (reqStatus.users = v)"
                            @close="closeDialog()"
                        >
                            <template #header>
                                <span class="text-sm">
                                    <TemplateRequirementsList
                                        :name="$t('common.citizen', 2)"
                                        :plural="$t('common.citizen', 2)"
                                        :specs="reqs.users!"
                                    />
                                </span>
                            </template>
                        </ClipboardCitizens>
                    </div>

                    <div v-if="reqs.vehicles">
                        <ClipboardVehicles
                            v-model:submit="submit"
                            :specs="reqs.vehicles!"
                            @statisfied="(v: boolean) => (reqStatus.vehicles = v)"
                            @close="closeDialog()"
                        >
                            <template #header>
                                <span class="text-sm">
                                    <TemplateRequirementsList
                                        :name="$t('common.vehicle', 2)"
                                        :plural="$t('common.vehicle', 2)"
                                        :specs="reqs.vehicles!"
                                    />
                                </span>
                            </template>
                        </ClipboardVehicles>
                    </div>

                    <div v-if="reqs.documents">
                        <ClipboardDocuments
                            v-model:submit="submit"
                            :specs="reqs.documents!"
                            @statisfied="(v: boolean) => (reqStatus.documents = v)"
                            @close="closeDialog()"
                        >
                            <template #header>
                                <span class="text-sm">
                                    <TemplateRequirementsList
                                        :name="$t('common.document', 2)"
                                        :plural="$t('common.document', 2)"
                                        :specs="reqs.documents!"
                                    />
                                </span>
                            </template>
                        </ClipboardDocuments>
                    </div>
                </div>
            </div>

            <template #footer>
                <UButtonGroup
                    v-if="template !== undefined && reqs !== undefined && steps.selectClipboard"
                    class="inline-flex w-full"
                >
                    <UButton class="flex-1" color="black" block @click="goBackDialog">
                        {{ $t('common.go_back') }}
                    </UButton>

                    <UButton class="flex-1" block :disabled="!readyToCreate" @click="clipboardDialog()">
                        {{ $t('common.create') }}
                    </UButton>
                </UButtonGroup>

                <UButton v-else class="flex-1" color="black" block @click="closeDialog">
                    {{ $t('common.close', 1) }}
                </UButton>
            </template>
        </UCard>
    </UModal>
</template>
