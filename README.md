# Grupo 5 - Entrega 1
Para la entrega 1, se desarrollaron los servicios solicitados segun el siguiente detalle:

- [Gestión de usuarios](#gestión-de-usuarios)
- [Gestión de trayectos](#gestión-de-trayectos)
- [Gestión de publicaciones](#gestión-de-publicaciones)
- [Gestión de ofertas](#gestión-de-ofertas)

# Gestión de usuarios

Este servicio corresponde a la creación de usuarios y su validación por medio de tokens.

## Índice

1. [Estructura](#estructura)
2. [Ejecución](#ejecución)
3. [Uso](#uso)
4. [Pruebas](#pruebas)
5. [Otras Caracteristicas](#otras-caracteristicas)
6. [Autor](#autor)


## Estructura
Segun lo recomendado se construyeron las siguientes carpetas
- **src:** LA cual contiene el codigo y logica del servicio
  - **src/models:** Esta carpeta contiene la capa de persistencia, donde se definen el modelo de la entidad usuario.
  - **src/commands:** Contiene el codigo para cada caso de uso, por ejemplo "crear usuario", "generar token, etc"
  - **src/blueprints:** Esta carpeta alberga la capa de aplicación de nuestro microservicio, encargada de definir y desarrollar cada servicio API que exponemos.
  - **src/errors:** Para devolver errores HTTP en los blueprints
- **tests:** En la cual estan implementados los tests unitarios
- **Dockerfile:** Define cómo se debe construir la imagen de la aplicación.
- **docker-compose.yml:** Encargado de montar la base de datos y la aplicación.
- **.env:** Almacena variables de entorno como configuraciones y credenciales de manera segura

## Ejecución
En la raiz del servicio utilizar como plantilla el archivo "env.template" y crear uno denominado "env.development", con la siguiente informacion:
* DB_USER=admin
* DB_PASSWORD=adminpassword
* DB_HOST=cnt_db_users
* DB_PORT=5432
* DB_NAME=user_db

Seguidamente, ejecutar el siguiente comando:
```bash
docker-compose up --build
```
## Uso
Despues de iniciarse los contenedores, puede ejecutar los endpoints de acuerdo a la siguiente detale:
El servicio de gestión de usuarios permite crear usuarios y validar la identidad de un usuario por medio de tokens.

  - **Creación de usuarios**
    - **Método:** POST  
    - **Ruta:** /users
    - **Parámetros:** N/A
    - **Encabezados:** N/A
    - **Cuerpo:**
```json
{
  "username": "nombre de usuario",
  "password": "contraseña del usuario",
  "email": "correo electrónico del usuario",
  "dni": "identificación",
  "fullName": "nombre completo del usuario",
  "phoneNumber": "número de teléfono"
}
```
  - **Actualización de usuarios**
    - **Método:** PATCH
    - **Ruta:** /users/{id}
    - **Parámetros:** id: identificador del usuario
    - **Encabezados:** N/A
    - **Cuerpo:**
```json
{
  "status": nuevo estado del usuario,
  "dni": identificación,
  "fullName": nombre completo del usuario,
  "phoneNumber": número de teléfono
}
```
  - **Generación de token**
    - **Método:** POST
    - **Ruta:** /users/auth
    - **Parámetros:** N/A
    - **Encabezados:** N/A
    - **Cuerpo:**
```json
{
  "username": nombre de usuario,
  "password": contraseña del usuario
}
```
  - **Consultar información del usuario**
    - **Método:** GET
    - **Ruta:** /users/name
    - **Parámetros:** N/A
    - **Encabezados:** Authorization: Bearer token
    - **Cuerpo:** N/A
  - **Consulta de salud del servicio**
    - **Método:** GET
    - **Ruta:** /users/ping
    - **Parámetros:** N/A
    - **Encabezados:** N/A
    - **Cuerpo:** N/A
  - **Restablecer base de datos**
    - **Método:** POST
    - **Ruta:** /users/reset
    - **Parámetros:** N/A
    - **Encabezados:** N/A
    - **Cuerpo:** N/A

## Pruebas
- Ejecutar la BD de pruebas
```
docker run --name postgres-pruebas-local \
  -e POSTGRES_USER=user_management_user \
  -e POSTGRES_PASSWORD=user_management_pass \
  -e POSTGRES_DB=user_management_db \
  -p 5432:5432 \
  -d postgres:13
```
- Dirigirse a la carpeta  "user_management"
- Activar el entorno virtual, puede utilizar los siguientes comandos en linux:
  - python3 -m venv venv
  - source venv/bin/activate
- Instalar las dependencias:
  - pip install -r requirements.txt
- Ejecutar las pruebas:
  - pytest --cov-fail-under=70 --cov=src

## Otras Caracteristicas
- El servicio fue desarrollado en  Python con Flask
- Para la persistencia en base de datos, se utilizó Postgresql

## Autor
Emerson Chaparro Ampa
- correo: e.chaparroa@uniandes.edu.co

# Gestión de trayectos

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


# Gestión de publicaciones

Este servicio corresponde a la creación de publicaciones de rutas, su consulta y su eliminación. 

## Índice

1. [Estructura](#estructura)
2. [Ejecución](#ejecución)
3. [Uso](#uso)
4. [Otras Caracteristicas](#Otras\_Caracteristicas)
5. [Autor](#autor)


## Estructura
Segun lo recomendado se construyeron las siguientes carpetas
- **src:** LA cual contiene el codigo y logica del servicio
  - **src/models:** Esta carpeta contiene la capa de persistencia, donde se definen el modelo de la entidad usuario.
  - **src/commands:** Contiene el codigo para cada caso de uso, por ejemplo "crear usuario", "generar token, etc"
  - **src/blueprints:** Esta carpeta alberga la capa de aplicación de nuestro microservicio, encargada de definir y desarrollar cada servicio API que exponemos.
  - **src/errors:** Para devolver errores HTTP en los blueprints
- **tests:** En la cual estan implementados los tests unitarios
- **Dockerfile:** Define cómo se debe construir la imagen de la aplicación.
- **docker-compose.yml:** Encargado de montar la base de datos y la aplicación.
- **.env:** Almacena variables de entorno como configuraciones y credenciales de manera segura

## Ejecución
En la raiz del servicio utilizar como plantilla el archivo "env.template" y crear uno denominado "env.development", con la siguiente informacion:
VERSION=1.0
* DB_USER=postgres
* DB_PASSWORD=postgres
* DB_HOST=post_db
* DB_PORT=5432
* DB_NAME=post_db
* USER_PATH="http://users:3000"

Seguidamente, ejecutar el siguiente comando:
```bash
docker-compose up --build
```
## Uso
Despues de iniciarse los contenedores, puede ejecutar los endpoints de acuerdo a la siguiente detale:

# Servicio de Gestión de Publicaciones

El servicio de gestión de publicaciones permite crear, buscar, eliminar y consultar publicaciones.

## 1. Creación de Publicación

Crea una publicación asociada al usuario al que pertenece el token.

- **Método:** `POST`
- **Ruta:** `/posts`
- **Parámetros:** N/A
- **Encabezados:**
  - `Authorization: Bearer token`
- **Cuerpo:**
  ```json
  {
    "routeId": "id del trayecto",
    "expireAt": "fecha y hora máxima en que se recibirán ofertas en formato ISO"
  }

## Pruebas
- Ejecutar la BD de pruebas
```
docker run --name postgres-pruebas-local \
  -e POSTGRES_USER=user_management_user \
  -e POSTGRES_PASSWORD=user_management_pass \
  -e POSTGRES_DB=user_management_db \
  -p 5432:5432 \
  -d postgres:13
```
- Dirigirse a la carpeta  "user_management"
- Activar el entorno virtual, puede utilizar los siguientes comandos en linux:
  - python3 -m venv venv
  - source venv/bin/activate
- Instalar las dependencias:
  - pip install -r requirements.txt
- Ejecutar las pruebas:
  - pytest --cov-fail-under=70 --cov=src

### Respuestas

| Código | Descripción                                                                | Cuerpo                                                                                                                                |
|--------|----------------------------------------------------------------------------|---------------------------------------------------------------------------------------------------------------------------------------|
| 401    | El token no es válido o está vencido.                                      | N/A                                                                                                                                   |
| 403    | No hay token en la solicitud.                                              | N/A                                                                                                                                   |
| 400    | Alguno de los campos no está presente en la solicitud o no tiene el formato esperado. | N/A                                                                                                                                   |
| 412    | La fecha de expiración no es en el futuro o no es válida.                  | `{ "msg": "La fecha expiración no es válida" }`                                                                                       |
| 201    | La publicación se ha creado con éxito.                                     | `{ "id": "id de la publicación", "userId": "id del usuario que creó la publicación", "createdAt": "fecha y hora de creación en formato ISO" }` |

### 2. Ver y filtrar publicaciones

Retorna el listado de publicaciones que coinciden con los parámetros brindados. Solo un usuario autorizado puede realizar esta operación.

- **Método:** GET
- **Ruta:** `/posts?expire={true|false}&route={routeId}&owner={id|me}`
- **Parámetros:** Todos los parámetros son opcionales, y su funcionamiento es de tipo AND.
  - `expire`: filtra por publicaciones expiradas o no expiradas.
  - `route`: id del trayecto que se desea usar para el envío.
  - `owner`: dueño de la publicación. Se reciben ids o el valor `me` que indica al usuario del token.

  En el caso de que ninguno esté presente se devolverá la lista de datos sin filtrar. Es decir, todas las publicaciones.

- **Encabezados:** 
  - `Authorization: Bearer token`

- **Cuerpo:** N/A

#### Respuestas

| Código | Descripción                                                                  | Cuerpo                                                                                                                              |
|--------|------------------------------------------------------------------------------|-------------------------------------------------------------------------------------------------------------------------------------|
| 401    | El token no es válido o está vencido.                                        | N/A                                                                                                                                 |
| 403    | No hay token en la solicitud.                                                | N/A                                                                                                                                 |
| 400    | Alguno de los campos de búsqueda no tiene el formato esperado.               | N/A                                                                                                                                 |
| 200    | Listado de publicaciones que corresponden a la búsqueda.                     | `[ { "id": "id de la publicación", "routeId": "id del trayecto", "userId": "id del usuario owner de la publicación", "expireAt": "fecha y hora máxima en que se recibirán ofertas en formato ISO", "createdAt": "fecha y hora de creación de la publicación en formato ISO" } ]` |

### 3. Consultar una publicación

Retorna una publicación, solo un usuario autorizado puede realizar esta operación.

- **Método:** GET
- **Ruta:** `/posts/{id}`
- **Parámetros:** 
  - `id`: identificador de la publicación

- **Encabezados:** 
  - `Authorization: Bearer token`

- **Cuerpo:** N/A

#### Respuestas

| Código | Descripción                                             | Cuerpo                                                                                                                            |
|--------|---------------------------------------------------------|-----------------------------------------------------------------------------------------------------------------------------------|
| 401    | El token no es válido o está vencido.                   | N/A                                                                                                                               |
| 403    | No hay token en la solicitud.                           | N/A                                                                                                                               |
| 400    | El id no es un valor string con formato UUID.           | N/A                                                                                                                               |
| 404    | La publicación con ese id no existe.                    | N/A                                                                                                                               |
| 200    | Publicación que corresponde al identificador.           | `{ "id": "id de la publicación", "routeId": "id del trayecto", "userId": "id del usuario owner de la publicación", "expireAt": "fecha y hora máxima en que se recibirán ofertas en formato ISO", "createdAt": "fecha y hora de creación de la publicación en formato ISO" }` |


### 4. Eliminar publicación

Elimina una publicación, solo un usuario autorizado puede realizar esta operación.

- **Método:** DELETE
- **Ruta:** `/posts/{id}`
- **Parámetros:** 
  - `id`: identificador de la publicación.

- **Encabezados:** 
  - `Authorization: Bearer token`

- **Cuerpo:** N/A

#### Respuestas

| Código | Descripción                                             | Cuerpo                              |
|--------|---------------------------------------------------------|-------------------------------------|
| 401    | El token no es válido o está vencido.                   | N/A                                 |
| 403    | No hay token en la solicitud.                           | N/A                                 |
| 400    | El id no es un valor string con formato UUID.           | N/A                                 |
| 404    | La publicación con ese id no existe.                    | N/A                                 |
| 200    | La publicación fue eliminada.                           | `{ "msg": "la publicación fue eliminada" }` |


### 5. Consulta de salud del servicio

Usado para verificar el estado del servicio.

- **Método:** GET
- **Ruta:** `/posts/ping`
- **Parámetros:** N/A
- **Encabezados:** N/A
- **Cuerpo:** N/A

#### Respuestas

| Código | Descripción                                    | Cuerpo |
|--------|------------------------------------------------|--------|
| 200    | Solo para confirmar que el servicio está arriba. | `pong` |

### 6. Restablecer base de datos

Usado para limpiar la base de datos del servicio.

- **Método:** POST
- **Ruta:** `/posts/reset`
- **Parámetros:** N/A
- **Encabezados:** N/A
- **Cuerpo:** N/A

#### Respuestas

| Código | Descripción                        | Cuerpo                                   |
|--------|------------------------------------|------------------------------------------|
| 200    | Todos los datos fueron eliminados. | `{ "msg": "Todos los datos fueron eliminados" }` |


## Otras Caracteristicas
- El servicio fue desarrollado en  Python con Flask
- Para la persistencia en base de datos, se utilizó Postgresql

## Autor
Cristian Arnulfo Arias Vargas
- correo: ca.ariasv1@uniandes.edu.co

# Gestión de ofertas
Este servicio corresponde a la creación de ofertas, su consulta y su eliminación. 
## Índice
1. Estructura
2. Ejecucion
3. Uso
4. Pruebas Postman 
5. Autor

## 1. Estructura

La estructura del proyecto es la siguiente:

- **auth/**: Realiza la validacion de verificacion del usuario.
- **globals/**: Son variables que pueden ser usadas muchas veces en el codigo.
- **db/**: Se realiza la conexión a la base de datos.
- **migration/**: Contiene consultas que no se pueden ejecutar directamente con GORM.
- **models/**: Aquí se encuentran las entidades o tablas. Estos modelos hacen uso de GORM.
- **routes/**: Contiene los manejadores de las rutas definidas en `main.go`. Aquí se almacena la lógica de negocio que se ejecuta en las rutas.
- **main.go**: Es el archivo principal de la aplicación.
- **go.mod**: Archivo que almacena los paquetes necesarios para que la aplicación funcione correctamente.
- **Dockerfile**: Define cómo se debe construir la imagen de la aplicación.

## 2. Ejecución

Este servicio es muy fácil de ejecutar. Simplemente ejecuta el siguiente comando en el root:

```bash
docker-compose up --build
```

## 3. Uso

1. **Realizar Solicitudes al Servidor**

   Una vez que los servicios estén en ejecución, puedes hacer solicitudes al servidor en la siguiente dirección:

    - **Base URL**: `http://localhost:3003`

   Aquí están los endpoints disponibles:

    - **Crear una oferta**
        - **Método**: POST
        - **Ruta**: `/offers`
        - **Descripción**: Crea una nueva oferta.
        - **Cuerpo de la solicitud**:
          ```json
          {
            "postId": "string",
            "userId": "string",
            "description": "string",
            "size": "LARGE|MEDIUM|SMALL",
            "fragile": true|false,
            "offer": 0.00
          }
          ```

    - **Obtener todas las ofertas**
        - **Método**: GET
        - **Ruta**: `/offers`
        - **Descripción**: Obtiene una lista de todas las ofertas.

    - **Obtener una oferta específica**
        - **Método**: GET
        - **Ruta**: `/offers/{id}`
        - **Descripción**: Obtiene los detalles de una oferta específica.

    - **Eliminar una oferta**
        - **Método**: DELETE
        - **Ruta**: `/offers/{id}`
        - **Descripción**: Elimina una oferta específica.

    - **Verificar la salud del servicio**
        - **Método**: GET
        - **Ruta**: `/offers/ping`
        - **Descripción**: Verifica que el servicio está en funcionamiento.

    - **Restablecer la base de datos**
        - **Método**: POST
        - **Ruta**: `/offers/reset`
        - **Descripción**: Restablece la base de datos a su estado inicial.

## 4. Pruebas Postman

## Pasos previos antes de ejecutar las pruebas

1. **Ejecutar Docker Compose:**
   Asegúrate de haber ejecutado el archivo `docker-compose` en el root para iniciar todos los servicios necesarios.

2. **Creación de usuario y autenticación:**
   - Visita la documentación del servicio de usuarios para el endpoint `http://localhost:3000/users` el cual permite crear un usuario.
   - Para obtener un token de autenticación, dirígete a la documentación del servicio de usuarios y revisa el siguiente endpoint `localhost:3000/users/auth`. Este token será necesario para realizar las solicitudes a otros endpoints.

## Endpoints de ofertas

### 1. Creación de una oferta

```bash
curl --location 'localhost:3003/offers' \
--header 'Content-Type: application/json' \
--header 'Authorization: Bearer TOKEN' \
--data '{
    "postId": "550e8400-e29b-41d4-a716-446655440000",
    "description": "Test Offer",
    "size": "LARGE",
    "fragile": true,
    "offer": 10.00
}'
```

### 2. Consultar todas las ofertas

```bash
curl --location 'localhost:3003/offers' \
--header 'Authorization: Bearer TOKEN'
```

### 3. Consultar una oferta por ID

```bash
curl --location 'localhost:3003/offers/ID_OFERTA' \
--header 'Authorization: Bearer TOKEN'
```

### 4. Eliminar una oferta

```bash
curl --location --request DELETE 'localhost:3003/offers/ID_OFERTA' \
--header 'Authorization: Bearer TOKEN'
```

### 5. Consultar el estado del servicio (Health Check)

```bash
curl --location 'localhost:3003/offers/ping'
```

### 6. Eliminar todos los datos de la base de datos

```bash
curl --location --request POST 'localhost:3003/offers/reset'
```

---

Asegúrate de reemplazar `TOKEN` con el token real obtenido durante el proceso de autenticación y `ID_OFERTA` con el ID de la oferta que deseas consultar o eliminar.

## 5. Autor

Juan Camilo Vallejos Guerrero  
Correo electrónico: j.vallejosg@uniandes.edu.co
