package dbUtil

import "C"
import (
	"cErr"
	"context"
	"errors"
	"fmt"
	"lang.yottadb.com/go/yottadb"
	"os"
)

func WrapTransaction(ctx context.Context, tpfn func(context.Context, uint64) cErr.Error) cErr.Error {
	var errStr yottadb.BufferT
	defer errStr.Free()
	errStr.Alloc(128)

	var custErr cErr.Error

	defer func() {
		if r := recover(); r != nil {
			fmt.Println("recovery from panic that rethrow from TpE - shut database down and close app")
			yottadb.Exit()
			os.Exit(1)
		}
	}()

	ytErr := yottadb.TpE(yottadb.NOTTP, &errStr, func(tptoken uint64, errstrp *yottadb.BufferT) (ret int32) {

		defer func() {
			if r := recover(); r != nil {
				errorStr := fmt.Sprintf("%v", r) // Seems to be only way to turn r back into a string
				howToHandle := RecoveryTypeCheck(errorStr)
				fmt.Println("In panic recover for", r, "- handle type", howToHandle)
				switch howToHandle {
				case FatalSignalPanicSync: // Synchronous panics don't create core files - create one
					// Turns out this core is mostly useless and generally only of a single
					// .. cgo thread that ran the call. But the process is here if we want
					// to disect it. A Go generated core has all the threads in it.
					// C.ydb_fork_n_core()	// we found it make system hang.
					fallthrough
				case FatalSignalPanicASync, FatalApplPanic:
					panic(errorStr)
				case NonFatalReturnAppl, NonFatalReturnNonAppl:
					// Note, if return YDB_OK here, any partial transaction udpates that have been made
					// will commit even through the transaction was cut short by the error leading to the
					// panic we are now looking at.
					custErr = cErr.WrapYT(errors.New(errorStr), 70000000, cErr.SystemErrorCat, "TP Rollback because panic", yottadb.YDB_TP_ROLLBACK)
					ret = yottadb.YDB_TP_ROLLBACK
				default:
					fmt.Println("Invalid howToHandle type:", howToHandle)
					os.Exit(-1)
				}
			}
		}()

		// set token
		ctx = context.WithValue(ctx, "YOTTATOKEN", tptoken)

		// execute function
		custErr = nil // for sure
		custErr = tpfn(ctx, tptoken)

		// calculate yotta error
		if custErr != nil {
			if custErr.YottaErrCode() == 0 {
				return yottadb.YDB_TP_ROLLBACK
			}
			return custErr.YottaErrCode() // must return this code, but set error object into context.
		}
		return yottadb.YDB_OK

	}, "", []string{})

	if ytErr != nil {
		fmt.Println("Err from TPE", ytErr) // <= possible ignore this one.
	}

	return custErr
}

// Whether to allow certain classes of panics to be treated as a panic or whether it should fail.
var PermissiveErrorMode bool = true

//var PermissiveErrorMode bool = false

// condition to indicates ydb error
const ydbFatalSignal string = "YDB: Fatal signal "
const ydbFatalError string = "YDB: "

// list of non fatal panic
var ydbNonFatalError = []string{
	"YDB: Name of routine to call cannot be null string",
}
var nonFatalError = []string{
	// "runtime error: index out of range",
	"runtime error: invalid memory address or nil pointer dereference",
}

// list of Text for fatal synchronous signal panics sent by Go.
// these panic will cause system to fork process before rethrow the panic
var fatalSignalText = [...]string{
	//"runtime error: invalid memory address or nil pointer dereference",
	//"runtime error: negative shift amount",
	//"runtime error: integer divide by zero",
	//"runtime error: integer overflow",
	//"runtime error: floating point error",
	//"fatal error: unexpected signal",
}

// Disposition types returned by recoverTypeCheck
const (
	FatalSignalPanicSync  = iota
	FatalSignalPanicASync = iota
	FatalApplPanic        = iota
	NonFatalReturnAppl    = iota
	NonFatalReturnNonAppl = iota
)

func RecoveryTypeCheck(panicRsn string) int {
	fmt.Println(panicRsn)
	// First separate into YDB panics and non-YDB panics
	if ydbFatalError == panicRsn[:len(ydbFatalError)] { // This is a YDB panic
		// See if this is a fatal signal - this should catch all of the asynchronous fatal signals
		if ydbFatalSignal == panicRsn[:len(ydbFatalSignal)] {
			return FatalSignalPanicASync
		}
		// Check for any application issue that we (for purposes of this example or system robustness) do not treat
		// as fatal if PermissiveErrorMode is enabled. Otherwise this error is a fatal application panic.
		if PermissiveErrorMode && isNonFatal(ydbNonFatalError, panicRsn) {
			return NonFatalReturnAppl
		}
		// Anthing left over is considered a fatal application panic
		return FatalApplPanic
	}
	// Check for our non-fatal application errors next again, only if we are operating in Permissive Mode. If there
	// are more than one we want to bypass, then their text strings should be in an  array much like the fatal signal
	// array is setup.
	if PermissiveErrorMode && isNonFatal(nonFatalError, panicRsn) {
		return NonFatalReturnNonAppl
	}
	// Look to see if we have one of the synchronous fatal signals. We have to check for each one of these.
	for _, val := range fatalSignalText {
		if (len(val) <= len(panicRsn)) && (val == panicRsn[:len(val)]) {
			return FatalSignalPanicSync
		}
	}
	// Anything left over at this point is a fatal application panic.
	return FatalApplPanic
}

func isNonFatal(list []string, panicRsn string) bool {
	for _, val := range list {
		if (len(val) <= len(panicRsn)) && (val == panicRsn[:len(val)]) {
			return true
		}
	}

	return false
}
