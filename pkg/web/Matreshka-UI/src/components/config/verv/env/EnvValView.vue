<script setup lang="ts">

import { EnvVar } from "@/models/configs/verv/env_vars/EnvVar.ts";
import { ConfigValue } from "@/models/shared/Values.ts";
import { StringArrayEnvVarClass } from "@/models/configs/verv/env_vars/MultipleValues.ts";

import ConfigField from "@/components/base/config/fields/ConfigInput.vue";
import Toggle from "@/components/base/config/fields/Toggle.vue";
import SelectValue from "@/components/base/config/fields/SelectValue.vue";

const model = defineModel<EnvVar>({
  required: true,
});

const envVarPrefix = 'ENVIRONMENT_'

</script>

<template>
  <div class="NodeField">
    <ConfigField
      v-if="model instanceof ConfigValue && ['string', 'number'].includes(typeof model.value) && model.enums.length == 0"
      v-model="model"
      :field-name="model.envName.slice(envVarPrefix.length)"
    />
    <SelectValue
      v-else-if="model instanceof ConfigValue && ['string'].includes(typeof model.value) && model.enums.length != 0"
      v-model="model"
      :options="model.enums"
      :label="model.envName.slice(envVarPrefix.length)"
      />
    <Toggle
      v-else-if="model instanceof ConfigValue && typeof model.value == 'boolean'"
      v-model="model"
      :field-name="model.envName.slice(envVarPrefix.length)" />

    <div
      v-else-if="model instanceof StringArrayEnvVarClass"
      class="Node outlined"
      :class="{'changed':model.isChanged() }"
    >
      <div class="NodeField Header horizontal">
        <div>{{ model.rootName.slice(envVarPrefix.length) }}</div>
        <img
          v-tooltip.bottom="'Array of values'"
          class="HintIcon"
          src="@/assets/svg/general/hint.svg"
          alt="Add" />
      </div>
      <div
        v-for="(_, index) in model.values"
        :key="index"
        class="Node"
        style="display:flex; flex-direction: row;"
      >
        <ConfigField
          v-model="model.values[index]"
          field-name=" "
        />
      </div>
    </div>
  </div>
</template>

<style scoped>
@import "@/assets/styles/config_display.css";

.HintIcon {
  width: 1em;
  aspect-ratio: 1;
}
</style>
