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
	key = ToSnakeCase(key) // Convert lowerCamelCase to snake_case for sql field to work
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

// Todo: Change this to a recursive function
// http://mindbowser.com/golang-go-with-gorm-2/
func includeFilter(res *Filter, tx *gorm.DB) *gorm.DB {
	if res.Include != nil {
		tx = include(res.Include, tx)
	}
	return tx
}

// Test Cases:
//	filter={"limit": 1, "include": "StudOther"}
//	filter={"limit": 1, "include": {"StudOther": "religion"}}
//	filter={"limit": 1, "include": ["StudDetail", "studOther"]}
//	filter={"limit": 1, "include": [{"StudOther": "religion"}]}
//	filter={"limit": 4, "include": ["studDetail", {"StudOther": "religion"}]}
func include(field interface{}, tx *gorm.DB) *gorm.DB {
	fmt.Println(field)
	switch reflect.TypeOf(field).Kind().String() {
	case "string":
		// fmt.Println("==>", field.(string))
		tx = tx.Preload(ToCamel(field.(string)))
	case "map":
		for ikey, ivalue := range field.(map[string]interface{}) {
			// fmt.Println("key=====>", ToCamel(ikey))
			tx = tx.Preload(ToCamel(ikey))

			if reflect.TypeOf(ivalue).Kind().String() == "string" {
				// fmt.Println("val=====>", ivalue)
				tx = tx.Preload(ToCamel(ikey) + "." + ToCamel(ivalue.(string))) // Nested Relationship : https://github.com/jinzhu/gorm/issues/392
			} else {
				return include(ivalue, tx)
			}
		}
	case "slice":
		for _, model := range field.([]interface{}) {
			if reflect.TypeOf(model).Kind().String() == "string" {
				// fmt.Println(">==>", model)
				tx = tx.Preload(ToCamel(model.(string)))
			} else if reflect.TypeOf(model).Kind().String() == "map" ||
				reflect.TypeOf(field).Kind().String() == "slice" {
				return include(model, tx)
			}
		}
	default:
		return nil
	}

	return tx
}

// func includeFilter(res *Filter, tx *gorm.DB) *gorm.DB {
// 	fmt.Println(res.Include)
// 	if res.Include != nil {
// 		switch reflect.TypeOf(res.Include).Kind().String() {
// 		case "string": // filter={"include": "emma"}
// 			fmt.Println(res.Include.(string))
// 			// tx = tx.Preload(ToCamel(res.Include.(string)))
// 		case "slice": // filter={"include": ["kofi", "emma"]}
// 			for _, model := range res.Include.([]interface{}) {
// 				// if reflect.TypeOf(model).Kind().String() == "string" {
// 				// fmt.Println(ToCamel(model.(string)))
// 				// 	// tx = tx.Preload(ToCamel(model.(string)))
// 				// }

// 				switch reflect.TypeOf(model).Kind().String() {
// 				case "string":
// 					fmt.Println(model.(string))
// 				case "slice":
// 					for _, model := range model.([]interface{}) {
// 						if reflect.TypeOf(model).Kind().String() == "string" {
// 							fmt.Println(ToCamel(model.(string)))
// 							// tx = tx.Preload(ToCamel(model.(string)))
// 						}
// 					}
// 				case "map":
// 					for jkey, jvalue := range model.(map[string]interface{}) {
// 						fmt.Println(ToCamel(jkey))
// 						fmt.Println(ToCamel(jvalue.(string)))
// 						// tx = tx.Preload(ToCamel(ikey))
// 					}
// 				}
// 			}
// 		case "map": // filter={"include": {"owner": "emma"}}
// 			// fmt.Println(res.Include)
// 			for ikey, ivalue := range res.Include.(map[string]interface{}) {
// 				fmt.Println(ToCamel(ikey))
// 				fmt.Println((ivalue))
// 				// tx = tx.Preload(ToCamel(ikey))

// 				switch reflect.TypeOf(ivalue).Kind().String() {
// 				case "string":
// 					fmt.Println(ivalue.(string))
// 				case "slice":
// 					for _, model := range ivalue.([]interface{}) {
// 						fmt.Println(model)
// 						// if reflect.TypeOf(model).Kind().String() == "string" {
// 						// 	fmt.Println(ToCamel(model.(string)))
// 						// 	// tx = tx.Preload(ToCamel(model.(string)))
// 						// }
// 					}
// 				case "map":
// 					for jkey, jvalue := range ivalue.(map[string]interface{}) {
// 						fmt.Println(ToCamel(jkey))
// 						fmt.Println(ToCamel(jvalue.(string)))
// 						// tx = tx.Preload(ToCamel(ikey))
// 					}
// 				}
// 			}
// 		}
// 	}

// 	return tx
// }

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
