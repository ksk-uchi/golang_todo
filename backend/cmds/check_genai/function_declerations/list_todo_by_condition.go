package check_genai_function_declerations

import (
	"fmt"
	"math/rand"

	"google.golang.org/genai"
)

var Filters = []*genai.FunctionDeclaration{
	{
		Name:        "get_all_todo_ids",
		Description: "システムに登録されているすべてのタスクのID(配列)を取得する",
		Parameters: &genai.Schema{
			Type:       genai.TypeObject,
			Properties: map[string]*genai.Schema{},
		},
		Response: &genai.Schema{
			Type: genai.TypeObject,
			Properties: map[string]*genai.Schema{
				"result": {
					Type: genai.TypeArray,
					Items: &genai.Schema{
						Type: genai.TypeInteger,
					},
				},
			},
		},
	},
	{
		Name:        "filter_for_sunny",
		Description: "与えられたIDの配列から、晴れの日に消化すべきタスクのID(配列)を返却する",
		Parameters: &genai.Schema{
			Type: genai.TypeObject,
			Properties: map[string]*genai.Schema{
				"ids": {
					Type: genai.TypeArray,
					Items: &genai.Schema{
						Type: genai.TypeInteger,
					},
				},
			},
		},
		Response: &genai.Schema{
			Type: genai.TypeObject,
			Properties: map[string]*genai.Schema{
				"result": {
					Type: genai.TypeArray,
					Items: &genai.Schema{
						Type: genai.TypeInteger,
					},
				},
			},
		},
	},
	{
		Name:        "filter_for_rainy",
		Description: "与えられたIDの配列から、雨の日に消化すべきタスクのID(配列)を返却する",
		Parameters: &genai.Schema{
			Type: genai.TypeObject,
			Properties: map[string]*genai.Schema{
				"ids": {
					Type: genai.TypeArray,
					Items: &genai.Schema{
						Type: genai.TypeInteger,
					},
				},
			},
		},
		Response: &genai.Schema{
			Type: genai.TypeObject,
			Properties: map[string]*genai.Schema{
				"result": {
					Type: genai.TypeArray,
					Items: &genai.Schema{
						Type: genai.TypeInteger,
					},
				},
			},
		},
	},
	{
		Name:        "filter_for_other_weather",
		Description: "与えられたIDの配列から、晴れ・雨 以外の日に消化すべきタスクのID(配列)を返却する",
		Parameters: &genai.Schema{
			Type: genai.TypeObject,
			Properties: map[string]*genai.Schema{
				"ids": {
					Type: genai.TypeArray,
					Items: &genai.Schema{
						Type: genai.TypeInteger,
					},
				},
			},
		},
		Response: &genai.Schema{
			Type: genai.TypeObject,
			Properties: map[string]*genai.Schema{
				"result": {
					Type: genai.TypeArray,
					Items: &genai.Schema{
						Type: genai.TypeInteger,
					},
				},
			},
		},
	},
	{
		Name:        "filter_for_male",
		Description: "与えられたIDの配列から、男性向けのタスクのID(配列)を返却する",
		Parameters: &genai.Schema{
			Type: genai.TypeObject,
			Properties: map[string]*genai.Schema{
				"ids": {
					Type: genai.TypeArray,
					Items: &genai.Schema{
						Type: genai.TypeInteger,
					},
				},
			},
		},
		Response: &genai.Schema{
			Type: genai.TypeObject,
			Properties: map[string]*genai.Schema{
				"result": {
					Type: genai.TypeArray,
					Items: &genai.Schema{
						Type: genai.TypeInteger,
					},
				},
			},
		},
	},
	{
		Name:        "filter_for_female",
		Description: "与えられたIDの配列から、女性向けのタスクのID(配列)を返却する",
		Parameters: &genai.Schema{
			Type: genai.TypeObject,
			Properties: map[string]*genai.Schema{
				"ids": {
					Type: genai.TypeArray,
					Items: &genai.Schema{
						Type: genai.TypeInteger,
					},
				},
			},
		},
		Response: &genai.Schema{
			Type: genai.TypeObject,
			Properties: map[string]*genai.Schema{
				"result": {
					Type: genai.TypeArray,
					Items: &genai.Schema{
						Type: genai.TypeInteger,
					},
				},
			},
		},
	},
	{
		Name:        "filter_for_children",
		Description: "与えられたIDの配列から、子供向けのタスクのID(配列)を返却する",
		Parameters: &genai.Schema{
			Type: genai.TypeObject,
			Properties: map[string]*genai.Schema{
				"ids": {
					Type: genai.TypeArray,
					Items: &genai.Schema{
						Type: genai.TypeInteger,
					},
				},
			},
		},
		Response: &genai.Schema{
			Type: genai.TypeObject,
			Properties: map[string]*genai.Schema{
				"result": {
					Type: genai.TypeArray,
					Items: &genai.Schema{
						Type: genai.TypeInteger,
					},
				},
			},
		},
	},
	{
		Name:        "filter_for_adults",
		Description: "与えられたIDの配列から、大人向けのタスクのID(配列)を返却する",
		Parameters: &genai.Schema{
			Type: genai.TypeObject,
			Properties: map[string]*genai.Schema{
				"ids": {
					Type: genai.TypeArray,
					Items: &genai.Schema{
						Type: genai.TypeInteger,
					},
				},
			},
		},
		Response: &genai.Schema{
			Type: genai.TypeObject,
			Properties: map[string]*genai.Schema{
				"result": {
					Type: genai.TypeArray,
					Items: &genai.Schema{
						Type: genai.TypeInteger,
					},
				},
			},
		},
	},
}

// RandomExcludeIDs は引数に ID (配列) をとります。
// 配列が空の場合はそのまま空配列を返します。
// 空ではない場合はランダムにいくつかの要素を除外した配列を新たに作成し返します。
func RandomExcludeIDs(ids []int) []int {
	if len(ids) == 0 {
		return []int{}
	}

	result := make([]int, 0)
	for _, id := range ids {
		// 50%の確率で要素を除外（除外しない要素をresultに追加）
		if rand.Intn(2) == 0 {
			result = append(result, id)
		}
	}

	fmt.Println("RandomExcludeIDs called.")

	return result
}

// GetAllTodoIDs はシステムに登録されているすべてのタスクのID（1から50）を返します。
func GetAllTodoIDs() []int {
	ids := make([]int, 50)
	for i := 0; i < 50; i++ {
		ids[i] = i + 1
	}
	fmt.Println("GetAllTodoIDs called.")
	return ids
}
