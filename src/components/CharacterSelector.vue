<script lang="ts">
import { defineComponent } from 'vue';
import { mapState } from 'vuex';
import CharacterCard from './CharacterCard.vue';

export default defineComponent({
    components: {
        CharacterCard,
    },
    computed: {
        ...mapState([
            'accessToken',
            'chars',
        ]),
    },
});
</script>

<template>
    <div class="grid place-items-center">
        <div class="flex w-full">
            <CharacterCard v-for="(char, index) in chars" :char="char" :charIndex="index" :key="char.getIdentifier()" />
        </div>
    </div>
</template>
