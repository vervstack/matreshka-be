package evon

import (
	"strings"

	"go.redsock.ru/evon"
	errors "go.redsock.ru/rerrors"

	"go.vervstack.ru/matreshka/internal/domain"
	"go.vervstack.ru/matreshka/internal/service/user_errors"
)

type Validator struct {
	validConfigNameSymbols map[rune]struct{}
	validEnvNameSymbols    map[rune]struct{}
}

func newValidator() Validator {
	v := Validator{
		validConfigNameSymbols: make(map[rune]struct{}),
		validEnvNameSymbols:    make(map[rune]struct{}),
	}

	for _, r := range []rune(`ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789-_`) {
		v.validConfigNameSymbols[r] = struct{}{}
		v.validEnvNameSymbols[r] = struct{}{}
	}

	for _, r := range []rune(`abcdefghijklmnopqrstuvwxyz`) {
		v.validConfigNameSymbols[r] = struct{}{}
	}

	for _, r := range []rune(`/{}[]`) {
		v.validEnvNameSymbols[r] = struct{}{}
	}

	return v
}

func (v Validator) IsConfigNameValid(name domain.ConfigName) error {
	actualName := name.PlainName()
	if len(actualName) < 3 {
		return errors.Wrap(user_errors.ErrValidation,
			"Service name must be at least 3 symbols long")
	}

	invalidChars := make(map[rune]struct{})

	for _, r := range []rune(actualName) {
		if _, ok := v.validConfigNameSymbols[r]; !ok {
			invalidChars[r] = struct{}{}
		}
	}

	if len(invalidChars) > 0 {
		sb := strings.Builder{}
		sb.WriteString("Name contains invalid character")
		if len(invalidChars) > 1 {
			sb.WriteString("s")
		}

		sb.WriteString(": ")

		idx := 0
		for r := range invalidChars {
			sb.WriteRune(r)
			if idx < len(invalidChars)-1 {
				sb.WriteRune(',')
			}

			idx++
		}

		return errors.Wrap(user_errors.ErrValidation, sb.String())
	}

	return nil
}

func (v Validator) IsEnvNameValid(envName string) error {
	if len(envName) < 3 {
		return errors.Wrap(user_errors.ErrValidation,
			"Variable name must be at least 3 symbols long")
	}

	invalidChars := make(map[rune]struct{})

	for _, r := range []rune(envName) {
		if _, ok := v.validEnvNameSymbols[r]; !ok {
			invalidChars[r] = struct{}{}
		}
	}

	if len(invalidChars) > 0 {
		sb := strings.Builder{}
		sb.WriteString("Variable name contains invalid character")
		if len(invalidChars) > 1 {
			sb.WriteString("s")
		}

		sb.WriteString(": ")

		idx := 0
		for r := range invalidChars {
			sb.WriteRune(r)
			if idx < len(invalidChars)-1 {
				sb.WriteRune(',')
			}

			idx++
		}

		return errors.Wrap(user_errors.ErrValidation, sb.String())
	}

	return nil
}

const (
	environmentSegment = "ENVIRONMENT"
)

type VervConfigValidationResult struct {
	EnvUpsert []domain.PatchUpdate
	Upsert    []domain.PatchUpdate

	Invalid []domain.PatchUpdate
}

func (v Validator) AsVerv(original *evon.Node, patch *domain.PatchConfigRequest) VervConfigValidationResult {
	vervV := newVervConfigValidator(*patch)

	vervV.parseEnvironmentChanges(original)

	return vervV
}

func newVervConfigValidator(patch domain.PatchConfigRequest) VervConfigValidationResult {
	vervV := VervConfigValidationResult{}

	for _, p := range patch.Upsert {
		if strings.HasPrefix(p.FieldName, environmentSegment) {
			vervV.EnvUpsert = append(vervV.EnvUpsert, p)
		} else {
			vervV.Upsert = append(vervV.Upsert, p)
		}
	}

	return vervV
}

