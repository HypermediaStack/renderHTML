package renderHTML

import (
	"fmt"
	"strings"
)

// #region AddContentFunc

type addContentFunc[T any] struct {
	el *element
	t  *T
}

func (p *addContentFunc[T]) set(el *element, t *T) {
	p.el = el
	p.t = t
}

// AddContent adds content to the current element.
//
// This attribute isn't represented in HTML. It's a method implemented by all
// elements within the "body" to facilitate writing code by including HTML
// elements or text within other elements.
func (p *addContentFunc[T]) AddContent(content ...any) *T {
	p.el.addContent(content...)
	return p.t
}

// #region ELEMENT
// An "element" is a basic HTML element that can have attributes, classes,
// styles, and content. It can be a void element (without a closing tag) or a
// non-void element (with a closing tag).

// type struct represents the <h1> element.
type element struct {
	tag string

	// An element that does not contain a closing tag is called an "void element"
	//
	// https://developer.mozilla.org/en-US/docs/Learn_web_development/Core/Structuring_content/Basic_HTML_syntax#void_elements
	hasClosingTag bool

	attributes []string
	classes    []string
	styles     []string
	content    []fmt.Stringer
}

// String returns HTML text of the current element.
// func (p *element) String() string {
// 	var s strings.Builder
// 	if p.tag == "" {
// 		s.WriteString(p.getContent())
// 		return s.String()
// 	}

// 	s.WriteString("<")
// 	s.WriteString(p.tag)
// 	s.WriteString(p.getAttributes())
// 	s.WriteString(p.getClasses())
// 	s.WriteString(p.getStyles())

// 	if p.hasClosingTag {
// 		s.WriteString(">")
// 		s.WriteString(p.getContent())
// 		s.WriteString("</")
// 		s.WriteString(p.tag)
// 		s.WriteString(">")
// 	} else {
// 		s.WriteString("/>")
// 	}

// 	return s.String()
// }

func (p *element) String() string {
	if p.tag == "" {
		// it is only used for UntaggedElement
		return fmt.Sprintf("%v", p.getContent())
	}

	if p.hasClosingTag {
		return fmt.Sprintf("<%[1]v%v%v%v>%v</%[1]v>", p.tag, p.getAttributes(), p.getClasses(), p.getStyles(), p.getContent())
	}

	return fmt.Sprintf("<%[1]v%v%v%v/>", p.tag, p.getAttributes(), p.getClasses(), p.getStyles())
}

func (p *element) addAttribute(attr string, value ...any) {
	attr = strings.TrimSpace(attr)
	if attr == "" {
		return
	}
	if value != nil {
		attr = fmt.Sprintf(`%v="%v"`, attr, value[0])
	}

	p.attributes = append(p.attributes, attr)
}

func (p *element) addClasses(class ...string) {
	for _, c := range class {
		c = strings.TrimSpace(c)
		if c != "" {
			p.classes = append(p.classes, c)
		}
	}
}

func (p *element) addStyles(style ...string) {
	for _, s := range style {
		s = strings.TrimSpace(s)
		if s != "" && s[len(s)-1] != ';' {
			s += ";"
		}

		p.styles = append(p.styles, s)
	}
}

func (p *element) addContent(content ...any) {
	for _, value := range content {
		switch v := value.(type) {
		case fmt.Stringer:
			p.content = append(p.content, v)
		case string:
			p.content = append(p.content, &rawStringEntity{content: v})
		case nil:
			continue
		default:
			p.content = append(p.content, rawString("%v", v))
		}
	}
}

func (p *element) getAttributes() string {
	if len(p.attributes) == 0 {
		return ""
	}
	return fmt.Sprintf(` %v`, strings.Join(p.attributes, " "))
}

func (p *element) getClasses() string {
	if len(p.classes) == 0 {
		return ""
	}
	return fmt.Sprintf(` class="%v"`, strings.Join(p.classes, " "))
}

func (p *element) getStyles() string {
	if len(p.styles) == 0 {
		return ""
	}
	return fmt.Sprintf(` style="%v"`, strings.Join(p.styles, " "))
}

func (p *element) getContent() string {
	var s strings.Builder
	for _, el := range p.content {
		s.WriteString(el.String())
	}

	return s.String()
}

func newElement(tag string, hasClosingTag bool, content ...any) *element {
	ne := &element{tag: tag, hasClosingTag: hasClosingTag}
	ne.addContent(content...)
	return ne
}

// #region UNTAGGED ELEMENT

// UntaggedElement is a container for elements. Is a element without tags.
//
// Note: It is not an official HTML element.
type UntaggedElement struct {
	*element
	*addContentFunc[UntaggedElement]
}

// #region Container

// Container creates a container for elements without opening and closing tags.
//
// Note: It is not an official HTML element.
func Container(content ...any) *UntaggedElement {
	ne := newElement("", false, content...)

	var ac = new(addContentFunc[UntaggedElement])
	var el = &UntaggedElement{ne, ac}
	ac.set(ne, el)

	return el
}

// #region Component

// The use of this method is to create server-side web components to group
// different HTML elements.
//
// Note: It is not an official HTML element.
func Component(content ...any) *UntaggedElement {
	ne := newElement("", false, content...)

	var ac = new(addContentFunc[UntaggedElement])
	var el = &UntaggedElement{ne, ac}
	ac.set(ne, el)

	return el
}

// #region MAIN ROOT
// The main and only element in an html document is precisely the
// <html> element.
//

// #region <html>

// HtmlElement represents the <html> element.
type HtmlElement struct {
	*element
	*attrGlobal[HtmlElement]
	*addContentFunc[HtmlElement]
}

// String returns HTML text of the current element.
func (p *HtmlElement) String() string {
	return fmt.Sprintf("<!DOCTYPE html><%[1]v%v%v%v>%v</%[1]v>", p.tag, p.getAttributes(), p.getClasses(), p.getStyles(), p.getContent())
}

// Html represents the root (top-level element) of an HTML document, so it is
// also referred to as the root element. All other elements must be descendants
// of this element.
//
// For more details, see the [documentation]
//
// [documentation]: https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/html
func Html(content ...any) *HtmlElement {
	ne := newElement("html", true, content...)

	var ga = new(attrGlobal[HtmlElement])
	var ac = new(addContentFunc[HtmlElement])
	var el = &HtmlElement{ne, ga, ac}
	ac.set(ne, el)
	ga.set(ne, el)

	return el
}

// #region METADATA
// Metadata contains information about the page. This includes information
// about styles, scripts and data to help software (search engines, browsers,
// etc.) use and render the page. Metadata for styles and scripts may be
// defined in the page or linked to another file that has the

// #region <base>

// BaseElement represents the <base> element.
type BaseElement struct {
	*element
	*attrHref[BaseElement]
	*attrTarget[BaseElement]
}

// Base specifies the base URL to use for all relative URLs in a document.
// There can be only one <base> element in a document.
//
// Warning: a <base> element must have an href attribute, a target attribute, or
// both. If at least one of these attributes are specified, the <base> element
// must come before other elements with attribute values that are URLs, such as
// a <link>'s href attribute.
//
// For more details, see the [documentation]
//
// [documentation]: https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/base
func Base() *BaseElement {
	ne := newElement("base", false)

	var a = new(attrHref[BaseElement])
	var b = new(attrTarget[BaseElement])
	var el = &BaseElement{ne, a, b}
	a.set(ne, el)
	b.set(ne, el)

	return el
}

// #region <head>

// HeadElement represents the <head> element.
type HeadElement struct {
	*element
	*addContentFunc[HeadElement]
}

// Head contains machine-readable information (metadata) about the document, like
// its title, scripts, and style sheets.
//
// For more details, see the [documentation]
//
// [documentation]: https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/head
func Head(content ...any) *HeadElement {
	ne := newElement("head", true, content...)

	var ac = new(addContentFunc[HeadElement])
	var el = &HeadElement{ne, ac}
	ac.set(ne, el)

	return el
}

// #region <link>

// LinkElement represents the <link> element.
type LinkElement struct {
	*element
	*attrAs[LinkElement]
	*attrCrossOrigin[LinkElement]
	*attrDisabled[LinkElement]
	*attrHref[LinkElement]
	*attrHrefLang[LinkElement]
	*attrIntegrity[LinkElement]
	*attrMedia[LinkElement]
	*attrRel[LinkElement]
	*attrSizes[LinkElement]
	*attrType[LinkElement]
}

// Link specifies relationships between the current document and an external
// resource.
//
// This element is most commonly used to link to CSS but is also
// used to establish site icons (both "favicon" style icons and icons for the
// home screen and apps on mobile devices) among other things.
//
// For more details, see the [documentation]
//
// [documentation]: https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/link
func Link() *LinkElement {
	ne := newElement("link", false)

	var a = new(attrAs[LinkElement])
	var b = new(attrCrossOrigin[LinkElement])
	var c = new(attrDisabled[LinkElement])
	var d = new(attrHref[LinkElement])
	var e = new(attrHrefLang[LinkElement])
	var f = new(attrIntegrity[LinkElement])
	var g = new(attrMedia[LinkElement])
	var h = new(attrRel[LinkElement])
	var i = new(attrSizes[LinkElement])
	var j = new(attrType[LinkElement])

	var el = &LinkElement{ne, a, b, c, d, e, f, g, h, i, j}

	a.set(ne, el)
	b.set(ne, el)
	c.set(ne, el)
	d.set(ne, el)
	e.set(ne, el)
	f.set(ne, el)
	g.set(ne, el)
	h.set(ne, el)
	i.set(ne, el)
	j.set(ne, el)

	return el
}

// #region <meta>

// MetaElement represents the <meta> element.
type MetaElement struct {
	*element
	*attrCharSet[MetaElement]
	*attrMetaContent[MetaElement]
	*attrHttpEquiv[MetaElement]
	*attrName[MetaElement]
}

// Meta represents metadata that cannot be represented by other HTML
// meta-related elements, like <base>, <link>, <script>, <style> and <title>.
//
// For more details, see the [documentation]
//
// [documentation]: https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/meta
func Meta() *MetaElement {
	ne := newElement("meta", false)

	var a = new(attrCharSet[MetaElement])
	var b = new(attrMetaContent[MetaElement])
	var c = new(attrHttpEquiv[MetaElement])
	var d = new(attrName[MetaElement])
	var el = &MetaElement{ne, a, b, c, d}
	a.set(ne, el)
	b.set(ne, el)
	c.set(ne, el)
	d.set(ne, el)

	return el
}

// #region <style>

// StyleElement represents the <style> element.
type StyleElement struct {
	*element
	*attrMedia[StyleElement]
	*attrType[StyleElement]
	*addContentFunc[StyleElement]
}

// Style contains style information for a document or part of a document. It
// contains CSS, which is applied to the contents of the document containing
// this element.
//
// For more details, see the [documentation]
//
// [documentation]: https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/style
func Style(content ...any) *StyleElement {
	ne := newElement("style", true, content...)

	var a = new(attrMedia[StyleElement])
	var b = new(attrType[StyleElement])
	var ac = new(addContentFunc[StyleElement])
	var el = &StyleElement{ne, a, b, ac}
	a.set(ne, el)
	b.set(ne, el)
	ac.set(ne, el)

	return el
}

// #region <title>

// TitleElement represents the <title> element.
type TitleElement struct {
	*element
	*addContentFunc[TitleElement]
}

// Title defines the document's title that is shown in a browser's title bar or
// a page's tab. It only contains text; tags within the element are ignored.
//
// For more details, see the [documentation]
//
// [documentation]: https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/title
func Title(content ...any) *TitleElement {
	ne := newElement("title", true, content...)

	var ac = new(addContentFunc[TitleElement])
	var el = &TitleElement{ne, ac}
	ac.set(ne, el)

	return el

}

// #region SECTIONING
// Content sectioning elements allow you to organize the document content into
// logical pieces. Use the sectioning elements to create a broad outline for
// your page content, including header and footer navigation, and heading
// elements to identify sections of content.

// #region <body>

// BodyElement represents the <body> element.
type BodyElement struct {
	*element
	*attrGlobal[BodyElement]
	*attrExternalAttributes[BodyElement]
	*attrOn[BodyElement]
	*addContentFunc[BodyElement]
}

// Body represents the content of an HTML document. There can be only one such
// element in a document.
//
// For more details, see the [documentation]
//
// [documentation]: https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/body
func Body(content ...any) *BodyElement {
	ne := newElement("body", true, content...)

	var ga = new(attrGlobal[BodyElement])
	var ea = new(attrExternalAttributes[BodyElement])
	var on = new(attrOn[BodyElement])
	var ac = new(addContentFunc[BodyElement])
	var el = &BodyElement{ne, ga, ea, on, ac}
	ga.set(ne, el)
	ea.set(ne, el)
	on.set(ne, el)
	ac.set(ne, el)

	return el
}

// #region <address>

// AddressElement represents the <address> element.
type AddressElement struct {
	*element
	*attrGlobal[AddressElement]
	*attrExternalAttributes[AddressElement]
	*attrOn[AddressElement]
	*addContentFunc[AddressElement]
}

// Address indicates that the enclosed HTML provides contact information for a
// person or people, or for an organization.
//
// For more details, see the [documentation]
//
// [documentation]: https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/address
func Address(content ...any) *AddressElement {
	ne := newElement("address", true, content...)

	var ga = new(attrGlobal[AddressElement])
	var ea = new(attrExternalAttributes[AddressElement])
	var on = new(attrOn[AddressElement])
	var ac = new(addContentFunc[AddressElement])
	var el = &AddressElement{ne, ga, ea, on, ac}
	ga.set(ne, el)
	ea.set(ne, el)
	on.set(ne, el)
	ac.set(ne, el)

	return el
}

// #region <article>

// ArticleElement represents the <article> element.
type ArticleElement struct {
	*element
	*attrGlobal[ArticleElement]
	*attrExternalAttributes[ArticleElement]
	*attrOn[ArticleElement]
	*addContentFunc[ArticleElement]
}

// Article represents a self-contained composition in a document, page,
// application, or site, which is intended to be independently distributable or
// reusable (e.g., in syndication). Examples include a forum post, a magazine
// or newspaper article, a blog entry, a product card, a user-submitted
// comment, an interactive widget or gadget, or any other independent item of
// content.
//
// For more details, see the [documentation]
//
// [documentation]: https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/article
func Article(content ...any) *ArticleElement {
	ne := newElement("article", true, content...)

	var ga = new(attrGlobal[ArticleElement])
	var ea = new(attrExternalAttributes[ArticleElement])
	var on = new(attrOn[ArticleElement])
	var ac = new(addContentFunc[ArticleElement])
	var el = &ArticleElement{ne, ga, ea, on, ac}
	ga.set(ne, el)
	ea.set(ne, el)
	on.set(ne, el)
	ac.set(ne, el)

	return el
}

// #region <aside>

// AsideElement represents the <aside> element.
type AsideElement struct {
	*element
	*attrGlobal[AsideElement]
	*attrExternalAttributes[AsideElement]
	*attrOn[AsideElement]
	*addContentFunc[AsideElement]
}

// Aside represents a portion of a document whose content is only indirectly
// related to the document's main content. Asides are frequently presented as
// sidebars or call-out boxes.
//
// For more details, see the [documentation]
//
// [documentation]: https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/aside
func Aside(content ...any) *AsideElement {
	ne := newElement("aside", true, content...)

	var ga = new(attrGlobal[AsideElement])
	var ea = new(attrExternalAttributes[AsideElement])
	var on = new(attrOn[AsideElement])
	var ac = new(addContentFunc[AsideElement])
	var el = &AsideElement{ne, ga, ea, on, ac}
	ga.set(ne, el)
	ea.set(ne, el)
	on.set(ne, el)
	ac.set(ne, el)

	return el
}

// #region <footer>

// FooterElement represents the <footer> element.
type FooterElement struct {
	*element
	*attrGlobal[FooterElement]
	*attrExternalAttributes[FooterElement]
	*attrOn[FooterElement]
	*addContentFunc[FooterElement]
}

// Footer represents a footer for its nearest ancestor sectioning content or
// sectioning root element. A <footer> typically contains information about
// the author of the section, copyright data, or links to related documents.
//
// For more details, see the [documentation]
//
// [documentation]: https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/footer
func Footer(content ...any) *FooterElement {
	ne := newElement("footer", true, content...)

	var ga = new(attrGlobal[FooterElement])
	var ea = new(attrExternalAttributes[FooterElement])
	var on = new(attrOn[FooterElement])
	var ac = new(addContentFunc[FooterElement])
	var el = &FooterElement{ne, ga, ea, on, ac}
	ga.set(ne, el)
	ea.set(ne, el)
	on.set(ne, el)
	ac.set(ne, el)

	return el
}

// #region <header>

// HeaderElement represents the <header> element.
type HeaderElement struct {
	*element
	*attrGlobal[HeaderElement]
	*attrExternalAttributes[HeaderElement]
	*attrOn[HeaderElement]
	*addContentFunc[HeaderElement]
}

// Header represents introductory content, typically a group of introductory or
// navigational aids. It may contain some heading elements but also a logo, a
// search form, an author name, and other elements.
//
// For more details, see the [documentation]
//
// [documentation]: https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/header
func Header(content ...any) *HeaderElement {
	ne := newElement("header", true, content...)

	var ga = new(attrGlobal[HeaderElement])
	var ea = new(attrExternalAttributes[HeaderElement])
	var on = new(attrOn[HeaderElement])
	var ac = new(addContentFunc[HeaderElement])
	var el = &HeaderElement{ne, ga, ea, on, ac}
	ga.set(ne, el)
	ea.set(ne, el)
	on.set(ne, el)
	ac.set(ne, el)

	return el
}

// #region <h1>

// H1Element represents the <h1> element.
type H1Element struct {
	*element
	*attrGlobal[H1Element]
	*attrExternalAttributes[H1Element]
	*attrOn[H1Element]
	*addContentFunc[H1Element]
}

// H1 represents one of six levels of section headings. <h1> is the highest
// section level and <h6> is the lowest.
//
// For more details, see the [documentation]
//
// [documentation]: https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/h1
func H1(content ...any) *H1Element {
	ne := newElement("h1", true, content...)

	var ga = new(attrGlobal[H1Element])
	var ea = new(attrExternalAttributes[H1Element])
	var on = new(attrOn[H1Element])
	var ac = new(addContentFunc[H1Element])
	var el = &H1Element{ne, ga, ea, on, ac}
	ga.set(ne, el)
	ea.set(ne, el)
	on.set(ne, el)
	ac.set(ne, el)

	return el
}

// #region <h2>

