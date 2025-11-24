import {
  DataType,
  EnvVar,
} from "@/models/configs/verv/env_vars/EnvVar.ts";
import { Node } from "@vervstack/matreshka";
import {
  BooleanEnvVarClass,
  FloatEnvVarClass,
  IntEnvVarClass,
  StringEnvVarClass,
} from "@/models/configs/verv/env_vars/SingleValue.ts";
import { StringArrayEnvVarClass } from "@/models/configs/verv/env_vars/MultipleValues.ts";

export default function mapEnvVar(root: Node): EnvVar[] {
  if (!root.innerNodes) {
    return [];
  }

  const envVars: EnvVar[] = [];

  root.innerNodes.map((node: Node) => {
    if (!node.name || !node.innerNodes) {
      return;
    }

    let name = node.name.split("_")[1];

    let singleValue = "";
    let multipleValues: string[] = [];

    let dataType = "";

    if (node.value?.startsWith("[") && node.value?.endsWith("]")) {
      multipleValues = node.value?.slice(1, -1).split(",");
    } else {
      singleValue = node.value || "123";
    }

    node.innerNodes.map((subNod: Node) => {
      if (!node.name || !subNod.name) return;

      switch (subNod.name.substring(node.name.length + 1)) {
        case "TYPE":
          dataType = subNod.value || "";
          break;
        default:
        // console.log("Unknown env var field: " + subNod.name);
      }

    });


    envVars.push(multipleValues.length > 0 ?
      mapMultipleValues(name, multipleValues) :
      mapSingleValueByDataType(name, singleValue, dataType));
  });

  return envVars;
}


function mapMultipleValues(name: string, valsArray: string[]): EnvVar {
  return new StringArrayEnvVarClass(name, valsArray);
}

function mapSingleValueByDataType(name: string, v: string, dt: string): EnvVar {
  switch (dt) {
    case DataType.INT:
      return new IntEnvVarClass(name, Number(v));
    case DataType.FLOAT:
      return new FloatEnvVarClass(name, Number(v));
    case DataType.BOOLEAN:
      return new BooleanEnvVarClass(name, v.toLowerCase() == "true");
    default:
      return new StringEnvVarClass(name, v);
  }
}
