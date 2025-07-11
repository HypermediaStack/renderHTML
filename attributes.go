package renderHTML

import (
	"fmt"
	"slices"
	"strconv"
	"strings"
)

// #region ATTRIBUTE
type attribute[T any] struct {
	el *element
	t  *T
}

func (p *attribute[T]) set(el *element, t *T) {
	p.el = el
	p.t = t
}

// #region EXTERNAL ATTRS

type attrExternalAttributes[T any] struct {
	attribute[T]
}

// AddAttributes adds attributes to the element. This method is useful for
// including attributes from external packages.
//
// It can receive strings or an object that implements the method:
//
//	interface {GetAttributes() []string} or
//	interface fmt.Stringer {String() string}
//
// The difference is that if the received object implements the
// {GetAttributes() []string} interface, each attribute will be included one by
// one in the affected element.
//
// If it implements the fmt.Stringer interface, which has the String() method,
// the affected element will include all external attributes as a single string.
func (p *attrExternalAttributes[T]) AddAttributes(attrs ...any) *T {
	for _, value := range attrs {
		switch v := value.(type) {
		case interface{ GetAttributes() []string }:
			for _, s := range v.GetAttributes() {
				p.el.addAttribute(s)
			}
		case fmt.Stringer:
			p.el.addAttribute(v.String())
		case string:
			p.el.addAttribute(v)
		default:
			continue
		}
	}

	return p.t
}

// func (p *attrExternalAttributes[T]) AddAttributes(o interface{ GetAttributes() []string }) *T {
// 	for _, attr := range o.GetAttributes() {
// 		p.el.addAttribute(attr)
// 	}

// 	return p.t
// }

// #region GLOBALS ATTRS (G)

type attrGlobal[T any] struct {
	attribute[T]
}

// #region G: accesskey

// AccessKey is an global attribute: specifies a shortcut key to activate or
// focus an element.
func (p *attrGlobal[T]) AccessKey(value string) *T {
	p.el.addAttribute("accesskey", value)
	return p.t
}

// #region G: aria-*

// Aria is an global attribute: Accessible Rich Internet Applications (ARIA)
// is a set of attributes that enhance the accessibility of web content and web
// applications, particularly for users with disabilities.
//
// ARIA attributes provide additional semantics and context that are not always
// available in HTML, allowing assistive technologies to better understand and
// interact with content.
func (p *attrGlobal[T]) Aria(name string, value string) *T {
	p.el.addAttribute("aria-"+name, value)
	return p.t
}

// #region G: autocapitalize

// AutoCapitalize is an global attribute: controls how text input is capitalized.
// The autocapitalize global attribute is an enumerated attribute that controls
// whether inputted text is automatically capitalized and, if so, in what manner.
//
// This is relevant to:
// <input> and <textarea> elements.
//
// Any element with contenteditable set on it.
// Doesn't affect behavior when typing on a physical keyboard. It affects the
// behavior of other input mechanisms such as virtual keyboards on mobile
// devices and voice input. This can assist users by making data entry quicker
// and easier, for example by automatically capitalizing the first letter of
// each sentence.
//
// Possible values are:
//   - none or off: do not automatically capitalize any text.
//   - sentences or on: automatically capitalize the first character of each sentence.
//   - words: automatically capitalize the first character of each word.
//   - characters: automatically capitalize every character.
func (p *attrGlobal[T]) AutoCapitalize(value string) *T {
	p.el.addAttribute("autocapitalize", value)
	return p.t
}

// #region G: autofocus

// AutoFocus indicates that an element is to be focused on page load, or as soon
// as the <dialog> it is part of is displayed. This attribute is a boolean,
// initially false.
func (p *attrGlobal[T]) AutoFocus(value ...bool) *T {
	if value == nil {
		p.el.addAttribute("autofocus")
	} else {
		p.el.addAttribute("autofocus", strconv.FormatBool(value[0]))
	}
	return p.t
}

// #region G: class

// Class is an global attribute: specifies one or more class names for an element.
func (p *attrGlobal[T]) Class(text ...string) *T {
	p.el.addClasses(text...)
	return p.t
}

// #region G: contenteditable

// ContentEditable is an global attribute: indicates whether the element's
// content is editable.
// The contenteditable global attribute is an enumerated attribute indicating
// if the element should be editable by the user. If so, the browser modifies
// its widget to allow editing.
//
// Possible values are:
//   - true or an empty string: which indicates that the element is editable.
//   - false: which indicates that the element is not editable.
//   - plaintext-only: which indicates that the element's raw text is editable, but rich text formatting is disabled.
func (p *attrGlobal[T]) ContentEditable(value ...string) *T {
	if value == nil {
		p.el.addAttribute("contenteditable")
		return p.t
	}

	var v = strings.TrimSpace(strings.ToLower(value[0]))
	if !slices.Contains([]string{"true", "false", "plaintext-only"}, v) {
		p.el.addAttribute("contenteditable")
	} else {
		p.el.addAttribute("contenteditable", v)
	}

	return p.t
}

// #region G: data-*

// Data (data-*) is an global attribute: specifies custom data attributes.
//
// Lets you attach custom attributes to an HTML element.
func (p *attrGlobal[T]) Data(name string, value string) *T {
	p.el.addAttribute("data-"+name, value)
	return p.t
}

// #region G: dir

// Dir is an global attribute: defines the text direction. Allowed values
// are ltr (Left-To-Right) or rtl (Right-To-Left).
// The dir global attribute is an enumerated attribute that indicates the
// directionality of the element's text.
//
// It can have the following values:
//   - ltr: which means left to right and is to be used for languages that are written from the left to the right (like English);
//   - rtl, which means right to left and is to be used for languages that are written from the right to the left (like Arabic);
//   - auto: which lets the user agent decide. It uses a basic algorithm as it parses the characters inside the element until it finds a character with a strong directionality, then applies that directionality to the whole element.
func (p *attrGlobal[T]) Dir(value string) *T {
	p.el.addAttribute("dir", value)
	return p.t
}

// #region G: draggable

