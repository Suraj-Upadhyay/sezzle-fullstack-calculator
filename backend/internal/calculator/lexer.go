package calculator

import (
	"errors"
	"fmt"
	"strconv"
	"unicode"
)

type tokenKind int

const (
	tokenNumber tokenKind = iota
	tokenPlus
	tokenMinus
	tokenTimes
	tokenDivide
	tokenPower
	tokenPercent
	tokenSqrt
	tokenOpenParen
	tokenCloseParen
	tokenEnd
)

type token struct {
	kind  tokenKind
	text  string
	value float64
	pos   int
}

var symbols = map[rune]tokenKind{
	'+': tokenPlus,
	'-': tokenMinus,
	'−': tokenMinus,
	'*': tokenTimes,
	'×': tokenTimes,
	'/': tokenDivide,
	'÷': tokenDivide,
	'^': tokenPower,
	'%': tokenPercent,
	'√': tokenSqrt,
	'(': tokenOpenParen,
	')': tokenCloseParen,
}

const sqrtWord = "sqrt"

func tokenize(input string) ([]token, error) {
	runes := []rune(input)
	var tokens []token

	for i := 0; i < len(runes); {
		r := runes[i]
		switch {
		case unicode.IsSpace(r):
			i++
		case isDigit(r) || r == '.':
			end := scanNumber(runes, i)
			text := string(runes[i:end])
			value, err := strconv.ParseFloat(text, 64)
			if errors.Is(err, strconv.ErrRange) {
				return nil, syntaxError(i, fmt.Sprintf(
					"number %q is out of range", text))
			}
			if err != nil {
				return nil, syntaxError(i, fmt.Sprintf(
					"invalid number %q", text))
			}
			tokens = append(tokens, token{tokenNumber, text, value, i + 1})
			i = end
		case hasWord(runes[i:], sqrtWord):
			tokens = append(tokens, token{tokenSqrt, sqrtWord, 0, i + 1})
			i += len(sqrtWord)
		default:
			kind, ok := symbols[r]
			if !ok {
				return nil, syntaxError(i, fmt.Sprintf(
					"unexpected character %q", r))
			}
			tokens = append(tokens, token{kind, string(r), 0, i + 1})
			i++
		}
	}

	return append(tokens, token{kind: tokenEnd, pos: len(runes) + 1}), nil
}

func scanNumber(runes []rune, i int) int {
	for i < len(runes) && (isDigit(runes[i]) || runes[i] == '.') {
		i++
	}
	if i == len(runes) || (runes[i] != 'e' && runes[i] != 'E') {
		return i
	}

	j := i + 1
	if j < len(runes) && (runes[j] == '+' || runes[j] == '-') {
		j++
	}
	if j == len(runes) || !isDigit(runes[j]) {
		return i
	}
	for j < len(runes) && isDigit(runes[j]) {
		j++
	}
	return j
}

func isDigit(r rune) bool {
	return r >= '0' && r <= '9'
}

func hasWord(runes []rune, word string) bool {
	return len(runes) >= len(word) && string(runes[:len(word)]) == word
}