// H2Element represents the <h2> element.
type H2Element struct {
	*element
	*attrGlobal[H2Element]
	*attrExternalAttributes[H2Element]
	*attrOn[H2Element]
	*addContentFunc[H2Element]
}

// H2 represents one of six levels of section headings. <h1> is the highest
// section level and <h6> is the lowest.
//
// For more details, see the [documentation]
//
// [documentation]: https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/h2
func H2(content ...any) *H2Element {
	ne := newElement("h2", true, content...)

	var ga = new(attrGlobal[H2Element])
	var ea = new(attrExternalAttributes[H2Element])
	var on = new(attrOn[H2Element])
	var ac = new(addContentFunc[H2Element])
	var el = &H2Element{ne, ga, ea, on, ac}
	ga.set(ne, el)
	ea.set(ne, el)
	on.set(ne, el)
	ac.set(ne, el)

	return el
}

// #region <h3>

// H3Element represents the <h3> element.
type H3Element struct {
	*element
	*attrGlobal[H3Element]
	*attrExternalAttributes[H3Element]
	*attrOn[H3Element]
	*addContentFunc[H3Element]
}

// H3 represents one of six levels of section headings. <h1> is the highest
// section level and <h6> is the lowest.
//
// For more details, see the [documentation]
//
// [documentation]: https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/h3
func H3(content ...any) *H3Element {
	ne := newElement("h3", true, content...)

	var ga = new(attrGlobal[H3Element])
	var ea = new(attrExternalAttributes[H3Element])
	var on = new(attrOn[H3Element])
	var ac = new(addContentFunc[H3Element])
	var el = &H3Element{ne, ga, ea, on, ac}
	ga.set(ne, el)
	ea.set(ne, el)
	on.set(ne, el)
	ac.set(ne, el)

	return el
}

// #region <h4>

// H4Element represents the <h4> element.
type H4Element struct {
	*element
	*attrGlobal[H4Element]
	*attrExternalAttributes[H4Element]
	*attrOn[H4Element]
	*addContentFunc[H4Element]
}

// H4 represents one of six levels of section headings. <h1> is the highest
// section level and <h6> is the lowest.
//
// For more details, see the [documentation]
//
// [documentation]: https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/h4
func H4(content ...any) *H4Element {
	ne := newElement("h4", true, content...)

	var ga = new(attrGlobal[H4Element])
	var ea = new(attrExternalAttributes[H4Element])
	var on = new(attrOn[H4Element])
	var ac = new(addContentFunc[H4Element])
	var el = &H4Element{ne, ga, ea, on, ac}
	ga.set(ne, el)
	ea.set(ne, el)
	on.set(ne, el)
	ac.set(ne, el)

	return el
}

// #region <h5>

// H5Element represents the <h5> element.
type H5Element struct {
	*element
	*attrGlobal[H5Element]
	*attrExternalAttributes[H5Element]
	*attrOn[H5Element]
	*addContentFunc[H5Element]
}

// H5 represents one of six levels of section headings. <h1> is the highest
// section level and <h6> is the lowest.
//
// For more details, see the [documentation]
//
// [documentation]: https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/h5
func H5(content ...any) *H5Element {
	ne := newElement("h5", true, content...)

	var ga = new(attrGlobal[H5Element])
	var ea = new(attrExternalAttributes[H5Element])
	var on = new(attrOn[H5Element])
	var ac = new(addContentFunc[H5Element])
	var el = &H5Element{ne, ga, ea, on, ac}
	ga.set(ne, el)
	ea.set(ne, el)
	on.set(ne, el)
	ac.set(ne, el)

	return el
}

// #region <h6>

// H6Element represents the <h6> element.
type H6Element struct {
	*element
	*attrGlobal[H6Element]
	*attrExternalAttributes[H6Element]
	*attrOn[H6Element]
	*addContentFunc[H6Element]
}

// H6 represents one of six levels of section headings. <h1> is the highest
// section level and <h6> is the lowest.
//
// For more details, see the [documentation]
//
// [documentation]: https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/h6
func H6(content ...any) *H6Element {
	ne := newElement("h6", true, content...)

	var ga = new(attrGlobal[H6Element])
	var ea = new(attrExternalAttributes[H6Element])
	var on = new(attrOn[H6Element])
	var ac = new(addContentFunc[H6Element])
	var el = &H6Element{ne, ga, ea, on, ac}
	ga.set(ne, el)
	ea.set(ne, el)
	on.set(ne, el)
	ac.set(ne, el)

	return el
}

// #region <hgroup>

// HgroupElement represents the <hgroup> element.
type HgroupElement struct {
	*element
	*attrGlobal[HgroupElement]
	*attrExternalAttributes[HgroupElement]
	*attrOn[HgroupElement]
	*addContentFunc[HgroupElement]
}

// Hgroup represents a heading grouped with any secondary content, such as
// subheadings, an alternative title, or a tagline.
//
// The <hgroup> element allows the grouping of a heading with any secondary
// content, such as subheadings, an alternative title, or tagline. Each of
// these types of content represented as a <p> element within the <hgroup>.
//
// The <hgroup> itself has no impact on the document outline of a web page.
// Rather, the single allowed heading within the <hgroup> contributes to the
// document outline.
//
// For more details, see the [documentation]
//
// [documentation]: https://developer.mozilla.org/en-US/docs/Web/HTML/Element/hgroup
func Hgroup(content ...any) *HgroupElement {
	ne := newElement("hgroup", true, content...)

	var ga = new(attrGlobal[HgroupElement])
	var ea = new(attrExternalAttributes[HgroupElement])
	var on = new(attrOn[HgroupElement])
	var ac = new(addContentFunc[HgroupElement])
	var el = &HgroupElement{ne, ga, ea, on, ac}
	ga.set(ne, el)
	ea.set(ne, el)
	on.set(ne, el)
	ac.set(ne, el)

	return el
}

// #region <main>

// MainElement represents the <main> element.
type MainElement struct {
	*element
	*attrGlobal[MainElement]
	*attrExternalAttributes[MainElement]
	*attrOn[MainElement]
	*addContentFunc[MainElement]
}

// Main represents the dominant content of the body of a document. The main
// content area consists of content that is directly related to or expands
// upon the central topic of a document, or the central functionality of an
// application.
//
// The content of a <main> element should be unique to the document. Content
// that is repeated across a set of documents or document sections such as
// sidebars, navigation links, copyright information, site logos, and search
// forms shouldn't be included unless the search form is the main function of
// the page.
//
// <main> doesn't contribute to the document's outline; that is, unlike
// elements such as <body>, headings such as h2, and such, <main> doesn't
// affect the DOM's concept of the structure of the page. It's strictly
// informative.
//
// For more details, see the [documentation]
//
// [documentation]: https://developer.mozilla.org/en-US/docs/Web/HTML/Element/main
func Main(content ...any) *MainElement {
	ne := newElement("main", true, content...)

	var ga = new(attrGlobal[MainElement])
	var ea = new(attrExternalAttributes[MainElement])
	var on = new(attrOn[MainElement])
	var ac = new(addContentFunc[MainElement])
	var el = &MainElement{ne, ga, ea, on, ac}
	ga.set(ne, el)
	ea.set(ne, el)
	on.set(ne, el)
	ac.set(ne, el)

	return el
}

// #region <nav>

// NavElement represents the <nav> element.
type NavElement struct {
	*element
	*attrGlobal[NavElement]
	*attrExternalAttributes[NavElement]
	*attrOn[NavElement]
	*addContentFunc[NavElement]
}

// Nav represents a section of a page whose purpose is to provide navigation
// links, either within the current document or to other documents. Common
// examples of navigation sections are menus, tables of contents, and indexes.
//
//   - It's not necessary for all links to be contained in a <nav> element. <nav> is intended only for a major block of navigation links; typically the <footer> element often has a list of links that don't need to be in a <nav> element.
//   - A document may have several <nav> elements, for example, one for site navigation and one for intra-page navigation. aria-labelledby can be used in such case to promote accessibility, see example.
//   - User agents, such as screen readers targeting disabled users, can use this element to determine whether to omit the initial rendering of navigation-only content.
//
// For more details, see the [documentation]
//
// [documentation]: https://developer.mozilla.org/en-US/docs/Web/HTML/Element/nav
func Nav(content ...any) *NavElement {
	ne := newElement("nav", true, content...)

	var ga = new(attrGlobal[NavElement])
	var ea = new(attrExternalAttributes[NavElement])
	var on = new(attrOn[NavElement])
	var ac = new(addContentFunc[NavElement])
	var el = &NavElement{ne, ga, ea, on, ac}
	ga.set(ne, el)
	ea.set(ne, el)
	on.set(ne, el)
	ac.set(ne, el)

	return el
}

// #region <section>

// SectionElement represents the <section> element.
type SectionElement struct {
	*element
	*attrGlobal[SectionElement]
	*attrExternalAttributes[SectionElement]
	*attrOn[SectionElement]
	*addContentFunc[SectionElement]
}

// Section represents a generic standalone section of a document, which doesn't
// have a more specific semantic element to represent it. Sections should
// always have a heading, with very few exceptions.
//
// Is a generic sectioning element, and should only be used if there isn't a
// more specific element to represent it. As an example, a navigation menu
// should be wrapped in a <nav> element, but a list of search results or a map
// display and its controls don't have specific elements, and could be put
// inside a <section>.
//
// For more details, see the [documentation]
//
// [documentation]: https://developer.mozilla.org/en-US/docs/Web/HTML/Element/section
func Section(content ...any) *SectionElement {
	ne := newElement("section", true, content...)

	var ga = new(attrGlobal[SectionElement])
	var ea = new(attrExternalAttributes[SectionElement])
	var on = new(attrOn[SectionElement])
	var ac = new(addContentFunc[SectionElement])
	var el = &SectionElement{ne, ga, ea, on, ac}
	ga.set(ne, el)
	ea.set(ne, el)
	on.set(ne, el)
	ac.set(ne, el)

	return el
}

// #region <search>

// SearchElement represents the <search> element.
type SearchElement struct {
	*element
	*attrGlobal[SearchElement]
	*attrExternalAttributes[SearchElement]
	*attrOn[SearchElement]
	*addContentFunc[SearchElement]
}

// Search represents a part that contains a set of form controls or other
// content related to performing a search or filtering operation.
//
// The <search> element is not for presenting search results. Rather, search
// or filtered results should be presented as part of the main content of
// that web page. That said, suggestions and links that are part of "quick
// search" functionality within the search or filtering functionality are
// appropriately nested within the contents of the <search> element as they
// are search features.
//
// For more details, see the [documentation]
//
// [documentation]: https://developer.mozilla.org/en-US/docs/Web/HTML/Element/search
func Search(content ...any) *SearchElement {
	ne := newElement("search", true, content...)

	var ga = new(attrGlobal[SearchElement])
	var ea = new(attrExternalAttributes[SearchElement])
	var on = new(attrOn[SearchElement])
	var ac = new(addContentFunc[SearchElement])
	var el = &SearchElement{ne, ga, ea, on, ac}
	ga.set(ne, el)
	ea.set(ne, el)
	on.set(ne, el)
	ac.set(ne, el)

	return el
}

// #region TEXT CONTENT
// Use HTML text content elements to organize blocks or sections of content
// placed between the opening <body> and closing </body> tags. Important for
// accessibility and SEO, these elements identify the purpose or structure of
// that content.

// #region <blockquote>

// BlockquoteElement represents the <blockquote> element.
type BlockquoteElement struct {
	*element
	*attrCite[BlockquoteElement]
	*attrGlobal[BlockquoteElement]
	*attrExternalAttributes[BlockquoteElement]
	*attrOn[BlockquoteElement]
	*addContentFunc[BlockquoteElement]
}

// Blockquote indicates that the enclosed text is an extended quotation. Usually,
// this is rendered visually by indentation. A URL for the source of the
// quotation may be given using the cite attribute, while a text representation
// of the source can be given using the <cite> element.
//
// For more details, see the [documentation]
//
// [documentation]: https://developer.mozilla.org/en-US/docs/Web/HTML/Element/blockquote
func Blockquote(content ...any) *BlockquoteElement {
	ne := newElement("blockquote", true, content...)

	var a = new(attrCite[BlockquoteElement])
	var ga = new(attrGlobal[BlockquoteElement])
	var ea = new(attrExternalAttributes[BlockquoteElement])
	var on = new(attrOn[BlockquoteElement])
	var ac = new(addContentFunc[BlockquoteElement])
	var el = &BlockquoteElement{ne, a, ga, ea, on, ac}
	a.set(ne, el)
	ga.set(ne, el)
	ea.set(ne, el)
	on.set(ne, el)
	ac.set(ne, el)

	return el
}

// #region <dd>

// DdElement represents the <dd> element.
type DdElement struct {
	*element
	*attrGlobal[DdElement]
	*attrExternalAttributes[DdElement]
	*attrOn[DdElement]
	*addContentFunc[DdElement]
}

// Dd provides the description, definition, or value for the preceding
// term (<dt>) in a description list (<dl>).
//
// For more details, see the [documentation]
//
// [documentation]: https://developer.mozilla.org/en-US/docs/Web/HTML/Element/dd
func Dd(content ...any) *DdElement {
	ne := newElement("dd", true, content...)

	var ga = new(attrGlobal[DdElement])
	var ea = new(attrExternalAttributes[DdElement])
	var on = new(attrOn[DdElement])
	var ac = new(addContentFunc[DdElement])
	var el = &DdElement{ne, ga, ea, on, ac}
	ga.set(ne, el)
	ea.set(ne, el)
	on.set(ne, el)
	ac.set(ne, el)

	return el
}

// #region <div>

// DivElement represents the <div> element.
type DivElement struct {
	*element
	*attrGlobal[DivElement]
	*attrExternalAttributes[DivElement]
	*attrOn[DivElement]
	*addContentFunc[DivElement]
}

// Div is the generic container for flow content. It has no effect on the
// content or layout until styled in some way using CSS (e.g., styling is
// directly applied to it, or some kind of layout model like flexbox is
// applied to its parent element).
//
// For more details, see the [documentation]
//
// [documentation]: https://developer.mozilla.org/en-US/docs/Web/HTML/Element/div
func Div(content ...any) *DivElement {
	ne := newElement("div", true, content...)

	var ga = new(attrGlobal[DivElement])
	var ea = new(attrExternalAttributes[DivElement])
	var on = new(attrOn[DivElement])
	var ac = new(addContentFunc[DivElement])
	var el = &DivElement{ne, ga, ea, on, ac}
	ga.set(ne, el)
	ea.set(ne, el)
	on.set(ne, el)
	ac.set(ne, el)

	return el
}

// #region <dl>

// DlElement represents the <dl> element.
type DlElement struct {
	*element
	*attrGlobal[DlElement]
	*attrExternalAttributes[DlElement]
	*attrOn[DlElement]
	*addContentFunc[DlElement]
}

// Dl represents a description list. The element encloses a list of groups of
// terms (specified using the <dt> element) and descriptions (provided by <dd>
// elements). Common uses for this element are to implement a glossary or to
// display metadata (a list of key-value pairs).
//
// For more details, see the [documentation]
//
// [documentation]: https://developer.mozilla.org/en-US/docs/Web/HTML/Element/dl
func Dl(content ...any) *DlElement {
	ne := newElement("dl", true, content...)

	var ga = new(attrGlobal[DlElement])
	var ea = new(attrExternalAttributes[DlElement])
	var on = new(attrOn[DlElement])
	var ac = new(addContentFunc[DlElement])
	var el = &DlElement{ne, ga, ea, on, ac}
	ga.set(ne, el)
	ea.set(ne, el)
	on.set(ne, el)
	ac.set(ne, el)

	return el
}

// #region <dt>

// DtElement represents the <dt> element.
type DtElement struct {
	*element
	*attrGlobal[DtElement]
	*attrExternalAttributes[DtElement]
	*attrOn[DtElement]
	*addContentFunc[DtElement]
}

// Dt specifies a term in a description or definition list, and as such must be
// used inside a <dl> element. It is usually followed by a <dd> element; however,
// multiple <dt> elements in a row indicate several terms that are all defined
// by the immediate next <dd> element.
//
// For more details, see the [documentation]
//
// [documentation]: https://developer.mozilla.org/en-US/docs/Web/HTML/Element/dt
func Dt(content ...any) *DtElement {
	ne := newElement("dt", true, content...)

	var ga = new(attrGlobal[DtElement])
	var ea = new(attrExternalAttributes[DtElement])
	var on = new(attrOn[DtElement])
	var ac = new(addContentFunc[DtElement])
	var el = &DtElement{ne, ga, ea, on, ac}
	ga.set(ne, el)
	ea.set(ne, el)
	on.set(ne, el)
	ac.set(ne, el)

	return el
}

// #region <figcaption>

// FigCaptionElement represents the <figcaption> element.
type FigCaptionElement struct {
	*element
	*attrGlobal[FigCaptionElement]
	*attrExternalAttributes[FigCaptionElement]
	*attrOn[FigCaptionElement]
	*addContentFunc[FigCaptionElement]
}

// FigCaption represents a caption or legend describing the rest of the contents
// of its parent <figure> element.
//
// For more details, see the [documentation]
//
// [documentation]: https://developer.mozilla.org/en-US/docs/Web/HTML/Element/figcaption
func FigCaption(content ...any) *FigCaptionElement {
	ne := newElement("figcaption", true, content...)

	var ga = new(attrGlobal[FigCaptionElement])
	var ea = new(attrExternalAttributes[FigCaptionElement])
	var on = new(attrOn[FigCaptionElement])
	var ac = new(addContentFunc[FigCaptionElement])
	var el = &FigCaptionElement{ne, ga, ea, on, ac}
	ga.set(ne, el)
	ea.set(ne, el)
	on.set(ne, el)
	ac.set(ne, el)

	return el
}

// #region <figure>

// FigureElement represents the <figure> element.
type FigureElement struct {
	*element
	*attrGlobal[FigureElement]
	*attrExternalAttributes[FigureElement]
	*attrOn[FigureElement]
	*addContentFunc[FigureElement]
}

// Figure represents self-contained content, potentially with an optional
// caption, which is specified using the <figcaption> element. The figure, its
// caption, and its contents are referenced as a single unit.
//
// Example:
//
//	<figure>
//	  <img src="/media/cc0-images/elephant-660-480.jpg" alt="Elephant at sunset" />
//	  <figcaption>An elephant at sunset</figcaption>
//	</figure>
//
// For more details, see the [documentation]
//
// [documentation]: https://developer.mozilla.org/en-US/docs/Web/HTML/Element/figure
func Figure(content ...any) *FigureElement {
	ne := newElement("figure", true, content...)

	var ga = new(attrGlobal[FigureElement])
	var ea = new(attrExternalAttributes[FigureElement])
	var on = new(attrOn[FigureElement])
	var ac = new(addContentFunc[FigureElement])
	var el = &FigureElement{ne, ga, ea, on, ac}
	ga.set(ne, el)
	ea.set(ne, el)
	on.set(ne, el)
	ac.set(ne, el)

	return el
}

