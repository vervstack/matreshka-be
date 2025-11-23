<script setup lang="ts">
import Dialog from "primevue/dialog";
import SpeedDial from "primevue/speeddial";
import { Component as VueComponent, ref, shallowRef, watch } from "vue";

import ConfigConstructorWidget from "@/widget/ConfigConstructorWidget.vue";
import Button from "primevue/button";

const isDialogOpen = ref<boolean>(false);
const newConfigDialog = shallowRef<VueComponent | undefined>();

watch(isDialogOpen, () => {
  if (!isDialogOpen.value) {
    newConfigDialog.value = undefined;
  }
});

watch(newConfigDialog, () => {
  if (newConfigDialog.value !== undefined) {
    isDialogOpen.value = true;
  }
});

const buttons = [
  {
    label: "New config",
    icon: "pi pi-box",
    command() {
      newConfigDialog.value = ConfigConstructorWidget;
    },
  },
];
</script>

<template>
  <!-- Help button at the bottom -->
  <SpeedDial
    :style="{ position: 'fixed', bottom: '2%', right: '2%' }"
    :tooltipOptions="{ event: 'hover', position: 'left' }"
    :model="buttons"
    direction="up"
    :radius="100"
  />

  <Dialog
    v-model:visible="isDialogOpen"
    modal
    header="New config"
    :dismissableMask="true"
    :pt="{
      root: 'border-none',
      mask: {
        style: 'backdrop-filter: blur(2px)',
      },
      pDialogContent: {
        overflow: '',
      },
    }"
    :style="{
      width: 'clamp(40em, 40vw, 60vh)',
      height: 'clamp(50vh, 60vh, 80vh)',
    }"
    :contentStyle="{ height: '100%', width: '100%' }"
    position="right"
  >
    <Component :is="newConfigDialog" />
  </Dialog>
</template>

<style scoped></style>
