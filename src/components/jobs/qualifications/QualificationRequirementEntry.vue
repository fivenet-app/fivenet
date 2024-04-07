<script lang="ts" setup>
import { Combobox, ComboboxButton, ComboboxInput, ComboboxOption, ComboboxOptions } from '@headlessui/vue';
import { CheckIcon, CloseIcon } from 'mdi-vue3';
import type { QualificationRequirement, QualificationShort } from '~~/gen/ts/resources/qualifications/qualifications';
import type { ListQualificationsResponse } from '~~/gen/ts/services/qualifications/qualifications';

const props = withDefaults(
    defineProps<{
        requirement: QualificationRequirement;
        readOnly?: boolean;
    }>(),
    {
        readOnly: false,
    },
);

const emits = defineEmits<{
    (e: 'update-qualification', qualification?: QualificationShort): void;
    (e: 'remove'): void;
}>();

const { $grpc } = useNuxtApp();

const queryQualificationRaw = ref('');
const queryQualification = computed(() => queryQualificationRaw.value.toLowerCase());

const { data, refresh } = useLazyAsyncData(`jobs-qualifications-0-${queryQualificationRaw.value}`, () => listQualifications());

async function listQualifications(): Promise<ListQualificationsResponse> {
    try {
        const call = $grpc.getQualificationsClient().listQualifications({
            pagination: {
                offset: 0,
            },
            search: queryQualification.value,
        });
        const { response } = await call;

        return response;
    } catch (e) {
        $grpc.handleError(e as RpcError);
        throw e;
    }
}

const filteredQualifications = computed(() => {
    const qualis =
        data.value?.qualifications.filter(
            (q) =>
                q.abbreviation.toLowerCase().includes(queryQualification.value) ||
                q.title.toLowerCase().includes(queryQualification.value),
        ) ?? [];

    if (selectedQualification.value && !qualis.find((q) => q.id === selectedQualification.value?.id)) {
        qualis.push({
            id: selectedQualification.value.id,
            job: selectedQualification.value.job,
            abbreviation: selectedQualification.value.abbreviation,
            title: selectedQualification.value.title,
            closed: selectedQualification.value.closed,
            content: '',
            creatorId: selectedQualification.value.creatorId,
            creatorJob: selectedQualification.value.creatorJob,
            requirements: selectedQualification.value.requirements,
            weight: selectedQualification.value.weight,
        });
    }

    return qualis;
});
const selectedQualification = ref<QualificationShort | undefined>(props.requirement.targetQualification);

watchDebounced(queryQualification, async () => refresh(), {
    debounce: 600,
    maxWait: 1750,
});

watch(selectedQualification, () => emits('update-qualification', selectedQualification.value));
</script>

<template>
    <div class="my-2 flex flex-row items-center">
        <Combobox v-model="selectedQualification" as="div" :disabled="readOnly" class="flex-1">
            <div class="relative">
                <ComboboxButton as="div">
                    <ComboboxInput
                        autocomplete="off"
                        :display-value="(qualification: any) => `${qualification?.abbreviation}: ${qualification?.title}`"
                        :class="readOnly ? 'disabled' : ''"
                        @change="queryQualificationRaw = $event.target.value"
                        @focusin="focusTablet(true)"
                        @focusout="focusTablet(false)"
                    />
                </ComboboxButton>

                <ComboboxOptions
                    class="absolute z-10 mt-1 max-h-44 w-full overflow-auto rounded-md bg-base-700 py-1 text-base sm:text-sm"
                >
                    <ComboboxOption
                        v-for="qualification in filteredQualifications"
                        :key="qualification.id"
                        v-slot="{ active, selected }"
                        :value="qualification"
                    >
                        <li :class="['relative cursor-default select-none py-2 pl-8 pr-4', active ? 'bg-primary-500' : '']">
                            <span :class="['block truncate', selected && 'font-semibold']">
                                {{ qualification.abbreviation }}: {{ qualification.title }}
                            </span>

                            <span
                                v-if="selected"
                                :class="[
                                    active ? 'text-neutral' : 'text-primary-500',
                                    'absolute inset-y-0 left-0 flex items-center pl-1.5',
                                ]"
                            >
                                <CheckIcon class="size-5" />
                            </span>
                        </li>
                    </ComboboxOption>
                </ComboboxOptions>
            </div>
        </Combobox>

        <UButton
            class="bg-primary-500 hover:bg-primary-400 ml-2 rounded-full p-1.5 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2"
            @click="$emit('remove')"
        >
            <CloseIcon class="size-5" />
        </UButton>
    </div>
</template>