// #region <hr>

// HrElement represents the <hr> element.
type HrElement struct {
	*element
	*attrGlobal[HrElement]
	*attrExternalAttributes[HrElement]
	*attrOn[HrElement]
}

// Hr represents a thematic break between paragraph-level elements: for
// example, a change of scene in a story, or a shift of topic within a section.
//
// Historically, this has been presented as a horizontal rule or line. While it
// may still be displayed as a horizontal rule in visual browsers, this element
// is now defined in semantic terms, rather than presentational terms, so if you
// wish to draw a horizontal line, you should do so using appropriate CSS.
//
// For more details, see the [documentation]
//
// [documentation]: https://developer.mozilla.org/en-US/docs/Web/HTML/Element/hr
func Hr() *HrElement {
	ne := newElement("hr", false)
	var ga = new(attrGlobal[HrElement])
	var ea = new(attrExternalAttributes[HrElement])
	var on = new(attrOn[HrElement])
	var el = &HrElement{ne, ga, ea, on}
	ga.set(ne, el)
	ea.set(ne, el)
	on.set(ne, el)

	return el
}

// #region <li>

// LiElement represents the <li> element.
type LiElement struct {
	*element
	*attrGlobal[LiElement]
	*attrExternalAttributes[LiElement]
	*attrOn[LiElement]
	*addContentFunc[LiElement]
}

// Value specifies the value of an element.
//
// This integer attribute indicates the current ordinal value of the list item
// as defined by the <ol> element. The only allowed value for this attribute is
// a number, even if the list is displayed with Roman numerals or letters. List
// items that follow this one continue numbering from the value set. The value
// attribute has no meaning for unordered lists (<ul>) or for menus (<menu>).
func (p *LiElement) Value(value int) *LiElement {
	p.addAttribute("value", value)
	return p
}

// Li represents an item in a list. It must be contained in a parent element:
// an ordered list (<ol>), an unordered list (<ul>), or a menu (<menu>). In
// menus and unordered lists, list items are usually displayed using bullet
// points. In ordered lists, they are usually displayed with an ascending
// counter on the left, such as a number or letter.
//
// For more details, see the [documentation]
//
// [documentation]: https://developer.mozilla.org/en-US/docs/Web/HTML/Element/li
func Li(content ...any) *LiElement {
	ne := newElement("li", true, content...)

	var ga = new(attrGlobal[LiElement])
	var ea = new(attrExternalAttributes[LiElement])
	var on = new(attrOn[LiElement])
	var ac = new(addContentFunc[LiElement])
	var el = &LiElement{ne, ga, ea, on, ac}
	ga.set(ne, el)
	ea.set(ne, el)
	on.set(ne, el)
	ac.set(ne, el)

	return el
}

// #region <menu>

// MenuElement represents the <menu> element.
type MenuElement struct {
	*element
	*attrGlobal[MenuElement]
	*attrExternalAttributes[MenuElement]
	*attrOn[MenuElement]
	*addContentFunc[MenuElement]
}

// Menu is an semantic alternative to <ul>, but treated by browsers (and
// exposed through the accessibility tree) as no different than <ul>.
//
// It represents an unordered list of items (which are represented by <li>
// elements).
//
// Example:
//
//	<menu>
//	  <li><button onclick="copy()">Copy</button></li>
//	  <li><button onclick="cut()">Cut</button></li>
//	  <li><button onclick="paste()">Paste</button></li>
//	</menu>
//
// For more details, see the [documentation]
//
// [documentation]: https://developer.mozilla.org/en-US/docs/Web/HTML/Element/menu
func Menu(content ...any) *MenuElement {
	ne := newElement("menu", true, content...)

	var ga = new(attrGlobal[MenuElement])
	var ea = new(attrExternalAttributes[MenuElement])
	var on = new(attrOn[MenuElement])
	var ac = new(addContentFunc[MenuElement])
	var el = &MenuElement{ne, ga, ea, on, ac}
	ga.set(ne, el)
	ea.set(ne, el)
	on.set(ne, el)
	ac.set(ne, el)

	return el
}

// #region <ol>

// OlElement represents the <ol> element.
type OlElement struct {
	*element
	*attrReversed[OlElement]
	*attrStart[OlElement]
	*attrGlobal[OlElement]
	*attrExternalAttributes[OlElement]
	*attrOn[OlElement]
	*addContentFunc[OlElement]
}

// Type sets the numbering type:
//   - a for lowercase letters.
//   - A for uppercase letters.
//   - i for lowercase Roman numerals.
//   - I for uppercase Roman numerals.
//   - 1 for numbers (default).
//
// The specified type is used for the entire list unless a different type
// attribute is used on an enclosed <li> element.
func (p *OlElement) Type(value string) *OlElement {
	p.addAttribute("type", value)
	return p
}

// Ol Represents an ordered list of items — typically rendered as a numbered list.
//
// For more details, see the [documentation]
//
// [documentation]: https://developer.mozilla.org/en-US/docs/Web/HTML/Element/ol
func Ol(content ...any) *OlElement {
	ne := newElement("ol", true, content...)

	var a = new(attrReversed[OlElement])
	var b = new(attrStart[OlElement])
	var ga = new(attrGlobal[OlElement])
	var ea = new(attrExternalAttributes[OlElement])
	var on = new(attrOn[OlElement])
	var ac = new(addContentFunc[OlElement])
	var el = &OlElement{ne, a, b, ga, ea, on, ac}
	a.set(ne, el)
	b.set(ne, el)
	ga.set(ne, el)
	ea.set(ne, el)
	on.set(ne, el)
	ac.set(ne, el)

	return el
}

// #region <p>

// PElement represents the <p> element.
type PElement struct {
	*element
	*attrGlobal[PElement]
	*attrExternalAttributes[PElement]
	*attrOn[PElement]
	*addContentFunc[PElement]
}

// P represents a paragraph. Paragraphs are usually represented in visual media
// as blocks of text separated from adjacent blocks by blank lines and/or
// first-line indentation, but HTML paragraphs can be any structural grouping
// of related content, such as images or form fields.
//
// For more details, see the [documentation]
//
// [documentation]: https://developer.mozilla.org/en-US/docs/Web/HTML/Element/p
func P(content ...any) *PElement {
	ne := newElement("p", true, content...)

	var ga = new(attrGlobal[PElement])
	var ea = new(attrExternalAttributes[PElement])
	var on = new(attrOn[PElement])
	var ac = new(addContentFunc[PElement])
	var el = &PElement{ne, ga, ea, on, ac}
	ga.set(ne, el)
	ea.set(ne, el)
	on.set(ne, el)
	ac.set(ne, el)

	return el
}

// #region <pre>

// PreElement represents the <pre> element.
type PreElement struct {
	*element
	*attrGlobal[PreElement]
	*attrExternalAttributes[PreElement]
	*attrOn[PreElement]
	*addContentFunc[PreElement]
}

// Pre represents preformatted text which is to be presented exactly as written
// in the HTML file. The text is typically rendered using a non-proportional, or
// monospaced, font. Whitespace inside this element is displayed as written.
//
// Example:
//
//	<figure>
//	  <pre role="img" aria-label="ASCII COW">
//	        ___________________________
//	    &lt; I'm an expert in my field. &gt;
//	        ---------------------------
//	            \   ^__^
//	             \  (oo)\_______
//	                (__)\       )\/\
//	                    ||----w |
//	                    ||     ||
//	  </pre>
//	  <figcaption id="cow-caption">
//	    A cow saying, "I'm an expert in my field." The cow is illustrated using
//	    preformatted text characters.
//	  </figcaption>
//	</figure>
//
// For more details, see the [documentation]
//
// [documentation]: https://developer.mozilla.org/en-US/docs/Web/HTML/Element/pre
func Pre(content ...any) *PreElement {
	ne := newElement("pre", true, content...)

	var ga = new(attrGlobal[PreElement])
	var ea = new(attrExternalAttributes[PreElement])
	var on = new(attrOn[PreElement])
	var ac = new(addContentFunc[PreElement])
	var el = &PreElement{ne, ga, ea, on, ac}
	ga.set(ne, el)
	ea.set(ne, el)
	on.set(ne, el)
	ac.set(ne, el)

	return el
}

// #region <ul>

// UlElement represents the <ul> element.
type UlElement struct {
	*element
	*attrGlobal[UlElement]
	*attrExternalAttributes[UlElement]
	*attrOn[UlElement]
	*addContentFunc[UlElement]
}

// Type sets the bullet style for the list.:
//   - disc.
//   - square.
//   - circle.
//
// The specified type is used for the entire list unless a different type
// attribute is used on an enclosed <li> element.
func (p *UlElement) Type(value string) *UlElement {
	p.addAttribute("type", value)
	return p
}

// Ul Represents an unordered list of items.
//
// For more details, see the [documentation]
//
// [documentation]: https://developer.mozilla.org/en-US/docs/Web/HTML/Element/ul
func Ul(content ...any) *UlElement {
	ne := newElement("ul", true, content...)

	var ga = new(attrGlobal[UlElement])
	var ea = new(attrExternalAttributes[UlElement])
	var on = new(attrOn[UlElement])
	var ac = new(addContentFunc[UlElement])
	var el = &UlElement{ne, ga, ea, on, ac}
	ga.set(ne, el)
	ea.set(ne, el)
	on.set(ne, el)
	ac.set(ne, el)

	return el
}

// #region INLINE TEXT SEMANTICS
// Use the HTML inline text semantic to define the meaning, structure, or style
// of a word, line, or any arbitrary piece of text.
//

// #region <a>

// AElement represents the <a> element.
type AElement struct {
	*element
	*attrDownload[AElement]
	*attrHref[AElement]
	*attrHrefLang[AElement]
	*attrPing[AElement]
	*attrRel[AElement]
	*attrTarget[AElement]
	*attrType[AElement]

	*attrGlobal[AElement]
	*attrExternalAttributes[AElement]
	*attrOn[AElement]
	*addContentFunc[AElement]
}

// A together with its href attribute, creates a hyperlink to web pages, files,
// email addresses, locations within the current page, or anything else a URL
// can address.
//
// For more details, see the [documentation]
//
// [documentation]: https://developer.mozilla.org/en-US/docs/Web/HTML/Element/a
func A(content ...any) *AElement {
	ne := newElement("a", true, content...)

	var b = new(attrDownload[AElement])
	var c = new(attrHref[AElement])
	var d = new(attrHrefLang[AElement])
	var e = new(attrPing[AElement])
	var f = new(attrRel[AElement])
	var g = new(attrTarget[AElement])
	var h = new(attrType[AElement])
	var ga = new(attrGlobal[AElement])
	var ea = new(attrExternalAttributes[AElement])
	var on = new(attrOn[AElement])
	var ac = new(addContentFunc[AElement])
	var el = &AElement{ne, b, c, d, e, f, g, h, ga, ea, on, ac}
	b.set(ne, el)
	c.set(ne, el)
	d.set(ne, el)
	e.set(ne, el)
	f.set(ne, el)
	g.set(ne, el)
	h.set(ne, el)
	ga.set(ne, el)
	ea.set(ne, el)
	on.set(ne, el)
	ac.set(ne, el)

	return el
}

// #region <abbr>

// AbbrElement represents the <abbr> element.
type AbbrElement struct {
	*element
	*attrGlobal[AbbrElement]
	*attrExternalAttributes[AbbrElement]
	*attrOn[AbbrElement]
	*addContentFunc[AbbrElement]
}

// Abbr represents an abbreviation or acronym.
//
// When including an abbreviation or acronym, provide a full expansion of the
// term in plain text on first use, along with the <abbr> to mark up the
// abbreviation. This informs the user what the abbreviation or acronym means.
//
// The optional title attribute can provide an expansion for the abbreviation
// or acronym when a full expansion is not present. This provides a hint to
// user agents on how to announce/display the content while informing all
// users what the abbreviation means. If present, title must contain this
// full description and nothing else.
//
// Example:
//
//	<p>Ashok's joke made me <abbr title="Laugh Out Loud">LOL</abbr> big time.</p>
//
// For more details, see the [documentation]
//
// [documentation]: https://developer.mozilla.org/en-US/docs/Web/HTML/Element/abbr
func Abbr(content ...any) *AbbrElement {
	ne := newElement("abbr", true, content...)

	var ga = new(attrGlobal[AbbrElement])
	var ea = new(attrExternalAttributes[AbbrElement])
	var on = new(attrOn[AbbrElement])
	var ac = new(addContentFunc[AbbrElement])
	var el = &AbbrElement{ne, ga, ea, on, ac}
	ga.set(ne, el)
	ea.set(ne, el)
	on.set(ne, el)
	ac.set(ne, el)

	return el
}

// #region <b>

// BElement represents the <b> element.
type BElement struct {
	*element
	*attrGlobal[BElement]
	*attrExternalAttributes[BElement]
	*attrOn[BElement]
	*addContentFunc[BElement]
}

// B is used to draw the reader's attention to the element's contents, which
// are not otherwise granted special importance. This was formerly known as
// the Boldface element, and most browsers still draw the text in boldface.
//
// However, you should not use <b> for styling text or granting importance.
// If you wish to create boldface text, you should use the CSS font-weight
// property. If you wish to indicate an element is of special importance, you
// should use the strong element.
//
// For more details, see the [documentation]
//
// [documentation]: https://developer.mozilla.org/en-US/docs/Web/HTML/Element/b
func B(content ...any) *BElement {
	ne := newElement("b", true, content...)

	var ga = new(attrGlobal[BElement])
	var ea = new(attrExternalAttributes[BElement])
	var on = new(attrOn[BElement])
	var ac = new(addContentFunc[BElement])
	var el = &BElement{ne, ga, ea, on, ac}
	ga.set(ne, el)
	ea.set(ne, el)
	on.set(ne, el)
	ac.set(ne, el)

	return el
}

// #region <bdi>

// BdiElement represents the <bdi> element.
type BdiElement struct {
	*element
	*attrGlobal[BdiElement]
	*attrExternalAttributes[BdiElement]
	*attrOn[BdiElement]
	*addContentFunc[BdiElement]
}

// Bdi tells the browser's bidirectional algorithm to treat the text it contains
// in isolation from its surrounding text. It's particularly useful when a
// website dynamically inserts some text and doesn't know the directionality
// of the text being inserted.
//
// For more details, see the [documentation]
//
// [documentation]: https://developer.mozilla.org/en-US/docs/Web/HTML/Element/bdi
func Bdi(content ...any) *BdiElement {
	ne := newElement("bdi", true, content...)

	var ga = new(attrGlobal[BdiElement])
	var ea = new(attrExternalAttributes[BdiElement])
	var on = new(attrOn[BdiElement])
	var ac = new(addContentFunc[BdiElement])
	var el = &BdiElement{ne, ga, ea, on, ac}
	ga.set(ne, el)
	ea.set(ne, el)
	on.set(ne, el)
	ac.set(ne, el)

	return el
}

// #region <bdo>

// BdoElement represents the <bdo> element.
type BdoElement struct {
	*element
	*attrGlobal[BdoElement]
	*attrExternalAttributes[BdoElement]
	*attrOn[BdoElement]
	*addContentFunc[BdoElement]
}

// Bdo overrides the current directionality of text, so that the text within is
// rendered in a different direction.
//
// For more details, see the [documentation]
//
// [documentation]: https://developer.mozilla.org/en-US/docs/Web/HTML/Element/bdo
func Bdo(content ...any) *BdoElement {
	ne := newElement("bdo", true, content...)

	var ga = new(attrGlobal[BdoElement])
	var ea = new(attrExternalAttributes[BdoElement])
	var on = new(attrOn[BdoElement])
	var ac = new(addContentFunc[BdoElement])
	var el = &BdoElement{ne, ga, ea, on, ac}
	ga.set(ne, el)
	ea.set(ne, el)
	on.set(ne, el)
	ac.set(ne, el)

	return el
}

// #region <br>

// BrElement represents the <br> element.
type BrElement struct {
	*element
	*attrGlobal[BrElement]
	*attrExternalAttributes[BrElement]
	*attrOn[BrElement]
}

// Br produces a line break in text (carriage-return). It is useful for writing
// a poem or an address, where the division of lines is significant.
//
// For more details, see the [documentation]
//
// [documentation]: https://developer.mozilla.org/en-US/docs/Web/HTML/Element/br
func Br() *BrElement {
	ne := newElement("br", false)
	var ga = new(attrGlobal[BrElement])
	var ea = new(attrExternalAttributes[BrElement])
	var on = new(attrOn[BrElement])
	var el = &BrElement{ne, ga, ea, on}
	ga.set(ne, el)
	ea.set(ne, el)
	on.set(ne, el)

	return el
}

// #region <cite>

// CiteElement represents the <cite> element.
type CiteElement struct {
	*element
	*attrGlobal[CiteElement]
	*attrExternalAttributes[CiteElement]
	*attrOn[CiteElement]
	*addContentFunc[CiteElement]
}

// Cite used to mark up the title of a cited creative work. The reference may
// be in an abbreviated form according to context-appropriate conventions
// related to citation metadata.
//
// For more details, see the [documentation]
//
// [documentation]: https://developer.mozilla.org/en-US/docs/Web/HTML/Element/cite
func Cite(content ...any) *CiteElement {
	ne := newElement("cite", true, content...)

	var ga = new(attrGlobal[CiteElement])
	var ea = new(attrExternalAttributes[CiteElement])
	var on = new(attrOn[CiteElement])
	var ac = new(addContentFunc[CiteElement])
	var el = &CiteElement{ne, ga, ea, on, ac}
	ga.set(ne, el)
	ea.set(ne, el)
	on.set(ne, el)
	ac.set(ne, el)

	return el
}

// #region <code>

// CodeElement represents the <code> element.
type CodeElement struct {
	*element
	*attrGlobal[CodeElement]
	*attrExternalAttributes[CodeElement]
	*attrOn[CodeElement]
	*addContentFunc[CodeElement]
}

// Code displays its contents styled in a fashion intended to indicate that the
// text is a short fragment of computer code. By default, the content text is
// displayed using the user agent's default monospace font.
//
// For more details, see the [documentation]
//
// [documentation]: https://developer.mozilla.org/en-US/docs/Web/HTML/Element/code
func Code(content ...any) *CodeElement {
	ne := newElement("code", true, content...)

	var ga = new(attrGlobal[CodeElement])
	var ea = new(attrExternalAttributes[CodeElement])
	var on = new(attrOn[CodeElement])
	var ac = new(addContentFunc[CodeElement])
	var el = &CodeElement{ne, ga, ea, on, ac}
	ga.set(ne, el)
	ea.set(ne, el)
	on.set(ne, el)
	ac.set(ne, el)

	return el
}

