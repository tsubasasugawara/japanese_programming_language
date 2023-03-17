package parser

import (
	"strconv"

	"jpl/ast"
	"jpl/token"
	"jpl/utils"
)

type Parser struct {
	curToken  *token.Token
	peekToken *token.Token

	errors []string
}

func newParser(head *token.Token) *Parser {
	p := &Parser{curToken: head, peekToken: head.Next}
	return p
}

func (p *Parser) nextToken() {
	p.curToken, p.peekToken = p.peekToken, p.peekToken.Next
}

func (p *Parser) curTokenIs(tokenKind token.TokenKind) bool {
	return p.curToken.Kind == tokenKind
}

func (p *Parser) consume(tokenKind token.TokenKind) bool {
	if p.curToken.Kind == token.NUMBER ||
		p.curToken.Kind == token.EOF ||
		p.curToken.Kind != tokenKind {
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

	p.appendError()
	return false
}

func (p *Parser) appendError() {
	p.errors = append(p.errors, p.curToken.Literal)
	p.nextToken()
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
		p.appendError()
		return nil
	}

	return ast.NewNumberNode(num)
}

func Parse(head *token.Token) *ast.Node {
	p := newParser(head)

	var n *ast.Node
	for !p.curTokenIs(token.EOF) {
		node := p.expr()
		if node != nil {
			n.Next = node
			n = n.Next
		}
	}

	return n.Next
}