// Draggable is an global attribute: specifies whether an element is draggable
// or not.
// The draggable global attribute is an enumerated attribute that indicates
// whether the element can be dragged, either with native browser behavior or
// the HTML Drag and Drop API.
//
// The draggable attribute may be applied to elements that strictly fall under
// the HTML namespace, which means that it cannot be applied to SVGs.
//
// For more information about what namespace declarations look like, and what
// they do, see Namespace crash course.
//
// Warning: This attribute is enumerated and not Boolean. A value of true or
// false is mandatory, and shorthand like <img draggable> is forbidden.
// The correct usage is <img draggable="false">
//
// Can have the following values:
//   - true: the element can be dragged.
//   - false: the element cannot be dragged.
func (p *attrGlobal[T]) Draggable(value bool) *T {
	p.el.addAttribute("draggable", value)
	return p.t
}

// #region G: enterkeyhint

// EnterKeyHint is an global attribute: hints what action label (or icon) to
// present for the enter key on virtual keyboards.
//
// The enterkeyhint attribute is an enumerated attribute and only accepts the
// following values:
//   - enterkeyhint="enter": Typically inserting a new line.
//   - enterkeyhint="done": Typically meaning there is nothing more to input and the input method editor (IME) will be closed.
//   - enterkeyhint="go": Typically meaning to take the user to the target of the text they typed.
//   - enterkeyhint="next": Typically taking the user to the next field that will accept text.
//   - enterkeyhint="previous": Typically taking the user to the previous field that will accept text.
//   - enterkeyhint="search": Typically taking the user to the results of searching for the text they have typed.
//   - enterkeyhint="send": Typically delivering the text to its target.
//
// Example:
//
//	<input enterkeyhint="go" />
//	<p contenteditable enterkeyhint="go">https://example.org</p>
func (p *attrGlobal[T]) EnterKeyHint() *T {
	p.el.addAttribute("enterkeyhint")
	return p.t
}

// #region G: hidden

// Hidden is an global attribute: is used to indicate that the content of an
// element should not be presented to the user.
func (p *attrGlobal[T]) Hidden() *T {
	p.el.addAttribute("hidden")
	return p.t
}

// #region G: id

// Id is an global attribute: specifies a unique id for an element.
func (p *attrGlobal[T]) Id(value string) *T {
	p.el.addAttribute("id", value)
	return p.t
}

// #region G: inert

// Inert is an global attribute: is a boolean value that makes the browser
// disregard user input events for the element. Useful when click events are
// present.
//
// Is a Boolean attribute indicating that the browser will ignore the element.
// With the inert attribute, all of the element's flat tree descendants (such
// as modal <dialog>s) that don't otherwise escape inertness are ignored. The
// inert attribute also makes the browser ignore input events sent by the user,
// including focus-related events and events from assistive technologies.
func (p *attrGlobal[T]) Inert(value ...bool) *T {
	if value == nil {
		p.el.addAttribute("inert")
	} else {
		p.el.addAttribute("inert", strconv.FormatBool(value[0]))
	}
	return p.t
}

// #region G: inputmode

// Inputmode is an global attribute: provides a hint to browsers about the
// type of virtual keyboard configuration to use when editing this element or
// its contents. Used primarily on <input> elements, but is usable on any element
// while in contenteditable mode.
//
// The attribute can have any of the following values:
//   - none: no virtual keyboard. For when the page implements its own keyboard input control.
//   - text (default value): standard input keyboard for the user's current locale.
//   - decimal: fractional numeric input keyboard containing the digits and decimal separator for the user's locale (typically . or ,). Devices may or may not show a minus key (-).
//   - numeric: numeric input keyboard, but only requires the digits 0–9. Devices may or may not show a minus key.
//   - tel: a telephone keypad input, including the digits 0–9, the asterisk (*), and the pound (#) key. Inputs that require a telephone number should typically use <input type="tel"> instead.
//   - search: a virtual keyboard optimized for search input. For instance, the return/submit key may be labeled "Search", along with possible other optimizations. Inputs that require a search query should typically use <input type="search"> instead.
//   - email: a virtual keyboard optimized for entering email addresses. Typically includes the @character as well as other optimizations. Inputs that require email addresses should typically use <input type="email"> instead.
//   - url: a keypad optimized for entering URLs. This may have the / key more prominent, for example. Enhanced features could include history access and so on. Inputs that require a URL should typically use <input type="url"> instead.
func (p *attrGlobal[T]) Inputmode(value string) *T {
	p.el.addAttribute("inputmode", value)
	return p.t
}

// #region G: itemprop

// ItemProp is an global attribute: specifies a property of an item in an item.
// The itemprop global attribute is used to add properties to an item.
//
// Every HTML element can have an itemprop attribute specified, and an itemprop
// consists of a name-value pair.
// Each name-value pair is called a property, and a group of one or more
// properties forms an item.
// Property values are either a string or a URL and can be associated with a
// very wide range of elements including <audio>, <embed>, <iframe>, <img>,
// <link>, <object>, <source>, <track>, and <video>.
func (p *attrGlobal[T]) ItemProp(value string) *T {
	p.el.addAttribute("itemprop", value)
	return p.t
}

// #region G: lang

// Lang is an global attribute: specifies the language of the element's content.
func (p *attrGlobal[T]) Lang(value string) *T {
	p.el.addAttribute("lang", value)
	return p.t
}

// #region G: role

// Role is an global attribute: defines an explicit role for an element for
// use by assistive technologies.
func (p *attrGlobal[T]) Role(value string) *T {
	p.el.addAttribute("role", value)
	return p.t
}

// #region G: spellcheck

// SpellCheck is an global attribute: indicates whether spell checking is
// allowed for the element.
// If this attribute is not set, its default value is element-type and browser-defined.
//
// This default value may also be inherited, which means that the element
// content will be checked for spelling errors only if its nearest ancestor
// has a spellcheck state of true.
//
// It may have the following values:
//   - empty string or true: which indicates that the element should be, if possible, checked for spelling errors.
//   - false: which indicates that the element should not be checked for spelling errors.
func (p *attrGlobal[T]) SpellCheck(value ...bool) *T {
	if value == nil {
		p.el.addAttribute("spellcheck")
	} else {
		p.el.addAttribute("spellcheck", strconv.FormatBool(value[0]))
	}
	return p.t
}

// #region G: style

// Style is an global attribute: specifies inline CSS styles for an element.
func (p *attrGlobal[T]) Style(text ...string) *T {
	p.el.addStyles(text...)
	return p.t
}

// #region G: tabindex

