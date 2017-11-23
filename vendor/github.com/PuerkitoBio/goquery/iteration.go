package goquery

// Each iterates over a Selection object, executing a function for each
// matched element. It returns the current Selection object. The function
// f is called for each element in the selection with the index of the
// element in that selection starting at 0, and a *Selection that contains
// only that element.
func (s *Selection) Each(f func(int, *Selection)) *Selection ***REMOVED***
	for i, n := range s.Nodes ***REMOVED***
		f(i, newSingleSelection(n, s.document))
	***REMOVED***
	return s
***REMOVED***

// EachWithBreak iterates over a Selection object, executing a function for each
// matched element. It is identical to Each except that it is possible to break
// out of the loop by returning false in the callback function. It returns the
// current Selection object.
func (s *Selection) EachWithBreak(f func(int, *Selection) bool) *Selection ***REMOVED***
	for i, n := range s.Nodes ***REMOVED***
		if !f(i, newSingleSelection(n, s.document)) ***REMOVED***
			return s
		***REMOVED***
	***REMOVED***
	return s
***REMOVED***

// Map passes each element in the current matched set through a function,
// producing a slice of string holding the returned values. The function
// f is called for each element in the selection with the index of the
// element in that selection starting at 0, and a *Selection that contains
// only that element.
func (s *Selection) Map(f func(int, *Selection) string) (result []string) ***REMOVED***
	for i, n := range s.Nodes ***REMOVED***
		result = append(result, f(i, newSingleSelection(n, s.document)))
	***REMOVED***

	return result
***REMOVED***
