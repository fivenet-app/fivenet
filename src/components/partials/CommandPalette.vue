<script lang="ts" setup>
import SvgIcon from '@jamescoyle/vue-icon';
import { mdiAccount, mdiFileDocument } from '@mdi/js';
import { onClickOutside, useMagicKeys, whenever } from '@vueuse/core';
import { Command } from 'vue-command-palette';
import '~/assets/css/command-palette.scss';

const { t } = useI18n();

const inputValue = ref('');

const groups = [
    {
        name: 'shortcuts',
        items: [
            {
                icon: mdiAccount,
                name: t(`commandpalette.groups.shortcuts.goto`, [t('common.citizen', 1), `CIT-`]),
                action: () => {
                    const id = inputValue.value.substring(inputValue.value.indexOf('-') + 1);
                    if (id.length > 0) {
                        navigateTo({ name: 'citizens-id', params: { id: id } });
                        visible.value = false;
                    } else {
                        inputValue.value = 'CIT-';
                    }
                },
            },
            {
                icon: mdiFileDocument,
                name: t(`commandpalette.groups.shortcuts.goto`, [t('common.document', 1), `DOC-`]),
                action: () => {
                    const id = inputValue.value.substring(inputValue.value.indexOf('-') + 1);
                    if (id.length > 0) {
                        navigateTo({ name: 'documents-id', params: { id: id } });
                        visible.value = false;
                    } else {
                        inputValue.value = 'DOC-';
                    }
                },
            },
        ],
    },
];
const dynamicItems = ref<
    {
        name: string;
        icon: string;
        group: string;
        action: () => any;
    }[]
>([]);

const visible = ref(false);
const target = ref(null);

const keys = useMagicKeys({
    passive: false,
    onEventFired(e) {
        if (
            ((e.metaKey || e.ctrlKey) && e.key === 'k' && e.type === 'keydown') ||
            (e.shiftKey && e.key === '/' && e.type === 'keydown')
        ) {
            e.preventDefault();
            visible.value = true;
        }
    },
});
const Escape = keys['Escape'];

whenever(Escape, () => {
    visible.value = false;
});

onClickOutside(target, () => {
    visible.value = false;
});

watch(inputValue, () => {
    const val = inputValue.value.toLowerCase();
    const item = {
        name: inputValue.value,
        icon: mdiAccount,
        group: 'shortcuts',
        action: () => {},
    };
    if (val.startsWith('cit-')) {
        item.icon = mdiAccount;
    } else if (val.startsWith('doc-')) {
        item.icon = mdiFileDocument;
    } else {
        dynamicItems.value.length = 0;
        return;
    }

    dynamicItems.value.pop();
    dynamicItems.value.push(item);
    console.log('DYNAMIC ITEMS', val, dynamicItems.value);
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

                <Command.Group v-if="dynamicItems && dynamicItems.length > 0" heading="Results">
                    <Command.Item v-for="item in dynamicItems" :data-value="item.name" @select="item.action">
                        <SvgIcon type="mdi" :path="item.icon" class="w-6 h-6" />
                        <div>{{ item.name }}</div>
                    </Command.Item>
                </Command.Group>

                <Command.Group v-for="group in groups" :heading="$t(`commandpalette.groups.${group.name}.heading`)">
                    <Command.Item v-for="item in group.items" :data-value="item.name" @select="item.action">
                        <SvgIcon type="mdi" :path="item.icon" class="w-6 h-6" />
                        <div>{{ item.name }}</div>
                    </Command.Item>
                </Command.Group>
            </Command.List>
        </template>
    </Command.Dialog>
</template>