// TabIndex is an global attribute: specifies the tab order of an element.
func (p *attrGlobal[T]) TabIndex(value int) *T {
	p.el.addAttribute("tabindex", value)
	return p.t
}

// #region G: title

// Title is an global attribute: specifies extra information about an element.
func (p *attrGlobal[T]) Title(value string) *T {
	p.el.addAttribute("title", value)
	return p.t
}

// #region G: translate

// Translate is an global attribute: specify whether an element's attribute
// values and the values of its Text node children are to be translated when
// the page is localized, or whether to leave them unchanged.
//
// The translate global attribute is an enumerated attribute that is used to
// specify whether an element's translatable attribute values and its Text node
// children should be translated when the page is localized, or whether to leave
// them unchanged.
//
// It can have the following values:
//   - empty string or yes: which indicates that the element should be translated when the page is localized.
//   - no: which indicates that the element must not be translated.
func (p *attrGlobal[T]) Translate(value string) *T {
	p.el.addAttribute("translate", value)
	return p.t
}

// #region HTML ATTRS
//
//
//

// #region accept
// <form>, <input>

type attrAccept[T any] struct {
	attribute[T]
}

// Accept specifies the types of files that the server accepts.
func (p *attrAccept[T]) Accept(value string) *T {
	p.el.addAttribute("accept", value)
	return p.t
}

// #region accept-charset
// <form>

type attrAcceptCharset[T any] struct {
	attribute[T]
}

// AcceptCharSet specifies the character encodings that are to be used for the
// form submission.
func (p *attrAcceptCharset[T]) AcceptCharSet(value string) *T {
	p.el.addAttribute("accept-charset", value)
	return p.t
}

// #region action
// <form>

type attrAction[T any] struct {
	attribute[T]
}

// Action specifies the URL where the form data will be sent.
func (p *attrAction[T]) Action(value string) *T {
	p.el.addAttribute("action", value)
	return p.t
}

// #region allow
// <iframe>

type attrAllow[T any] struct {
	attribute[T]
}

// Allow specifies a feature-policy for the iframe.
func (p *attrAllow[T]) Allow(value string) *T {
	p.el.addAttribute("allow", value)
	return p.t
}

// #region alt
// <area>, <img>, <input>

type attrAlt[T any] struct {
	attribute[T]
}

// Alt provides alternative text for an image.
func (p *attrAlt[T]) Alt(value string) *T {
	p.el.addAttribute("alt", value)
	return p.t
}

// #region as
// <link>

type attrAs[T any] struct {
	attribute[T]
}

// As specifies the type of content being loaded by the link.
// This attribute is required when rel="preload" has been set on the <link>
// element, optional when rel="modulepreload" has been set, and otherwise should
// not be used.
//
// It specifies the type of content being loaded by the <link>, which
// is necessary for request matching, application of correct content security
// policy, and setting of correct Accept request header.
//
// Furthermore, rel="preload" uses this as a signal for request prioritization.
func (p *attrAs[T]) As(value string) *T {
	p.el.addAttribute("as", value)
	return p.t
}

// #region async
// <script>

type attrAsync[T any] struct {
	attribute[T]
}

// Async specifies that the script is to be executed asynchronously.
func (p *attrAsync[T]) Async() *T {
	p.el.addAttribute("async")
	return p.t
}

// #region autocomplete
// <form>, <input>, <select>, <textarea>

type attrAutoComplete[T any] struct {
	attribute[T]
}

// AutoComplete controls whether the browser should autocomplete input fields.
// The autocomplete attribute provides a hint to the user agent specifying how
// to, or indeed whether to, prefill a form control.
//
// The attribute value is either the keyword off or on, or an ordered list
// of space-separated tokens.
//
// Example:
//
//	<input autocomplete="off" />
//	<input autocomplete="on" />
//	<input autocomplete="shipping street-address" />
//	<input autocomplete="section-user1 billing postal-code" />
func (p *attrAutoComplete[T]) AutoComplete(value string) *T {
	p.el.addAttribute("autocomplete", value)
	return p.t
}

// #region autoplay
// <audio>, <video>

type attrAutoPlay[T any] struct {
	attribute[T]
}

// AutoPlay specifies that the audio/video should start playing automatically.
func (p *attrAutoPlay[T]) AutoPlay() *T {
	p.el.addAttribute("autoplay")
	return p.t
}

type attrAbbr[T any] struct {
	attribute[T]
}

// Abbr is a short, abbreviated description of the header cell's content
// provided as an alternative label to use for the header cell when referencing
// the cell in other contexts. Some user-agents, such as speech readers, may
// present this description before the content itself.
func (p *attrAbbr[T]) Abbr() *T {
	p.el.addAttribute("abbr")
	return p.t
}

// #region capture
// <input>

type attrCapture[T any] struct {
	attribute[T]
}

// Capture media capture input method in file upload controls.
func (p *attrCapture[T]) Capture(value string) *T {
	p.el.addAttribute("capture", value)
	return p.t
}

// #region charset
// <meta>

type attrCharSet[T any] struct {
	attribute[T]
}

// CharSet specifies the character encoding for the HTML document. Its value
// must be an ASCII case-insensitive match for the string "utf-8", because
// UTF-8 is the only valid encoding for HTML5 documents.
func (p *attrCharSet[T]) CharSet(value string) *T {
	p.el.addAttribute("charset", value)
	return p.t
}

// #region checked
// <input>

type attrChecked[T any] struct {
	attribute[T]
}

// Checked specifies that an input element should be pre-selected.
func (p *attrChecked[T]) Checked() *T {
	p.el.addAttribute("checked")
	return p.t
}

// #region cite
// <blockquote>, <del>, <ins>, <q>

type attrCite[T any] struct {
	attribute[T]
}

// Cite contains a URI which points to the source of the quote or change.
func (p *attrCite[T]) Cite(value string) *T {
	p.el.addAttribute("cite", value)
	return p.t
}

// #region cols
// <textarea>

type attrCols[T any] struct {
	attribute[T]
}

// Cols specifies the number of visible columns in a text area.
func (p *attrCols[T]) Cols(value int) *T {
	p.el.addAttribute("cols", value)
	return p.t
}

// #region colspan
// <td>, <th>

type attrColSpan[T any] struct {
	attribute[T]
}

// ColSpan specifies the number of columns a table cell should span.
func (p *attrColSpan[T]) ColSpan(value int) *T {
	p.el.addAttribute("colspan", value)
	return p.t
}

