package matreshka

import (
	"reflect"
	"sort"
	"strings"

	"go.redsock.ru/evon"
	errors "go.redsock.ru/rerrors"

	"go.vervstack.ru/matreshka/internal/utils/cases"
	"go.vervstack.ru/matreshka/pkg/matreshka/environment"
)

var ErrNotAPointer = errors.New("not a pointer")

type Environment []*environment.Variable

func (a *Environment) MarshalEnv(prefix string) ([]*evon.Node, error) {
	if prefix != "" {
		prefix += evon.ObjectSplitter
	}

	out := make([]*evon.Node, 0, len(*a))
	for _, v := range *a {
		pref := prefix + strings.NewReplacer(" ", evon.FieldSplitter, evon.ObjectSplitter, evon.FieldSplitter).Replace(strings.ToUpper(v.Name))

		root := &evon.Node{
			Name:  pref,
			Value: v.Value.String(),
			InnerNodes: []*evon.Node{{
				Name:  pref + "_TYPE",
				Value: v.Type,
			},
			},
		}

		if v.Enum.Value() != nil {
			root.InnerNodes = append(root.InnerNodes,
				&evon.Node{
					Name:  pref + "_ENUM",
					Value: v.EnumString(),
				})
		}
		out = append(out, root)
	}

	sort.Slice(out, func(i, j int) bool {
		return out[i].Name < out[j].Name
	})

	return out, nil
}

func (a *Environment) UnmarshalEnv(rootNode *evon.Node) error {
	*a = make([]*environment.Variable, len(rootNode.InnerNodes))

	for idx := range rootNode.InnerNodes {
		(*a)[idx] = &environment.Variable{}

		err := (*a)[idx].UnmarshalEnv(rootNode.InnerNodes[idx])
		if err != nil {
			return errors.Wrap(err, "error unmarshalling environment variable")
		}
	}

	sort.Slice(*a, func(i, j int) bool {
		return (*a)[i].Name < (*a)[j].Name
	})

	return nil
}

func (a *Environment) ParseToStruct(dst any) error {
	dstRef := reflect.ValueOf(dst)
	if dstRef.Kind() != reflect.Ptr {
		return errors.Wrap(ErrNotAPointer, "expected destination to be a pointer ")
	}

	dstRef = dstRef.Elem()
	numFields := dstRef.NumField()

	dstMapping := make(map[string]reflect.Value)

	for i := 0; i < numFields; i++ {
		field := dstRef.Type().Field(i)
		dstMapping[field.Name] = dstRef.Field(i)
	}

	for _, env := range *a {
		name := env.Name
		name = strings.ReplaceAll(name, " ", evon.ObjectSplitter)
		name = cases.SnakeToPascal(name)
		v, ok := dstMapping[name]
		if !ok {
			return errors.Wrap(ErrNotFound, "field with name "+name+" can't be found in target struct")
		}

		val := reflect.ValueOf(env.Value.Value())
		if val.Type().AssignableTo(v.Type()) {
			v.Set(val)
		} else if val.Type().ConvertibleTo(v.Type()) {
			v.Set(val.Convert(v.Type()))
		} else {
			return errors.Wrapf(ErrUnexpectedType, "field %s has type %s but value has type %s", name, v.Type(), val.Type())
		}
	}

	return nil
}
