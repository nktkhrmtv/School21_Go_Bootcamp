package main

import (
	"errors"
	"fmt"
	"unsafe"
)

func getElement(arr []int, idx int) (int, error) {
	if len(arr) == 0 {
		return 0, errors.New("срез пуст")
	}
	if idx < 0 {
		return 0, errors.New("индекс не может быть отрицательным")
	}
	if idx >= len(arr) {
		return 0, errors.New("индекс вне допустимых границ")
	}

	firstElementPtr := &arr[0]

	elementPtr := uintptr(unsafe.Pointer(firstElementPtr)) + uintptr(idx)*unsafe.Sizeof(arr[0])
	element := *(*int)(unsafe.Pointer(elementPtr))

	return element, nil
}

func main() {
	arr := []int{10, 20, 30, 40, 50}

	element, err := getElement(arr, 4)
	if err != nil {
		fmt.Println("Ошибка:", err)
	} else {
		fmt.Println("Элемент:", element)
	}
}