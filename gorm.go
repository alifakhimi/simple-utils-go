package simutils

import (
	"context"
	"fmt"
	"net/url"
	"reflect"
	"regexp"
	"strings"
	"sync"

	"github.com/spf13/cast"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"
)

var (
	// sortRegex defines the regex pattern for sorting
	sortRegex = regexp.MustCompile(`^(\w+)(?::(asc|desc))?$`)

	mapURLToDBOperator = map[string]string{
		"eq":       "=",
		"neq":      "<>",
		"gt":       ">",
		"gte":      ">=",
		"lt":       "<",
		"lte":      "<=",
		"like":     "like",
		"nlike":    "NOT like",
		"similar":  "similar to",
		"nsimilar": "NOT similar to",
		"in":       "IN",
		"nin":      "NOT IN",
	}

	mapURLToDBOrder = map[string]string{"asc": "ASC", "desc": "DESC"}
)

// GetTableName get table name using gorm
func GetTableName(val interface{}) string {
	tblSchema, _ := schema.Parse(val, &sync.Map{}, schema.NamingStrategy{})
	return tblSchema.Table
}

// GetID get id from database.Model
func GetID(val interface{}) PID {
	reflectValue := reflect.Indirect(reflect.ValueOf(val))
	for reflectValue.Kind() == reflect.Ptr || reflectValue.Kind() == reflect.Interface {
		reflectValue = reflect.Indirect(reflectValue)
	}

	var value reflect.Value

	switch reflectValue.Kind() {
	case reflect.Struct:
		value = reflectValue.FieldByName("ID")
	}

	return Parse(fmt.Sprintf("%s", value))
}

func ParseFilters(db *gorm.DB, driver DatabaseDriver, filters map[string][]FilterValue, mapKeyToColumn map[string][]string) (*gorm.DB, error) {
	var (
		err error
	)

	for fk, fvs := range filters {
		for _, fv := range fvs {
			var (
				query []string
				args  []interface{}
			)

			cols, ok := mapKeyToColumn[fk]

			if !ok || fv.Value == nil || len(fmt.Sprintf("%v", fv.Value)) == 0 {
				continue
			}

			if strings.ToLower(fv.Operator) == "in" {
				for _, col := range cols {
					ins := strings.Split(fv.Value.(string), ",")
					if len(ins) > 2000 {
						values := make([]string, len(ins))
						for i, s := range ins {
							values[i] = fmt.Sprintf("('%s')", s)
						}

						query = append(query, fmt.Sprintf("%s %s (%s as tbl(id))", col, mapURLToDBOperator[fv.Operator], fmt.Sprintf("select * from (values %s)", strings.Join(values, ","))))
						// args = append(args, ins)
					} else {
						query = append(query, fmt.Sprintf("%s %s (?)", col, mapURLToDBOperator[fv.Operator]))
						args = append(args, ins)
					}

					/*
						limit := 1000
						for offset := 0; offset < len(inArray); offset += limit {
							query = append(query, fmt.Sprintf("%s %s (select id from (values ))", col, mapURLToDBOperator[fv.Operator]))

							if offset+limit > len(inArray) {
								limit = len(inArray) - offset
							}

							args = append(args, inArray[offset:offset+limit])
						}*/
				}
			} else {
				for _, col := range cols {
					query = append(query, fmt.Sprintf("%s %s ?", col, mapURLToDBOperator[fv.Operator]))
					args = append(args, CorrectSimilarChars(driver, fv.Value))
				}
			}

			if fv.Or {
				db = db.Or(strings.Join(query, " OR "), args...)
			} else {
				db = db.Where(strings.Join(query, " OR "), args...)
			}
		}
	}

	return db, err
}

func ParseSorts(db *gorm.DB, sorts []SortValue, mapKeyToColumn map[string][]string) (*gorm.DB, error) {
	var (
		err error
	)

	for _, sv := range sorts {
		cols, ok := mapKeyToColumn[sv.Key]
		if !ok || len(cols) == 0 {
			continue
		}

		db = db.Order(fmt.Sprintf("%s %s", cols[0], mapURLToDBOrder[sv.Order]))
	}

	return db, err
}

func BuildGormQuery(ctx *context.Context, db *gorm.DB, queryParams url.Values) *gorm.DB {
	// init gorm db
	qb := db

	for field, values := range queryParams {
		// continue if value is empty
		if len(values) == 0 {
			continue
		}

		// extend case if more advanced query is necessary
		// all matches values will reach the default case
		switch field {
		case "search":
			qb = qb.Where("search like ?", values[0])
			for i := 1; i < len(values); i++ {
				qb = qb.Or("search like ?", values[i])
			}
		case "limit":
			qb = qb.Limit(cast.ToInt(values[0]))
		case "offset":
			qb = qb.Offset(cast.ToInt(values[0]))
		case "sort":
			orderByColumns := []clause.OrderByColumn{}
			for _, param := range values {
				param = strings.ToLower(param)
				if sortRegex.MatchString(param) {
					matches := sortRegex.FindStringSubmatch(param)
					orderByColumns = append(orderByColumns,
						clause.OrderByColumn{
							Column: clause.Column{
								Name: matches[1],
							},
							Desc: matches[2] == "desc",
						},
					)
				}
			}
			qb = qb.Order(clause.OrderBy{Columns: orderByColumns})
		case "includes":
			for _, inc := range values {
				qb = qb.Preload(inc)
			}
		default:
			if len(values) == 1 {
				qb = qb.Where(fmt.Sprintf("%s = ?", field), values)
			} else {
				qb = qb.Where(fmt.Sprintf("%s in ?", field), values)
			}
		}
	}

	return qb
}
