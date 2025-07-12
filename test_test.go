package renderHTML

import (
	"fmt"
	"testing"
)

// go test -benchmem -run=^$ -bench ^BenchmarkStatistics$ github.com/fabianpallares/renderHTML -count=10
func BenchmarkStatistics(b *testing.B) {
	for i := 0; i < b.N; i++ {
		one()
	}
}

func TestOne(t *testing.T) {
	fmt.Println(one())
}

func TestTwo(t *testing.T) {
	fmt.Println(two())
}

func two() fmt.Stringer {
	return Html().Lang("en").AddContent(
		Head(
			// title
			Title("Testing app"),

			// favicon
			Link().Rel("shortcut icon").Href("media/imgs/favicon.ico"),

			// metadata
			Meta().CharSet("UTF-8"),
			Meta().Name("viewport").Content("width=device-width, initial-scale=1.0"),
			Meta().HttpEquiv("X-UA-Compatible").Content("IE=edge"),
			Meta().HttpEquiv("pragma").Content("no-cache"),
			Meta().HttpEquiv("expires").Content("-1"),

			// css
			Link().Rel("stylesheet").Href("css/styles.css"),

			// javascript
			Script().Defer().Src("js/htmx.min.js"),
		),
		Body(
			H1("Page title"),
			P("Hello world!"),
		),
	)
}

func one() fmt.Stringer {
	return Html(
		Head(
			// title
			Title("Page title"),

			// metainformation
			Meta().CharSet("UTF-8"),
			Meta().Name("viewport").Content("width=device-width, initial-scale=1.0"),
			Meta().HttpEquiv("X-UA-Compatible").Content("IE=edge"),
			Meta().HttpEquiv("pragma").Content("no-cache"),
			Meta().HttpEquiv("expires").Content("-1"),

			// page icon
			Link().Rel("shortcut icon").Href("media/images/favicon.ico"),

			// libraries
			Link().Rel("stylesheet").Href("css/fontawesome/css/all.min.css"),
			Link().Rel("stylesheet").Href("css/bulma/bulma.min.css"),
			Link().Rel("stylesheet").Href("css/styles.css"),

			// custom styles
			Style().Type("text/css").Media("screen").AddContent(`p{color:green;}span{color:blue;background-color: #f2f2f2;}`),
		),
		Body(
			H1("Title").Class("title", "is-1"),
			H2("Sub-Title"),
			Article(
				H3("Bla bla bla").AddContent(RawString(" un poquito mas"), RawString(" y algo mas!")),
			),
			Div(
				Dl(
					Dt("Denim (semigloss finish)"),
					Dd("Ceiling"),
					Hr(),
					Dt("Denim (eggshell finish)"),
					Dt("Evening Sky (eggshell finish)"),
					Dd("Layered on the walls"),
					Span("a text").Id("id-span"),
				),
			),
			Div().Id("id-div").AddContent(
				Select().Id("sel-colors").Name("sel-colors").AddContent(
					Option("Red").Value("red"),
					Option("Green").Value("green").Selected(),
					Option("Blue").Value("blue"),
				),
			),
			A("Save").Href("http://www.foo.com").Target("_blank").AddAttributes("un atributo"),
			Br(),
			Form().Method("post").Action("http://www.foo.com").AddContent(
				Button("Go").Type("submit").Id("btn-go").Name("btn-go"),
				Progress("80").Max("100").Value("10"),
				Textarea().MinLength(10).MaxLength(30).SpellCheck(false),
				Input().Type("text").Id("txt-name").Name("txt-name").Placeholder("Name..."),
			),
			Div(
				Container(
					"Hola!",
					Button("Save").Type("submit").Id("btn-save").Name("btn-save"),
					Span("aditional text").Class("text-one"),
					nil,
					true,
					123,
					456.78,
					"the last text",
					Hr(),
				),
			),
		),
	)
}
