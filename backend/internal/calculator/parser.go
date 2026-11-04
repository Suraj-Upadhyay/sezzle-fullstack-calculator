package calculator

import "fmt"

const maxDepth = 64

// node is either a number leaf (op == "") or an operation applied to args.
type node struct {
	value float64
	op    Operation
	args  []*node
}

func leaf(value float64) *node {
	return &node{value: value}
}

func apply(op Operation, args ...*node) *node {
	return &node{op: op, args: args}
}

// parse builds a syntax tree with the usual precedence, lowest first:
//
//	expression = term { ("+" | "-") term }
//	term       = unary { ("*" | "/") unary }
//	unary      = ("+" | "-") unary | power
//	power      = postfix [ "^" unary ]
//	postfix    = primary { "%" }
//	primary    = number | "(" expression ")" | "√" unary
func parse(input string) (*node, error) {
	tokens, err := tokenize(input)
	if err != nil {
		return nil, err
	}
	if tokens[0].kind == tokenEnd {
		return nil, &SyntaxError{Position: 1, Reason: "expression is empty"}
	}

	p := &parser{tokens: tokens}
	tree, err := p.expression()
	if err != nil {
		return nil, err
	}
	if t := p.peek(); t.kind != tokenEnd {
		return nil, unexpected(t)
	}
	return tree, nil
}

type parser struct {
	tokens []token
	pos    int
	depth  int
}

func (p *parser) peek() token {
	return p.tokens[p.pos]
}

func (p *parser) next() token {
	t := p.tokens[p.pos]
	if t.kind != tokenEnd {
		p.pos++
	}
	return t
}

func (p *parser) expression() (*node, error) {
	return p.binary(p.term, map[tokenKind]Operation{
		tokenPlus:  Add,
		tokenMinus: Subtract,
	})
}

func (p *parser) term() (*node, error) {
	return p.binary(p.unary, map[tokenKind]Operation{
		tokenTimes:  Multiply,
		tokenDivide: Divide,
	})
}

func (p *parser) binary(
	operand func() (*node, error), ops map[tokenKind]Operation,
) (*node, error) {
	left, err := operand()
	if err != nil {
		return nil, err
	}
	for {
		op, ok := ops[p.peek().kind]
		if !ok {
			return left, nil
		}
		p.next()
		right, err := operand()
		if err != nil {
			return nil, err
		}
		left = apply(op, left, right)
	}
}

func (p *parser) unary() (*node, error) {
	p.depth++
	defer func() { p.depth-- }()
	if p.depth > maxDepth {
		return nil, &SyntaxError{
			Position: p.peek().pos,
			Reason:   "expression is nested too deeply",
		}
	}

	switch t := p.peek(); t.kind {
	case tokenPlus:
		p.next()
		return p.unary()
	case tokenMinus:
		p.next()
		operand, err := p.unary()
		if err != nil {
			return nil, err
		}
		return apply(Multiply, leaf(-1), operand), nil
	default:
		return p.power()
	}
}

func (p *parser) power() (*node, error) {
	base, err := p.postfix()
	if err != nil || p.peek().kind != tokenPower {
		return base, err
	}
	p.next()
	exponent, err := p.unary()
	if err != nil {
		return nil, err
	}
	return apply(Power, base, exponent), nil
}

func (p *parser) postfix() (*node, error) {
	n, err := p.primary()
	if err != nil {
		return nil, err
	}
	for p.peek().kind == tokenPercent {
		p.next()
		n = apply(Divide, n, leaf(100))
	}
	return n, nil
}

func (p *parser) primary() (*node, error) {
	switch t := p.next(); t.kind {
	case tokenNumber:
		return leaf(t.value), nil
	case tokenOpenParen:
		inner, err := p.expression()
		if err != nil {
			return nil, err
		}
		if closing := p.next(); closing.kind != tokenCloseParen {
			return nil, missingCloseParen(closing)
		}
		return inner, nil
	case tokenSqrt:
		operand, err := p.unary()
		if err != nil {
			return nil, err
		}
		return apply(SquareRoot, operand), nil
	default:
		return nil, unexpected(t)
	}
}

func unexpected(t token) error {
	if t.kind == tokenEnd {
		return &SyntaxError{t.pos, "unexpected end of expression"}
	}
	return &SyntaxError{t.pos, fmt.Sprintf("unexpected %q", t.text)}
}

func missingCloseParen(t token) error {
	if t.kind == tokenEnd {
		return &SyntaxError{t.pos, `missing ")"`}
	}
	return &SyntaxError{t.pos, fmt.Sprintf(`expected ")" but found %q`, t.text)}
}

func syntaxError(index int, reason string) error {
	return &SyntaxError{Position: index + 1, Reason: reason}
}
