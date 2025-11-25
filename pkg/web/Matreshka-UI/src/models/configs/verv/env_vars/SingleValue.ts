import { Component } from "vue";
import EnvVariableView from "@/components/config/verv/env/EnvVariableView.vue";
import { EnvVar } from "@/models/configs/verv/env_vars/EnvVar.ts";
import { ConfigValue } from "@/models/shared/Values.ts";

export class StringEnvVarClass extends ConfigValue<string> implements EnvVar {
  constructor(envName: string, value: string) {
    super(envName, value);
  }

  getComponent(): Component {
    return EnvVariableView;
  }
}

export class IntEnvVarClass extends ConfigValue<number> implements EnvVar {
  constructor(envName: string, value: number) {
    super(envName, value);
  }

  getComponent(): Component {
    return EnvVariableView;
  }
}

export class FloatEnvVarClass extends ConfigValue<number> implements EnvVar {
  constructor(envName: string, value: number) {
    super(envName, value);
  }

  getComponent(): Component {
    return EnvVariableView;
  }
}

export class BooleanEnvVarClass extends ConfigValue<boolean> implements EnvVar {
  constructor(envName: string, value: boolean) {
    super(envName, value);
  }

  getComponent(): Component {
    return EnvVariableView;
  }
}
