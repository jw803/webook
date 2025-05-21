package mongox

type PaginationInfo[T any] struct {
	Items            []T    `json:"items"`
	NextCursorString string `json:"next_cursor"`
	PrevCursorString string `json:"previous_cursor"`
}

type PaginatedItem interface {
	GetID() string
}
