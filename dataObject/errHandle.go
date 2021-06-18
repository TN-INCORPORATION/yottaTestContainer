package dataObject

import (
	"cErr"
	"fmt"
	"lang.yottadb.com/go/yottadb"
	"runtime"
)

func HandleYTError(err error) cErr.Error {
	if err == nil {
		return nil
	}
	if ydbErr, ok := err.(*yottadb.YDBError); ok {
		yottaErrCode := yottadb.ErrorCode(ydbErr)
		switch yottaErrCode {
		case yottadb.YDB_TP_RESTART:
			// If an application uses transactions, TP_RESTART must be handled inside the transaction callback;
			// it is here. For completeness, but ensure that one modifies this routine as needed, or copies bits
			// from it. A transaction must be restarted; this can happen if some other process modifies a value
			// we read before we commit the transaction.
			return cErr.WrapYT(err, 70000000, cErr.SystemErrorCat, "TP restart", int32(yottaErrCode))
		case yottadb.YDB_TP_ROLLBACK:
			// If an application uses transactions, TP_ROLLBACK must be handled inside the transaction callback;
			// it is here for completeness, but ensure that one modifies this routine as needed, or copies bits
			// from it. The transaction should be aborted; this can happen if a subtransaction return YDB_TP_ROLLBACK
			// This return will be a bit more situational.
			return cErr.WrapYT(err, 70000000, cErr.SystemErrorCat, "TP rollback", int32(yottaErrCode))
		case yottadb.YDB_ERR_CALLINAFTERXIT:
			// The database engines was told to close, yet we tried to perform an operation. Either reopen the
			// database, or exit the program. Since the behavior of this depends on how your program should behave,
			// it is commented out so that a panic is raised.
			return cErr.WrapYT(err, 70000000, cErr.SystemErrorCat, "Database is close", int32(yottaErrCode))
		case yottadb.YDB_ERR_NODEEND:
			// This should be detected seperately, and handled by the looping function; calling a more generic error
			// checker should be done to check for other errors that can be encountered.
			return cErr.WrapYT(err, 70000000, cErr.SystemErrorCat, "Node End", int32(yottaErrCode))
		default:
			var stdErr string
			_, file, line, ok := runtime.Caller(1)
			if ok {
				stdErr = fmt.Sprintf("Assertion failure in %v at line %v with error (%d): %v", file, line, yottadb.ErrorCode(err), err)
			} else {
				stdErr = fmt.Sprintf("Assertion failure (%d): %v", yottadb.ErrorCode(err), err)
			}
			return cErr.WrapYT(err, 70000000, cErr.SystemErrorCat, stdErr, int32(yottaErrCode))
		}
	} else {
		// todo : or should panic? think again when handle ytErrCd later
		return cErr.Wrap(err, 70000000, cErr.SystemErrorCat, "unidentified yottaDB error")
	}
}
