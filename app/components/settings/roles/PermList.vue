<script lang="ts" setup>
import RoleViewAttr from '~/components/settings/roles/RoleViewAttr.vue';
import type { RoleAttribute } from '~~/gen/ts/resources/permissions/attributes';
import type { Permission } from '~~/gen/ts/resources/permissions/permissions';
import { isEmptyAttributes } from './helpers';

const props = defineProps<{
    permissions: Permission[];
    attributes: RoleAttribute[];
    disabled?: boolean;
}>();

const { t } = useI18n();

const permCategories = ref<Set<string>>(new Set());
const permStates = ref(new Map<number, boolean | undefined>());

const attrList = ref<RoleAttribute[]>([]);

const accordionCategories = computed(() =>
    [...permCategories.value.entries()].map((category) => {
        return {
            category: category[0],
            label: t(`perms.${category[1]}.category`),
            disabled: props.disabled,
        };
    }),
);

async function genPermissionCategories(): Promise<void> {
    permCategories.value.clear();

    props.permissions.forEach((perm) => {
        permCategories.value.add(perm.category);
    });
}

async function propogateRolePermissionStates(): Promise<void> {
    permStates.value.clear();

    props.permissions.forEach((perm) => permStates.value.set(perm.id, true));

    props.attributes.forEach((attr) => {
        const idx = attrList.value.findIndex((a) => a.attrId === attr.attrId);
        if (idx > -1 && attrList.value[idx]) {
            attrList.value[idx].value = attr.value;
        } else {
            attrList.value.push(attr);
        }
    });
}

async function setFromProps(): Promise<void> {
    await genPermissionCategories();
    await propogateRolePermissionStates();
}

setFromProps();
watch(props, setFromProps);
</script>

<template>
    <div class="w-full">
        <div class="px-1 sm:px-2">
            <div class="flex flex-col gap-2">
                <UAccordion :items="accordionCategories" multiple default-open :unmount="true">
                    <template #item="{ item: category }">
                        <div class="flex flex-col divide-y divide-gray-100 dark:divide-gray-800">
                            <div
                                v-for="perm in permissions.filter((p) => p.category === category.category)"
                                :key="perm.id"
                                class="flex flex-col gap-1"
                            >
                                <div class="flex flex-row items-center gap-2">
                                    <div class="flex-1">
                                        <p class="text-highlighted" :title="`${$t('common.id')}: ${perm.id}`">
                                            {{ $t(`perms.${perm.category}.${perm.name}.key`) }}
                                        </p>
                                        <p class="text-base-500">
                                            {{ $t(`perms.${perm.category}.${perm.name}.description`) }}
                                        </p>
                                    </div>

                                    <UButtonGroup class="inline-flex flex-initial">
                                        <UButton color="green" variant="solid" icon="i-mdi-check" disabled />
                                    </UButtonGroup>
                                </div>

                                <template v-for="(attr, idx) in attrList" :key="attr.attrId">
                                    <RoleViewAttr
                                        v-if="attr.permissionId === perm.id && !isEmptyAttributes(attr.maxValues)"
                                        v-model="attrList[idx]!"
                                        :permission="perm"
                                        disabled
                                        default-open
                                    />
                                </template>
                            </div>
                        </div>
                    </template>
                </UAccordion>
            </div>
        </div>
    </div>
</template>
