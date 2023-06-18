<script lang="ts" setup>
import {
    Disclosure,
    DisclosureButton,
    DisclosurePanel,
    Listbox,
    ListboxButton,
    ListboxOption,
    ListboxOptions,
} from '@headlessui/vue';
import SvgIcon from '@jamescoyle/vue-icon';
import { mdiCheck, mdiChevronDown } from '@mdi/js';
import { Job, JobGrade } from '~~/gen/ts/resources/jobs/jobs';
import { AttributeValues, Permission, RoleAttribute } from '~~/gen/ts/resources/permissions/permissions';

const props = defineProps<{
    attribute: RoleAttribute;
    states: Map<bigint, AttributeValues | undefined>;
    jobs: Job[];
    disabled?: boolean;
    permission: Permission;
}>();

const emit = defineEmits<{
    (e: 'update:states', payload: Map<bigint, AttributeValues | undefined>): void;
}>();

const states = ref<typeof props.states>(props.states);
const id = ref<bigint>(props.attribute.attrId);
const maxValues = props.attribute.maxValues;

const jobGrades = ref<Map<string, JobGrade>>(new Map());

const validValues = ref<AttributeValues | undefined>(props.attribute.validValues);
if (!states.value.has(id.value)) {
    switch (lowercaseFirstLetter(props.attribute.type)) {
        case 'stringList':
            states.value.set(id.value, {
                validValues: {
                    oneofKind: 'stringList',
                    stringList: {
                        strings: [],
                    },
                },
            });
            break;

        case 'jobList':
            states.value.set(id.value, {
                validValues: {
                    oneofKind: 'jobList',
                    jobList: {
                        strings: [],
                    },
                },
            });
            break;

        case 'jobGradeList':
            states.value.set(id.value, {
                validValues: {
                    oneofKind: 'jobGradeList',
                    jobGradeList: {
                        jobs: {},
                    },
                },
            });
            break;
    }
}

const state: AttributeValues = states.value.get(id.value)!;

async function toggleStringListValue(value: string): Promise<void> {
    if (state.validValues.oneofKind !== 'stringList') {
        return;
    }

    const array = state.validValues.stringList.strings;
    if (array.indexOf(value) < 0) {
        array.push(value);
    } else {
        array.splice(array.indexOf(value), 1);
    }

    state.validValues.stringList.strings = array;
    states.value.set(id.value, state);
    emit('update:states', states.value);
}

async function toggleJobListValue(value: string): Promise<void> {
    if (state.validValues.oneofKind !== 'jobList') {
        return;
    }

    const array = state.validValues.jobList.strings;
    if (array.indexOf(value) < 0) {
        array.push(value);
    } else {
        array.splice(array.indexOf(value), 1);
    }

    state.validValues.jobList.strings = array;
    states.value.set(id.value, state);
    emit('update:states', states.value);
}

async function toggleJobGradeValue(job: Job, checked: boolean): Promise<void> {
    if (state.validValues.oneofKind !== 'jobGradeList') {
        return;
    }

    const map = state.validValues.jobGradeList.jobs;
    if (checked && !map[job.name]) {
        map[job.name] = 1;
        jobGrades.value.set(job.name, job.grades[0]);
    } else if (!checked && map[job.name]) {
        delete map[job.name];
        jobGrades.value.set(job.name, job.grades[0]);
    } else {
        return;
    }

    state.validValues.jobGradeList.jobs = map;
    states.value.set(id.value, state);
    emit('update:states', states.value);
}

async function updateJobGradeValue(job: Job, grade: JobGrade): Promise<void> {
    if (state.validValues.oneofKind !== 'jobGradeList') {
        return;
    }

    const map = state.validValues.jobGradeList.jobs;

    map[job.name] = grade.grade;
    jobGrades.value.set(job.name, job.grades[grade.grade - 1]);

    state.validValues.jobGradeList.jobs = map;
    states.value.set(id.value, state);
    emit('update:states', states.value);
}

