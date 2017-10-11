package main

import (
	"math"
	"strconv"
	"errors"
)

type Queue []float64

/*
 * Pointer to queue in order to pass by reference
 */
func (queue *Queue) Pop() (num float64, err error) {
	if len(*queue) == 0 {
		return
	}
	num = (*queue)[len(*queue)-1]

	/* removes last element from queue */
	*queue = (*queue)[:len(*queue)-1]
	return
}

func (queue *Queue) Push(num float64) (err error) {
	*queue = append(*queue, num)
	return
}

/*
 * Format for input should resemble:
 * ["123", "5555", "+"]
 */
func calc(input []string) (float64, error) {
	var queue Queue
	/* Maybe do something about this monstrostity later */
	for _, num := range input {
		result, err := strconv.ParseFloat(num, 64)
		if err != nil {
			switch num {
			case "+":
				/*
				 * I don't like assigning num1 and num2 for every single
				 * case.
				 * Not sure how to correctly handle it while preserving
				 * error handling
				 * But obviously this current state is absolutely disgusting
				 */
				num1, err := queue.Pop()
				if err != nil {
					return -1, err
				}
				num2, err := queue.Pop()
				if err != nil {
					return -1, err
				}
				queue.Push(num1 + num2)
			case "-":
				num1, err := queue.Pop()
				if err != nil {
					return -1, err
				}
				num2, err := queue.Pop()
				if err != nil {
					return -1, err
				}
				queue.Push(num2 - num1)
			case "*":
				num1, err := queue.Pop()
				if err != nil {
					return -1, err
				}
				num2, err := queue.Pop()
				if err != nil {
					return -1, err
				}
				queue.Push(num1 * num2)
			case "/":
				num1, err := queue.Pop()
				if err != nil {
					return -1, err
				}
				num2, err := queue.Pop()
				if err != nil {
					return -1, err
				}
				if num1 == 0 {
					err = errors.New("Zero divisor")
					return -1, err
				}
				queue.Push(num2 / num1)
			case "%":
				num1, err := queue.Pop()
				if err != nil {
					return -1, err
				}
				num2, err := queue.Pop()
				if err != nil {
					return -1, err
				}
				queue.Push(math.Mod(num2, num1))
			case "^":
				num1, err := queue.Pop()
				if err != nil {
					return -1, err
				}
				num2, err := queue.Pop()
				if err != nil {
					return -1, err
				}
				queue.Push(math.Pow(num2, num1))
			default:
				return -1, err
			}
		} else {
			queue.Push(result)
		}
	}
	return queue[0], nil
}

/* Hate using strings like this */
func shuntingYard(input []string) (output []string, err error) {
	operator := make([]string, 0, 0)
	output = make([]string, 0, 0)

	precedence := map[string]int{
		"+": 0,
		"-": 0,
		"*": 1,
		"/": 1,
		"%": 1,
		"^": 2,
		"(": -1,
	}

	for _, num := range input {
		/* Not sure if i should be using ParseFloat
		 * for its side effect
		 */
		_, err = strconv.ParseFloat(num, 64)
		if err != nil {
			switch num {
			/*
			 * Parenthesis are handled first in an attempt to ensure
			 * they aren't accidentally pushed to output
			 */
			case "(":
				operator = append(operator, num)
			case ")":
				var leftFound bool
				for i := len(operator) - 1; i >= 0; i-- {
					if operator[i] == "(" {
						leftFound = true
						break
					}
					output = append(output, operator[i])
					operator = operator[:len(operator)-1]
				}
				if leftFound {
					operator = operator[:len(operator)-1]
				} else {
					err = errors.New("Mismatched parenthesis")
					output = nil
					return
				}
			case "+":
				fallthrough
			case "-":
				fallthrough
			case "*":
				fallthrough
			case "/":
				fallthrough
			case "%":
				fallthrough
			case "^":
				for i := len(operator) - 1; i >= 0; i-- {
					if precedence[num] <= precedence[operator[i]] {
						output = append(output, operator[i])
						operator = operator[:len(operator)-1]
					} else {
						break
					}
				}
				operator = append(operator, num)

			default:
				err = errors.New("Invalid input: " + num)
				output = nil
				return
			}
		} else {
			output = append(output, num)
		}
	}
	for i := len(operator) - 1; i >= 0; i-- {
		output = append(output, operator[i])
		operator = operator[:len(operator)-1]
	}
	return
}
