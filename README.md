# openapi-httprequest [![Build Status](https://travis-ci.org/mhemmings/openapi-httprequest.svg?branch=master)](https://travis-ci.org/mhemmings/openapi-httprequest)

Tooling to support use of [httprequest](https://github.com/go-httprequest/httprequest) with OpenAPI specifications.

## Usage

At the moment, the command line tool generates a httprequest server from an Open API specification:

`$ openapi-httprequest spec.yaml`

The specification can also be a web url, for example the famous "pet store" example:

`$ openapi-httprequest https://raw.githubusercontent.com/OAI/OpenAPI-Specification/master/examples/v3.0/petstore-expanded.yaml`

The generated code can be ran out of the box using the `--serve` flag. At this point, you basically have the equivalent of a running mock API server (see comment on examples below)

For more docs: `$ openapi-httprequest help`

## Limitations / Bugs

- Requests in the generated API are just the blank values. It would be relatively easy to use the defined example values here from the OpenAPI doc.
- OpenAPI docs are hard to parse (e.g. due to recursive references). Only a subset of common features are supported, and will be added when needed.
- The code is a little scrappy. There must be a more "Go way" of doing things?
- Tests...