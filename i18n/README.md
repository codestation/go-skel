# Generate and update translations

First install the required tool with the following command:

```shell
go install github.com/nicksnyder/go-i18n/v2/goi18n@latest
```

Then go to the `translations` directory.

```shell
goi18n merge -format yaml active.*.yaml
```


```shell
goi18n merge -format yaml active.*.yaml
```

```shell
goi18n merge -format yaml active.*.yaml translate.*.yaml
```