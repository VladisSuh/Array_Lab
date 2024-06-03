package main

import (
	"bufio"
	"fmt"
	"math"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

var ArrayMap = make(map[string][]int)

func ParseCommand(command string) {
	command = strings.Replace(command, ",", "", -1)
	parts := strings.Fields(command)
	switch strings.ToLower(parts[0]) {
	case "load":
		Load(parts[1], parts[2])
	case "save":
		Save(parts[1], parts[2])
	case "rand":
		count, _ := strconv.Atoi(parts[2])
		lb, _ := strconv.Atoi(parts[3])
		rb, _ := strconv.Atoi(parts[4])
		Rand(parts[1], count, lb, rb)
	case "concat":
		Concat(parts[1], parts[2])
	case "free":
		Free(parts[1])
	case "remove":
		index, _ := strconv.Atoi(parts[2])
		count, _ := strconv.Atoi(parts[3])
		Remove(parts[1], index, count)
	case "copy":
		start, _ := strconv.Atoi(parts[2])
		end, _ := strconv.Atoi(parts[3])
		Copy(parts[1], start, end, parts[4])
	case "sort":
		Sort(parts[1], parts[2])
	case "shuffle":
		Shuffle(parts[1])
	case "stats":
		Stats(parts[1])
	case "print":
		if parts[2] == "all" {
			PrintAll(parts[1])
		} else {
			start, _ := strconv.Atoi(parts[2])
			if len(parts) == 3 {
				PrintSingle(parts[1], start)
			} else {
				end, _ := strconv.Atoi(parts[3])
				PrintRange(parts[1], start, end)
			}
		}
	default:
		fmt.Println("Неизвестная команда")
	}
}

func Load(name, filename string) {
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("Ошибка открытия файла:", err)
		return
	}
	defer file.Close()

	var numbers []int
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		str := scanner.Text()
		for _, char := range str {
			num, err := strconv.Atoi(string(char))
			if err == nil {
				numbers = append(numbers, num)
			}
		}
	}
	ArrayMap[name] = numbers
}

func Save(name, filename string) {
	numbers, ok := ArrayMap[name]
	if !ok {
		fmt.Println("Массив не существует:", name)
		return
	}

	content := ""
	for _, num := range numbers {
		content += fmt.Sprintf("%d\n", num)
	}

	err := os.WriteFile(filename, []byte(content), 0644)
	if err != nil {
		fmt.Println("Ошибка записи в файл:", err)
	}
}

func Rand(name string, count, lb, rb int) {
	rand.New(rand.NewSource(time.Now().UnixNano()))
	var numbers []int
	for i := 0; i < count; i++ {
		numbers = append(numbers, rand.Intn(rb-lb+1)+lb)
	}
	ArrayMap[name] = numbers
}

func Concat(name1, name2 string) {
	arr1, exists1 := ArrayMap[name1]
	arr2, exists2 := ArrayMap[name2]
	if !exists1 || !exists2 {
		if !exists1 {
			fmt.Printf("Массив %s не существует.\n", name1)
		}
		if !exists2 {
			fmt.Printf("Массив %s не существует.\n", name2)
		}
		return
	}

	arr1 = append(arr1, arr2...)
	ArrayMap[name1] = arr1

	fmt.Printf("Массивы %s и %s сконкатенированы.\n", name1, name2)
}

func Free(name string) {
	ArrayMap[name] = []int{}

	fmt.Printf("Массив %s очищен.\n", name)
}

func Remove(name string, index, count int) {
	arr, exists := ArrayMap[name]
	if !exists {
		fmt.Printf("Массив %s не существует.\n", name)
		return
	}

	if index < 0 || index >= len(arr) || index+count > len(arr) {
		fmt.Println("Ошибка: Индекс за пределами границ.")
		return
	}

	arr = append(arr[:index], arr[index+count:]...)
	ArrayMap[name] = arr

	fmt.Printf("Удалены %d элемент(ов) из массива %s, начиная с индекса %d.\n", count, name, index)
}

func Copy(name string, start, end int, dest string) {
	sourceArr, exists := ArrayMap[name]
	if !exists {
		fmt.Printf("Исходный массив %s не существует.\n", name)
		return
	}

	if start < 0 || end >= len(sourceArr) || start > end {
		fmt.Println("Ошибка: Неверный начальный или конечный индекс.")
		return
	}

	copiedPart := make([]int, end-start+1)
	copy(copiedPart, sourceArr[start:end+1])

	ArrayMap[dest] = copiedPart

	fmt.Printf("Скопированы элементы из массива %s в массив %s с индекса %d до %d.\n", name, dest, start, end)
}

func quickSort(arr []int, low, high int, ascending bool) {
	if low < high {
		pi := partition(arr, low, high, ascending)
		quickSort(arr, low, pi-1, ascending)
		quickSort(arr, pi+1, high, ascending)
	}
}

