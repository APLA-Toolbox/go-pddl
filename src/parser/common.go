package parser

import (
	"fmt"
	"log"

	"github.com/guilyx/go-pddl/src/config"
	"github.com/guilyx/go-pddl/src/lexer"
	"github.com/guilyx/go-pddl/src/models"
)

type ParserToolbox struct {
	Lexer         *lexer.Lexer
	Peeks         []*lexer.ScannedToken
	nPeeks        int
	Configuration *config.Config
}

func NewParserToolbox(config *config.Config, lx *lexer.Lexer) (*ParserToolbox, error) {
	if lx == nil || config == nil {
		return nil, fmt.Errorf("Failed to create new parser: config or lexer is nil")
	}
	return &ParserToolbox{
		Configuration: config,
		Lexer:         lx,
		Peeks:         make([]*lexer.ScannedToken, config.MaxPeek),
	}, nil
}

func (p *ParserToolbox) NewPddlError(format string, args ...interface{}) *models.PddlError {
	loc, err := p.Locate()
	if err != nil {
		return &models.PddlError{
			Location: nil,
			Error:    fmt.Errorf(format, args...),
		}
	}
	return &models.PddlError{
		Location: loc,
		Error:    fmt.Errorf(format, args...),
	}
}

func (p *ParserToolbox) Next() (*lexer.ScannedToken, error) {
	if p == nil || p.Lexer == nil || p.Peeks == nil {
		return nil, fmt.Errorf("Failed to get the next lexical token")
	}
	if p.nPeeks == 0 {
		tk, err := p.Lexer.ScanToken()
		if err != nil {
			return nil, fmt.Errorf("Failed to get the next lexical token from the parser: %v", err)
		}
		return tk, nil
	}
	t := p.Peeks[0]
	p.Peeks = p.Peeks[1:]
	p.nPeeks -= 1
	return t, nil
}

func (p *ParserToolbox) Locate() (*models.Location, error) {
	if p == nil || p.Lexer == nil || p.Lexer.CurrentLocator == nil {
		return nil, fmt.Errorf("Failed to get the locate the parser")
	}
	return &models.Location{
		Path: p.Lexer.Name,
		Line: p.Lexer.CurrentLocator.LineNumber,
	}, nil
}

func (p *ParserToolbox) ExpectsType(tokenType lexer.Token) (*lexer.ScannedToken, *models.PddlError) {
	if p == nil || p.Lexer == nil || p.Lexer.CurrentLocator == nil {
		return nil, &models.PddlError{
			Location: nil,
			Error:    fmt.Errorf("Expects failed: critical pointers are nil"),
		}
	}
	tk, err := p.Next()
	if err != nil {
		return nil, p.NewPddlError("Expects failed: %v", err)
	}
	if tk.Type != tokenType {
		actualTyp, err := tk.Type.ToString()
		if err != nil {
			return nil, p.NewPddlError("Expects failed: %v", err)
		}
		expectedTyp, err := tokenType.ToString()
		if err != nil {
			return nil, p.NewPddlError("Expects failed: %v", err)
		}
		return nil, p.NewPddlError("Expected [%s], got [%s]", expectedTyp, actualTyp)
	}
	return tk, nil
}

func (p *ParserToolbox) ExpectsText(text string) (*lexer.ScannedToken, *models.PddlError) {
	if p == nil || p.Lexer == nil || p.Lexer.CurrentLocator == nil {
		return nil, &models.PddlError{
			Location: nil,
			Error:    fmt.Errorf("Expects failed: critical pointers are nil"),
		}
	}
	tk, err := p.Next()
	if err != nil {
		return nil, p.NewPddlError("Expects failed: %v", err)
	}
	if tk.Text != text {
		return nil, p.NewPddlError("Expected [%s], got [%s]", text, tk.Text)
	}
	return tk, nil
}

func (p *ParserToolbox) Expects(args ...string) *models.PddlError {
	if p == nil || p.Lexer == nil || p.Lexer.CurrentLocator == nil {
		return &models.PddlError{
			Location: nil,
			Error:    fmt.Errorf("Expects failed: critical pointers are nil"),
		}
	}
	for _, val := range args {
		tk, err := p.Next()
		if err != nil {
			return p.NewPddlError("Expects failed: %v", err)
		}
		if tk.Text != val {
			return p.NewPddlError("Expected [%s], got [%s]", val, tk.Text)
		}
	}
	return nil
}

