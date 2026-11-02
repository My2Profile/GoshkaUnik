package main

type Counter struct{ n int }

func (c Counter) Get() int { return c.n }
func (c *Counter) Inc()    { c.n++ }
func NewCounter() Counter  { return Counter{} }
func main() {
	var c Counter
	c.Inc()         // 1 compiles
	(&c).Inc()      // 2 compiles
	Counter{}.Get() // compiles
	// Counter{}.Inc() // compiler error
	// NewCounter().Inc() // 5
	m := map[string]Counter{"a": {}}
	m["a"].Get() // 6
	//	m["a"].Inc() // 7
	s := []Counter{{}}
	s[0].Inc() // 8
	pm := map[string]*Counter{"a": {}}
	pm["a"].Inc()
}
