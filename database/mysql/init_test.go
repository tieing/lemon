package mysql

import (
	"fmt"
	"testing"
)

func TestAccDB(t *testing.T) {
	rs := &Rs{Name: "1111", R: &Rs{Name: "1"}}
	fmt.Println(rs.Name)
	A(rs)
	fmt.Println(rs.Name, rs.R.Name)

}

type Rs struct {
	Name string
	R    *Rs
}

func A(rs *Rs) {

	locaRs := &Rs{
		Name: "2222",
		R:    &Rs{Name: "2"},
	}

	*rs = *locaRs

}
