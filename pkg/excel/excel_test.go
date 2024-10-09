package excel

import (
	"os"
	"testing"
)

func TestReadExcelToUsers(t *testing.T) {
	buf, err := os.ReadFile("./软工学生名单.xlsx")
	if err != nil {
		t.Fatal(err)
	}

	us, err := ReadExcelToUsers(buf)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%v", *us[0])
}