// #region <data>

// DataElement represents the <data> element.
type DataElement struct {
	*element
	*attrValue[DataElement]

	*attrGlobal[DataElement]
	*attrExternalAttributes[DataElement]
	*attrOn[DataElement]
	*addContentFunc[DataElement]
}

// Data links a given piece of content with a machine-readable translation.
// If the content is time- or date-related, the <time> element must be used.
//
// For more details, see the [documentation]
//
// [documentation]: https://developer.mozilla.org/en-US/docs/Web/HTML/Element/data
func Data(content ...any) *DataElement {
	ne := newElement("data", true, content...)

	var a = new(attrValue[DataElement])
	var ga = new(attrGlobal[DataElement])
	var ea = new(attrExternalAttributes[DataElement])
	var on = new(attrOn[DataElement])
	var ac = new(addContentFunc[DataElement])
	var el = &DataElement{ne, a, ga, ea, on, ac}
	a.set(ne, el)
	ga.set(ne, el)
	ea.set(ne, el)
	on.set(ne, el)
	ac.set(ne, el)

	return el
}

// #region <dfn>

// DfnElement represents the <dfn> element.
type DfnElement struct {
	*element
	*attrGlobal[DfnElement]
	*attrExternalAttributes[DfnElement]
	*attrOn[DfnElement]
	*addContentFunc[DfnElement]
}

// Dfn is used to indicate the term being defined within the context of a
// definition phrase or sentence. The ancestor <p> element, the <dt>/<dd>
// pairing, or the nearest section ancestor of the <dfn> element, is considered
// to be the definition of the term.
//
// Example:
//
//	<p>
//	  The <strong>HTML Definition element (<dfn>&lt;dfn&gt;</dfn>)</strong> is used
//	  to indicate the term being defined within the context of a definition phrase
//	  or sentence.
//	</p>
//
// For more details, see the [documentation]
//
// [documentation]: https://developer.mozilla.org/en-US/docs/Web/HTML/Element/dfn
func Dfn(content ...any) *DfnElement {
	ne := newElement("dfn", true, content...)

	var ga = new(attrGlobal[DfnElement])
	var ea = new(attrExternalAttributes[DfnElement])
	var on = new(attrOn[DfnElement])
	var ac = new(addContentFunc[DfnElement])
	var el = &DfnElement{ne, ga, ea, on, ac}
	ga.set(ne, el)
	ea.set(ne, el)
	on.set(ne, el)
	ac.set(ne, el)

	return el
}

// #region <em>

// EmElement represents the <em> element.
type EmElement struct {
	*element
	*attrGlobal[EmElement]
	*attrExternalAttributes[EmElement]
	*attrOn[EmElement]
	*addContentFunc[EmElement]
}

// Em marks text that has stress emphasis. The <em> element can be nested, with
// each nesting level indicating a greater degree of emphasis.
//
// Example:
//
//	<p>Get out of bed <em>now</em>!</p>
//	<p>We <em>had</em> to do something about it.</p>
//	<p>This is <em>not</em> a drill!</p>
//
// For more details, see the [documentation]
//
// [documentation]: https://developer.mozilla.org/en-US/docs/Web/HTML/Element/em
func Em(content ...any) *EmElement {
	ne := newElement("em", true, content...)

	var ga = new(attrGlobal[EmElement])
	var ea = new(attrExternalAttributes[EmElement])
	var on = new(attrOn[EmElement])
	var ac = new(addContentFunc[EmElement])
	var el = &EmElement{ne, ga, ea, on, ac}
	ga.set(ne, el)
	ea.set(ne, el)
	on.set(ne, el)
	ac.set(ne, el)

	return el
}

// #region <i>

// IElement represents the <i> element.
type IElement struct {
	*element
	*attrGlobal[IElement]
	*attrExternalAttributes[IElement]
	*attrOn[IElement]
	*addContentFunc[IElement]
}

// I represents a range of text that is set off from the normal text for some
// reason, such as idiomatic text, technical terms, and taxonomical designations,
// among others.
//
// Historically, these have been presented using italicized type, which is the
// original source of the <i> naming of this element.
//
// For more details, see the [documentation]
//
// [documentation]: https://developer.mozilla.org/en-US/docs/Web/HTML/Element/i
func I(content ...any) *IElement {
	ne := newElement("i", true, content...)

	var ga = new(attrGlobal[IElement])
	var ea = new(attrExternalAttributes[IElement])
	var on = new(attrOn[IElement])
	var ac = new(addContentFunc[IElement])
	var el = &IElement{ne, ga, ea, on, ac}
	ga.set(ne, el)
	ea.set(ne, el)
	on.set(ne, el)
	ac.set(ne, el)

	return el
}

// #region <kbd>

// KbdElement represents the <kbd> element.
type KbdElement struct {
	*element
	*attrGlobal[KbdElement]
	*attrExternalAttributes[KbdElement]
	*attrOn[KbdElement]
	*addContentFunc[KbdElement]
}

// Kbd represents a span of inline text denoting textual user input from a
// keyboard, voice input, or any other text entry device. By convention, the
// user agent defaults to rendering the contents of a <kbd> element using its
// default monospace font, although this is not mandated by the HTML standard.
//
// For more details, see the [documentation]
//
// [documentation]: https://developer.mozilla.org/en-US/docs/Web/HTML/Element/kbd
func Kbd(content ...any) *KbdElement {
	ne := newElement("kbd", true, content...)

	var ga = new(attrGlobal[KbdElement])
	var ea = new(attrExternalAttributes[KbdElement])
	var on = new(attrOn[KbdElement])
	var ac = new(addContentFunc[KbdElement])
	var el = &KbdElement{ne, ga, ea, on, ac}
	ga.set(ne, el)
	ea.set(ne, el)
	on.set(ne, el)
	ac.set(ne, el)

	return el
}

// #region <mark>

// MarkElement represents the <mark> element.
type MarkElement struct {
	*element
	*attrGlobal[MarkElement]
	*attrExternalAttributes[MarkElement]
	*attrOn[MarkElement]
	*addContentFunc[MarkElement]
}

// Mark represents text which is marked or highlighted for reference or notation
// purposes due to the marked passage's relevance in the enclosing context.
//
// For more details, see the [documentation]
//
// [documentation]: https://developer.mozilla.org/en-US/docs/Web/HTML/Element/mark
func Mark(content ...any) *MarkElement {
	ne := newElement("mark", true, content...)

	var ga = new(attrGlobal[MarkElement])
	var ea = new(attrExternalAttributes[MarkElement])
	var on = new(attrOn[MarkElement])
	var ac = new(addContentFunc[MarkElement])
	var el = &MarkElement{ne, ga, ea, on, ac}
	ga.set(ne, el)
	ea.set(ne, el)
	on.set(ne, el)
	ac.set(ne, el)

	return el
}

// #region <q>

// QElement represents the <q> element.
type QElement struct {
	*element
	*attrCite[QElement]

	*attrGlobal[QElement]
	*attrExternalAttributes[QElement]
	*attrOn[QElement]
	*addContentFunc[QElement]
}

// Q indicates that the enclosed text is a short inline quotation. Most modern
// browsers implement this by surrounding the text in quotation marks.
//
// This element is intended for short quotations that don't require paragraph
// breaks; for long quotations use the <blockquote> element.
//
// For more details, see the [documentation]
//
// [documentation]: https://developer.mozilla.org/en-US/docs/Web/HTML/Element/q
func Q(content ...any) *QElement {
	ne := newElement("q", true, content...)

	var a = new(attrCite[QElement])
	var ga = new(attrGlobal[QElement])
	var ea = new(attrExternalAttributes[QElement])
	var on = new(attrOn[QElement])
	var ac = new(addContentFunc[QElement])
	var el = &QElement{ne, a, ga, ea, on, ac}
	a.set(ne, el)
	ga.set(ne, el)
	ea.set(ne, el)
	on.set(ne, el)
	ac.set(ne, el)

	return el
}

// #region <rp>

// RpElement represents the <rp> element.
type RpElement struct {
	*element
	*attrGlobal[RpElement]
	*attrExternalAttributes[RpElement]
	*attrOn[RpElement]
	*addContentFunc[RpElement]
}

// Rp is used to provide fall-back parentheses for browsers that do not support
// the display of ruby annotations using the <ruby> element. One <rp> element
// should enclose each of the opening and closing parentheses that wrap the
// <rt> element that contains the annotation's text.
//
// For more details, see the [documentation]
//
// [documentation]: https://developer.mozilla.org/en-US/docs/Web/HTML/Element/rp
func Rp(content ...any) *RpElement {
	ne := newElement("rp", true, content...)

	var ga = new(attrGlobal[RpElement])
	var ea = new(attrExternalAttributes[RpElement])
	var on = new(attrOn[RpElement])
	var ac = new(addContentFunc[RpElement])
	var el = &RpElement{ne, ga, ea, on, ac}
	ga.set(ne, el)
	ea.set(ne, el)
	on.set(ne, el)
	ac.set(ne, el)

	return el
}

// #region <rt>

// RtElement represents the <rt> element.
type RtElement struct {
	*element
	*attrGlobal[RtElement]
	*attrExternalAttributes[RtElement]
	*attrOn[RtElement]
	*addContentFunc[RtElement]
}

// Rt specifies the ruby text component of a ruby annotation, which is used to
// provide pronunciation, translation, or transliteration information for East
// Asian typography. The <rt> element must always be contained within a <ruby>
// element.
//
// For more details, see the [documentation]
//
// [documentation]: https://developer.mozilla.org/en-US/docs/Web/HTML/Element/rt
func Rt(content ...any) *RtElement {
	ne := newElement("rt", true, content...)

	var ga = new(attrGlobal[RtElement])
	var ea = new(attrExternalAttributes[RtElement])
	var on = new(attrOn[RtElement])
	var ac = new(addContentFunc[RtElement])
	var el = &RtElement{ne, ga, ea, on, ac}
	ga.set(ne, el)
	ea.set(ne, el)
	on.set(ne, el)
	ac.set(ne, el)

	return el
}

// #region <s>

// SElement represents the <s> element.
type SElement struct {
	*element
	*attrGlobal[SElement]
	*attrExternalAttributes[SElement]
	*attrOn[SElement]
	*addContentFunc[SElement]
}

// s renders text with a strikethrough, or a line through it. Use the <s>
// element to represent things that are no longer relevant or no longer accurate.
//
// However, <s> is not appropriate when indicating document edits; for that,
// use the <del> and <ins> elements, as appropriate.
//
// For more details, see the [documentation]
//
// [documentation]: https://developer.mozilla.org/en-US/docs/Web/HTML/Element/s
func S(content ...any) *SElement {
	ne := newElement("s", true, content...)

	var ga = new(attrGlobal[SElement])
	var ea = new(attrExternalAttributes[SElement])
	var on = new(attrOn[SElement])
	var ac = new(addContentFunc[SElement])
	var el = &SElement{ne, ga, ea, on, ac}
	ga.set(ne, el)
	ea.set(ne, el)
	on.set(ne, el)
	ac.set(ne, el)

	return el
}

// #region <samp>

// SampElement represents the <samp> element.
type SampElement struct {
	*element
	*attrGlobal[SampElement]
	*attrExternalAttributes[SampElement]
	*attrOn[SampElement]
	*addContentFunc[SampElement]
}

// Samp is used to enclose inline text which represents sample (or quoted)
// output from a computer program. Its contents are typically rendered using
// the browser's default monospaced font (such as Courier or Lucida Console)
//
// For more details, see the [documentation]
//
// [documentation]: https://developer.mozilla.org/en-US/docs/Web/HTML/Element/samp
func Samp(content ...any) *SampElement {
	ne := newElement("samp", true, content...)

	var ga = new(attrGlobal[SampElement])
	var ea = new(attrExternalAttributes[SampElement])
	var on = new(attrOn[SampElement])
	var ac = new(addContentFunc[SampElement])
	var el = &SampElement{ne, ga, ea, on, ac}
	ga.set(ne, el)
	ea.set(ne, el)
	on.set(ne, el)
	ac.set(ne, el)

	return el
}

// #region <small>

// SmallElement represents the <small> element.
type SmallElement struct {
	*element
	*attrGlobal[SmallElement]
	*attrExternalAttributes[SmallElement]
	*attrOn[SmallElement]
	*addContentFunc[SmallElement]
}

// Small represents side-comments and small print, like copyright and legal text,
// independent of its styled presentation. By default, it renders text within
// it one font size smaller, such as from small to x-small.
//
// For more details, see the [documentation]
//
// [documentation]: https://developer.mozilla.org/en-US/docs/Web/HTML/Element/small
func Small(content ...any) *SmallElement {
	ne := newElement("small", true, content...)

	var ga = new(attrGlobal[SmallElement])
	var ea = new(attrExternalAttributes[SmallElement])
	var on = new(attrOn[SmallElement])
	var ac = new(addContentFunc[SmallElement])
	var el = &SmallElement{ne, ga, ea, on, ac}
	ga.set(ne, el)
	ea.set(ne, el)
	on.set(ne, el)
	ac.set(ne, el)

	return el
}

// #region <span>

// SpanElement represents the <span> element.
type SpanElement struct {
	*element
	*attrGlobal[SpanElement]
	*attrExternalAttributes[SpanElement]
	*attrOn[SpanElement]
	*addContentFunc[SpanElement]
}

// Span is a generic inline container for phrasing content, which does not
// inherently represent anything. It can be used to group elements for
// styling purposes (using the class or id attributes), or because they share
// attribute values, such as lang. It should be used only when no other
// semantic element is appropriate. <span> is very much like a div element, but
// div is a block-level element whereas a <span> is an inline-level element.
//
// For more details, see the [documentation]
//
// [documentation]: https://developer.mozilla.org/en-US/docs/Web/HTML/Element/span
func Span(content ...any) *SpanElement {
	ne := newElement("span", true, content...)

	var ga = new(attrGlobal[SpanElement])
	var ea = new(attrExternalAttributes[SpanElement])
	var on = new(attrOn[SpanElement])
	var ac = new(addContentFunc[SpanElement])
	var el = &SpanElement{ne, ga, ea, on, ac}
	ga.set(ne, el)
	ea.set(ne, el)
	on.set(ne, el)
	ac.set(ne, el)

	return el
}

// #region <strong>

// StrongElement represents the <strong> element.
type StrongElement struct {
	*element
	*attrGlobal[StrongElement]
	*attrExternalAttributes[StrongElement]
	*attrOn[StrongElement]
	*addContentFunc[StrongElement]
}

// Strong indicates that its contents have strong importance, seriousness, or
// urgency. Browsers typically render the contents in bold type.
//
// For more details, see the [documentation]
//
// [documentation]: https://developer.mozilla.org/en-US/docs/Web/HTML/Element/strong
func Strong(content ...any) *StrongElement {
	ne := newElement("strong", true, content...)

	var ga = new(attrGlobal[StrongElement])
	var ea = new(attrExternalAttributes[StrongElement])
	var on = new(attrOn[StrongElement])
	var ac = new(addContentFunc[StrongElement])
	var el = &StrongElement{ne, ga, ea, on, ac}
	ga.set(ne, el)
	ea.set(ne, el)
	on.set(ne, el)
	ac.set(ne, el)

	return el
}

// #region <sub>

// SubElement represents the <sub> element.
type SubElement struct {
	*element
	*attrGlobal[SubElement]
	*attrExternalAttributes[SubElement]
	*attrOn[SubElement]
	*addContentFunc[SubElement]
}

// Sub specifies inline text which should be displayed as subscript for solely
// typographical reasons. Subscripts are typically rendered with a lowered
// baseline using smaller text.
//
// For more details, see the [documentation]
//
// [documentation]: https://developer.mozilla.org/en-US/docs/Web/HTML/Element/sub
func Sub(content ...any) *SubElement {
	ne := newElement("sub", true, content...)

	var ga = new(attrGlobal[SubElement])
	var ea = new(attrExternalAttributes[SubElement])
	var on = new(attrOn[SubElement])
	var ac = new(addContentFunc[SubElement])
	var el = &SubElement{ne, ga, ea, on, ac}
	ga.set(ne, el)
	ea.set(ne, el)
	on.set(ne, el)
	ac.set(ne, el)

	return el
}

// #region <sup>

// SupElement represents the <sup> element.
type SupElement struct {
	*element
	*attrGlobal[SupElement]
	*attrExternalAttributes[SupElement]
	*attrOn[SupElement]
	*addContentFunc[SupElement]
}

// Sup specifies inline text which is to be displayed as superscript for solely
// typographical reasons. Superscripts are usually rendered with a raised
// baseline using smaller text.
//
// For more details, see the [documentation]
//
// [documentation]: https://developer.mozilla.org/en-US/docs/Web/HTML/Element/sup
func Sup(content ...any) *SupElement {
	ne := newElement("sup", true, content...)

	var ga = new(attrGlobal[SupElement])
	var ea = new(attrExternalAttributes[SupElement])
	var on = new(attrOn[SupElement])
	var ac = new(addContentFunc[SupElement])
	var el = &SupElement{ne, ga, ea, on, ac}
	ga.set(ne, el)
	ea.set(ne, el)
	on.set(ne, el)
	ac.set(ne, el)

	return el
}

// #region <time>

// TimeElement represents the <time> element.
type TimeElement struct {
	*element
	*attrDateTime[TimeElement]

	*attrGlobal[TimeElement]
	*attrExternalAttributes[TimeElement]
	*attrOn[TimeElement]
	*addContentFunc[TimeElement]
}

// Time represents a specific period in time. It may include the datetime
// attribute to translate dates into machine-readable format, allowing for
// better search engine results or custom features such as reminders.
//
// For more details, see the [documentation]
//
// [documentation]: https://developer.mozilla.org/en-US/docs/Web/HTML/Element/time
func Time(content ...any) *TimeElement {
	ne := newElement("time", true, content...)

	var a = new(attrDateTime[TimeElement])
	var ga = new(attrGlobal[TimeElement])
	var ea = new(attrExternalAttributes[TimeElement])
	var on = new(attrOn[TimeElement])
	var ac = new(addContentFunc[TimeElement])
	var el = &TimeElement{ne, a, ga, ea, on, ac}
	a.set(ne, el)
	ga.set(ne, el)
	ea.set(ne, el)
	on.set(ne, el)
	ac.set(ne, el)

	return el
}

