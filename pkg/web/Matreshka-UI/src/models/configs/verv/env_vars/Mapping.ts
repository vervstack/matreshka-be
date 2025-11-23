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

  const envVars: EnvVar[] = [...root.innerNodes.map(() => new StringEnvVarClass("", ""))];

  root.innerNodes.map((node: Node) => {
    if (!node.name || !node.innerNodes) {
      return;
    }

    const parts = node.name.split("_");
    const orderIndex = Number(parts[1].slice(1, parts[1].length - 1));

    let name = "";

    let singleValue = "";
    let multipleValues: Node[] = [];
    let valuePrefix = node.name;

    let dataType = "";

    node.innerNodes.map((subNod: Node) => {
      if (!node.name || !subNod.name) return;

      switch (subNod.name.substring(node.name.length + 1)) {
        case "NAME":
          name = subNod.value || "";
          valuePrefix = subNod.name;
          break;
        case "VALUE":
          if (subNod.value) {
            singleValue = subNod.value;
          }

          if (subNod.innerNodes?.length != 0) {
            multipleValues = subNod.innerNodes || [];
          }
          break;
        case "TYPE":
          dataType = subNod.value || "";
          break;
        default:
          // console.log("Unknown env var field: " + subNod.name);
      }

    });


    envVars[orderIndex] = multipleValues.length > 0 ?
      mapMultipleValues(name, valuePrefix, multipleValues) :
      mapSingleValueByDataType(name, singleValue, dataType);
  });
  console.log(envVars);
  return envVars;
}


function mapMultipleValues(name: string, prefix: string, nodes: Node[]): EnvVar {

  let valsArray: string[] = [...nodes.map(() => "")];

  nodes.map((n: Node) => {
    if (!n.name || !n.value) return;
    // Extracting index from name
    // 1 for last symbol
    // 1 for underscore
    // 1 for index opening bracket "["
    const orderIndex = Number(n.name.slice(prefix.length + 3, n.name.length - 1));

    // As for now, string values and int arrays ar basically same due to range syntax (1:3 - 1,2,3)
    // in future different approach to handle range values and slices might appear
    valsArray[orderIndex] = n.value;
  });

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
