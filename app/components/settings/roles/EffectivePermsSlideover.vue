<script lang="ts" setup>
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import DataNoDataBlock from '~/components/partials/data/DataNoDataBlock.vue';
import DataPendingBlock from '~/components/partials/data/DataPendingBlock.vue';
import RoleViewAttr from '~/components/settings/roles/RoleViewAttr.vue';
import type { RoleAttribute } from '~~/gen/ts/resources/permissions/attributes';
import type { Permission } from '~~/gen/ts/resources/permissions/permissions';
import type { GetEffectivePermissionsResponse } from '~~/gen/ts/services/settings/settings';
import { isEmptyAttributes } from './helpers';

const props = defineProps<{
    roleId: number;
}>();

const { t } = useI18n();

const { $grpc } = useNuxtApp();

const { isOpen } = useSlideover();

async function getEffectivePermissions(roleId: number): Promise<GetEffectivePermissionsResponse> {
    try {
        const call = $grpc.settings.settings.getEffectivePermissions({
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
    pending: loading,
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
    <USlideover :ui="{ width: 'w-screen sm:max-w-3xl' }">
        <UCard
            class="flex flex-1 flex-col"
            :ui="{
                body: {
                    base: 'flex-1 min-h-[calc(100dvh-(2*var(--header-height)))] max-h-[calc(100dvh-(2*var(--header-height)))] overflow-y-auto',
                    padding: 'px-1 py-2 sm:p-2',
                },
                ring: '',
                divide: 'divide-y divide-gray-100 dark:divide-gray-800',
            }"
        >
            <template #header>
                <div class="flex flex-col gap-1">
                    <div class="flex items-center justify-between">
                        <h3 class="inline-flex gap-2 text-2xl font-semibold leading-6">
                            {{ $t('common.effective_permissions') }}: {{ role?.role?.jobLabel! }} -
                            {{ role?.role?.jobGradeLabel }} ({{ role?.role?.grade }})
                        </h3>

                        <UButton variant="link" trailing-icon="i-mdi-refresh" color="primary" @click="refresh()" />

                        <UButton class="-my-1" color="gray" variant="ghost" icon="i-mdi-window-close" @click="isOpen = false" />
                    </div>
                </div>
            </template>

            <div>
                <UAlert
                    class="mb-2"
                    :ui="{
                        icon: { base: 'size-6' },
                    }"
                    icon="i-mdi-information-outline"
                >
                    <template #description>
                        {{ $t('components.settings.role_view.effective_permissions.description') }}
                    </template>
                </UAlert>

                <UDivider class="mb-2" />

                <DataPendingBlock v-if="loading" :message="$t('common.loading', [$t('common.role')])" />
                <DataErrorBlock
                    v-else-if="error"
                    :title="$t('common.unable_to_load', [$t('common.role')])"
                    :error="error"
                    :retry="refresh"
                />
                <DataNoDataBlock v-else-if="!role" :type="$t('common.role')" icon="i-mdi-comment-text-multiple" />

                <div v-else class="flex flex-col gap-2">
                    <UCard
                        v-for="category in permCategoriesSorted"
                        :key="category.category"
                        :ui="{ body: { padding: 'px-2 py-3 sm:px-4 sm:p-2' } }"
                    >
                        <template #header>
                            <h3 class="text-xl" :title="`ID: ${role.role?.id}`">
                                {{ category.label }}
                            </h3>
                        </template>

                        <div class="flex flex-col divide-y divide-gray-100 dark:divide-gray-800">
                            <div
                                v-for="perm in permList.filter((p) => p.category === category.category)"
                                :key="perm.id"
                                class="flex flex-col gap-1"
                            >
                                <div class="flex flex-row items-center gap-2">
                                    <div class="flex-1">
                                        <p class="text-gray-900 dark:text-white" :title="`${$t('common.id')}: ${perm.id}`">
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
                                            color="black"
                                            :variant="
                                                !permStates.has(perm.id) || permStates.get(perm.id) === undefined
                                                    ? 'solid'
                                                    : 'soft'
                                            "
                                            icon="i-mdi-minus"
                                            :disabled="true"
                                        />

                                        <UButton
                                            color="error"
                                            :variant="
                                                permStates.get(perm.id) !== undefined && !permStates.get(perm.id)
                                                    ? 'solid'
                                                    : 'soft'
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
            </div>

            <template #footer>
                <UButtonGroup class="inline-flex w-full">
                    <UButton class="flex-1" color="black" block @click="isOpen = false">
                        {{ $t('common.close', 1) }}
                    </UButton>
                </UButtonGroup>
            </template>
        </UCard>
    </USlideover>
</template>
