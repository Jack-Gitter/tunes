package db

import (
	"fmt"
	"reflect"
	"time"
)


func PatchQueryBuilder(table string, setParams map[string]any, whereParams map[string]any, returning []string) (string, []any) {

    paramCount := 1
    values := []any{}
    query := fmt.Sprintf(`UPDATE %s SET`, table)

    for key, value := range setParams {
        if !reflect.ValueOf(value).IsNil() {
            query += fmt.Sprintf(` %s = $%d,`, key, paramCount)
            paramCount+=1
            values = append(values, getReflectValue(reflect.ValueOf(value).Elem()))
        } 
    }

    if query[len(query)-1:] == "," {
        query = query[:len(query)-1]
    }

    if len(whereParams) != 0 {
        query += ` WHERE`
        for key, value := range whereParams {
            query += fmt.Sprintf(` %s = $%d AND`, key, paramCount)
            paramCount+=1
            values = append(values, value)
        }
        query = query[:len(query)-3]
    }


    if len(returning) != 0 {
        query += ` RETURNING `
        for _, value := range returning {
            query += fmt.Sprintf(`%s, `, value)
        }
        query = query[:len(query)-1]
    }

    return query[:len(query)-1], values

}

func getReflectValue(val reflect.Value) any {
    switch val.Kind() {
    case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
        return val.Int()
    case reflect.String:
        return val.String()    
    }

    if v, ok := val.Interface().(time.Time); ok {
        return v
    }

    panic("we don't know this reflected value")
}