func (p *ParserToolbox) PeekNth(n int) (*lexer.ScannedToken, *models.PddlError) {
	if p == nil || p.Lexer == nil || p.Lexer.CurrentLocator == nil {
		return nil, &models.PddlError{
			Location: nil,
			Error:    fmt.Errorf("Peek nth failed: critical pointers are nil"),
		}
	}
	if n > p.Configuration.MaxPeek {
		panic("Max peeking threshold surpassed")
	}
	for ; p.nPeeks < n; p.nPeeks++ {
		tk, err := p.Lexer.ScanToken()
		if err != nil {
			return nil, p.NewPddlError("Failed to peek at %dth token: %v", n, err)
		}
		p.Peeks[p.nPeeks] = tk
	}
	return p.Peeks[n-1], nil
}

func (p *ParserToolbox) Peek() (*lexer.ScannedToken, *models.PddlError) {
	if p == nil || p.Lexer == nil || p.Lexer.CurrentLocator == nil {
		return nil, &models.PddlError{
			Location: nil,
			Error:    fmt.Errorf("Peek failed: critical pointers are nil"),
		}
	}
	return p.PeekNth(1)
}

func (p *ParserToolbox) Junk(n int) *models.PddlError {
	if p == nil || p.Lexer == nil || p.Lexer.CurrentLocator == nil {
		return &models.PddlError{
			Location: nil,
			Error:    fmt.Errorf("Junk failed: critical pointers are nil"),
		}
	}
	for i := 0; i < n; i++ {
		_, err := p.Next()
		if err != nil {
			return p.NewPddlError("Failed to junk %d tokens: %v", n, err)
		}
	}
	return nil
}

func (p *ParserToolbox) AcceptsToken(tokenType lexer.Token) (*lexer.ScannedToken, bool, *models.PddlError) {
	if p == nil || p.Lexer == nil || p.Lexer.CurrentLocator == nil {
		return nil, false, &models.PddlError{
			Location: nil,
			Error:    fmt.Errorf("Failed to check if token is accepted: critical pointers are nil"),
		}
	}
	tk, err := p.Peek()
	if err != nil {
		return nil, false, p.NewPddlError("Failed to check if token is accepted: %s", err.Error.Error())
	}
	if tk.Type != tokenType {
		return &lexer.ScannedToken{}, false, nil
	}
	tk, err2 := p.Next()
	if err2 != nil {
		return nil, false, p.NewPddlError("Failed to check if token is accepted: %v", err2)
	}
	return tk, true, nil
}

func (p *ParserToolbox) Accepts(strings ...string) (bool, *models.PddlError) {
	if p == nil || p.Lexer == nil || p.Lexer.CurrentLocator == nil {
		return false, &models.PddlError{
			Location: nil,
			Error:    fmt.Errorf("Accepts failed: critical pointers are nil"),
		}
	}
	if len(strings) > p.Configuration.MaxPeek {
		panic("Max peeking threshold surpassed")
	}
	for i := range strings {
		tk, err := p.PeekNth(i + 1)
		if err != nil {
			return false, p.NewPddlError("Failed to check if [%s] is accepted: %s", strings[i], err.Error.Error())
		}
		if tk.Text != strings[i] {
			return false, nil
		}
	}
	err := p.Junk(len(strings))
	if err != nil {
		return false, p.NewPddlError("Failed to check if strings are accepted: %v", err)
	}
	return true, nil
}

func (p *ParserToolbox) parseNamesAppend(tokenType lexer.Token) ([]*models.Name, *models.PddlError) {
	n, err := p.parseName(tokenType)
	if err != nil {
		return nil, p.NewPddlError("Failed to append parsed name: %v", err)
	}
	names := []*models.Name{n}
	ns, err := p.parseMultipleNames(tokenType)
	if err != nil {
		return nil, p.NewPddlError("Failed to append parsed name: %v", err)
	}
	names = append(names, ns...)
	return names, nil
}

