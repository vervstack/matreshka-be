<script setup lang="ts">
import SelectButton from "primevue/selectbutton";
import { ConfigValue } from "@/models/shared/Values.ts";

const model = defineModel<ConfigValue<string>>(
  {
    required: true,
  },
);
const props = defineProps({
  options: {
    type: Array<string>,
    default: [],
  },
  label: {
    type: String,
    default: "",
  },
});
</script>

<template>
  <div class="SelectValueContainer">
    <SelectButton v-model="model.value" :options="props.options" />
    <div class="floating-label">{{ label }}</div>
  </div>
</template>

<style>
.p-togglebutton {
  font-size: 1em;
}

.SelectValueContainer {
  display: flex;
  flex-direction: column;
  border: 1px solid var(--basic-element-color);
  border-radius: var(--border-radius);
  padding: 0.25em;
  position: relative;
}

.floating-label {
  position: absolute;
  top: -0.75em;
  left: 1em;

  padding: 0 0.5em;
  font-size: 0.75em;

  transform: translateY(50%);
  pointer-events: none;
}
</style>
