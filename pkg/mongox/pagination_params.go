package mongox

import (
	"github.com/samber/lo"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	ASC  = "asc"
	DESC = "desc"
)

type PaginationParamsBuilder struct {
	query          primitive.M
	limit          int64
	sortOrder      string
	paginatedField string
	next           CursorString
	previous       CursorString
}

type PaginationParams struct {
	Query          primitive.M
	Limit          int64
	SortOrder      string
	PaginatedField string
	Next           CursorString
	Previous       CursorString
}

func NewPaginationParamsBuilder() *PaginationParamsBuilder {
	return &PaginationParamsBuilder{
		limit:          24,
		sortOrder:      DESC,
		paginatedField: "_id",
	}
}

func (b *PaginationParamsBuilder) SetQuery(query primitive.M) *PaginationParamsBuilder {
	b.query = query
	return b
}

func (b *PaginationParamsBuilder) SetLimit(limit int64) *PaginationParamsBuilder {
	if limit < 0 {
		limit = 0
	}
	b.limit = limit
	return b
}

func (b *PaginationParamsBuilder) SetSortOrder(sortOrder string) *PaginationParamsBuilder {
	isValid := lo.Contains([]string{DESC, ASC}, sortOrder)
	if isValid {
		b.sortOrder = sortOrder
	}
	return b
}

func (b *PaginationParamsBuilder) SetPaginatedField(paginatedField string) *PaginationParamsBuilder {
	b.paginatedField = paginatedField
	return b
}

func (b *PaginationParamsBuilder) Build() *PaginationParams {
	return &PaginationParams{
		Query:          b.query,
		Limit:          b.limit,
		SortOrder:      b.sortOrder,
		PaginatedField: b.paginatedField,
		Next:           b.next,
		Previous:       b.previous,
	}
}

func (p *PaginationParams) IsFirstPage() bool {
	return p.Previous == "" && p.Next == ""
}

func (p *PaginationParams) PointsNext() bool {
	return p.Previous == "" && p.Next != ""
}

func (p *PaginationParams) HasPagination(itemsNum int) bool {
	return itemsNum > int(p.Limit)
}
