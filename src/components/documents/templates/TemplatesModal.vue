<script lang="ts" setup>
import { Dialog, DialogPanel, DialogTitle, TransitionChild, TransitionRoot } from '@headlessui/vue';
import { CloseIcon, PencilIcon } from 'mdi-vue3';
import Citizens from '~/components/clipboard/modal/Citizens.vue';
import Documents from '~/components/clipboard/modal/Documents.vue';
import Vehicles from '~/components/clipboard/modal/Vehicles.vue';
import { useClipboardStore } from '~/store/clipboard';
import { TemplateRequirements, TemplateShort } from '~~/gen/ts/resources/documents/templates';
import List from '~/components/documents/templates/List.vue';
import TemplateRequirementsList from '~/components/documents/templates/TemplateRequirementsList.vue';

const clipboardStore = useClipboardStore();

const props = withDefaults(
    defineProps<{
        open: boolean;
        autoFill?: boolean;
    }>(),
    {
        autoFill: false,
    },
);

const emit = defineEmits<{
    (e: 'close'): void;
}>();

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

    emit('close');
}

async function templateSelected(t: TemplateShort): Promise<void> {
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
}
</script>

<template>
    <TransitionRoot as="template" :show="open">
        <Dialog as="div" class="relative z-30" @close="closeDialog">
            <TransitionChild
                as="template"
                enter="ease-out duration-300"
                enter-from="opacity-0"
                enter-to="opacity-100"
                leave="ease-in duration-200"
                leave-from="opacity-100"
                leave-to="opacity-0"
            >
                <div class="fixed inset-0 transition-opacity bg-opacity-75 bg-base-900" />
            </TransitionChild>

            <div class="fixed inset-0 z-30 overflow-y-auto">
                <div class="flex min-h-full items-end justify-center p-4 text-center sm:items-center sm:p-0">
                    <TransitionChild
                        as="template"
                        enter="ease-out duration-300"
                        enter-from="opacity-0 translate-y-4 sm:translate-y-0 sm:scale-95"
                        enter-to="opacity-100 translate-y-0 sm:scale-100"
                        leave="ease-in duration-200"
                        leave-from="opacity-100 translate-y-0 sm:scale-100"
                        leave-to="opacity-0 translate-y-4 sm:translate-y-0 sm:scale-95"
                    >
                        <div>
                            <DialogPanel
                                class="relative px-4 pt-5 pb-4 overflow-hidden text-left transition-all transform rounded-lg bg-base-800 text-neutral sm:my-8 w-full sm:max-w-7xl sm:min-w-full sm:p-6"
                            >
                                <div class="absolute right-0 top-0 pr-4 pt-4 block">
                                    <button
                                        type="button"
                                        class="rounded-md bg-neutral text-gray-400 hover:text-gray-500 focus:outline-none focus:ring-2 focus:ring-primary-500 focus:ring-offset-2"
                                        @click="$emit('close')"
                                    >
                                        <span class="sr-only">{{ $t('common.close') }}</span>
                                        <CloseIcon class="h-6 w-6" aria-hidden="true" />
                                    </button>
                                </div>
                                <div v-if="steps.selectTemplate">
                                    <div>
                                        <div
                                            class="mx-auto flex h-12 w-12 items-center justify-center rounded-full bg-base-700"
                                        >
                                            <PencilIcon class="h-6 w-6 text-primary-500" aria-hidden="true" />
                                        </div>
                                        <div class="mt-3 text-center sm:mt-5">
                                            <DialogTitle as="h3" class="text-base font-semibold leading-6">
                                                {{ $t('common.document', 1) }}
                                                {{ $t('common.template', 2) }}
                                            </DialogTitle>
                                            <div class="mt-2">
                                                <NuxtLink
                                                    :to="{
                                                        name: 'documents-create',
                                                    }"
                                                    type="button"
                                                    class="w-full rounded-md bg-base-500 py-2.5 px-3.5 text-sm font-semibold text-neutral hover:bg-base-400"
                                                >
                                                    {{ $t('components.documents.templates.templates_modal.no_template') }}
                                                </NuxtLink>
                                                <div class="pt-4">
                                                    <List @selected="(t: TemplateShort) => templateSelected(t)" />
                                                </div>
                                            </div>
                                        </div>
                                    </div>
                                    <div class="gap-2 mt-5 sm:mt-4 sm:flex">
                                        <button
                                            type="button"
                                            class="flex-1 rounded-md bg-base-500 py-2.5 px-3.5 text-sm font-semibold text-neutral hover:bg-base-400"
                                            @click="closeDialog"
                                        >
                                            {{ $t('common.close', 1) }}
                                        </button>
                                    </div>
                                </div>
                                <div v-else-if="template !== undefined && reqs !== undefined && steps.selectClipboard">
                                    <div class="mx-auto flex h-12 w-12 items-center justify-center rounded-full bg-gray-100">
                                        <PencilIcon class="h-6 w-6 text-primary-600" aria-hidden="true" />
                                    </div>
                                    <div class="mt-3 text-center sm:mt-5">
                                        <DialogTitle as="h3" class="text-base font-semibold leading-6 text-neutral">
                                            {{ $t('common.template', 1) }}:
                                            {{ template.title }}
                                        </DialogTitle>
                                        <div class="mt-2 text-neutral">
                                            <div v-if="reqs.users">
                                                <p>
                                                    <TemplateRequirementsList
                                                        :name="$t('common.citizen', 2)"
                                                        :specs="reqs.users!"
                                                    />
                                                </p>

                                                <Citizens
                                                    :submit.sync="submit"
                                                    :specs="reqs.users!"
                                                    @statisfied="(v: boolean) => (reqStatus.users = v)"
                                                />
                                            </div>
                                            <div v-if="reqs.vehicles">
                                                <p>
                                                    <TemplateRequirementsList
                                                        :name="$t('common.vehicle', 2)"
                                                        :specs="reqs.vehicles!"
                                                    />
                                                </p>

                                                <Vehicles
                                                    :submit.sync="submit"
                                                    :specs="reqs.vehicles!"
                                                    @statisfied="(v: boolean) => (reqStatus.vehicles = v)"
                                                />
                                            </div>
                                            <div v-if="reqs.documents">
                                                <p>
                                                    <TemplateRequirementsList
                                                        :name="$t('common.document', 2)"
                                                        :specs="reqs.documents!"
                                                    />
                                                </p>

                                                <Documents
                                                    :submit.sync="submit"
                                                    :specs="reqs.documents!"
                                                    @statisfied="(v: boolean) => (reqStatus.documents = v)"
                                                />
                                            </div>
                                        </div>
                                    </div>
                                    <div class="mt-5 sm:mt-6 sm:grid sm:grid-flow-row-dense sm:gap-3">
                                        <button
                                            type="button"
                                            class="mt-3 inline-flex w-full justify-center rounded-md bg-neutral px-3 py-2 text-sm font-semibold text-gray-900 shadow-sm ring-1 ring-inset ring-gray-300 hover:bg-gray-50 sm:col-start-1 sm:mt-0"
                                            @click="goBackDialog"
                                        >
                                            {{ $t('common.go_back') }}
                                        </button>
                                        <button
                                            type="button"
                                            class="inline-flex w-full justify-center rounded-md px-3 py-2 text-sm font-semibold text-neutral shadow-sm focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 sm:col-start-2"
                                            :disabled="!readyToCreate"
                                            :class="[
                                                !readyToCreate
                                                    ? 'disabled bg-base-500 hover:bg-base-400 focus-visible:outline-base-500'
                                                    : 'bg-primary-500 hover:bg-primary-400 focus-visible:outline-primary-500',
                                            ]"
                                            @click="clipboardDialog()"
                                        >
                                            {{ $t('common.create') }}
                                        </button>
                                    </div>
                                </div>
                            </DialogPanel>
                        </div>
                    </TransitionChild>
                </div>
            </div>
        </Dialog>
    </TransitionRoot>
</template>
