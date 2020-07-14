## Adminio openAPI v3 specification

### What is OpenAPI?

The OpenAPI Specification (OAS) defines a standard, programming language-agnostic interface description for REST APIs, which allows both humans and computers to discover and understand the capabilities of a service without requiring access to source code, additional documentation, or inspection of network traffic. [more info](http://spec.openapis.org/oas/v3.0.3)

### How to use it?

you can open and try out with this tools:
- [Insomnia designer](https://insomnia.rest/products/designer/)
- [Online swagger edtor](https://editor.swagger.io/)

or can generate standalone html file with [redoc-cli](https://github.com/Redocly/redoc/blob/master/cli/README.md)
```
$ npm install -g redoc-cli
$ redoc-cli bundle -o index.html openapi_v3.yaml
```
