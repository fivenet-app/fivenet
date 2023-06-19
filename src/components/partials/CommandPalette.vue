<script lang="ts" setup>
import SvgIcon from '@jamescoyle/vue-icon';
import {
    mdiAccount,
    mdiAccountPlus,
    mdiCalendar,
    mdiLabelMultiple,
    mdiLabelOff,
    mdiListStatus,
    mdiPriorityHigh,
} from '@mdi/js';
import { onClickOutside, useMagicKeys, whenever } from '@vueuse/core';
import { Command } from 'vue-command-palette';

const visible = ref(false);
const target = ref(null);

const keys = useMagicKeys();
const ControlK = keys['Ctrl+K'];
const Escape = keys['Escape'];

whenever(ControlK, () => {
    visible.value = true;
});

whenever(Escape, () => {
    visible.value = false;
});

onClickOutside(target, () => {
    visible.value = false;
});

const handleSelectItem = (item: any) => {
    console.log(item);
};

const groups = [
    {
        label: 'Test',
        items: [
            {
                icon: mdiAccountPlus,
                label: 'Assign to...',
                shortcut: ['A'],
                perform: () => {
                    console.log('action');
                },
            },
            {
                icon: mdiAccount,
                label: 'Assign to me',
                shortcut: ['I'],
            },
            {
                icon: mdiListStatus,
                label: 'Change status...',
                shortcut: ['S'],
            },
            {
                icon: mdiPriorityHigh,
                label: 'Change priority...',
                shortcut: ['P'],
            },
            {
                icon: mdiLabelMultiple,
                label: 'Change labels...',
                shortcut: ['L'],
            },
            {
                icon: mdiLabelOff,
                label: 'Remove label...',
                shortcut: ['⇧', 'L'],
            },
            {
                icon: mdiCalendar,
                label: 'Set due date...',
                shortcut: ['⇧', 'D'],
            },
        ],
    },
];
</script>

<template>
    <Command.Dialog ref="target" :visible="visible" theme="custom">
        <template #header>
            <Command.Input placeholder="Type a command or search..." />
        </template>
        <template #body>
            <Command.List>
                <Command.Empty>No results found.</Command.Empty>

                <Command.Group v-for="g in groups" :heading="g.label">
                    <Command.Item
                        v-for="item in g.items"
                        :data-value="item.label"
                        :shortcut="item.shortcut"
                        :perform="item.perform"
                        @select="handleSelectItem"
                    >
                        <SvgIcon type="mdi" :path="item.icon" class="w-6 h-6" />
                        <div>{{ item.label }}</div>
                        <div command-linear-shortcuts>
                            <kbd v-for="key in item.shortcut" key="key">{{ key }}</kbd>
                        </div>
                    </Command.Item>
                </Command.Group>
            </Command.List>
        </template>
    </Command.Dialog>
</template>

<style>
@import '~/assets/css/command-palette.scss';
</style>
