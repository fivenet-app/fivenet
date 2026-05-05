<script lang="ts" setup>
import RoleViewAttr from '~/components/settings/roles/RoleViewAttr.vue';
import {
    buildPermissionGroups,
    getPermissionNamespaceLabel,
    getPermissionServiceLabel,
    type PermissionNamespaceGroup,
    type PermissionServiceGroup,
} from '~/components/settings/permissions';
import type { RoleAttribute } from '~~/gen/ts/resources/permissions/attributes/attributes';
import type { Permission } from '~~/gen/ts/resources/permissions/permissions/permissions';
import { isEmptyAttributes } from './helpers';

const props = defineProps<{
    permissions: Permission[];
    attributes: RoleAttribute[];
    disabled?: boolean;
}>();

const { t, te } = useI18n();

const permCategories = ref<PermissionNamespaceGroup[]>([]);
const permStates = ref(new Map<number, boolean | undefined>());

const attrList = ref<RoleAttribute[]>([]);

const accordionCategories = computed(() =>
    permCategories.value.map((namespace) => {
        const services = namespace.services.map((service) => ({
            ...service,
            label: getPermissionServiceLabel(service.namespace, service.service, t, te),
        }));
        const singleService = services.length === 1 ? services[0] : undefined;

        return {
            ...namespace,
            services,
            singleService,
            label: singleService?.label ?? getPermissionNamespaceLabel(namespace.namespace, t, te),
        };
    }),
);

async function genPermissionCategories(): Promise<void> {
    permCategories.value = buildPermissionGroups(props.permissions);
}

function getPermissionsForService(service: PermissionServiceGroup): Permission[] {
    return props.permissions.filter((perm) => perm.namespace === service.namespace && perm.service === service.service);
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
                <UAccordion :items="accordionCategories" type="multiple" default-open>
                    <template #content="{ item: namespace }">
                        <div v-if="namespace.singleService" class="flex flex-col divide-y divide-default">
                            <div
                                v-for="perm in getPermissionsForService(namespace.singleService)"
                                :key="perm.id"
                                class="flex flex-col gap-1"
                            >
                                <div class="flex flex-row items-center gap-2">
                                    <div class="flex-1">
                                        <p class="text-highlighted" :title="`${$t('common.id')}: ${perm.id}`">
                                            {{ $t(`perms.${perm.namespace}.${perm.service}.${perm.name}.key`) }}
                                        </p>
                                        <p class="text-base-500">
                                            {{ $t(`perms.${perm.namespace}.${perm.service}.${perm.name}.description`) }}
                                        </p>
                                    </div>

                                    <UFieldGroup class="inline-flex flex-initial">
                                        <UButton color="success" variant="solid" icon="i-mdi-check" disabled />
                                    </UFieldGroup>
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

                        <UAccordion v-else class="p-1" :items="namespace.services" type="multiple" default-open>
                            <template #content="{ item: service }">
                                <div class="flex flex-col divide-y divide-default">
                                    <div
                                        v-for="perm in getPermissionsForService(service)"
                                        :key="perm.id"
                                        class="flex flex-col gap-1"
                                    >
                                        <div class="flex flex-row items-center gap-2">
                                            <div class="flex-1">
                                                <p class="text-highlighted" :title="`${$t('common.id')}: ${perm.id}`">
                                                    {{ $t(`perms.${perm.namespace}.${perm.service}.${perm.name}.key`) }}
                                                </p>
                                                <p class="text-base-500">
                                                    {{ $t(`perms.${perm.namespace}.${perm.service}.${perm.name}.description`) }}
                                                </p>
                                            </div>

                                            <UFieldGroup class="inline-flex flex-initial">
                                                <UButton color="success" variant="solid" icon="i-mdi-check" disabled />
                                            </UFieldGroup>
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
                    </template>
                </UAccordion>
            </div>
        </div>
    </div>
</template>
