<script lang="ts">
import { defineComponent } from 'vue';
import CharacterSelectorCard from './CharacterSelectorCard.vue';
import { XCircleIcon } from '@heroicons/vue/20/solid';
import { GetCharactersRequest } from '@arpanet/gen/services/auth/auth_pb';
import { User } from '@arpanet/gen/resources/users/users_pb';
import { getAccountClient, handleGRPCError } from '../../grpc';
import { RpcError } from 'grpc-web';
import { mapActions } from 'vuex';

export default defineComponent({
    components: {
        CharacterSelectorCard,
        XCircleIcon,
    },
    data() {
        return {
            'chars': [] as Array<User>,
        };
    },
    methods: {
        ...mapActions([
            'updateActiveChar',
        ]),
        async fetchCharacters() {
            return getAccountClient().
                getCharacters(new GetCharactersRequest(), null).
                then((resp) => {
                    this.chars = resp.getCharsList();
                }).catch((err: RpcError) => {
                    handleGRPCError(err, this.$route);
                });
        },
    },
    beforeMount() {
        this.updateActiveChar(null);

        // Fetch user's characters
        this.fetchCharacters();
    },
});

</script>

<template>
    <div v-if="chars.length <= 0" class="rounded-md bg-red-50 p-4">
        <div class="flex">
            <div class="flex-shrink-0">
                <XCircleIcon class="h-5 w-5 text-red-400" aria-hidden="true" />
            </div>
            <div class="ml-3">
                <h3 class="text-sm font-medium text-red-800">Unable to load your characters!</h3>
                <div class="mt-2 text-sm text-red-700">
                    <p>Please try again a few minutes.</p>
                </div>
            </div>
        </div>
    </div>
    <ul v-else role="list" class="grid grid-cols-1 gap-6 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4">
        <CharacterSelectorCard v-for="char in chars" :char="char" :key="char.getUserid()" />
    </ul>
</template>
