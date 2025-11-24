package environment

import (
	"fmt"
	"slices"
	"strconv"
	"strings"

	"go.redsock.ru/evon"
	errors "go.redsock.ru/rerrors"
	"gopkg.in/yaml.v3"
)

const (
	mathMinus      = "-"
	rangeDelimiter = ":"
)

type intValue struct {
	v int
}

func (v *intValue) Val() any {
	return v.v
}

func (v *intValue) YamlValue() any {
	return v.v
}

func (v *intValue) EvonValue() string {
	return fmt.Sprintf("%d", v.v)
}

type intSliceValue struct {
	v []int
}

func (v *intSliceValue) Val() any {
	return v.v
}

func (v *intSliceValue) YamlValue() any {
	node := &yaml.Node{
		Kind:  yaml.SequenceNode,
		Style: yaml.FlowStyle,
	}

	if len(v.v) == 0 {
		return node
	}

	for _, r := range v.stringArray() {
		node.Content = append(node.Content,
			&yaml.Node{
				Kind:  yaml.ScalarNode,
				Value: r,
			})
	}

	return node
}

func (v *intSliceValue) EvonValue() string {
	return "[" + strings.Join(v.stringArray(), ",") + "]"
}

func (v *intSliceValue) isEnum(val typedValue) error {
	valuesToValidate := make([]int, 0, 1)

	switch actualValue := val.(type) {
	case *intValue:
		valuesToValidate = append(valuesToValidate, actualValue.v)
	case *intSliceValue:
		valuesToValidate = append(valuesToValidate, actualValue.v...)
	default:
		return errors.Wrapf(ErrUnexpectedType, "Expected Int or Slice of Ints but got %T", val)
	}

	for _, valueToValidate := range valuesToValidate {
		if !slices.Contains(v.v, valueToValidate) {
			return errors.Wrapf(ErrEnumValidationFailed, "got %d", valueToValidate)
		}
	}

	return nil
}

func toIntVariable(val any) (typedValue, error) {
	switch switchValue := val.(type) {
	case string:
		if switchValue[0] == '[' {
			v, err := extractIntSliceFromString(switchValue)
			return &intSliceValue{v: v}, err
		}

		rangeSeparator := strings.Index(switchValue, mathMinus)
		if rangeSeparator != -1 {
			v, err := extractIntRange(rangeSeparator, switchValue)
			return &intSliceValue{v: v}, err
		}

		v, err := strconv.Atoi(switchValue)
		return &intValue{v: v}, err

	case []interface{}:
		v, err := anySliceToIntSlice(switchValue)
		return &intSliceValue{v: v}, err
	case []int:
		return &intSliceValue{v: switchValue}, nil
	case []int8:
		return &intSliceValue{v: toIntSlice(switchValue)}, nil
	case []int16:
		return &intSliceValue{v: toIntSlice(switchValue)}, nil
	case []int32:
		return &intSliceValue{v: toIntSlice(switchValue)}, nil
	case []int64:
		return &intSliceValue{v: toIntSlice(switchValue)}, nil
	default:
		v, err := anyToInt(val)
		return &intValue{v: v}, err
	}
}

func intValueFromNode(node *yaml.Node) (typedValue, error) {
	if node.Kind == yaml.ScalarNode {
		if strings.HasPrefix(node.Value, "[") {
			return intSliceFromYamlNode(node)
		}

		i, err := strconv.Atoi(node.Value)
		return &intValue{i}, err
	}

	if node.Kind == yaml.SequenceNode {
		return intSliceFromYamlNode(node)
	}

	return nil, errors.New("Expected Int OR Int Slice type, got yaml %s", node.Tag)
}

func intSliceFromYamlNode(node *yaml.Node) (*intSliceValue, error) {
	intSlice := &intSliceValue{}

	var values []string

	if node.Kind == yaml.SequenceNode {
		for _, c := range node.Content {
			values = append(values, c.Value)
		}
	} else {
		stringedSlice := node.Value[1:]
		stringedSlice = stringedSlice[:len(stringedSlice)-1]
		values = strings.Split(stringedSlice, ",")
	}

	for _, v := range values {

		i := 0
		var err error
		rangeIndex := strings.Index(v, rangeDelimiter)
		if rangeIndex == -1 {
			i, err = strconv.Atoi(v)
			if err != nil {
				return nil, errors.Wrapf(err, "could not parse int. Expected single int value. Got '%s'", v)
			}

			intSlice.v = append(intSlice.v, i)
		} else {
			var firstInt, lastInt int
			firstInt, err = strconv.Atoi(v[:rangeIndex])
			if err != nil {
				return nil, errors.Wrap(err, "error parsing first int in sequence")
			}

			lastInt, err = strconv.Atoi(v[rangeIndex+1:])
			if err != nil {
				return nil, errors.Wrap(err, "error parsing last int in sequence")
			}

			for ; firstInt <= lastInt; firstInt++ {
				intSlice.v = append(intSlice.v, firstInt)
			}
		}
	}

	return intSlice, nil
}

