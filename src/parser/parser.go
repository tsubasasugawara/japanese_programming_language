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

func (p *Parser) expr() *ast.Node {
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
	node := p.primary()

	for {
		if p.consume(token.ASTERISK) {
			node = ast.NewNodeBinop(ast.MUL, node, p.primary())
		} else if p.consume(token.SLASH) {
			node = ast.NewNodeBinop(ast.DIV, node, p.primary())
		} else {
			return node
		}
	}
}

func (p *Parser) primary() *ast.Node {
	if p.consume(token.LPAREN) {
		node := p.expr()
		p.expect(token.RPAREN)
		return node
	}

	str := utils.ToLower(p.curToken.Literal)
	num, err := strconv.Atoi(str)
	if err != nil {
		p.appendError("Expect Number.")
		return nil
	}
	p.nextToken()
	return ast.NewNumberNode(num)
}

func Parse(head *token.Token) (*ast.Program, []string) {
	p := newParser(head)
	program := ast.NewProgram()

	for !p.curTokenIs(token.EOF) {
		node := p.expr()
		if node != nil {
			program.Nodes = append(program.Nodes, node)
		}
	}

	return program, p.Errors
}
