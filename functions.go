/*
 * @author    Emmanuel Kofi Bessah
 * @email     ekbessah@uew.edu.gh
 * @created   Sat Jun 30 2018 11:41:21
 * @copyright © 2018 University of Education, Winneba
 */

package juggler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"reflect"
	"strings"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
)

// Filter specify criteria for the returned data set
type Filter struct {
	Where   interface{}
	Limit   interface{}
	Offset  interface{}
	Order   interface{}
	Include interface{}
	Fields  interface{}

	Valid bool
}

// FilterRequest find, findByID, count anad exists request
type FilterRequest struct {
	Filter, Body     interface{}
	Method, RawQuery string
	Params           map[string]string
}

// GetFilterParamMap returns a map of the filter request
func GetFilterParamMap(r *http.Request) (Filter, error) {
	filter := Filter{Valid: true}
	urlQuery := r.URL.RawQuery
	// fmt.Println(urlQuery)
	// fmt.Println(strings.Trim(urlQuery, " "))
	// urlStr := `members?filter=%7B%22where%22%3A%20%7B%22and%22%3A%20%5B%7B%22period%22%3A%201804%7D%2C%20%7B%22deleted%22%3A%200%7D%2C%20%7B%22approved%22%3A%201%7D%5D%7D%7D'`
	// var dat map[string]interface{}
	if strings.Trim(urlQuery, " ") == "" {
		filter.Valid = false
		return filter, nil
	}

	// fmt.Println(urlQuery)
	decode, err := url.PathUnescape(urlQuery)

	// fmt.Println(decode)

	if err != nil {
		return filter, err
	}

	if strings.Index(decode, "filter=") < 0 || strings.Index(decode, "filter=") > 0 {
		return filter, fmt.Errorf("Invalid query prefix provided")
	}

	uJSON := strings.TrimLeft(decode, "filter=")
	// println(uJSON)
	byt := []byte(uJSON)

	if err := json.Unmarshal(byt, &filter); err != nil {
		return filter, fmt.Errorf("Value is not an object")
	}

	// Get vars from request and check whether we're finding by id
	// If so do not filter with WHERE, ORDER, LIMIT and OFFSET
	vars := mux.Vars(r)
	if _, ok := vars["id"]; ok {
		filter.Where = nil
		filter.Order = nil
		filter.Limit = nil
		filter.Offset = nil
	}

	return filter, nil
}

// GetFilterParamMapFromJSONString returns a map of the filter request when given json string
func GetFilterParamMapFromJSONString(jsonstr string) (Filter, error) {
	filter := Filter{Valid: true}

	byt := []byte(jsonstr)

	if err := json.Unmarshal(byt, &filter); err != nil {
		return filter, fmt.Errorf("Value is not an object")
	}

	return filter, nil
}

// FilterQuery takes request and GORM transaction instance
func FilterQuery(f Filter, tx *gorm.DB) (*gorm.DB, error) {
	if f.Where != nil {
		datwhere := f.Where.(map[string]interface{})
		for key, value := range datwhere {
			// fmt.Println(key, "=>", value)
			// fmt.Println(reflect.ValueOf(value).Kind().String())
			if reflect.ValueOf(value).Kind().String() == "map" {
				datwhere = value.(map[string]interface{})
				for ikey, ivalue := range datwhere {
					// fmt.Println(ikey, "==>", ivalue)
					// fmt.Println(key, "==>--", value)

					tx = whereFilter(key, ivalue, ikey, tx)

					if reflect.ValueOf(ivalue).Kind().String() == "map" {
						datwhere = ivalue.(map[string]interface{})

						for jkey, jvalue := range datwhere {
							fmt.Println(jkey, "===>", jvalue)
						}
					}
				}
			} else if reflect.ValueOf(value).Kind().String() == "slice" {
				// Handle AND and OR arrays
				arrwhere := value.([]interface{})
				for _, isvalue := range arrwhere {
					if reflect.ValueOf(isvalue).Kind().String() == "map" {
						datwhere = isvalue.(map[string]interface{})
						for ikey, ivalue := range datwhere {
							// fmt.Println(ikey, "==>", ivalue)
							// fmt.Println(key, "==>--", value)

							if key == "and" || key == "or" {
								tx = whereFilter(ikey, ivalue, ikey, tx, key)
							}

							if reflect.ValueOf(ivalue).Kind().String() == "map" {
								datwhere = ivalue.(map[string]interface{})

								for jkey, jvalue := range datwhere {
									// fmt.Println(jkey, "==>>", jvalue)
									// fmt.Println(ikey, "==>>--", ivalue)
									tx = whereFilter(ikey, jvalue, jkey, tx)
								}
							}
						}
					}
				}
			} else {
				tx = whereFilter(key, value, key, tx)
			}
		}
	}

	tx = orderFilter(&f, tx)
	tx = offsetFilter(&f, tx)
	tx = limitFilter(&f, tx)
	tx = fieldFilter(&f, tx)
	tx = includeFilter(&f, tx)

	// tx.Find(&mod)
	// fmt.Println(mod)
	return tx, nil
}
