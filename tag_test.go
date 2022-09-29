package schema

import (
	"sort"
	"strings"
	"testing"
)

func _sameStrSlice(a []string, b []string) bool {
	result := true

	sort.Strings(a)
	sort.Strings(b)

	if len(a) != len(b) {
		return false
	}

	for i, v := range a {
		if v != b[i] {
			result = false
		}
	}

	return result
}

func _sameIntSlice(a []int64, b []int64) bool {
	result := true
	if len(a) != len(b) {
		return false
	}

	a64 := []int{}
	b64 := []int{}
	for _, v := range a {
		a64 = append(a64, int(v))
	}
	for _, v := range b {
		b64 = append(b64, int(v))
	}

	sort.Ints(a64)
	sort.Ints(b64)

	for i, v := range a {
		if v != b[i] {
			result = false
		}
	}

	return result
}

func _sameUintSlice(a []uint64, b []uint64) bool {
	result := true
	if len(a) != len(b) {
		return false
	}

	a64 := []int{}
	b64 := []int{}
	for _, v := range a {
		a64 = append(a64, int(v))
	}
	for _, v := range b {
		b64 = append(b64, int(v))
	}

	sort.Ints(a64)
	sort.Ints(b64)

	for i, v := range a {
		if v != b[i] {
			result = false
		}
	}

	return result
}

func TestTagStructure(t *testing.T) {

	type Test struct {
		Input    string
		Expected Tag
	}

	tests := []Test{
		{
			Input:    "",
			Expected: Tag{Values: map[string]string{}},
		},
		{
			Input:    " \r\n\t",
			Expected: Tag{Values: map[string]string{}},
		},
		{
			Input:    "field_01:this is a test,field_02:this 'is another' test,field_03:1,field_04:123.321,field_05:1 2 3 4,field_06:1.1 1.2 1.3 1.4",
			Expected: Tag{Values: map[string]string{}},
		},
	}

	for i, test := range tests {
		tag := ParseTag(test.Input)

		if strings.Trim(test.Input, " \n\r\t") == "" && !tag.IsEmpty() {
			t.Errorf("Test %d fail, Expected Empty Tag", i)
		} else if strings.Trim(test.Input, " \n\r\t") != "" && tag.IsEmpty() {
			t.Errorf("Test %d fail, Expected Not Empty Tag", i)
		} else if !tag.IsEmpty() {

			// String tests
			gotStr := tag.GetString("field_01")
			wantStr := "this is a test"
			if gotStr != wantStr {
				t.Errorf("Test %d fail, Unexpected value for field \"field_01\":\nGot:  %s\nWant: %s", i, gotStr, wantStr)
			}
			gotStr = tag.GetString("field_02")
			wantStr = "this 'is another' test"
			if gotStr != wantStr {
				t.Errorf("Test %d fail, Unexpected value for field \"field_03\":\nGot:  %s\nWant: %s", i, gotStr, wantStr)
			}
			gotSliceStr := tag.GetSliceString("field_02", ' ')
			wantSliceStr := []string{"this", "is another", "test"}
			if !_sameStrSlice(gotSliceStr, wantSliceStr) {
				t.Errorf("Test %d fail, Unexpected value for field \"field_03\":\nGot:  %v\nWant: %v", i, gotSliceStr, wantSliceStr)
			}

			// Int tests
			gotInt := tag.GetInt("field_03")
			wantInt := int64(1)
			if gotInt != wantInt {
				t.Errorf("Test %d fail, Unexpected value for field \"field_03\":\nGot:  %d\nWant: %d", i, gotInt, wantInt)
			}
			gotInt = tag.GetInt("field_01")
			if gotInt != 0 {
				t.Errorf("Test %d fail, Unexpected value for field \"field_01\":\nGot:  %d\nWant: 0", i, gotInt)
			}
			gotInt = tag.GetInt("field_03")
			if gotInt == 0 {
				t.Errorf("Test %d fail, Unexpected value for field \"field_03\":\nGot:  0\nWant: 1", i)
			}
			gotSliceInt := tag.GetSliceInt("field_05", ' ')
			wantSliceInt := []int64{1, 2, 3, 4}
			if !_sameIntSlice(gotSliceInt, wantSliceInt) {
				t.Errorf("Test %d fail, Unexpected value for field \"field_05\":\nGot:  %v\nWant: %v", i, gotSliceInt, wantSliceInt)
			}
			gotSliceInt = tag.GetSliceInt("field_99", ' ')
			if len(gotSliceInt) > 0 {
				t.Errorf("Test %d fail, Unexpected value for field \"field_99\":\nGot:  %v\nWant: []", i, gotSliceInt)
			}

			// Uint tests
			gotUint := tag.GetUint("field_03")
			wantUint := uint64(1)
			if gotUint != wantUint {
				t.Errorf("Test %d fail, Unexpected value for field \"field_02\":\nGot:  %d\nWant: %d", i, gotUint, wantUint)
			}
			gotUint = tag.GetUint("field_01")
			if gotUint != 0 {
				t.Errorf("Test %d fail, Unexpected value for field \"field_02\":\nGot:  %d\nWant: 0", i, gotUint)
			}
			gotSliceUint := tag.GetSliceUint("field_05", ' ')
			wantSliceUint := []uint64{1, 2, 3, 4}
			if !_sameUintSlice(gotSliceUint, wantSliceUint) {
				t.Errorf("Test %d fail, Unexpected value for field \"field_05\":\nGot:  %v\nWant: %v", i, gotSliceUint, wantSliceUint)
			}
			gotSliceUint = tag.GetSliceUint("field_99", ' ')
			if len(gotSliceInt) > 0 {
				t.Errorf("Test %d fail, Unexpected value for field \"field_99\":\nGot:  %v\nWant: []", i, gotSliceUint)
			}

			// Float tests
			gotFloat := tag.GetFloat("field_04")
			wantFloat := float64(123.321)
			if gotFloat != wantFloat {
				t.Errorf("Test %d fail, Unexpected value for field \"field_02\":\nGot:  %f\nWant: %f", i, gotFloat, gotFloat)
			}
			gotFloat = tag.GetFloat("field_01")
			if gotFloat != 0 {
				t.Errorf("Test %d fail, Unexpected value for field \"field_02\":\nGot:  %f\nWant: 0.0", i, gotFloat)
			}
		}
	}
}
