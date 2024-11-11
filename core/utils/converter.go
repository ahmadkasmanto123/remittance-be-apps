package utils

import (
	"encoding/json"
	"love-remittance-be-apps/lib/model"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func Merge(m1 map[string]interface{}, m2 map[string]interface{}) map[string]interface{} {
	merge := make(map[string]interface{})
	for k, v := range m1 {
		merge[k] = v
	}
	for k, v := range m2 {
		merge[k] = v
	}
	return merge
}

func ToMap(data []byte) map[string]interface{} {
	var s fiber.Map
	json.Unmarshal(data, &s)
	return s
}

func StringToMap(data any) map[string]interface{} {
	var s fiber.Map
	out, _ := json.Marshal(data)
	json.Unmarshal(out, &s)
	return s
}

func JsonToObject[T any](data []byte) (*T, []model.ErrorData) {
	var result T
	err := json.Unmarshal(data, &result)
	if err != nil {
		return nil, append([]model.ErrorData{}, model.ErrorData{
			Title:       "error json",
			Description: err.Error(),
		})
	}
	// request := &model.PaymentMethodStandard{}
	errv := NewValidator().Validate(result)
	if errv != nil {
		return nil, errv
	}

	return &result, nil
}

func StringToArray(data string) []string {
	// fmt.Println(time.Now().Format("02 January 2006 | 15:04:05.00"))
	dataSatu := strings.ReplaceAll(data, " ", "")
	newData := strings.ReplaceAll(dataSatu, "[", "")
	newData2 := strings.ReplaceAll(newData, "]", "")
	newData3 := strings.ReplaceAll(newData2, "\"", "")
	splitData := strings.Split(newData3, ",")
	// fmt.Println(time.Now().Format("02 January 2006 | 15:04:05.00"))
	return splitData
}
func AnyToAny[T any](data any) *T {
	var result T
	out, _ := json.Marshal(data)
	json.Unmarshal(out, &result)
	return &result
}

func GetParam(c *fiber.Ctx) model.Param {
	qLimit := c.Query("limit")
	var limit int
	if qLimit == "all" {
		limit = 100
	} else {
		limit = c.QueryInt("limit")
	}
	offset := c.QueryInt("offset")
	searchs := c.Query("search")
	sort := c.Query("sort")
	result := model.Param{Limit: &limit,
		Offset: &offset,
		Search: &searchs,
		Sort:   &sort}

	dd, _ := json.Marshal(result)
	var ress fiber.Map
	json.Unmarshal(dd, &ress)
	// fmt.Println("disini")
	// fmt.Printf("%v", &ress)

	return result
}
