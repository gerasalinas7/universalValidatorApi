# Universal Validator API

Microservicio REST genérico desarrollado en Go. Recibe un tipo de dato y un valor, y devuelve si cumple las reglas de validación.

## Requisitos

- Go 1.22 o superior

## Ejecutar localmente

```bash
go run .
```

El servicio queda disponible en `http://localhost:8080`.

También se puede configurar el puerto:

```bash
PORT=9000 go run .
```

## Endpoints

### Health check

```http
GET /health
```

### Tipos soportados

```http
GET /api/v1/validate/types
```

Respuesta:

```json
{
  "types": ["email", "cuit", "phone", "url", "password"]
}
```

### Validar un dato

```http
POST /api/v1/validate
Content-Type: application/json
```

Request:

```json
{
  "type": "email",
  "value": "alumno@gmail.com"
}
```

Respuesta válida:

```json
{
  "valid": true,
  "type": "email",
  "value": "alumno@gmail.com",
  "errors": []
}
```

Respuesta inválida:

```json
{
  "valid": false,
  "type": "email",
  "value": "alumno@gmail",
  "errors": ["Email address has an invalid format"]
}
```

## Ejemplos con curl

```bash
curl http://localhost:8080/health

curl http://localhost:8080/api/v1/validate/types

curl -X POST http://localhost:8080/api/v1/validate \
  -H "Content-Type: application/json" \
  -d '{"type":"email","value":"alumno@gmail.com"}'

curl -X POST http://localhost:8080/api/v1/validate \
  -H "Content-Type: application/json" \
  -d '{"type":"password","value":"ClaveSegura1!"}'
```

## Validaciones incluidas

- `email`: formato de email usando la biblioteca estándar de Go.
- `cuit`: 11 dígitos y dígito verificador según módulo 11.
- `phone`: formato telefónico general; no verifica que el número esté asignado.
- `url`: URL absoluta con esquema HTTP o HTTPS.
- `password`: mínimo 8 caracteres, mayúscula, minúscula, número y carácter especial.

> Importante: estas validaciones comprueban formato y reglas básicas. No verifican identidad, titularidad, existencia real de cuentas ni que un teléfono esté activo.

## Ejecutar tests

```bash
go test ./...
```

## Docker

```bash
docker build -t universal-validator-api .
docker run --rm -p 8080:8080 universal-validator-api
```

## Contrato de errores

- `200`: validación ejecutada. El campo `valid` indica el resultado.
- `400`: JSON inválido, campos requeridos ausentes o tipo no soportado.
- `GET /health`: permite comprobar que el proceso está activo.
