package main

import (
	"errors"
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
	"math"
	"reflect"
	"strconv"
	"unsafe"
)

// SimpleChaincode example simple Chaincode implementation
type Auction struct {
	Base_price  int
	Heightest_n uint64
	Heightest_p int
	Second_n    uint64
	Second_p    int
}

func (a *Auction) Init(stub shim.ChaincodeStubInterface) pb.Response {
	fmt.Println("auction Init")
	_, args := stub.GetFunctionAndParameters()
	var Name string // Entities
	var err error

	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 4")
	}

	// Initialize the chaincode
	Name = args[0]

	base_price, err := strconv.Atoi(args[1])
	if err != nil {
		return shim.Error("Expecting integer value for asset holding")
	}
	fmt.Printf("Auction = %s, base_price = %d\n", Name, base_price)

	a.Base_price = base_price
	b := a.auction2Bytes()

	println("Bytes:",b)

	// Write the state to the ledger
	err = stub.PutState(Name,b)
	if err != nil {
		fmt.Println("Fail to put the data")
		return shim.Error(err.Error())
	} else {
		fmt.Println("Init successfully!")
	}
	return shim.Success(nil)
}

func (a *Auction) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	fmt.Println("ex02 Invoke")
	function, args := stub.GetFunctionAndParameters()
	if function == "invoke" {
		// Make payment of X units from A to B
		return a.invoke(stub, args)
	} else if function == "query" {
		// Deletes an entity from its state
		return a.query(stub, args)
	}

	return shim.Error("Invalid invoke function name. Expecting \"invoke\" \"delete\" \"query\"")
}

// Transaction makes payment of X units from A to B
func (a *Auction) invoke(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var AuctionName string
	var Name uint64
	var price int
	var err error

	if len(args) != 3 {
		return shim.Error("Incorrect number of arguments. Expecting 3")
	}

	AuctionName = args[0]
	Name,err = String2Uint64(args[1])
	price, err = strconv.Atoi(args[2])
	if err != nil {
		return shim.Error("Expecting integer value for asset holding")
	}

	Avalbytes, err := stub.GetState(AuctionName)
	if err != nil {
		return shim.Error("Failed to get state")
	}
	if Avalbytes == nil {
		return shim.Error("Entity not found")
	}
	a.bytes2Auction(Avalbytes)

	if a.Base_price > price {
		return shim.Error("Your price must heigher than " + strconv.Itoa(price))
	} else if price > a.Second_p && price < a.Heightest_p {
		a.Second_n = Name
		a.Second_p = price
	} else if price > a.Heightest_p{
		a.Second_n = a.Heightest_n
		a.Second_p= a.Heightest_p
		a.Heightest_n = Name
		a.Heightest_p = price
	}

	err = stub.PutState(AuctionName, a.auction2Bytes())
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(nil)
}

func (a *Auction) query(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	//TODO:Only the owner of the auction can open the result.

	var Auction string // Entities
	var err error

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting name of the person to query")
	}

	Auction = args[0]

	Avalbytes, err := stub.GetState(Auction)
	if err != nil {
		jsonResp := "{\"Error\":\"Failed to get state for " + Auction + "\"}"
		return shim.Error(jsonResp)
	}

	if Avalbytes == nil {
		jsonResp := "{\"Error\":\"Nil amount for " + Auction + "\"}"
		return shim.Error(jsonResp)
	}

	a.bytes2Auction(Avalbytes)
	fmt.Printf("Query Response:\n")
	a.print(Auction)

	return shim.Success([]byte(Uint642String(a.Second_n)))
}

var sizeOfMyStruct = int(unsafe.Sizeof(Auction{}))

func (a *Auction) auction2Bytes() []byte {
	var x reflect.SliceHeader
	x.Len = sizeOfMyStruct
	x.Cap = sizeOfMyStruct
	x.Data = uintptr(unsafe.Pointer(a))
	return *(*[]byte)(unsafe.Pointer(&x))
}

func (a *Auction) bytes2Auction(b []byte) {
	a = (*Auction)(unsafe.Pointer(
		(*reflect.SliceHeader)(unsafe.Pointer(&b)).Data,
	))
}

func (a *Auction) print(name string){
	fmt.Printf("Auction = %s, base_price = %d\n",name,a.Base_price)
	fmt.Printf("Heightest participant = %s , price = %d\n",Uint642String(a.Heightest_n),a.Heightest_p)
	fmt.Printf("Second participant = %s , price = %d\n",Uint642String(a.Second_n),a.Second_p)
}

func String2Uint64(str string) (uint64, error) {
	var value uint64 = 0

	if len(str) > 13 {
		return 0, errors.New("string is too long to be a valid name")
	}

	if len(str) == 0 {
		return 0, nil
	}

	var n = int(math.Min(float64(len(str)), 12.0))

	for i := 0; i < n; i++ {
		value <<= 5
		v, err := Char2Value(str[i])
		if err != nil {
			return 0, err
		}
		value |= uint64(v)
	}
	value <<= uint(4 + 5*(12-n))
	if len(str) == 13 {
		v, err := Char2Value(str[12])
		if err != nil {
			return 0, err
		}
		if v > 0x0F {
			return 0, errors.New("thirteenth character in name cannot be a letter that comes after j")
		}
		value |= uint64(v)
	}
	return value, nil
}

func Char2Value(c uint8) (uint8, error) {
	if c == '.' {
		return 0, nil
	} else if c >= '1' && c <= '5' {
		return (c - '1') + 1, nil
	} else if c >= 'a' && c <= 'z' {
		return (c - 'a') + 6, nil
	}

	return 0, errors.New("character is not in allowed character set for names")
}

func Uint642String(c uint64) string {
	const charmap = ".12345abcdefghijklmnopqrstuvwxyz"
	var mask uint64 = 0xF800000000000000
	var str []uint8
	var i int

	for i = 0; i < 13; i++ {
		if c == 0 {
			return string(str)
		}

		if i == 12 {
			index := (c & mask) >> 60
			str = append(str, charmap[index])
		} else {
			index := (c & mask) >> 59
			str = append(str, charmap[index])
		}

		c <<= 5
	}
	return string(str)

}


func main() {
	err := shim.Start(new(Auction))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}