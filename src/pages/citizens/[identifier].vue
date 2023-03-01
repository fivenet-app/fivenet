<script lang="ts">
import { defineComponent } from 'vue';
import { useRoute } from 'vue-router/auto';
import { ref, watch } from 'vue';

export default defineComponent({
  setup() {
    const route = useRoute();
    const userData = ref();

    // fetch the user information when params change
    watch(() => route.params.identifier,
      async newIdentifier => {
        userData.value = await fetchUser(newIdentifier);
      });
  },
});
</script>

<route lang="json">
{
  "name": "citizens-byid",
  "meta": {
    "requiresAuth": false
  }
}
</route>

<template>
  TODO
</template>
