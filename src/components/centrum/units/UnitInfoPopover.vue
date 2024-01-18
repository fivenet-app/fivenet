<script lang="ts" setup>
import { Float } from '@headlessui-float/vue';
import { Popover, PopoverButton, PopoverPanel } from '@headlessui/vue';
import { TimerIcon } from 'mdi-vue3';
import PhoneNumberBlock from '~/components/partials/citizens/PhoneNumberBlock.vue';
import { useCentrumStore } from '~/store/centrum';
import { DispatchAssignment } from '~~/gen/ts/resources/centrum/dispatches';
import { Unit } from '~~/gen/ts/resources/centrum/units';

const centrumStore = useCentrumStore();
const { timeCorrection } = storeToRefs(centrumStore);

withDefaults(
    defineProps<{
        unit: Unit | undefined;
        initialsOnly?: boolean;
        textClass?: unknown;
        buttonClass?: unknown;
        badge?: boolean;
        assignment?: DispatchAssignment;
    }>(),
    {
        initialsOnly: false,
        badge: false,
        assignment: undefined,
        textClass: '' as any,
        buttonClass: '' as any,
    },
);
</script>

<template>
    <template v-if="!unit">
        <span class="inline-flex items-center">
            <slot name="before" />
            <span>N/A</span>
            <slot name="after" />
        </span>
    </template>
    <Popover v-else class="relative">
        <Float portal auto-placement :offset="16">
            <PopoverButton class="inline-flex items-center" :class="buttonClass">
                <slot name="before" />
                <span :class="textClass">
                    <template v-if="badge">
                        <span
                            class="inline-flex items-center justify-center rounded-md bg-gray-400/10 px-2 py-1 text-xs font-medium text-gray-400 ring-1 ring-inset ring-gray-400/20"
                        >
                            <template v-if="!initialsOnly"> {{ unit.name }} ({{ unit.initials }}) </template>
                            <template v-else>
                                {{ unit.initials }}
                            </template>
                            <template v-if="assignment?.expiresAt">
                                <TimerIcon class="ml-0.5 h-5 w-5 fill-warn-600" />
                            </template>
                        </span>
                    </template>
                    <template v-else>
                        <template v-if="!initialsOnly"> {{ unit.name }} ({{ unit.initials }}) </template>
                        <template v-else>
                            {{ unit.initials }}
                        </template>
                    </template>
                </span>
                <slot name="after" />
            </PopoverButton>

            <PopoverPanel
                class="absolute z-5 w-64 min-w-fit max-w-[18rem] rounded-lg border border-gray-600 bg-gray-800 text-sm text-gray-400 shadow-sm transition-opacity"
            >
                <div class="p-3">
                    <p class="text-base font-semibold leading-none text-gray-900 dark:text-neutral">
                        {{ unit.name }} ({{ unit.initials }})
                    </p>
                    <p v-if="assignment?.expiresAt" class="inline-flex items-center justify-center text-sm font-normal">
                        {{
                            useLocaleTimeAgo(toDate(assignment.expiresAt, timeCorrection), {
                                showSecond: true,
                                updateInterval: 1_000,
                            }).value
                        }}
                    </p>
                    <p class="text-base font-medium leading-none text-gray-800 dark:text-neutral">
                        {{ $t('common.members') }}
                    </p>
                    <ul class="text-sm font-normal">
                        <li v-for="user in unit.users" :key="user.userId" class="inline-flex items-center">
                            {{ user.user?.firstname }}
                            {{ user.user?.lastname }}
                            <PhoneNumberBlock class="ml-0.5" :number="user.user?.phoneNumber" :hide-number="true" />
                        </li>
                    </ul>
                </div>
            </PopoverPanel>
        </Float>
    </Popover>
</template>