// #region content
// <meta>

type attrMetaContent[T any] struct {
	attribute[T]
}

// Content specifies the value of a meta element. A value associated with
// "http-equiv" or "name" depending on the context.
func (p *attrMetaContent[T]) Content(value string) *T {
	p.el.addAttribute("content", value)
	return p.t
}

// #region controls
// <audio>, <video>

type attrControls[T any] struct {
	attribute[T]
}

// Controls specifies that audio/video controls should be displayed.
func (p *attrControls[T]) Controls() *T {
	p.el.addAttribute("controls")
	return p.t
}

// #region controlslist
// <video>

type attrControlsList[T any] struct {
	attribute[T]
}

// ControlsList when specified, helps the browser select what controls to show
// for the video element whenever the browser shows its own set of controls
// (that is, when the controls attribute is specified).
//
// The allowed values are nodownload, nofullscreen and noremoteplayback.
// Use the disablepictureinpicture attribute if you want to disable the
// Picture-In-Picture mode (and the control).
func (p *attrControlsList[T]) ControlsList() *T {
	p.el.addAttribute("controlslist")
	return p.t
}

// #region coords
// <area>

type attrCoords[T any] struct {
	attribute[T]
}

// Coords details the coordinates of the shape attribute in size, shape, and
// placement of an <area>. This attribute must not be used if shape is set to
// default.
func (p *attrCoords[T]) Coords(value string) *T {
	p.el.addAttribute("coords", value)
	return p.t
}

// #region crossorigin
// <audio>, <img>, <link>, <script>, <video>

type attrCrossOrigin[T any] struct {
	attribute[T]
}

// CrossOrigin how the element handles cross-origin requests.
func (p *attrCrossOrigin[T]) CrossOrigin(value string) *T {
	p.el.addAttribute("crossorigin", value)
	return p.t
}

// #region data
// <object>

type attrData[T any] struct {
	attribute[T]
}

// Data (data="text") specifies the URL of the resource.
func (p *attrData[T]) DataURL(value string) *T {
	p.el.addAttribute("data", value)
	return p.t
}

// #region datetime
// <del>, <ins>, <time>

type attrDateTime[T any] struct {
	attribute[T]
}

// DateTime specifies a specific date and time.
func (p *attrDateTime[T]) DateTime(value string) *T {
	p.el.addAttribute("datetime", value)
	return p.t
}

// #region decoding
// <img>

type attrDecoding[T any] struct {
	attribute[T]
}

// Decoding indicates the preferred method to decode the image.
func (p *attrDecoding[T]) Decoding(value string) *T {
	p.el.addAttribute("decoding", value)
	return p.t
}

// #region default
// <track>

type attrDefault[T any] struct {
	attribute[T]
}

// Default indicates that the track should be enabled unless the user's
// preferences indicate something different.
func (p *attrDefault[T]) Default() *T {
	p.el.addAttribute("default")
	return p.t
}

// #region defer
// <script>

type attrDefer[T any] struct {
	attribute[T]
}

// Defer indicates that the script should be executed after the page has been
// parsed.
func (p *attrDefer[T]) Defer() *T {
	p.el.addAttribute("defer")
	return p.t
}

// #region dirname
// <input>, <textarea>

type attrDirName[T any] struct {
	attribute[T]
}

// DirName defines the text direction. Allowed values are ltr (Left-To-Right)
// or rtl (Right-To-Left).
// The dirname attribute is an enumerated attribute that indicates the
// directionality of the element's text.
//
// It can have the following values:
//   - ltr: which means left to right and is to be used for languages that are written from the left to the right (like English);
//   - rtl, which means right to left and is to be used for languages that are written from the right to the left (like Arabic);
//   - auto: which lets the user agent decide. It uses a basic algorithm as it parses the characters inside the element until it finds a character with a strong directionality, then applies that directionality to the whole element.
func (p *attrDirName[T]) DirName(value string) *T {
	p.el.addAttribute("dirname", value)
	return p.t
}

// #region disabled
// <button>, <fieldset>, <input>, <optgroup>, <option>, <select>, <textarea>

type attrDisabled[T any] struct {
	attribute[T]
}

// Disabled specifies that an element should be disabled.
func (p *attrDisabled[T]) Disabled() *T {
	p.el.addAttribute("disabled")
	return p.t
}

// #region disablepictureinpicture
// <video>

type attrDisablePictureInPicture[T any] struct {
	attribute[T]
}

// DisablePictureInPicture prevents the browser from suggesting a
// Picture-in-Picture context menu or to request Picture-in-Picture
// automatically in some cases.
func (p *attrDisablePictureInPicture[T]) DisablePictureInPicture() *T {
	p.el.addAttribute("disablepictureinpicture")
	return p.t
}

// #region disableremoteplayback
// <video>

type attrDisableRemotePlayBack[T any] struct {
	attribute[T]
}

// DisablePictureAndPicture ia an boolean attribute used to disable the
// capability of remote playback in devices that are attached using wired
// (HDMI, DVI, etc.) and wireless technologies (Miracast, Chromecast, DLNA,
// AirPlay, etc.).
func (p *attrDisableRemotePlayBack[T]) DisableRemotePlayback(value bool) *T {
	p.el.addAttribute("disableremoteplayback", value)
	return p.t
}

// #region download
// <a>, <area>

type attrDownload[T any] struct {
	attribute[T]
}

// Download specifies that the target should be downloaded when a user clicks
// on the hyperlink.
func (p *attrDownload[T]) Download(value string) *T {
	p.el.addAttribute("download", value)
	return p.t
}

// #region enctype
// <form>

type attrEncType[T any] struct {
	attribute[T]
}

// EncType specifies how the form data should be encoded when submitting to the
// server. Defines the content type of the form data when the method is POST.
//
// If the value of the method attribute is post, enctype is the MIME type of
// the form submission.
//
// Possible values:
//   - application/x-www-form-urlencoded: the default value.
//   - multipart/form-data: use this if the form contains <input> elements with type=file.
//   - text/plain: useful for debugging purposes.
//
// This value can be overridden by formenctype attributes on <button>, <input type="submit">, or <input type="image"> elements.
func (p *attrEncType[T]) EncType(value string) *T {
	p.el.addAttribute("enctype", value)
	return p.t
}

