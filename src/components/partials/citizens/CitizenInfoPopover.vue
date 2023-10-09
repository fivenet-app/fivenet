<script lang="ts" setup>
import { Float } from '@headlessui-float/vue';
import { Popover, PopoverButton, PopoverPanel } from '@headlessui/vue';
import { AccountIcon } from 'mdi-vue3';
import { User, UserShort } from '~~/gen/ts/resources/users/users';
import PhoneNumber from './PhoneNumber.vue';

defineProps<{
    user: User | UserShort | undefined;
    noPopover?: boolean;
    textClass?: unknown;
    buttonClass?: unknown;
}>();
</script>

<template>
    <template v-if="!user">
        <span class="inline-flex items-center">
            <slot name="before" />
            <span>N/A</span>
            <slot name="after" />
        </span>
    </template>
    <template v-else-if="noPopover">
        <span class="inline-flex items-center">
            <slot name="before" />
            <NuxtLink :to="{ name: 'citizens-id', params: { id: user.userId } }">
                {{ user.firstname }} {{ user.lastname }}
            </NuxtLink>
            <span v-if="user.phoneNumber">
                <PhoneNumber v-if="user.phoneNumber" :number="user.phoneNumber" :hide-number="true" :show-label="false" />
            </span>
            <slot name="after" />
        </span>
    </template>
    <Popover v-else class="relative">
        <Float portal auto-placement :offset="16">
            <PopoverButton class="inline-flex items-center" :class="buttonClass">
                <slot name="before" />
                <span :class="textClass"> {{ user.firstname }} {{ user.lastname }} </span>
                <slot name="after" />
            </PopoverButton>

            <PopoverPanel
                class="absolute z-5 w-64 max-w-[18rem] min-w-fit text-sm text-gray-400 transition-opacity bg-gray-800 border border-gray-600 rounded-lg shadow-sm"
            >
                <div class="p-3">
                    <div class="flex items-center gap-2 mb-2">
                        <NuxtLink
                            :to="{ name: 'citizens-id', params: { id: user.userId } }"
                            class="inline-flex items-center text-primary-500 hover:text-primary-400"
                        >
                            <AccountIcon class="w-6 h-6" />
                            <span class="ml-1">{{ $t('common.profile') }}</span>
                        </NuxtLink>
                        <PhoneNumber
                            v-if="user.phoneNumber"
                            :number="user.phoneNumber"
                            :hide-number="true"
                            :show-label="true"
                        />
                    </div>
                    <p class="text-base font-semibold leading-none text-gray-900 dark:text-white">
                        <NuxtLink :to="{ name: 'citizens-id', params: { id: user.userId } }">
                            {{ user.firstname }} {{ user.lastname }}
                        </NuxtLink>
                    </p>
                    <p v-if="user.jobLabel" class="text-sm font-normal">
                        {{ user.jobLabel }}
                        <span v-if="user.jobGradeLabel"> ({{ $t('common.rank') }}: {{ user.jobGradeLabel }})</span>
                    </p>
                </div>
            </PopoverPanel>
        </Float>
    </Popover>
</template>
