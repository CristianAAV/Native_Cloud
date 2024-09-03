import requests
import os
from .base_command import BaseCommannd
from ..errors.errors import InvalidToken, FaultToken

class ValidateToken(BaseCommannd):
    def __init__(self, token):
        self.token = token

    def execute(self):
        if not self.token:
            raise FaultToken()
        user_url = os.environ.get("USERS_PATH")         
        headers = {'Authorization': self.token}
        # Verifica si el valor del token es None, vacío, o si la clave no está presente
        try:
            response = requests.get(f'{user_url}/users/me', headers=headers)
            print(response)
            response.raise_for_status()
            if response.status_code == 401:
                raise InvalidToken()
            if response.status_code == 403:
                raise FaultToken()
            data = response.json()
            user_id = data.get('id')
            if not user_id:
                raise InvalidToken()
            return user_id

        except requests.HTTPError as http_err:
            if response.status_code == 401:
                raise InvalidToken()
            raise http_err
        except Exception as err:
            return f"Error inesperado: {err}"
