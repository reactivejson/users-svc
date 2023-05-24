// internal/utils/pagination.go
package utils

import (
	"strconv"

	"github.com/pkg/errors"
)

// ParsePageAndPageSize parses the page and page size query parameters.
func ParsePageAndPageSize(pageStr, pageSizeStr string) (int, int, error) {
	page, err := strconv.Atoi(pageStr)
	if err != nil {
		return 0, 0, errors.Wrap(err, "invalid page number")
	}

	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil {
		return 0, 0, errors.Wrap(err, "invalid page size")
	}

	return page, pageSize, nil
}
