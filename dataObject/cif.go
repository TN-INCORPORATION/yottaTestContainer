package dataObject

import (
	"cErr"
	"fmt"
	"lang.yottadb.com/go/yottadb"
	"strconv"
	"strings"
)

const CifGlobal = "^ZCIF"

type CIF struct {
	ID        int    `json:"id"`
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
	Nickname  string `json:"nickname,omitempty"`
	Salary    int    `json:"salary,omitempty"`
}

func (c *CIF) Load(tptoken uint64, y YottaCon, id int) cErr.Error {

	var val1 string
	var err error

	c.ID = id
	err = setAccountKey(tptoken, c.ID, &y.Key, &y.ErrStr)
	ctErr := HandleYTError(err)
	if ctErr != nil {
		return ctErr
	}

	err = y.Key.ValST(tptoken, &y.ErrStr, &y.Data)
	ctErr = HandleYTError(err)
	if ctErr != nil {
		return ctErr
	}
	val1, err = y.Data.ValStr(tptoken, &y.ErrStr)
	ctErr = HandleYTError(err)
	if ctErr != nil {
		return ctErr
	}
	values := strings.Split(val1, "|")

	c.ID = id
	c.Firstname = values[0]
	c.Lastname = values[1]
	c.Nickname = values[2]
	c.Salary, _ = strconv.Atoi(values[3]) // ignore error to make it easier

	fmt.Println("Load : ", c.ID, c.Firstname, c.Lastname, c.Nickname, c.Salary)
	return nil
}

func (c *CIF) Save(tptoken uint64, y YottaCon) cErr.Error {

	// it should have code to check mandatory field here

	// insert
	err := setAccountKey(tptoken, c.ID, &y.Key, &y.ErrStr)
	ctErr := HandleYTError(err)
	if ctErr != nil {
		return ctErr
	}
	stringData := c.getTextData()
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

	fmt.Println("Save : ", c.ID, c.Firstname, c.Lastname, c.Nickname, c.Salary)
	return nil
}

func (c *CIF) getTextData() string {
	// return strconv.Itoa(c.ID) + "|" + c.Firstname + "|" + c.Lastname + "|" + c.Nickname
	return c.Firstname + "|" + c.Lastname + "|" + c.Nickname + "|" + strconv.Itoa(c.Salary)
}

func setAccountKey(tptoken uint64, id int, key *yottadb.KeyT, errStr *yottadb.BufferT) error {
	var err error
	err = key.Varnm.SetValStr(tptoken, errStr, CifGlobal)
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