// #region for
// <label>, <output>

type attrFor[T any] struct {
	attribute[T]
}

// For specifies which form element a label is bound to.
func (p *attrFor[T]) For(value string) *T {
	p.el.addAttribute("for", value)
	return p.t
}

// #region <form>
// <button>, <fieldset>, <input>, <label>, <meter>, <object>, <output>, <progress>, <select>, <textarea>

type attrForm[T any] struct {
	attribute[T]
}

// Form associates an element with a form. Indicates the form that is the owner
// of the element.
func (p *attrForm[T]) Form(value string) *T {
	p.el.addAttribute("form", value)
	return p.t
}

// #region formaction
// <input>, <button>

type attrFormAction[T any] struct {
	attribute[T]
}

// FormAction specifies the URL for form submission.
func (p *attrFormAction[T]) FormAction(value string) *T {
	p.el.addAttribute("formaction", value)
	return p.t
}

// #region formenctype
// <button>, <input>

type attrFormEncType[T any] struct {
	attribute[T]
}

// FormEncType specifies how the form data should be encoded when submitting to
// the server.
// If the button/input is a submit button (e.g. type="submit"), this attribute
// sets the encoding type to use during form submission.
//
// If this attribute is specified, it overrides the enctype attribute of the
// button's form owner.
func (p *attrFormEncType[T]) FormEncType(value string) *T {
	p.el.addAttribute("formenctype", value)
	return p.t
}

// #region formmethod
// <button>, <input>

type attrFormMethod[T any] struct {
	attribute[T]
}

// FormMethod specifies the HTTP method to use when submitting form data.
// If the button/input is a submit button (e.g. type="submit"), this attribute
// sets the submission method to use during form submission (GET, POST, etc.).
//
// If this attribute is specified, it overrides the method attribute of the
// button's form owner.
func (p *attrFormMethod[T]) FormMethod(value string) *T {
	p.el.addAttribute("formmethod", value)
	return p.t
}

// #region formnovalidate
// <button>, <input>

type attrFormNoValidate[T any] struct {
	attribute[T]
}

// FormNoValidate specifies that the form should not be validated when submitted.
// If the button/input is a submit button (e.g. type="submit"), this boolean
// attribute specifies that the form is not to be validated when it is submitted.
//
// If this attribute is specified, it overrides the novalidate attribute of the
// button's form owner.
//
// This attribute is also available on <input type="image"> and
// <input type="submit"> elements.
func (p *attrFormNoValidate[T]) FormNoValidate() *T {
	p.el.addAttribute("formnovalidate")
	return p.t
}

// #region formtarget
// <button>, <input>

type attrFormTarget[T any] struct {
	attribute[T]
}

// FormTarget specifies where to display the response after submitting the form.
// If the button/input is a submit button (e.g. type="submit"), this attribute
// specifies the browsing context (for example, tab, window, or inline frame) in
// which to display the response that is received after submitting the form.
//
// If this attribute is specified, it overrides the target attribute of the
// button's form owner.
func (p *attrFormTarget[T]) FormTarget(value string) *T {
	p.el.addAttribute("formtarget", value)
	return p.t
}

// #region headers
// <td>, <th>

type attrHeaders[T any] struct {
	attribute[T]
}

// Headers indicates IDs of the <th> elements which applies to this element.
// Contains a list of space-separated strings, each corresponding to the id
// attribute of the <th> elements that provide headings for this table cell.
func (p *attrHeaders[T]) Headers(value string) *T {
	p.el.addAttribute("headers", value)
	return p.t
}

// #region high
// <meter>

type attrHigh[T any] struct {
	attribute[T]
}

// High specifies the range of the gauge. Indicates the lower bound of the upper
// range.
func (p *attrHigh[T]) High(value string) *T {
	p.el.addAttribute("high", value)
	return p.t
}

// #region href
// <a>, <area>, <base>, <link>

type attrHref[T any] struct {
	attribute[T]
}

// Href specifies the URL of the page the link goes to.
func (p *attrHref[T]) Href(value string) *T {
	p.el.addAttribute("href", value)
	return p.t
}

// #region hreflang
// <a>, <link>

type attrHrefLang[T any] struct {
	attribute[T]
}

// HrefLang specifies the language of the linked document.
func (p *attrHrefLang[T]) HrefLang(value string) *T {
	p.el.addAttribute("hreflang", value)
	return p.t
}

// #region http-equiv
// <meta>

type attrHttpEquiv[T any] struct {
	attribute[T]
}

// HttpEquiv provides an HTTP header for the value of the content attribute.
// Defines a pragma directive.
func (p *attrHttpEquiv[T]) HttpEquiv(value string) *T {
	p.el.addAttribute("http-equiv", value)
	return p.t
}

// #region integrity
// <link>, <script>

type attrIntegrity[T any] struct {
	attribute[T]
}

// Integrity specifies a Subresource Integrity value that allows browsers to
// verify what they fetch.
func (p *attrIntegrity[T]) Integrity(value string) *T {
	p.el.addAttribute("integrity", value)
	return p.t
}

// #region ismap
// <img>

type attrIsmap[T any] struct {
	attribute[T]
}

// Ismap specifies that an image is part of a client-side image map.
func (p *attrIsmap[T]) Ismap() *T {
	p.el.addAttribute("ismap")
	return p.t
}

// #region kind
// <track>

type attrKind[T any] struct {
	attribute[T]
}

// Kind represents how the text track is meant to be used. If omitted the
// default kind is subtitles. If the attribute contains an invalid value, it
// will use metadata.
//
// The following keywords are allowed:
//   - subtitles: subtitles provide translation of content that cannot be understood by the viewer. For example speech or text that is not English in an English language film. Subtitles may contain additional content, usually extra background information. For example the text at the beginning of the Star Wars films, or the date, time, and location of a scene.
//   - captions: closed captions provide a transcription and possibly a translation of audio. It may include important non-verbal information such as music cues or sound effects. It may indicate the cue's source (e.g. music, text, character). Suitable for users who are deaf or when the sound is muted.
//   - chapters: chapter titles are intended to be used when the user is navigating the media resource.
//   - metadata: tracks used by scripts. Not visible to the user.
func (p *attrKind[T]) Kind(value string) *T {
	p.el.addAttribute("kind", value)
	return p.t
}

