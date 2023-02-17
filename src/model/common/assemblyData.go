package common

import (
	"database/sql"
)

// AssemblyData 组装查询的数据并返回数组
func AssemblyData(rows *sql.Rows) []map[string]string {
	//获取列名
	columns, _ := rows.Columns()
	//定义一个切片,长度是字段的个数,切片里面的元素类型是sql.RawBytes
	values := make([]sql.RawBytes, len(columns))
	//定义一个切片,元素类型是interface{} 接口
	scanArgs := make([]interface{}, len(values))
	for i := range values {
		//把sql.RawBytes类型的地址存进去了
		scanArgs[i] = &values[i]
	}
	//获取字段值
	result := make([]map[string]string, 0)
	for rows.Next() {
		res := make(map[string]string)
		_ = rows.Scan(scanArgs...)
		for i, col := range values {
			res[columns[i]] = string(col)
		}
		result = append(result, res)
	}
	return result
}
