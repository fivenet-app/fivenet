<script lang="ts" setup>
import { Popover, PopoverButton, PopoverPanel } from '@headlessui/vue';
import { User, UserShort } from '~~/gen/ts/resources/users/users';

defineProps<{
    user: User | UserShort;
}>();
</script>

<template>
    <Popover class="relative">
        <PopoverButton>{{ user.firstname }} {{ user.lastname }}</PopoverButton>

        <PopoverPanel class="absolute z-10">
            <div class="grid grid-cols-2">
                <a href="/analytics">Analytics</a>
                <a href="/engagement">Engagement</a>
                <a href="/security">Security</a>
                <a href="/integrations">Integrations</a>
            </div>
        </PopoverPanel>
    </Popover>
</template>
