<script setup lang="ts">
import { ref } from "vue";

import { ConfigValue } from "@/models/shared/Values.ts";

import Inputer from "@/components/base/Inputer.vue";
import Button from "@/components/base/config/Button.vue";
import RollbackIcon from "@/assets/svg/general/rollback.svg";

const model = defineModel<ConfigValue<string>>({
  required: true,
});

const emit = defineEmits(["rollback", "showOriginal", "showActual"]);

const originalNameRef = ref<string>(model.value.getOriginalName());
const originalValueRef = ref<string>(model.value.getOriginalValue());

const isPreparingToRevert = ref<boolean>(false);

function rollback() {
  emit("rollback");
  isPreparingToRevert.value = false;
}

function showOriginal() {
  emit("showOriginal");
  isPreparingToRevert.value = true;
}

function showActual() {
  emit("showActual");
  isPreparingToRevert.value = false;
}

function shouldShowOldName(): boolean {
  return isPreparingToRevert.value && model.value.isNameChanged();
}

function isNew(): boolean {
  return model.value.isNew && !model.value.isMuted;
}

function shouldShowValue(): boolean {
  return model.value.value !== "";
}

function shouldShowOldValue(): boolean {
  return model.value.isValueChanged() && isPreparingToRevert.value;
}

</script>

<template>
  <div
    class="KeyValueInputer"
  >
    <div class="Field show"
         :class="{
            changed: model.isNameChanged(),
            new: isNew(),
          }"
    >
      <Inputer
        v-model="model.envName"
      />
    </div>
    <div
      class="Field"
      :class="{show: shouldShowOldName()}"
    >
      <Inputer
        disabled
        v-model="originalNameRef"
      />
    </div>
    <p
      v-if="model.value !== '' || model.getOriginalValue() != '' && isPreparingToRevert"
      class="Colon"
    >:</p>

    <!--    For value when this is a leaf-->
    <div
      class="Field"
      :class="{
        show: shouldShowValue(),
        changed: model.isValueChanged(),
        new: isNew(),
      }"
    >
      <Inputer
        v-model="model.value"
      />
    </div>
    <!--    For name when this is a node-->
    <div
      class="Field"
      :class="{show: shouldShowOldValue()}"
    >
      <Inputer
        disabled
        v-model="originalValueRef"
      />
    </div>


    <Transition name="RollbackButton" mode="in-out">
      <Button
        v-if="model.isNameChanged() || model.isValueChanged() || isPreparingToRevert"
        class="RollBackButton"
        title="Rollback to original"
        @click="rollback"
        :icon="RollbackIcon"

        @mouseenter="showOriginal"
        @mouseleave="showActual"
      />
    </Transition>
  </div>
</template>

<style scoped>
.KeyValueInputer {
  width: 100%;
  height: 100%;

  display: flex;
  justify-content: center;
  align-items: center;
  box-sizing: border-box;
}

.Field {
  flex: 0;
  overflow: hidden;
  transition: 0.5s ease-in-out;
  border: none;
}

.Colon {
  margin: 0 0.25em 0 0.25em;
}

.show {
  flex: 1;
}

.RollBackButton {
  width: 1.75em;
  height: 1.75em;
  margin-left: 0.25em;

  overflow: hidden;
  border: black solid 1px;
  border-radius: var(--border-radius);
}

.show.changed {
  border-right: solid 4px;
  border-color: var(--warn);
  border-radius: var(--border-radius);
}

.show.new {
  border-right: solid 4px;
  border-color: var(--good);
  border-radius: var(--border-radius);
}

.RollbackButton-enter-active,
.RollbackButton-leave-active {
  transition: all 0.3s ease;
}

.RollbackButton-enter-to,
.RollbackButton-leave-from {
  transform: translateX(0);
  opacity: 1;
}

.RollbackButton-enter-from,
.RollbackButton-leave-to {
  transform: translateX(-10%);
  opacity: 0;
}
</style>