// #region <u>

// UElement represents the <u> element.
type UElement struct {
	*element
	*attrGlobal[UElement]
	*attrExternalAttributes[UElement]
	*attrOn[UElement]
	*addContentFunc[UElement]
}

// U represents a span of inline text which should be rendered in a way that
// indicates that it has a non-textual annotation. This is rendered by default
// as a simple solid underline but may be altered using CSS.
//
// For more details, see the [documentation]
//
// [documentation]: https://developer.mozilla.org/en-US/docs/Web/HTML/Element/u
func U(content ...any) *UElement {
	ne := newElement("u", true, content...)

	var ga = new(attrGlobal[UElement])
	var ea = new(attrExternalAttributes[UElement])
	var on = new(attrOn[UElement])
	var ac = new(addContentFunc[UElement])
	var el = &UElement{ne, ga, ea, on, ac}
	ga.set(ne, el)
	ea.set(ne, el)
	on.set(ne, el)
	ac.set(ne, el)

	return el
}

// #region <var>

// VarElement represents the <var> element.
type VarElement struct {
	*element
	*attrGlobal[VarElement]
	*attrExternalAttributes[VarElement]
	*attrOn[VarElement]
	*addContentFunc[VarElement]
}

// Var represents the name of a variable in a mathematical expression or a
// programming context. It's typically presented using an italicized version
// of the current typeface, although that behavior is browser-dependent.
//
// For more details, see the [documentation]
//
// [documentation]: https://developer.mozilla.org/en-US/docs/Web/HTML/Element/var
func Var(content ...any) *VarElement {
	ne := newElement("var", true, content...)

	var ga = new(attrGlobal[VarElement])
	var ea = new(attrExternalAttributes[VarElement])
	var on = new(attrOn[VarElement])
	var ac = new(addContentFunc[VarElement])
	var el = &VarElement{ne, ga, ea, on, ac}
	ga.set(ne, el)
	ea.set(ne, el)
	on.set(ne, el)
	ac.set(ne, el)

	return el
}

// #region <wbr>

// WbrElement represents the <wbr> element.
type WbrElement struct {
	*element
	*attrGlobal[WbrElement]
	*attrExternalAttributes[WbrElement]
	*attrOn[WbrElement]
}

// Wbr represents a word break opportunity—a position within text where the
// browser may optionally break a line, though its line-breaking rules would
// not otherwise create a break at that location.
//
// For more details, see the [documentation]
//
// [documentation]: https://developer.mozilla.org/en-US/docs/Web/HTML/Element/wbr
func Wbr() *WbrElement {
	ne := newElement("wbr", false)

	var ga = new(attrGlobal[WbrElement])
	var ea = new(attrExternalAttributes[WbrElement])
	var on = new(attrOn[WbrElement])
	var el = &WbrElement{ne, ga, ea, on}
	ga.set(ne, el)
	ea.set(ne, el)
	on.set(ne, el)

	return el
}

// #region IMAGE AND MULTIMEDIA
// HTML supports various multimedia resources such as images, audio, and video.
//
//

// #region <area>

// AreaElement represents the <area> element.
type AreaElement struct {
	*element
	*attrAlt[AreaElement]
	*attrCoords[AreaElement]
	*attrDownload[AreaElement]
	*attrHref[AreaElement]
	*attrPing[AreaElement]
	*attrRel[AreaElement]
	*attrShape[AreaElement]
	*attrTarget[AreaElement]

	*attrGlobal[AreaElement]
	*attrExternalAttributes[AreaElement]
	*attrOn[AreaElement]
}

// Area defines an area inside an image map that has predefined clickable areas.
// An image map allows geometric areas on an image to be associated with
// hyperlink.
//
// For more details, see the [documentation]
//
// [documentation]: https://developer.mozilla.org/en-US/docs/Web/HTML/Element/area
func Area() *AreaElement {
	ne := newElement("area", false)

	var a = new(attrAlt[AreaElement])
	var b = new(attrCoords[AreaElement])
	var c = new(attrDownload[AreaElement])
	var d = new(attrHref[AreaElement])
	var e = new(attrPing[AreaElement])
	var f = new(attrRel[AreaElement])
	var g = new(attrShape[AreaElement])
	var h = new(attrTarget[AreaElement])
	var ga = new(attrGlobal[AreaElement])
	var ea = new(attrExternalAttributes[AreaElement])
	var on = new(attrOn[AreaElement])
	var el = &AreaElement{ne, a, b, c, d, e, f, g, h, ga, ea, on}
	a.set(ne, el)
	b.set(ne, el)
	c.set(ne, el)
	d.set(ne, el)
	e.set(ne, el)
	f.set(ne, el)
	g.set(ne, el)
	h.set(ne, el)
	ga.set(ne, el)
	ea.set(ne, el)
	on.set(ne, el)

	return el
}

// #region <img>

// ImgElement represents the <img> element.
type ImgElement struct {
	*element
	*attrAlt[ImgElement]
	*attrCrossOrigin[ImgElement]
	*attrDecoding[ImgElement]
	*attrIsmap[ImgElement]
	*attrLoading[ImgElement]
	*attrSrc[ImgElement]
	*attrSrcSet[ImgElement]

	*attrGlobal[ImgElement]
	*attrExternalAttributes[ImgElement]
	*attrOn[ImgElement]
}

// Img embeds an image into the document.
//
// For more details, see the [documentation]
//
// [documentation]: https://developer.mozilla.org/en-US/docs/Web/HTML/Element/img
func Img() *ImgElement {
	ne := newElement("img", false)

	var a = new(attrAlt[ImgElement])
	var b = new(attrCrossOrigin[ImgElement])
	var c = new(attrDecoding[ImgElement])
	var d = new(attrIsmap[ImgElement])
	var e = new(attrLoading[ImgElement])
	var f = new(attrSrc[ImgElement])
	var g = new(attrSrcSet[ImgElement])
	var ga = new(attrGlobal[ImgElement])
	var ea = new(attrExternalAttributes[ImgElement])
	var on = new(attrOn[ImgElement])
	var el = &ImgElement{ne, a, b, c, d, e, f, g, ga, ea, on}
	a.set(ne, el)
	b.set(ne, el)
	c.set(ne, el)
	d.set(ne, el)
	e.set(ne, el)
	f.set(ne, el)
	g.set(ne, el)
	ga.set(ne, el)
	ea.set(ne, el)
	on.set(ne, el)

	return el
}

// #region <map>

// MapElement represents the <map> element.
type MapElement struct {
	*element
	*attrName[MapElement]

	*attrGlobal[MapElement]
	*attrExternalAttributes[MapElement]
	*attrOn[MapElement]
	*addContentFunc[MapElement]
}

// Map is used with <area> elements to define an image map (a clickable link area).
//
// For more details, see the [documentation]
//
// [documentation]: https://developer.mozilla.org/en-US/docs/Web/HTML/Element/map
func Map(content ...any) *MapElement {
	ne := newElement("map", true, content...)

	var a = new(attrName[MapElement])
	var ga = new(attrGlobal[MapElement])
	var ea = new(attrExternalAttributes[MapElement])
	var on = new(attrOn[MapElement])
	var ac = new(addContentFunc[MapElement])
	var el = &MapElement{ne, a, ga, ea, on, ac}
	a.set(ne, el)
	ga.set(ne, el)
	ea.set(ne, el)
	on.set(ne, el)
	ac.set(ne, el)

	return el
}

// #region <track>

// TrackElement represents the <track> element.
type TrackElement struct {
	*element
	*attrDefault[TrackElement]
	*attrKind[TrackElement]
	*attrLabel[TrackElement]
	*attrSrc[TrackElement]
	*attrSrcLang[TrackElement]

	*attrGlobal[TrackElement]
	*attrExternalAttributes[TrackElement]
	*attrOn[TrackElement]
}

// Track is used as a child of the media elements, audio and video. It lets you
// specify timed text tracks (or time-based data), for example to automatically
// handle subtitles. The tracks are formatted in WebVTT format (.vtt files)—Web
// Video Text Tracks.
//
// For more details, see the [documentation]
//
// [documentation]: https://developer.mozilla.org/en-US/docs/Web/HTML/Element/track
func Track() *TrackElement {
	ne := newElement("track", false)

	var a = new(attrDefault[TrackElement])
	var b = new(attrKind[TrackElement])
	var c = new(attrLabel[TrackElement])
	var d = new(attrSrc[TrackElement])
	var e = new(attrSrcLang[TrackElement])
	var ga = new(attrGlobal[TrackElement])
	var ea = new(attrExternalAttributes[TrackElement])
	var on = new(attrOn[TrackElement])
	var el = &TrackElement{ne, a, b, c, d, e, ga, ea, on}
	a.set(ne, el)
	b.set(ne, el)
	c.set(ne, el)
	d.set(ne, el)
	e.set(ne, el)
	ga.set(ne, el)
	ea.set(ne, el)
	on.set(ne, el)
	return el
}

// #region <video>

// VideoElement represents the <video> element.
type VideoElement struct {
	*element
	*attrAutoPlay[VideoElement]
	*attrControls[VideoElement]
	*attrControlsList[VideoElement]
	*attrCrossOrigin[VideoElement]
	*attrDisablePictureInPicture[VideoElement]
	*attrDisableRemotePlayBack[VideoElement]
	*attrLoop[VideoElement]
	*attrMuted[VideoElement]
	*attrPlaysInLine[VideoElement]
	*attrPoster[VideoElement]
	*attrPreLoad[VideoElement]
	*attrSrc[VideoElement]

	*attrGlobal[VideoElement]
	*attrExternalAttributes[VideoElement]
	*attrOn[VideoElement]
	*addContentFunc[VideoElement]
}

// Video embeds a media player which supports video playback into the document.
// You can also use <video> for audio content, but the audio element may provide
// a more appropriate user experience.
//
// For more details, see the [documentation]
//
// [documentation]: https://developer.mozilla.org/en-US/docs/Web/HTML/Element/video
func Video(content ...any) *VideoElement {
	ne := newElement("video", true, content...)

	var a = new(attrAutoPlay[VideoElement])
	var b = new(attrControls[VideoElement])
	var c = new(attrControlsList[VideoElement])
	var d = new(attrCrossOrigin[VideoElement])
	var e = new(attrDisablePictureInPicture[VideoElement])
	var f = new(attrDisableRemotePlayBack[VideoElement])
	var g = new(attrLoop[VideoElement])
	var h = new(attrMuted[VideoElement])
	var i = new(attrPlaysInLine[VideoElement])
	var j = new(attrPoster[VideoElement])
	var k = new(attrPreLoad[VideoElement])
	var l = new(attrSrc[VideoElement])
	var ga = new(attrGlobal[VideoElement])
	var ea = new(attrExternalAttributes[VideoElement])
	var on = new(attrOn[VideoElement])
	var ac = new(addContentFunc[VideoElement])
	var el = &VideoElement{ne, a, b, c, d, e, f, g, h, i, j, k, l, ga, ea, on, ac}
	a.set(ne, el)
	b.set(ne, el)
	c.set(ne, el)
	d.set(ne, el)
	e.set(ne, el)
	f.set(ne, el)
	g.set(ne, el)
	h.set(ne, el)
	i.set(ne, el)
	j.set(ne, el)
	k.set(ne, el)
	l.set(ne, el)
	ga.set(ne, el)
	ea.set(ne, el)
	on.set(ne, el)
	ac.set(ne, el)

	return el
}

// #region EMBEDDED CONTENT
// In addition to regular multimedia content, HTML can include a variety of
// other content, even if it's not always easy to interact with.
//

// #region <embed>

// EmbedElement represents the <embed> element.
type EmbedElement struct {
	*element
	*attrType[EmbedElement]
	*attrSrc[EmbedElement]

	*attrGlobal[EmbedElement]
	*attrExternalAttributes[EmbedElement]
	*attrOn[EmbedElement]
}

// Embed is an external content at the specified point in the document. This
// content is provided by an external application or other source of interactive
// content such as a browser plug-in.
//
// For more details, see the [documentation]
//
// [documentation]: https://developer.mozilla.org/en-US/docs/Web/HTML/Element/embed
func Embed() *EmbedElement {
	ne := newElement("embed", false)

	var a = new(attrType[EmbedElement])
	var b = new(attrSrc[EmbedElement])
	var ga = new(attrGlobal[EmbedElement])
	var ea = new(attrExternalAttributes[EmbedElement])
	var on = new(attrOn[EmbedElement])
	var el = &EmbedElement{ne, a, b, ga, ea, on}
	a.set(ne, el)
	b.set(ne, el)
	ga.set(ne, el)
	ea.set(ne, el)
	on.set(ne, el)
	return el
}

// #region <iframe>

// IframeElement represents the <iframe> element.
type IframeElement struct {
	*element
	*attrAllow[IframeElement]
	*attrLoading[IframeElement]
	*attrName[IframeElement]
	*attrSandbox[IframeElement]
	*attrSrc[IframeElement]

	*attrGlobal[IframeElement]
	*attrExternalAttributes[IframeElement]
	*attrOn[IframeElement]
	*addContentFunc[IframeElement]
}

// Iframe represents a nested browsing context, embedding another HTML page
// into the current one.
//
// For more details, see the [documentation]
//
// [documentation]: https://developer.mozilla.org/en-US/docs/Web/HTML/Element/iframe
func Iframe(content ...any) *IframeElement {
	ne := newElement("iframe", true, content...)

	var a = new(attrAllow[IframeElement])
	var b = new(attrLoading[IframeElement])
	var c = new(attrName[IframeElement])
	var d = new(attrSandbox[IframeElement])
	var e = new(attrSrc[IframeElement])
	var ga = new(attrGlobal[IframeElement])
	var ea = new(attrExternalAttributes[IframeElement])
	var on = new(attrOn[IframeElement])
	var ac = new(addContentFunc[IframeElement])
	var el = &IframeElement{ne, a, b, c, d, e, ga, ea, on, ac}
	a.set(ne, el)
	b.set(ne, el)
	c.set(ne, el)
	d.set(ne, el)
	e.set(ne, el)
	ga.set(ne, el)
	ea.set(ne, el)
	on.set(ne, el)
	ac.set(ne, el)

	return el
}

// #region <object>

// ObjectElement represents the <object> element.
type ObjectElement struct {
	*element
	*attrData[ObjectElement]
	*attrForm[ObjectElement]
	*attrName[ObjectElement]
	*attrType[ObjectElement]

	*attrGlobal[ObjectElement]
	*attrExternalAttributes[ObjectElement]
	*attrOn[ObjectElement]
}

// Object represents an external resource, which can be treated as an image, a
// nested browsing context, or a resource to be handled by a plugin.
//
// For more details, see the [documentation]
//
// [documentation]: https://developer.mozilla.org/en-US/docs/Web/HTML/Element/object
func Object() *ObjectElement {
	ne := newElement("object", false)

	var a = new(attrData[ObjectElement])
	var b = new(attrForm[ObjectElement])
	var c = new(attrName[ObjectElement])
	var d = new(attrType[ObjectElement])
	var ga = new(attrGlobal[ObjectElement])
	var ea = new(attrExternalAttributes[ObjectElement])
	var on = new(attrOn[ObjectElement])
	var el = &ObjectElement{ne, a, b, c, d, ga, ea, on}
	a.set(ne, el)
	b.set(ne, el)
	c.set(ne, el)
	d.set(ne, el)
	ga.set(ne, el)
	ea.set(ne, el)
	on.set(ne, el)
	return el
}

// #region <picture>

// PictureElement represents the <picture> element.
type PictureElement struct {
	*element
	*attrGlobal[PictureElement]
	*attrExternalAttributes[PictureElement]
	*attrOn[PictureElement]
	*addContentFunc[PictureElement]
}

// Picture contains zero or more <source> elements and one <img> element to
// offer alternative versions of an image for different display/device scenarios.
//
// For more details, see the [documentation]
//
// [documentation]: https://developer.mozilla.org/en-US/docs/Web/HTML/Element/picture
func Picture(content ...any) *PictureElement {
	ne := newElement("picture", true, content...)

	var ga = new(attrGlobal[PictureElement])
	var ea = new(attrExternalAttributes[PictureElement])
	var on = new(attrOn[PictureElement])
	var ac = new(addContentFunc[PictureElement])
	var el = &PictureElement{ne, ga, ea, on, ac}
	ga.set(ne, el)
	ea.set(ne, el)
	on.set(ne, el)
	ac.set(ne, el)

	return el
}

// #region <source>

// SourceElement represents the <source> element.
type SourceElement struct {
	*element
	*attrType[SourceElement]
	*attrSrc[SourceElement]
	*attrSrcSet[SourceElement]
	*attrSizes[SourceElement]
	*attrMedia[SourceElement]

	*attrGlobal[SourceElement]
	*attrExternalAttributes[SourceElement]
	*attrOn[SourceElement]
}

// Surce specifies multiple media resources for the picture, the audio element,
// or the video element. It is a void element, meaning that it has no content
// and does not have a closing tag.
//
// content in multiple file formats in order to provide compatibility with a
// It is commonly used to offer the same media
// broad range of browsers given their differing support for image file formats
// and media file formats.
//
// For more details, see the [documentation]
//
// [documentation]: https://developer.mozilla.org/en-US/docs/Web/HTML/Element/source
func Source() *SourceElement {
	ne := newElement("source", false)

	var a = new(attrType[SourceElement])
	var b = new(attrSrc[SourceElement])
	var c = new(attrSrcSet[SourceElement])
	var d = new(attrSizes[SourceElement])
	var e = new(attrMedia[SourceElement])
	var ga = new(attrGlobal[SourceElement])
	var ea = new(attrExternalAttributes[SourceElement])
	var on = new(attrOn[SourceElement])
	var el = &SourceElement{ne, a, b, c, d, e, ga, ea, on}
	a.set(ne, el)
	b.set(ne, el)
	c.set(ne, el)
	d.set(ne, el)
	e.set(ne, el)
	ga.set(ne, el)
	ea.set(ne, el)
	on.set(ne, el)
	return el
}

// #region SVG AND MathML
// You can embed SVG and MathML content directly into HTML documents, using
// the <svg> and <math> elements.
//

// #region <svg>

// SvgElement represents the <svg> element.
type SvgElement struct {
	*element
	*attrGlobal[SvgElement]
	*attrExternalAttributes[SvgElement]
	*attrOn[SvgElement]
	*addContentFunc[SvgElement]
}

