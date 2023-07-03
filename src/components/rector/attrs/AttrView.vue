<script lang="ts" setup>
import { Disclosure, DisclosureButton, DisclosurePanel } from '@headlessui/vue';
import SvgIcon from '@jamescoyle/vue-icon';
import { mdiChevronDown } from '@mdi/js';
import { RpcError } from '@protobuf-ts/runtime-rpc/build/types';
import Divider from '~/components/partials/elements/Divider.vue';
import { useNotificationsStore } from '~/store/notifications';
import { Job } from '~~/gen/ts/resources/jobs/jobs';
import { AttributeValues, Permission, Role, RoleAttribute } from '~~/gen/ts/resources/permissions/permissions';
import { AttrsUpdate } from '~~/gen/ts/services/rector/rector';
import AttrViewAttr from './AttrViewAttr.vue';

const { $grpc } = useNuxtApp();

const notifications = useNotificationsStore();

const props = defineProps<{
    roleId: bigint;
}>();

const role = ref<Role>();

const permList = ref<Permission[]>([]);
const permCategories = ref<Set<string>>(new Set());

const attrList = ref<RoleAttribute[]>([]);
const attrStates = ref<Map<bigint, AttributeValues | undefined>>(new Map());

const jobs = ref<Job[]>([]);

async function getRole(): Promise<void> {
    return new Promise(async (res, rej) => {
        try {
            const call = $grpc.getRectorClient().getRole({
                id: props.roleId,
            });
            const { response } = await call;

            role.value = response.role;

            return res();
        } catch (e) {
            $grpc.handleError(e as RpcError);
            return rej(e as RpcError);
        }
    });
}

async function prepareAttributeData(): Promise<void> {
    return new Promise(async (res, rej) => {
        try {
            const call = $grpc.getCompletorClient().completeJobs({});
            const { response } = await call;

            jobs.value = response.jobs;

            return res();
        } catch (e) {
            $grpc.handleError(e as RpcError);
            return rej(e as RpcError);
        }
    });
}

async function getPermissions(): Promise<void> {
    return new Promise(async (res, rej) => {
        try {
            const call = $grpc.getRectorClient().getPermissions({
                roleId: props.roleId,
            });
            const { response } = await call;

            permList.value = response.permissions;
            attrList.value = response.attributes;

            genPermissionCategories();

            attrStates.value.clear();
            attrList.value.forEach((attr) => {
                attrStates.value.set(attr.attrId, attr.maxValues);
            });

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

async function updateLimits(): Promise<void> {
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

            initializeRoleView();

            return res();
        } catch (e) {
            $grpc.handleError(e as RpcError);
            return rej(e as RpcError);
        }
    });
}

async function initializeRoleView(): Promise<void> {
    await getRole();
    await prepareAttributeData();
    await getPermissions();
}

onMounted(async () => {
    initializeRoleView();
});
</script>

<template>
    <div class="py-4 max-w-7xl mx-auto">
        <div class="px-2 sm:px-6 lg:px-8">
            <div v-if="role">
                <h2 class="text-3xl text-white">{{ role?.jobLabel! }}</h2>
                <Divider label="Permissions" />
                <div class="py-2 flex flex-col gap-4">
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
                                <SvgIcon
                                    :class="[open ? 'upsidedown' : '', 'h-6 w-6 transition-transform']"
                                    aria-hidden="true"
                                    type="mdi"
                                    :path="mdiChevronDown"
                                />
                            </span>
                        </DisclosureButton>
                        <DisclosurePanel
                            class="px-4 pb-2 border-2 border-t-0 rounded-b-lg transition-colors border-inherit -mt-2"
                        >
                            <div class="flex flex-col gap-2 max-w-4xl mx-auto my-2">
                                <div
                                    v-for="(perm, idx) in permList.filter(
                                        (p) =>
                                            p.category === category &&
                                            attrList.filter((a) => a.permissionId === p.id).length > 0
                                    )"
                                    :key="perm.id?.toString()"
                                    class="flex flex-col gap-2"
                                >
                                    <div class="flex flex-row gap-4">
                                        <div class="flex flex-1 flex-col my-auto">
                                            <span class="truncate lg:max-w-full max-w-xs">
                                                {{ $t(`perms.${perm.category}.${perm.name}.key`) }}
                                            </span>
                                            <span class="text-base-500 truncate lg:max-w-full max-w-xs">
                                                {{ $t(`perms.${perm.category}.${perm.name}.description`) }}
                                            </span>
                                        </div>
                                    </div>
                                    <AttrViewAttr
                                        v-for="attr in attrList.filter((a) => a.permissionId === perm.id)"
                                        :attribute="attr"
                                        :permission="perm"
                                        v-model:states="attrStates"
                                        :jobs="jobs"
                                    />
                                    <div
                                        v-if="idx !== permList.filter((p) => p.category === category).length - 1"
                                        class="w-full border-t border-neutral/20"
                                    />
                                </div>
                            </div>
                        </DisclosurePanel>
                    </Disclosure>
                    <button
                        type="button"
                        @click="updateLimits()"
                        class="inline-flex px-3 py-2 text-center justify-center transition-colors font-semibold rounded-md bg-primary-500 text-neutral hover:bg-primary-400 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-success-600"
                    >
                        {{ $t('common.save', 1) }}
                    </button>
                </div>
            </div>
        </div>
    </div>
</template>
