package lexer

import (
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"
)

type Token int

type ScannedToken struct {
	Type Token
	Text string
}

type LexerLocator struct {
	Position   int
	LineNumber int
}

type Lexer struct {
	Name    string
	Text    string
	Start   int
	CurrentLocator *LexerLocator
	Width   int
}

const (
	TOKEN_EOF   Token = Token(EOF)
	TOKEN_OPEN  Token = '('
	TOKEN_CLOSE Token = ')'
	TOKEN_MINUS Token = '-'
	TOKEN_EQUAL Token = '='
	TOKEN_ERROR Token = iota + 255
	TOKEN_NAME
	TOKEN_VARIABLE_NAME
	TOKEN_CATEGORY_NAME
	TOKEN_NUMBER
	EOF         = -1
	WHITE_SPACE = " \t\n\r"
	RETURN      = '\n'
)

var (
	TokenNames = map[Token]string{}
	RuneTokens = map[rune]Token{}
)

func InitializeTokens() {
	TokenNames[TOKEN_ERROR] = "error"
	TokenNames[TOKEN_OPEN] = "'('"
	TokenNames[TOKEN_CLOSE] = "')'"
	TokenNames[TOKEN_MINUS] = "'-'"
	TokenNames[TOKEN_EQUAL] = "'='"
	TokenNames[TOKEN_NAME] = "name"
	TokenNames[TOKEN_CATEGORY_NAME] = ":name"
	TokenNames[TOKEN_VARIABLE_NAME] = "?name"
	TokenNames[TOKEN_NUMBER] = "number"

	RuneTokens['('] = TOKEN_OPEN
	RuneTokens[')'] = TOKEN_CLOSE
	RuneTokens['-'] = TOKEN_MINUS
	RuneTokens['='] = TOKEN_EQUAL
}

func (t *Token) ToString() (string, error) {
	if t == nil {
		return "", fmt.Errorf("Failed to convert token to string: token is nil")
	}
	return TokenNames[*t], nil
}

func NewLexer(name string, text string) (*Lexer, error) {
	if name == "" || text == "" {
		return nil, fmt.Errorf("Failed to build new lexer: name and text not specified")
	}
	return &Lexer{
		Name: name,
		Text: text,
		CurrentLocator: &LexerLocator{
			LineNumber: 1,
		},
	}, nil
}

func (l *Lexer) Next() (rune, error) {
	if l == nil {
		return rune(0), fmt.Errorf("Failed to get next rune: lexer is nil")
	}
	if l.CurrentLocator == nil {
		return rune(0), fmt.Errorf("Failed to get next rune: lexer locator is nil")
	}
	if l.CurrentLocator.Position >= len(l.Text) {
		l.Width = 0
		return EOF, nil
	}
	r, width := utf8.DecodeRuneInString(l.Text[l.CurrentLocator.Position:])
	l.Width = width
	l.CurrentLocator.Position += width
	if r == RETURN {
		l.CurrentLocator.LineNumber += 1
	}
	return r, nil
}

func (l *Lexer) Backup() error {
	if l == nil {
		return fmt.Errorf("Can't back up last rune: lexer is nil")
	}
	if l.CurrentLocator == nil {
		return fmt.Errorf("Can't back up last rune: lexer locator is nil")
	}
	backedupRuneStart := l.CurrentLocator.Position - l.Width
	backedupRuneEnd := l.CurrentLocator.Position
	if strings.HasPrefix(l.Text[backedupRuneStart:backedupRuneEnd], "\n") {
		// If our location prefix is a return line, we go back one line
		l.CurrentLocator.LineNumber -= 1
	}
	l.CurrentLocator.Position -= l.Width
	return nil
}

func (l *Lexer) Peek() (rune, error) {
	if l == nil {
		return rune(0), fmt.Errorf("Can't peek next rune: lexer is nil")
	}
	if l.CurrentLocator == nil {
		return rune(0), fmt.Errorf("Can't peek next rune: lexer locator is nil")
	}
	r, err := l.Next()
	if err != nil {
		return rune(0), fmt.Errorf("Can't peek next rune: %v", err)
	}
	return r, nil
}

func (l *Lexer) Clear() error {
	if l == nil {
		return fmt.Errorf("Can't clear lexer: lexer is nil")
	}
	if l.CurrentLocator == nil {
		return fmt.Errorf("Can't clear lexer: lexer locator is nil")
	}
	l.Start = l.CurrentLocator.Position
	return nil
}

// Returns true if the next rune is among the input runes
func (l *Lexer) Accepts(runes string) (bool, error) {
	if l == nil {
		return false, fmt.Errorf("Failed to check if lexer accepts next rune: lexer is nil")
	}
	r, err := l.Next()
	if err != nil {
		return false, fmt.Errorf("Failed to check if lexer accepts next rune: %v", err)
	}
	if strings.IndexRune(runes, r) >= 0 {
		return true, nil
	}
	err = l.Backup()
	if err != nil {
		return false, fmt.Errorf("Failed to check if lexer accepts next rune: %v", err)
	}
	return false, nil
}

// Returns true if all the next consecutive accepted runes are consumed
func (l *Lexer) AcceptsSequence(runes string) (bool, error) {
	if l == nil {
		return false, fmt.Errorf("Failed to check if lexer accepts next rune: lexer is nil")
	}
	var onStreak bool
	accepted, err := l.Accepts(runes)
	if err != nil {
		return false, fmt.Errorf("Failed to run sequence acceptor: %v", err)
	}
	for accepted {
		onStreak = true
		accepted, err = l.Accepts(runes) 
		if err != nil {
			return false, fmt.Errorf("Failed to run sequence acceptor: %v", err)
		}
	}
	return onStreak, nil
}

