#./blueprints/post.py
from flask import Flask, jsonify, request, Blueprint
from sqlalchemy.orm import sessionmaker
from sqlalchemy import create_engine
import uuid
from datetime import datetime, timezone
from ..models.post import PostJsonSchema,Post
from ..commands.token import ValidateToken
from ..commands.fields import ValidateFields
from ..commands.create import CreatePost
from ..commands.getposts import GetPosts
from ..commands.getpost import GetPost
from ..commands.deletepost import DeletePost
from ..models.model import Base
from ..database import Session,engine
from ..errors.errors  import InvalidToken, FaultToken, Invalid, MissingField,InvalidId, InvalidPost

import os


operations_blueprint = Blueprint('operations', __name__)

@operations_blueprint.route('/posts', methods=['POST'])
def posts():
    auth_header = request.headers.get('Authorization')
    user_id = ValidateToken(auth_header).execute()
    data = request.get_json()
    dataPost = ValidateFields(data, user_id).execute()
    new_post = CreatePost(dataPost, user_id).execute()
    return jsonify({
        "userId": new_post['userId'],  
        "createdAt": new_post['createdAt'], 
        "id": new_post['id']
    }), 201

    

@operations_blueprint.route('/posts', methods=['GET'])
def get_posts():
    auth_header = request.headers.get('Authorization')
    user_id = ValidateToken(auth_header).execute()
    data=request.args.to_dict()
    result = GetPosts(data,user_id).execute()
    return jsonify(result), 200


@operations_blueprint.route('/posts/<id>', methods=['GET'])
def getpost(id):
    auth_header = request.headers.get('Authorization')
     # Verificar si el encabezado está presente
    ValidateToken(auth_header).execute()   
    post=GetPost(id).execute()
    return jsonify(post), 200

@operations_blueprint.route('/posts/<id>', methods=['DELETE'])
def deletepost(id):
    auth_header = request.headers.get('Authorization')
     # Verificar si el encabezado está presente
    ValidateToken(auth_header).execute()   
    DeletePost(id).execute()
    return jsonify({"msg": "la publicación fue eliminada"}), 200

@operations_blueprint.route('/posts/ping', methods=['GET'])
def ping():
    return ("pong"), 200


@operations_blueprint.route('/posts/reset', methods=['POST'])
def reset():
    session = Session()
    try:
        # Eliminar todos los registros de la tabla Post
        session.query(Post).delete()
        session.commit()
        return jsonify({'msg': "Todos los datos fueron eliminados"}), 200    
    except Exception as e:
        session.rollback()
        return jsonify({'msg': 'Error al eliminar los datos', 'error': str(e)}), 500
    finally:
        session.close()