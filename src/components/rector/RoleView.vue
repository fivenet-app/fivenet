<script lang="ts" setup>
import { Permission, Role, RoleAttribute } from '@fivenet/gen/resources/permissions/permissions_pb';
import { RpcError } from 'grpc-web';
import { UpdateRolePermsRequest, DeleteRoleRequest, GetPermissionsRequest, GetRoleRequest, PermsUpdate, PermItem } from '@fivenet/gen/services/rector/rector_pb';
import { ChevronDownIcon, CheckIcon, XMarkIcon, MinusIcon } from '@heroicons/vue/24/solid';
import { TrashIcon } from '@heroicons/vue/20/solid';
import Divider from '~/components/partials/Divider.vue';
import RoleViewAttr from '~/components/rector/RoleViewAttr.vue';
import { useNotificationsStore } from '~/store/notifications';
import { Disclosure, DisclosureButton, DisclosurePanel } from '@headlessui/vue'

const { $grpc } = useNuxtApp();

const notifications = useNotificationsStore();

const { t } = useI18n();

const props = defineProps({
    roleId: {
        type: Number,
        required: true,
    }
});

const role = ref<Role>();

async function getRole(): Promise<void> {
    const req = new GetRoleRequest();
    req.setId(props.roleId);

    try {
        const resp = await $grpc.getRectorClient().getRole(req, null);
        // TODO Take care of attributes

        role.value = resp.getRole();
    } catch (e) {
        $grpc.handleRPCError(e as RpcError);
        return;
    }
}

async function deleteRole(): Promise<void> {
    return new Promise(async (res, rej) => {
        const req = new DeleteRoleRequest();
        req.setId(props.roleId);

        try {
            await $grpc.getRectorClient().
                deleteRole(req, null);

            notifications.dispatchNotification({
                title: t('notifications.rector.role_deleted.title'),
                content: t('notifications.rector.role_deleted.content'),
                type: 'success'
            });
            await navigateTo({ name: 'rector-roles' });
            return res();
        } catch (e) {
            $grpc.handleRPCError(e as RpcError);
            return rej(e as RpcError);
        }
    });
}

const permList = ref<Permission[]>([]);
const permCategories = ref<Set<string>>(new Set());
const permStates = ref<Map<number, boolean | undefined>>(new Map());

const attrList = ref<RoleAttribute[]>([]);
const attrStates = ref<Map<number, (string | number)[]>>(new Map());

async function getPermissions(): Promise<void> {
    const req = new GetPermissionsRequest();

    try {
        const resp = await $grpc.getRectorClient().getPermissions(req, null);
        permList.value = resp.getPermissionsList();
        attrList.value = resp.getAttributesList();

        genPermissionCategories();
    } catch (e) {
        $grpc.handleRPCError(e as RpcError);
        return;
    }
}

async function genPermissionCategories(): Promise<void> {
    permList.value.forEach(perm => {
        permCategories.value.add(perm.getCategory());
    });
}

async function propogatePermissionStates(): Promise<void> {
    role.value?.getPermissionsList().forEach(perm => {
        permStates.value.set(perm.getId(), Boolean(perm.getVal()));
    });
}

async function updatePermissionState(perm: number, state: boolean | undefined): Promise<void> {
    permStates.value.set(perm, state);
}

async function updatePermissions(): Promise<void> {
    const currentPermissions = role.value?.getPermissionsList().map(p => p.getId()) ?? [];

    const permsToUpdate: PermItem[] = [];
    const permsToRemove: number[] = [];
    permStates.value.forEach((state, perm) => {
        if (state !== undefined) {
            const p = role.value?.getPermissionsList().find(v => v.getId() == perm);

            if (p?.getVal() != state) {
                const item = new PermItem();
                permsToUpdate.push(item.setId(perm).setVal(state));
            }
        }
        if (state === undefined && currentPermissions.includes(perm)) permsToRemove.push(perm);
    });

    if (permsToUpdate.length == 0 && permsToRemove.length == 0) return;

    const req = new UpdateRolePermsRequest();
    req.setId(props.roleId);

    const perms = new PermsUpdate();
    perms.setToUpdateList(permsToUpdate);
    perms.setToRemoveList(permsToRemove);
    req.setPerms(perms);

    try {
        await $grpc.getRectorClient().updateRolePerms(req, null);

        notifications.dispatchNotification({
            title: t('notifications.rector.role_updated.title'),
            content: t('notifications.rector.role_updated.content'),
            type: 'success'
        });
    } catch (e) {
        $grpc.handleRPCError(e as RpcError);
        return;
    }
}

