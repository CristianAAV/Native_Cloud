# Offer Management

Es el servicio de gestión de ofertas.

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