package dto

// SortData ...
type SortData struct {
	Item  string
	Order string
}

// RangeData ...
type RangeData struct {
	From int
	To   int
}

// FilterData ...
type FilterData struct {
	Item  string
	Value string
}
