<script lang="ts" setup>
import SvgIcon from '@jamescoyle/vue-icon';
import { mdiAccount, mdiFileDocument } from '@mdi/js';
import { onClickOutside, useMagicKeys, whenever } from '@vueuse/core';
import { Command } from 'vue-command-palette';
import '~/assets/css/command-palette.scss';

const groups = [
    {
        name: 'shortcuts',
        items: [
            {
                icon: mdiAccount,
                name: 'goto-citizen',
                action: () => {
                    const id = inputValue.value.substring(inputValue.value.indexOf('-') + 1);
                    if (id.length > 0) {
                        navigateTo({ name: 'citizens-id', params: { id: id } });
                    }

                    visible.value = false;
                },
            },
            {
                icon: mdiFileDocument,
                name: 'goto-doc',
                action: () => {
                    const id = inputValue.value.substring(inputValue.value.indexOf('-') + 1);
                    if (id.length > 0) {
                        navigateTo({ name: 'documents-id', params: { id: id } });
                    }

                    visible.value = false;
                },
            },
        ],
    },
];

const visible = ref(false);
const target = ref(null);

const keys = useMagicKeys({
    passive: false,
    onEventFired(e) {
        if ((e.metaKey || e.ctrlKey) && e.key === 'k' && e.type === 'keydown') {
            e.preventDefault();
            visible.value = true;
        }
    },
});
const Escape = keys['Escape'];

const inputValue = ref('');

whenever(Escape, () => {
    visible.value = false;
});

onClickOutside(target, () => {
    visible.value = false;
});
</script>

<template>
    <Command.Dialog ref="target" :visible="visible" theme="custom">
        <template #header>
            <Command.Input :placeholder="$t('commandpalette.input')" v-model:value="inputValue" />
        </template>
        <template #body>
            <Command.List>
                <Command.Empty>{{ $t('commandpalette.empty') }}</Command.Empty>

                <Command.Group v-for="group in groups" :heading="$t(`commandpalette.groups.${group.name}.heading`)">
                    <Command.Item v-for="item in group.items" :data-value="item.name" @select="item.action">
                        <SvgIcon type="mdi" :path="item.icon" class="w-6 h-6" />
                        <div>{{ $t(`commandpalette.groups.${group.name}.${item.name}`) }}</div>
                    </Command.Item>
                </Command.Group>
            </Command.List>
        </template>
    </Command.Dialog>
</template>
