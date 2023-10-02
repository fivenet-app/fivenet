<script lang="ts" setup>
import { Disclosure, DisclosureButton, DisclosurePanel } from '@headlessui/vue';
import { RpcError } from '@protobuf-ts/runtime-rpc/build/types';
import { ChevronDownIcon } from 'mdi-vue3';
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import DataNoDataBlock from '~/components/partials/data/DataNoDataBlock.vue';
import DataPendingBlock from '~/components/partials/data/DataPendingBlock.vue';
import Divider from '~/components/partials/elements/Divider.vue';
import { useNotificatorStore } from '~/store/notificator';
import { AttributeValues, Permission, Role, RoleAttribute } from '~~/gen/ts/resources/permissions/permissions';
import { AttrsUpdate } from '~~/gen/ts/services/rector/rector';
import AttrViewAttr from './AttrViewAttr.vue';

const props = defineProps<{
    roleId: bigint;
}>();

const { $grpc } = useNuxtApp();

const notifications = useNotificatorStore();

const {
    data: role,
    pending,
    refresh,
    error,
} = useLazyAsyncData(`rector-attrs-${props.roleId.toString()}`, () => getRole(props.roleId));

const changed = ref(false);

const permList = ref<Permission[]>([]);
const permCategories = ref<Set<string>>(new Set());

const attrList = ref<RoleAttribute[]>([]);
const attrStates = ref<Map<bigint, AttributeValues | undefined>>(new Map());

async function getRole(id: bigint): Promise<Role> {
    return new Promise(async (res, rej) => {
        try {
            const call = $grpc.getRectorClient().getRole({
                id: id,
            });
            const { response } = await call;

            if (response.role === undefined) return rej();

            return res(response.role!);
        } catch (e) {
            $grpc.handleError(e as RpcError);
            return rej(e as RpcError);
        }
    });
}

async function getPermissions(roleId: bigint): Promise<void> {
    return new Promise(async (res, rej) => {
        try {
            const call = $grpc.getRectorClient().getPermissions({
                roleId: roleId,
            });
            const { response } = await call;

            permList.value = response.permissions;
            attrList.value = response.attributes;

            genPermissionCategories();

            return res();
        } catch (e) {
            $grpc.handleError(e as RpcError);
            return rej(e as RpcError);
        }
    });
}

async function genPermissionCategories(): Promise<void> {
    permCategories.value.clear();

    permList.value.forEach((perm) => {
        if (attrList.value.filter((a) => a.permissionId === perm.id).length > 0) {
            permCategories.value.add(perm.category);
        }
    });
}

async function updateAttrs(): Promise<void> {
    return new Promise(async (res, rej) => {
        const attrs: AttrsUpdate = {
            toRemove: [],
            toUpdate: [],
        };
        attrStates.value.forEach((state, attr) => {
            if (state !== undefined) {
                attrs.toUpdate.push({
                    roleId: role.value!.id,
                    attrId: attr,
                    category: '',
                    key: '',
                    name: '',
                    permissionId: 0n,
                    type: '',
                    maxValues: state,
                });
            } else if (state === undefined) {
                attrs.toRemove.push({
                    roleId: role.value!.id,
                    attrId: attr,
                    category: '',
                    key: '',
                    name: '',
                    permissionId: 0n,
                    type: '',
                });
            }
        });

        if (attrs.toUpdate.length === 0 && attrs.toRemove.length === 0) {
            changed.value = false;
            return res();
        }

        try {
            await $grpc.getRectorClient().updateRoleLimits({
                roleId: props.roleId,
                attrs: attrs,
            });

            notifications.dispatchNotification({
                title: { key: 'notifications.rector.role_updated.title', parameters: [] },
                content: { key: 'notifications.rector.role_updated.content', parameters: [] },
                type: 'success',
            });

            changed.value = false;
            refresh();

            return res();
        } catch (e) {
            $grpc.handleError(e as RpcError);
            return rej(e as RpcError);
        }
    });
}

