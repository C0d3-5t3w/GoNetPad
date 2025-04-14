package tools

import (
	"image/color"
	"regexp"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type SyntaxHighlighter struct {
	language string
}

func NewSyntaxHighlighter(language string) *SyntaxHighlighter {
	return &SyntaxHighlighter{
		language: strings.ToLower(language),
	}
}

func DetectLanguage(filename string, content string) string {
	if strings.HasSuffix(filename, ".go") {
		return "go"
	} else if strings.HasSuffix(filename, ".js") {
		return "javascript"
	} else if strings.HasSuffix(filename, ".ts") {
		return "typescript"
	} else if strings.HasSuffix(filename, ".html") {
		return "html"
	} else if strings.HasSuffix(filename, ".css") {
		return "css"
	}

	if strings.Contains(content, "package ") && strings.Contains(content, "import ") {
		return "go"
	} else if strings.Contains(content, "<html") || strings.Contains(content, "<!DOCTYPE html") {
		return "html"
	} else if strings.Contains(content, "{") && strings.Contains(content, "}") {
		if strings.Contains(content, "interface ") || strings.Contains(content, ": string") {
			return "typescript"
		}
		if strings.Contains(content, "function") || strings.Contains(content, "var ") {
			return "javascript"
		}
		if strings.Contains(content, "@media") || strings.Contains(content, "px") {
			return "css"
		}
	}

	return "text"
}

type tokenType int

const (
	tokenPlain tokenType = iota
	tokenKeyword
	tokenString
	tokenNumber
	tokenComment
	tokenFunction
	tokenTag
	tokenAttribute
)

var tokenStyles = map[tokenType]color.RGBA{
	tokenKeyword:   {R: 170, G: 13, B: 145, A: 255},
	tokenString:    {R: 196, G: 26, B: 22, A: 255},
	tokenNumber:    {R: 28, G: 0, B: 207, A: 255},
	tokenComment:   {R: 63, G: 127, B: 95, A: 255},
	tokenFunction:  {R: 66, G: 113, B: 174, A: 255},
	tokenTag:       {R: 0, G: 104, B: 129, A: 255},
	tokenAttribute: {R: 156, G: 101, B: 0, A: 255},
}

type token struct {
	tokenType tokenType
	text      string
}

func (sh *SyntaxHighlighter) HighlightCode(code string) *fyne.Container {
	tokens := sh.tokenize(code)
	var textObjects []fyne.CanvasObject

	for _, t := range tokens {
		text := canvas.NewText(t.text, tokenStyles[t.tokenType])
		text.TextStyle = fyne.TextStyle{Monospace: true}
		textObjects = append(textObjects, text)
	}

	return container.NewVBox(textObjects...)
}

func (sh *SyntaxHighlighter) HighlightCodeAsRichText(code string) string {
	tokens := sh.tokenize(code)
	var highlightedText strings.Builder

	for _, t := range tokens {
		switch t.tokenType {
		case tokenKeyword:
			highlightedText.WriteString("[Keyword]" + t.text + "[/Keyword]")
		case tokenString:
			highlightedText.WriteString("[String]" + t.text + "[/String]")
		case tokenComment:
			highlightedText.WriteString("[Comment]" + t.text + "[/Comment]")
		default:
			highlightedText.WriteString(t.text)
		}
	}

	return highlightedText.String()
}

func (sh *SyntaxHighlighter) HighlightCodeAsRichTextSegments(code string) []widget.RichTextSegment {
	tokens := sh.tokenize(code)
	var segments []widget.RichTextSegment

	for _, t := range tokens {
		style := richTextStyleForToken(t.tokenType)
		segments = append(segments, &widget.TextSegment{
			Text:  t.text,
			Style: style,
		})
	}

	return segments
}

func (sh *SyntaxHighlighter) ApplySyntaxHighlighting(code string) string {
	tokens := sh.tokenize(code)
	var highlightedText strings.Builder

	for _, t := range tokens {
		switch t.tokenType {
		case tokenKeyword:
			highlightedText.WriteString(t.text)
		case tokenString:
			highlightedText.WriteString(t.text)
		case tokenComment:
			highlightedText.WriteString(t.text)
		default:
			highlightedText.WriteString(t.text)
		}
	}

	return highlightedText.String()
}

func (sh *SyntaxHighlighter) tokenize(code string) []token {
	switch sh.language {
	case "go":
		return sh.tokenizeGo(code)
	case "javascript":
		return sh.tokenizeJavaScript(code)
	case "typescript":
		return sh.tokenizeTypeScript(code)
	case "html":
		return sh.tokenizeHTML(code)
	case "css":
		return sh.tokenizeCSS(code)
	default:
		return []token{{tokenPlain, code}}
	}
}

func (sh *SyntaxHighlighter) tokenizeGo(code string) []token {
	var tokens []token

	goKeywords := []string{
		"break", "default", "func", "interface", "select",
		"case", "defer", "go", "map", "struct",
		"chan", "else", "goto", "package", "switch",
		"const", "fallthrough", "if", "range", "type",
		"continue", "for", "import", "return", "var",
	}

	lines := strings.Split(code, "\n")
	for _, line := range lines {
		if len(line) == 0 {
			tokens = append(tokens, token{tokenPlain, "\n"})
			continue
		}

		if strings.HasPrefix(strings.TrimSpace(line), "//") {
			tokens = append(tokens, token{tokenComment, line + "\n"})
			continue
		}

		if strings.Contains(line, "/*") {
			tokens = append(tokens, token{tokenComment, line + "\n"})
			continue
		}

		var currentWord strings.Builder
		for i, char := range line {
			if strings.ContainsRune(" \t():;,{}[]", char) {
				word := currentWord.String()
				if word != "" {
					isKeyword := false
					for _, keyword := range goKeywords {
						if word == keyword {
							tokens = append(tokens, token{tokenKeyword, word})
							isKeyword = true
							break
						}
					}

					if !isKeyword {
						if i < len(line)-1 && line[i] == '(' {
							tokens = append(tokens, token{tokenFunction, word})
						} else {
							if regexp.MustCompile(`^\d+$`).MatchString(word) {
								tokens = append(tokens, token{tokenNumber, word})
							} else {
								tokens = append(tokens, token{tokenPlain, word})
							}
						}
					}
				}
				tokens = append(tokens, token{tokenPlain, string(char)})
				currentWord.Reset()
			} else {
				currentWord.WriteRune(char)
			}
		}

		if currentWord.Len() > 0 {
			word := currentWord.String()
			isKeyword := false
			for _, keyword := range goKeywords {
				if word == keyword {
					tokens = append(tokens, token{tokenKeyword, word})
					isKeyword = true
					break
				}
			}

			if !isKeyword {
				tokens = append(tokens, token{tokenPlain, word})
			}
		}

		tokens = append(tokens, token{tokenPlain, "\n"})
	}

	return tokens
}

func (sh *SyntaxHighlighter) tokenizeJavaScript(code string) []token {
	var tokens []token

	jsKeywords := []string{
		"break", "case", "catch", "class", "const", "continue",
		"debugger", "default", "delete", "do", "else", "export",
		"extends", "finally", "for", "function", "if", "import",
		"in", "instanceof", "new", "return", "super", "switch",
		"this", "throw", "try", "typeof", "var", "void", "while",
		"with", "yield", "let", "async", "await",
	}

	lines := strings.Split(code, "\n")
	for _, line := range lines {
		if len(line) == 0 {
			tokens = append(tokens, token{tokenPlain, "\n"})
			continue
		}

		if strings.HasPrefix(strings.TrimSpace(line), "//") {
			tokens = append(tokens, token{tokenComment, line + "\n"})
			continue
		}

		words := strings.Fields(line)
		for i, word := range words {
			isKeyword := false
			for _, keyword := range jsKeywords {
				if word == keyword {
					tokens = append(tokens, token{tokenKeyword, word})
					isKeyword = true
					break
				}
			}

			if !isKeyword {
				if i < len(words)-1 && strings.HasPrefix(words[i+1], "(") {
					tokens = append(tokens, token{tokenFunction, word})
				} else {
					tokens = append(tokens, token{tokenPlain, word})
				}
			}

			if i < len(words)-1 {
				tokens = append(tokens, token{tokenPlain, " "})
			}
		}

		tokens = append(tokens, token{tokenPlain, "\n"})
	}

	return tokens
}

func (sh *SyntaxHighlighter) tokenizeTypeScript(code string) []token {
	var tokens []token

	tsKeywords := []string{
		"break", "case", "catch", "class", "const", "continue",
		"debugger", "default", "delete", "do", "else", "export",
		"extends", "finally", "for", "function", "if", "import",
		"in", "instanceof", "new", "return", "super", "switch",
		"this", "throw", "try", "typeof", "var", "void", "while",
		"with", "yield", "let", "async", "await",
		"interface", "implements", "namespace", "module", "declare",
		"type", "enum", "private", "protected", "public", "readonly",
	}

	lines := strings.Split(code, "\n")
	for _, line := range lines {
		if len(line) == 0 {
			tokens = append(tokens, token{tokenPlain, "\n"})
			continue
		}

		if strings.HasPrefix(strings.TrimSpace(line), "//") {
			tokens = append(tokens, token{tokenComment, line + "\n"})
			continue
		}

		words := strings.Fields(line)
		for i, word := range words {
			isKeyword := false
			for _, keyword := range tsKeywords {
				if word == keyword {
					tokens = append(tokens, token{tokenKeyword, word})
					isKeyword = true
					break
				}
			}

			if !isKeyword {
				if i < len(words)-1 && strings.HasPrefix(words[i+1], "(") {
					tokens = append(tokens, token{tokenFunction, word})
				} else {
					tokens = append(tokens, token{tokenPlain, word})
				}
			}

			if i < len(words)-1 {
				tokens = append(tokens, token{tokenPlain, " "})
			}
		}

		tokens = append(tokens, token{tokenPlain, "\n"})
	}

	return tokens
}

func (sh *SyntaxHighlighter) tokenizeHTML(code string) []token {
	var tokens []token

	tagRegex := regexp.MustCompile(`<[^>]+>`)
	htmlCode := code

	matches := tagRegex.FindAllStringIndex(htmlCode, -1)
	lastIndex := 0

	for _, match := range matches {
		start, end := match[0], match[1]

		if start > lastIndex {
			tokens = append(tokens, token{tokenPlain, htmlCode[lastIndex:start]})
		}

		tagContent := htmlCode[start:end]
		if strings.HasPrefix(tagContent, "<!--") {
			tokens = append(tokens, token{tokenComment, tagContent})
		} else {
			tokens = append(tokens, token{tokenTag, tagContent})
		}

		lastIndex = end
	}

	if lastIndex < len(htmlCode) {
		tokens = append(tokens, token{tokenPlain, htmlCode[lastIndex:]})
	}

	return tokens
}

func (sh *SyntaxHighlighter) tokenizeCSS(code string) []token {
	var tokens []token

	lines := strings.Split(code, "\n")
	inComment := false

	for _, line := range lines {
		if len(line) == 0 {
			tokens = append(tokens, token{tokenPlain, "\n"})
			continue
		}

		if strings.Contains(line, "/*") {
			inComment = true
			tokens = append(tokens, token{tokenComment, line + "\n"})
			continue
		}

		if strings.Contains(line, "*/") {
			inComment = false
			tokens = append(tokens, token{tokenComment, line + "\n"})
			continue
		}

		if inComment {
			tokens = append(tokens, token{tokenComment, line + "\n"})
			continue
		}

		if strings.Contains(line, "{") && !strings.Contains(line, "}") {
			tokens = append(tokens, token{tokenTag, line + "\n"})
			continue
		}

		if strings.Contains(line, ":") {
			parts := strings.Split(line, ":")
			if len(parts) >= 2 {
				tokens = append(tokens, token{tokenAttribute, parts[0] + ":"})
				tokens = append(tokens, token{tokenPlain, strings.Join(parts[1:], ":") + "\n"})
				continue
			}
		}

		tokens = append(tokens, token{tokenPlain, line + "\n"})
	}

	return tokens
}

func GenerateSyntaxHighlightedRichText(code string, language string) *widget.RichText {
	rt := widget.NewRichText()

	highlighter := NewSyntaxHighlighter(language)
	tokens := highlighter.tokenize(code)

	for _, t := range tokens {
		style := richTextStyleForToken(t.tokenType)
		rt.Segments = append(rt.Segments, &widget.TextSegment{
			Text:  t.text,
			Style: style,
		})
	}

	return rt
}

func richTextStyleForToken(tokenType tokenType) widget.RichTextStyle {
	style := widget.RichTextStyle{
		TextStyle: fyne.TextStyle{Monospace: true},
	}

	switch tokenType {
	case tokenKeyword:
		style.TextStyle.Bold = true
		style.ColorName = theme.ColorNamePrimary
	case tokenString:
		style.TextStyle.Italic = true
		style.ColorName = theme.ColorNameError
	case tokenComment:
		style.TextStyle.Italic = true
		style.ColorName = theme.ColorNameDisabled
	case tokenNumber:
		style.ColorName = theme.ColorNameSuccess
	case tokenFunction:
		style.ColorName = theme.ColorNamePrimary
	case tokenTag:
		style.TextStyle.Bold = true
		style.ColorName = theme.ColorNameButton
	case tokenAttribute:
		style.ColorName = theme.ColorNameWarning
	}

	return style
}
