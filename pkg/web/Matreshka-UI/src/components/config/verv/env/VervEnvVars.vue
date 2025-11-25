<script setup lang="ts">
import { EnvVar } from "@/models/configs/verv/env_vars/EnvVar.ts";
import { ConfigValue } from "@/models/shared/Values.ts";
import { StringArrayEnvVarClass } from "@/models/configs/verv/env_vars/MultipleValues.ts";

import SelectButton from "@/components/base/config/SelectButton.vue";
import { Ref, ref, watch } from "vue";

const envVars = defineModel<EnvVar[]>({ default: [] });

function isSimpleValue(v: EnvVar): boolean {
  return (v instanceof ConfigValue && ["string", "number"].includes(typeof v.value))
    || v instanceof StringArrayEnvVarClass;
}

function isToggle(v: EnvVar): boolean {
  return (v instanceof ConfigValue && typeof v.value == "boolean");
}

const groupsOfValues: Ref<number[][]> = ref(groupByType());

enum groupingType {
  byType = "By type",
  byName = "By Name"
}

const groupingOptions: string[] = [groupingType.byType, groupingType.byName];
let selectedGrouping = ref(groupingOptions[0]);

watch(selectedGrouping, (newV) => {
  switch (newV) {
    case groupingType.byType:
      groupsOfValues.value = groupByType();
      break;
    case groupingType.byName:
      groupsOfValues.value = groupByName();
      break;
  }
});

function groupByType(): number[][] {
  return [
    envVars.value
      .map((v, idx) => isSimpleValue(v) ? idx : undefined)
      .filter(v => v != undefined),
    envVars.value
      .map((v, idx) => isToggle(v) ? idx : undefined)
      .filter(v => v != undefined),
  ];
}

function groupByName(): number[][] {
  const groups = new Map<string, number[]>();

  envVars.value.forEach((v, idx) => {
    const name = v.getOriginalName();
    const prefix = name.split("-")[0];

    if (!groups.has(prefix)) {
      groups.set(prefix, []);
    }
    groups.get(prefix)?.push(idx);
  });

  const otherGroup: number[] = [];
  const mainGroups: number[][] = [];

  groups.forEach((indices) => {
    if (indices.length === 1) {
      otherGroup.push(...indices);
    } else {
      mainGroups.push(indices.sort((a, b) =>
        envVars.value[a].getOriginalName()
          .localeCompare(envVars.value[b].getOriginalName()),
      ));
    }
  });

  return otherGroup.length ? [...mainGroups, otherGroup.sort((a, b) =>
    envVars.value[a].getOriginalName()
      .localeCompare(envVars.value[b].getOriginalName()),
  )] : mainGroups;
}


</script>

<template>
  <div class="Node">
    <div v-if="envVars.length == 0">No environment variables are defined</div>

    <div
      v-else
      class="NodeField horizontal blockHeader"
    >
      <div>Environment variables:</div>
      <div
        class="groupingWrapper"
      >
        <div> Group by</div>
        <SelectButton
          v-model="selectedGrouping"
          :options="groupingOptions"
        />
      </div>
    </div>

    <div
      class="Node outlined"
      v-for="(group, j) in groupsOfValues"
      :key="j">
      <div
        class="NodeField"
        v-for="(i) in group"
        :key="envVars[i].getOriginalName()">
        <component
          :is="envVars[i].getComponent()"
          v-model="envVars[i]"
        />
      </div>
    </div>
  </div>
</template>

<style scoped>
@import "@/assets/styles/config_display.css";

.blockHeader {
  justify-content: space-between;
}

.groupingWrapper {
  display: flex;
  flex-direction: row;
  align-items: center;
  gap: 1em;
}
</style>
