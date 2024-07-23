package record

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/frankill/gotools/array"
)

type RecordAggType func(r *Record, rowIndex [][]int) (string, []any)

func recordAgg(fun func(slice ...[]int) int) func(filedName string) RecordAggType {

	return func(filedName string) RecordAggType {
		return func(r *Record, rowIndex [][]int) (string, []any) {

			data := r.SelectName(filedName).(Int)
			res := make([]any, len(rowIndex))

			for i := 0; i < len(rowIndex); i++ {

				res[i] = fun(array.ArrayChoose(rowIndex[i], data))

			}
			return filedName + "_sum", res
		}
	}
}

func RSum(filedName string) func(r *Record, rowIndex [][]int) (string, []any) {

	return recordAgg(array.ASum[[]int, int])(filedName)

}

func RMean(filedName string) func(r *Record, rowIndex [][]int) (string, []any) {

	return func(r *Record, rowIndex [][]int) (string, []any) {

		data := r.SelectName(filedName).(Int)
		res := make([]any, len(rowIndex))

		for i := 0; i < len(rowIndex); i++ {
			tmp := array.ArrayChoose(rowIndex[i], data)
			res[i] = array.ASum(tmp) / (len(tmp))

		}
		return filedName + "_mean", res
	}
}

func RMax(filedName string) func(r *Record, rowIndex [][]int) (string, []any) {

	return recordAgg(array.AMax[[]int, int])(filedName)
}

func RMin(filedName string) func(r *Record, rowIndex [][]int) (string, []any) {

	return recordAgg(array.AMin[[]int, int])(filedName)
}

type Record struct {
	name      string
	fieldNum  int
	fieldName []string
	rowNum    int
	typ       []int
	offer     []int
	int       []Int
	string    []String
	bool      []Bool
	float     []Float
}

func (r *Record) Head() {
	fmt.Printf("%s\t|", strings.Join(r.fieldName, "\t|"))
	fmt.Println()
	for i := 0; i < array.Min(10, r.rowNum); i++ {
		for j := 0; j < r.fieldNum; j++ {
			switch r.typ[j] {
			case 0:
				if i < len(r.int[r.offer[j]]) {
					fmt.Printf("%d ", r.int[r.offer[j]][i])
				}
			case 1:
				if i < len(r.string[r.offer[j]]) {
					fmt.Printf("%s ", r.string[r.offer[j]][i])
				}
			case 2:
				if i < len(r.bool[r.offer[j]]) {
					fmt.Printf("%t ", r.bool[r.offer[j]][i])
				}
			case 3:
				if i < len(r.float[r.offer[j]]) {
					fmt.Printf("%f ", r.float[r.offer[j]][i])
				}
			}
			fmt.Printf("\t|")
		}
		fmt.Println()
	}
}

func NewRecord(name string, length int) *Record {
	return &Record{
		name:      name,
		fieldName: []string{},
		typ:       make([]int, 0, length),
		offer:     make([]int, 0, length),
		int:       []Int{},
		string:    []String{},
		bool:      []Bool{},
		float:     []Float{},
	}
}

type GroupRecord struct {
	Record
}

func NewGroupRecord(name string, length int) *GroupRecord {
	return &GroupRecord{
		Record: *NewRecord(name, length),
	}
}

func (r *Record) WriteCsv(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	data := make([][]string, r.fieldNum)

	for i := 0; i < r.fieldNum; i++ {
		data[i] = r.SelectInt(i).ToString()
	}

	writer := csv.NewWriter(file)
	defer writer.Flush()

	writer.WriteAll(data)

	return nil
}

func (r *Record) Group(by []string, funs ...RecordAggType) *GroupRecord {

	if len(by) == 0 {
		log.Println("字段名称数组不能为空")
		return nil
	}

	if len(funs) == 0 {
		log.Println("至少需要一个函数作为参数")
		return nil
	}

	indexs := array.MapIntersect(array.ArrayMap(func(x ...string) map[string][]int {

		return array.GroupLocation(r.SelectName(x[0]).ToString())

	}, by)...)

	cols := array.ArrayZip(indexs.First...)

	res := NewGroupRecord(r.name, len(cols)+len(funs))

	for i := 0; i < len(cols); i++ {
		res.AddStringField(by[i], cols[i]...)
	}

	for _, fun := range funs {
		colname, data := fun(r, indexs.Second)
		res.AddField(colname, data...)
	}

	return res
}

func (r *Record) RowName() int {
	return r.rowNum
}

