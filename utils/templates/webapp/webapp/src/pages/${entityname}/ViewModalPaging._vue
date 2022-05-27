<template>
  <MirzaModal id="modalPaging" ref="modalPaging" title="Paging Config" @submit="submitDataPaging">
    <label class="form-label">Item Size Per Page</label>
    <input type="number" class="form-control" placeholder="Page Size" v-model="state.filter.size">
  </MirzaModal>
</template>

<script setup>

import MirzaModal from "../../components/modal/MirzaModal.vue";
import {state} from "./state.js";
import {ref} from "vue";

const modalPaging = ref()

const emit = defineEmits(["submit"])

const submitDataPaging = () => {
  state.filter.page = 1
  emit("submit")
  hideModal()
}

const showModal = () => {
  modalPaging.value.showModal()
}

const hideModal = () => {
  modalPaging.value.hideModal()
}

defineExpose({showModal, hideModal})

</script>

<style scoped>

</style>