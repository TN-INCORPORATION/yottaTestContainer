package dataObject

import (
	"cErr"
	"fmt"
	"os"
	"strconv"
	"time"
)

const NseqGlobal = "^ZNSEQ"

func GetNextSeq(tptoken uint64, y YottaCon) (int, cErr.Error) {
	var val1 string

	// next sequence
	ytErr := y.Key.Varnm.SetValStr(tptoken, &y.ErrStr, NseqGlobal)
	ctErr := HandleYTError(ytErr)
	if ctErr != nil {
		return 0, ctErr
	}
	ytErr = y.Key.Subary.SetElemUsed(tptoken, &y.ErrStr, 0)
	ctErr = HandleYTError(ytErr)
	if ctErr != nil {
		return 0, ctErr
	}

	// it should use IncrST, but to make it easier to conflict, seperate command instead
	ytErr = y.Key.ValST(tptoken, &y.ErrStr, &y.Data)
	ctErr = HandleYTError(ytErr)
	if ctErr != nil {
		return 0, ctErr
	}
	val1, ytErr = y.Data.ValStr(tptoken, &y.ErrStr)
	ctErr = HandleYTError(ytErr)
	if ctErr != nil {
		return 0, ctErr
	}

	if val1 == "" {
		val1 = "0"
	}
	id, _ := strconv.Atoi(val1) // ignore error handling to make it easy
	id++

	// delay
	if delay, err := strconv.Atoi(os.Getenv("DELAY")); err == nil {
		fmt.Printf("Delay %d\n", delay)
		time.Sleep(time.Duration(delay) * time.Millisecond)
	}
	// update sequence back
	ytErr = y.Data.SetValStr(tptoken, &y.ErrStr, strconv.Itoa(id))
	ctErr = HandleYTError(ytErr)
	if ctErr != nil {
		return 0, ctErr
	}
	ytErr = y.Key.SetValST(tptoken, &y.ErrStr, &y.Data)
	ctErr = HandleYTError(ytErr)
	if ctErr != nil {
		return 0, ctErr
	}

	return id, nil
}
