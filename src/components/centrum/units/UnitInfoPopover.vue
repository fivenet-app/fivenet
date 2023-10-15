<script lang="ts" setup>
import { Float } from '@headlessui-float/vue';
import { Popover, PopoverButton, PopoverPanel } from '@headlessui/vue';
import PhoneNumber from '~/components/partials/citizens/PhoneNumber.vue';
import { Unit } from '~~/gen/ts/resources/dispatch/units';

withDefaults(
    defineProps<{
        unit: Unit | undefined;
        initialsOnly?: boolean;
        textClass?: unknown;
        buttonClass?: unknown;
        badge?: boolean;
    }>(),
    {
        initialsOnly: false,
        badge: false,
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
                            class="inline-flex items-center rounded-md bg-gray-400/10 px-2 py-1 text-xs font-medium text-gray-400 ring-1 ring-inset ring-gray-400/20 justify-center"
                        >
                            <template v-if="!initialsOnly"> {{ unit.name }} ({{ unit.initials }}) </template>
                            <template v-else>
                                {{ unit.initials }}
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
                class="absolute z-5 w-64 max-w-[18rem] min-w-fit text-sm text-gray-400 transition-opacity bg-gray-800 border border-gray-600 rounded-lg shadow-sm"
            >
                <div class="p-3">
                    <p class="text-base font-semibold leading-none text-gray-900 dark:text-neutral">
                        {{ unit.name }} ({{ unit.initials }})
                    </p>
                    <p v-for="user in unit.users" class="text-sm font-normal inline-flex items-center justify-center">
                        {{ user.user?.firstname }}
                        {{ user.user?.lastname }}
                        <PhoneNumber :number="user.user?.phoneNumber" :hide-number="true" />
                    </p>
                </div>
            </PopoverPanel>
        </Float>
    </Popover>
</template>
