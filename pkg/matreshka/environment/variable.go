package environment

import (
	"fmt"
	"reflect"
	"strings"

	"go.redsock.ru/evon"
	errors "go.redsock.ru/rerrors"
	"go.redsock.ru/toolbox"
	"gopkg.in/yaml.v3"
)

var (
	ErrUnknownEnvVariableType = errors.New("unknown environment variable type")
	ErrEnumValidationFailed   = errors.New("value doesn't satisfy enum")
	ErrNotAEnumable           = errors.New("not a enumable type. Only string | int values can be enum")
	ErrUnexpectedType         = errors.New("unexpected type")
)

type variableType string

const (
	VariableTypeInt      variableType = "int"
	VariableTypeStr      variableType = "string"
	VariableTypeBool     variableType = "bool"
	VariableTypeFloat    variableType = "float"
	VariableTypeDuration variableType = "duration"
)

type Variable struct {
	Name    string       `yaml:"name"`
	Type    variableType `yaml:"type"`
	Enum    TypedEnum    `yaml:"enum,omitempty"`
	Value   Value        `yaml:"value"`
	Comment string       `yaml:"comment,omitempty"`
}

type opt func(*Variable)

func MustNewVariable(name string, val any, opts ...opt) *Variable {
	v, err := NewVariable(name, val, opts...)
	if err != nil {
		panic(err)
	}

	return v
}

func NewVariable(name string, val any, opts ...opt) (*Variable, error) {
	out := &Variable{
		Name: name,
	}

	for _, o := range opts {
		o(out)
	}

	out.Type = toolbox.Coalesce(out.Type, GetType(val))
	if out.Type == "" {
		return nil, errors.Wrap(ErrUnknownEnvVariableType)
	}

	var err error
	out.Value.val, err = mapVariableTypeToTypedValueConstructor[out.Type](val)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to convert value to typed constructor for field %s", out.Name)
	}

	return out, nil
}

var (
	mapReflectTypeToVariableType = map[reflect.Kind]variableType{
		reflect.String:  VariableTypeStr,
		reflect.Bool:    VariableTypeBool,
		reflect.Float64: VariableTypeFloat,
		reflect.Float32: VariableTypeFloat,
		reflect.Int:     VariableTypeInt,
		reflect.Int8:    VariableTypeInt,
		reflect.Int16:   VariableTypeInt,
		reflect.Int32:   VariableTypeInt,
		reflect.Int64:   VariableTypeInt,
		reflect.Uint:    VariableTypeInt,
		reflect.Uint8:   VariableTypeInt,
		reflect.Uint16:  VariableTypeInt,
		reflect.Uint32:  VariableTypeInt,
		reflect.Uint64:  VariableTypeInt,
	}
	mapVariableTypeToTypedValueConstructor = map[variableType]func(in any) (typedValue, error){
		VariableTypeStr:      toStringValue,
		VariableTypeInt:      toIntVariable,
		VariableTypeFloat:    toFloatVariable,
		VariableTypeBool:     toBoolValue,
		VariableTypeDuration: toDuration,
	}
	mapVariableTypeToYamlNodeParser = map[variableType]func(node *yaml.Node) (typedValue, error){
		VariableTypeStr:      stringValueFromNode,
		VariableTypeInt:      intValueFromNode,
		VariableTypeFloat:    fromFloatNode,
		VariableTypeBool:     fromBoolNode,
		VariableTypeDuration: fromDurationNode,
	}
	mapVariableTypeToEnumYamlNodeParser = map[variableType]func(node *yaml.Node) (typedEnum, error){
		VariableTypeStr: func(node *yaml.Node) (typedEnum, error) {
			return stringSliceFromYamlNode(node)
		},
		VariableTypeInt: func(node *yaml.Node) (typedEnum, error) {
			return intSliceFromYamlNode(node)
		},
	}
	mapVariableTypeToEnumEvonNodeParser = map[variableType]func(node *evon.Node) (typedEnum, error){
		VariableTypeStr: func(node *evon.Node) (typedEnum, error) {
			return stringsSliceFromEvonNode(node)
		},
		VariableTypeInt: func(node *evon.Node) (typedEnum, error) {
			return intSliceFromEvonNode(node)
		},
	}
)

