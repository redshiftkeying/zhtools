package main

import "github.com/nats-io/nats.go"

func main() {
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		panic(err)
	}
	defer nc.Close()

}

// var t0 uint64 = 0
// var timeStep uint64 = 30
// var currentUnixTime uint64 = uint64(time.Now().Unix())
// var key string = "7q66vogpkka6tln7qmr32beac4vhqepc" // sample secret key
//
// t := calculateT(currentUnixTime, t0, timeStep)
// hash := computeHmacSha1(key, t)
// otpCode := truncate(hash, 6)
//
// fmt.Println(otpCode)
