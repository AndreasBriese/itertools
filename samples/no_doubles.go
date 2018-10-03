// no_doubles.go
//

package main

import (
	"fmt"
	iter "github.com/AndreasBriese/itertools"
	// "sync"
	// "time"
)

const SAMPLELEN = 12

var emitter = make(chan string)

func main() {
	// initialize
	s1 := []string{"Anna", "Anna", "Andy", "Tony", "Tony", "Anna", "Susi", "Susi", "Emi", "Roy", "Anna", "Anna"}

	condFn := func(prev, actual interface{}) bool {
		return (prev.(string) != actual.(string))
	}

	fmt.Println("original:", s1)
	// ~~ no additional memory
	fmt.Print("no doubles:  ")
	seq := iter.ToIterIf(s1)
	for step, v, ex := seq.DoubleCompNext(condFn), interface{}(""), false; ; {
		v, ex = step()
		if ex {
			break
		}
		fmt.Printf("%v ", v)
	}
	fmt.Println()

	// this uses additional memory
	fmt.Println("no doubles:", iter.ToIterIf(s1).DoubleComp(condFn).ToList())

	// this uses additional memory
	fmt.Println("doubles:", iter.ToIterIf(s1).DoubleComp(func(prev, actual interface{}) bool {
		return (prev.(string) == actual.(string))
	}).ToList())

	// s2 := iter.ToIterString([]string{"Anna", "Anna", "Andy", "Tony", "Tony", "Anna", "Susi", "Susi", "Emi", "Roy", "Anna", "Anna"})
	// for step, v, ex := s2.DoubleCompNext(func(prev, actual string) bool {
	// 	return (prev != actual)
	// }), "", false; ; {
	// 	v, ex = step()
	// 	if ex {
	// 		break
	// 	}
	// 	fmt.Printf("%v ", v)
	// }
	// fmt.Println()

	s1 = []string{"A", "A", "B", "T", "T", "A", "Su", "Su", "E", "Roy", "A", "A"}

	fmt.Println("original:", s1)
	// ~~ no additional memory
	fmt.Print("no doubles:  ")
	seq = iter.ToIterIf(s1)
	for step, v, ex := seq.DoubleCompNext(condFn), interface{}(""), false; ; {
		v, ex = step()
		if ex {
			break
		}
		fmt.Printf("%v ", v)
	}
	fmt.Println()

	// this uses additional memory
	fmt.Println("no doubles:", iter.ToIterIf(s1).DoubleComp(condFn).ToList())

	// this uses additional memory
	fmt.Println("doubles:", iter.ToIterIf(s1).DoubleComp(func(prev, actual interface{}) bool {
		return (prev.(string) == actual.(string))
	}).ToList())

	// s2 = iter.ToIterString([]string{"Anna", "Andy", "Tony", "Tony", "Anna", "Susi", "Susi", "Emi", "Roy", "Anna", "Anna"})
	// for step, v, ex := s2.DoubleCompNext(func(prev, actual string) bool {
	// 	return (prev != actual)
	// }), "", false; ; {
	// 	v, ex = step()
	// 	if ex {
	// 		break
	// 	}
	// 	fmt.Printf("%v ", v)
	// }
	// fmt.Println()

}
