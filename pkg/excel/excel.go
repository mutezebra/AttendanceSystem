package excel

import (
	"bytes"
	"fmt"
	"github.com/mutezebra/ClassroomRandomRollCallSystem/pkg/errno"
	"github.com/pkg/errors"
	"github.com/xuri/excelize/v2"
)

type User struct {
	StudentNumber string
	Name          string
}

func ReadExcelToUsers(data []byte) ([]*User, error) {
	buffer := bytes.NewBuffer(data)
	f, err := excelize.OpenReader(buffer)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Failed when read file as excel,err: %v", err))
	}

	cols, err := f.GetCols("Sheet1")
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Failed when get cols from excel file,err: %v", err))
	}

	if len(cols) != 2 || len(cols[0]) == 0 || len(cols[1]) == 0 {
		return nil, errno.New(errno.WrongExcelFormat, "wrong excel format")
	}

	stuN := cols[0]
	names := cols[1]
	if stuN[0] != "学号" || names[0] != "姓名" {
		return nil, errno.New(errno.WrongExcelFormat, "wrong excel format")
	}

	users := make([]*User, len(stuN)-1)
	for i := 1; i < len(stuN); i++ {
		users[i-1] = &User{
			StudentNumber: stuN[i],
			Name:          names[i],
		}
	}
	return users, nil
}
