/*
Package utils implements utility methods for the program.
*/
package utils

func Check(e error) {
	if e != nil {
		panic(e)
	}
}
