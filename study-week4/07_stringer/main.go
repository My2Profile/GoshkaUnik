// STUDY FILE (analog of Task 7: Tour of Go "Exercise: Stringers").
// Domain: Version [3]int instead of IPAddr [4]byte. Run: go run ./07_stringer
//
// BIG PICTURE
//   - fmt.Stringer is an interface: type Stringer interface{ String() string }
//   - fmt.Println / Printf("%v") check whether a value has a String() method
//     and, if so, print that instead of the raw representation.
//   - Methods can be declared on ANY named type defined in your package,
//     not only structs (here: an array type).
package main

import "fmt"

type Version [3]int

func (v Version) String() string {
	return fmt.Sprintf("v%d.%d.%d", v[0], v[1], v[2])
}

// Gotcha: String() on a POINTER receiver is only in the method set of *Tag.
type Tag struct{ Label string }

func (t *Tag) String() string { return "#" + t.Label }

func main() {
	v := Version{1, 4, 2}
	fmt.Println(v) // v1.4.2  (not [1 4 2])

	releases := map[string]Version{"stable": {1, 4, 2}, "beta": {2, 0, 0}}
	for name, ver := range releases { // map order is random; Println still uses String()
		fmt.Println(name, ver)
	}

	t := Tag{"go"}
	fmt.Println(t)  // {go}   : Tag value has no String() in its method set
	fmt.Println(&t) // #go    : *Tag does
}
