package main

import (
	"fmt"
	"testing"
	"github.com/hyperledger/fabric/core/chaincode/shim"
)

func checkInit(t *testing.T, stub *shim.MockStub, args [][]byte) {
	res := stub.MockInit("1", args)
	if res.Status != shim.OK {
		fmt.Println("Init failed", string(res.Message))
		t.FailNow()
	}
}

func checkState(t *testing.T, stub *shim.MockStub, name string, a Auction) {
	bytes := stub.State[name]
	if bytes == nil {
		fmt.Println("State", name, "failed to get value")
		t.FailNow()
	}
	if string(bytes) != string(a.auction2Bytes()) {
		fmt.Println("State value", name, "was not", a, "as expected")
		t.FailNow()
	}
}

func checkQuery(t *testing.T, stub *shim.MockStub, name string, a Auction) {
	res := stub.MockInvoke("1", [][]byte{[]byte("query"), []byte(name)})
	if res.Status != shim.OK {
		fmt.Println("Query", name, "failed", string(res.Message))
		t.FailNow()
	}
	if res.Payload == nil {
		fmt.Println("Query", name, "failed to get value")
		t.FailNow()
	}
	/*if string(res.Payload) != string(a.auction2Bytes()) {
		fmt.Println("Query value", name, "was not", a, "as expected")
		t.FailNow()*/
	a.bytes2Auction(res.Payload)
	fmt.Println("Query value", name,"Auction",a)

	}


func checkInvoke(t *testing.T, stub *shim.MockStub, args [][]byte) {
	res := stub.MockInvoke("1", args)
	if res.Status != shim.OK {
		fmt.Println("Invoke", args, "failed", string(res.Message))
		t.FailNow()
	}
}


func TestInit(t *testing.T) {
	scc := new(Auction)
	stub := shim.NewMockStub("ex02", scc)

	// Init A=123 B=234
	checkInit(t, stub, [][]byte{[]byte("init"), []byte("CAT"), []byte("999999")})

	c := Auction{
		Base_price:999999,
	}
	checkState(t, stub, "A", c)
}

func TestQuery(t *testing.T) {
	scc := new(Auction)
	stub := shim.NewMockStub("ex02", scc)

	// Init A=345 B=456
	checkInit(t, stub, [][]byte{[]byte("init"), []byte("CAT"), []byte("999999")})

	// Query A
	c := Auction{
		Base_price:999999,
	}

	checkQuery(t, stub, "CAT", c)

}

func TestInvoke(t *testing.T) {
	scc := new(Auction)
	stub := shim.NewMockStub("ex02", scc)

	// Init A=567 B=678
	checkInit(t, stub, [][]byte{[]byte("init"), []byte("CAT"), []byte("999999")})

	// Invoke A->B for 123
	checkInvoke(t, stub, [][]byte{[]byte("invoke"), []byte("CAT"), []byte("fly"), []byte("9999991")})
	name,_:= String2Uint64("fly")
	s := Auction{
		Base_price:999999,
		Heightest_p:9999991,
		Heightest_n:name,
	}
	checkQuery(t, stub, "CAT", s)

	checkInvoke(t, stub, [][]byte{[]byte("invoke"), []byte("CAT"), []byte("pig"), []byte("9999990")})
	name1,_ := String2Uint64("pig")
	sc := Auction{
		Base_price:999999,
		Heightest_p:9999991,
		Heightest_n:name,
		Second_n:name1,
		Second_p:9999990,
	}
	checkQuery(t, stub, "CAT", sc)
}