func (p *ParserToolbox) parseMultipleNames(tokenType lexer.Token) ([]*models.Name, *models.PddlError) {
	ids := []*models.Name{}
	for tk, ok, err := p.AcceptsToken(tokenType); ok; tk, ok, err = p.AcceptsToken(tokenType) {
		if err != nil {
			return nil, p.NewPddlError("Failed to parse multiple names: %v", err.Error)
		}
		l, err2 := p.Locate()
		if err2 != nil {
			return nil, p.NewPddlError("Failed to parse multiple names: %v", err2)
		}
		ids = append(ids, &models.Name{
			Name:     tk.Text,
			Location: l,
		})
	}
	return ids, nil
}

func (p *ParserToolbox) parseName(tokenType lexer.Token) (*models.Name, *models.PddlError) {
	t, err := p.ExpectsType(tokenType)
	if err != nil {
		return nil, p.NewPddlError("Failed to parse name: %v", err.Error)
	}
	loc, err2 := p.Locate()
	if err2 != nil {
		return nil, &models.PddlError{
			Location: nil,
			Error:    err2,
		}
	}

	return &models.Name{
		Name:     t.Text,
		Location: loc,
	}, nil
}

func (p *ParserToolbox) parseFunctionTypedList() (funs []*models.Function) {
	for {
		var fs []*models.Function
		pk, _ := p.Peek()
		for pk.Type == lexer.TOKEN_OPEN {
			fs = append(fs, p.parseAtomicFunc())
			pk, _ = p.Peek()
		}
		if len(fs) == 0 {
			break
		}
		typ := p.parseFunctionType()
		for i := range fs {
			fs[i].Types = typ
		}
		funs = append(funs, fs...)
	}
	return
}

func (p *ParserToolbox) parseFunctionType() (typ []*models.TypeName) {
	ok, _ := p.Accepts("-")
	if !ok {
		return
	}
	ls, _ := p.Locate()
	n, _ := p.ExpectsText("number")
	return []*models.TypeName{
		&models.TypeName{
			Name: &models.Name{
				Location: ls,
				Name:     n.Text,
			},
		}}
}

func (p *ParserToolbox) parseActionDef() (act *models.Action) {
	p.Expects("(", ":action")
	defer p.Expects(")")
	act.Name, _ = p.parseName(lexer.TOKEN_VARIABLE_NAME)
	act.Params = p.parseActionParams()
	ok, _ := p.Accepts(":precondition")
	if ok {
		ok2, _ := p.Accepts("(", ")")
		if !ok2 {
			act.Precondition = parsePreGd(p)
		}
	}
	ok, _ = p.Accepts(":effect")
	if ok {
		ok2, _ := p.Accepts("(", ")")
		if !ok2 {
			act.Effect, _ = p.parseEffect()
		}
	}
	return
}

func (p *ParserToolbox) parseActionParams() (parms []*models.TypedEntry) {
	p.Expects(":parameters", "(")
	defer p.Expects(")")
	te, _ := p.parseTypedListString(lexer.TOKEN_VARIABLE_NAME)
	return te
}

func parsePreGd(p *ParserToolbox) models.Formula {
	ok, _ := p.Accepts("(", "and")
	ok2, _ := p.Accepts("(", "forall")

	switch {
	case ok:
		f, _ := p.parseAndGd(parsePreGd)
		return f
	case ok2:
		f := p.parseForAllGd(parsePreGd)
		return f
	}
	return parsePrefGd(p)
}

func (p *ParserToolbox) parseDomainName() (*models.Name, *models.PddlError) {
	err := p.Expects("(", "domain")
	if err != nil {
		return nil, p.NewPddlError("Failed to parse domain name: %v", err.Error)
	}
	defer p.Expects(")")
	return p.parseName(lexer.TOKEN_NAME)
}

