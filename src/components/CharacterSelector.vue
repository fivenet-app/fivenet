<script lang="ts">
import { defineComponent } from 'vue';
import CharacterSelectorCard from './CharacterSelectorCard.vue';
import { AccountServiceClient } from '@arpanet/gen/auth/AuthServiceClientPb';
import { GetCharactersRequest } from '@arpanet/gen/auth/auth_pb';
import { Character } from '@arpanet/gen/common/character_pb';
import authInterceptor from '../grpcauth';
import { RpcError } from 'grpc-web';
import { mapActions } from 'vuex';
import config from '../config';

export default defineComponent({
    components: {
        CharacterSelectorCard,
    },
    data() {
        return {
            'client': new AccountServiceClient(config.apiProtoURL, null, {
                unaryInterceptors: [authInterceptor],
                streamInterceptors: [authInterceptor],
            }),
            'chars': [] as Array<Character>,
        };
    },
    methods: {
        ...mapActions([
            'updateActiveChar',
            'updateActiveCharIdentifier',
        ]),
        async fetchCharacters() {
            return this.client.
                getCharacters(new GetCharactersRequest(), null).
                then((resp) => {
                    this.chars = resp.getCharsList();
                }).catch((err: RpcError) => {
                    authInterceptor.handleError(err, this.$route);
                });
        },
    },
    beforeMount() {
        this.updateActiveChar(null);
        this.updateActiveCharIdentifier(null);

        // Fetch user's characters
        this.fetchCharacters();
    },
});

</script>

<template>
    <div class="grid place-items-center">
        <div class="flex w-full">
            <CharacterSelectorCard v-for="char in chars" :char="char" :identifier="char.getIdentifier()"
                :key="char.getIdentifier()" />
        </div>
    </div>
</template>
