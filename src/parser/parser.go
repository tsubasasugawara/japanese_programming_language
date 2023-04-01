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

func (p *Parser) nextToken() {
	p.curToken = p.curToken.Next
}

func (p *Parser) curTokenIs(tokenKind token.TokenKind) bool {
	return p.curToken != nil && p.curToken.Kind == tokenKind
}

func (p *Parser) expect(tokenKind token.TokenKind) bool {
	if p.curTokenIs(tokenKind) {
		p.nextToken()
		return true
	}

	return false
}

func (p *Parser) error(category string, format string, args ...interface{}) {
	err := ast.NewSyntaxError(category, format, args...)
	p.Errors = append(p.Errors, err)
}

/*--------------------------------- 構文解析 ---------------------------------*/

func (p *Parser) program() *ast.Node {
	if p.expect(token.FUNC) {
		return p.parseFunction()
	}

	return p.stmt()
}

func (p *Parser) stmt() *ast.Node {
	if p.expect(token.IF) {
		return p.parseIfStmt()
	} else if p.expect(token.LBRACE) {
		return p.parseBlockStmt()
	}

	node := p.expr()
	if node == nil {
		return node
	}

	p.expect(token.THEN)
	if p.expect(token.FOR) {
		return p.parseForStmt(node)
	}

	if p.expect(token.RETURN) {
		return ast.NewNodeBinop(ast.RETURN, node, nil)
	}

	return node
}

func (p *Parser) expr() *ast.Node {
	return p.assign()
}

func (p *Parser) assign() *ast.Node {
	node := p.equality()

	if p.expect(token.ASSIGN) {
		node = ast.NewNodeBinop(ast.ASSIGN, node, p.assign())
	} else if p.expect(token.PA) {
		node = ast.NewNodeBinop(ast.PA, node, p.add())
	} else if p.expect(token.MA) {
		node = ast.NewNodeBinop(ast.MA, node, p.add())
	} else if p.expect(token.AA) {
		node = ast.NewNodeBinop(ast.AA, node, p.add())
	} else if p.expect(token.SA) {
		node = ast.NewNodeBinop(ast.SA, node, p.add())
	}

	return node
}

func (p *Parser) equality() *ast.Node {
	node := p.relational()

	for {
		if p.expect(token.EQ) {
			node = ast.NewNodeBinop(ast.EQ, node, p.relational())
		} else if p.expect(token.NOT_EQ) {
			node = ast.NewNodeBinop(ast.NOT_EQ, node, p.relational())
		} else {
			return node
		}
	}
}

func (p *Parser) relational() *ast.Node {
	node := p.add()

	for {
		if p.expect(token.GT) {
			node = ast.NewNodeBinop(ast.GT, node, p.add())
		} else if p.expect(token.GE) {
			node = ast.NewNodeBinop(ast.GE, node, p.add())
		} else if p.expect(token.LT) {
			node = ast.NewNodeBinop(ast.GT, p.add(), node)
		} else if p.expect(token.LE) {
			node = ast.NewNodeBinop(ast.GE, p.add(), node)
		} else {
			return node
		}
	}
}

func (p *Parser) add() *ast.Node {
	node := p.mul()

	for {
		if p.expect(token.PLUS) {
			node = ast.NewNodeBinop(ast.ADD, node, p.mul())
		} else if p.expect(token.MINUS) {
			node = ast.NewNodeBinop(ast.SUB, node, p.mul())
		} else {
			return node
		}
	}
}

func (p *Parser) mul() *ast.Node {
	node := p.unary()

	for {
		if p.expect(token.ASTERISK) {
			node = ast.NewNodeBinop(ast.MUL, node, p.unary())
		} else if p.expect(token.SLASH) {
			node = ast.NewNodeBinop(ast.DIV, node, p.unary())
		} else if p.expect(token.CALET) {
			node = ast.NewNodeBinop(ast.EXPONENT, node, p.unary())
		} else if p.expect(token.PARCENT) {
			node = ast.NewNodeBinop(ast.MODULUS, node, p.unary())
		} else {
			return node
		}
	}
}

func (p *Parser) unary() *ast.Node {
	if p.expect(token.PLUS) {
		return ast.NewNodeBinop(ast.ADD, ast.NewIntegerNode(0), p.primary())
	} else if p.expect(token.MINUS) {
		return ast.NewNodeBinop(ast.SUB, ast.NewIntegerNode(0), p.primary())
	}
	return p.primary()
}

func (p *Parser) primary() *ast.Node {
	if p.expect(token.LPAREN) {
		return p.parseParen()
	}

	if p.curTokenIs(token.IDENT) {
		node := p.parseIdentifier()
		if p.expect(token.LPAREN) {
			return p.parseCallFunc(node)
		}
		return node
	}

	if p.curToken != nil && p.curTokenIs(token.INTEGER) {
		return p.parseInteger()
	}

	if p.curToken != nil && p.curToken.Kind == token.ILLEGAL {
		p.error(ast.ILLEGAL_CHARACTER, "対応していない文字 = \"%s\"", p.curToken.Literal)
		p.nextToken()
	} else {
		p.error(ast.UNEXPECTED_TOKEN, "式が必要です。")
	}
	return nil
}