// Returns token with input type, token.text is the text between start and current position of the lexer
func (l *Lexer) CreateToken(t Token) (*ScannedToken, error) {
	if l == nil {
		return nil, fmt.Errorf("Failed to create token from lexer: lexer is nil")
	}
	if l.CurrentLocator == nil {
		return nil, fmt.Errorf("Failed to create token from lexer: lexer locator is nil")
	}
	tk := &ScannedToken{
		Type: t,
		Text: l.Text[l.Start:l.CurrentLocator.Position],
	}
	return tk, nil
}

func (l *Lexer) TokenError(format string, args ...interface{}) (*ScannedToken, error) {
	if l == nil {
		return nil, fmt.Errorf("Can't generate token error from lexer: lexer is nil")
	}
	return &ScannedToken{
		Type: TOKEN_ERROR, 
		Text: fmt.Sprintf(format, args...),	
	}, nil
}

func (l *Lexer) GetNameToken(t Token) (*ScannedToken, error) {
	if l == nil {
		return nil, fmt.Errorf("Can't get name from lexer: lexer is nil")
	}
	r, err := l.Next()
	if err != nil {
		return nil, fmt.Errorf("Can't get name from lexer: %v", err)
	}
	for unicode.IsLetter(r) || unicode.IsDigit(r) || r == '_' || r == '-' {
		r, err = l.Next()
		if err != nil {
			return nil, fmt.Errorf("Can't get name from lexer: %v", err)
		}
	}
	err = l.Backup()
	if err != nil {
		return nil, fmt.Errorf("Can't get name from lexer: %v", err)
	}
	tk, err := l.CreateToken(t)
	if err != nil {
		return nil, fmt.Errorf("Can't get name from lexer: %v", err)
	}
	return tk, nil
}

func (l *Lexer) GetNumberToken() (*ScannedToken, error) {
	if l == nil {
		return nil, fmt.Errorf("Can't get number from lexer: lexer is nil")
	}
	l.AcceptsSequence("-")
	digits := "0123456789"
	l.AcceptsSequence(digits)
	l.Accepts(".")
	l.AcceptsSequence(digits)
	
	ok, err := l.Accepts("eE")
	if err != nil {
		return nil, fmt.Errorf("Can't get number from lexer: %v", err)
	}
	if ok {
		l.Accepts("-")
		l.AcceptsSequence(digits)
	}
	tk, err := l.CreateToken(TOKEN_NUMBER)
	if err != nil {
		return nil, fmt.Errorf("Can't get number from lexer: %v", err)
	}
	return tk, nil
}

func (l *Lexer) GetCommentToken() error {
	if l == nil {
		return fmt.Errorf("Can't get comment from lexer: lexer is nil")
	}
	for t, err := l.Next(); t != '\n'; t, err = l.Next() {
		if err != nil {
			return fmt.Errorf("Can't get comment from lexer: %v", err)
		}
	}
	err := l.Clear()
	if err != nil {
		return fmt.Errorf("Can't get comment from lexer: %v", err)
	}
	return nil
}

func (l *Lexer) ScanToken() (*ScannedToken, error) {
	if l == nil {
		return nil, fmt.Errorf("Can't generate token error from lexer: lexer is nil")
	}
	for {
		r, err := l.Next()
		if err != nil {
			return nil, fmt.Errorf("Failed to scan token: %v", err)
		}
		rp, err := l.Peek()
		if err != nil {
			return nil, fmt.Errorf("Failed to scan token: %v", err)
		}
		if r == '-' && (unicode.IsDigit(rp) || rp == '-') {
			tk, err := l.GetNumberToken()
			if err != nil {
				return nil, fmt.Errorf("Failed to scan token: %v", err)
			}
			return tk, nil
		}
		if rType, ok := RuneTokens[r]; ok {
			tk, err := l.CreateToken(rType)
			if err != nil {
				return nil, fmt.Errorf("Failed to scan token: %v", err)
			}
			return tk, nil
		}
		switch {
		case r == EOF:
			tk, err := l.CreateToken(EOF)
			if err != nil {
				return nil, fmt.Errorf("Failed to scan token: %v", err)
			}
			return tk, nil
		case unicode.IsSpace(r):
			err = l.Clear()
			if err != nil {
				return nil, fmt.Errorf("Failed to scan token: %v", err)
			}
		case r == ';':
			err = l.GetCommentToken()
			if err != nil {
				return nil, fmt.Errorf("Failed to scan token: %v", err)
			}
		case r == '?': 
			tk, err := l.GetNameToken(TOKEN_VARIABLE_NAME)
			if err != nil {
				return nil, fmt.Errorf("Failed to scan token: %v", err)
			}
			return tk, nil
		case r == ':':
			tk, err := l.GetNameToken(TOKEN_CATEGORY_NAME) 
			if err != nil {
				return nil, fmt.Errorf("Failed to scan token: %v", err)
			}
			return tk, nil
		case unicode.IsLetter(r):
			tk, err := l.GetNameToken(TOKEN_NAME)
			if err != nil {
				return nil, fmt.Errorf("Failed to scan token: %v", err)
			}
			return tk, nil
		case unicode.IsDigit(r):
			tk, err := l.GetNumberToken()
			if err != nil {
				return nil, fmt.Errorf("Failed to scan token: %v", err)
			}
			return tk, err
		default:
			return l.TokenError("Unhandle token in PDDL: %c", r)
		}
	}
}