func (p *ParserToolbox) parseRequirements() ([]*models.Name, *models.PddlError) {
	reqs := []*models.Name{}
	ok, err := p.Accepts("(", ":requirements")
	if ok {
		defer p.Expects(")")
		tk, err := p.Peek()
		if err != nil {
			return nil, p.NewPddlError("Failed to parse requirements: %v", err.Error)
		}
		typeTok := tk.Type
		for typeTok == lexer.TOKEN_CATEGORY_NAME {
			n, err := p.parseName(lexer.TOKEN_CATEGORY_NAME)
			if err != nil {
				return nil, p.NewPddlError("Failed to parse requirements: %v", err.Error)
			}
			reqs = append(reqs, n)
			tk, err = p.Peek()
			if err != nil {
				return nil, p.NewPddlError("Failed to parse requirements: %v", err.Error)
			}
			typeTok = tk.Type
		}
	}
	if err != nil {
		return nil, p.NewPddlError("Failed to parse requirements: %v", err.Error)
	}
	return reqs, nil
}

func (p *ParserToolbox) parseTypedListString(tokenType lexer.Token) ([]*models.TypedEntry, *models.PddlError) {
	typedList := []*models.TypedEntry{}
	for {
		ids, err := p.parseMultipleNames(tokenType)
		if err != nil {
			return nil, p.NewPddlError("Failed to parse typed string: %v", err.Error)
		}
		tk, err := p.Peek()
		if err != nil {
			return nil, p.NewPddlError("Failed to parse typed string: %v", err.Error)
		}
		if len(ids) == 0 && tk.Type == lexer.TOKEN_MINUS {
			log.Printf("Permissive error")
			fmt.Printf("Permissive error")
		} else if len(ids) == 0 {
			break
		}
		t, err := p.parseType()
		if err != nil {
			return nil, p.NewPddlError("Failed to parse typed list string: %v", err.Error)
		}
		for _, id := range ids {
			typedList = append(typedList, &models.TypedEntry{
				Name:  id,
				Types: t,
			})
		}
	}
	return typedList, nil
}

func (p *ParserToolbox) parseFuncsDef() []*models.Function {
	ok, _ := p.Accepts("(", ":functions")
	if ok {
		defer p.Expects(")")
		return p.parseFunctionTypedList()
	}
	return nil
}

func (p *ParserToolbox) parseActionsDef() (acts []*models.Action) {
	tk, _ := p.Peek()
	for tk.Type == lexer.TOKEN_OPEN {
		act := p.parseActionDef()
		acts = append(acts, act)
		tk, _ = p.Peek()
	}
	return
}

func (p *ParserToolbox) parseType() ([]*models.TypeName, *models.PddlError) {
	typeNames := []*models.TypeName{}
	ok, err := p.Accepts("-")
	if !ok {
		return typeNames, nil
	}
	if err != nil {
		return nil, p.NewPddlError("Failed to parse type: %v", err.Error)
	}
	ok, err = p.Accepts("(")
	if !ok {
		n, err := p.parseName(lexer.TOKEN_NAME)
		if err != nil {
			return nil, p.NewPddlError("Failed to parse type: %v", err.Error)
		}
		return []*models.TypeName{
			{
				Name: n,
			},
		}, nil
	}
	if err != nil {
		return nil, p.NewPddlError("Failed to parse type: %v", err.Error)
	}
	err = p.Expects("either")
	if err != nil {
		return nil, p.NewPddlError("Failed to parse type: %v", err.Error)
	}
	defer p.Expects(")")
	ns, err := p.parseNamesAppend(lexer.TOKEN_NAME)
	if err != nil {
		return nil, p.NewPddlError("Failed to parse type: %v", err.Error)
	}
	for _, id := range ns {
		typeNames = append(typeNames, &models.TypeName{
			Name: id,
		})
	}
	return typeNames, nil
}

func (p *ParserToolbox) parseTypesDefinition() ([]*models.Type, *models.PddlError) {
	types := []*models.Type{}
	ok, err := p.Accepts("(", ":types")
	if ok {
		defer p.Expects(")")
		tls, err := p.parseTypedListString(lexer.TOKEN_NAME)
		if err != nil {
			return nil, p.NewPddlError("Failed to parse types definition: %v", err.Error)
		}
		for _, tp := range tls {
			types = append(types, &models.Type{
				TypedEntry: tp,
			})
		}
	}
	if err != nil {
		return nil, p.NewPddlError("Failed to parse types definition: %v", err.Error)
	}
	return types, nil
}