// PreserveAspectRatio represents how the svg fragment must be deformed if it
// is displayed with a different aspect ratio.
//
// Value type: (none| xMinYMin| xMidYMin| xMaxYMin| xMinYMid| xMidYMid|
// xMaxYMid| xMinYMax| xMidYMax| xMaxYMax) (meet|slice).
// Default value: xMidYMid meet; Animatable: yes
func (p *SvgElement) PreserveAspectRatio(value string) *SvgElement {
	p.addAttribute("preserveAspectRatio", value)
	return p
}

// Viewbox is the SVG viewport coordinates for the current SVG fragment.
//
// Value type: <list-of-numbers> ; Default value: none; Animatable: yes
func (p *SvgElement) Viewbox(value string) *SvgElement {
	p.addAttribute("viewbox", value)
	return p
}

// X is the displayed x coordinate of the svg container. No effect on outermost
// svg elements. Value type: <length>|<percentage> ; Default value: 0;
// Animatable: yes
func (p *SvgElement) X(value string) *SvgElement {
	p.addAttribute("x", value)
	return p
}

// Y is the displayed y coordinate of the svg container. No effect on outermost
// svg elements. Value type: <length>|<percentage> ; Default value: 0;
// Animatable: yes
func (p *SvgElement) Y(value string) *SvgElement {
	p.addAttribute("y", value)
	return p
}

// Svg is an container defining a new coordinate system and viewport. It is
// used as the outermost element of SVG documents, but it can also be used to
// embed an SVG fragment inside an SVG or HTML document.
//
// For more details, see the [documentation]
//
// [documentation]: https://developer.mozilla.org/en-US/docs/Web/SVG/Element/svg
func Svg(content ...any) *SvgElement {
	ne := newElement("svg", true, content...)

	var ga = new(attrGlobal[SvgElement])
	var ea = new(attrExternalAttributes[SvgElement])
	var on = new(attrOn[SvgElement])
	var ac = new(addContentFunc[SvgElement])
	var el = &SvgElement{ne, ga, ea, on, ac}
	ga.set(ne, el)
	ea.set(ne, el)
	on.set(ne, el)
	ac.set(ne, el)

	return el
}

// #region SCRIPTING
// To create dynamic content and Web applications, HTML supports the use of
// scripting languages, most prominently JavaScript. Certain elements support
// this capability.

// #region <canvas>

// CanvasElement represents the <canvas> element.
type CanvasElement struct {
	*element
	*attrGlobal[CanvasElement]
	*attrExternalAttributes[CanvasElement]
	*attrOn[CanvasElement]
	*addContentFunc[CanvasElement]
}

// Canvas is an container element to use with either the canvas scripting API
// or the WebGL API to draw graphics and animations.
//
// For more details, see the [documentation]
//
// [documentation]: https://developer.mozilla.org/en-US/docs/Web/HTML/Element/canvas
func Canvas(content ...any) *CanvasElement {
	ne := newElement("canvas", true, content...)

	var ga = new(attrGlobal[CanvasElement])
	var ea = new(attrExternalAttributes[CanvasElement])
	var on = new(attrOn[CanvasElement])
	var ac = new(addContentFunc[CanvasElement])
	var el = &CanvasElement{ne, ga, ea, on, ac}
	ga.set(ne, el)
	ea.set(ne, el)
	on.set(ne, el)
	ac.set(ne, el)

	return el
}

// #region <noscript>

// NoScriptElement represents the <noscript> element.
type NoScriptElement struct {
	*element
	*attrGlobal[NoScriptElement]
	*attrExternalAttributes[NoScriptElement]
	*attrOn[NoScriptElement]
	*addContentFunc[NoScriptElement]
}

// NoScript defines a section of HTML to be inserted if a script type on the
// page is unsupported or if scripting is currently turned off in the browser.
//
// For more details, see the [documentation]
//
// [documentation]: https://developer.mozilla.org/en-US/docs/Web/HTML/Element/noscript
func NoScript(content ...any) *NoScriptElement {
	ne := newElement("noscript", true, content...)

	var ga = new(attrGlobal[NoScriptElement])
	var ea = new(attrExternalAttributes[NoScriptElement])
	var on = new(attrOn[NoScriptElement])
	var ac = new(addContentFunc[NoScriptElement])
	var el = &NoScriptElement{ne, ga, ea, on, ac}
	ga.set(ne, el)
	ea.set(ne, el)
	on.set(ne, el)
	ac.set(ne, el)

	return el
}

// #region <script>

// ScriptElement represents the <script> element.
type ScriptElement struct {
	*element
	*attrAsync[ScriptElement]
	*attrCrossOrigin[ScriptElement]
	*attrDefer[ScriptElement]
	*attrIntegrity[ScriptElement]
	*attrSrc[ScriptElement]
	*attrType[ScriptElement]

	*attrGlobal[ScriptElement]
	*attrExternalAttributes[ScriptElement]
	*attrOn[ScriptElement]
	*addContentFunc[ScriptElement]
}

// Script is used to embed executable code or data; this is typically used to
// embed or refer to JavaScript code. The <script> element can also be used
// with other languages, such as WebGL's GLSL shader programming language and JSON.
//
// For more details, see the [documentation]
//
// [documentation]: https://developer.mozilla.org/en-US/docs/Web/HTML/Element/script
func Script(content ...any) *ScriptElement {
	ne := newElement("script", true, content...)

	var a = new(attrAsync[ScriptElement])
	var b = new(attrCrossOrigin[ScriptElement])
	var c = new(attrDefer[ScriptElement])
	var d = new(attrIntegrity[ScriptElement])
	var e = new(attrSrc[ScriptElement])
	var f = new(attrType[ScriptElement])
	var ga = new(attrGlobal[ScriptElement])
	var ea = new(attrExternalAttributes[ScriptElement])
	var on = new(attrOn[ScriptElement])
	var ac = new(addContentFunc[ScriptElement])
	var el = &ScriptElement{ne, a, b, c, d, e, f, ga, ea, on, ac}
	a.set(ne, el)
	b.set(ne, el)
	c.set(ne, el)
	d.set(ne, el)
	e.set(ne, el)
	f.set(ne, el)
	ga.set(ne, el)
	ea.set(ne, el)
	on.set(ne, el)
	ac.set(ne, el)

	return el
}

// #region DEMARCATING EDITS
// These elements let you provide indications that specific parts of the text
// have been altered.
//

// #region <del>

// DelElement represents the <del> element.
type DelElement struct {
	*element
	*attrCite[DelElement]
	*attrDateTime[DelElement]

	*attrGlobal[DelElement]
	*attrExternalAttributes[DelElement]
	*attrOn[DelElement]
	*addContentFunc[DelElement]
}

// Del represents a range of text that has been deleted from a document. This
// can be used when rendering "track changes" or source code diff information,
// for example. The <ins> element can be used for the opposite purpose: to
// indicate text that has been added to the document.
//
// Example:
//
//	<p><del>This text has been deleted</del>, here is the rest of the paragraph.</p>
//	<del><p>This paragraph has been deleted.</p></del>
//
// For more details, see the [documentation]
//
// [documentation]: https://developer.mozilla.org/en-US/docs/Web/HTML/Element/del
func Del(content ...any) *DelElement {
	ne := newElement("del", true, content...)

	var a = new(attrCite[DelElement])
	var b = new(attrDateTime[DelElement])
	var ga = new(attrGlobal[DelElement])
	var ea = new(attrExternalAttributes[DelElement])
	var on = new(attrOn[DelElement])
	var ac = new(addContentFunc[DelElement])
	var el = &DelElement{ne, a, b, ga, ea, on, ac}
	a.set(ne, el)
	b.set(ne, el)
	ga.set(ne, el)
	ea.set(ne, el)
	on.set(ne, el)
	ac.set(ne, el)

	return el
}

// #region <ins>

// InsElement represents the <ins> element.
type InsElement struct {
	*element
	*attrCite[InsElement]
	*attrDateTime[InsElement]

	*attrGlobal[InsElement]
	*attrExternalAttributes[InsElement]
	*attrOn[InsElement]
	*addContentFunc[InsElement]
}

// Ins represents a range of text that has been added to a document. You can
// use the <del> element to similarly represent a range of text that has been
// deleted from the document.
//
// Example:
//
//	<ins>This text has been inserted</ins>
//
// For more details, see the [documentation]
//
// [documentation]: https://developer.mozilla.org/en-US/docs/Web/HTML/Element/ins
func Ins(content ...any) *InsElement {
	ne := newElement("ins", true, content...)

	var a = new(attrCite[InsElement])
	var b = new(attrDateTime[InsElement])
	var ga = new(attrGlobal[InsElement])
	var ea = new(attrExternalAttributes[InsElement])
	var on = new(attrOn[InsElement])
	var ac = new(addContentFunc[InsElement])
	var el = &InsElement{ne, a, b, ga, ea, on, ac}
	a.set(ne, el)
	b.set(ne, el)
	ga.set(ne, el)
	ea.set(ne, el)
	on.set(ne, el)
	ac.set(ne, el)

	return el
}

// #region TABLE CONTENT
// The elements here are used to create and handle tabular data.
//
//

// #region <caption>

// CaptionElement represents the <caption> element.
type CaptionElement struct {
	*element
	*attrGlobal[CaptionElement]
	*attrExternalAttributes[CaptionElement]
	*attrOn[CaptionElement]
	*addContentFunc[CaptionElement]
}

// Caption specifies the caption (or title) of a table.
//
// Example:
//
//	<table>
//	  <caption>
//	    This is an caption of a table
//	  </caption>
//	  <tr>
//	    <td></td>
//	    <th scope="col" class="heman">He-Man</th>
//	    <th scope="col" class="skeletor">Skeletor</th>
//	  </tr>
//	...
//
// For more details, see the [documentation]
//
// [documentation]: https://developer.mozilla.org/en-US/docs/Web/HTML/Element/caption
func Caption(content ...any) *CaptionElement {
	ne := newElement("caption", true, content...)

	var ga = new(attrGlobal[CaptionElement])
	var ea = new(attrExternalAttributes[CaptionElement])
	var on = new(attrOn[CaptionElement])
	var ac = new(addContentFunc[CaptionElement])
	var el = &CaptionElement{ne, ga, ea, on, ac}
	ga.set(ne, el)
	ea.set(ne, el)
	on.set(ne, el)
	ac.set(ne, el)

	return el
}

// #region <col>

// ColElement represents the <col> element.
type ColElement struct {
	*element
	*attrSpan[ColElement]
	*attrGlobal[ColElement]
	*attrExternalAttributes[ColElement]
	*attrOn[ColElement]
}

// Col defines one or more columns in a column group represented by its
// implicit or explicit parent <colgroup> element.
//
// The <col> element is only valid as a child of a <colgroup> element that has
// no span attribute defined.
//
// Example:
//
//	<table>
//	  <caption>
//	    Superheros and sidekicks
//	  </caption>
//	  <colgroup>
//	    <col />
//	    <col span="2" class="batman" />
//	    <col span="2" class="flash" />
//	  </colgroup>
//
// For more details, see the [documentation]
//
// [documentation]: https://developer.mozilla.org/en-US/docs/Web/HTML/Element/col
func Col() *ColElement {
	ne := newElement("col", false)

	var a = new(attrSpan[ColElement])
	var ga = new(attrGlobal[ColElement])
	var ea = new(attrExternalAttributes[ColElement])
	var on = new(attrOn[ColElement])
	var el = &ColElement{ne, a, ga, ea, on}
	a.set(ne, el)
	ga.set(ne, el)
	ea.set(ne, el)
	on.set(ne, el)

	return el
}

// #region <colgroup>

// ColGroupElement represents the <colgroup> element.
type ColGroupElement struct {
	*element
	*attrSpan[ColGroupElement]

	*attrGlobal[ColGroupElement]
	*attrExternalAttributes[ColGroupElement]
	*attrOn[ColGroupElement]
	*addContentFunc[ColGroupElement]
}

// ColGroup defines a group of columns within a table.
//
// Example:
//
//	<table>
//	  <caption>
//	    Superheros and sidekicks
//	  </caption>
//	  <colgroup>
//	    <col />
//	    <col span="2" class="batman" />
//	    <col span="2" class="flash" />
//	  </colgroup>
//
// For more details, see the [documentation]
//
// [documentation]: https://developer.mozilla.org/en-US/docs/Web/HTML/Element/colgroup
func ColGroup(content ...any) *ColGroupElement {
	ne := newElement("colgroup", true, content...)

	var a = new(attrSpan[ColGroupElement])
	var ga = new(attrGlobal[ColGroupElement])
	var ea = new(attrExternalAttributes[ColGroupElement])
	var on = new(attrOn[ColGroupElement])
	var ac = new(addContentFunc[ColGroupElement])
	var el = &ColGroupElement{ne, a, ga, ea, on, ac}
	a.set(ne, el)
	ga.set(ne, el)
	ea.set(ne, el)
	on.set(ne, el)
	ac.set(ne, el)

	return el
}

// #region <table>

// TableElement represents the <table> element.
type TableElement struct {
	*element
	*attrGlobal[TableElement]
	*attrExternalAttributes[TableElement]
	*attrOn[TableElement]
	*addContentFunc[TableElement]
}

// Table represents tabular data—that is, information presented in a
// two-dimensional table comprised of rows and columns of cells containing data.
//
// For more details, see the [documentation]
//
// [documentation]: https://developer.mozilla.org/en-US/docs/Web/HTML/Element/table
func Table(content ...any) *TableElement {
	ne := newElement("table", true, content...)

	var ga = new(attrGlobal[TableElement])
	var ea = new(attrExternalAttributes[TableElement])
	var on = new(attrOn[TableElement])
	var ac = new(addContentFunc[TableElement])
	var el = &TableElement{ne, ga, ea, on, ac}
	ga.set(ne, el)
	ea.set(ne, el)
	on.set(ne, el)
	ac.set(ne, el)

	return el
}

// #region <tbody>

// TbodyElement represents the <tbody> element.
type TbodyElement struct {
	*element
	*attrGlobal[TbodyElement]
	*attrExternalAttributes[TbodyElement]
	*attrOn[TbodyElement]
	*addContentFunc[TbodyElement]
}

// Tbody encapsulates a set of table rows (<tr> elements), indicating that they
// comprise the body of a table's (main) data.
//
// For more details, see the [documentation]
//
// [documentation]: https://developer.mozilla.org/en-US/docs/Web/HTML/Element/tbody
func Tbody(content ...any) *TbodyElement {
	ne := newElement("tbody", true, content...)

	var ga = new(attrGlobal[TbodyElement])
	var ea = new(attrExternalAttributes[TbodyElement])
	var on = new(attrOn[TbodyElement])
	var ac = new(addContentFunc[TbodyElement])
	var el = &TbodyElement{ne, ga, ea, on, ac}
	ga.set(ne, el)
	ea.set(ne, el)
	on.set(ne, el)
	ac.set(ne, el)

	return el
}

// #region <td>

// TdElement represents the <td> element.
type TdElement struct {
	*element
	*attrColSpan[TdElement]
	*attrHeaders[TdElement]
	*attrRowSpan[TdElement]

	*attrGlobal[TdElement]
	*attrExternalAttributes[TdElement]
	*attrOn[TdElement]
	*addContentFunc[TdElement]
}

// Td is an child of the <tr> element, it defines a cell of a table that
// contains data.
//
// For more details, see the [documentation]
//
// [documentation]: https://developer.mozilla.org/en-US/docs/Web/HTML/Element/td
func Td(content ...any) *TdElement {
	ne := newElement("td", true, content...)

	var a = new(attrColSpan[TdElement])
	var b = new(attrHeaders[TdElement])
	var c = new(attrRowSpan[TdElement])
	var ga = new(attrGlobal[TdElement])
	var ea = new(attrExternalAttributes[TdElement])
	var on = new(attrOn[TdElement])
	var ac = new(addContentFunc[TdElement])
	var el = &TdElement{ne, a, b, c, ga, ea, on, ac}
	a.set(ne, el)
	b.set(ne, el)
	c.set(ne, el)
	ga.set(ne, el)
	ea.set(ne, el)
	on.set(ne, el)
	ac.set(ne, el)

	return el
}

// #region <tfoot>

// TfootElement represents the <tfoot> element.
type TfootElement struct {
	*element
	*attrGlobal[TfootElement]
	*attrExternalAttributes[TfootElement]
	*attrOn[TfootElement]
	*addContentFunc[TfootElement]
}

// Tfoot encapsulates a set of table rows (<tr> elements), indicating that they
// comprise the foot of a table with information about the table's columns.
// This is usually a summary of the columns, e.g., a sum of the given numbers
// in a column.
//
// For more details, see the [documentation]
//
// [documentation]: https://developer.mozilla.org/en-US/docs/Web/HTML/Element/tfoot
func Tfoot(content ...any) *TfootElement {
	ne := newElement("tfoot", true, content...)

	var ga = new(attrGlobal[TfootElement])
	var ea = new(attrExternalAttributes[TfootElement])
	var on = new(attrOn[TfootElement])
	var ac = new(addContentFunc[TfootElement])
	var el = &TfootElement{ne, ga, ea, on, ac}
	ga.set(ne, el)
	ea.set(ne, el)
	on.set(ne, el)
	ac.set(ne, el)

	return el
}

// #region <th>

// ThElement represents the <th> element.
type ThElement struct {
	*element
	*attrAbbr[ThElement]
	*attrColSpan[ThElement]
	*attrHeaders[ThElement]
	*attrRowSpan[ThElement]
	*attrScope[ThElement]

	*attrGlobal[ThElement]
	*attrExternalAttributes[ThElement]
	*attrOn[ThElement]
	*addContentFunc[ThElement]
}

// Th is an child of the <tr> element, it defines a cell as the header of a
// group of table cells. The nature of this group can be explicitly defined by
// the scope and headers attributes.
//
// For more details, see the [documentation]
//
// [documentation]: https://developer.mozilla.org/en-US/docs/Web/HTML/Element/th
func Th(content ...any) *ThElement {
	ne := newElement("th", true, content...)

	var a = new(attrAbbr[ThElement])
	var b = new(attrColSpan[ThElement])
	var c = new(attrHeaders[ThElement])
	var d = new(attrRowSpan[ThElement])
	var e = new(attrScope[ThElement])
	var ga = new(attrGlobal[ThElement])
	var ea = new(attrExternalAttributes[ThElement])
	var on = new(attrOn[ThElement])
	var ac = new(addContentFunc[ThElement])
	var el = &ThElement{ne, a, b, c, d, e, ga, ea, on, ac}
	a.set(ne, el)
	b.set(ne, el)
	c.set(ne, el)
	d.set(ne, el)
	e.set(ne, el)
	ga.set(ne, el)
	ea.set(ne, el)
	on.set(ne, el)
	ac.set(ne, el)

	return el
}

// #region <thead>

