<script lang="ts" setup>
import { Disclosure, DisclosureButton, DisclosurePanel, Combobox, ComboboxButton, ComboboxInput, ComboboxOption, ComboboxOptions } from '@headlessui/vue'
import { AttributeValues, JobGradeList, RoleAttribute, StringList } from '@fivenet/gen/resources/permissions/permissions_pb';
import { ChevronDownIcon, CheckIcon } from '@heroicons/vue/24/solid';
import { Job } from '@fivenet/gen/resources/jobs/jobs_pb';
import { Map as protobufMap } from 'google-protobuf';

const props = defineProps<{
    attribute: RoleAttribute,
    states: Map<number, AttributeValues | undefined>,
    jobs: Job[],
    disabled?: boolean,
}>();

const emit = defineEmits<{
    (e: 'update:states', payload: Map<number, AttributeValues | undefined>): void,
}>();

const states = ref<typeof props.states>(props.states);
const id = ref<number>(props.attribute.getAttrId());

const validValues = ref<AttributeValues | undefined>(props.attribute.getValidValues());
const type = ref<string | undefined>(props.attribute.getType());

const tmp = ref<{ job: string, grade: number }[]>([]);

function getState(): AttributeValues {
    if (!states.value.has(id.value)) states.value.set(id.value, new AttributeValues());

    return states.value.get(id.value)!;
}

async function toggleListValue(value: string): Promise<void> {
    const state = getState();
    const list = state.getStringList() ?? new StringList();
    const array = list.getStringsList();

    if (array.indexOf(value) < 0) {
        array.push(value);
    } else {
        array.splice(array.indexOf(value), 1);
    }

    list.setStringsList(array);
    state.setStringList(list);
    states.value.set(id.value, state);
    emit('update:states', states.value);
}

async function updateJobGradeValue(job: string, grade: number): Promise<void> {
    const state = getState();
    const list = state.getJobGradeList() ?? new JobGradeList();
    const map = list.getJobsMap();

    if (grade === 0) {
        map.del(job);
    } else {
        map.set(job, grade -1);
    }

    state.setJobGradeList(list);
    states.value.set(id.value, state);
    emit('update:states', states.value);
}

onMounted(() => {
    props.jobs.forEach(job => {
        tmp.value.push({
            job: job.getName(),
            grade: 0,
        })
    })
});
</script>

<style scoped>
.upsidedown {
    transform: rotate(180deg);
}
</style>

<template>
    <div v-if="$props.attribute">
        <Disclosure as="div"
            :class="[$props.disabled ? 'border-neutral/10 text-base-300' : 'hover:border-neutral/70 border-neutral/20 text-neutral']"
            v-slot="{ open }">
            <DisclosureButton :disabled="$props.disabled"
                :class="[open ? 'rounded-t-lg border-b-0' : 'rounded-lg', $props.disabled ? 'cursor-not-allowed' : '', ' flex w-full items-start justify-between text-left border-2 p-2 border-inherit transition-colors']">
                <span class="text-base leading-7 transition-colors">
                    Options
                </span>
                <span class="ml-6 flex h-7 items-center">
                    <ChevronDownIcon :class="[open ? 'upsidedown' : '', 'h-6 w-6 transition-transform']"
                        aria-hidden="true" />
                </span>
            </DisclosureButton>
            <DisclosurePanel class="px-4 pb-2 border-2 border-t-0 rounded-b-lg transition-colors border-inherit -mt-2">
                <div class="flex flex-col gap-2 max-w-4xl mx-auto my-2">
                    <div v-if="type === 'StringList' && validValues" class="flex flex-row gap-4 flex-wrap">
                        <div v-for="value in validValues.getStringList()?.getStringsList()" :key="value"
                            class="flex flex-row flex-initial flex-nowrap">
                            <input :id="value" :name="value" type="checkbox"
                                :checked="!!getState().getStringList()?.getStringsList().find(v => v === value)"
                                @click="toggleListValue(value)"
                                class="h-4 w-4 my-auto rounded border-base-300 text-primary-500 focus:ring-primary-500" />
                            <span class="ml-1">{{ value }}</span>
                        </div>
                    </div>
                    <div v-else-if="type === 'JobList'" class="flex flex-row gap-4 flex-wrap">
                        <div v-for="job in props.jobs" :key="job.getName()" class="flex flex-row flex-initial flex-nowrap">
                            <input :id="job.getName()" :name="job.getName()" type="checkbox"
                                :checked="!!getState().getJobList()?.getStringsList().find(v => v === job.getName())"
                                @click="toggleListValue(job.getName())"
                                class="h-4 w-4 my-auto rounded border-base-300 text-primary-500 focus:ring-primary-500" />
                            <span class="ml-1">{{ job.getLabel() }}</span>
                        </div>
                    </div>
                    <div v-else-if="type === 'JobRankList'" class="flex flex-col gap-2">
                        <div v-for="job in props.jobs" :key="job.getName()"
                            class="flex flex-row flex-initial flex-nowrap gap-2">
                            <span class="flex-1">{{ job.getLabel() }}</span>
                            <span class="flex-1">{{
                                job.getGradesList()[
                                    getState().getJobGradeList()?.getJobsMap().get(job.getName()) ?? -1
                                ]?.getLabel() ?? '-'
                            }}</span>
                            <input id="markerSize" name="markerSize" type="range"
                                class="h-1.5 flex-1 cursor-grab rounded-full my-auto accent-primary-500" min="0"
                                :max="job.getGradesList().length" step="1"
                                :value="(getState().getJobGradeList()?.getJobsMap().get(job.getName()) ?? -1) + 1"
                                @change="updateJobGradeValue(job.getName(), ($event.target as any).value)" />
                        </div>
                    </div>
                    <div v-else>{{ type }} {{ validValues }}</div>
                </div>
            </DisclosurePanel>
        </Disclosure>
    </div>
</template>
