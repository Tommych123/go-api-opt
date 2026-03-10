package sum

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
)

// Request типизированный формат входных данных
type Request struct {
	Numbers []int64 `json:"numbers"`
	Repeat  int     `json:"repeat"`
}

type Result struct {
	Sum   int64
	Count int
}

// ParseAndSum оставляем для бенчмарка и других вызовов с []byte
// В оптимизированной версии делаем decode через json.Decoder в типизированную структуру
func ParseAndSum(body []byte) (Result, error) {
	return DecodeAndSum(bytes.NewReader(body))
}

// DecodeAndSum оптимизированный путь без map[string]any и float64
func DecodeAndSum(r io.Reader) (Result, error) {
	var req Request

	dec := json.NewDecoder(r)
	dec.DisallowUnknownFields()

	if err := dec.Decode(&req); err != nil {
		return Result{}, fmt.Errorf("decode json: %w", err)
	}

	if len(req.Numbers) == 0 {
		return Result{}, errors.New("numbers must not be empty")
	}

	if req.Repeat < 1 {
		req.Repeat = 1
	}

	var total int64
	for i := 0; i < req.Repeat; i++ {
		total = 0
		for _, n := range req.Numbers {
			total += n
		}
	}

	return Result{Sum: total, Count: len(req.Numbers)}, nil
}