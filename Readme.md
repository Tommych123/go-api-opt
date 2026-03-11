# Оптимизация простого API-сервиса с профилировкой

Простой HTTP API (`/sum`) + нагрузка + оптимизация по CPU и памяти на основе профилирования  
Инструменты: `net/http/pprof`, `go test -bench`, `benchstat`, `trace`  
Артефакты и сравнение “до/после” сохранены в репозитории

---

## Что реализовано в проекте

- API: `POST /sum`, `GET /healthz`
- Нагрузка: `cmd/loadgen`
- pprof: admin порт `:6060`
- Bench: `internal/sum` + результаты в `bench/`
- Trace: файлы и отчёты в `profiles/*`

---

## Запуск

### Сервис
```bash
go run ./cmd/api
```
- API: `http://127.0.0.1:8080`
- pprof: `http://127.0.0.1:6060/debug/pprof/`

### Нагрузка
```bash
go run ./cmd/loadgen -workers 16 -n 500000 -numbers 1000 -repeat 200
```

---

## Профилировка pprof

Снимать во время нагрузки

```bash
curl -o cpu.pprof    "http://127.0.0.1:6060/debug/pprof/profile?seconds=10"
curl -o heap.pprof   "http://127.0.0.1:6060/debug/pprof/heap?gc=1"
curl -o allocs.pprof "http://127.0.0.1:6060/debug/pprof/allocs?gc=1"

go tool pprof cpu.pprof
go tool pprof heap.pprof
go tool pprof allocs.pprof
```

Сохранённые артефакты:
- `profiles/baseline/{raw,report}`
- `profiles/optimized/{raw,report}`

---

## Бенчмарки и benchstat

Запуск:
```bash
go test ./internal/sum -run '^$' -bench '^(BenchmarkParseAndSum|BenchmarkDecodeAndSum)$' -benchmem -count=10
```

Сравнение:
```bash
benchstat bench/baseline.txt bench/optimized.txt
```

Артефакты:
- `bench/baseline.txt`
- `bench/optimized.txt`
- `bench/benchstat.txt`

### Результаты benchstat (baseline → optimized)

- `allocs/op`: **2015 → 18** (≈ **-99.11%**)
- `B/op`: **58.50KiB → 24.88KiB** (≈ **-57.48%**)
- `sec/op`: **124.2µs → 131.3µs** (≈ **+5.68%**)

---

## Trace-анализ

Снять trace:
```bash
curl -o trace.out "http://127.0.0.1:6060/debug/pprof/trace?seconds=5"
```

Открыть:
```bash
go build -o bin/api ./cmd/api
./bin/api
go tool trace bin/api trace.out
```

Scheduler delay (pprof из trace):
- baseline total delay: **26956.76 ms**
- optimized total delay: **19042.12 ms** (≈ **-29.3%**)

---

## Нагрузка: результаты “до/после”

Параметры: `workers=16`, `requests=500000`, `numbers=1000`, `repeat=200`

Baseline:
- **rps=13397.9**
- p50=863.583µs, p95=2.904916ms, p99=4.24425ms

Optimized:
- run1: **rps=18864.7**, p50=658.334µs, p95=2.069542ms, p99=3.133667ms
- run2: **rps=19542.8**, p50=646.875µs, p95=1.971333ms, p99=2.936583ms

---

## Что оптимизировали

Baseline:
- `json.Unmarshal` в `map[string]any`, числа как `float64` → конвертации
- `io.ReadAll` в handler
- много аллокаций

Optimized:
- типизированный `Request{Numbers []int64, Repeat int}`
- потоковый decode из `r.Body` (`json.Decoder`) без `io.ReadAll`
- `ParseAndSum([]byte)` через `json.Unmarshal` в struct

---

## История коммитов

История отражает процесс: baseline → нагрузка/pprof → измерения → оптимизация → сравнение → артефакты
