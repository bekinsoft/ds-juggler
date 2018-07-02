package juggler

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/jinzhu/gorm"
)

// opk = operator key
func whereFilter(key string, value interface{}, opk string, tx *gorm.DB, v ...string) *gorm.DB {
	opr := ""
	if operators[opk] == "" {
		if opk != key {
			panic("Invalid operator specified '" + opk + "'") // Raise proper error instead of panic. Or remove this from block
		}
		opk = "="
	}
	opr = " " + strings.ToUpper(operators[opk])
	// fmt.Println(key+opk+" ?", value)
	switch reflect.ValueOf(value).Kind().String() {
	case "string":
		if v != nil && v[0] == "or" {
			return tx.Or(key+opr+" ?", value.(string))
		}
		tx = tx.Where(key+opr+" ?", value.(string))
	case "float64":
		if v != nil && v[0] == "or" {
			return tx.Or(key+opr+" ?", value.(float64))
		}
		tx = tx.Where(key+opr+" ?", value.(float64))
	case "bool":
		if v != nil && v[0] == "or" {
			return tx.Or(key+opr+" ?", value.(bool))
		}
		tx = tx.Where(key+opr+" ?", value.(bool))
	case "slice":
		// array := value.([]interface{})
		// fmt.Println("Slice", array[0], reflect.ValueOf(array[0]).Kind().String())
		tx = tx.Where(key+opr+" (?)", value.([]interface{}))
		// tx = tx.Where(key+opr+" (?)", []string{"36644da7-77eb-11e8-aa14-02420aff035a", "70ba7eca-7961-11e8-81f5-34363bd1eeea"})
	default:
		// fmt.Println("============", value)
	}

	return tx
}

func limitFilter(res *Filter, tx *gorm.DB) *gorm.DB {
	datLimit := res.Limit
	if datLimit != nil && reflect.TypeOf(datLimit).Kind().String() == "float64" {
		tx = tx.Limit(datLimit)
	}
	return tx
}

func offsetFilter(res *Filter, tx *gorm.DB) *gorm.DB {
	datOffset := res.Offset
	if datOffset != nil && reflect.TypeOf(datOffset).Kind().String() == "float64" {
		tx = tx.Offset(datOffset)
	}
	return tx
}

func orderFilter(res *Filter, tx *gorm.DB) *gorm.DB {
	if res.Order != nil {
		datOrder := res.Order
		// USE COMMA-SEPERATED LIST FOR ORDER
		// EXAMPLE: ... "order": "username DESC, fullname, counter ASC "}
		// NOTE: by default all fields are in ascending so using the ASC is optional
		if reflect.TypeOf(datOrder).Kind().String() == "string" {
			tx = tx.Order(datOrder)
			// fmt.Println(datOrder)
		} else if reflect.TypeOf(datOrder).Kind().String() == "slice" {
			// May not be implemented
		}
	}

	return tx
}

func fieldFilter(res *Filter, tx *gorm.DB) *gorm.DB {
	if res.Fields != nil {
		datFields := res.Fields.(map[string]interface{})
		for key, value := range datFields {
			fmt.Println(key, "=>", value)
			// Will be implemented later
		}
	}

	return tx
}
