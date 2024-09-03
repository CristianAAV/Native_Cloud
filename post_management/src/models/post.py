# ./models/post.py
import uuid
from datetime import datetime, timezone
from marshmallow import Schema, fields
from sqlalchemy import Column, String, DateTime
from .model import Model, Base

# Extender la clase Model proporcionada
class Post(Model, Base):
    __tablename__ = 'posts'    

    # Definición de columnas
    id = Column(String(36), primary_key=True, default=lambda: str(uuid.uuid4()))
    routeId = Column(String(36), nullable=False)
    userId = Column(String(36), nullable=False)
    expireAt = Column(DateTime, nullable=False)
    createdAt = Column(DateTime, nullable=False, default=lambda: datetime.now(timezone.utc).strftime('%Y-%m-%dT%H:%M:%S'))

    # Constructor
    def __init__(self, routeId, userId, expireAt):
        Model.__init__(self)
        self.routeId = routeId
        self.userId = userId
        self.expireAt = expireAt

# Especificar los campos que estarán presentes al serializar el objeto como JSON
class PostJsonSchema(Schema):
    id = fields.Str()
    routeId = fields.Str()
    userId = fields.Str()
    expireAt = fields.DateTime()
    createdAt = fields.DateTime()

class PostSchema(Schema):
    userId = fields.Str()
    routeId = fields.Str()
    expireAt = fields.DateTime()