func intSliceFromEvonNode(node *evon.Node) (*intSliceValue, error) {
	switch v := node.Value.(type) {
	case string:
		if strings.HasPrefix(v, "[") {
			v = v[1:]
			if strings.HasSuffix(v, "]") {
				v = v[:len(v)-1]
			}
		}

		splited := strings.Split(v, ",")
		isv := &intSliceValue{
			v: make([]int, 0, len(splited)),
		}

		for _, oneV := range splited {
			intV, err := strconv.Atoi(oneV)
			if err != nil {
				return nil, errors.Wrap(err, "error parsing int value from string")
			}
			isv.v = append(isv.v, intV)
		}

		return isv, nil

	case []int:
		return &intSliceValue{v: v}, nil

	default:
		return nil, errors.Wrap(ErrUnexpectedType, "expected string or slice of int")
	}

}
func extractIntSliceFromString(switchValue string) ([]int, error) {
	separatedVals := strings.Split(switchValue[1:len(switchValue)-1], ",")

	anyVals := make([]any, 0, len(separatedVals))
	for _, v := range separatedVals {
		anyVals = append(anyVals, v)
	}

	return anySliceToIntSlice(anyVals)
}

func anySliceToIntSlice(value []any) ([]int, error) {
	out := make([]int, 0, len(value))

	for _, v := range value {
		switch v := v.(type) {
		case string:
			rangeSeparator := strings.Index(v, rangeDelimiter)
			if rangeSeparator != -1 {
				rng, err := extractIntRange(rangeSeparator, v)
				if err != nil {
					return nil, errors.Wrap(err, "error converting value to int range")
				}

				out = append(out, rng...)
			} else {
				newInt, err := strconv.Atoi(v)
				if err != nil {
					return nil, errors.Wrap(err, "error converting value to int")
				}

				out = append(out, newInt)
			}
		default:
			val, err := anyToInt(v)
			if err != nil {
				return nil, errors.Wrap(err, "error converting any to int")
			}

			out = append(out, val)
		}
	}

	return out, nil
}

func anyToInt(val any) (int, error) {
	switch switchValue := val.(type) {
	case int:
		return switchValue, nil
	case int8:
		return int(switchValue), nil
	case int16:
		return int(switchValue), nil
	case int32:
		return int(switchValue), nil
	case int64:
		return int(switchValue), nil
	case uint:
		return int(switchValue), nil
	case uint8:
		return int(switchValue), nil
	case uint16:
		return int(switchValue), nil
	case uint32:
		return int(switchValue), nil
	case uint64:
		return int(switchValue), nil
	default:
		return 0, errors.New(fmt.Sprintf("can't cast %T to int", val))
	}
}

func extractIntRange(rangeSeparatorIdx int, strValue string) ([]int, error) {
	firstNumber, err := strconv.Atoi(strValue[:rangeSeparatorIdx])
	if err != nil {
		return nil, errors.Wrap(err, "error parsing first number of range to int")
	}

	secondNumber, err := strconv.Atoi(strValue[rangeSeparatorIdx+1:])
	if err != nil {
		return nil, errors.Wrap(err, "error parsing second number of range to int")
	}

	out := make([]int, 0, secondNumber-firstNumber)

	for i := firstNumber; i <= secondNumber; i++ {
		out = append(out, i)
	}

	return out, nil
}

func toIntSlice[T int | int8 | int16 | int32 | int64](v []T) []int {
	out := make([]int, 0, len(v))
	for _, v := range v {
		out = append(out, int(v))
	}

	return out
}

func (v *intSliceValue) stringArray() []string {
	if len(v.v) == 0 {
		return []string{}
	}

	convertToRange := func(start, end int) string {
		newRange := strconv.Itoa(start)
		if start != end {
			newRange += rangeDelimiter + strconv.Itoa(end)
		}

		return newRange
	}

	prev := v.v[0]
	rangeStart := prev

	out := make([]string, 0, min(len(v.v), 16))

	for _, v := range v.v[1:] {
		if v-prev != 1 {
			out = append(out, convertToRange(rangeStart, prev))

			prev = v
			rangeStart = v
		}
		prev = v
	}

	out = append(out, convertToRange(rangeStart, prev))

	return out
}