// #region label
// <optgroup>, <option>, <track>

type attrLabel[T any] struct {
	attribute[T]
}

// Label specifies the label of a track element.
func (p *attrLabel[T]) Label(value string) *T {
	p.el.addAttribute("label", value)
	return p.t
}

// #region loading
// <img>, <iframe>

type attrLoading[T any] struct {
	attribute[T]
}

// Loading indicates if the element should be loaded lazily (loading="lazy")
// or loaded immediately (loading="eager").
func (p *attrLoading[T]) Loading(value string) *T {
	p.el.addAttribute("loading", value)
	return p.t
}

// #region list
// <input>

type attrList[T any] struct {
	attribute[T]
}

// List associates an input field with a datalist element.
func (p *attrList[T]) List(value string) *T {
	p.el.addAttribute("list", value)
	return p.t
}

// #region loop
// <audio>, <marquee>, <video>

type attrLoop[T any] struct {
	attribute[T]
}

// Loop specifies that the media should start over again when it reaches the end.
func (p *attrLoop[T]) Loop() *T {
	p.el.addAttribute("loop")
	return p.t
}

// #region low
// <meter>

type attrLow[T any] struct {
	attribute[T]
}

// Low specifies the lower bound of the gauge.
func (p *attrLow[T]) Low(value int) *T {
	p.el.addAttribute("low", value)
	return p.t
}

// #region max
// <input>, <meter>, <progress>

type attrMax[T any] struct {
	attribute[T]
}

// Max specifies the maximum value of an element.
// The max attribute defines the maximum value that is acceptable and valid for
// the input containing the attribute.
//
// If the value of the element is greater than this, the element fails validation.
// This value must be greater than or equal to the value of the min attribute.
// If the max attribute is present but is not specified or is invalid, no max
// value is applied.
// If the max attribute is valid and a non-empty value is greater than the
// maximum allowed by the max attribute, constraint validation will prevent
// form submission.
//
// Valid for the numeric input types, including the date, month, week, time,
// datetime-local, number and range types, and both the <progress> and <meter>
// elements, the max attribute is a number that specifies the most positive
// value a form control to be considered valid.
func (p *attrMax[T]) Max(value string) *T {
	p.el.addAttribute("max", value)
	return p.t
}

// #region maxlength
// <input>, <textarea>

type attrMaxLength[T any] struct {
	attribute[T]
}

// MaxLength specifies the maximum number of characters allowed in an input field.
func (p *attrMaxLength[T]) MaxLength(value int) *T {
	p.el.addAttribute("maxlength", value)
	return p.t
}

// #region minlength
// <input>, <textarea>

type attrMinLength[T any] struct {
	attribute[T]
}

// MinLength specifies the minimum number of characters allowed in an input field.
func (p *attrMinLength[T]) MinLength(value int) *T {
	p.el.addAttribute("minlength", value)
	return p.t
}

// #region media
// <a>, <area>, <link>, <source>, <style>

type attrMedia[T any] struct {
	attribute[T]
}

// Media specifies what media/device the linked resource is optimized for.
func (p *attrMedia[T]) Media(value string) *T {
	p.el.addAttribute("media", value)
	return p.t
}

// #region method
// <form>

type attrMethod[T any] struct {
	attribute[T]
}

// Method specifies the HTTP method to use when submitting form data.
// The HTTP method to submit the form with.
//
// The only allowed methods/values are (case insensitive):
//   - post: the POST method; form data sent as the request body.
//   - get (default): the GET; form data appended to the action URL with a ? separator. Use this method when the form has no side effects.
//   - dialog: when the form is inside a <dialog>, closes the dialog and causes a submit event to be fired on submission, without submitting data or clearing the form.
//
// This value is overridden by formmethod attributes on <button>, <input type="submit">, or <input type="image"> elements.
func (p *attrMethod[T]) Method(value string) *T {
	p.el.addAttribute("method", value)
	return p.t
}

// #region min
// <input>, <meter>

type attrMin[T any] struct {
	attribute[T]
}

// Min specifies the minimum value of an element.
// It is valid for the input types including: date, month, week, time,
// datetime-local, number and range types, and the <meter> element.
func (p *attrMin[T]) Min(value string) *T {
	p.el.addAttribute("min", value)
	return p.t
}

// #region multiple
// <input>, <select>

type attrMultiple[T any] struct {
	attribute[T]
}

// Multiple specifies that multiple options can be selected.
func (p *attrMultiple[T]) Multiple() *T {
	p.el.addAttribute("multiple")
	return p.t
}

// #region muted
// <audio>, <video>

type attrMuted[T any] struct {
	attribute[T]
}

// Muted specifies that the audio/video should be muted.
func (p *attrMuted[T]) Muted() *T {
	p.el.addAttribute("muted")
	return p.t
}

// #region name
// <button>, <form>, <fieldset>, <iframe>, <input>, <object>, <output>, <select>,
// <textarea>, <map>, <meta>

type attrName[T any] struct {
	attribute[T]
}

// Name specifies a name for an element.
func (p *attrName[T]) Name(value string) *T {
	p.el.addAttribute("name", value)
	return p.t
}

// #region novalidate
// <form>

type attrNoValidate[T any] struct {
	attribute[T]
}

// NoValidate indicates that the form shouldn't be validated when submitted.
func (p *attrNoValidate[T]) NoValidate() *T {
	p.el.addAttribute("novalidate")
	return p.t
}

// #region open
// <details>, <dialog>

type attrOpen[T any] struct {
	attribute[T]
}

// Open indicates whether the contents are currently visible (in the case of
// a <details> element) or whether the dialog is active and can be interacted
// with (in the case of a <dialog> element).
func (p *attrOpen[T]) Open(value ...bool) *T {
	if value == nil {
		p.el.addAttribute("open")
	} else {
		p.el.addAttribute("open", strconv.FormatBool(value[0]))
	}

	return p.t
}

// #region optimum
// <meter>

type attrOptimum[T any] struct {
	attribute[T]
}

// Optimum indicates the optimal numeric value.
func (p *attrOptimum[T]) Optimum(value int) *T {
	p.el.addAttribute("optimum", value)
	return p.t
}

// #region pattern
// <input>

type attrPattern[T any] struct {
	attribute[T]
}

