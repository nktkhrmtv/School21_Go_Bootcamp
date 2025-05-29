package minCoinsFuncs

import (
	"math"
)

/*
Documemtation: 
    For gen use godoc: 
    (bash)  godoc -http=:6060
    http://localhost:6060/pkg/ Documentation in Browser: pkg/day07/minCoinsFuncs	

    For save HTML use godoc:
    1) (bash src/)  godoc -url "http://localhost:6060/pkg/day07/minCoinsFuncs/" > docs.html 
    2) (bash src/)  zip -r docs.zip docs.html static/ 
*/                 


// MinCoins возвращает набор монет минимального размера, сумма которых равна val.
// Используется жадный алгоритм, который не всегда гарантирует оптимальное решение.
// Например, для val = 6 и coins = [1, 3, 4] результат будет [4, 1, 1], хотя оптимальное решение — [3, 3].
//
// Параметры:
//   - val: целевое значение (сумма монет).
//   - coins: список номиналов монет (может содержать дубликаты и не быть отсортированным).
//
// Возвращает:
//   - Слайс монет, сумма которых равна val. Если решение не найдено, возвращает пустой слайс.
func MinCoins(val int, coins []int) []int {
    res := make([]int, 0)
    i := len(coins) - 1
    for i >= 0 {
        for val >= coins[i] {
            val -= coins[i]
            res = append(res, coins[i])
        }
        i -= 1
    }
    return res
}

// MinCoins2 возвращает набор монет минимального размера, сумма которых равна val.
// Используется динамическое программирование, что гарантирует оптимальное решение.
// Оптимизации:
//   - Удаление дубликатов из списка монет.
//   - Использование массива prev для восстановления набора монет.
//   - Использование math.MaxInt32, чтобы избежать потенциальных проблем с переполнением.
//   - Прерывание цикла, если монета больше текущей суммы.
// Параметры:
//   - val: целевое значение (сумма монет).
//   - coins: список номиналов монет (может содержать дубликаты и не быть отсортированным).
//  
// Возвращает:
//   - Слайс монет, сумма которых равна val. Если решение не найдено, возвращает пустой слайс.
func MinCoins2(val int, coins []int) []int {
    if val == 0 || len(coins) == 0 {
        return []int{}
    }

    uniqueCoins := make([]int, 0)
    seen := make(map[int]bool)
    for _, coin := range coins {
        if !seen[coin] && coin > 0 {
            seen[coin] = true
            uniqueCoins = append(uniqueCoins, coin)
        }
    }

    dp := make([]int, val+1)
    for i := 1; i <= val; i++ {
        dp[i] = math.MaxInt32 
    }

    prev := make([]int, val+1)

    for i := 1; i <= val; i++ {
        for _, coin := range uniqueCoins {
            if coin > i {
				break 
			}
            if coin <= i && dp[i-coin]+1 < dp[i] {
                dp[i] = dp[i-coin] + 1
                prev[i] = coin
            }
        }
    }

    if dp[val] == math.MaxInt32 {
        return []int{} 
    }

    res := make([]int, 0)
    for val > 0 {
        res = append(res, prev[val])
        val -= prev[val]
    }

    return res
}
