<script lang="ts" setup>
import { RpcError } from '@protobuf-ts/runtime-rpc';
import { CloseCircleIcon } from 'mdi-vue3';
import { OAuth2Account, OAuth2Provider } from '~~/gen/ts/resources/accounts/oauth2';
import OAuth2ConnectButton from '~/components/auth/account/OAuth2ConnectButton.vue';

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
    try {
        await $grpc.getAuthClient().deleteOAuth2Connection({
            provider: provider.name,
        });

        emit('disconnected', provider.name);
    } catch (e) {
        $grpc.handleError(e as RpcError);
        throw e;
    }
}
</script>

<template>
    <div class="mt-3 overflow-hidden bg-base-800 text-neutral shadow sm:rounded-lg">
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
                <div
                    v-for="provider in providers"
                    :key="provider.name"
                    class="py-4 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-6 sm:py-5"
                >
                    <dt class="text-sm font-medium">
                        <NuxtLink :external="true" :to="provider.homepage" target="_blank">
                            {{ provider.label }}
                        </NuxtLink>
                    </dt>
                    <dd class="mt-1 text-sm sm:col-span-2 sm:mt-0">
                        <div
                            v-if="getProviderConnection(provider.name) !== undefined"
                            class="flex items-center justify-between"
                        >
                            <img
                                :src="getProviderConnection(provider.name)!.avatar"
                                alt="Avatar"
                                class="h-10 w-auto rounded-full bg-base-800 fill-base-300 text-base-300 ring-2 ring-neutral hover:fill-base-100 hover:text-base-100 hover:transition-colors"
                            />
                            <span class="text-left" :title="`ID: ${getProviderConnection(provider.name)?.externalId}`">
                                {{ getProviderConnection(provider.name)?.username }}
                            </span>

                            <button @click="disconnect(provider)">
                                <CloseCircleIcon class="mx-auto h-5 w-5 text-neutral" />
                                {{ $t('common.disconnect') }}
                            </button>
                        </div>
                        <div v-else>
                            <OAuth2ConnectButton :provider="provider" />
                        </div>
                    </dd>
                </div>
            </dl>
        </div>
    </div>
</template>