func (v *Variable) UnmarshalYAML(node *yaml.Node) error {
	var value, enum *yaml.Node

	for cIdx := 0; cIdx < len(node.Content); cIdx += 2 {
		fieldName := node.Content[cIdx].Value
		switch fieldName {
		case "name":
			v.Name = node.Content[cIdx+1].Value
		case "type":
			v.Type = variableType(node.Content[cIdx+1].Value)
		case "value":
			value = node.Content[cIdx+1]
		case "enum":
			enum = node.Content[cIdx+1]
		}
	}

	parser, exists := mapVariableTypeToYamlNodeParser[v.Type]
	if !exists {
		return ErrUnknownEnvVariableType
	}
	val, err := parser(value)
	if err != nil {
		return errors.Wrap(err, "error parsing yaml value node")
	}
	v.Value.val = val

	if enum != nil {
		constructor, ok := mapVariableTypeToEnumYamlNodeParser[v.Type]
		if !ok {
			return errors.Wrap(ErrNotAEnumable, "error unmarshalling yaml enum")
		}

		v.Enum.v, err = constructor(enum)
		if err != nil {
			return errors.Wrap(err, "error parsing yaml enums node")
		}

		err = v.Enum.v.isEnum(v.Value.val)
		if err != nil {
			return errors.Wrap(err)
		}
	}

	return nil
}

func (v *Variable) UnmarshalEnv(node *evon.Node) (err error) {
	v.Name = node.Name
	splitIdx := strings.LastIndex(v.Name, evon.ObjectSplitter)
	if splitIdx != -1 {
		v.Name = v.Name[splitIdx+1:]
	}

	v.Name = strings.ReplaceAll(v.Name, evon.FieldSplitter, evon.ObjectSplitter)
	v.Name = strings.ToLower(v.Name)

	var enum *evon.Node
	for _, n := range node.InnerNodes {
		switch n.Name[len(node.Name)+1:] {
		case "TYPE":
			v.Type = variableType(fmt.Sprint(n.Value))
		case "ENUM":
			enum = n
		default:

		}
	}
	if v.Type == "" {
		v.Type = VariableTypeStr
	}

	v.Value.val, err = mapVariableTypeToTypedValueConstructor[v.Type](node.Value)
	if err != nil {
		return errors.Wrap(err)
	}

	if enum != nil {
		constructor, ok := mapVariableTypeToEnumEvonNodeParser[v.Type]
		if !ok {
			return errors.Wrap(ErrNotAEnumable, "error unmarshalling evon enum")
		}

		v.Enum.v, err = constructor(enum)
		if err != nil {
			return errors.Wrap(err)
		}

		err = v.Enum.v.isEnum(v.Value.val)
		if err != nil {
			return errors.Wrap(err)
		}
	}

	return nil
}

func (v *Variable) EnumString() string {
	if v.Enum.v == nil {
		return ""
	}

	return v.Enum.v.EvonValue()
}

// MapVariableToGoType - maps variable onto golang's type name and import path
func MapVariableToGoType(variable Variable) (typeName string, importName string) {
	switch variable.Type {
	case VariableTypeInt:
		typeName = "int"
	case VariableTypeStr:
		typeName = "string"
	case VariableTypeBool:
		typeName = "bool"
	case VariableTypeFloat:
		typeName = "float64"
	case VariableTypeDuration:
		typeName = "time.Duration"
		importName = "time"
	default:
		return "any", ""
	}

	varRef := reflect.ValueOf(variable.Value)
	if varRef.Kind() == reflect.Slice {
		typeName = "[]" + typeName
	}

	return typeName, importName
}
