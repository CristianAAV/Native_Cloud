# commands/getposts.py
import uuid
from datetime import datetime, timezone
from .base_command import BaseCommannd
from ..models.post import Post, PostJsonSchema
from ..database import Session
from ..errors.errors import Invalid


class GetPosts(BaseCommannd):
    def __init__(self, data, user_id):
        self.data = data
        self.user_id = user_id

    def execute(self):

        # Obtiene los parámetros de la URL
        expire_at_param = self.data.get('expire')
        route_id = self.data.get('route')
        owner_param = self.data.get('owner')

        # Inicializa la sesión de la base de datos
        session = Session()
        query = session.query(Post)

        # Filtra por dueño (owner)
        if owner_param:
            if owner_param == 'me':
                owner_idstr = self.user_id
            else:
                try:
                    owner_id = uuid.UUID(owner_param)
                    owner_idstr=str(owner_id)
                except ValueError:
                    session.close()
                    raise Invalid()

            query = query.filter(Post.userId == owner_idstr)

        # Filtra por ID de ruta (route)
        if route_id:
            try:
                route_uuid = uuid.UUID(route_id)
                route_id = str(route_uuid)
                query = query.filter(Post.routeId == route_id)
            except ValueError:
                session.close()
                raise Invalid()
            
        # Filtra por expiración (expire)
        if expire_at_param:
            if expire_at_param.lower() == 'true':
                query = query.filter(Post.expireAt < datetime.now(timezone.utc))
            elif expire_at_param.lower() == 'false':
                query = query.filter(Post.expireAt >= datetime.now(timezone.utc))
            else:
                session.close()
                raise Invalid()

        # Ejecuta la consulta
        posts = query.all()
        session.close()

        # Serializa los resultados
        result = PostJsonSchema(many=True).dump(posts)
        return result