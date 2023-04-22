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

func (p *Parser) peekTokenIs(tokenKind token.TokenKind) bool {
	return p.curToken != nil && p.curToken.Next != nil && p.curToken.Next.Kind == tokenKind
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

func (p *Parser) program() ast.Node {
	if p.expect(token.FUNC) {
		return p.parseFunction()
	}

	return p.stmt()
}

func (p *Parser) stmt() ast.Stmt {
	if p.expect(token.IF) {
		return p.parseIfStmt()
	} else if p.curTokenIs(token.LBRACE) {
		return p.parseBlockStmt()
	}

	node := p.expr()
	if node == nil {
		return &ast.ExprStmt{Expr: nil}
	}

	if p.expect(token.FOREACH) {
		return p.parseForEachStmt(node)
	}

	p.expect(token.THEN)
	if p.expect(token.FOR) {
		return p.parseForStmt(node)
	}
	
	if p.expect(token.RETURN) {
		return &ast.ReturnStmt{Token: p.curToken, Value: node}
	}

	return &ast.ExprStmt{Expr: node}
}

func (p *Parser) expr() ast.Expr {
	return p.assign()
}

func (p *Parser) assign() ast.Expr {
	node := p.logical()

	if p.expect(token.ASSIGN) {
		node = &ast.InfixExpr{Token: p.curToken, Left: node, Operator: ast.ASSIGN, Right: p.assign()}
	} else if p.expect(token.PA) {
		node = &ast.InfixExpr{Token: p.curToken, Left: node, Operator: ast.PA, Right: p.add()}
	} else if p.expect(token.MA) {
		node = &ast.InfixExpr{Token: p.curToken, Left: node, Operator: ast.MA, Right: p.add()}
	} else if p.expect(token.AA) {
		node = &ast.InfixExpr{Token: p.curToken, Left: node, Operator: ast.AA, Right: p.add()}
	} else if p.expect(token.SA) {
		node = &ast.InfixExpr{Token: p.curToken, Left: node, Operator: ast.SA, Right: p.add()}
	}

	return node
}

func (p *Parser) logical() ast.Expr {
	node := p.equality()

	if p.expect(token.AND) {
		node = &ast.InfixExpr{Token: p.curToken, Left: node, Operator: ast.AND, Right: p.logical()}
	} else if p.expect(token.OR) {
		node = &ast.InfixExpr{Token: p.curToken, Left: node, Operator: ast.OR, Right: p.logical()}
	}

	return node
}

func (p *Parser) equality() ast.Expr {
	node := p.relational()

	for {
		if p.expect(token.EQ) {
			node = &ast.InfixExpr{Token: p.curToken, Left: node, Operator: ast.EQ, Right: p.relational()}
		} else if p.expect(token.NOT_EQ) {
			node = &ast.InfixExpr{Token: p.curToken, Left: node, Operator: ast.NOT_EQ, Right: p.relational()}
		} else {
			return node
		}
	}
}

func (p *Parser) relational() ast.Expr {
	node := p.list()

	for {
		if p.expect(token.GT) {
			node = &ast.InfixExpr{Token: p.curToken, Left: node, Operator: ast.GT, Right: p.add()}
		} else if p.expect(token.GE) {
			node = &ast.InfixExpr{Token: p.curToken, Left: node, Operator: ast.GE, Right: p.add()}
		} else if p.expect(token.LT) {
			node = &ast.InfixExpr{Token: p.curToken, Left: p.add(), Operator: ast.GT, Right: node}
		} else if p.expect(token.LE) {
			node = &ast.InfixExpr{Token: p.curToken, Left: p.add(), Operator: ast.GE, Right: node}
		} else {
			return node
		}
	}
}

func (p *Parser) list() ast.Expr {
	node := p.add()

	if p.expect(token.RANGE) {
		return &ast.InfixExpr{Token: p.curToken, Left: node, Operator: ast.RANGE, Right: p.add()}
	}

	return node
}

func (p *Parser) add() ast.Expr {
	node := p.mul()

	for {
		if p.expect(token.PLUS) {
			node = &ast.InfixExpr{Token: p.curToken, Left: node, Operator: ast.ADD, Right: p.mul()}
		} else if p.expect(token.MINUS) {
			node = &ast.InfixExpr{Token: p.curToken, Left: node, Operator: ast.SUB, Right: p.mul()}
		} else {
			return node
		}
	}
}

func (p *Parser) mul() ast.Expr {
	node := p.unary()

	for {
		if p.expect(token.ASTERISK) {
			node = &ast.InfixExpr{Token: p.curToken, Left: node, Operator: ast.MUL, Right: p.unary()}
		} else if p.expect(token.SLASH) {
			node = &ast.InfixExpr{Token: p.curToken, Left: node, Operator: ast.DIV, Right: p.unary()}
		} else if p.expect(token.CALET) {
			node = &ast.InfixExpr{Token: p.curToken, Left: node, Operator: ast.EXPONENT, Right: p.unary()}
		} else if p.expect(token.PARCENT) {
			node = &ast.InfixExpr{Token: p.curToken, Left: node, Operator: ast.MODULUS, Right: p.unary()}
		} else {
			return node
		}
	}
}

func (p *Parser) unary() ast.Expr {
	if p.expect(token.PLUS) {
		return &ast.PrefixExpr{Token: p.curToken, Operator: ast.ADD, Right: p.primary()}
	} else if p.expect(token.MINUS) {
		return &ast.PrefixExpr{Token: p.curToken, Operator: ast.SUB, Right: p.primary()}
	} else if p.expect(token.NOT) {
		return &ast.PrefixExpr{Token: p.curToken, Operator: ast.NOT, Right: p.primary()}
	}

	return p.primary()
}

func (p *Parser) primary() ast.Expr {
	if p.expect(token.LPAREN) {
		return p.parseParen()
	}

	if p.curTokenIs(token.IDENT) {
		node := p.parseIdentifier()

		if p.expect(token.LPAREN) {
			return p.parseCallFunc(node.Name)
		}

		if p.curTokenIs(token.L_SQUARE_BRACE) {
			return p.parseCallArray(node)
		}

		return node
	}

	if p.curTokenIs(token.LBRACE) {
		return p.parseListElements()
	}

	if p.expect(token.DOUBLE_QUOTES) {
		return p.parseString()
	}

	if p.curToken == nil {
		p.error(ast.UNEXPECTED_TOKEN, "式が必要です。")
	}

	switch p.curToken.Kind {
	case token.INTEGER:
		if p.peekTokenIs(token.POINT) {
			return p.parseFloat()
		} else {
			return p.parseInteger()
		}
	case token.TRUE, token.FALSE:
		return p.parseBoolean()
	}

	if p.curToken.Kind == token.ILLEGAL {
		p.error(ast.ILLEGAL_CHARACTER, "対応していない文字 = \"%s\"", p.curToken.Literal)
		p.nextToken()
	} else {
		p.error(ast.UNEXPECTED_TOKEN, "予期しない文字が検出されました。 取得した文字=%s", p.curToken.Literal)
		p.nextToken()
	}
	return nil
}

// 宣言された関数の、引数を解析する
func (p *Parser) parseFunctionParams() []*ast.Ident {
	params := []*ast.Ident{}

	if !p.expect(token.LPAREN) {
		p.error(ast.UNEXPECTED_TOKEN, "予期しない文字が検出されました。 取得した文字=%s", p.curToken.Literal)
		return nil
	}

	if p.curTokenIs(token.RPAREN) {
		p.nextToken()
		return params
	}

	if p.curTokenIs(token.IDENT) {
		params = append(params, &ast.Ident{Token: p.curToken, Name: p.curToken.Literal})
		p.nextToken()
		for p.expect(token.COMMA) {
			param := &ast.Ident{Token: p.curToken, Name: p.curToken.Literal}
			params = append(params, param)
			p.nextToken()
		}
	}

	if !p.expect(token.RPAREN) {
		p.error(ast.MISSING_RPAREN, "括弧を閉じてください。")
		return nil
	}

	return params
}

// 関数　<関数名>(引数...) {}
func (p *Parser) parsePreIdentFunc() *ast.FuncStmt {
	funcNode := &ast.FuncStmt{Token: p.curToken, Name: p.curToken.Literal}
	p.nextToken()

	params := p.parseFunctionParams()
	if params == nil {
		return nil
	}

	funcNode.Params = params
	funcNode.Body = p.parseBlockStmt()
	return funcNode
}

// 関数 (引数...)<関数名> {}
func (p *Parser) parsePostIdentFunc() *ast.FuncStmt {
	params := p.parseFunctionParams()
	if params == nil {
		return nil
	}

	if !p.curTokenIs(token.IDENT) {
		p.error(ast.MISSING_FUNCTION_NAME, "関数名が必要です。");
		return nil
	}

	funcNode := &ast.FuncStmt{Token: p.curToken, Name: p.curToken.Literal}
	p.nextToken()
	funcNode.Body = p.parseBlockStmt()
	funcNode.Params = params
	return funcNode
}

func (p *Parser) parseFunction() *ast.FuncStmt {
	if p.curTokenIs(token.IDENT) {
		return p.parsePreIdentFunc()
	} else if p.curTokenIs(token.LPAREN) {
		return p.parsePostIdentFunc()
	}

	p.error(ast.UNEXPECTED_TOKEN, "予期しない文字が検出されました。 取得した文字=%s", p.curToken.Literal);
	return nil
}

func (p *Parser) parseIfStmt() *ast.IfStmt {
	node := &ast.IfStmt{Token: p.curToken}
	node.Condition = p.expr()
	if node.Condition == nil {
		return nil
	}
	p.expect(token.THEN)

	node.Body = p.parseBlockStmt()
	if node.Body == nil {
		return nil
	}

	if p.expect(token.ELSE) {
		node.Else = p.parseBlockStmt()
	}
	return node
}

func (p *Parser) parseBlockStmt() *ast.BlockStmt {
	if !p.expect(token.LBRACE) {
		p.error(ast.UNEXPECTED_TOKEN,"予期しない文字が検出されました。 取得した文字=%s", p.curToken.Literal)
		return nil
	}

	node := &ast.BlockStmt{Token: p.curToken, List: []ast.Stmt{}}

	for !p.curTokenIs(token.RBRACE) && !p.curTokenIs(token.EOF) {
		node.List = append(node.List, p.stmt())
	}

	if !p.expect(token.RBRACE) {
		p.error(ast.MISSING_RBRACE, "括弧を閉じてください。")
		return nil
	}

	return node
}

func (p *Parser) parseForStmt(node ast.Expr) *ast.ForStmt {
	return &ast.ForStmt{Token: p.curToken, Condition: node, Body: p.parseBlockStmt()}
}

func (p *Parser) parseForEachStmt(node ast.Expr) *ast.ForEachStmt {
	return &ast.ForEachStmt{Token: p.curToken, Array: node, Body: p.parseBlockStmt()}
}

func (p *Parser) parseParen() ast.Expr {
	if p.expect(token.RPAREN) {
		if p.curTokenIs(token.IDENT) {
			node := &ast.CallExpr{Token: p.curToken, Name: p.curToken.Literal}
			p.nextToken()
			return node
		}
		p.error(ast.UNEXPECTED_TOKEN, "式が必要です。")
		return nil
	}

	expressions := []ast.Expr{}
	expr := p.expr()
	if expr == nil {
		return nil
	}

	expressions = append(expressions, expr)
	for p.expect(token.COMMA) {
		e := p.expr()
		if e == nil {
			return nil
		}
		expressions  = append(expressions, e)
	}

	if p.curTokenIs(token.EOF) {
		p.error(ast.MISSING_RPAREN, "括弧を閉じてください。")
		return nil
	}

	if p.expect(token.RPAREN) {
		if len(expressions) == 1 && !p.curTokenIs(token.IDENT) {
			return expressions[0]
		} else if p.curTokenIs(token.IDENT) {
			node := &ast.CallExpr{Token: p.curToken, Name: p.curToken.Literal}
			p.nextToken()
			for _, v := range expressions {
				node.Params = append(node.Params, v)
			}
			return node
		}
	}

	p.error(ast.UNEXPECTED_TOKEN, "予期しない文字が検出されました。 取得した文字=%s", p.curToken.Literal)
	return nil
}

func (p *Parser) parseInteger() *ast.Integer {
	str := utils.ToLower(p.curToken.Literal)
	num, _ := strconv.ParseInt(str, 10, 64)
	node := &ast.Integer{Token: p.curToken, Value: num}
	p.nextToken()
	return node
}

func (p *Parser) parseFloat() *ast.Float {
	node := &ast.Float{Token: p.curToken}

	intStr := utils.ToLower(p.curToken.Literal)
	p.nextToken()
	if !p.expect(token.POINT) {
		p.error(ast.UNEXPECTED_TOKEN, "予期しない文字が検出されました。 取得した文字=%s", p.curToken.Literal)
	}
	floatStr := utils.ToLower(p.curToken.Literal)
	p.nextToken()

	node.Value, _ = strconv.ParseFloat(intStr + "." + floatStr, 64)
	return node
}

func (p *Parser) parseBoolean() *ast.Boolean {
	node := &ast.Boolean{Token: p.curToken, Value: false}
	if node.Token.Kind == token.TRUE {
		node.Value = true
	}
	p.nextToken()
	return node
}

func (p *Parser) parseString() *ast.String {
	node := &ast.String{Token: p.curToken, Value: []rune(p.curToken.Literal)}
	p.nextToken()
	if !p.expect(token.DOUBLE_QUOTES) {
		p.error(ast.MISSING_DOUBLE_QOUTES, "引用符が閉じられていません。")
		return nil
	}
	return node
}

func (p *Parser) parseIdentifier() *ast.Ident {
	node := &ast.Ident{Token: p.curToken, Name: p.curToken.Literal}
	p.nextToken()
	return node
}

func (p *Parser) parseCallArray(node *ast.Ident) *ast.IndexExpr {
	array := &ast.IndexExpr{Token: p.curToken, Ident: node}
	p.nextToken()

	index := p.expr()
	if index == nil {
		return nil
	}
	if !p.expect(token.R_SQUARE_BRACE) {
		p.error(ast.MISSING_R_SQUARE_BRACE, "括弧が閉じていません。")
		return nil
	}
	array.IndexList = append(array.IndexList, index)

	for p.expect(token.L_SQUARE_BRACE) {
		array.IndexList = append(array.IndexList, p.expr())
		if !p.expect(token.R_SQUARE_BRACE) {
			p.error(ast.MISSING_R_SQUARE_BRACE, "括弧が閉じていません。")
			return nil
		}
	}

	return array
}

// parseIdentifierで返ってきた識別子を引数に使用する
func (p *Parser) parseCallFunc(identifier string) *ast.CallExpr {
	node := &ast.CallExpr{Token: p.curToken, Name: identifier}

	expr := p.expr()
	if expr == nil {
		return nil
	}

	node.Params = append(node.Params, expr)
	for p.expect(token.COMMA) {
		node.Params = append(node.Params, p.expr())
	}

	if p.expect(token.RPAREN) {
		return node
	}

	p.error(ast.MISSING_RPAREN, "括弧を閉じてください。")
	return nil
}

func (p *Parser) parseListElements() *ast.ArrayExpr {
	array := &ast.ArrayExpr{Token: p.curToken, Elements: []ast.Expr{}}

	if !p.expect(token.LBRACE) {
		p.error(ast.UNEXPECTED_TOKEN, "予期しない文字が検出されました。 取得した文字=%s", p.curToken.Literal)
		return nil
	}

	if p.expect(token.RBRACE) {
		return array
	}

	expr := p.expr()
	if expr == nil {
		return nil
	}

	array.Elements = append(array.Elements, expr)
	for p.expect(token.COMMA) {
		array.Elements = append(array.Elements, p.expr())
	}

	if !p.expect(token.RBRACE) {
		p.error(ast.MISSING_R_SQUARE_BRACE, "括弧が閉じていません。")
		return nil
	}

	return array
}
