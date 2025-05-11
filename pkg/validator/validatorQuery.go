package validator

type Query struct {
	Page     int
	PageSize int
	Search   string
	Filter   string
	Filters  []string
}

func ValidatorQuery(v *Validator, q Query) {
	v.Check(q.Page >= 0, "page", "must be positive integer")
	v.Check(q.Page < 1000, "page", "must be less than 1000")
	v.Check(q.PageSize > 0, "page size", "must be greater than 0")
	v.Check(q.PageSize < 100, "page size", "must be less than 100")
	v.Check(In(q.Filter, append(q.Filters, "")...), "filter", "invalid filter value")
}

func (f Query) Limit() int {
	return f.PageSize
}

func (f Query) Offset() int {
	return (f.Page - 1) * f.PageSize
}