async function initializeRoleView(): Promise<void> {
    await getPermissions(props.roleId);

    attrStates.value.clear();
    attrList.value.forEach((attr) => {
        attrStates.value.set(attr.attrId, attr.maxValues);
    });

    role.value?.attributes.forEach((attr) => {
        attrStates.value.set(attr.attrId, attr.maxValues);
    });
}

watch(role, async () => initializeRoleView());

watch(props, () => {
    if (role.value === null || role.value.id !== props.roleId) {
        refresh();
    }
});
</script>

<template>
    <div class="py-4 w-full">
        <div class="px-1 sm:px-2 lg:px-4">
            <DataPendingBlock v-if="pending" :message="$t('common.loading', [$t('common.attributes', 2)])" />
            <DataErrorBlock
                v-else-if="error"
                :title="$t('common.unable_to_load', [$t('common.attributes', 2)])"
                :retry="refresh"
            />
            <DataNoDataBlock v-else-if="!role" :type="$t('common.attributes', 2)" />
            <div v-else>
                <h2 class="text-3xl text-white" :title="`ID: ${role.id}`">{{ role?.jobLabel! }}</h2>
                <Divider label="Permissions" />
                <div class="py-2 flex flex-col gap-4">
                    <button
                        type="button"
                        @click="updateAttrs()"
                        :disabled="!changed"
                        class="inline-flex px-3 py-2 text-center justify-center transition-colors font-semibold rounded-md text-neutral focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2"
                        :class="
                            !changed
                                ? 'disabled bg-base-500 hover:bg-base-400 focus-visible:outline-base-500'
                                : 'bg-primary-500 hover:bg-primary-400 focus-visible:outline-primary-500'
                        "
                    >
                        {{ $t('common.save', 1) }}
                    </button>

                    <Disclosure
                        as="div"
                        class="text-neutral hover:border-neutral/70 border-neutral/20"
                        v-for="category in permCategories"
                        :key="category"
                        v-slot="{ open }"
                    >
                        <DisclosureButton
                            :class="[
                                open ? 'rounded-t-lg border-b-0' : 'rounded-lg',
                                'flex w-full items-start justify-between text-left border-2 p-2 border-inherit transition-colors',
                            ]"
                        >
                            <span class="text-base font-semibold leading-7">
                                {{ $t(`perms.${category}.category`) }}
                            </span>
                            <span class="ml-6 flex h-7 items-center">
                                <ChevronDownIcon
                                    :class="[open ? 'upsidedown' : '', 'h-6 w-6 transition-transform']"
                                    aria-hidden="true"
                                />
                            </span>
                        </DisclosureButton>
                        <DisclosurePanel
                            class="px-4 pb-2 border-2 border-t-0 rounded-b-lg transition-colors border-inherit -mt-2"
                        >
                            <div class="flex flex-col gap-2 mx-auto my-2">
                                <div
                                    v-for="(perm, idx) in permList.filter(
                                        (p) =>
                                            p.category === category &&
                                            attrList.filter((a) => a.permissionId === p.id).length > 0,
                                    )"
                                    :key="perm.id?.toString()"
                                    class="flex flex-col gap-2"
                                >
                                    <div class="flex flex-row gap-4">
                                        <div class="flex flex-1 flex-col my-auto">
                                            <span class="truncate">
                                                {{ $t(`perms.${perm.category}.${perm.name}.key`) }}
                                            </span>
                                            <span class="text-base-500 truncate">
                                                {{ $t(`perms.${perm.category}.${perm.name}.description`) }}
                                            </span>
                                        </div>
                                    </div>
                                    <AttrViewAttr
                                        v-for="attr in attrList.filter((a) => a.permissionId === perm.id)"
                                        :attribute="attr"
                                        :permission="perm"
                                        v-model:states="attrStates"
                                        @changed="changed = true"
                                    />
                                    <div
                                        v-if="idx !== permList.filter((p) => p.category === category).length - 1"
                                        class="w-full border-t border-neutral/20"
                                    />
                                </div>
                            </div>
                        </DisclosurePanel>
                    </Disclosure>
                </div>
            </div>
        </div>
    </div>
</template>
