package testdata

import (
	"encoding/csv"
	"os"
	"reflect"
	"strconv"
	"strings"
	"unicode"
)

type csvUnmarshaller struct {
	source string
}

func (u *csvUnmarshaller) unmarshal(m interface{}) error {
	records, err := readRecords(u.source)
	if err != nil {
		return err
	}
	if len(records) == 1 {
		return nil
	}
	o := reflect.ValueOf(m)
	if o.Kind() == reflect.Ptr {
		o = o.Elem()
	}
	if o.Kind() == reflect.Slice {
		elemType := o.Type().Elem()
		if elemType.Kind() == reflect.Slice && elemType.Elem().Kind() == reflect.String {
			recordsVal := reflect.ValueOf(records)
			o.Set(recordsVal)
			return nil
		}
		isAddr := false
		if elemType.Kind() == reflect.Ptr {
			elemType = elemType.Elem()
			isAddr = true
		}
		if elemType.Kind() == reflect.Struct {
			content := records[1:]
			rows := make([]reflect.Value, len(content))
			fields := structFields(records[0], elemType)
			for i, record := range content {
				v := reflect.New(elemType).Elem()
				if err := unmarshalRecord(record, fields, &v); err != nil {
					return err
				}
				if isAddr {
					v = v.Addr()
				}
				rows[i] = v
			}
			o.Set(reflect.Append(o, rows...))
		}
	}
	return nil
}

func readRecords(filepath string) ([][]string, error) {
	f, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	r := csv.NewReader(f)
	records, err := r.ReadAll()
	if err != nil {
		return nil, err
	}
	return records, nil
}

func stripAllWhitespaces(s string) string {
	return strings.Map(func(r rune) rune {
		if unicode.IsSpace(r) {
			return -1
		}
		return r
	}, s)
}

func lookupField(columnName string, t reflect.Type) string {
	column := strings.ToLower(stripAllWhitespaces(columnName))
	field, found := t.FieldByNameFunc(func(name string) bool {
		ft, _ := t.FieldByName(name)
		tag, present := ft.Tag.Lookup("csv")
		if present {
			return columnName == tag
		}
		return column == strings.ToLower(name)
	})
	if found {
		return field.Name
	}
	return ""
}

func structFields(header []string, t reflect.Type) []string {
	fieldnames := make([]string, len(header))
	for i, columnName := range header {
		fieldnames[i] = lookupField(columnName, t)
	}
	return fieldnames
}

func unmarshalRecord(record []string, fields []string, s *reflect.Value) error {
	for i, field := range fields {
		if len(field) == 0 {
			continue
		}
		fieldVal := s.FieldByName(field)
		if !fieldVal.IsValid() {
			continue
		}
		if fieldVal.Kind() == reflect.Ptr {
			fieldVal = fieldVal.Elem()
		}
		t := fieldVal.Type()
		switch fieldVal.Kind() {
		case reflect.String:
			fieldVal.SetString(record[i])
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			val, err := strconv.ParseInt(record[i], 0, t.Bits())
			if err != nil {
				return err
			}
			fieldVal.SetInt(val)
		case reflect.Bool:
			val, err := strconv.ParseBool(record[i])
			if err != nil {
				return err
			}
			fieldVal.SetBool(val)
		case reflect.Float32, reflect.Float64:
			val, err := strconv.ParseFloat(record[i], t.Bits())
			if err != nil {
				return err
			}
			fieldVal.SetFloat(val)
		}
	}
	return nil
}

func newCSVUnmarshaller(sourceFilepath string) unmarshaller {
	return &csvUnmarshaller{sourceFilepath}
}
