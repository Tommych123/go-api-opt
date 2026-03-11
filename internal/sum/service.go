package sum

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
)

type Request struct {
	Numbers []int64 `json:"numbers"`
	Repeat  int     `json:"repeat"`
}

type Result struct {
	Sum   int64
	Count int
}

// ParseAndSum быстрый путь для []byte
func ParseAndSum(body []byte) (Result, error) {
	var req Request
	if err := json.Unmarshal(body, &req); err != nil {
		return Result{}, fmt.Errorf("decode json: %w", err)
	}
	return sumReq(req)
}

// DecodeAndSum потоковый путь для HTTP без io.ReadAll
func DecodeAndSum(r io.Reader) (Result, error) {
	var req Request

	dec := json.NewDecoder(r)
	dec.DisallowUnknownFields()

	if err := dec.Decode(&req); err != nil {
		return Result{}, fmt.Errorf("decode json: %w", err)
	}

	return sumReq(req)
}

func sumReq(req Request) (Result, error) {
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
