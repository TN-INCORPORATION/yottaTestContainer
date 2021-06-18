// Package provides core error handling capabilities.
package cErr

import (
	"fmt"
	"strconv"
	"strings"
)

const (
	// SystemErrorCat represents an internal system error
	SystemErrorCat int8 = iota
	// BadRequestErrorCat represents a server was unable to process
	BadRequestErrorCat int8 = iota
)

//event error category
const (
	//ok defines a non error state , error cat rage form 1000++
	OK int8 = iota + 100
	//retry means the event will be execute again
	Retry
	//skip means the event ignore error and wil be process next record.
	Skip
	//Unprocessable means the message cant be processed by the system in its current state ,Event Will be stop
	Unprocessable
)
const (
	InsertDuplicateKey = "duplicate key"
	QueryNotFound      = "data not found"
	UpdateNotFound     = "cannot found Record to update"
)


//Error Type
const (
	TypeInvalidInput    = "INVLDINPUT"
	TypeDataSubmitted   = "SUBMITTED"
	TypeOther           = "OTHER"
	TypeCommon          = "COMMON"
	TypeDataUpdate      = "DATAUPDATE"
	TypeInquiryResponse = "INQUIRYRESPONSE"
)

// coreError is a custom error which store additional information than the
// standard error
type coreError struct {
	// code is unique error number
	code int32
	// category specifies the category of known errors.  It is used
	// to map to errors at other layers
	category int8
	// desc is a description of the error
	desc string
	//causedBy is the root error if error is being wrapped
	causedBy error
	// yotta error code
	yottaErrCode int32
}

// Error is an interfacw with addition domain error information
type Error interface {
	error
	// Code returns the error code
	Code() int32
	// Category returns error category
	Category() int8
	// Desc returns the description of the error
	Desc() string
	// CausedBy returns the root error for wrapped errors
	CausedBy() error
	// IsInsertDuplicateKey return flag of error in case insert duplicate key
	IsInsertDuplicateKey() bool
	// IsQueryNotFound return flag of error in case query not found
	IsQueryNotFound() bool
	// String Concate Detail
	String() string
	// yotta error code
	YottaErrCode() int32
}

// New creates a new Error object with code, category and description
func New(code int32, cat int8, desc string) Error {
	return &coreError{code: code, category: cat, desc: desc}
}

// New creates a new Error object with only code
func NewByErrorCode(code int32,params ...string) Error {
	errMap := GetErrMap(code,params...)
	return &coreError{code: code, category: errMap.Category, desc: errMap.Description}
}

func WrapDErrorByCode(err error, code int32, params ...string) Error {
	errMap := GetErrMap(code, params...)

	return Wrap(err, errMap.Code, errMap.Category, errMap.Description)
}

// Wrap creates a new error while keeping reference to the caused by error
func Wrap(err error, code int32, cat int8, desc string) Error {
	return &coreError{code: code, category: cat, desc: desc, causedBy: err}
}

// Wrap creates a new error while keeping reference to the caused by error - for yottaDB
func WrapYT(err error, code int32, cat int8, desc string, yottaErrCode int32) Error {
	return &coreError{code: code, category: cat, desc: desc, causedBy: err, yottaErrCode: yottaErrCode}
}

// Error implements the error
func (e *coreError) Error() string {
	//append cause if present
	cause := ""
	if e.causedBy != nil {
		cause = fmt.Sprintf(", Caused by: %s", e.causedBy.Error())
	}
	return fmt.Sprintf("%s%s", e.desc, cause)
}

// Code returns the error code, satisfying the Error interface
func (e *coreError) Code() int32 {
	return e.code
}

// Category returns the error category
func (e *coreError) Category() int8 {
	return e.category
}

// Desc returns the error description
func (e *coreError) Desc() string {
	return e.desc
}

// CausedBy returns the wrapped error if there is one
func (e *coreError) CausedBy() error {
	return e.causedBy
}

// IsInsertDuplicateKey return flag of error in case insert duplicate key
// * WARNING * Need to handle case error != nil before use this
func (e *coreError) IsInsertDuplicateKey() bool {
	if e==nil {
		panic("Need to handle case error != nil before use this")
	}
	if strings.Contains(e.Error(), InsertDuplicateKey) {
		return true
	}
	return false
}

// IsQueryNotFound return flag of error in case query not found
// * WARNING * Need to handle case error != nil before use this
func (e *coreError) IsQueryNotFound() bool {
	if e==nil {
		panic("Need to handle case error != nil before use this")
	}
	if strings.Contains(e.Error(), QueryNotFound) {
		return true
	}
	return false
}

// return Yotta Error Code
func (e *coreError) YottaErrCode() int32 {
	return e.yottaErrCode
}

// String return String Error
func (e *coreError) String() string {

	return "Error Code:"+strconv.Itoa(int(e.code))+
		" ,Error Desc:"+e.desc+
		" ,Error Cat:"+strconv.Itoa(int(e.category))+
		" ,Error Cause by:"+e.causedBy.Error()

}
