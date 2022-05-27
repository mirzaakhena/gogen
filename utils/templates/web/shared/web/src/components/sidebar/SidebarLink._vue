<template>

  <li>
    <router-link :to="to" class="nav-link" :class="{ active: isActive }" aria-current="page">
      <svg class="bi me-2" width="16" height="16">
        <use :href="icon"/>
      </svg>
      <transition name="fade">
        <span v-if="!collapsed">
          <slot/>
        </span>
      </transition>
    </router-link>
  </li>

</template>

<script setup>
import {computed} from "vue";
import {useRoute} from 'vue-router'

const props = defineProps({
  to: {type: String, required: true},
  icon: {type: String, required: true},
})

const isActive = computed(() => useRoute().path.startsWith(props.to))

const collapsed = false
</script>

<style scoped>

.nav-link {
  color: black;
}

.link {
  display: flex;
  align-items: center;

  cursor: pointer;
  position: relative;
  font-weight: 400;
  user-select: none;

  margin: 0.1em 0;
  padding: 0.4em;
  border-radius: 0.25em;
  height: 1.5em;

  /*color: white;*/
  text-decoration: none;
}

.link:hover {
  background-color: var(--sidebar-item-hover);
}

.link.active {
  background-color: var(--sidebar-item-active);
}

.link .icon {
  flex-shrink: 0;
  width: 25px;
  margin-right: 10px;
}
</style>