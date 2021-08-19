package sort

import (
	"encoding/json"
	"io"
)

const (
	BubbleSort    string = "bubble"
	InsertionSort string = "insertion"
	MergeSort     string = "merge"
)

type SortAlgorithm func(elements []int) *StepList

type Step struct {
	ID   int   `json:"id"`
	List []int `json:"list"`
}

type StepList struct {
	Steps []Step `json:"steps"`
}

func newStepList() *StepList {
	return &StepList{
		Steps: make([]Step, 0),
	}
}

func (s *StepList) AddStep(elements []int) {
	temp := make([]int, len(elements))
	copy(temp, elements)
	s.Steps = append(s.Steps, Step{ID: len(s.Steps) + 1, List: temp})
}

type Input struct {
	Elements  []int  `json:"elements"`
	Algorithm string `json:"algorithm"`
}

func DecodeInput(reader io.Reader) (*Input, error) {
	input := new(Input)
	err := json.NewDecoder(reader).Decode(input)
	if err != nil {
		return nil, err
	}
	return input, nil
}

func swap(arr []int, i int, j int) {
	temp := arr[i]
	arr[i] = arr[j]
	arr[j] = temp
}

func bubbleSort(elements []int) *StepList {
	steps := newStepList()
	for i := 0; i < len(elements)-1; i++ {
		for j := 0; j < len(elements)-i-1; j++ {
			if elements[j] > elements[j+1] {
				swap(elements, j+1, j)
				steps.AddStep(elements)
			}
		}
	}
	return steps
}

func insertionSort(elements []int) *StepList {
	steps := newStepList()
	for i := 0; i < len(elements); i++ {
		key := elements[i]
		j := i - 1
		/* Move elements of arr[0..i-1], that are
		   greater than key, to one position ahead
		   of their current position */
		for j >= 0 && elements[j] > key {
			elements[j+1] = elements[j]
			j = j - 1
		}
		elements[j+1] = key
		steps.AddStep(elements)
	}
	return steps
}

func merge(a []int, b []int) []int {
	final := []int{}
	i := 0
	j := 0
	for i < len(a) && j < len(b) {
		if a[i] < b[j] {
			final = append(final, a[i])
			i++
		} else {
			final = append(final, b[j])
			j++
		}
	}
	for ; i < len(a); i++ {
		final = append(final, a[i])
	}
	for ; j < len(b); j++ {
		final = append(final, b[j])
	}
	return final
}

func doMergeSort(items []int, steps *StepList) []int {
	if len(items) < 2 {
		return items
	}
	steps.AddStep(items)
	first := doMergeSort(items[:len(items)/2], steps)
	second := doMergeSort(items[len(items)/2:], steps)
	sorted := merge(first, second)
	steps.AddStep(sorted)
	return sorted
}

func mergeSort(elements []int) *StepList {
	s := newStepList()
	s.AddStep(doMergeSort(elements, s))
	return s
}

func Do(input *Input) *StepList {
	algorithms := map[string]SortAlgorithm{
		BubbleSort:    bubbleSort,
		InsertionSort: insertionSort,
		MergeSort:     mergeSort,
	}
	list := algorithms[input.Algorithm](input.Elements)
	return list
}