// TheadElement represents the <thead> element.
type TheadElement struct {
	*element
	*attrGlobal[TheadElement]
	*attrExternalAttributes[TheadElement]
	*attrOn[TheadElement]
	*addContentFunc[TheadElement]
}

// Thead encapsulates a set of table rows (<tr> elements), indicating that they
// comprise the head of a table with information about the table's columns.
//
// This is usually in the form of column headers (<th> elements).
//
// For more details, see the [documentation]
//
// [documentation]: https://developer.mozilla.org/en-US/docs/Web/HTML/Element/thead
func Thead(content ...any) *TheadElement {
	ne := newElement("thead", true, content...)

	var ga = new(attrGlobal[TheadElement])
	var ea = new(attrExternalAttributes[TheadElement])
	var on = new(attrOn[TheadElement])
	var ac = new(addContentFunc[TheadElement])
	var el = &TheadElement{ne, ga, ea, on, ac}
	ga.set(ne, el)
	ea.set(ne, el)
	on.set(ne, el)
	ac.set(ne, el)

	return el
}

// #region <tr>

// TrElement represents the <tr> element.
type TrElement struct {
	*element
	*attrGlobal[TrElement]
	*attrExternalAttributes[TrElement]
	*attrOn[TrElement]
	*addContentFunc[TrElement]
}

// Tr defines a row of cells in a table. The row's cells can then be
// established using a mix of <td> (data cell) and <th> (header cell) elements.
//
// For more details, see the [documentation]
//
// [documentation]: https://developer.mozilla.org/en-US/docs/Web/HTML/Element/tr
func Tr(content ...any) *TrElement {
	ne := newElement("tr", true, content...)

	var ga = new(attrGlobal[TrElement])
	var ea = new(attrExternalAttributes[TrElement])
	var on = new(attrOn[TrElement])
	var ac = new(addContentFunc[TrElement])
	var el = &TrElement{ne, ga, ea, on, ac}
	ga.set(ne, el)
	ea.set(ne, el)
	on.set(ne, el)
	ac.set(ne, el)

	return el
}

// #region FORMS
// HTML provides several elements that can be used together to create forms that the user can fill out and submit to the
// website or application.
//

// #region <button>

// ButtonElement represents the <button> element.
type ButtonElement struct {
	*element
	*attrDisabled[ButtonElement]
	*attrForm[ButtonElement]
	*attrFormAction[ButtonElement]
	*attrFormEncType[ButtonElement]
	*attrFormMethod[ButtonElement]
	*attrFormNoValidate[ButtonElement]
	*attrFormTarget[ButtonElement]
	*attrName[ButtonElement]
	*attrPopoverTarget[ButtonElement]
	*attrPopoverTargetAction[ButtonElement]
	*attrValue[ButtonElement]

	*attrGlobal[ButtonElement]
	*attrExternalAttributes[ButtonElement]
	*attrOn[ButtonElement]
	*addContentFunc[ButtonElement]
}

// Type specifies the type of an element.
// The default behavior of the button.
//
// Possible values are:
//   - submit: the button submits the form data to the server. This is the default if the attribute is not specified for buttons associated with a <form>, or if the attribute is an empty or invalid value.
//   - reset: the button resets all the controls to their initial values, like <input type="reset">. (This behavior tends to annoy users.)
//   - button: the button has no default behavior, and does nothing when pressed by default. It can have client-side scripts listen to the element's events, which are triggered when the events occur.
func (p *ButtonElement) Type(value string) *ButtonElement {
	p.addAttribute("type", value)
	return p
}

// Button is an interactive element activated by a user with a mouse, keyboard,
// finger, voice command, or other assistive technology. Once activated, it
// performs an action, such as submitting a form or opening a dialog.
//
// For more details, see the [documentation]
//
// [documentation]: https://developer.mozilla.org/en-US/docs/Web/HTML/Element/button
func Button(content ...any) *ButtonElement {
	ne := newElement("button", true, content...)

	var a = new(attrDisabled[ButtonElement])
	var b = new(attrForm[ButtonElement])
	var c = new(attrFormAction[ButtonElement])
	var d = new(attrFormEncType[ButtonElement])
	var e = new(attrFormMethod[ButtonElement])
	var f = new(attrFormNoValidate[ButtonElement])
	var g = new(attrFormTarget[ButtonElement])
	var h = new(attrName[ButtonElement])
	var i = new(attrPopoverTarget[ButtonElement])
	var j = new(attrPopoverTargetAction[ButtonElement])
	var k = new(attrValue[ButtonElement])
	var ga = new(attrGlobal[ButtonElement])
	var ea = new(attrExternalAttributes[ButtonElement])
	var on = new(attrOn[ButtonElement])
	var ac = new(addContentFunc[ButtonElement])
	var el = &ButtonElement{ne, a, b, c, d, e, f, g, h, i, j, k, ga, ea, on, ac}
	a.set(ne, el)
	b.set(ne, el)
	c.set(ne, el)
	d.set(ne, el)
	e.set(ne, el)
	f.set(ne, el)
	g.set(ne, el)
	h.set(ne, el)
	i.set(ne, el)
	j.set(ne, el)
	k.set(ne, el)
	ga.set(ne, el)
	ea.set(ne, el)
	on.set(ne, el)
	ac.set(ne, el)

	return el
}

// #region <datalist>

// DataListElement represents the <datalist> element.
type DataListElement struct {
	*element
	*attrGlobal[DataListElement]
	*attrExternalAttributes[DataListElement]
	*attrOn[DataListElement]
	*addContentFunc[DataListElement]
}

// DataList contains a set of <option> elements that represent the permissible
// or recommended options available to choose from within other controls.
//
// For more details, see the [documentation]
//
// [documentation]: https://developer.mozilla.org/en-US/docs/Web/HTML/Element/datalist
func DataList(content ...any) *DataListElement {
	ne := newElement("datalist", true, content...)

	var ga = new(attrGlobal[DataListElement])
	var ea = new(attrExternalAttributes[DataListElement])
	var on = new(attrOn[DataListElement])
	var ac = new(addContentFunc[DataListElement])
	var el = &DataListElement{ne, ga, ea, on, ac}
	ga.set(ne, el)
	ea.set(ne, el)
	on.set(ne, el)
	ac.set(ne, el)

	return el
}

// #region <fieldset>

// FieldSetElement represents the <fieldset> element.
type FieldSetElement struct {
	*element
	*attrDisabled[FieldSetElement]
	*attrForm[FieldSetElement]
	*attrName[FieldSetElement]

	*attrGlobal[FieldSetElement]
	*attrExternalAttributes[FieldSetElement]
	*attrOn[FieldSetElement]
	*addContentFunc[FieldSetElement]
}

// FieldSet is used to group several controls as well as labels (<label>) within
// a web form.
//
// For more details, see the [documentation]
//
// [documentation]: https://developer.mozilla.org/en-US/docs/Web/HTML/Element/fieldset
func FieldSet(content ...any) *FieldSetElement {
	ne := newElement("fieldset", true, content...)

	var a = new(attrDisabled[FieldSetElement])
	var b = new(attrForm[FieldSetElement])
	var c = new(attrName[FieldSetElement])
	var ga = new(attrGlobal[FieldSetElement])
	var ea = new(attrExternalAttributes[FieldSetElement])
	var on = new(attrOn[FieldSetElement])
	var ac = new(addContentFunc[FieldSetElement])
	var el = &FieldSetElement{ne, a, b, c, ga, ea, on, ac}
	a.set(ne, el)
	b.set(ne, el)
	c.set(ne, el)
	ga.set(ne, el)
	ea.set(ne, el)
	on.set(ne, el)
	ac.set(ne, el)

	return el
}

// #region <form>

// FormElement represents the <form> element.
type FormElement struct {
	*element
	*attrAcceptCharset[FormElement]
	*attrAutoComplete[FormElement]
	*attrName[FormElement]
	*attrRel[FormElement]
	*attrAction[FormElement]
	*attrEncType[FormElement]
	*attrMethod[FormElement]
	*attrNoValidate[FormElement]
	*attrTarget[FormElement]

	*attrGlobal[FormElement]
	*attrExternalAttributes[FormElement]
	*attrOn[FormElement]
	*addContentFunc[FormElement]
}

// Form represents a document section containing interactive controls for
// submitting information.
//
// For more details, see the [documentation]
//
// [documentation]: https://developer.mozilla.org/en-US/docs/Web/HTML/Element/form
func Form(content ...any) *FormElement {
	ne := newElement("form", true, content...)

	var a = new(attrAcceptCharset[FormElement])
	var b = new(attrAutoComplete[FormElement])
	var c = new(attrName[FormElement])
	var d = new(attrRel[FormElement])
	var e = new(attrAction[FormElement])
	var f = new(attrEncType[FormElement])
	var g = new(attrMethod[FormElement])
	var h = new(attrNoValidate[FormElement])
	var i = new(attrTarget[FormElement])
	var ga = new(attrGlobal[FormElement])
	var ea = new(attrExternalAttributes[FormElement])
	var on = new(attrOn[FormElement])
	var ac = new(addContentFunc[FormElement])
	var el = &FormElement{ne, a, b, c, d, e, f, g, h, i, ga, ea, on, ac}
	a.set(ne, el)
	b.set(ne, el)
	c.set(ne, el)
	d.set(ne, el)
	e.set(ne, el)
	f.set(ne, el)
	g.set(ne, el)
	h.set(ne, el)
	i.set(ne, el)
	ga.set(ne, el)
	ea.set(ne, el)
	on.set(ne, el)
	ac.set(ne, el)

	return el
}

// #region <input>

// InputElement represents the <input> element.
type InputElement struct {
	*element
	*attrAccept[InputElement]
	*attrAlt[InputElement]
	*attrAutoComplete[InputElement]
	*attrCapture[InputElement]
	*attrChecked[InputElement]
	*attrDirName[InputElement]
	*attrDisabled[InputElement]
	*attrForm[InputElement]
	*attrFormAction[InputElement]
	*attrFormEncType[InputElement]
	*attrFormMethod[InputElement]
	*attrFormNoValidate[InputElement]
	*attrFormTarget[InputElement]
	*attrList[InputElement]
	*attrMax[InputElement]
	*attrMaxLength[InputElement]
	*attrMin[InputElement]
	*attrMinLength[InputElement]
	*attrName[InputElement]
	*attrPattern[InputElement]
	*attrPlaceholder[InputElement]
	*attrPopoverTarget[InputElement]
	*attrPopoverTargetAction[InputElement]
	*attrReadOnly[InputElement]
	*attrRequired[InputElement]
	*attrSize[InputElement]
	*attrSrc[InputElement]
	*attrStep[InputElement]
	*attrValue[InputElement]

	*attrGlobal[InputElement]
	*attrExternalAttributes[InputElement]
	*attrOn[InputElement]
}

// Type specifies the type of an element.
//
// The available types are as follows:
//   - button: a push button with no default behavior displaying the value of the value attribute, empty by default.
//   - checkbox: A check box allowing single values to be selected/deselected.
//   - color: a control for specifying a color; opening a color picker when active in supporting browsers.
//   - date: a control for entering a date (year, month, and day, with no time). Opens a date picker or numeric wheels for year, month, day when active in supporting browsers.
//   - datetime-local: a control for entering a date and time, with no time zone. Opens a date picker or numeric wheels for date- and time-components when active in supporting browsers.
//   - email: a field for editing an email address. Looks like a text input, but has validation parameters and relevant keyboard in supporting browsers and devices with dynamic keyboards.
//   - file: a control that lets the user select a file. Use the accept attribute to define the types of files that the control can select.
//   - hidden: a control that is not displayed but whose value is submitted to the server. There is an example in the next column, but it's hidden!
//   - image: a graphical submit button. Displays an image defined by the src attribute. The alt attribute displays if the image src is missing.
//   - month: a control for entering a month and year, with no time zone.
//   - number: a control for entering a number. Displays a spinner and adds default validation. Displays a numeric keypad in some devices with dynamic keypads.
//   - password: a single-line text field whose value is obscured. Will alert user if site is not secure.
//   - radio: a radio button, allowing a single value to be selected out of multiple choices with the same name value.
//   - range: a control for entering a number whose exact value is not important. Displays as a range widget defaulting to the middle value. Used in conjunction min and max to define the range of acceptable values.
//   - reset: a button that resets the contents of the form to default values. Not recommended.
//   - search: a single-line text field for entering search strings. Line-breaks are automatically removed from the input value. May include a delete icon in supporting browsers that can be used to clear the field. Displays a search icon instead of enter key on some devices with dynamic keypads.
//   - submit: a button that submits the form.
//   - tel: a control for entering a telephone number. Displays a telephone keypad in some devices with dynamic keypads.
//   - text: the default value. A single-line text field. Line-breaks are automatically removed from the input value.
//   - time: a control for entering a time value with no time zone.
//   - url: a field for entering a URL. Looks like a text input, but has validation parameters and relevant keyboard in supporting browsers and devices with dynamic keyboards.
//   - week: a control for entering a date consisting of a week-year number and a week number with no time zone.
func (p *InputElement) Type(value string) *InputElement {
	p.addAttribute("type", value)
	return p
}

// Input is used to create interactive controls for web-based forms to accept
// data from the user; a wide variety of types of input data and control
// widgets are available, depending on the device and user agent.
//
// The <input> element is one of the most powerful and complex in all of HTML
// due to the sheer number of combinations of input types and attributes.
//
// For more details, see the [documentation]
//
// [documentation]: https://developer.mozilla.org/en-US/docs/Web/HTML/Element/input
func Input() *InputElement {
	ne := newElement("input", false)

	var a = new(attrAccept[InputElement])
	var b = new(attrAlt[InputElement])
	var c = new(attrAutoComplete[InputElement])
	var d = new(attrCapture[InputElement])
	var e = new(attrChecked[InputElement])
	var f = new(attrDirName[InputElement])
	var g = new(attrDisabled[InputElement])
	var h = new(attrForm[InputElement])
	var i = new(attrFormAction[InputElement])
	var j = new(attrFormEncType[InputElement])
	var k = new(attrFormMethod[InputElement])
	var l = new(attrFormNoValidate[InputElement])
	var m = new(attrFormTarget[InputElement])
	var n = new(attrList[InputElement])
	var o = new(attrMax[InputElement])
	var p = new(attrMaxLength[InputElement])
	var q = new(attrMin[InputElement])
	var r = new(attrMinLength[InputElement])
	var s = new(attrName[InputElement])
	var t = new(attrPattern[InputElement])
	var u = new(attrPlaceholder[InputElement])
	var v = new(attrPopoverTarget[InputElement])
	var w = new(attrPopoverTargetAction[InputElement])
	var aa = new(attrReadOnly[InputElement])
	var ab = new(attrRequired[InputElement])
	var ac = new(attrSize[InputElement])
	var ad = new(attrSrc[InputElement])
	var ae = new(attrStep[InputElement])
	var af = new(attrValue[InputElement])

	var ga = new(attrGlobal[InputElement])
	var ea = new(attrExternalAttributes[InputElement])
	var on = new(attrOn[InputElement])
	var el = &InputElement{ne, a, b, c, d, e, f, g, h, i, j, k, l, m, n, o, p, q, r, s, t, u, v, w, aa, ab, ac, ad, ae, af, ga, ea, on}
	a.set(ne, el)
	b.set(ne, el)
	c.set(ne, el)
	d.set(ne, el)
	e.set(ne, el)
	f.set(ne, el)
	g.set(ne, el)
	h.set(ne, el)
	i.set(ne, el)
	j.set(ne, el)
	k.set(ne, el)
	l.set(ne, el)
	m.set(ne, el)
	n.set(ne, el)
	o.set(ne, el)
	p.set(ne, el)
	q.set(ne, el)
	r.set(ne, el)
	s.set(ne, el)
	t.set(ne, el)
	u.set(ne, el)
	v.set(ne, el)
	w.set(ne, el)
	aa.set(ne, el)
	ab.set(ne, el)
	ac.set(ne, el)
	ad.set(ne, el)
	ae.set(ne, el)
	af.set(ne, el)

	ga.set(ne, el)
	ea.set(ne, el)
	on.set(ne, el)

	return el
}

// #region <label>

// LabelElement represents the <label> element.
type LabelElement struct {
	*element
	*attrFor[LabelElement]

	*attrGlobal[LabelElement]
	*attrExternalAttributes[LabelElement]
	*attrOn[LabelElement]
	*addContentFunc[LabelElement]
}

// Label represents a caption for an item in a user interface.
//
// For more details, see the [documentation]
//
// [documentation]: https://developer.mozilla.org/en-US/docs/Web/HTML/Element/label
func Label(content ...any) *LabelElement {
	ne := newElement("label", true, content...)

	var a = new(attrFor[LabelElement])
	var ga = new(attrGlobal[LabelElement])
	var ea = new(attrExternalAttributes[LabelElement])
	var on = new(attrOn[LabelElement])
	var ac = new(addContentFunc[LabelElement])
	var el = &LabelElement{ne, a, ga, ea, on, ac}
	a.set(ne, el)
	ga.set(ne, el)
	ea.set(ne, el)
	on.set(ne, el)
	ac.set(ne, el)

	return el
}

// #region <legend>

// LegendElement represents the <legend> element.
type LegendElement struct {
	*element
	*attrGlobal[LegendElement]
	*attrExternalAttributes[LegendElement]
	*attrOn[LegendElement]
	*addContentFunc[LegendElement]
}

// Legend represents a caption for the content of its parent <fieldset>.
//
// For more details, see the [documentation]
//
// [documentation]: https://developer.mozilla.org/en-US/docs/Web/HTML/Element/legend
func Legend(content ...any) *LegendElement {
	ne := newElement("legend", true, content...)

	var ga = new(attrGlobal[LegendElement])
	var ea = new(attrExternalAttributes[LegendElement])
	var on = new(attrOn[LegendElement])
	var ac = new(addContentFunc[LegendElement])
	var el = &LegendElement{ne, ga, ea, on, ac}
	ga.set(ne, el)
	ea.set(ne, el)
	on.set(ne, el)
	ac.set(ne, el)

	return el
}

// #region meter

// MeterElement represents the <h1> element.
type MeterElement struct {
	*element
	*attrMin[MeterElement]
	*attrMax[MeterElement]
	*attrLow[MeterElement]
	*attrHigh[MeterElement]
	*attrOptimum[MeterElement]
	*attrForm[MeterElement]

	*attrGlobal[MeterElement]
	*attrExternalAttributes[MeterElement]
	*attrOn[MeterElement]
	*addContentFunc[MeterElement]
}

