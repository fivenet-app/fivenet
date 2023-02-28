<script lang="ts">
import { defineComponent } from 'vue';

import Navbar from '../../components/Navbar.vue';
import Footer from '../../components/Footer.vue';

export default defineComponent({
  components: {
    Navbar,
    Footer,
  },
});
</script>

<route lang="json">
{
  "name": "documents",
  "meta": {
    "requiresAuth": true
  }
}
</route>

<template>
  <Navbar />
  <div class="container mx-auto px-4">
    <div class="text-sm breadcrumbs">
      <ul>
        <li>
          <a>
            <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" class="w-4 h-4 mr-2 stroke-current">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                d="M3 7v10a2 2 0 002 2h14a2 2 0 002-2V9a2 2 0 00-2-2h-6l-2-2H5a2 2 0 00-2 2z"></path>
            </svg>
            Overview
          </a>
        </li>
        <li>
          <a>
            <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" class="w-4 h-4 mr-2 stroke-current">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                d="M3 7v10a2 2 0 002 2h14a2 2 0 002-2V9a2 2 0 00-2-2h-6l-2-2H5a2 2 0 00-2 2z"></path>
            </svg>
            Documents
          </a>
        </li>
      </ul>
    </div>
    <h1>DOCUMENTS</h1>
  </div>
  <Footer />
</template>