func (r *Record) FieldName() []string {
	return r.fieldName
}
func (r *Record) Name() string {
	return r.name
}
func (r *Record) Type() []string {
	return array.ArrayMap(func(x ...int) string {
		return r.GetType(x[0])
	}, r.typ)
}

func (r *Record) GetType(index int) string {
	switch r.typ[index] {
	case 0:
		return "Int"
	case 1:
		return "String"
	case 2:
		return "Bool"
	case 3:
		return "Float"
	default:
		return "Unknown"
	}
}

func (r *Record) GetType2(name string) string {
	index := array.MatchOne([]string{name}, r.fieldName)[0]
	return r.GetType(index)
}

func (r *Record) Slice(index []int) *Record {

	res := NewRecord(r.name, r.fieldNum)

	for i := 0; i < r.fieldNum; i++ {

		switch r.typ[i] {
		case 0:
			res.AddIntField(r.fieldName[i], array.ArrayChoose(index, r.int[r.offer[i]])...)
		case 1:
			res.AddStringField(r.fieldName[i], array.ArrayChoose(index, r.string[r.offer[i]])...)
		case 2:
			res.AddBoolField(r.fieldName[i], array.ArrayChoose(index, r.bool[r.offer[i]])...)
		case 3:
			res.AddFloatField(r.fieldName[i], array.ArrayChoose(index, r.float[r.offer[i]])...)
		}

	}
	return res

}

func (r *Record) SelectName(name string) Column {
	index := array.MatchOne([]string{name}, r.fieldName)[0]

	if index == -1 {
		log.Println("field not found")
		return nil
	}

	switch r.typ[index] {
	case 0:
		return r.int[r.offer[index]]
	case 1:
		return r.string[r.offer[index]]
	case 2:
		return r.bool[r.offer[index]]
	case 3:
		return r.float[r.offer[index]]
	default:
		return nil
	}
}

func (r *Record) SelectInt(index int) Column {
	index--

	if index < 0 || index >= r.fieldNum {
		log.Println("index out of range , index muts gte 1 ")
		return nil
	}

	switch r.typ[index] {
	case 0:
		return r.int[r.offer[index]]
	case 1:
		return r.string[r.offer[index]]
	case 2:
		return r.bool[r.offer[index]]
	case 3:
		return r.float[r.offer[index]]
	default:
		return nil
	}
}

func (r *Record) AddField(name string, x ...any) {
	r.rowNum = len(x)
	switch x[0].(type) {
	case int:
		r.AddIntField(name, array.ArrayToGeneric[int](x)...)
	case string:
		r.AddStringField(name, array.ArrayToGeneric[string](x)...)
	case bool:
		r.AddBoolField(name, array.ArrayToGeneric[bool](x)...)
	case float64:
		r.AddFloatField(name, array.ArrayToGeneric[float64](x)...)
	default:
		log.Println("type not supported")
	}
}

func (r *Record) AddIntField(name string, x ...int) {
	r.fieldNum++
	r.fieldName = append(r.fieldName, name)
	r.typ = append(r.typ, 0)
	r.offer = append(r.offer, len(r.int))
	r.int = append(r.int, Int(x))
}

func (r *Record) AddStringField(name string, x ...string) {
	r.fieldNum++
	r.fieldName = append(r.fieldName, name)
	r.typ = append(r.typ, 1)
	r.offer = append(r.offer, len(r.string))
	r.string = append(r.string, String(x))
}

func (r *Record) AddBoolField(name string, x ...bool) {
	r.fieldNum++
	r.fieldName = append(r.fieldName, name)
	r.typ = append(r.typ, 2)
	r.offer = append(r.offer, len(r.bool))
	r.bool = append(r.bool, Bool(x))
}

func (r *Record) AddFloatField(name string, x ...float64) {
	r.fieldNum++
	r.fieldName = append(r.fieldName, name)
	r.typ = append(r.typ, 3)
	r.offer = append(r.offer, len(r.float))
	r.float = append(r.float, Float(x))
}

type Column interface {
	ToString() []string
}

type Int []int

func (i Int) ToString() []string {
	return array.ArrayMap(func(x ...int) string {
		return strconv.Itoa(x[0])
	}, i)
}

type String []string

func (i String) ToString() []string {
	return i
}

type Bool []bool

func (i Bool) ToString() []string {
	return array.ArrayMap(func(x ...bool) string {
		if x[0] {
			return "true"
		}
		return "false"
	}, i)
}

type Float []float64

func (i Float) ToString() []string {
	return array.ArrayMap(func(x ...float64) string {
		return strconv.FormatFloat(x[0], 'f', 6, 64)
	}, i)
}
