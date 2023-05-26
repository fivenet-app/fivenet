<script lang="ts" setup>
import { Disclosure, DisclosureButton, DisclosurePanel, Listbox, ListboxButton, ListboxOption, ListboxOptions } from '@headlessui/vue'
import { AttributeValues, JobGradeList, RoleAttribute, StringList } from '~~/gen/ts/resources/permissions/permissions';
import { ChevronDownIcon, CheckIcon } from '@heroicons/vue/24/solid';
import { Job, JobGrade } from '~~/gen/ts/resources/jobs/jobs';

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
const id = ref<number>(props.attribute.attrId);

const validValues = ref<AttributeValues | undefined>(props.attribute.validValues);

const type = ref<'stringList' | 'jobList' | 'jobGradeList' | undefined>();
switch (props.attribute.type) {
    case 'stringList':
        type.value = 'stringList';
        break;

    case 'jobList':
        type.value = 'jobList';
        break;

    case 'jobGradeList':
        type.value = 'jobGradeList';
        break;

    default:
        break;
}

const jobGrades = ref<Map<string, JobGrade>>(new Map());

function getState(): AttributeValues {
    if (!states.value.has(id.value)) states.value.set(id.value, {
        validValues: {
            oneofKind: type,
        },
    });

    return states.value.get(id.value)!;
}

async function toggleListValue(value: string): Promise<void> {
    const state = getState();
    const list = state.validValues ?? new StringList();
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

async function toggleJobListValue(value: string): Promise<void> {
    const state = getState();
    const list = state.getJobList() ?? new StringList();
    const array = list.getStringsList();

    if (array.indexOf(value) < 0) {
        array.push(value);
    } else {
        array.splice(array.indexOf(value), 1);
    }

    list.setStringsList(array);
    state.setJobList(list);
    states.value.set(id.value, state);
    emit('update:states', states.value);
}

async function toggleJobGradeValue(job: Job, checked: boolean): Promise<void> {
    const state = getState();
    const list = state.getJobGradeList() ?? new JobGradeList();
    const map = list.getJobsMap();

    if (checked && !map.has(job.name)) {
        map.set(job.name, 1);
        jobGrades.value.set(job.name, job.getGradesList()[0]);
    } else if (!checked && map.has(job.name)) {
        map.del(job.name);
        jobGrades.value.set(job.name, job.getGradesList()[0]);
    } else {
        return;
    }

    state.setJobGradeList(list);
    states.value.set(id.value, state);
    emit('update:states', states.value);
}

async function updateJobGradeValue(job: Job, grade: JobGrade): Promise<void> {
    const state = getState();
    const list = state.getJobGradeList() ?? new JobGradeList();
    const map = list.getJobsMap();

    map.set(job.name, grade.getGrade());
    jobGrades.value.set(job.name, job.grades[grade.grade - 1]);

    state.setJobGradeList(list);
    states.value.set(id.value, state);
    emit('update:states', states.value);
}

onMounted(() => {
    if (type.value === 'JobGradeList') {
        props.jobs.forEach(job => {
            jobGrades.value.set(job.name, job.grades[(getState().getJobGradeList()?.getJobsMap().get(job.name) ?? 1) - 1]);
        });
    }
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
                    {{ $t(`attrs.${attribute.category}.${attribute.name}.${attribute.key}`) }}
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
                        <div v-for="job in props.jobs.filter(j => !validValues?.getJobList()?.getStringsList().length || validValues.getJobList()?.getStringsList().includes(j.name))"
                            :key="job.name" class="flex flex-row flex-initial flex-nowrap">
                            <input :id="job.name" :name="job.name" type="checkbox"
                                :checked="!!getState().getJobList()?.getStringsList().find(v => v === job.name)"
                                @click="toggleJobListValue(job.name)"
                                class="h-4 w-4 my-auto rounded border-base-300 text-primary-500 focus:ring-primary-500" />
                            <span class="ml-1">{{ job.label }}</span>
                        </div>
                    </div>
                    <div v-else-if="type === 'JobGradeList'" class="flex flex-col gap-2">
                        <div v-for="job in props.jobs" :key="job.name" class="flex flex-row flex-initial flex-nowrap gap-2">
                            <input :id="job.name" :name="job.name" type="checkbox"
                                :checked="!!getState().getJobGradeList()?.getJobsMap().has(job.name)"
                                @change="toggleJobGradeValue(job, ($event.target as any).checked)"
                                class="h-4 w-4 my-auto rounded border-base-300 text-primary-500 focus:ring-primary-500" />
                            <span class="flex-1 my-auto">{{ job.label }}</span>
                            <Listbox as="div" class="flex-1" :model-value="jobGrades.get(job.name)"
                                @update:model-value="updateJobGradeValue(job, $event)"
                                :disabled="!getState().getJobGradeList()?.getJobsMap().has(job.name)">
                                <div class="relative">
                                    <ListboxButton
                                        class="block pl-3 text-left w-full rounded-md border-0 py-1.5 bg-base-700 text-neutral placeholder:text-base-200 disabled:bg-base-800 disabled:text-neutral/50 disabled:cursor-not-allowed focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6">
                                        <span class="block truncate">{{ jobGrades.get(job.name)?.label }}</span>
                                        <span class="absolute inset-y-0 right-0 flex items-center pr-2 pointer-events-none">
                                            <ChevronDownIcon class="w-5 h-5 text-gray-400" aria-hidden="true" />
                                        </span>
                                    </ListboxButton>

                                    <transition leave-active-class="transition duration-100 ease-in"
                                        leave-from-class="opacity-100" leave-to-class="opacity-0">
                                        <ListboxOptions
                                            class="absolute z-10 w-full py-1 mt-1 overflow-auto text-base rounded-md bg-base-700 max-h-60 sm:text-sm">
                                            <ListboxOption as="template" v-for="grade in job.grades" :key="grade.grade"
                                                :value="grade" v-slot="{ active, selected }">
                                                <li
                                                    :class="[active ? 'bg-primary-500' : '', 'text-neutral relative cursor-default select-none py-2 pl-8 pr-4']">
                                                    <span
                                                        :class="[selected ? 'font-semibold' : 'font-normal', 'block truncate']">{{
                                                            grade.label
                                                        }}</span>

                                                    <span v-if="selected"
                                                        :class="[active ? 'text-neutral' : 'text-primary-500', 'absolute inset-y-0 left-0 flex items-center pl-1.5']">
                                                        <CheckIcon class="w-5 h-5" aria-hidden="true" />
                                                    </span>
                                                </li>
                                            </ListboxOption>
                                        </ListboxOptions>
                                    </transition>
                                </div>
                            </Listbox>
                        </div>
                    </div>
                    <div v-else>{{ type }} {{ validValues }}</div>
                </div>
            </DisclosurePanel>
        </Disclosure>
    </div>
</template>
