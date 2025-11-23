<script setup lang="ts">

import Button from "@/components/base/config/Button.vue";

const model = defineModel<
  {
    isOpen: Boolean,
    options: {
      title?: string,
      icon?: string,
      value: string
    }[]
  }>({
  default: {
    isOpen: false,
    options: [],
  },
});

defineProps({
  title: {
    type: String,
    default: "",
  },
  icon: {
    type: String,
  },
  borderless: {
    type: Boolean,
    default: false,
  },
  disabled: {
    type: Boolean,
    default: false,
  },
  openOnHover: {
    type: Boolean,
    default: false,
  },
});

const emit = defineEmits(["select"]);

function select(opt: string) {
  emit("select", opt);
}

</script>

<template>
  <div
    class="DialButtonWrapper">
    <Button
      :title="title"
      :icon="icon"
      :borderless="borderless"
      :disabled="disabled"
      rounded
      @click="model.isOpen = !model.isOpen"

    />

    <TransitionGroup
      name="subButton">
      <div
        v-for="(opt, idx) in model.options"
        class="SubButton"
        v-if="model.isOpen"
        :style="{
          top: 0,
          right: 2+idx*2.2+'em',
        }"
      >
        <Button
          :title="opt.title"
          :icon="opt.icon"
          @click="select(opt.value)"
        />
      </div>
    </TransitionGroup>
  </div>
</template>

<style scoped>

.DialButtonWrapper {
  position: relative;
}

.SubButton {
  position: absolute;

  width: 2em;
  height: 2em;
  transition: all 0.25s ease-in-out;
}

.SubButton:hover {
  transform: scale(1.05);
}

.subButton-enter-active,
.subButton-leave-active {
  transition: all 0.3s ease;
}

.subButton-enter-to,
.subButton-leave-from {
  transform: translateX(0);
  opacity: 1;
}

.subButton-enter-from,
.subButton-leave-to {
  transform: translateX(10%);
  opacity: 0;
}
</style>