// 宣言された関数の、引数を解析する
func (p *Parser) parseFunctionParams() []*ast.Node {
	params := []*ast.Node{}
	for p.curTokenIs(token.IDENT) && !p.curTokenIs(token.EOF) {
		ident := ast.NewIdentNode(p.curToken.Literal)
		params = append(params, ident)
		p.nextToken()
		p.expect(token.COMMA)
	}
	return params
}

// 関数　<関数名>(引数...) {}
func (p *Parser) parsePreIdentFunc() *ast.Node {
	funcNode := p.parseIdentifier()

	if !p.expect(token.LPAREN) {
		p.error(ast.UNEXPECTED_TOKEN, "括弧が必要です。")
		return nil
	}

	params := p.parseFunctionParams()

	if !p.expect(token.RPAREN) {
		p.error(ast.MISSING_RPAREN, "括弧を閉じてください。")
		return nil
	}

	funcNode.NodeKind = ast.FUNC
	funcNode.Params = params
	funcNode.Body = p.stmt()
	return funcNode
}

// 関数 (引数...)<関数名> {}
func (p *Parser) parsePostIdentFunc() *ast.Node {
	params := p.parseFunctionParams()
	if !p.expect(token.RPAREN) {
		p.error(ast.MISSING_RPAREN, "括弧を閉じてください。")
		return nil
	}

	if !p.curTokenIs(token.IDENT) {
		p.error(ast.MISSING_FUNCTION_NAME, "関数名が必要です。");
		return nil
	}

	funcNode := p.parseIdentifier()
	funcNode.NodeKind = ast.FUNC
	funcNode.Body = p.stmt()
	funcNode.Params = params
	return funcNode
}

func (p *Parser) parseFunction() *ast.Node {
	if p.curTokenIs(token.IDENT) {
		return p.parsePreIdentFunc()
	} else if p.expect(token.LPAREN) {
		return p.parsePostIdentFunc()
	}

	p.error(ast.UNEXPECTED_TOKEN, "予期しない文字が検出されました。 取得した文字=%s", p.curToken.Literal);
	return nil
}

func (p *Parser) parseIfStmt() *ast.Node {
	node := ast.NewNode(ast.IF)
	node.Condition = p.expr()
	if node.Condition == nil {
		return nil
	}
	p.expect(token.THEN)

	node.Then = p.stmt()
	if node.Then == nil {
		return nil
	}

	if p.expect(token.ELSE) {
		node.Else = p.stmt()
	}
	return node
}

func (p *Parser) parseBlockStmt() *ast.Node {
	node := ast.NewNode(ast.BLOCK)
	for !p.expect(token.RBRACE) {
		if p.curToken == nil || p.curTokenIs(token.EOF) {
			p.error(ast.MISSING_RBRACE, "括弧を閉じてください。")
			return nil
		}
		node.Stmts = append(node.Stmts, p.stmt())
	}
	return node
}

func (p *Parser) parseForStmt(node *ast.Node) *ast.Node {
	fNode := ast.NewNode(ast.FOR)
	fNode.Condition = node
	fNode.Then = p.stmt()
	return fNode
}

func (p *Parser) parseParen() *ast.Node {
	if p.expect(token.RPAREN) {
		if p.curTokenIs(token.IDENT) {
			node := p.parseIdentifier()
			node.NodeKind = ast.CALL
			return node
		}
		p.error(ast.UNEXPECTED_TOKEN, "式が必要です。")
		return nil
	}

	expressions := []*ast.Node{}
	for !p.curTokenIs(token.RPAREN) && !p.curTokenIs(token.EOF) {
		expressions  = append(expressions, p.expr())
		p.expect(token.COMMA)
	}

	if p.curTokenIs(token.EOF) {
		p.error(ast.MISSING_RPAREN, "括弧を閉じてください。")
		return nil
	}

	if p.expect(token.RPAREN) {
		if len(expressions) == 1 && !p.curTokenIs(token.IDENT) {
			return expressions[0]
		} else if p.curTokenIs(token.IDENT) {
			node := p.parseIdentifier()
			node.NodeKind = ast.CALL
			for _, v := range expressions {
				node.Params = append(node.Params, v)
			}
			return node
		}
	}
	
	p.error(ast.UNEXPECTED_TOKEN, "予期しない文字が検出されました。 取得した文字=%s", p.curToken.Literal)
	return nil
}

func (p *Parser) parseInteger() *ast.Node {
	str := utils.ToLower(p.curToken.Literal)
	num, _ := strconv.Atoi(str)
	node := ast.NewIntegerNode(num)
	p.nextToken()
	return node
}

func (p *Parser) parseIdentifier() *ast.Node {
	node := ast.NewIdentNode(p.curToken.Literal)
	p.nextToken()
	return node
}

// parseIdentifierで返ってきたnodeを引数に使用する
func (p *Parser) parseCallFunc(node *ast.Node) *ast.Node {
	node.NodeKind = ast.CALL

	for !p.curTokenIs(token.RPAREN) && !p.curTokenIs(token.EOF) {
		node.Params = append(node.Params, p.expr())
		p.expect(token.COMMA)
	}

	if p.expect(token.RPAREN) {
		return node
	}

	p.error(ast.MISSING_RPAREN, "括弧を閉じてください。")
	return nil
}