onMounted(async () => {
    await getRole();
    await getPermissions();
    propogatePermissionStates();
});
</script>

<style scoped>
.upsidedown {
    transform: rotate(180deg);
}
</style>

<template>
    <div class="py-4 max-w-7xl mx-auto">
        <div class="px-2 sm:px-6 lg:px-8">
            <div v-if="role">
                <h2 class="text-3xl text-white">
                    {{ role?.getJobLabel()! }} - {{ role?.getJobGradeLabel() }}
                    <button v-can="'RectorService.DeleteRole'" @click="deleteRole()">
                        <TrashIcon class="w-6 h-6 mx-auto text-neutral" />
                    </button>
                </h2>
                <Divider label="Permissions" />
                <div class="py-2 flex flex-col gap-4">
                    <Disclosure as="div" class="text-neutral hover:border-neutral/70 border-neutral/20"
                        v-for="category in permCategories" :key="category" v-slot="{ open }">
                        <DisclosureButton
                            :class="[open ? 'rounded-t-lg border-b-0' : 'rounded-lg', 'flex w-full items-start justify-between text-left border-2 p-2 border-inherit transition-colors']">
                            <span class="text-base font-semibold leading-7">
                                {{ $t(`perms.${category}.category`) }}
                            </span>
                            <span class="ml-6 flex h-7 items-center">
                                <ChevronDownIcon :class="[open ? 'upsidedown' : '', 'h-6 w-6 transition-transform']"
                                    aria-hidden="true" />
                            </span>
                        </DisclosureButton>
                        <DisclosurePanel
                            class="px-4 pb-2 border-2 border-t-0 rounded-b-lg transition-colors border-inherit -mt-2">
                            <div class="flex flex-col gap-2 max-w-4xl mx-auto my-2">
                                <div v-for="(perm, idx) in permList.filter(p => p.getCategory() === category)" :key="perm.getId()">
                                    <div class="flex flex-row gap-4">
                                        <div class="flex flex-1 flex-col my-auto">
                                            <span class="truncate lg:max-w-full max-w-xs">
                                                {{ $t(`perms.${perm.getCategory()}.${perm.getName()}.key`) }}
                                            </span>
                                            <span class="text-base-500 truncate lg:max-w-full max-w-xs">
                                                {{ $t(`perms.${perm.getCategory()}.${perm.getName()}.description`) }}
                                            </span>
                                        </div>
                                        <div class="flex flex-initial flex-row max-h-8 my-auto">
                                            <button :data-active="permStates.has(perm.getId()) ? permStates.get(perm.getId()) : false"
                                                @click="updatePermissionState(perm.getId(), true)"
                                                class="transition-colors rounded-l-lg p-1 bg-success-600/50 data-[active=true]:bg-success-600 text-base-300 data-[active=true]:text-neutral hover:bg-success-600/70">
                                                <CheckIcon class="w-6 h-6" />
                                            </button>
                                            <button :data-active="!permStates.has(perm.getId()) || permStates.get(perm.getId()) === undefined"
                                                @click="updatePermissionState(perm.getId(), undefined)"
                                                class="transition-colors p-1 bg-base-700 data-[active=true]:bg-base-500 text-base-300 data-[active=true]:text-neutral hover:bg-base-600">
                                                <MinusIcon class="w-6 h-6" />
                                            </button>
                                            <button
                                                :data-active="permStates.has(perm.getId()) ? (permStates.get(perm.getId()) !== undefined && !permStates.get(perm.getId())) : false"
                                                @click="updatePermissionState(perm.getId(), false)"
                                                class="transition-colors rounded-r-lg p-1 bg-error-600/50 data-[active=true]:bg-error-600 text-base-300 data-[active=true]:text-neutral hover:bg-error-600/70">
                                                <XMarkIcon class="w-6 h-6" />
                                            </button>
                                        </div>
                                    </div>
                                    <RoleViewAttr :attribute="attrList.find(a => a.getPermissionId() === perm.getId())" v-model:states="attrStates" :disabled="permStates.get(perm.getId()) !== true" />
                                    <div v-if="idx !== permList.filter(p => p.getCategory() === category).length - 1"
                                        class="w-full border-t border-neutral/20 mt-2" />
                                </div>
                            </div>
                        </DisclosurePanel>
                    </Disclosure>
                    <button type="button" @click="updatePermissions()"
                        class="inline-flex px-3 py-2 text-center justify-center transition-colors font-semibold rounded-md bg-primary-600 text-neutral hover:bg-primary-500 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-success-600">
                        {{ $t('common.save', 1) }}
                    </button>
                </div>
            </div>
        </div>
    </div>
</template>
