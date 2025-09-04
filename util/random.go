package util

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

const alphabet  = "abcdefghijklmnopqrstuvwxyz"


var seededRand = rand.New(rand.NewSource(time.Now().UnixNano()))


func RandomInt(min,max int64)int64{
	return min + seededRand.Int63n(max-min+1) //interger btween min and max
}


func RandomString(n int)string {
	var sb strings.Builder
	k := len(alphabet)

	for i:= 0; i <n ; i++ {
		c:=alphabet[seededRand.Intn(k)]
		sb.WriteByte(c)
	}
	return sb.String()
}


func RandomOwner ()string{
	return RandomString(10)
}

func RandomMoney () int64{
	return RandomInt(0,5000)
}

func RandomCurrency ()string{
	currencies := []string{EUR,USD,CAD}
	n := len(currencies)
	return  currencies[seededRand.Intn(n)]
}

func RandomEmail() string {
	return fmt.Sprintf("%s@email.com", RandomString(6))
}