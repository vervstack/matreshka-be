<script setup lang="ts">
import { EnvVar } from "@/models/configs/verv/env_vars/EnvVar.ts";
import { ConfigValue } from "@/models/shared/Values.ts";
import { StringArrayEnvVarClass } from "@/models/configs/verv/env_vars/MultipleValues.ts";

const envVars = defineModel<EnvVar[]>({ default: [] });

// Values that enters via input
const typingValuesIndexes: number[] = envVars.value
  .map((v, idx) => {
    if ((
        v instanceof ConfigValue && ["string", "number"].includes(typeof v.value))
      || v instanceof StringArrayEnvVarClass) return idx;
  }).filter(v => v != undefined);

const toggleValuesIndexes: number[] = envVars.value
  .map((v, idx) => {
    if (v instanceof ConfigValue && typeof v.value == "boolean") return idx;
  }).filter(v => v != undefined);

</script>

<template>
  <div class="Node">
    <div v-if="envVars.length == 0">No environment variables are defined</div>

    <div v-else>Environment variables:</div>
    <div class="Node" v-for="(i) in typingValuesIndexes" :key="envVars[i].getOriginalName()">
      <component
        :is="envVars[i].getComponent()"
        v-model="envVars[i]"
      />
    </div>
    <div class="Node" v-for="(i) in toggleValuesIndexes" :key="envVars[i].getOriginalName()">
      <component
        :is="envVars[i].getComponent()"
        v-model="envVars[i]"
      />
    </div>
  </div>
</template>

<style scoped>
@import "@/assets/styles/config_display.css";

</style>
