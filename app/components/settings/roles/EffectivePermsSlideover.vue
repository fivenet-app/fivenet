<script lang="ts" setup>
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import DataNoDataBlock from '~/components/partials/data/DataNoDataBlock.vue';
import DataPendingBlock from '~/components/partials/data/DataPendingBlock.vue';
import RoleViewAttr from '~/components/settings/roles/RoleViewAttr.vue';
import { getSettingsSettingsClient } from '~~/gen/ts/clients';
import type { RoleAttribute } from '~~/gen/ts/resources/permissions/attributes';
import type { Permission } from '~~/gen/ts/resources/permissions/permissions';
import type { GetEffectivePermissionsResponse } from '~~/gen/ts/services/settings/settings';
import { isEmptyAttributes } from './helpers';

const settingsSettingsClient = await getSettingsSettingsClient();

const props = defineProps<{
    roleId: number;
}>();

defineEmits<{
    close: [boolean];
}>();

const { t } = useI18n();

async function getEffectivePermissions(roleId: number): Promise<GetEffectivePermissionsResponse> {
    try {
        const call = settingsSettingsClient.getEffectivePermissions({
            roleId: roleId,
        });
        const { response } = await call;

        permList.value = response.permissions;
        attrList.value = response.attributes;

        await genPermissionCategories();

        return response;
    } catch (e) {
        handleGRPCError(e as RpcError);
        throw e;
    }
}

const {
    data: role,
    status,
    refresh,
    error,
} = useLazyAsyncData(`settings-roles-${props.roleId}-effective`, () => getEffectivePermissions(props.roleId));

const permList = ref<Permission[]>([]);
const permCategories = ref<Set<string>>(new Set());
const permStates = ref(new Map<number, boolean | undefined>());

const attrList = ref<RoleAttribute[]>([]);

async function genPermissionCategories(): Promise<void> {
    permCategories.value.clear();

    permList.value.forEach((perm) => {
        permCategories.value.add(perm.category);
    });
}

async function propogatePermissionStates(): Promise<void> {
    permStates.value.clear();

    permList.value.forEach((perm) => {
        permStates.value.set(perm.id, Boolean(perm.val));
    });
}

async function initializeRoleView(): Promise<void> {
    await propogatePermissionStates();
}

watch(role, async () => initializeRoleView());

const permCategoriesSorted = computed(() =>
    [...permCategories.value.entries()].map((category) => {
        return {
            label: t(`perms.${category[1]}.category`),
            category: category[0],
        };
    }),
);
</script>

<template>
    <USlideover
        :title="`${$t('common.effective_permissions')}: ${role?.role?.jobLabel!} - ${role?.role?.jobGradeLabel} (${role?.role?.grade})`"
    >
        <template #actions>
            <UButton variant="link" trailing-icon="i-mdi-refresh" color="primary" @click="refresh()" />
        </template>

        <template #body>
            <UAlert
                class="mb-2"
                :ui="{
                    icon: 'size-6',
                }"
                icon="i-mdi-information-outline"
            >
                <template #description>
                    {{ $t('components.settings.role_view.effective_permissions.description') }}
                </template>
            </UAlert>

            <USeparator class="mb-2" />

            <DataPendingBlock v-if="isRequestPending(status)" :message="$t('common.loading', [$t('common.role')])" />
            <DataErrorBlock
                v-else-if="error"
                :title="$t('common.unable_to_load', [$t('common.role')])"
                :error="error"
                :retry="refresh"
            />
            <DataNoDataBlock v-else-if="!role" :type="$t('common.role')" icon="i-mdi-comment-text-multiple" />

            <div v-else class="flex flex-col gap-2">
                <UCard v-for="category in permCategoriesSorted" :key="category.category">
                    <template #header>
                        <h3 class="text-xl" :title="`ID: ${role.role?.id}`">
                            {{ category.label }}
                        </h3>
                    </template>

                    <div class="flex flex-col divide-y divide-default">
                        <div
                            v-for="perm in permList.filter((p) => p.category === category.category)"
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
                                    <UButton
                                        color="green"
                                        :variant="permStates.get(perm.id) ? 'solid' : 'soft'"
                                        icon="i-mdi-check"
                                        :disabled="true"
                                    />

                                    <UButton
                                        color="neutral"
                                        :variant="
                                            !permStates.has(perm.id) || permStates.get(perm.id) === undefined ? 'solid' : 'soft'
                                        "
                                        icon="i-mdi-minus"
                                        :disabled="true"
                                    />

                                    <UButton
                                        color="error"
                                        :variant="
                                            permStates.get(perm.id) !== undefined && !permStates.get(perm.id) ? 'solid' : 'soft'
                                        "
                                        icon="i-mdi-close"
                                        :disabled="true"
                                    />
                                </UButtonGroup>
                            </div>

                            <template v-for="(attr, idx) in attrList" :key="attr.attrId">
                                <RoleViewAttr
                                    v-if="attr.permissionId === perm.id && !isEmptyAttributes(attr.maxValues)"
                                    v-model="attrList[idx]!"
                                    :permission="perm"
                                    disabled
                                    default-open
                                    hide-fine-grained
                                />
                            </template>
                        </div>
                    </div>
                </UCard>
            </div>
        </template>

        <template #footer>
            <UButtonGroup class="inline-flex w-full">
                <UButton class="flex-1" color="neutral" block :label="$t('common.close', 1)" @click="$emit('close', false)" />
            </UButtonGroup>
        </template>
    </USlideover>
</template>
