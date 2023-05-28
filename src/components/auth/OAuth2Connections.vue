<script lang="ts" setup>
import { XCircleIcon } from '@heroicons/vue/24/solid';
import { RpcError } from 'grpc-web';
import { OAuth2Account, OAuth2Provider } from '~~/gen/ts/resources/accounts/oauth2';
import OAuth2ConnectButton from './OAuth2ConnectButton.vue';

const { $grpc } = useNuxtApp();

const props = defineProps<{
    providers: OAuth2Provider[];
    connections: OAuth2Account[];
}>();

const emit = defineEmits<{
    (e: 'disconnected', provider: string): void;
}>();

function getProviderConnection(provider: string): undefined | OAuth2Account {
    return props.connections.find((v) => v.providerName === provider);
}

async function disconnect(provider: OAuth2Provider): Promise<void> {
    return new Promise(async (res, rej) => {
        try {
            await $grpc.getAuthClient().deleteOAuth2Connection({
                provider: provider.name,
            });

            emit('disconnected', provider.name);

            return res();
        } catch (e) {
            $grpc.handleError(e as RpcError);
            return rej(e as RpcError);
        }
    });
}
</script>

<template>
    <div class="overflow-hidden bg-base-800 shadow sm:rounded-lg text-neutral mt-3">
        <div class="px-4 py-5 sm:px-6">
            <h3 class="text-base font-semibold leading-6">
                {{ $t('components.auth.oauth2_connections.title') }}
            </h3>
            <p class="mt-1 max-w-2xl text-sm">
                {{ $t('components.auth.oauth2_connections.subtitle') }}
            </p>
        </div>
        <div v-if="providers && providers.length > 0" class="border-t border-base-400 px-4 py-5 sm:p-0">
            <dl class="sm:divide-y sm:divide-base-400">
                <div class="py-4 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-6 sm:py-5" v-for="prov in providers">
                    <dt class="text-sm font-medium">
                        <NuxtLink :external="true" :to="prov.homepage" target="_blank">
                            {{ prov.label }}
                        </NuxtLink>
                    </dt>
                    <dd class="mt-1 text-sm sm:col-span-2 sm:mt-0">
                        <div v-if="getProviderConnection(prov.name) !== undefined" class="flex items-center justify-between">
                            <img
                                :src="getProviderConnection(prov.name)!.avatar"
                                alt="Avatar"
                                class="w-auto h-10 rounded-full hover:transition-colors text-base-300 bg-base-800 fill-base-300 hover:text-base-100 hover:fill-base-100"
                            />
                            <span class="text-left" :title="`ID: ${getProviderConnection(prov.name)?.externalId}`">
                                {{ getProviderConnection(prov.name)?.username }}
                            </span>

                            <button @click="disconnect(prov)">
                                <XCircleIcon class="w-6 h-6 mx-auto text-neutral" />
                                {{ $t('common.disconnect') }}
                            </button>
                        </div>
                        <div v-else>
                            <OAuth2ConnectButton :provider="prov" />
                        </div>
                    </dd>
                </div>
            </dl>
        </div>
    </div>
</template>
