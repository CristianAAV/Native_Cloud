# src/main.py
from dotenv import load_dotenv, find_dotenv

env_file = find_dotenv('../.env.development')
loaded=load_dotenv(env_file)

from flask import Flask, jsonify
from .blueprints.posts import operations_blueprint
from .models.model import Base
from .database import engine
from .errors.errors import ApiError
import os


app = Flask(__name__)
app.config['ENV'] = 'development'
app.config['DEBUG'] = True

app.register_blueprint(operations_blueprint)

Base.metadata.create_all(engine)
if not loaded:
    print("Error: No se pudieron cargar las variables de entorno")


@app.errorhandler(ApiError)
def handle_exception(err):
    response = {
      "msg": err.description
    }
    return jsonify(response), err.code


# print("Variables de Entorno:")
# print(f"VERSION: {os.getenv('VERSION')}")
# print(f"DB_USER: {os.getenv('DB_USER')}")
# print(f"DB_PASSWORD: {os.getenv('DB_PASSWORD')}")
# print(f"DB_HOST: {os.getenv('DB_HOST')}")
# print(f"DB_PORT: {os.getenv('DB_PORT')}")
# print(f"DB_NAME: {os.getenv('DB_NAME')}")
# print(f"USER_PATH: {os.getenv('USER_PATH')}")
# Verifica si las variables se cargaron correctamente
print(f"VERSION: {os.getenv('VERSION')}")