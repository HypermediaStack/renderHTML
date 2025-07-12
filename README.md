# renderHTML

This Go package **generates HTML from the server**.

Designed for those who love writing clean, readable code that aligns with the practices of **hypermedia systems**.

Inspired by the simple spirit of the web, it avoids unnecessary dependencies and focuses on making your views as clear and maintainable as your backend.

## Installation

```sh
go get -u github.com/hypermediastack/renderHTML
```

## What does it do?

- Allows you to **create HTML element trees** in a declarative and type-safe way with Go.
- Feels as natural as writing native HTML, but with all the advantages of strong typing, autocomplete, and safe refactoring.
- Makes it easy to manage attributes, classes, styles, and events with intuitive methods.
- Perfect for server-side applications that generate dynamic HTML without depending on complex JavaScript frameworks.

## Example usage
```go
package main

import (
    . "renderHtml"
)

func viewRoot() fmt.Stringer {
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
```

The above Golang code produces the following HTML output:

```html
<!DOCTYPE html><html lang="en"><head><title>Testing app</title><link rel="shortcut icon" href="media/imgs/favicon.ico"/><meta charset="UTF-8"/><meta name="viewport" content="width=device-width, initial-scale=1.0"/><meta http-equiv="X-UA-Compatible" content="IE=edge"/><meta http-equiv="pragma" content="no-cache"/><meta http-equiv="expires" content="-1"/><link rel="stylesheet" href="css/styles.css"/><script defer src="js/htmx.min.js"></script></head><body><h1>Page title</h1><p>Hello world!</p></body></html>
```

## Advantages of using renderHTML

- **100% typed:** Forget about magic strings and silly errors from unclosed tags.
- **Readable and maintainable:** Your HTML code is a tree of Go functions and methods, easy to read and refactor.
- **Total flexibility:** You can add attributes, classes, styles, and dynamic content to any element, using chainable methods.
- **Native events:** Supports `on*` attributes to handle native HTML events.
- **Extensible:** You can add custom attributes and your own structures if you need to extend the package.
- **No magic, no dependencies:** Pure Go, no weird reflection or code injection.
- **No third-party packages:** Uses only Go's standard library, no external dependencies.
- **Hypermedia philosophy:** Designed to build web systems following the classic, proven best practices of the web, keeping full control and transparency over every response.

## How it works

**1.** You can import the package without a prefix:  
```import . "renderHTML"```

This way, every time you write an HTML element name, you don’t have to prepend the package name.

**2.** All HTML attributes are represented. Each one has help text and a reference link to [MDN Web Docs](https://developer.mozilla.org).

**3.** When you type a period (`"."`), your editor will suggest all the attributes available for that element. You can stop worrying about typos. Plus, you don’t have to remember all attributes for every HTML element; this package handles it for you.

**4.** All non-void elements have a method that is not native to HTML, called `AddAttributes(attrs ...any)`.
This method is useful for including external attributes, such as those used by [htmx](https://htmx.org).

Example:

```go
Button("Save").Type("submit").AddAttributes(
    "hx-post='/customer'",
    "hx-target('#div-one')",
    "hx-swap('innerHTML transition:true')",
)
Div().Id("div-one")
```

**5.** All non-void elements have two mechanisms for adding content:
- The traditional way: using the element parameter.
    Example:
    ```go
    H1("Title")
    Div(
        P("Line one"),
        P("Line two"),
    )
    ```

- The second way: using the `AddContent(content ...any)` method.

    Frequently, elements require their own attributes like `id`, `class`, `name`, etc.  
    For this reason, using `AddContent(...)` helps you see the content much like you would in native HTML.

    ```go
    H1().Class("title").Id("page-title").AddContent("Title")
    Div().Class("card").Id("card-one").AddContent(
        P().Id("p-one").Class("help").AddContent("Line one"),
        P("Line two").Id("p-two").Class("help"),
    )
    ```

## Contributing

Suggestions are welcome!

If you have new ideas, find a bug, or simply want to share your hypermedia experience, you are more than welcome.

Please reach out by email: fabianpallares@gmail.com.

## License

This project is licensed under the MIT License.  
See the [LICENSE](LICENSE) file for more information.

> **renderHTML:** Made with love, so the web stays simple, powerful, and fun.
