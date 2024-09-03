# Gestión de Usuarios

Este servicio corresponde a la creación de usuarios y su validación por medio de tokens.

## Índice

1. [Estructura](#estructura)
2. [Ejecución](#ejecución)
3. [Uso](#uso)
4. [Pruebas](#pruebas)
5. [Otras Caracteristicas](#Otras\_Caracteristicas)
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