func (p *ParserToolbox) parseConstantsDefinition() ([]*models.TypedEntry, *models.PddlError) {
	ok, err := p.Accepts("(", ":constants")
	if ok {
		defer p.Expects(")")
		tls, err := p.parseTypedListString(lexer.TOKEN_NAME)
		if err != nil {
			return nil, p.NewPddlError("Failed to parse constants definition: %v", err.Error)
		}
		return tls, nil
	}
	if err != nil {
		return nil, p.NewPddlError("Failed to parse constants definition: %v", err.Error)
	}
	return nil, nil
}

func (p *ParserToolbox) parsePredicatesDefinition() ([]*models.Predicate, *models.PddlError) {
	ok, err := p.Accepts("(", ":predicates")
	if ok {
		defer p.Expects(")")
		pd, err := p.parseAtomicPred()
		if err != nil {
			return nil, p.NewPddlError("Failed to parse predicates definition: %v", err.Error)
		}
		preds := []*models.Predicate{pd}
		tk, err := p.Peek()
		if err != nil {
			return nil, p.NewPddlError("Failed to parse predicates definition: %v", err.Error)
		}
		for tk.Type == lexer.TOKEN_OPEN {
			pd, err = p.parseAtomicPred()
			if err != nil {
				return nil, p.NewPddlError("Failed to parse predicates definition: %v", err.Error)
			}
			tk, err = p.Peek()
			if err != nil {
				return nil, p.NewPddlError("Failed to parse predicates definition: %v", err.Error)
			}
			preds = append(preds, pd)
		}
	}
	if err != nil {
		return nil, p.NewPddlError("Failed to parse predicates definition: %v", err.Error)
	}
	return nil, nil
}

func (p *ParserToolbox) parseAtomicPred() (*models.Predicate, *models.PddlError) {
	err := p.Expects("(")
	if err != nil {
		return nil, p.NewPddlError("Failed to parse atomic predicates: %v", err.Error)
	}
	defer p.Expects(")")
	pname, err := p.parseName(lexer.TOKEN_NAME)
	if err != nil {
		return nil, p.NewPddlError("Failed to parse atomic predicates: %v", err.Error)
	}
	params, err := p.parseTypedListString(lexer.TOKEN_VARIABLE_NAME)
	if err != nil {
		return nil, p.NewPddlError("Failed to parse atomic predicates: %v", err.Error)
	}
	return &models.Predicate{
		Name:       pname,
		Parameters: params,
	}, nil
}

func (p *ParserToolbox) parseAtomicFunc() *models.Function {
	p.Expects("(")
	defer p.Expects(")")
	n, _ := p.parseName(lexer.TOKEN_NAME)
	ps, _ := p.parseTypedListString(lexer.TOKEN_VARIABLE_NAME)
	return &models.Function{
		Name:   n,
		Params: ps,
	}
}

func (p *ParserToolbox) parseFunctioninit() (*models.FunctionInit, *models.PddlError) {
	fi := &models.FunctionInit{}
	ok, err := p.Accepts("(")
	if err != nil {
		return nil, p.NewPddlError("Failed to parse function head: %v", err.Error)
	}
	fi.Name, err = p.parseName(lexer.TOKEN_NAME)
	if err != nil {
		return nil, p.NewPddlError("Failed to parse function head: %v", err.Error)
	}
	if ok {
		fi.Terms, err = p.parseTerms()
		if err != nil {
			return nil, p.NewPddlError("Failed to parse function head: %v", err.Error)
		}
		err = p.Expects(")")
		if err != nil {
			return nil, p.NewPddlError("Failed to parse function head: %v", err.Error)
		}
	}
	return fi, nil
}

func (p *ParserToolbox) parseAssign() (*models.AssignNode, *models.PddlError) {
	err := p.Expects("(")
	if err != nil {
		return nil, p.NewPddlError("Failed to parse assignment operation: %v", err.Error)
	}
	defer p.Expects(")")
	assignNode := &models.AssignNode{}
	assignNode.Operation, err = p.parseName(lexer.TOKEN_NAME)
	if err != nil {
		return nil, p.NewPddlError("Failed to parse assignment operation: %v", err.Error)
	}
	assignNode.AssignedTo, err = p.parseFunctioninit()
	if err != nil {
		return nil, p.NewPddlError("Failed to parse assignment operation: %v", err)
	}
	if n, ok, err := p.AcceptsToken(lexer.TOKEN_NAME); ok {
		if err != nil {
			return nil, p.NewPddlError("Failed to parse assignment operation: %v", err)
		}
		assignNode.IsNumber = true
		assignNode.Number = n.Text
	} else {
		assignNode.FunctionInit, err = p.parseFunctioninit()
		if err != nil {
			return nil, p.NewPddlError("Failed to parse assignment operation: %v", err)
		}
	}
	return assignNode, nil
}

