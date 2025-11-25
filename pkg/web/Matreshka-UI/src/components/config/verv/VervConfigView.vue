<script setup lang="ts">
import VervConfig from "@/models/configs/verv/VervConfig.ts";

import AppInfo from "@/components/config/verv/app_info/AppInfo.vue";
import ResourcesInfo from "@/components/config/verv/resource/AppResources.vue";
import ServersInfo from "@/components/config/verv/server/ServersInfo.vue";
import VervEnvVars from "@/components/config/verv/env/VervEnvVars.vue";

const model = defineModel<VervConfig>({
  required: true,
});
</script>

<template>
  <div class="Content">
    <div
      class="ContentBlock"
      :style="{
        borderColor: model.appInfo.isChanged() ? 'var(--warn)' : 'var(--good)',
      }"
    >
      <AppInfo v-model="model.appInfo" />
    </div>
    <div
      class="ContentBlock"
      :style="{
        borderColor:
          model?.getChangedDataSourcesNames().length != 0 ? 'var(--warn)' : 'var(--good)',
      }"
    >
      <ResourcesInfo v-model="model.dataSources" />
    </div>
    <div
      class="ContentBlock"
      :style="{
        borderColor:
        model?.getChangedServersNames().length != 0 ? 'var(--warn)' : 'var(--good)',
      }"
    >
      <ServersInfo v-model="model.servers" />
    </div>

    <div
      class="ContentBlock"
      :style="{
        borderColor:
        model?.getChangedEnvVarNames().length != 0 ? 'var(--warn)' : 'var(--good)',
      }"
    >
      <VervEnvVars v-model="model.envVars" />
    </div>
  </div>
</template>

<style scoped>
.Content {
  display: flex;
  justify-content: center;
  flex-direction: column;
  gap: 1em;
}

.ContentBlock {
  border: solid;
  border-radius: 16px;
  padding: 0.75em;
}
</style>