// Meter represents either a scalar value within a known range or a fractional value.
//
// For more details, see the [documentation]
//
// [documentation]: https://developer.mozilla.org/en-US/docs/Web/HTML/Element/meter
func Meter(content ...any) *MeterElement {
	ne := newElement("meter", true, content...)

	var a = new(attrMin[MeterElement])
	var b = new(attrMax[MeterElement])
	var c = new(attrLow[MeterElement])
	var d = new(attrHigh[MeterElement])
	var e = new(attrOptimum[MeterElement])
	var f = new(attrForm[MeterElement])
	var ga = new(attrGlobal[MeterElement])
	var ea = new(attrExternalAttributes[MeterElement])
	var on = new(attrOn[MeterElement])
	var ac = new(addContentFunc[MeterElement])
	var el = &MeterElement{ne, a, b, c, d, e, f, ga, ea, on, ac}
	a.set(ne, el)
	b.set(ne, el)
	c.set(ne, el)
	d.set(ne, el)
	e.set(ne, el)
	f.set(ne, el)
	ga.set(ne, el)
	ea.set(ne, el)
	on.set(ne, el)
	ac.set(ne, el)

	return el
}

// #region <optgroup>

// OptGroupElement represents the <optgroup> element.
type OptGroupElement struct {
	*element
	*attrDisabled[OptGroupElement]
	*attrLabel[OptGroupElement]

	*attrGlobal[OptGroupElement]
	*attrExternalAttributes[OptGroupElement]
	*attrOn[OptGroupElement]
	*addContentFunc[OptGroupElement]
}

// OptGroup creates a grouping of options within a <select> element.
//
// For more details, see the [documentation]
//
// [documentation]: https://developer.mozilla.org/en-US/docs/Web/HTML/Element/optgroup
func OptGroup(content ...any) *OptGroupElement {
	ne := newElement("optgroup", true, content...)

	var a = new(attrDisabled[OptGroupElement])
	var b = new(attrLabel[OptGroupElement])
	var ga = new(attrGlobal[OptGroupElement])
	var ea = new(attrExternalAttributes[OptGroupElement])
	var on = new(attrOn[OptGroupElement])
	var ac = new(addContentFunc[OptGroupElement])
	var el = &OptGroupElement{ne, a, b, ga, ea, on, ac}
	a.set(ne, el)
	b.set(ne, el)
	ga.set(ne, el)
	ea.set(ne, el)
	on.set(ne, el)
	ac.set(ne, el)

	return el
}

// #region <option>

// OptionElement represents the <option> element.
type OptionElement struct {
	*element
	*attrDisabled[OptionElement]
	*attrLabel[OptionElement]
	*attrSelected[OptionElement]
	*attrValue[OptionElement]

	*attrGlobal[OptionElement]
	*attrExternalAttributes[OptionElement]
	*attrOn[OptionElement]
	*addContentFunc[OptionElement]
}

// Option is used to define an item contained in a select, an <optgroup>, or
// a <datalist> element. As such, <option> can represent menu items in popups
// and other lists of items in an HTML document.
//
// For more details, see the [documentation]
//
// [documentation]: https://developer.mozilla.org/en-US/docs/Web/HTML/Element/option
func Option(content ...any) *OptionElement {
	ne := newElement("option", true, content...)

	var a = new(attrDisabled[OptionElement])
	var b = new(attrLabel[OptionElement])
	var c = new(attrSelected[OptionElement])
	var d = new(attrValue[OptionElement])
	var ga = new(attrGlobal[OptionElement])
	var ea = new(attrExternalAttributes[OptionElement])
	var on = new(attrOn[OptionElement])
	var ac = new(addContentFunc[OptionElement])
	var el = &OptionElement{ne, a, b, c, d, ga, ea, on, ac}
	a.set(ne, el)
	b.set(ne, el)
	c.set(ne, el)
	d.set(ne, el)
	ga.set(ne, el)
	ea.set(ne, el)
	on.set(ne, el)
	ac.set(ne, el)

	return el
}

// #region <output>

// OutputElement represents the <output> element.
type OutputElement struct {
	*element
	*attrFor[OutputElement]
	*attrForm[OutputElement]
	*attrName[OutputElement]

	*attrGlobal[OutputElement]
	*attrExternalAttributes[OutputElement]
	*attrOn[OutputElement]
	*addContentFunc[OutputElement]
}

// Output is a container element into which a site or app can inject the results
// of a calculation or the outcome of a user action.
//
// For more details, see the [documentation]
//
// [documentation]: https://developer.mozilla.org/en-US/docs/Web/HTML/Element/output
func Output(content ...any) *OutputElement {
	ne := newElement("output", true, content...)

	var a = new(attrFor[OutputElement])
	var b = new(attrForm[OutputElement])
	var c = new(attrName[OutputElement])
	var ga = new(attrGlobal[OutputElement])
	var ea = new(attrExternalAttributes[OutputElement])
	var on = new(attrOn[OutputElement])
	var ac = new(addContentFunc[OutputElement])
	var el = &OutputElement{ne, a, b, c, ga, ea, on, ac}
	a.set(ne, el)
	b.set(ne, el)
	c.set(ne, el)
	ga.set(ne, el)
	ea.set(ne, el)
	on.set(ne, el)
	ac.set(ne, el)

	return el
}

// #region <progress>

// ProgressElement represents the <progress> element.
type ProgressElement struct {
	*element
	*attrMax[ProgressElement]
	*attrValue[ProgressElement]
	*attrGlobal[ProgressElement]
	*attrExternalAttributes[ProgressElement]
	*attrOn[ProgressElement]
	*addContentFunc[ProgressElement]
}

// Progress displays an indicator showing the completion progress of a task,
// typically displayed as a progress bar.
//
// For more details, see the [documentation]
//
// [documentation]: https://developer.mozilla.org/en-US/docs/Web/HTML/Element/progress
func Progress(content ...any) *ProgressElement {
	ne := newElement("progress", true, content...)

	var a = new(attrMax[ProgressElement])
	var b = new(attrValue[ProgressElement])
	var ga = new(attrGlobal[ProgressElement])
	var ea = new(attrExternalAttributes[ProgressElement])
	var on = new(attrOn[ProgressElement])
	var ac = new(addContentFunc[ProgressElement])
	var el = &ProgressElement{ne, a, b, ga, ea, on, ac}
	a.set(ne, el)
	b.set(ne, el)
	ga.set(ne, el)
	ea.set(ne, el)
	on.set(ne, el)
	ac.set(ne, el)

	return el
}

// #region <select>

// SelectElement represents the <select> element.
type SelectElement struct {
	*element
	*attrAutoComplete[SelectElement]
	*attrDisabled[SelectElement]
	*attrForm[SelectElement]
	*attrMultiple[SelectElement]
	*attrName[SelectElement]
	*attrRequired[SelectElement]
	*attrSize[SelectElement]

	*attrGlobal[SelectElement]
	*attrExternalAttributes[SelectElement]
	*attrOn[SelectElement]
	*addContentFunc[SelectElement]
}

// Select represents a control that provides a menu of options.
//
// Example:
//
//	<label for="hr-select">Your favorite food</label>
//	<select name="foods" id="hr-select">
//	  <option value="">Choose a food</option>
//	  <hr />
//	  <optgroup label="Fruit">
//	    <option value="apple">Apples</option>
//	    <option value="banana">Bananas</option>
//	    <option value="cherry">Cherries</option>
//	    <option value="damson">Damsons</option>
//	  </optgroup>
//	  <hr />
//	  <optgroup label="Vegetables">
//	    <option value="artichoke">Artichokes</option>
//	    <option value="broccoli">Broccoli</option>
//	    <option value="cabbage">Cabbages</option>
//	  </optgroup>
//	  <hr />
//	  <optgroup label="Meat">
//	    <option value="beef">Beef</option>
//	    <option value="chicken">Chicken</option>
//	    <option value="pork">Pork</option>
//	  </optgroup>
//	  <hr />
//	  <optgroup label="Fish">
//	    <option value="cod">Cod</option>
//	    <option value="haddock">Haddock</option>
//	    <option value="salmon">Salmon</option>
//	    <option value="turbot">Turbot</option>
//	  </optgroup>
//	</select>
//
// For more details, see the [documentation]
//
// [documentation]: https://developer.mozilla.org/en-US/docs/Web/HTML/Element/select
func Select(content ...any) *SelectElement {
	ne := newElement("select", true, content...)

	var a = new(attrAutoComplete[SelectElement])
	var b = new(attrDisabled[SelectElement])
	var c = new(attrForm[SelectElement])
	var d = new(attrMultiple[SelectElement])
	var e = new(attrName[SelectElement])
	var f = new(attrRequired[SelectElement])
	var g = new(attrSize[SelectElement])

	var ga = new(attrGlobal[SelectElement])
	var ea = new(attrExternalAttributes[SelectElement])
	var on = new(attrOn[SelectElement])
	var ac = new(addContentFunc[SelectElement])
	var el = &SelectElement{ne, a, b, c, d, e, f, g, ga, ea, on, ac}
	a.set(ne, el)
	b.set(ne, el)
	c.set(ne, el)
	d.set(ne, el)
	e.set(ne, el)
	f.set(ne, el)
	g.set(ne, el)
	ga.set(ne, el)
	ea.set(ne, el)
	on.set(ne, el)
	ac.set(ne, el)

	return el
}

// #region <textarea>

// TextareaElement represents the <textarea> element.
type TextareaElement struct {
	*element
	*attrAutoComplete[TextareaElement]
	*attrCols[TextareaElement]
	*attrDirName[TextareaElement]
	*attrDisabled[TextareaElement]
	*attrForm[TextareaElement]
	*attrMaxLength[TextareaElement]
	*attrMinLength[TextareaElement]
	*attrName[TextareaElement]
	*attrPlaceholder[TextareaElement]
	*attrReadOnly[TextareaElement]
	*attrRequired[TextareaElement]
	*attrRows[TextareaElement]
	*attrWrap[TextareaElement]

	*attrGlobal[TextareaElement]
	*attrExternalAttributes[TextareaElement]
	*attrOn[TextareaElement]
	*addContentFunc[TextareaElement]
}

// Textarea represents a multi-line plain-text editing control, useful when you
// want to allow users to enter a sizeable amount of free-form text, for
// example, a comment on a review or feedback form.
//
// For more details, see the [documentation]
//
// [documentation]: https://developer.mozilla.org/en-US/docs/Web/HTML/Element/textarea
func Textarea(content ...any) *TextareaElement {
	ne := newElement("textarea", true, content...)

	var a = new(attrAutoComplete[TextareaElement])
	var b = new(attrCols[TextareaElement])
	var c = new(attrDirName[TextareaElement])
	var d = new(attrDisabled[TextareaElement])
	var e = new(attrForm[TextareaElement])
	var f = new(attrMaxLength[TextareaElement])
	var g = new(attrMinLength[TextareaElement])
	var h = new(attrName[TextareaElement])
	var i = new(attrPlaceholder[TextareaElement])
	var j = new(attrReadOnly[TextareaElement])
	var k = new(attrRequired[TextareaElement])
	var l = new(attrRows[TextareaElement])
	var m = new(attrWrap[TextareaElement])
	var ga = new(attrGlobal[TextareaElement])
	var ea = new(attrExternalAttributes[TextareaElement])
	var on = new(attrOn[TextareaElement])
	var ac = new(addContentFunc[TextareaElement])
	var el = &TextareaElement{ne, a, b, c, d, e, f, g, h, i, j, k, l, m, ga, ea, on, ac}
	a.set(ne, el)
	b.set(ne, el)
	c.set(ne, el)
	d.set(ne, el)
	e.set(ne, el)
	f.set(ne, el)
	g.set(ne, el)
	h.set(ne, el)
	i.set(ne, el)
	j.set(ne, el)
	k.set(ne, el)
	l.set(ne, el)
	m.set(ne, el)
	ga.set(ne, el)
	ea.set(ne, el)
	on.set(ne, el)
	ac.set(ne, el)

	return el
}

// #region INTERACTIVE ELEMENTS
// HTML offers a selection of elements that help to create interactive user
// interface objects.
//

// #region <details>

// DetailsElement represents the <details> element.
type DetailsElement struct {
	*element
	*attrOpen[DetailsElement]
	*attrName[DetailsElement]

	*attrGlobal[DetailsElement]
	*attrExternalAttributes[DetailsElement]
	*attrOn[DetailsElement]
	*addContentFunc[DetailsElement]
}

// Details creates a disclosure widget in which information is visible only
// when the widget is toggled into an "open" state. A summary or label must be
// provided using the <summary> element.
//
// A <details> widget can be in one of two states. The default closed state
// displays only the triangle and the label inside <summary> (or a user
// agent-defined default string if no <summary>).
//
// When the user clicks on the widget or focuses it then presses the space bar,
// it "twists" open, revealing its contents. The common use of a triangle which
// rotates or twists around to represent opening or closing the widget is why
// these are sometimes called "twisty".
//
// You can use CSS to style the disclosure widget, and you can programmatically
// open and close the widget by setting/removing its open attribute.
// Unfortunately, at this time, there's no built-in way to animate the
// transition between open and closed.
//
// By default when closed, the widget is only tall enough to display the
// disclosure triangle and summary. When open, it expands to display the
// details contained within.
//
// Example:
//
//	<details>
//	  <summary>Details</summary>
//	  Something small enough to escape casual notice.
//	</details>
//
// For more details, see the [documentation]
//
// [documentation]: https://developer.mozilla.org/en-US/docs/Web/HTML/Element/details
func Details(content ...any) *DetailsElement {
	ne := newElement("details", true, content...)

	var a = new(attrOpen[DetailsElement])
	var b = new(attrName[DetailsElement])
	var ga = new(attrGlobal[DetailsElement])
	var ea = new(attrExternalAttributes[DetailsElement])
	var on = new(attrOn[DetailsElement])
	var ac = new(addContentFunc[DetailsElement])
	var el = &DetailsElement{ne, a, b, ga, ea, on, ac}
	a.set(ne, el)
	b.set(ne, el)
	ga.set(ne, el)
	ea.set(ne, el)
	on.set(ne, el)
	ac.set(ne, el)

	return el
}

// #region <dialog>

// DialogElement represents the <dialog> element.
type DialogElement struct {
	*element
	*attrOpen[DialogElement]

	*attrGlobal[DialogElement]
	*attrExternalAttributes[DialogElement]
	*attrOn[DialogElement]
	*addContentFunc[DialogElement]
}

// dialog represents a dialog box or other interactive component, such as a
// dismissible alert, inspector, or subwindow.
//
// Example:
//
//	<dialog open>
//	  <p>Greetings, one and all!</p>
//	  <form method="dialog">
//	    <button>OK</button>
//	  </form>
//	</dialog>
//
// For more details, see the [documentation]
//
// [documentation]: https://developer.mozilla.org/en-US/docs/Web/HTML/Element/dialog
func Dialog(content ...any) *DialogElement {
	ne := newElement("dialog", true, content...)

	var a = new(attrOpen[DialogElement])
	var ga = new(attrGlobal[DialogElement])
	var ea = new(attrExternalAttributes[DialogElement])
	var on = new(attrOn[DialogElement])
	var ac = new(addContentFunc[DialogElement])
	var el = &DialogElement{ne, a, ga, ea, on, ac}
	a.set(ne, el)
	ga.set(ne, el)
	ea.set(ne, el)
	on.set(ne, el)
	ac.set(ne, el)

	return el
}

// #region <summary>

// SummaryElement represents the <summary> element.
type SummaryElement struct {
	*element
	*attrGlobal[SummaryElement]
	*attrExternalAttributes[SummaryElement]
	*attrOn[SummaryElement]
	*addContentFunc[SummaryElement]
}

// Summary specifies a summary, caption, or legend for a details element's
// disclosure box. Clicking the <summary> element toggles the state of the
// parent <details> element open and closed.
//
// For more details, see the [documentation]
//
// [documentation]: https://developer.mozilla.org/en-US/docs/Web/HTML/Element/summary
func Summary(content ...any) *SummaryElement {
	ne := newElement("summary", true, content...)

	var ga = new(attrGlobal[SummaryElement])
	var ea = new(attrExternalAttributes[SummaryElement])
	var on = new(attrOn[SummaryElement])
	var ac = new(addContentFunc[SummaryElement])
	var el = &SummaryElement{ne, ga, ea, on, ac}
	ga.set(ne, el)
	ea.set(ne, el)
	on.set(ne, el)
	ac.set(ne, el)

	return el
}

// #region WEB COMPONENTS
// Web Components is an HTML-related technology that makes it possible to,
// essentially, create and use custom elements as if it were regular HTML.
// In addition, you can create custom versions of standard HTML elements.
//

// #region <slot>

// SlotElement represents the <slot> element.
type SlotElement struct {
	*element
	*attrName[SlotElement]

	*attrGlobal[SlotElement]
	*attrExternalAttributes[SlotElement]
	*attrOn[SlotElement]
	*addContentFunc[SlotElement]
}

// Slot specifies a slot.
//
// For more details, see the [documentation]
//
// [documentation]: https://developer.mozilla.org/en-US/docs/Web/HTML/Element/slot
func Slot(content ...any) *SlotElement {
	ne := newElement("slot", true, content...)

	var a = new(attrName[SlotElement])
	var ga = new(attrGlobal[SlotElement])
	var ea = new(attrExternalAttributes[SlotElement])
	var on = new(attrOn[SlotElement])
	var ac = new(addContentFunc[SlotElement])
	var el = &SlotElement{ne, a, ga, ea, on, ac}
	a.set(ne, el)
	ga.set(ne, el)
	ea.set(ne, el)
	on.set(ne, el)
	ac.set(ne, el)

	return el
}

// #region <template>

// TemplateElement represents the <template> element.
type TemplateElement struct {
	*element

	*attrGlobal[TemplateElement]
	*attrExternalAttributes[TemplateElement]
	*attrOn[TemplateElement]
	*addContentFunc[TemplateElement]
}

// Template specifies a template.
// A mechanism for holding HTML that is not to be rendered immediately when a
// page is loaded but may be instantiated subsequently during runtime using
// JavaScript.
//
// For more details, see the [documentation]
//
// [documentation]: https://developer.mozilla.org/en-US/docs/Web/HTML/Element/template
func Template(content ...any) *TemplateElement {
	ne := newElement("template", true, content...)

	var ga = new(attrGlobal[TemplateElement])
	var ea = new(attrExternalAttributes[TemplateElement])
	var on = new(attrOn[TemplateElement])
	var ac = new(addContentFunc[TemplateElement])
	var el = &TemplateElement{ne, ga, ea, on, ac}
	ga.set(ne, el)
	ea.set(ne, el)
	on.set(ne, el)
	ac.set(ne, el)

	return el
}
