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


function labeler(val: string): string {
  // Displays '*' over original value in case other value was selected
  if (val != model.value.value &&
        model.value.getOriginalValue() === val) {
      return val+'*'
  }

  return val;
}

</script>

<template>
  <div
    class="SelectValueContainer"
    :class="{'isChanged': model.isChanged()}"
  >
    <SelectButton
      v-model="model.value"
      :options="props.options"
      :optionLabel="labeler"
    />
    <div class="FloatingLabel">{{ label }}</div>
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
  padding: 0.5em;
  position: relative;

  &.isChanged {
    border-color: var(--value-changed-outline);
  }
}

.FloatingLabel {
  position: absolute;
  top: -1em;
  left: 1em;

  padding: 0 0.5em;
  font-size: 0.75em;

  transform: translateY(50%);
  pointer-events: none;

  background: var(--background);
  color: var(--text-color-secondary);
}
</style>
