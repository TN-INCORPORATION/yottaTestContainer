package dataObject

import (
	"cErr"
	"fmt"
	"lang.yottadb.com/go/yottadb"
	"strconv"
)

const PicGlobal = "^ZPIC"

type Pic struct {
	ID        int
	Data 	  string
}

func (p *Pic) Save(tptoken uint64, y YottaCon) cErr.Error {

	// it should have code to check mandatory field here

	// insert
	err := setPicKey(tptoken, p.ID, &y.Key, &y.ErrStr)
	ctErr := HandleYTError(err)
	if ctErr != nil {
		return ctErr
	}
	stringData := p.getTextData()
	err = y.Data.SetValStr(tptoken, &y.ErrStr, stringData)
	ctErr = HandleYTError(err)
	if ctErr != nil {
		return ctErr
	}
	err = y.Key.SetValST(tptoken, &y.ErrStr, &y.Data)
	ctErr = HandleYTError(err)
	if ctErr != nil {
		return ctErr
	}

	fmt.Println("Save Picture : ", p.ID)
	return nil
}

func (p *Pic) getTextData() string {
	// return strconv.Itoa(c.ID) + "|" + c.Firstname + "|" + c.Lastname + "|" + c.Nickname
	return p.Data
}

func setPicKey(tptoken uint64, id int, key *yottadb.KeyT, errStr *yottadb.BufferT) error {
	var err error
	err = key.Varnm.SetValStr(tptoken, errStr, PicGlobal)
	if err != nil {
		return err
	}
	err = key.Subary.SetElemUsed(tptoken, errStr, 1)
	if err != nil {
		return err
	}
	err = key.Subary.SetValStr(tptoken, errStr, 0, strconv.Itoa(id))
	if err != nil {
		return err
	}
	return nil
}
