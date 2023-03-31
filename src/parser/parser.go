package parser

import (
	"strconv"

	"jpl/ast"
	"jpl/token"
	"jpl/utils"
)

type Parser struct {
	curToken  *token.Token

	Errors []ast.Error
}

func newParser(head *token.Token) *Parser {
	p := &Parser{curToken: head}
	return p
}

func (p *Parser) nextToken() {
	p.curToken = p.curToken.Next
}

func (p *Parser) curTokenIs(tokenKind token.TokenKind) bool {
	return p.curToken != nil && p.curToken.Kind == tokenKind
}

func (p *Parser) consume(tokenKind token.TokenKind) bool {
	if p.curTokenIs(token.INTEGER) ||
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

	return false
}

func (p *Parser) appendError(err ast.Error) {
	p.Errors = append(p.Errors, err)
}

func (p *Parser) program() *ast.Node {
	if p.consume(token.FUNC) {
		if !p.curTokenIs(token.IDENT) {
			err := ast.NewSyntaxError(ast.MISSING_FUNCTION_NAME, "関数名が必要です。")
			p.appendError(err);
			return nil
		}
		funcNode := ast.NewNode(ast.FUNC)
		funcNode.Ident = p.curToken.Literal
		p.nextToken()

		if !p.expect(token.LPAREN) {
			err := ast.NewSyntaxError(ast.UNEXPECTED_TOKEN, "括弧が必要です。")
			p.appendError(err)
			return nil
		}
		for p.curTokenIs(token.IDENT) {
			ident := ast.NewIdentNode(p.curToken.Literal)
			funcNode.Params = append(funcNode.Params, ident)
			p.nextToken()
			p.consume(token.COMMA)
		}
		if !p.expect(token.RPAREN) {
			err := ast.NewSyntaxError(ast.MISSING_RPAREN, "括弧を閉じてください。")
			p.appendError(err)
			return nil
		}

		funcNode.Body = p.stmt()
		return funcNode
	}

	return p.stmt()
}

func (p *Parser) stmt() *ast.Node {
	if p.consume(token.IF) {
		node := ast.NewNode(ast.IF)
		node.Condition = p.expr()
		if node.Condition == nil {
			return nil
		}
		p.consume(token.THEN)

		node.Then = p.stmt()
		if node.Then == nil {
			return nil
		}

		if p.consume(token.ELSE) {
			node.Else = p.stmt()
		}
		return node
	} else if p.consume(token.LBRACE) {
		node := ast.NewNode(ast.BLOCK)
		for !p.consume(token.RBRACE) {
			if p.curToken == nil || p.curTokenIs(token.EOF) {
				err := ast.NewSyntaxError(ast.MISSING_RBRACE, "括弧を閉じてください。")
				p.appendError(err)
				return nil
			}
			node.Stmts = append(node.Stmts, p.stmt())
		}
		return node
	}

	node := p.expr()

	if p.consume(token.THEN) && p.consume(token.FOR) {
		fNode := ast.NewNode(ast.FOR)
		fNode.Condition = node
		fNode.Then = p.stmt()
		return fNode
	}

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
	} else if p.consume(token.PA) {
		node = ast.NewNodeBinop(ast.PA, node, p.add())
	} else if p.consume(token.MA) {
		node = ast.NewNodeBinop(ast.MA, node, p.add())
	} else if p.consume(token.AA) {
		node = ast.NewNodeBinop(ast.AA, node, p.add())
	} else if p.consume(token.SA) {
		node = ast.NewNodeBinop(ast.SA, node, p.add())
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
		} else if p.consume(token.CALET) {
			node = ast.NewNodeBinop(ast.EXPONENT, node, p.unary())
		} else if p.consume(token.PARCENT) {
			node = ast.NewNodeBinop(ast.MODULUS, node, p.unary())
		} else {
			return node
		}
	}
}

func (p *Parser) unary() *ast.Node {
	if p.consume(token.PLUS) {
		return ast.NewNodeBinop(ast.ADD, ast.NewIntegerNode(0), p.primary())
	} else if p.consume(token.MINUS) {
		return ast.NewNodeBinop(ast.SUB, ast.NewIntegerNode(0), p.primary())
	}
	return p.primary()
}

func (p *Parser) primary() *ast.Node {
	if p.consume(token.LPAREN) {
		if p.consume(token.RPAREN) {
			err := ast.NewSyntaxError(ast.UNEXPECTED_TOKEN, "式が必要です。")
			p.appendError(err)
			return nil
		}

		node := p.expr()
		if node == nil {
			if !p.consume(token.RPAREN) {
				err := ast.NewSyntaxError(ast.MISSING_RPAREN, "括弧を閉じてください。")
				p.appendError(err)
			}
			return nil
		}

		if p.expect(token.RPAREN) {
			return node
		}

		if p.curTokenIs(token.EOF) {
			err := ast.NewSyntaxError(ast.MISSING_RPAREN, "括弧を閉じてください。")
			p.appendError(err)
			return nil
		}

		err := ast.NewSyntaxError(ast.UNEXPECTED_TOKEN, "予期しない文字が検出されました。 取得した文字=%s", p.curToken.Literal)
		p.appendError(err)
		return node
	}

	if p.curTokenIs(token.IDENT) {
		identifier := p.curToken.Literal
		p.nextToken()

		if p.consume(token.LPAREN) {
			node := ast.NewNode(ast.CALL)
			node.Ident = identifier

			for !p.curTokenIs(token.RPAREN) && !p.curTokenIs(token.EOF) {
				node.Params = append(node.Params, p.expr()) 
				p.consume(token.COMMA)
			}

			if p.expect(token.RPAREN) {
				return node
			}

			err := ast.NewSyntaxError(ast.MISSING_RPAREN, "括弧を閉じてください。")
			p.appendError(err)
			return nil
		}

		node := ast.NewIdentNode(identifier)
		return node
	}

	if p.curToken != nil && p.curToken.Kind != token.EOF && p.curToken.Kind != token.ILLEGAL {
		str := utils.ToLower(p.curToken.Literal)
		num, err := strconv.Atoi(str)
		var node *ast.Node
		if err != nil {
			e := ast.NewSyntaxError(ast.UNEXPECTED_TOKEN, "数値が必要です。 取得した文字=%s", str)
			p.appendError(e)
		} else {
			node = ast.NewIntegerNode(num)
		}
		p.nextToken()
		return node
	}

	if p.curToken != nil && p.curToken.Kind == token.ILLEGAL {
		err := ast.NewSyntaxError(ast.ILLEGAL_CHARACTER, "対応していない文字 = \"%s\"", p.curToken.Literal)
		p.appendError(err)
		p.nextToken()
	} else {
		err := ast.NewSyntaxError(ast.UNEXPECTED_TOKEN, "式が必要です。")
		p.appendError(err)
	}
	return nil
}
	
func Parse(head *token.Token) (*ast.Program, []ast.Error) {
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
