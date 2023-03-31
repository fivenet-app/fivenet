<script lang="ts" setup>
import { DocumentTemplateShort, TemplateRequirements } from '@fivenet/gen/resources/documents/templates_pb';
import { Dialog, DialogPanel, DialogTitle, TransitionChild, TransitionRoot } from '@headlessui/vue';
import { PencilIcon } from '@heroicons/vue/24/solid';
import { ref, watch } from 'vue';
import ClipboardModalDocuments from '../clipboard/ClipboardModalDocuments.vue';
import ClipboardModalUsers from '../clipboard/ClipboardModalUsers.vue';
import ClipboardModalVehicles from '../clipboard/ClipboardModalVehicles.vue';
import TemplatesList from './TemplatesList.vue';

const router = useRouter();

const props = defineProps({
    open: {
        required: true,
        type: Boolean,
    },
    autoFill: {
        required: false,
        type: Boolean,
        default: false,
    },
});

const emit = defineEmits<{
    (e: 'close'): void,
}>();

const template = ref<undefined | DocumentTemplateShort>();
const reqs = ref<undefined | TemplateRequirements>();

const steps = ref<{ selectTemplate: boolean, selectClipboard: boolean }>({
    selectTemplate: true,
    selectClipboard: false,
});

const reqStatus = ref<{
    documents: boolean,
    users: boolean,
    vehicles: boolean,
}>({ documents: false, users: false, vehicles: false });

const readyToCreate = ref(false);

watch(reqStatus.value, (v) => {
    readyToCreate.value = (reqStatus.value.documents && reqStatus.value.users && reqStatus.value.vehicles);

    // Auto redirect users when the requirements are matched
    if (readyToCreate && props.autoFill) {
        clipboardDialog();
    }
});

function closeDialog() {
    template.value = undefined;
    steps.value.selectTemplate = true;
    steps.value.selectClipboard = false;

    emit('close');
}

function templateSelected(t: DocumentTemplateShort) {
    if (t) {
        template.value = t;
        if (t.getSchema()) {
            reqs.value = t.getSchema()?.getRequirements();
            if (reqs.value) {
                reqStatus.value.documents = !reqs.value?.hasDocuments();
                reqStatus.value.users = !reqs.value?.hasUsers();
                reqStatus.value.vehicles = !reqs.value?.hasVehicles();
            } else {
                reqStatus.value.documents = false;
                reqStatus.value.users = false;
                reqStatus.value.vehicles = false;
            }

            steps.value.selectTemplate = false;
            steps.value.selectClipboard = true;
        } else {
            router.push({ name: 'documents-create', query: { templateId: template.value?.getId() } });
        }
    } else {
        reqStatus.value.documents = false;
        reqStatus.value.users = false;
        reqStatus.value.vehicles = false;

        template.value = undefined;
        reqs.value = undefined;
    }
}

function goBackDialog() {
    steps.value.selectTemplate = true;
    steps.value.selectClipboard = false;
}

const submit = ref(false);

function clipboardDialog() {
    submit.value = true;
    router.push({ name: 'documents-create', query: { templateId: template.value?.getId() } });
}
</script>

