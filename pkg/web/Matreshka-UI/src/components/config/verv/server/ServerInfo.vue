<script setup lang="ts">
import FileServerHandlerInfo from "@/components/config/verv/server/FileServerHandlerInfo.vue";
import GrpcHandlerInfo from "@/components/config/verv/server/GrpcHandlerInfo.vue";
import { ServerClass } from "@/models/configs/verv/servers/Servers.ts";

const cfg = defineModel<ServerClass>({ required: true });
</script>

<template>
  <div v-if="!cfg">No information on server</div>

  <div v-else>
    <div v-for="(_, idx) in cfg.grpc" class="Node">
      <div class="NodeField">
        {{ "Grpc" }}
      </div>
      <GrpcHandlerInfo v-model="cfg.grpc[idx]" />
    </div>

    <div v-for="(_, idx) in cfg.fs" class="Node">
      <div class="NodeField">
        {{ "File server" }}
      </div>
      <FileServerHandlerInfo v-model="cfg.fs[idx]" />
    </div>
  </div>
</template>

<style scoped>
@import "@/assets/styles/config_display.css";
</style>