func (v *VervConfigValidationResult) parseEnvironmentChanges(original *evon.Node) {
	nodeStorage := evon.NodesToStorage(original)

	newEnvValues := make(map[string]domain.PatchUpdate)
	typesMap := make(map[string]domain.PatchUpdate)
	enumMap := make(map[string]domain.PatchUpdate)

	envUpsert := make([]domain.PatchUpdate, 0, len(newEnvValues))

	for _, valuePatch := range v.EnvUpsert {
		// already exists -> simply update value
		_, ok := nodeStorage[valuePatch.FieldName]
		if ok {
			envUpsert = append(envUpsert, valuePatch)
			continue
		}

		if strings.HasSuffix(valuePatch.FieldName, "_TYPE") {
			typesMap[valuePatch.FieldName[:len(valuePatch.FieldName)-5]] = valuePatch
			continue
		}

		if strings.HasSuffix(valuePatch.FieldName, "_ENUM") {
			enumMap[valuePatch.FieldName[:len(valuePatch.FieldName)-5]] = valuePatch
			continue
		}

		newEnvValues[valuePatch.FieldName] = valuePatch
	}

	for key, patchVal := range newEnvValues {
		envUpsert = append(envUpsert, patchVal)

		typeVal, ok := typesMap[key]
		if ok {
			envUpsert = append(envUpsert, typeVal)
		}

		enumVal, ok := enumMap[key]
		if ok {
			envUpsert = append(envUpsert, enumVal)
		}
	}

	v.EnvUpsert = envUpsert
}

func (v Validator) AsEvon(original evon.NodeStorage, patch *domain.PatchConfigRequest) (err error) {
	for i := range patch.Upsert {
		patch.Upsert[i].FieldName, err = v.normalizeAndValidateEnvName(patch.Upsert[i].FieldName)
		if err != nil {
			return errors.Wrap(err)
		}
	}

	for i := range patch.RenameTo {
		originalNode := original[patch.RenameTo[i].OldName]
		if originalNode == nil {
			continue
		}

		patch.RenameTo[i].OldName, err = v.normalizeAndValidateEnvName(patch.RenameTo[i].OldName)
		if err != nil {
			return errors.Wrap(err)
		}

		patch.RenameTo[i].NewName, err = v.normalizeAndValidateEnvName(patch.RenameTo[i].NewName)
		if err != nil {
			return errors.Wrap(err)
		}

		patch.RenameTo = append(patch.RenameTo,
			walkAndRename(originalNode.InnerNodes, patch.RenameTo[i].OldName, patch.RenameTo[i].NewName)...)
	}

	for i := range patch.Delete {
		originalNode := original[patch.Delete[i]]
		if originalNode == nil {
			continue
		}

		patch.Delete[i], err = v.normalizeAndValidateEnvName(patch.Delete[i])
		if err != nil {
			return errors.Wrap(err)
		}

		deletedChildren := childrenAsPlainSlice(originalNode)
		for _, c := range deletedChildren {
			patch.Delete = append(patch.Delete, c.Name)
		}
	}
	return nil
}

func walkAndRename(children []*evon.Node, oldName, newName string) []domain.PatchRename {
	out := make([]domain.PatchRename, 0, len(children))
	for _, child := range children {
		out = append(out, domain.PatchRename{
			OldName: child.Name,
			NewName: replaceAtStart(child.Name, oldName, newName),
		})
		out = append(out, walkAndRename(child.InnerNodes, oldName, newName)...)
	}

	return out
}

func childrenAsPlainSlice(root *evon.Node) []*evon.Node {
	out := make([]*evon.Node, 0, len(root.InnerNodes))
	for _, child := range root.InnerNodes {
		out = append(out, child)
		out = append(out, childrenAsPlainSlice(child)...)
	}

	return out
}

func replaceAtStart(str, old, new string) string {
	if !strings.HasPrefix(str, old) {
		return str
	}

	return new + str[len(old):]

}

func (v Validator) normalizeAndValidateEnvName(envName string) (string, error) {
	envName = strings.ToUpper(envName)
	return envName, v.IsEnvNameValid(envName)
}
