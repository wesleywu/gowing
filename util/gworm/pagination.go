package gworm

import "strings"

type SortDirection uint8

const (
	SortDirection_ASC  SortDirection = 0
	SortDirection_DESC SortDirection = 1
)

// Enum value maps for MultiType.
var (
	SortDirection_name = map[SortDirection]string{
		0: "ASC",
		1: "DESC",
	}
	SortDirection_value = map[string]SortDirection{
		"ASC":  0,
		"DESC": 1,
	}
)

type PageRequest struct {
	Number uint32

	Offset uint64
	Size   uint32
	Sorts  []SortParam
}

type SortParam struct {
	Property  string
	Direction SortDirection
}

type PageInfo struct {
	Offset           uint64      // the offset of the first element in current page of the paging request
	Size             uint32      // page size of the paging request, may be larger than NumberOfElements
	Sorts            []SortParam // the sorting parameters of the paging request
	Number           uint32      // page number of current page, starting from 1
	NumberOfElements uint32      // number of elements in current page
	TotalElements    uint64      // total number of elements for current request when without paging
	TotalPages       uint32      // total number of pages of the paging request
	First            bool        // whether current page is first page
	Last             bool        // whether current page is last page
}

func (pr PageRequest) Of(page, size uint32, sort ...SortParam) PageRequest {
	if page < 1 {
		page = 1
	}
	if size <= 0 {
		size = 10
	}
	pr.Number = page
	pr.Offset = uint64(page-1) * uint64(size)
	pr.Size = size
	pr.Sorts = sort
	return pr
}

func (pr PageRequest) AddSortByString(sort string) PageRequest {
	props := strings.Split(sort, ",")
	for _, prop := range props {
		parts := strings.Split(strings.TrimSpace(prop), " ")
		switch len(parts) {
		case 1:
			pr.AddSort(SortParam{
				Property:  parts[0],
				Direction: SortDirection_ASC,
			})
		case 2:
			direction, ok := SortDirection_value[strings.ToUpper(parts[1])]
			if !ok {
				direction = SortDirection_ASC
			}
			pr.AddSort(SortParam{
				Property:  parts[0],
				Direction: direction,
			})
		default:
			continue
		}
	}
	return pr
}

func (pr PageRequest) AddSort(sort SortParam) PageRequest {
	pr.Sorts = append(pr.Sorts, sort)
	return pr
}

func (pr PageRequest) OrderString() string {
	sb := strings.Builder{}
	for i, s := range pr.Sorts {
		if i > 0 {
			sb.WriteString(",")
		}
		sb.WriteString(s.Property)
		sb.WriteString(" ")
		sb.WriteString(SortDirection_name[s.Direction])
	}
	return sb.String()
}

func (pi *PageInfo) From(pageRequest PageRequest, numberOfElement uint32, totalElements uint64) {
	pi.Offset = pageRequest.Offset
	pi.Size = pageRequest.Size
	pi.Sorts = pageRequest.Sorts
	pi.Number = pageRequest.Number
	pi.NumberOfElements = numberOfElement
	pi.TotalElements = totalElements
	if pi.Size > 0 {
		pi.TotalPages = uint32(totalElements / uint64(pi.Size))
	} else {
		pi.TotalPages = 0
	}
	pi.First = pi.Number == 1
	pi.Last = (pi.Offset + uint64(pi.NumberOfElements)) >= pi.TotalElements
}
