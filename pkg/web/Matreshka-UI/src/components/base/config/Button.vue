<script setup lang="ts">

defineProps({
  label: {
    type: String,
    default: "",
  },
  icon: {
    type: String,
  },
  disabled: {
    type: Boolean,
    default: false,
  },
  rounded: {
    type: Boolean,
    default: false,
  },
  borderless: {
    type: Boolean,
    default: false,
  },
  soft: {
    type: Boolean,
    default: true,
  },
});

const emit = defineEmits(["click"]);

function buttonClicked() {
  emit("click");
}
</script>

<template>
  <button
    class="simple-button"
    :class="{
    'rounded':rounded,
    'borderless': borderless,
    'soft': soft,
    }"
    :onclick="buttonClicked"
  >
    <img
      class="icon"
      v-if="icon"
      :src="icon"
      alt="no icon"
    />
    <p v-else>
      {{ label }}
    </p>
  </button>
</template>

<style scoped>
.simple-button {
  width: 100%;
  height: 100%;

  background: none;
  box-sizing: border-box;
  border: black solid 1px;
  display: flex;
  justify-content: center;
  align-items: center;

  font-size: inherit;
}

.simple-button:hover {
  cursor: pointer;
}

.rounded {
  border-radius: 50%;
}

.soft {
  border-radius: var(--border-radius);
}

.borderless {
  border: none;
}

.icon {
  width: 100%;
  height: 100%;
}
</style>