func (p *ParserToolbox) parseForAllEffect(nestedFormula func(*ParserToolbox) (models.Formula, *models.PddlError)) (models.Formula, *models.PddlError) {
	defer p.Expects(")")
	loc, err := p.Locate()
	if err != nil {
		return nil, p.NewPddlError("Failed to parse for all effect: %v", err.Error)
	}
	qv := p.parseQuantVariables()
	f, err2 := nestedFormula(p)
	if err2 != nil {
		return nil, p.NewPddlError("Failed to parse for all effect: %v", err.Error)
	}
	return &models.ForAllNode{
		QuantNode: &models.QuantNode{
			Variables: qv,
			UnaryNode: &models.UnaryNode{
				Node: &models.Node{
					Location: loc,
				},
				Formula: f,
			},
		},
	}, nil
}

func (p *ParserToolbox) parseAndGd(nested func(*ParserToolbox) models.Formula) (models.Formula, *models.PddlError) {
	defer p.Expects(")")
	ps, err := p.parseFormulaStar(nested)
	if err != nil {
		return nil, p.NewPddlError("Failed to parse and grounded: %v", err.Error)
	}
	l, err2 := p.Locate()
	if err2 != nil {
		return nil, p.NewPddlError("Failed to parse and grounded: %v", err.Error)
	}
	return &models.AndNode{
		MultiNode: &models.MultiNode{
			Node: models.Node{
				Location: l,
			},
			Formula: ps,
		},
	}, nil
}

func (p *ParserToolbox) parseFormulaStar(nested func(*ParserToolbox) models.Formula) (fs []models.Formula, err *models.PddlError) {
	tk, err := p.Peek()
	if err != nil {
		return nil, p.NewPddlError("Failed to parse formula star")
	}
	for tk.Type == lexer.TOKEN_OPEN {
		fs = append(fs, nested(p))
	}
	return
}

func (p *ParserToolbox) parseWhenEffect(nestedFormula func(*ParserToolbox) models.Formula) (models.Formula, *models.PddlError) {
	defer p.Expects(")")
	loc, err := p.Locate()
	if err != nil {
		return nil, p.NewPddlError("Failed to parse when effect: %v", err.Error)
	}
	cond := parseGd(p)
	return &models.WhenNode{
		Condition: cond,
		UnaryNode: &models.UnaryNode{
			Node: &models.Node{
				Location: loc,
			},
			Formula: nestedFormula(p),
		},
	}, nil
}

func parsePrefGd(p *ParserToolbox) models.Formula {
	return parseGd(p)
}

func parseOrGd(p *ParserToolbox, nested func(*ParserToolbox) models.Formula) models.Formula {
	defer p.Expects(")")
	f, _ := p.parseFormulaStar(nested)
	l, _ := p.Locate()
	return &models.OrNode{
		&models.MultiNode{
			Node: models.Node{
				Location: l,
			},
			Formula: f,
		},
	}
}

func (p *ParserToolbox) parseNotGd() models.Formula {
	defer p.Expects(")")
	l, _ := p.Locate()
	return &models.NotNode{
		&models.UnaryNode{
			Node: &models.Node{
				Location: l,
			},
			Formula: parseGd(p),
		},
	}
}

func (p *ParserToolbox) parseImplyGd() models.Formula {
	defer p.Expects(")")
	l, _ := p.Locate()
	return &models.ImplyNode{
		BinaryNode: &models.BinaryNode{
			Node: models.Node{
				Location: l,
			},
			Left:  parseGd(p),
			Right: parseGd(p),
		},
	}
}

func (p *ParserToolbox) parseForAllGd(nested func(*ParserToolbox) models.Formula) models.Formula {
	defer p.Expects(")")
	l, _ := p.Locate()
	return &models.ForAllNode{
		QuantNode: &models.QuantNode{
			Variables: p.parseQuantVariables(),
			UnaryNode: &models.UnaryNode{
				Node: &models.Node{
					Location: l,
				},
				Formula: nested(p),
			},
		},
	}
}

