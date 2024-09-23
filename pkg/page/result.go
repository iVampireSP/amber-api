package page

import (
	"fmt"
	"math"
)

// PagedResult 定义一个泛型的分页结果结构体
type PagedResult[T any] struct {
	Data       []T // 泛型类型的切片，用于存储实际的数据
	Page       int // 当前页码
	PageSize   int // 每页大小
	TotalCount int // 数据总条数
	TotalPages int // 总页数
}

// CalculateTotalPages 计算总页数
func (p *PagedResult[T]) CalculateTotalPages() {
	if p.PageSize == 0 {
		p.TotalPages = 0
	} else {
		p.TotalPages = int(math.Ceil(float64(p.TotalCount) / float64(p.PageSize)))
	}
}

// 获取指定页的数据
func (p *PagedResult[T]) GetDataForPage(page int) []T {
	if page < 1 || page > p.TotalPages {
		return nil
	}

	startIndex := (page - 1) * p.PageSize
	endIndex := startIndex + p.PageSize
	if endIndex > p.TotalCount {
		endIndex = p.TotalCount
	}

	return p.Data[startIndex:endIndex]
}

// 示例用法
func main() {
	// 假设我们有一些数据
	var data = []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

	// 创建一个分页结果实例
	pagedResult := &PagedResult[int]{
		Data:       data,
		Page:       1,
		PageSize:   5,
		TotalCount: len(data),
	}

	// 计算总页数
	pagedResult.CalculateTotalPages()

	// 获取第一页的数据
	firstPageData := pagedResult.GetDataForPage(1)
	fmt.Println("First Page Data:", firstPageData)

	// 获取第二页的数据
	secondPageData := pagedResult.GetDataForPage(2)
	fmt.Println("Second Page Data:", secondPageData)

	// 输出总页数
	fmt.Printf("Total Pages: %d\n", pagedResult.TotalPages)
}