// Pattern specifies a regular expression that the input value must match.
// The pattern attribute is an attribute of the text, tel, email, url, password,
// and search input types.
//
// The pattern attribute, when specified, is a regular expression which the
// input's value must match for the value to pass constraint validation.
//
// It must be a valid JavaScript regular expression, as used by the RegExp type,
// and as documented in our guide on regular expressions.
//
// The pattern's regular expression is compiled with the 'v' flag.
// This makes the regular expression unicode-aware, and also changes how
// character classes are interpreted.
//
// This allows character class set intersection and subtraction operations,
// and in addition to ] and \, the following characters must be escaped using
// a \ backslash if they represent literal characters: (, ), [, {, }, /, -, |.
func (p *attrPattern[T]) Pattern(value string) *T {
	p.el.addAttribute("pattern", value)
	return p.t
}

// #region ping
// <a>, <area>

type attrPing[T any] struct {
	attribute[T]
}

// Ping specifies a space-separated list of URLs to be notified if a user
// follows the hyperlink.
func (p *attrPing[T]) Ping(value string) *T {
	p.el.addAttribute("ping", value)
	return p.t
}

// #region placeholder
// <input>, <textarea>

type attrPlaceholder[T any] struct {
	attribute[T]
}

// Placeholder specifies a short hint that describes the expected value of an
// input field.
func (p *attrPlaceholder[T]) Placeholder(value string) *T {
	p.el.addAttribute("placeholder", value)
	return p.t
}

// #region playsinline
// <video>

type attrPlaysInLine[T any] struct {
	attribute[T]
}

// PlaysInLine indicating that the video is to be played "inline"; that is,
// within the element's playback area. Note that the absence of this attribute
// does not imply that the video will always be played in fullscreen.
func (p *attrPlaysInLine[T]) PlaysInLine(value bool) *T {
	p.el.addAttribute("playsinline", value)
	return p.t
}

// #region poster
// <video>

type attrPoster[T any] struct {
	attribute[T]
}

// Poster specifies an image to be shown while the video is downloading, or
// until the user hits the play button.
func (p *attrPoster[T]) Poster(value string) *T {
	p.el.addAttribute("poster", value)
	return p.t
}

// #region preload
// <audio>, <video>

type attrPreLoad[T any] struct {
	attribute[T]
}

// PreLoad specifies if and how the media file should be loaded when the page loads.
func (p *attrPreLoad[T]) PreLoad(value string) *T {
	p.el.addAttribute("preload", value)
	return p.t
}

// #region popovertarget
// <button>

type attrPopoverTarget[T any] struct {
	attribute[T]
}

// PopoverTarget turns a element into a popover control; takes the ID of the
// popover element to control as its value.
//
// See the [Popover API landing page] for more details.
//
// [Popover API landing page]: https://developer.mozilla.org/en-US/docs/Web/API/Popover_API
func (p *attrPopoverTarget[T]) PopoverTarget(value string) *T {
	p.el.addAttribute("popovertarget", value)
	return p.t
}

// #region popovertargetaction
// <button>

type attrPopoverTargetAction[T any] struct {
	attribute[T]
}

// PopoverTargetAction specifies the action to be performed on a popover element.
//
// Possible values are:
//   - "hide": the button will hide a shown popover. If you try to hide an already hidden popover, no action will be taken.
//   - "show": the button will show a hidden popover. If you try to show an already showing popover, no action will be taken.
//   - "toggle": the button will toggle a popover between showing and hidden. If the popover is hidden, it will be shown; if the popover is showing, it will be hidden. If popovertargetaction is omitted, "toggle" is the default action that will be performed by the control button.
//
// [Popover API landing page]: https://developer.mozilla.org/en-US/docs/Web/API/Popover_API
func (p *attrPopoverTargetAction[T]) PopoverTargetAction(value string) *T {
	p.el.addAttribute("popovertargetaction", value)
	return p.t
}

// #region readonly
// <input>, <textarea>

type attrReadOnly[T any] struct {
	attribute[T]
}

// ReadOnly specifies that an input field is read-only.
func (p *attrReadOnly[T]) ReadOnly() *T {
	p.el.addAttribute("readonly")
	return p.t
}

// #region rel
// <a>, <area>, <link>

type attrRel[T any] struct {
	attribute[T]
}

// Rel specifies the relationship between the current document and the linked
// document.
func (p *attrRel[T]) Rel(value string) *T {
	p.el.addAttribute("rel", value)
	return p.t
}

// #region required
// <input>, <select>, <textarea>

type attrRequired[T any] struct {
	attribute[T]
}

// Required specifies that an input field must be filled out before submitting
// the form.
func (p *attrRequired[T]) Required() *T {
	p.el.addAttribute("required")
	return p.t
}

// #region reversed
// <ol>

type attrReversed[T any] struct {
	attribute[T]
}

// Reversed indicates whether the list should be displayed in a descending order
// instead of an ascending order.
func (p *attrReversed[T]) Reversed() *T {
	p.el.addAttribute("reversed")
	return p.t
}

// #region rows
// <textarea>

type attrRows[T any] struct {
	attribute[T]
}

// Rows specifies the number of visible rows in a text area.
func (p *attrRows[T]) Rows(value int) *T {
	p.el.addAttribute("rows", value)
	return p.t
}

// #region rowspan
// <td>, <th>

type attrRowSpan[T any] struct {
	attribute[T]
}

// RowSpan specifies the number of rows a table cell should span.
func (p *attrRowSpan[T]) RowSpan(value int) *T {
	p.el.addAttribute("rowspan", value)
	return p.t
}

// #region sandbox
// <iframe>

type attrSandbox[T any] struct {
	attribute[T]
}

// Sandbox enables an extra set of restrictions for the content in an iframe.
func (p *attrSandbox[T]) Sandbox(value string) *T {
	p.el.addAttribute("sandbox", value)
	return p.t
}

// #region scope
// <th>

type attrScope[T any] struct {
	attribute[T]
}

// Scope specifies whether a header cell is a header for a row, column, or group
// of rows or columns.
//
// Defines the cells that the header (defined in the <th>) element relates to.
//
// Possible enumerated values are:
//   - row: the header relates to all cells of the row it belongs to.
//   - col: the header relates to all cells of the column it belongs to.
//   - rowgroup: the header belongs to a rowgroup and relates to all of its cells
//   - colgroup: the header belongs to a colgroup and relates to all of its cells.
//
// If the scope attribute is not specified, or its value is not row, col, rowgroup, or colgroup, then browsers automatically select the set of cells to which the header cell applies.
func (p *attrScope[T]) Scope(value string) *T {
	p.el.addAttribute("scope", value)
	return p.t
}

