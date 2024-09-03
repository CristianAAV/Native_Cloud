# commands/getpost.py
import uuid
from .base_command import BaseCommannd
from ..models.post import Post, PostJsonSchema
from ..database import Session
from ..errors.errors import InvalidPost,InvalidId

class GetPost(BaseCommannd):
    def __init__(self, id, ):
        self.id = id
        
    def execute(self):
        # Inicializa la sesión de la base de datos
        session = Session()
        query = session.query(Post)
        try:
            post_id = uuid.UUID(self.id)
            post_id_str=str(post_id)   
        except ValueError:
            session.close()
            raise InvalidId()
        query = query.filter(Post.id == post_id_str)

        if len(query.all()) <= 0:
            session.close()
            raise InvalidPost()

        query = session.query(Post).filter_by(id=post_id_str).one()
        schema = PostJsonSchema()
        post = schema.dump(query)

        session.close()

        return post