onMounted(() => {
    if (state.validValues.oneofKind === 'jobGradeList') {
        props.jobs.forEach((job) => {
            if (state.validValues.oneofKind !== 'jobGradeList') {
                return;
            }
            if (maxValues && maxValues.validValues.oneofKind === 'jobGradeList') {
                if (!maxValues.validValues.jobGradeList.jobs[job.name]) {
                    return;
                }
            }
            jobGrades.value.set(job.name, job.grades[(state.validValues?.jobGradeList.jobs[job.name] ?? 1) - 1]);
        });
    }
});
</script>

<template>
    <div v-if="$props.attribute">
        <Disclosure
            as="div"
            :class="[
                $props.disabled ? 'border-neutral/10 text-base-300' : 'hover:border-neutral/70 border-neutral/20 text-neutral',
            ]"
            v-slot="{ open }"
        >
            <DisclosureButton
                :disabled="$props.disabled"
                :class="[
                    open ? 'rounded-t-lg border-b-0' : 'rounded-lg',
                    $props.disabled ? 'cursor-not-allowed' : '',
                    ' flex w-full items-start justify-between text-left border-2 p-2 border-inherit transition-colors',
                ]"
            >
                <span class="text-base leading-7 transition-colors">
                    {{ $t(`attrs.${attribute.category}.${attribute.name}.${attribute.key}`) }}
                </span>
                <span class="ml-6 flex h-7 items-center">
                    <SvgIcon
                        :class="[open ? 'upsidedown' : '', 'h-6 w-6 transition-transform']"
                        aria-hidden="true"
                        type="mdi"
                        :path="mdiChevronDown"
                    />
                </span>
            </DisclosureButton>
            <DisclosurePanel class="px-4 pb-2 border-2 border-t-0 rounded-b-lg transition-colors border-inherit -mt-2">
                <div class="flex flex-col gap-2 max-w-4xl mx-auto my-2">
                    <div
                        v-if="
                            state.validValues.oneofKind === 'stringList' &&
                            maxValues?.validValues &&
                            maxValues?.validValues.oneofKind === 'stringList'
                        "
                        class="flex flex-row gap-4 flex-wrap"
                    >
                        <div
                            v-for="value in maxValues.validValues.stringList.strings"
                            :key="value"
                            class="flex flex-row flex-initial flex-nowrap"
                        >
                            <input
                                :id="value"
                                :name="value"
                                type="checkbox"
                                :checked="!!state.validValues.stringList.strings.find((v) => v === value)"
                                @click="toggleStringListValue(value)"
                                class="h-4 w-4 my-auto rounded border-base-300 text-primary-500 focus:ring-primary-500"
                            />
                            <span class="ml-1">{{
                                $t(`perms.${permission.category}.${permission.name}.attrs.${value.replaceAll('.', '_')}`)
                            }}</span>
                        </div>
                    </div>
                    <div
                        v-else-if="
                            state.validValues.oneofKind === 'jobList' &&
                            maxValues?.validValues &&
                            maxValues?.validValues.oneofKind === 'jobList'
                        "
                        class="flex flex-row gap-4 flex-wrap"
                    >
                        <div
                            v-for="job in props.jobs.filter(
                                (j) =>
                                    maxValues?.validValues.oneofKind === 'jobList' &&
                                    (!maxValues?.validValues.jobList?.strings.length ||
                                        maxValues.validValues?.jobList?.strings.includes(j.name))
                            )"
                            :key="job.name"
                            class="flex flex-row flex-initial flex-nowrap"
                        >
                            <input
                                :id="job.name"
                                :name="job.name"
                                type="checkbox"
                                :checked="!!state.validValues.jobList?.strings.find((v) => v === job.name)"
                                @click="toggleJobListValue(job.name)"
                                class="h-4 w-4 my-auto rounded border-base-300 text-primary-500 focus:ring-primary-500"
                            />
                            <span class="ml-1">{{ job.label }}</span>
                        </div>
                    </div>
                    <div
                        v-else-if="
                            state.validValues.oneofKind === 'jobGradeList' &&
                            maxValues?.validValues &&
                            maxValues.validValues.oneofKind === 'jobGradeList'
                        "
                        class="flex flex-col gap-2"
                    >
                        <div
                            v-for="job in props.jobs.filter(
                                (j) =>
                                    maxValues &&
                                    maxValues.validValues.oneofKind === 'jobGradeList' &&
                                    maxValues.validValues.jobGradeList.jobs[j.name]
                            )"
                            :key="job.name"
                            class="flex flex-row flex-initial flex-nowrap gap-2"
                        >
                            <input
                                :id="job.name"
                                :name="job.name"
                                type="checkbox"
                                :checked="!!state.validValues?.jobGradeList.jobs[job.name]"
                                @change="toggleJobGradeValue(job, ($event.target as any).checked)"
                                class="h-4 w-4 my-auto rounded border-base-300 text-primary-500 focus:ring-primary-500"
                            />
                            <span class="flex-1 my-auto">{{ job.label }}</span>
                            <Listbox
                                as="div"
                                class="flex-1"
                                :model-value="jobGrades.get(job.name)"
                                @update:model-value="updateJobGradeValue(job, $event)"
                                :disabled="!state.validValues.jobGradeList?.jobs[job.name]"
                            >
                                <div class="relative">
                                    <ListboxButton
                                        class="block pl-3 text-left w-full rounded-md border-0 py-1.5 bg-base-700 text-neutral placeholder:text-base-200 disabled:bg-base-800 disabled:text-neutral/50 disabled:cursor-not-allowed focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6"
                                    >
                                        <span class="block truncate">{{ jobGrades.get(job.name)?.label }}</span>
                                        <span class="absolute inset-y-0 right-0 flex items-center pr-2 pointer-events-none">
                                            <SvgIcon
                                                class="w-5 h-5 text-gray-400"
                                                aria-hidden="true"
                                                type="mdi"
                                                :path="mdiChevronDown"
                                            />
                                        </span>
                                    </ListboxButton>

                                    <transition
                                        leave-active-class="transition duration-100 ease-in"
                                        leave-from-class="opacity-100"
                                        leave-to-class="opacity-0"
                                    >
                                        <ListboxOptions
                                            class="absolute z-10 w-full py-1 mt-1 overflow-auto text-base rounded-md bg-base-700 max-h-60 sm:text-sm"
                                        >
                                            <ListboxOption
                                                as="template"
                                                v-for="grade in job.grades.filter(
                                                    (g) =>
                                                        maxValues &&
                                                        maxValues.validValues.oneofKind === 'jobGradeList' &&
                                                        maxValues.validValues.jobGradeList.jobs[job.name] + 1 > g.grade
                                                )"
                                                :key="grade.grade"
                                                :value="grade"
                                                v-slot="{ active, selected }"
                                            >
                                                <li
                                                    :class="[
                                                        active ? 'bg-primary-500' : '',
                                                        'text-neutral relative cursor-default select-none py-2 pl-8 pr-4',
                                                    ]"
                                                >
                                                    <span
                                                        :class="[selected ? 'font-semibold' : 'font-normal', 'block truncate']"
                                                    >
                                                        {{ grade.label }}
                                                    </span>

                                                    <span
                                                        v-if="selected"
                                                        :class="[
                                                            active ? 'text-neutral' : 'text-primary-500',
                                                            'absolute inset-y-0 left-0 flex items-center pl-1.5',
                                                        ]"
                                                    >
                                                        <SvgIcon
                                                            class="w-5 h-5"
                                                            aria-hidden="true"
                                                            type="mdi"
                                                            :path="mdiCheck"
                                                        />
                                                    </span>
                                                </li>
                                            </ListboxOption>
                                        </ListboxOptions>
                                    </transition>
                                </div>
                            </Listbox>
                        </div>
                    </div>
                    <div v-else>{{ state.validValues.oneofKind }} {{ validValues }}</div>
                </div>
            </DisclosurePanel>
        </Disclosure>
    </div>
</template>
