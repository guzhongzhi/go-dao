package dao

import "math"

func NewPagination() *Pagination {
	p := &Pagination{}
	p.SetPage(1).SetPageSize(20).SetTotal(0)
	return p
}

type Pagination struct {
	page     *int64
	pageSize *int64
	total    *int64
}

func (s *Pagination) SetPage(v int64) *Pagination {
	s.page = &v
	return s
}

func (s *Pagination) SetPageSize(v int64) *Pagination {
	s.pageSize = &v
	return s
}

func (s *Pagination) SetTotal(v int64) *Pagination {
	s.total = &v
	return s
}

func (s *Pagination) Page() int64 {
	if s.page == nil {
		return 1
	}
	return *s.page
}

func (s *Pagination) PageSize() int64 {
	if s.pageSize != nil {
		return *s.pageSize
	}
	return 20
}

func (s *Pagination) Total() int64 {
	if s.total == nil {
		return 0
	}
	return *s.total
}

func (s *Pagination) TotalPage() int64 {
	return int64(math.Ceil(float64(s.Total()) / float64(s.PageSize())))
}

func (s *Pagination) Offset() int64 {
	return s.PageSize() * (s.Page() - 1)
}
