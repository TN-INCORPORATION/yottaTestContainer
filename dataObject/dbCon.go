package dataObject

import (
	"lang.yottadb.com/go/yottadb"
)

var tptoken uint64 = yottadb.NOTTP

type YottaCon struct {
	Key yottadb.KeyT
	Data,Incrval yottadb.BufferT
	ErrStr yottadb.BufferT
}

func (y *YottaCon) Alloc(keyVarSiz,keyNumSubs,keySubSiz,dataSiz,incSiz,errSiz uint32) {
	y.Key.Alloc(keyVarSiz, keyNumSubs, keySubSiz)
	y.Data.Alloc(dataSiz)
	y.ErrStr.Alloc(errSiz)
	if incSiz > 0 {
		y.Incrval.Alloc(incSiz)
	}
}

func (y *YottaCon) Free() {
	y.Key.Free()
	y.Data.Free()
	y.ErrStr.Free()
	y.Incrval.Free()
}

func GetCon(keyVarSiz,keyNumSubs,keySubSiz,dataSiz,incSiz,errSiz uint32) YottaCon {
	var dbCon YottaCon
	dbCon.Alloc(keyVarSiz, keyNumSubs,keySubSiz, dataSiz, incSiz,yottadb.YDB_MAX_ERRORMSG)
	return dbCon
}