func (p *ParserToolbox) parseQuantVariables() []*models.TypedEntry {
	p.Expects("(")
	defer p.Expects(")")
	te, _ := p.parseTypedListString(lexer.TOKEN_VARIABLE_NAME)
	return te
}

func (p *ParserToolbox) parseExistsGd(nested func(*ParserToolbox) models.Formula) models.Formula {
	defer p.Expects(")")
	loc, _ := p.Locate()
	return &models.ExistsNode{
		QuantNode: &models.QuantNode{
			Variables: p.parseQuantVariables(),
			UnaryNode: &models.UnaryNode{
				Node: &models.Node{
					Location: loc,
				},
				Formula: nested(p),
			},
		},
	}
}

func parseGd(p *ParserToolbox) models.Formula {
	ok, _ := p.Accepts("(", "and")
	ok2, _ := p.Accepts("(", "or")
	ok3, _ := p.Accepts("(", "not")
	ok4, _ := p.Accepts("(", "imply")
	ok5, _ := p.Accepts("(", "exists")
	ok6, _ := p.Accepts("(", "forall")

	switch {
	case ok:
		x, _ := p.parseAndGd(parseGd)
		return x
	case ok2:
		x := parseOrGd(p, parseGd)
		return x
	case ok3:
		x := p.parseNotGd()
		if lit, ok := x.(*models.LiteralNode); ok {
			lit.Negative = !lit.Negative
			return lit
		}
		return x
	case ok4:
		x := p.parseImplyGd()
		return x
	case ok5:
		x := p.parseExistsGd(parseGd)
		return x
	case ok6:
		x := p.parseForAllGd(parseGd)
		return x
	}

	x, _ := p.parseLitteral(false)
	return x
}

func parsePEffect(p *ParserToolbox) models.Formula {
	tk, err := p.PeekNth(2)
	if err != nil {
		return nil
	}
	if _, ok := models.AssignOps[tk.Text]; ok {
		tk2, err := p.Peek()
		if err != nil {
			return nil
		}
		if tk2.Type == lexer.TOKEN_OPEN {
			n, err := p.parseAssign()
			if err != nil {
				return nil
			}
			return n
		}
	}
	ln, err := p.parseLitteral(true)
	if err != nil {
		return nil
	}
	return ln
}

func (p *ParserToolbox) parseLitteral(effect bool) (*models.LiteralNode, *models.PddlError) {
	lit := &models.LiteralNode{}
	ok, err := p.Accepts("(", "not")
	if ok {
		lit.Negative = true
		defer p.Expects(")")
	}
	if err != nil {
		return nil, p.NewPddlError("Failed to parse litteral: %v", err.Error)
	}
	err = p.Expects("(")
	if err != nil {
		return nil, p.NewPddlError("Failed to parse litteral: %v", err.Error)
	}
	defer p.Expects(")")

	l, err2 := p.Locate()
	if err2 != nil {
		return nil, p.NewPddlError("Failed to parse litteral: %v", err.Error)
	}
	lit.IsEffect = effect
	lit.Node = &models.Node{
		Location: l,
	}
	ok, err = p.Accepts("=")
	if err != nil {
		return nil, p.NewPddlError("Failed to parse litteral: %v", err.Error)
	}
	if ok {
		lit.Predicate = &models.Name{
			Name:     "=",
			Location: lit.Node.Location,
		}
	} else {
		lit.Predicate, err = p.parseName(lexer.TOKEN_NAME)
		if err != nil {
			return nil, p.NewPddlError("Failed to parse litteral: %v", err.Error)
		}
	}
	lit.Terms, err = p.parseTerms()
	if err != nil {
		return nil, p.NewPddlError("Failed to parse litteral: %v", err.Error)
	}
	return lit, nil
}

