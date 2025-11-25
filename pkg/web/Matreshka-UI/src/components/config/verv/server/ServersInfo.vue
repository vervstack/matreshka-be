<script setup lang="ts">
import ConfigField from "@/components/base/config/fields/ConfigInput.vue";
import ServerInfo from "@/components/config/verv/server/ServerInfo.vue";
import Server from "@/models/configs/verv/servers/Servers.ts";

const servers = defineModel<Server[]>({ default: [] });
</script>

<template>
  <div class="Node">
    <div v-if="servers.length == 0">No servers defined</div>
    <div v-else>Servers:</div>
    <div class="Node" v-for="(s, i) in servers" :key="s.name">
      <div class="NodeField">
        {{ s.name }}
      </div>
      <div class="NodeField PortField">
        <ConfigField v-model="s.port" field-name="Port" />
      </div>
      <div class="NodeField">
        <ServerInfo v-model="servers[i]" />
      </div>
    </div>
  </div>
</template>

<style scoped>
@import "@/assets/styles/config_display.css";

.PortField {
  flex-direction: row;
  align-items: center;
}
</style>
