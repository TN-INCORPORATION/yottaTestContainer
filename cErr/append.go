package cErr

import (
	"fmt"
	"strings"
)

//ErrorList is a list implementation of a domain error
type ErrorList []Error

func (err ErrorList) Error() string {

	if len(err) == 1 {
		return fmt.Sprintf("%s", err[0])
	}

	errs := make([]string, len(err))
	for i, e := range err {
		errs[i] = fmt.Sprintf("%s\n", e)
	}

	return fmt.Sprintf("%d errors: \n\t%s", len(errs), strings.Join(errs, "\n\t"))
}

// Code returns the error code of ErrorList
func (err ErrorList) Code() int32 {
	if len(err) > 0 {
		return err[0].Code()
	}
	return 0
}

// Category returns error category of ErrorList
func (err ErrorList) Category() int8 {
	if len(err) > 0 {
		return err[0].Category()
	}
	return 0
}

//Desc returns the description of the ErrorList
func (err ErrorList) Desc() string {
	return err.Error()
}

//CausedBy returns the root error for wrapped errors
func (err ErrorList) CausedBy() error {
	if len(err) > 0 {
		return err[0].CausedBy()
	}
	return nil
}

func (err ErrorList) IsInsertDuplicateKey() bool {
	panic("not support")
	return false
}

func (err ErrorList) IsQueryNotFound() bool {
	panic("not support")
	return false
}

func (err ErrorList) YottaErrCode() int32 {
	panic("not support")
	return 0
}

func (err ErrorList) String() string {
	panic("not support")
	return ""
}

//Append will handle multi error
func Append(err Error, errs ...Error) Error {

	//check the type of the incomeing error
	switch err := err.(type) {

	case ErrorList:
		if err == nil {
			err = ErrorList{}
		}

		for _, e := range errs {
			switch e := e.(type) {
			case ErrorList:
				if e != nil {
					err = append(err, e...)
				}
			default:
				if e != nil {
					err = append(err, e)
				}
			}
		}
		if len(err) == 0 {
			return nil
		}
		return err
	default:
		errList := make(ErrorList, 0, len(errs)+1)
		if err != nil {
			errList = append(errList, err)
		}

		errList = append(errList, errs...)
		return Append(ErrorList{}, errList...)
	}

}