func (p *ParserToolbox) parseTerms() ([]*models.Term, *models.PddlError) {
	terms := []*models.Term{}
	for {
		l, err := p.Locate()
		if err != nil {
			return nil, p.NewPddlError("Failed to parse terms: %v", err.Error)
		}
		t, ok, err2 := p.AcceptsToken(lexer.TOKEN_NAME)
		if err2 != nil {
			return nil, p.NewPddlError("Failed to parse terms: %v", err.Error)
		}
		if ok {
			terms = append(terms, &models.Term{
				Name: &models.Name{
					Name:     t.Text,
					Location: l,
				},
			})
			continue
		}
		t, ok, err2 = p.AcceptsToken(lexer.TOKEN_VARIABLE_NAME)
		if err2 != nil {
			return nil, p.NewPddlError("Failed to parse terms: %v", err.Error)
		}
		if ok {
			terms = append(terms, &models.Term{
				Name: &models.Name{
					Name:     t.Text,
					Location: l,
				},
			})
			continue
		}
	}
}

func (p *ParserToolbox) parseEffect() (models.Formula, *models.PddlError) {
	ok, err := p.Accepts("(", "and")
	if ok {
		f, err := p.parseAndEffect(parseConditionalEffect)
		if err != nil {
			return nil, p.NewPddlError("Failed to parse effect: %v", err.Error)
		}
		return f, nil
	}
	if err != nil {
		return nil, p.NewPddlError("Failed to parse effect: %v", err.Error)
	}
	f := parseConditionalEffect(p)
	return f, nil
}

func (p *ParserToolbox) parseAndEffect(nestedFormula func(p *ParserToolbox) models.Formula) (models.Formula, *models.PddlError) {
	defer p.Expects(")")
	fs, err := p.parseFormulaStar(nestedFormula)
	if err != nil {
		return nil, p.NewPddlError("Failed to parse and effect: %v", err.Error)
	}
	l, err2 := p.Locate()
	if err2 != nil {
		return nil, p.NewPddlError("Failed to parse and effect: %v", err.Error)
	}
	return &models.AndNode{
		MultiNode: &models.MultiNode{
			Node: models.Node{
				Location: l,
			},
			Formula: fs,
		},
	}, nil
}

func parseConditionalEffect(p *ParserToolbox) models.Formula {
	ok, err := p.Accepts("(", "and")
	if err != nil {
		return nil
	}
	if ok {
		f, err := p.parseAndEffect(parsePEffect)
		if err != nil {
			return nil
		}
		return f
	}
	f := parsePEffect(p)
	return f
}

func (p *ParserToolbox) parseProbName() *models.Name {
	p.Expects("(", "problem")
	defer p.Expects(")")
	tk, _ := p.parseName(lexer.TOKEN_NAME)
	return tk
}

func (p *ParserToolbox) parseProbDomain() *models.Name {
	p.Expects("(", ":domain")
	defer p.Expects(")")
	tk, _ := p.parseName(lexer.TOKEN_NAME)
	return tk
}

func (p *ParserToolbox) parseObjsDecl() []*models.TypedEntry {
	if ok, _ := p.Accepts("(", ":objects"); ok {
		defer p.Expects(")")
		te, _ := p.parseTypedListString(lexer.TOKEN_NAME)
		return te
	}
	return nil
}

func (p *ParserToolbox) parseInit() (els []models.Formula) {
	p.Expects("(", ":init")
	defer p.Expects(")")
	tk, _ := p.Peek()
	if tk.Type == lexer.TOKEN_OPEN {
		els = append(els, p.parseInitEl())
		tk, _ = p.Peek()
	}
	return
}

func (p *ParserToolbox) parseInitEl() models.Formula {
	loc, _ := p.Locate()
	at, _ := p.parseFunctioninit()
	n, _ := p.ExpectsType(lexer.TOKEN_NAME)
	if ok, _ := p.Accepts("(", "="); ok {
		defer p.Expects(")")
		return &models.AssignNode{
			Node:     &models.Node{
				Location: loc,
			},
			Operation: &models.Name{
				Name: "=",
				Location: loc,
			},
			AssignedTo: at,
			IsInit: true,
			IsNumber: true,
			Number: n.Text,
		}
	}
	ln, _ := p.parseLitteral(false)
	return ln
}

func (p *ParserToolbox) parseGoal() models.Formula {
	p.Expects("(", ":goal")
	defer p.Expects(")")
	return parsePreGd(p)
}
