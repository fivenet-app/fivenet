<script lang="ts" setup>
import { RpcError } from '@protobuf-ts/runtime-rpc/build/types';
import { SelectIcon } from 'mdi-vue3';
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import DataNoDataBlock from '~/components/partials/data/DataNoDataBlock.vue';
import DataPendingBlock from '~/components/partials/data/DataPendingBlock.vue';
import { Role } from '~~/gen/ts/resources/permissions/permissions';
import AttrRolesListEntry from './AttrRolesListEntry.vue';
import AttrView from './AttrView.vue';

const { $grpc } = useNuxtApp();

const { data: roles, pending, refresh, error } = useLazyAsyncData('rector-roles', () => getRoles());

async function getRoles(): Promise<Role[]> {
    return new Promise(async (res, rej) => {
        try {
            const call = $grpc.getRectorClient().getRoles({
                lowestRank: true,
            });
            const { response } = await call;

            return res(response.roles);
        } catch (e) {
            $grpc.handleError(e as RpcError);
            return rej(e as RpcError);
        }
    });
}

const selectedRole = ref<Role | undefined>();
</script>

<template>
    <div class="py-2">
        <div class="px-1 sm:px-2 lg:px-4">
            <div class="flex flex-col lg:flex-row">
                <div class="flow-root mt-2 basis-1/3">
                    <div class="mx-0 -my-2 overflow-x-auto">
                        <div class="inline-block min-w-full py-2 align-middle px-1">
                            <DataPendingBlock v-if="pending" :message="$t('common.loading', [$t('common.role', 2)])" />
                            <DataErrorBlock
                                v-else-if="error"
                                :title="$t('common.unable_to_load', [$t('common.role', 2)])"
                                :retry="refresh"
                            />
                            <DataNoDataBlock v-else-if="roles && roles.length === 0" :type="$t('common.role', 2)" />
                            <div v-else>
                                <table class="min-w-full divide-y divide-base-600">
                                    <thead>
                                        <tr>
                                            <th scope="col" class="py-3.5 px-2 text-left text-sm font-semibold text-neutral">
                                                {{ $t('common.name') }}
                                            </th>
                                            <th
                                                scope="col"
                                                class="relative py-3.5 pl-3 pr-4 sm:pr-0 text-right text-sm font-semibold text-neutral"
                                            >
                                                {{ $t('common.action', 2) }}
                                            </th>
                                        </tr>
                                    </thead>
                                    <tbody class="divide-y divide-base-800">
                                        <AttrRolesListEntry
                                            v-for="role in roles"
                                            :role="role"
                                            @selected="selectedRole = role"
                                            :class="selectedRole?.id === role.id ? 'bg-base-800' : ''"
                                        />
                                    </tbody>
                                    <thead>
                                        <tr>
                                            <th scope="col" class="py-3.5 px-2 text-left text-sm font-semibold text-neutral">
                                                {{ $t('common.name') }}
                                            </th>
                                            <th
                                                scope="col"
                                                class="relative py-3.5 pl-3 pr-4 sm:pr-0 text-right text-sm font-semibold text-neutral"
                                            >
                                                {{ $t('common.action', 2) }}
                                            </th>
                                        </tr>
                                    </thead>
                                </table>
                            </div>
                        </div>
                    </div>
                </div>
                <div class="flex basis-2/3 w-full ml-2">
                    <template v-if="selectedRole">
                        <AttrView :role-id="selectedRole.id" @deleted="refresh()" />
                    </template>
                    <template v-else>
                        <DataNoDataBlock :icon="SelectIcon" :message="$t('common.none_selected', [$t('common.job', 2)])" />
                    </template>
                </div>
            </div>
        </div>
    </div>
</template>
