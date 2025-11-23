<script setup lang="ts">
import { ref } from "vue";

export interface IInputer {
  doFocus: () => void;
}

defineProps({
  placeholder: {
    type: String,
    default: "",
  },
  floatingLabel: {
    type: String,
    default: "",
  },
  disabled: {
    type: Boolean,
    default: false
  }
});

const model = defineModel<string>({
  required: true,
});

const inputElementRef = ref<HTMLElement>();

function doFocus() {
  inputElementRef.value?.focus();
}

defineExpose<IInputer>({ doFocus });
</script>

<template>
    <input
      :disabled="disabled"
      ref="inputElementRef"
      v-model="model"
      type="text"
      class="input"
      :placeholder="placeholder"
    />
</template>

<style scoped>

.input {
  width: 100%;
  height: 100%;
  padding: 0.5rem 0.75rem;
  font-size: 0.75em;
  border: 1px solid #d1d5db;
  border-radius: var(--border-radius);
  outline: none;
}

.input:focus {
  border-color: var(--focus);
  box-shadow: 1px var(--focus);
}

</style>
