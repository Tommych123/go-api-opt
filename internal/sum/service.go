package sum

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
)

type Result struct {
	Sum   int64
	Count int
}

// ParseAndSum читает numbers и repeat из JSON и считает сумму
func ParseAndSum(body []byte) (Result, error) {
	var payload map[string]any
	if err := json.Unmarshal(body, &payload); err != nil {
		return Result{}, fmt.Errorf("decode json: %w", err)
	}

	rawNums, ok := payload["numbers"]
	if !ok {
		return Result{}, errors.New("missing field: numbers")
	}

	// numbers ожидаем как массив
	arr, ok := rawNums.([]any)
	if !ok {
		return Result{}, errors.New("numbers must be an array")
	}

	// Делаем отдельный []int64 и копируем в него значения(для оптимизации)
	nums := make([]int64, 0, len(arr))
	for _, item := range arr {
		f, ok := item.(float64)
		// Проверяем что число целое
		if !ok || math.Trunc(f) != f {
			return Result{}, errors.New("numbers must contain only integers")
		}
		nums = append(nums, int64(f))
	}

	// repeat увеличивает CPU нагрузку чтобы профили были заметнее
	repeat := 1
	if rawRepeat, ok := payload["repeat"]; ok {
		if f, ok := rawRepeat.(float64); ok && math.Trunc(f) == f {
			repeat = int(f)
		}
	}
	if repeat < 1 {
		repeat = 1
	}

	var total int64
	for i := 0; i < repeat; i++ {
		// Считаем сумму заново на каждом проходе
		total = 0
		for _, n := range nums {
			total += n
		}
	}

	return Result{
		Sum:   total,
		Count: len(nums),
	}, nil
}