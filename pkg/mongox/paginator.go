package mongox

import (
	"context"
	"errors"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Paginator[T PaginatedItem] interface {
	BuildQueries(ctx context.Context, p PaginationParams) (queries []bson.M, sort bson.D, err error)
	Paginate(params PaginationParams, items []T) (PaginationInfo[T], error)
}

type paginator[T PaginatedItem] struct {
	col *mongo.Collection
}

const (
	lessThan    = "$lt"
	greaterThan = "$gt"
)

func NewPaginator[T PaginatedItem](col *mongo.Collection) Paginator[T] {
	return &paginator[T]{
		col: col,
	}
}

func (pg paginator[T]) BuildQueries(ctx context.Context, p PaginationParams) ([]bson.M, bson.D, error) {
	comparisonOp, sortDir := pg.getPaginationOperator(p)

	queries := []bson.M{p.Query}
	sort := bson.D{{Key: "_id", Value: sortDir}}
	if p.IsFirstPage() {
		return queries, sort, nil
	}

	var cursor Cursor
	if p.PointsNext() {
		nextCursorValues, err := p.Next.Decode()
		if err != nil {
			return []bson.M{}, nil, fmt.Errorf("next cursor parse failed: %s", err)
		}
		cursor = nextCursorValues
	} else {
		previousCursorValues, err := p.Previous.Decode()
		if err != nil {
			return []bson.M{}, nil, fmt.Errorf("previous cursor parse failed: %s", err)
		}
		cursor = previousCursorValues
	}

	var cursorQuery bson.M
	cursorQuery, err := pg.generateCursorQuery(comparisonOp, cursor)
	if err != nil {
		return []bson.M{}, nil, err
	}

	queries = append(queries, cursorQuery)

	sort = bson.D{{Key: "_id", Value: sortDir}}

	return queries, sort, nil
}

func (pg paginator[T]) getPaginationOperator(params PaginationParams) (string, int) {
	pointNext := params.Previous == ""
	sortOrder := params.SortOrder
	if pointNext && sortOrder == ASC {
		return greaterThan, 1
	}
	if pointNext && sortOrder == DESC {
		return lessThan, -1
	}
	if !pointNext && sortOrder == ASC {
		return lessThan, -1
	}
	if !pointNext && sortOrder == DESC {
		return greaterThan, 1
	}

	return "", -1
}

func (pg paginator[T]) generateCursorQuery(comparisonOp string, cursor Cursor) (bson.M, error) {
	if cursor.UniqueTimeSeriesFieldValue == "" {
		return nil, errors.New("invalid Cursor Value")
	}

	if comparisonOp != lessThan && comparisonOp != greaterThan {
		return nil, errors.New("invalid comparison operator specified: only $lt and $gt are allowed")
	}

	cuesorObjectID, err := primitive.ObjectIDFromHex(cursor.UniqueTimeSeriesFieldValue)
	if err != nil {
		return nil, errors.New("invalid id")
	}
	query := bson.M{"_id": bson.M{comparisonOp: cuesorObjectID}}
	return query, nil
}

func (pg paginator[T]) Paginate(params PaginationParams, items []T) (PaginationInfo[T], error) {
	var err error
	pagination := PaginationInfo[T]{}
	var nextCur string
	var prevCur string

	if len(items) == 0 {
		return pagination, err
	}

	var lastItemIndex int

	if params.IsFirstPage() {
		if params.HasPagination(len(items)) {
			// 如果一頁是10筆，因為會故意多要一筆，如果是11筆，代表有下一頁
			nextCur, err = pg.createCursor(items[len(items)-2])
			if err != nil {
				return pagination, err
			}
			lastItemIndex = len(items) - 1
			items = items[:lastItemIndex]
		} else {
			// 如果一頁是10筆，因為會故意多要一筆，如果是10筆剛剛好，代表沒下一頁
			lastItemIndex = len(items)
			items = items[:lastItemIndex]
		}
	} else {
		if params.PointsNext() {
			// 因為是往下一頁要，所以一定會有往上一頁的cursor
			prevCur, err = pg.createCursor(items[0])
			if err != nil {
				return pagination, err
			}
			// if pointing next, it always has prev but it might not have next
			if params.HasPagination(len(items)) {
				// 如果一頁是10筆，因為會故意多要一筆，如果是11筆，代表有下一頁
				nextCur, err = pg.createCursor(items[len(items)-2])
				if err != nil {
					return pagination, err
				}
				lastItemIndex = len(items) - 1
				items = items[:lastItemIndex]
			} else {
				// 如果一頁是10筆，因為會故意多要一筆，如果是10筆剛剛好，代表沒下一頁
				lastItemIndex = len(items)
				items = items[:lastItemIndex]
			}
		} else {
			// 因為是往上一頁要，所以一定會有往下一頁的cursor
			if params.HasPagination(len(items)) {
				lastItemIndex = len(items) - 1
				items = items[:lastItemIndex]
				items = pg.reverse(items)
				prevCur, err = pg.createCursor(items[0])
				if err != nil {
					return pagination, err
				}
			} else {
				lastItemIndex = len(items)
				items = items[:lastItemIndex]
				items = pg.reverse(items)
			}
			nextCur, err = pg.createCursor(items[len(items)-1])
			if err != nil {
				return pagination, err
			}

		}
	}
	if err != nil {
		return pagination, err
	}

	// 為了去推算這一頁是否有下一頁，因此會故意去多要一筆，因此在這邊需要去對最後一筆截斷
	return PaginationInfo[T]{
		Items:            items,
		NextCursorString: nextCur,
		PrevCursorString: prevCur,
	}, nil
}

func (pg paginator[T]) createCursor(item T) (string, error) {
	cursorData := NewCursor(item.GetID())
	cursorString, err := cursorData.EncodeCursor()
	if err != nil {
		return "", err
	}
	return cursorString, nil
}

func (pg paginator[T]) generatePager(nextToken string, prevToken string) PaginationInfo[T] {
	return PaginationInfo[T]{
		NextCursorString: nextToken,
		PrevCursorString: prevToken,
	}
}

func (pg paginator[T]) reverse(s []T) []T {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
	return s
}
