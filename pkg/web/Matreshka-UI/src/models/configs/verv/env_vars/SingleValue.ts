import { Component } from "vue";
import EnvValView from "@/components/config/verv/env/EnvValView.vue";
import { EnvVar } from "@/models/configs/verv/env_vars/EnvVar.ts";
import { ConfigValue } from "@/models/shared/Values.ts";

export class StringEnvVarClass extends ConfigValue<string> implements EnvVar {
  constructor(envName: string, value: string) {
    super(envName, value);
  }

  getComponent(): Component {
    return EnvValView;
  }
}

export class IntEnvVarClass extends ConfigValue<number> implements EnvVar {
  constructor(envName: string, value: number) {
    super(envName, value);
  }

  getComponent(): Component {
    return EnvValView;
  }
}

export class FloatEnvVarClass extends ConfigValue<number> implements EnvVar {
  constructor(envName: string, value: number) {
    super(envName, value);
  }

  getComponent(): Component {
    return EnvValView;
  }
}

export class BooleanEnvVarClass extends ConfigValue<boolean> implements EnvVar {
  constructor(envName: string, value: boolean) {
    super(envName, value);
  }

  getComponent(): Component {
    return EnvValView;
  }
}
