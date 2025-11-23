<script setup lang="ts">
import { EnvVar } from "@/models/configs/verv/env_vars/EnvVar.ts";
import { ConfigValue } from "@/models/shared/Values.ts";
import { StringArrayEnvVarClass } from "@/models/configs/verv/env_vars/MultipleValues.ts";

const envVars = defineModel<EnvVar[]>({ default: [] });

// Values that enters via input
const typingValues: EnvVar[] = envVars.value
  .filter(v => (
      v instanceof ConfigValue && ["string", "number"].includes(typeof v.value))
    || v instanceof StringArrayEnvVarClass);

const toggleValues: EnvVar[] = envVars.value
  .filter(v => v instanceof ConfigValue && typeof v.value == "boolean");

</script>

<template>
  <div class="Node">
    <div v-if="envVars.length == 0">No environment variables are defined</div>

    <div v-else>Environment variables:</div>
    <div class="Node" v-for="(ev, i) in typingValues" :key="ev.getOriginalName()">
      <component
        :is="ev.getComponent()"
        v-model="typingValues[i]"
      />
    </div>
    <div class="Node" v-for="(ev, i) in toggleValues" :key="ev.getOriginalName()">
      <component
        :is="ev.getComponent()"
        v-model="toggleValues[i]"
      />
    </div>
  </div>
</template>

<style scoped>
@import "@/assets/styles/config_display.css";

</style>
