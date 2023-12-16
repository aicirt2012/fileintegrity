# License file generation

Install [go-licenses](https://github.com/google/go-licenses) and execute the following command:

```bash
go-licenses report . --template doc/license/license.md.tpl > doc/license/license.md
```