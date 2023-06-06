# Traducciones

Antes de empezar se tiene que instalar la herramienta `gotext`

```shell
go install golang.org/x/text/cmd/gotext@master
```

## Como agregar nuevas traducciones

* Ejecutar `go generate` para extraer los nuevos mensajes no traducidos.
* Sí se muestran traducciones incompletas seguir con los siguientes pasos.
* Copiar `locales/es/out.gotext.json` a `locales/es/messages.gotext.json`.
* Traducir los mensajes faltantes en `locales/es/messages.gotext.json`.
* Ejecutar `go generate` nuevamente para generar los mensajes traducidos.