<template>
    <TransitionRoot as="template" :show="open">
        <Dialog as="div" class="relative z-10" @close="closeDialog">
            <TransitionChild as="template" enter="ease-out duration-300" enter-from="opacity-0" enter-to="opacity-100"
                leave="ease-in duration-200" leave-from="opacity-100" leave-to="opacity-0">
                <div class="fixed inset-0 bg-gray-500 bg-opacity-75 transition-opacity" />
            </TransitionChild>

            <div class="fixed inset-0 z-10 overflow-y-auto">
                <div class="flex min-h-full items-end justify-center p-4 text-center sm:items-center sm:p-0">
                    <TransitionChild as="template" enter="ease-out duration-300"
                        enter-from="opacity-0 translate-y-4 sm:translate-y-0 sm:scale-95"
                        enter-to="opacity-100 translate-y-0 sm:scale-100" leave="ease-in duration-200"
                        leave-from="opacity-100 translate-y-0 sm:scale-100"
                        leave-to="opacity-0 translate-y-4 sm:translate-y-0 sm:scale-95">
                        <div>
                            <DialogPanel
                                class="relative transform overflow-hidden rounded-lg bg-gray-800 px-4 pt-5 pb-4 text-left shadow-xl transition-all sm:my-8 sm:w-full sm:max-w-5xl sm:p-6">
                                <div v-if="steps.selectTemplate">
                                    <div>
                                        <div
                                            class="mx-auto flex h-12 w-12 items-center justify-center rounded-full bg-gray-100">
                                            <PencilIcon class="h-6 w-6 text-indigo-600" aria-hidden="true" />
                                        </div>
                                        <div class="mt-3 text-center sm:mt-5">
                                            <DialogTitle as="h3" class="text-base font-semibold leading-6 text-white">
                                                Document Templates
                                            </DialogTitle>
                                            <div class="mt-2 text-white">
                                                <NuxtLink :to="{ name: 'documents-create' }" type="button"
                                                    class="mb-5 mt-3 inline-flex w-full justify-center rounded-md bg-white px-3 py-2 text-sm font-semibold text-gray-900 shadow-sm ring-1 ring-inset ring-gray-300 hover:bg-gray-50 sm:col-start-1 sm:mt-0">
                                                    No Template
                                                </NuxtLink>
                                                <TemplatesList
                                                    @selected="(t: DocumentTemplateShort) => templateSelected(t)" />
                                            </div>
                                        </div>
                                    </div>
                                    <div class="mt-5 sm:mt-6 sm:grid sm:grid-flow-row-dense sm:gap-3">
                                        <button type="button"
                                            class="mt-3 inline-flex w-full justify-center rounded-md bg-white px-3 py-2 text-sm font-semibold text-gray-900 shadow-sm ring-1 ring-inset ring-gray-300 hover:bg-gray-50 sm:col-start-1 sm:mt-0"
                                            @click="closeDialog" ref="cancelButtonRef">
                                            Close
                                        </button>
                                    </div>
                                </div>
                                <div v-else-if="!autoFill && steps.selectClipboard && template && reqs">
                                    <div
                                        class="mx-auto flex h-12 w-12 items-center justify-center rounded-full bg-gray-100">
                                        <PencilIcon class="h-6 w-6 text-indigo-600" aria-hidden="true" />
                                    </div>
                                    <div class="mt-3 text-center sm:mt-5">
                                        <DialogTitle as="h3" class="text-base font-semibold leading-6 text-white">
                                            Template: {{ template.getTitle() }}
                                        </DialogTitle>
                                        <div class="mt-2 text-white">
                                            <div v-if="reqs.getUsers()">
                                                <p>
                                                    <span>
                                                        <span v-if="reqs.getUsers()?.getRequired()">
                                                            Requires
                                                        </span>
                                                        <span v-if="reqs.getUsers()?.getMax() == reqs.getUsers()?.getMin()">
                                                            1 User
                                                        </span>
                                                        <span v-else>
                                                            {{ (reqs.getUsers()?.getMin() === 0 &&
                                                                reqs.getUsers()?.getRequired()) ? reqs.getUsers()?.getMax() :
                                                                reqs.getUsers()?.getMin() }} User(s)
                                                            <span v-if="reqs.getUsers()?.getMax()! > 0">&nbsp;(Max: {{
                                                                reqs.getUsers()?.getMax() }})</span>
                                                        </span>
                                                    </span>
                                                </p>

                                                <ClipboardModalUsers :submit.sync="submit" :showSelect="true"
                                                    :min="reqs.getUsers()?.getMin()" :max="reqs.getUsers()?.getMax()"
                                                    @statisfied="(v) => reqStatus.users = v" />
                                            </div>
                                            <div v-if="reqs.getVehicles()">
                                                <p>
                                                    <span>
                                                        <span v-if="reqs.getVehicles()?.getRequired()">
                                                            Requires
                                                        </span>
                                                        <span
                                                            v-if="reqs.getVehicles()?.getMax() == reqs.getVehicles()?.getMin()">
                                                            1 Vehicle
                                                        </span>
                                                        <span v-else>
                                                            {{ (reqs.getVehicles()?.getMin() === 0 &&
                                                                reqs.getVehicles()?.getRequired()) ?
                                                                reqs.getVehicles()?.getMax() : reqs.getVehicles()?.getMin() }}
                                                            Vehicle(s)
                                                            <span v-if="reqs.getVehicles()?.getMax()! > 0">&nbsp;(Max: {{
                                                                reqs.getVehicles()?.getMax() }})</span>
                                                        </span>
                                                    </span>
                                                </p>

                                                <ClipboardModalVehicles :submit.sync="submit" :showSelect="true"
                                                    :min="reqs.getVehicles()?.getMin()" :max="reqs.getVehicles()?.getMax()"
                                                    @statisfied="(v) => reqStatus.vehicles = v" />
                                            </div>
                                            <div v-if="reqs.getDocuments()">
                                                <p>
                                                    <span>
                                                        <span v-if="reqs.getDocuments()?.getRequired()">
                                                            Requires
                                                        </span>
                                                        <span
                                                            v-if="reqs.getDocuments()?.getMax() == reqs.getDocuments()?.getMin()">
                                                            1 User
                                                        </span>
                                                        <span v-else>
                                                            {{ (reqs.getDocuments()?.getMin() === 0 &&
                                                                reqs.getDocuments()?.getRequired()) ?
                                                                reqs.getDocuments()?.getMax() : reqs.getDocuments()?.getMin() }}
                                                            User(s)
                                                            <span v-if="reqs.getDocuments()?.getMax()! > 0">&nbsp;(Max: {{
                                                                reqs.getDocuments()?.getMax() }})</span>
                                                        </span>
                                                    </span>
                                                </p>

                                                <ClipboardModalDocuments :submit.sync="submit" :showSelect="true"
                                                    :min="reqs.getDocuments()?.getMin()"
                                                    :max="reqs.getDocuments()?.getMax()"
                                                    @statisfied="(v) => reqStatus.documents = v" />
                                            </div>
                                        </div>
                                    </div>
                                    <div class="mt-5 sm:mt-6 sm:grid sm:grid-flow-row-dense sm:gap-3">
                                        <button type="button"
                                            class="mt-3 inline-flex w-full justify-center rounded-md bg-white px-3 py-2 text-sm font-semibold text-gray-900 shadow-sm ring-1 ring-inset ring-gray-300 hover:bg-gray-50 sm:col-start-1 sm:mt-0"
                                            @click="goBackDialog">
                                            Go Back
                                        </button>
                                        <button type="button"
                                            class="inline-flex w-full justify-center rounded-md bg-indigo-600 px-3 py-2 text-sm font-semibold text-white shadow-sm hover:bg-indigo-500 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-indigo-600 sm:col-start-2"
                                            @click="clipboardDialog()" :disabled="!readyToCreate">
                                            Create
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