// #region selected
// <option>

type attrSelected[T any] struct {
	attribute[T]
}

// Selected specifies that an option should be pre-selected when the page loads.
func (p *attrSelected[T]) Selected() *T {
	p.el.addAttribute("selected")
	return p.t
}

// #region shape
// <a>, <area>

type attrShape[T any] struct {
	attribute[T]
}

// Shape requests a value that is advisable to obtain from official documentation.
func (p *attrShape[T]) Shape() *T {
	p.el.addAttribute("shape")
	return p.t
}

// #region size
// <input>, <area>

type attrSize[T any] struct {
	attribute[T]
}

// Size specifies the width of an input element or the number of visible options
// in a select element.
func (p *attrSize[T]) Size(value int) *T {
	p.el.addAttribute("size", value)
	return p.t
}

// #region sizes
// <link>, <img>, <source>

type attrSizes[T any] struct {
	attribute[T]
}

// Sizes specifies a list of source sizes that describe the final rendered
// width of the image. Allowed if the parent of <source> is <picture>. Not
// allowed if the parent is <audio> or <video>.
//
// The list consists of source sizes separated by commas. Each source size is
// media condition-length pair. Before laying the page out, the browser uses
// this information to determine which image defined in srcset to display.
//
// Note that sizes will take effect only if width descriptors are provided with
// srcset, not pixel density descriptors (i.e., 200w should be used instead of 2x).
func (p *attrSizes[T]) Sizes(value string) *T {
	p.el.addAttribute("sizes", value)
	return p.t
}

// #region span
// <col>, <colgroup>

type attrSpan[T any] struct {
	attribute[T]
}

// Span specifies the number of columns a table cell should span.
func (p *attrSpan[T]) Span(value int) *T {
	p.el.addAttribute("span", value)
	return p.t
}

// #region src
// <audio>, <embed>, <iframe>, <img>, <input>, <script>, <source>, <track>, <video>

type attrSrc[T any] struct {
	attribute[T]
}

// Src specifies the URL of the media file or script to be used.
func (p *attrSrc[T]) Src(value string) *T {
	p.el.addAttribute("src", value)
	return p.t
}

// #region srclang
// <track>

type attrSrcLang[T any] struct {
	attribute[T]
}

// SrcLang specifies the language of the track text data. It must be a valid
// BCP 47 language tag. If the kind attribute is set to subtitles, then srcLang
// must be defined.
func (p *attrSrcLang[T]) SrcLang(value string) *T {
	p.el.addAttribute("srclang", value)
	return p.t
}

// #region srcset
// <img>, <source>

type attrSrcSet[T any] struct {
	attribute[T]
}

// SrcSet specifies multiple sources for responsive images.
func (p *attrSrcSet[T]) SrcSet(value string) *T {
	p.el.addAttribute("srcset", value)
	return p.t
}

// #region start
// <ol>

type attrStart[T any] struct {
	attribute[T]
}

// Start specifies the legal number intervals for an input field.
func (p *attrStart[T]) Start(value int) *T {
	p.el.addAttribute("start", value)
	return p.t
}

// #region step
// <input>

type attrStep[T any] struct {
	attribute[T]
}

// Step specifies the legal number intervals for an input field.
func (p *attrStep[T]) Step(value int) *T {
	p.el.addAttribute("step", value)
	return p.t
}

// #region target
// <a>, <area>, <base>, <form>

type attrTarget[T any] struct {
	attribute[T]
}

// Target specifies where to open the linked document (in the case of an
// <a> element) or where to display the response received (in the case of a
// <form> element).
//
// The following keywords have special meanings:
//   - _self (default): show the result in the current browsing context.
//   - _blank: show the result in a new, unnamed browsing context.
//   - _parent: show the result in the parent browsing context of the current one, if the current page is inside a frame. If there is no parent, acts the same as _self.
//   - _top: show the result in the topmost browsing context (the browsing context that is an ancestor of the current one and has no parent). If there is no parent, acts the same as _self.
func (p *attrTarget[T]) Target(value string) *T {
	p.el.addAttribute("target", value)
	return p.t
}

// #region type
// <button>, <input>, <embed>, <object>, <ol>, <script>, <source>, <style>,
// <menu>, <link>

type attrType[T any] struct {
	attribute[T]
}

// Type specifies the type of an element.
func (p *attrType[T]) Type(value string) *T {
	p.el.addAttribute("type", value)
	return p.t
}

// #region value
// <button>, <data>, <input>, <li>, <meter>, <option>, <progress>

type attrValue[T any] struct {
	attribute[T]
}

// Value specifies the value of an element.
func (p *attrValue[T]) Value(value string) *T {
	p.el.addAttribute("value", value)
	return p.t
}

// #region wrap
// <textarea>

type attrWrap[T any] struct {
	attribute[T]
}

// Wrap specifies how the text in a text area is wrapped.
// Indicates how the control should wrap the value for form submission.
//
// Possible values are:
//   - hard: the browser automatically inserts line breaks (CR+LF) so that each line is no longer than the width of the control; the cols attribute must be specified for this to take effect.
//   - soft: the browser ensures that all line breaks in the entered value are a CR+LF pair, but no additional line breaks are added to the value.
func (p *attrWrap[T]) Wrap(value string) *T {
	p.el.addAttribute("wrap", value)
	return p.t
}

// #region event handlers (on*)

type attrOn[T any] struct {
	attribute[T]
}

// On specifies a JavaScript code to be executed when the event occurs.
//
// Example:
//
//	<button type="button" onclick="alert('Hello world!');">Hello</button>
//
// For more details, see the [event reference] or [documentation]
//
// [event reference]: https://developer.mozilla.org/en-US/docs/Web/Events
// [documentation]: https://developer.mozilla.org/en-US/docs/Web/API/Element#events
func (p *attrOn[T]) On(eventName, script string) *T {
	eventName = strings.ToLower(eventName)
	if eventName[:2] == "on" {
		p.el.addAttribute(eventName, script)
	} else {
		p.el.addAttribute("on"+eventName, script)
	}

	return p.t
}
