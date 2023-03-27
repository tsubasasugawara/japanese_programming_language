package parser

import (
	"strconv"

	"jpl/ast"
	"jpl/token"
	"jpl/utils"
)

type Parser struct {
	curToken  *token.Token

	Errors []string
}

func newParser(head *token.Token) *Parser {
	p := &Parser{curToken: head}
	return p
}

func (p *Parser) nextToken() {
	p.curToken = p.curToken.Next
}

func (p *Parser) curTokenIs(tokenKind token.TokenKind) bool {
	return p.curToken.Kind == tokenKind
}

func (p *Parser) consume(tokenKind token.TokenKind) bool {
	if p.curTokenIs(token.NUMBER) ||
		p.curTokenIs(token.EOF) ||
		!p.curTokenIs(tokenKind) {
			return false
		}

	p.nextToken()
	return true
}

func (p *Parser) expect(tokenKind token.TokenKind) bool {
	if p.curTokenIs(tokenKind) {
		p.nextToken()
		return true
	}

	p.appendError("Syntax Error in expect()")
	return false
}

func (p *Parser) appendError(err string) {
	p.Errors = append(p.Errors, err)
}

func (p *Parser) program() *ast.Node {
	return p.stmt()
}

func (p *Parser) stmt() *ast.Node {
	if p.consume(token.IF) {
		node := ast.NewNode(ast.IF)
		node.Condition = p.expr()
		p.consume(token.THEN)

		node.Then = p.stmt()
		if p.consume(token.ELSE) {
			node.Else = p.stmt()
		}
		return node
	}

	node := p.expr()

	if p.consume(token.RETURN) {
		node = ast.NewNodeBinop(ast.RETURN, node, nil)
	}

	return node
}

func (p *Parser) expr() *ast.Node {
	return p.assign()
}

func (p *Parser) assign() *ast.Node {
	node := p.equality()

	if p.consume(token.ASSIGN) {
		node = ast.NewNodeBinop(ast.ASSIGN, node, p.assign())
	}

	return node
}

func (p *Parser) equality() *ast.Node {
	node := p.relational()

	for {
		if p.consume(token.EQ) {
			node = ast.NewNodeBinop(ast.EQ, node, p.relational())
		} else if p.consume(token.NOT_EQ) {
			node = ast.NewNodeBinop(ast.NOT_EQ, node, p.relational())
		} else {
			return node
		}
	}
}

func (p *Parser) relational() *ast.Node {
	node := p.add()

	for {
		if p.consume(token.GT) {
			node = ast.NewNodeBinop(ast.GT, node, p.add())
		} else if p.consume(token.GE) {
			node = ast.NewNodeBinop(ast.GE, node, p.add())
		} else if p.consume(token.LT) {
			node = ast.NewNodeBinop(ast.GT, p.add(), node)
		} else if p.consume(token.LE) {
			node = ast.NewNodeBinop(ast.GE, p.add(), node)
		} else {
			return node
		}
	}
}

func (p *Parser) add() *ast.Node {
	node := p.mul()

	for {
		if p.consume(token.PLUS) {
			node = ast.NewNodeBinop(ast.ADD, node, p.mul())
		} else if p.consume(token.MINUS) {
			node = ast.NewNodeBinop(ast.SUB, node, p.mul())
		} else {
			return node
		}
	}
}

func (p *Parser) mul() *ast.Node {
	node := p.unary()

	for {
		if p.consume(token.ASTERISK) {
			node = ast.NewNodeBinop(ast.MUL, node, p.unary())
		} else if p.consume(token.SLASH) {
			node = ast.NewNodeBinop(ast.DIV, node, p.unary())
		} else {
			return node
		}
	}
}

func (p *Parser) unary() *ast.Node {
	if p.consume(token.PLUS) {
		return ast.NewNodeBinop(ast.ADD, ast.NewNumberNode(0), p.primary())
	} else if p.consume(token.MINUS) {
		return ast.NewNodeBinop(ast.SUB, ast.NewNumberNode(0), p.primary())
	}
	return p.primary()
}

func (p *Parser) primary() *ast.Node {
	if p.consume(token.LPAREN) {
		node := p.expr()
		p.expect(token.RPAREN)
		return node
	}

	if p.curTokenIs(token.IDENT) {
		node := ast.NewIdentNode(p.curToken.Literal)
		p.nextToken()
		return node
	}

	str := utils.ToLower(p.curToken.Literal)
	num, err := strconv.Atoi(str)
	if err != nil {
		p.appendError("Expect Number.")
		p.nextToken()
		return nil
	}
	p.nextToken()
	return ast.NewNumberNode(num)
}

func Parse(head *token.Token) (*ast.Program, []string) {
	p := newParser(head)
	program := ast.NewProgram()

	for !p.curTokenIs(token.EOF) {
		node := p.program()
		if node != nil {
			program.Nodes = append(program.Nodes, node)
		}
	}

	return program, p.Errors
}
