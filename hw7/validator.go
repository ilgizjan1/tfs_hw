package homework

import (
	"fmt"
	"github.com/pkg/errors"
	"reflect"
	"strconv"
	"strings"
)

var ErrNotStruct = errors.New("wrong argument given, should be a struct")
var ErrInvalidValidatorSyntax = errors.New("invalid validator syntax")
var ErrValidateForUnexportedFields = errors.New("validation for unexported field is not allowed")
var ErrInvalidElementType = errors.New("invalid element type")
var ErrInvalidIn = errors.New("the value is not contained in the set In")
var ErrInvalidMax = errors.New("the value is greater than the max")
var ErrInvalidMin = errors.New("the value is less than the min")
var ErrInvalidLen = errors.New("the length is not equal to len")

type ValidationError struct {
	Err error
}

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	var s string
	for _, err := range v {
		s += fmt.Sprintf("%s", err.Err)
	}
	return s
}

func (v *ValidationErrors) Add(err error) {
	if err == nil {
		return
	}
	*v = append(*v, ValidationError{
		Err: err,
	})

}

func Validate(v any) error {
	var errorSlice ValidationErrors
	value := reflect.ValueOf(v)
	if value.Type().Kind() != reflect.Struct {
		return ErrNotStruct
	}
	for i := 0; i < value.Type().NumField(); i++ {
		field := value.Type().Field(i)
		fieldValue := value.Field(i)
		validator, ok := field.Tag.Lookup("validate")
		if !ok {
			continue
		}
		if !field.IsExported() {
			errorSlice.Add(ErrValidateForUnexportedFields)
			continue
		}
		kind, err := getKind(field)
		if err != nil {
			errorSlice.Add(err)
			continue
		}
		switch kind {
		case reflect.Struct:
			err := Validate(fieldValue.Interface())
			if err != nil {
				errorSlice = append(errorSlice, err.(ValidationErrors)...)
			}
		default:
			splittedValidator := strings.Split(validator, ":")
			tagHandler := newTagHandler(kind, fieldValue, field.Name, splittedValidator[0], splittedValidator[1])
			errorSlice = append(errorSlice, tagHandler.validate(field.Type)...)
		}
	}
	if len(errorSlice) == 0 {
		return nil
	}
	return errorSlice
}

type tagHandler struct {
	kind      reflect.Kind
	value     reflect.Value
	fieldName string
	tagName   string
	tagValue  string
}

func newTagHandler(kind reflect.Kind, value reflect.Value, fieldName string, tagName string, tagValue string) tagHandler {
	return tagHandler{
		kind:      kind,
		value:     value,
		fieldName: fieldName,
		tagName:   tagName,
		tagValue:  tagValue,
	}
}

func (t *tagHandler) validate(fieldType reflect.Type) ValidationErrors {
	value := t.value
	var errorSlice ValidationErrors
	switch fieldType.Kind() {
	case reflect.Slice:
		for i := 0; i < value.Len(); i++ {
			t.value = value.Index(i)
			errorSlice.Add(t.processor())
		}
	default:
		errorSlice.Add(t.processor())
	}
	if len(errorSlice) == 0 {
		return nil
	}
	return errorSlice
}

func (t *tagHandler) processor() error {
	if t.tagName == "in" {
		in := strings.Split(t.tagValue, ",")
		return t.validateIn(in)
	}
	val, err := strconv.Atoi(t.tagValue)
	if err != nil {
		return ErrInvalidValidatorSyntax
	}
	switch t.tagName {
	case "len":
		return t.validateLen(val)
	case "max":
		return t.validateMax(val)
	case "min":
		return t.validateMin(val)
	default:
		return fmt.Errorf("%s : %w\n", t.fieldName, ErrInvalidValidatorSyntax)
	}
}

func (t tagHandler) validateLen(length int) error {
	if length < 0 {
		return ErrInvalidValidatorSyntax
	}
	if len(t.value.String()) != length {
		return fmt.Errorf("%s : %w\n", t.fieldName, ErrInvalidLen)
	}
	return nil
}

func (t tagHandler) validateMin(min int) error {
	var val int
	switch t.kind {
	case reflect.String:
		val = len(t.value.String())
	default:
		val = int(t.value.Int())
	}
	if val < min {
		return fmt.Errorf("%s : %w\n", t.fieldName, ErrInvalidMin)
	}
	return nil
}

func (t tagHandler) validateMax(max int) error {
	var val int
	switch t.kind {
	case reflect.String:
		if max < 0 {
			return fmt.Errorf("%s : %w\n", t.fieldName, ErrInvalidValidatorSyntax)
		}
		val = len(t.value.String())
	default:
		val = int(t.value.Int())
	}
	if val > max {
		return fmt.Errorf("%s : %w\n", t.fieldName, ErrInvalidMax)
	}
	return nil
}

func (t tagHandler) validateIn(strSlice []string) error {
	if len(strSlice) == 1 && len(strSlice[0]) == 0 {
		return fmt.Errorf("%s : %w\n", t.fieldName, ErrInvalidValidatorSyntax)
	}
	var err error
	switch t.kind {
	case reflect.String:
		err = contains[string](strSlice, t.value.String())
	default:
		var slice []int
		slice, err = converter(strSlice)
		if err != nil {
			return fmt.Errorf("%s : %w\n", t.fieldName, err)
		}
		err = contains(slice, int(t.value.Int()))
	}
	if err != nil {
		return fmt.Errorf("%s : %w\n", t.fieldName, err)
	}
	return nil
}

func converter(strSlice []string) ([]int, error) {
	intSlice := make([]int, len(strSlice))
	var err error
	for i, val := range strSlice {
		intSlice[i], err = strconv.Atoi(val)
		if err != nil {
			return nil, ErrInvalidValidatorSyntax
		}
	}
	return intSlice, nil
}

func contains[T comparable](t []T, needed T) error {
	for _, v := range t {
		if v == needed {
			return nil
		}
	}
	return ErrInvalidIn
}

func getKind(field reflect.StructField) (reflect.Kind, error) {
	switch field.Type.Kind() {
	case reflect.Struct:
		return reflect.Struct, nil
	case reflect.Int:
		return reflect.Int, nil
	case reflect.String:
		return reflect.String, nil
	case reflect.Slice:
		switch field.Type.String() {
		case "[]int":
			return reflect.Int, nil
		case "[]string":
			return reflect.String, nil
		}
	}
	return reflect.Invalid, fmt.Errorf("%s : %w", field.Name, ErrInvalidElementType)
}
