# commands/create.py
from .base_command import BaseCommannd
from ..models.post import Post, PostJsonSchema,PostSchema
from ..database import Session


class CreatePost(BaseCommannd):
    def __init__(self, data,user_id):
        self.data = data
        self.user_id = user_id
    def execute(self):
        userid = self.user_id
        route_id = self.data['routeId']
        expire_at = self.data['expireAt']        
        save_post = Post(
            userId=userid,
            routeId=route_id,
            expireAt=expire_at
        )
        session = Session()

        session.add(save_post)
        session.commit()

        new_post = PostJsonSchema().dump(save_post)
        session.close()
        return new_post

