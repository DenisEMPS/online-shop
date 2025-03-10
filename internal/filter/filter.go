package filter

import (
	"fmt"
	"strings"
)

const (
	DataTypeStr  = "string"
	DataTypeInt  = "int"
	DataTypeDate = "date"
	DataTypeBool = "int"

	OperatorEq            = "eq"
	OperatorNotEq         = "neq"
	OperatorLowerThan     = "lt"
	OperatorLowerThanEq   = "lte"
	OperatorGreaterThan   = "gt"
	OperatorGreaterThanEq = "gte"
	OperatorBetween       = "between"
	OperatorLike          = "like"
)

type options struct {
	isToApply bool
	limit     int
	fields    []Field
}

type Field struct {
	Name     string
	Value    string
	Operator string
	Type     string
}

func NewOptions(limit int) *options {
	return &options{limit: limit}
}

type Options interface {
	GetLimit() int
	IsToApply() bool
	AddField(name, operator, value, dtype string) error
	Fileds() []Field
}

func (o *options) GetLimit() int {
	return o.limit
}

func (o *options) IsToApply() bool {
	return o.isToApply
}

func (o *options) AddField(name, operator, value, dtype string) error {
	if err := validateOperator(operator); err != nil {
		return err
	}

	o.isToApply = true

	o.fields = append(o.fields, Field{
		Name:     name,
		Value:    value,
		Operator: operator,
		Type:     dtype,
	})
	return nil
}

func (o *options) Fileds() []Field {
	return o.fields
}

func validateOperator(operator string) error {
	switch operator {
	case OperatorEq:
	case OperatorNotEq:
	case OperatorLowerThan:
	case OperatorLowerThanEq:
	case OperatorGreaterThan:
	case OperatorGreaterThanEq:
	case OperatorBetween:
	case OperatorLike:
	default:
		return fmt.Errorf("invalid operator")
	}
	return nil
}

func BuildQuery(options Options) (string, []interface{}, error) {
	var build strings.Builder
	var params []interface{}

	build.WriteString("WHERE ")

	for indx, field := range options.Fileds() {
		build.WriteString(fmt.Sprintf("%s ", field.Name))

		switch field.Operator {
		case OperatorEq:
			build.WriteString(fmt.Sprintf("= $%d", len(params)+1))
			params = append(params, field.Value)

		case OperatorNotEq:
			build.WriteString(fmt.Sprintf("!= $%d", len(params)+1))
			params = append(params, field.Value)

		case OperatorGreaterThan:
			build.WriteString(fmt.Sprintf("> $%d", len(params)+1))
			params = append(params, field.Value)

		case OperatorGreaterThanEq:
			build.WriteString(fmt.Sprintf(">= $%d", len(params)+1))
			params = append(params, field.Value)

		case OperatorLowerThan:
			build.WriteString(fmt.Sprintf("< $%d", len(params)+1))
			params = append(params, field.Value)

		case OperatorLowerThanEq:
			build.WriteString(fmt.Sprintf("<= $%d", len(params)+1))
			params = append(params, field.Value)

		case OperatorLike:
			build.WriteString(fmt.Sprintf("ILIKE $%d", len(params)+1))
			params = append(params, "%"+field.Value+"%")

		case OperatorBetween:
			rangeValues := strings.Split(field.Value, ":")
			if len(rangeValues) == 2 {
				build.WriteString(fmt.Sprintf("BETWEEN $%d AND $%d", len(params)+1, len(params)+2))
				params = append(params, rangeValues[0], rangeValues[1])
			} else {
				return "", nil, fmt.Errorf("invalid format for BETWEEN operator")
			}

		default:
			return "", nil, fmt.Errorf("invalid operator: %s", field.Operator)
		}

		if indx+1 < len(options.Fileds()) {
			build.WriteString(" AND ")
		}
	}

	return build.String(), params, nil
}
