package analytics

import (
	"log"
	"time"
)

// Track calculates the time taken for a function to run
//
// For example:
// 	func doSomething() {
// 		defer Track("doSomething", time.Now())
// 	}
func Track(name string, start time.Time, params ...interface{}) {
	elapsed := time.Since(start)
	log.Printf("%s took %s\n", name, elapsed)
	log.Println(params)
}
