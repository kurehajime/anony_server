package main

import (
	"net/http"
	"regexp"
	"unicode/utf8"

	ipa "github.com/ikawaha/kagome-dict/ipa"
	"github.com/ikawaha/kagome/v2/tokenizer"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Server struct{}

func (h Server) Query(ctx echo.Context) error {
	req := new(Req)
	if err := ctx.Bind(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return ctx.JSON(http.StatusOK, Res{
		Answer: anony(req.Query, false),
	})
}
func (h Server) OpenApi(ctx echo.Context) error {
	ctx.Response().Header().Set("Content-Type", "text/yaml; charset=utf-8")
	return ctx.File("openapi.yaml")
}
func (h Server) AiPlugin(ctx echo.Context) error {
	ctx.Response().Header().Set("Content-Type", "text/json; charset=utf-8")
	return ctx.File("ai-plugin.json")
}

var t *tokenizer.Tokenizer

func init() {
	var err error
	t, err = tokenizer.New(ipa.Dict())
	if err != nil {
		panic(err)
	}
}

func anony(text string, single bool) string {
	tokens := t.Tokenize(text)
	var rText string
	var IniCount int
	for j := 0; j < len(tokens); j++ {
		tk := tokens[j]
		ft := tk.Features()
		if len(ft) > 7 {
			if ft[2] == "人名" && ft[1] == "固有名詞" {
				if IniCount == 0 {
					rText += Word2initial(ft[7])
				} else if IniCount == 1 && single == false {
					rText += "・"
					rText += Word2initial(ft[7])
				}
				IniCount++
			} else {
				rText += tk.Surface
				IniCount = 0
			}
		} else if len(ft) > 0 {
			rText += tk.Surface
			IniCount = 0
		}
	}
	return rText
}

func Word2initial(kana string) string {
	r, _ := utf8.DecodeRune([]byte(kana))
	ini := string(r)

	ini = regexp.MustCompile("[アァ]").ReplaceAllString(ini, "A")
	ini = regexp.MustCompile("[イィ]").ReplaceAllString(ini, "I")
	ini = regexp.MustCompile("[ウゥ]").ReplaceAllString(ini, "U")
	ini = regexp.MustCompile("[エェ]").ReplaceAllString(ini, "E")
	ini = regexp.MustCompile("[オォ]").ReplaceAllString(ini, "O")
	ini = regexp.MustCompile("[カキクケコ]").ReplaceAllString(ini, "K")
	ini = regexp.MustCompile("[サシスセソ]").ReplaceAllString(ini, "S")
	ini = regexp.MustCompile("[タツテト]").ReplaceAllString(ini, "T")
	ini = regexp.MustCompile("[チ]").ReplaceAllString(ini, "T")
	ini = regexp.MustCompile("[ナニヌネノ]").ReplaceAllString(ini, "N")
	ini = regexp.MustCompile("[ハヒヘホ]").ReplaceAllString(ini, "H")
	ini = regexp.MustCompile("[フ]").ReplaceAllString(ini, "F")
	ini = regexp.MustCompile("[マミムメモ]").ReplaceAllString(ini, "M")
	ini = regexp.MustCompile("[ヤユヨ]").ReplaceAllString(ini, "Y")
	ini = regexp.MustCompile("[ラリルレロ]").ReplaceAllString(ini, "R")
	ini = regexp.MustCompile("[ワヲ]").ReplaceAllString(ini, "W")
	ini = regexp.MustCompile("[ン]").ReplaceAllString(ini, "N")
	ini = regexp.MustCompile("[ガギグゲゴ]").ReplaceAllString(ini, "G")
	ini = regexp.MustCompile("[ザズゼゾ]").ReplaceAllString(ini, "Z")
	ini = regexp.MustCompile("[ダヂヅデド]").ReplaceAllString(ini, "D")
	ini = regexp.MustCompile("[ジ]").ReplaceAllString(ini, "J")
	ini = regexp.MustCompile("[パピプペポ]").ReplaceAllString(ini, "P")
	ini = regexp.MustCompile("[バビブベボ]").ReplaceAllString(ini, "B")

	return ini
}

func main() {
	e := echo.New()
	s := Server{}
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.GET, echo.PUT, echo.POST, echo.DELETE},
	}))

	e.GET("/.well-known/ai-plugin.json", s.AiPlugin)
	e.GET("/openapi.yaml", s.OpenApi)

	RegisterHandlers(e, s)

	e.Logger.Fatal(e.Start(":3333"))
}
