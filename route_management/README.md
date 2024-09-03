# Nombre del Proyecto

El servicio de gestión de trayectos permite crear trayectos (rutas) para ser usados por las publicaciones. Este fue desarrollado en Golang con el framework Gin y el ORM gorm.

## Índice

1. [Estructura](#estructura)
2. [Ejecución](#ejecución)
3. [Uso](#uso)
4. [Pruebas](#pruebas)
5. [Otras caracteristicas](#Otras\_Características)
6. [Autor](#autor)

## Estructura

Describe la estructura de archivos de la carpeta, puedes usar una estructura de arbol para ello:
```
route_management
├── mocks # Carpeta con archivo con mocks de servicios para ser usados en las pruebas
│
├── model # Carpeta con modelos de gorm con sus respectivas funciones para transformar los datos a la base de datos
│
├── routes # Carpeta con archivos que tienen las funciones que las rutas que usa la aplicación
│   └── delete_routes.go # archivo con las rutas que sean de tipo DELETE
│   └── get_routes.go # archivo con las rutas que sean de tipo GET
│   └── post_routes.go # archivo con las rutas que sean de tipo POST
|
├── utils # Carpeta con utilidades compartidas por diferentes 
│   └── auth.go # archivo con la autenticación y sus respectivas interfaces e utilidades
│   └── utils.go # archivo con la utilidades usadas por varias rutas y partes de la aplicación
│
├── Dockerfile # Dockerfile base con la configuración para desplegar el servicio de manera segura
│
├── go.mod # archivo con las librerías usadas por la aplicación y la versión de golang
│
├── go.sum # archivo con las urls de las librerías y sus dependecias usadas por la aplicación
│
├── init.sql # sql usado para inicializar la base de datos
|
├── main_test.go # archivo con las pruebas de la aplicación a nivel de main.go
|
├── main.go # archivo con el entrypoint de la aplicación. genera la base de datos, el servidor y las rutas. 
|
├── test.sh # archivo para correr las pruebas en local o en los pipelines. 
|
└── README.md # usted está aquí
```
## Ejecución

Proporciona instrucciones claras y concisas sobre cómo instalar y ejecutar este proyecto. Si es necesario, incluye ejemplos de código o comandos.

```bash
go mod download
go build -o main .
./main
...
```

## Uso

Para consumir la API se necesita lo siguiente: 

### Variables de entorno

```bash
DB_USER=routes_management_user
DB_PASSWORD=routes_management_pass
DB_NAME=routes_management_db
DB_HOST=routes_db
DB_PORT=5432
CONFIG_PORT=3000
USERS_PATH="http://users:3000"
```

### Servicios

Tener una instancia del sevicio de users corriendo para poder obtener y validar los tokens. Si un token no es válido no se podrá usar el servicio de manera efectiva. 

El servicio de gestión de trayectos permite crear trayectos (rutas) para ser usados por las publicaciones.

#### 1. Creación de trayecto

Crea un trayecto con los datos brindados, solo un usuario autorizado puede realizar esta operación.

<table>
<tr>
<td> Método </td>
<td> POST </td>
</tr>
<tr>
<td> Ruta </td>
<td> <strong>/routes</strong> </td>
</tr>
<tr>
<td> Parámetros </td>
<td> N/A </td>
</tr>
<tr>
<td> Encabezados </td>
<td>

```Authorization: Bearer token```
</td>
</tr>
<tr>
<td> Cuerpo </td>
<td>

```json
{
  "flightId": código del vuelo,
  "sourceAirportCode": código del aeropuerto de origen,
  "sourceCountry": nombre del país de origen,
  "destinyAirportCode": código del aeropuerto de destino,
  "destinyCountry": nombre del país de destino,
  "bagCost": costo de envío de maleta,
  "plannedStartDate": fecha y hora de inicio del trayecto,
  "plannedEndDate": fecha y hora de finalización del trayecto
}
```
</td>
</tr>
</table>

##### Respuestas

<table>
<tr>
<th> Código </th>
<th> Descripción </th>
<th> Cuerpo </th>
</tr>
<tbody>
<tr>
<td> 401 </td>
<td>El token no es válido o está vencido.</td>
<td> N/A </td>
</tr>
<tr>
<td> 403 </td>
<td>No hay token en la solicitud</td>
<td> N/A </td>
</tr>
<tr>
<td> 400 </td>
<td>En el caso que alguno de los campos no esté presente en la solicitud.</td>
<td> N/A </td>
</tr>
<tr>
<td> 412 </td>
<td>En el caso que el flightId ya existe.</td>
<td> N/A </td>
</tr>
<tr>
<td> 412 </td>
<td> En el caso que la fecha de inicio y fin del trayecto no sean válidas; fechas en el pasado o no consecutivas.</td>
<td> 

```json
{
    "msg": "Las fechas del trayecto no son válidas"
}
```
</td>
</tr>
<tr>
<td> 201 </td>
<td>En el caso que el trayecto se haya creado con éxito.</td>
<td>

```json
{
  "id": id del trayecto,
  "createdAt": fecha y hora de creación del trayecto en formato ISO
}
```
</td>
</tr>
</tbody>
</table>

#### 2. Ver y filtrar trayectos

Retorna todos los trayectos o aquellos que corresponden a los parámetros de búsqueda. Solo un usuario autorizado puede realizar esta operación.

<table>
<tr>
<td> Método </td>
<td> GET </td>
</tr>
<tr>
<td> Ruta </td>
<td> <strong>/routes?flight={flightId}</strong> </td>
</tr>
<tr>
<td> Parámetros </td>
<td> flightId: id del vuelo, este campo es opcional </td>
</tr>
<tr>
<td> Encabezados </td>
<td>

```Authorization: Bearer token```
</td>
</tr>
<tr>
<td> Cuerpo </td>
<td> N/A </td>
</tr>
</table>

##### Respuestas

<table>
<tr>
<th> Código </th>
<th> Descripción </th>
<th> Cuerpo </th>
</tr>
<tbody>
<tr>
<td> 401 </td>
<td>El token no es válido o está vencido.</td>
<td> N/A </td>
</tr>
<tr>
<td> 403 </td>
<td>No hay token en la solicitud</td>
<td> N/A </td>
</tr>
<tr>
<td> 400 </td>
<td>En el caso que alguno de los campos de búsqueda no tenga el formato esperado.</td>
<td> N/A </td>
</tr>
<tr>
<td> 200 </td>
<td>Listado de trayectos.</td>
<td>

```json
[
  {
    "id": id del trayecto
    "flightId": código del vuelo,
    "sourceAirportCode": código del aeropuerto de origen,
    "sourceCountry": nombre del país de origen,
    "destinyAirportCode": código del aeropuerto de destino,
    "destinyCountry": nombre del país de destino,
    "bagCost": costo de envío de maleta,
    "plannedStartDate": fecha y hora de inicio del trayecto,
    "plannedEndDate": fecha y hora de finalización del trayecto,
    "createdAt": fecha y hora de creación del trayecto en formato ISO
  }
]
```
</td>
</tr>
</tbody>
</table>


#### 3. Consultar un trayecto

Retorna un trayecto, solo un usuario autorizado puede realizar esta operación.

<table>
<tr>
<td> Método </td>
<td> GET </td>
</tr>
<tr>
<td> Ruta </td>
<td> <strong>/routes/{id}</strong> </td>
</tr>
<tr>
<td> Parámetros </td>
<td> id: identificador del trayecto </td>
</tr>
<tr>
<td> Encabezados </td>
<td>

```Authorization: Bearer token```
</td>
</tr>
<tr>
<td> Cuerpo </td>
<td> N/A </td>
</tr>
</table>

##### Respuestas

<table>
<tr>
<th> Código </th>
<th> Descripción </th>
<th> Cuerpo </th>
</tr>
<tbody>
<tr>
<td> 401 </td>
<td>El token no es válido o está vencido.</td>
<td> N/A </td>
</tr>
<tr>
<td> 403 </td>
<td>No hay token en la solicitud</td>
<td> N/A </td>
</tr>
<tr>
<td> 400 </td>
<td>El id no es un valor string con formato uuid.</td>
<td> N/A </td>
</tr>
</tr>
<tr>
<td> 404 </td>
<td>El trayecto con ese id no existe.</td>
<td> N/A </td>
</tr>
<tr>
<td> 200 </td>
<td>Trayecto que corresponde al identificador.</td>
<td>

```json
  {
    "id": id del trayecto
    "flightId": código del vuelo,
    "sourceAirportCode": código del aeropuerto de origen,
    "sourceCountry": nombre del país de origen,
    "destinyAirportCode": código del aeropuerto de destino,
    "destinyCountry": nombre del país de destino,
    "bagCost": costo de envío de maleta,
    "plannedStartDate": fecha y hora de inicio del trayecto,
    "plannedEndDate": fecha y hora de finalización del trayecto,
    "createdAt": fecha y hora de creación del trayecto en formato ISO
  }
```
</td>
</tr>
</tbody>
</table>


#### 4. Eliminar trayecto

Elimina el trayecto, solo un usuario autorizado puede realizar esta operación.

<table>
<tr>
<td> Método </td>
<td> DELETE </td>
</tr>
<tr>
<td> Ruta </td>
<td> <strong>/routes/{id}</strong> </td>
</tr>
<tr>
<td> Parámetros </td>
<td> id: identificador del trayecto </td>
</tr>
<tr>
<td> Encabezados </td>
<td>

```Authorization: Bearer token```
</td>
</tr>
<tr>
<td> Cuerpo </td>
<td> N/A </td>
</tr>
</table>

##### Respuestas

<table>
<tr>
<th> Código </th>
<th> Descripción </th>
<th> Cuerpo </th>
</tr>
<tbody>
<tr>
<td> 401 </td>
<td>El token no es válido o está vencido.</td>
<td> N/A </td>
</tr>
<tr>
<td> 403 </td>
<td>No hay token en la solicitud</td>
<td> N/A </td>
</tr>
<tr>
<td> 400 </td>
<td>El id no es un valor string con formato uuid.</td>
<td> N/A </td>
</tr>
</tr>
<tr>
<td> 404 </td>
<td>El trayecto con ese id no existe.</td>
<td> N/A </td>
</tr>
<tr>
<td> 200 </td>
<td>El trayecto fue eliminado.</td>
<td>

```json
  {
    "msg": "el trayecto fue eliminado"
  }
```
</td>
</tr>
</tbody>
</table>


#### 5. Consulta de salud del servicio

Usado para verificar el estado del servicio.

<table>
<tr>
<td> Método </td>
<td> GET </td>
</tr>
<tr>
<td> Ruta </td>
<td> <strong>/routes/ping</strong> </td>
</tr>
<tr>
<td> Parámetros </td>
<td> N/A </td>
</tr>
<tr>
<td> Encabezados </td>
<td>N/A</td>
</tr>
<tr>
<td> Cuerpo </td>
<td> N/A </td>
</tr>
</table>

##### Respuestas

<table>
<tr>
<th> Código </th>
<th> Descripción </th>
<th> Cuerpo </th>
</tr>
<tbody>
<tr>
<td> 200 </td>
<td> Solo para confirmar que el servicio está arriba.</td>
<td>

```pong```
</td>
</tr>
</tbody>
</table>

#### 6. Restablecer base de datos

Usado para limpiar la base de datos del servicio.

<table>
<tr>
<td> Método </td>
<td> POST </td>
</tr>
<tr>
<td> Ruta </td>
<td> <strong>/routes/reset</strong> </td>
</tr>
<tr>
<td> Parámetros </td>
<td> N/A </td>
</tr>
<tr>
<td> Encabezados </td>
<td>N/A</td>
</tr>
<tr>
<td> Cuerpo </td>
<td> N/A </td>
</tr>
</table>

##### Respuestas

<table>
<tr>
<th> Código </th>
<th> Descripción </th>
<th> Cuerpo </th>
</tr>
<tbody>
<tr>
<td> 200 </td>
<td> Todos los datos fueron eliminados.</td>
<td>

```
{"msg": "Todos los datos fueron eliminados"}
```
</td>
</tr>
</tbody>
</table>

## Pruebas

Las pruebas se ejecutan corriendo el archivo `test.sh`.

Las pruebas se hicieron usando la librería [sqlMock para golang](https://pkg.go.dev/github.com/DATA-DOG/go-sqlmock@v1.5.2), para mockear los llamados a base datos, y se utilizó el patrón de inyección de dependencias para simular los llamados al autenticador (que deberían ser al servicio de users). El sentido de esto es para limitar la cantidad de pruebas que se hacen a código que no es el del servicio, sino las conexiones con terceros. La inyección de dependencias se hizo creando mocks en [este archivo](/route_management/mocks/mock_authenticator.go).

## Otras Características

Requerimientos

- Golang 1.23.0
- Docker
- Postgresql

## Autor

Andres Felipe Losada Luna

