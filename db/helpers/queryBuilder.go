package helpers

import "fmt"


func PatchQueryBuilder(table string, setParams map[string]any, whereParams map[string]any) string {

    paramCount := 0
    query := fmt.Sprintf(`UPDATE %s SET`, table)

    for key, value := range setParams {
        if value != nil {
            query += fmt.Sprintf(` %s = $%d,`, key, paramCount)
            paramCount+=1
        } 
    }

    if query[len(query)-1:] == "," {
        query = query[:len(query)-1]
    }

    if len(whereParams) == 0 {
        return query
    }

    query += ` WHERE`

    for key, value := range whereParams {
        if value != nil {
            query += fmt.Sprintf(` %s = $%d,`, key, paramCount)
            paramCount+=1
        }
    }

    return query[:len(query)-1]

}
