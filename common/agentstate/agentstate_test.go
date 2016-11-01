package agentstate

import (
	//"fmt"
	"flag"
	"os"
	"testing"
)

var TS *State

func TestMain(M *testing.M) {
	flag.Parse()
	TS = NewState()
	os.Exit(M.Run())
}

func TestAdd(t *testing.T) {

	rc := TS.Add("S1", "T1", 1)
	if !rc {
		t.Fail()
	}
	results := TS.ListDistribution("T1")

	if len(results) != 1 {
		t.Fail()
	}

	if results[0].NodeName != "S1" {
		t.Fail()
	}

}

func TestDelInvalid(t *testing.T) {

	rc := TS.Del("S10", "T1")

	if rc {
		t.Fail()
	}
}

func TestDelValid(t *testing.T) {

	rc := TS.Del("S1", "T1")

	if !rc {
		t.Fail()
	}
}

func TestCanfitInvalid(t *testing.T) {
	TS.Add("S1", "T1", 2)

	if TS.Canfit("S1", "T1", 2) {
		//If it returns true then throw an error
		t.Fail()
	}

}

func TestCanfitValid(t *testing.T) {
	if !TS.Canfit("S1", "T2", 2) {
		//If it returns true then throw an error
		t.Fail()
	}
}
