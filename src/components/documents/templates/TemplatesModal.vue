<script lang="ts" setup>
import ClipboardCitizens from '~/components/clipboard/modal/ClipboardCitizens.vue';
import ClipboardDocuments from '~/components/clipboard/modal/ClipboardDocuments.vue';
import ClipboardVehicles from '~/components/clipboard/modal/ClipboardVehicles.vue';
import TemplateRequirementsList from '~/components/documents/templates/TemplateRequirementsList.vue';
import TemplatesList from '~/components/documents/templates/TemplatesList.vue';
import { useClipboardStore } from '~/store/clipboard';
import { TemplateRequirements, TemplateShort } from '~~/gen/ts/resources/documents/templates';

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

function closeDialog(): void {
    template.value = undefined;
    steps.value.selectTemplate = true;
    steps.value.selectClipboard = false;

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
            await navigateTo({
                name: 'documents-create',
                query: { templateId: template.value?.id },
            });
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
    await navigateTo({
        name: 'documents-create',
        query: { templateId: template.value?.id },
    });

    isOpen.value = false;
}
</script>

<template>
    <UModal :ui="{ width: 'w-full sm:max-w-5xl' }">
        <UCard :ui="{ ring: '', divide: 'divide-y divide-gray-100 dark:divide-gray-800' }">
            <template #header>
                <div class="flex items-center justify-between">
                    <h3 class="text-2xl font-semibold leading-6">
                        {{ $t('common.document', 1) }}
                        {{ $t('common.template', 2) }}
                    </h3>

                    <UButton color="gray" variant="ghost" icon="i-mdi-window-close" class="-my-1" @click="closeDialog()" />
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
            <template v-else-if="template !== undefined && reqs !== undefined && steps.selectClipboard">
                <div class="text-center">
                    <h3 class="text-base font-semibold leading-6">
                        {{ $t('common.template', 1) }}:
                        {{ template.title }}
                    </h3>
                    <div>
                        <div v-if="reqs.users">
                            <p>
                                <TemplateRequirementsList :name="$t('common.citizen', 2)" :specs="reqs.users!" />
                            </p>

                            <ClipboardCitizens
                                :submit.sync="submit"
                                :specs="reqs.users!"
                                @statisfied="(v: boolean) => (reqStatus.users = v)"
                                @close="closeDialog()"
                            />
                        </div>
                        <div v-if="reqs.vehicles">
                            <p>
                                <TemplateRequirementsList :name="$t('common.vehicle', 2)" :specs="reqs.vehicles!" />
                            </p>

                            <ClipboardVehicles
                                :submit.sync="submit"
                                :specs="reqs.vehicles!"
                                @statisfied="(v: boolean) => (reqStatus.vehicles = v)"
                                @close="closeDialog()"
                            />
                        </div>
                        <div v-if="reqs.documents">
                            <p>
                                <TemplateRequirementsList :name="$t('common.document', 2)" :specs="reqs.documents!" />
                            </p>

                            <ClipboardDocuments
                                :submit.sync="submit"
                                :specs="reqs.documents!"
                                @statisfied="(v: boolean) => (reqStatus.documents = v)"
                                @close="closeDialog()"
                            />
                        </div>
                    </div>
                </div>
            </template>

            <template #footer>
                <UButtonGroup
                    v-if="template !== undefined && reqs !== undefined && steps.selectClipboard"
                    class="inline-flex w-full"
                >
                    <UButton color="black" block class="flex-1" @click="goBackDialog">
                        {{ $t('common.go_back') }}
                    </UButton>
                    <UButton block class="flex-1" :disabled="!readyToCreate" @click="clipboardDialog()">
                        {{ $t('common.create') }}
                    </UButton>
                </UButtonGroup>

                <UButton v-else color="black" block class="flex-1" @click="closeDialog">
                    {{ $t('common.close', 1) }}
                </UButton>
            </template>
        </UCard>
    </UModal>
</template>
