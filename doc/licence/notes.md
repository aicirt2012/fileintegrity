# License file generation

Install [go-licences](https://github.com/google/go-licenses) and execute the following command:

```bash
go-licenses report . --template doc/licences/license.tpl > doc/licences/licence.md
```