func partition(arr []int, low, high int, ascending bool) int {
	pivot := arr[high]
	i := low - 1
	for j := low; j < high; j++ {
		if ascending {
			if arr[j] < pivot {
				i++
				arr[i], arr[j] = arr[j], arr[i]
			}
		} else {
			if arr[j] > pivot {
				i++
				arr[i], arr[j] = arr[j], arr[i]
			}
		}
	}
	arr[i+1], arr[high] = arr[high], arr[i+1]
	return i + 1
}

func partitionShuffle(arr []int, low, high int) int {
	rand.New(rand.NewSource(time.Now().UnixNano()))
	randomPivotIndex := rand.Intn(high-low+1) + low
	arr[randomPivotIndex], arr[high] = arr[high], arr[randomPivotIndex]

	i := low - 1
	for j := low; j < high; j++ {
		if true {
			i++
			arr[i], arr[j] = arr[j], arr[i]
		}
	}
	arr[i+1], arr[high] = arr[high], arr[i+1]
	return i + 1
}

func Sort(name, order string) {
	arr, exists := ArrayMap[name]
	if !exists {
		fmt.Printf("Массив %s не существует.\n", name)
		return
	}

	ascending := order == "+"

	quickSort(arr, 0, len(arr)-1, ascending)

	ArrayMap[name] = arr

	if ascending {
		fmt.Printf("Массив %s отсортирован в порядке возрастания.\n", name)
	} else {
		fmt.Printf("Массив %s отсортирован в порядке убывания.\n", name)
	}
}

func shuffle(arr []int, low, high int) {
	if low < high {
		pi := partitionShuffle(arr, low, high)
		shuffle(arr, low, pi-1)
		shuffle(arr, pi+1, high)
	}
}

func Shuffle(name string) {
	arr, exists := ArrayMap[name]
	if !exists {
		fmt.Printf("Массив %s не существует.\n", name)
		return
	}

	shuffle(arr, 0, len(arr)-1)

	fmt.Printf("Массив %s перемешан.\n", name)
}

func Stats(name string) {
	arr, exists := ArrayMap[name]
	if !exists {
		fmt.Printf("Массив %s не существует.\n", name)
		return
	}
	if len(arr) == 0 {
		fmt.Println("Массив пустой.")
		return
	}

	maxVal, minVal := arr[0], arr[0]
	maxIndex, minIndex := 0, 0
	sum := 0
	frequency := make(map[int]int)
	maxFreqVal, maxFreq := arr[0], 0

	for i, val := range arr {
		if val > maxVal {
			maxVal, maxIndex = val, i
		}
		if val < minVal {
			minVal, minIndex = val, i
		}

		sum += val

		frequency[val]++
		if frequency[val] > maxFreq || (frequency[val] == maxFreq && val > maxFreqVal) {
			maxFreqVal, maxFreq = val, frequency[val]
		}
	}

	avg := float64(sum) / float64(len(arr))

	maxDeviation := 0.0
	for _, val := range arr {
		deviation := math.Abs(float64(val) - avg)
		if deviation > maxDeviation {
			maxDeviation = deviation
		}
	}

	fmt.Printf("Статистика массива %s:\n", name)
	fmt.Printf("Размер: %d\n", len(arr))
	fmt.Printf("Максимальное значение: %d (Индекс: %d)\n", maxVal, maxIndex)
	fmt.Printf("Минимальное Значение: %d (Индекс: %d)\n", minVal, minIndex)
	fmt.Printf("Наиболее часто встречающийся элемент: %d (Частота: %d)\n", maxFreqVal, maxFreq)
	fmt.Printf("Среднее значение: %.2f\n", avg)
	fmt.Printf("Максимальное отклонение от среднего: %.2f\n", maxDeviation)
}

func PrintAll(name string) {
	arr, exists := ArrayMap[name]
	if !exists {
		fmt.Printf("Массив %s не существует.\n", name)
		return
	}
	if len(arr) == 0 {
		fmt.Println("Массив пуст.")
		return
	}

	fmt.Printf("Все элементы массива %s: ", name)
	for _, value := range arr {
		fmt.Printf("%d ", value)
	}
	fmt.Println()
}

func PrintSingle(name string, index int) {
	arr, exists := ArrayMap[name]
	if !exists {
		fmt.Printf("Массив %s не существует\n", name)
		return
	}
	if index < 0 || index >= len(arr) {
		fmt.Printf("Ошибка: Индекс %d за границами для массива %s.\n", index, name)
		return
	}

	fmt.Printf("Элемент для индекса %d массива %s: %d\n", index, name, arr[index])
}

func PrintRange(name string, start, end int) {
	arr, exists := ArrayMap[name]
	if !exists {
		fmt.Printf("Массив %s не существует.\n", name)
		return
	}
	if start < 0 || end >= len(arr) || start > end {
		fmt.Println("Ошибка, неправильные границы массива.")
		return
	}

	fmt.Printf("Элементы массива %s с индекса %d до %d: ", name, start, end)
	for i := start; i <= end; i++ {
		fmt.Printf("%d ", arr[i])
	}
	fmt.Println()
}

func main() {
	file, err := os.Open("commands.txt")
	if err != nil {
		fmt.Println("Ошибка в открытии файла", err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		ParseCommand(scanner.Text())
	}
}
