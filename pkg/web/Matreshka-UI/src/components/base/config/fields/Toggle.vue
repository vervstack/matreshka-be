<script setup lang="ts">
import ToggleSwitch from "primevue/toggleswitch";
import { ConfigValue } from "@/models/shared/Values.ts";

const model = defineModel<ConfigValue<boolean>>({ required: true });
defineProps({
  fieldName: {
    type: String,
  },
})

</script>

<template>
  <div
    class="ConfigToggle"
    :class="{ 'changed': model.isChanged() }"
  >
    <div>{{ fieldName || model.envName }}:</div>
    <ToggleSwitch v-model="model.value">
      <template #handle="{ checked }">
        <i :class="['!text-xs pi', { 'pi-check': checked, 'pi-times': !checked }]" />
      </template>
    </ToggleSwitch>
  </div>
</template>

<style scoped>
.ConfigToggle {
  display: flex;
  flex-direction: row;
  width: 100%;
  gap: 1em;
  align-items: center;

  border: 1px solid var(--basic-element-color);
  padding: 0.5em;
  border-radius: var(--border-radius);

  &.changed {
    border-color: var(--value-changed-outline);
  }
}
</style>
