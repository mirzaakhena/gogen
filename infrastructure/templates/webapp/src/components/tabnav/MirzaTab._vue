<template>
  <ul class="nav nav-tabs justify-content-start mb-3">
    <MirzaTabLink v-for="header in headers" :key="header.to" :to="header.to" :icon="header.icon" :name="header.name"></MirzaTabLink>
  </ul>
  <div class="tab-content">
    <div class="tab-pane fade show active">
      <router-view />
    </div>
  </div>
</template>

<script setup>

import MirzaTabLink from "./MirzaTabLink.vue";

defineProps({
  headers: {type: Array, required: true}
})

</script>

<style scoped>

</style>