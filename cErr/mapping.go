package cErr

import (
	"strconv"
	"strings"
)



type ErrorDetail struct {
	Code          int32 `json:"code"`
	Category      int8 `json:"category"`
	Description   string `json:"description"`
	ErrorType     string `json:"error_type"`
}


var errorMap = make(map[int32]ErrorDetail)

func InitErrorMapping(errorDetails []ErrorDetail) {

	for _,detail:= range  errorDetails {
		errorMap[detail.Code]=detail
	}

}

func GetErrMap(code int32, params ...string) ErrorDetail {
	v, ok := errorMap[code]
	if !ok { // default mapping
		return ErrorDetail{
			Code:        code,
			Category:    BadRequestErrorCat,
			ErrorType:   TypeOther,
			Description: replaceDynamicValue("[errorMap]Unknown Mapping : ~p0", params...),
		}
	}

	// Default ErrorType
	if v.ErrorType == "" {
		v.ErrorType = TypeOther
	}
	v.Description = replaceDynamicValue(v.Description, params...)

	return v
}

func replaceDynamicValue(des string, p ...string) string {

	output := des
	i := 0
	for i < 20 && strings.Index(output, "~p") > -1 {
		if i < len(p) && strings.Index(output, "~p"+strconv.Itoa(i)) > -1 {
			output = strings.Replace(output, "~p"+strconv.Itoa(i), p[i], -1)
		} else {
			output = strings.Replace(output, "~p"+strconv.Itoa(i), "", -1)
		}
		i++
	}

	return output
}

