package excel

import (
	"bytes"
	"fmt"
	"github.com/mutezebra/ClassroomRandomRollCallSystem/pkg/errno"
	"github.com/pkg/errors"
	"github.com/xuri/excelize/v2"
)

type ImportUser struct {
	StudentNumber string
	Name          string
	PhoneNumber   string
}

type ExportUser struct {
	UID      int64
	Password string
}

func ReadExcelToUsers(data []byte) ([]*ImportUser, error) {
	buffer := bytes.NewBuffer(data)
	f, err := excelize.OpenReader(buffer)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Failed when read file as excel,err: %v", err))
	}

	cols, err := f.GetCols("Sheet1")
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Failed when get cols from excel file,err: %v", err))
	}

	if len(cols) != 3 || len(cols[0]) == 0 || len(cols[1]) == 0 {
		return nil, errno.New(errno.WrongExcelFormat, "wrong excel format")
	}

	stuN := cols[0]
	names := cols[1]
	phoneNumbers := cols[2]
	if stuN[0] != "学号" || names[0] != "姓名" || phoneNumbers[0] != "手机号" {
		return nil, errno.New(errno.WrongExcelFormat, "wrong excel format")
	}

	users := make([]*ImportUser, len(stuN)-1)
	for i := 1; i < len(stuN); i++ {
		users[i-1] = &ImportUser{
			StudentNumber: stuN[i],
			Name:          names[i],
			PhoneNumber:   phoneNumbers[i],
		}
	}
	return users, nil
}
