package helpers

import (
	"html"
	"html/template"

	"github.com/frustra/bbcode"
)

var BBCodeCompiler = bbcode.NewCompiler(true, true)

func init() {

	BBCodeCompiler.SetTag("h1", func(node *bbcode.BBCodeNode) (*bbcode.HTMLTag, bool) {
		out := bbcode.NewHTMLTag("")
		out.Name = "h1"
		return out, true
	})
	BBCodeCompiler.SetTag("h2", func(node *bbcode.BBCodeNode) (*bbcode.HTMLTag, bool) {
		out := bbcode.NewHTMLTag("")
		out.Name = "h2"
		return out, true
	})
	BBCodeCompiler.SetTag("h3", func(node *bbcode.BBCodeNode) (*bbcode.HTMLTag, bool) {
		out := bbcode.NewHTMLTag("")
		out.Name = "h3"
		return out, true
	})
	BBCodeCompiler.SetTag("b", func(node *bbcode.BBCodeNode) (*bbcode.HTMLTag, bool) {
		out := bbcode.NewHTMLTag("")
		out.Name = "strong"
		return out, true
	})
	BBCodeCompiler.SetTag("strike", func(node *bbcode.BBCodeNode) (*bbcode.HTMLTag, bool) {
		out := bbcode.NewHTMLTag("")
		out.Name = "span"
		out.Attrs["style"] = "text-decoration: line-through;"
		return out, true
	})
	BBCodeCompiler.SetTag("list", func(node *bbcode.BBCodeNode) (*bbcode.HTMLTag, bool) {
		out := bbcode.NewHTMLTag("")
		out.Name = "ul"
		return out, true
	})
	BBCodeCompiler.SetTag("*", func(node *bbcode.BBCodeNode) (*bbcode.HTMLTag, bool) {
		out := bbcode.NewHTMLTag("")
		out.Name = "li"
		return out, true
	})
	BBCodeCompiler.SetTag("table", func(node *bbcode.BBCodeNode) (*bbcode.HTMLTag, bool) {
		out := bbcode.NewHTMLTag("")
		out.Name = "table"
		out.Attrs["class"] = "table table-hover table-striped table-sm table-nonfluid mb-0"
		return out, true
	})
	BBCodeCompiler.SetTag("tr", func(node *bbcode.BBCodeNode) (*bbcode.HTMLTag, bool) {
		out := bbcode.NewHTMLTag("")
		out.Name = "tr"
		return out, true
	})
	BBCodeCompiler.SetTag("th", func(node *bbcode.BBCodeNode) (*bbcode.HTMLTag, bool) {
		out := bbcode.NewHTMLTag("")
		out.Name = "th"
		return out, true
	})
	BBCodeCompiler.SetTag("td", func(node *bbcode.BBCodeNode) (*bbcode.HTMLTag, bool) {
		out := bbcode.NewHTMLTag("")
		out.Name = "td"
		return out, true
	})
	BBCodeCompiler.SetTag("hr", func(node *bbcode.BBCodeNode) (*bbcode.HTMLTag, bool) {
		out := bbcode.NewHTMLTag("")
		out.Name = "hr"
		return out, true
	})
	BBCodeCompiler.SetTag("previewyoutube", func(node *bbcode.BBCodeNode) (*bbcode.HTMLTag, bool) {
		out := bbcode.NewHTMLTag("")
		out.Name = "iframe"
		out.Attrs["width"] = "560"
		out.Attrs["height"] = "315"
		out.Attrs["src"] = "https://www.youtube.com/embed/" + RegexYouTubeID.FindString(node.GetOpeningTag().Value)
		out.Attrs["frameborder"] = "0"
		out.Attrs["allow"] = "accelerometer; autoplay; encrypted-media; gyroscope; picture-in-picture"
		out.Attrs["allowfullscreen"] = "allowfullscreen"
		return out, true
	})
}

func RenderHTMLAndBBCode(in string) template.HTML {

	return template.HTML(html.UnescapeString(BBCodeCompiler.Compile(in)))
}
