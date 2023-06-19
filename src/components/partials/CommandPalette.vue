<script lang="ts" setup>
import SvgIcon from '@jamescoyle/vue-icon';
import { mdiAccount, mdiFileDocument } from '@mdi/js';
import { onClickOutside, useMagicKeys, whenever } from '@vueuse/core';
import { Command } from 'vue-command-palette';
import '~/assets/css/command-palette.scss';

const items = [
    {
        icon: mdiAccount,
        label: 'Go to Citizen: CIT-...',
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
        label: 'Go to Document: DOC-...',
        action: () => {
            const id = inputValue.value.substring(inputValue.value.indexOf('-') + 1);
            if (id.length > 0) {
                navigateTo({ name: 'documents-id', params: { id: id } });
            }

            visible.value = false;
        },
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

                <Command.Item v-for="item in items" :data-value="item.label" @select="item.action">
                    <SvgIcon type="mdi" :path="item.icon" class="w-6 h-6" />
                    <div>{{ item.label }}</div>
                </Command.Item>
            </Command.List>
        </template>
    </Command.Dialog>
</template>
