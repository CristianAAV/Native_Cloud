# commands/fields.py
from datetime import datetime, timezone
import requests
from dateutil import parser
import os
from .base_command import BaseCommannd
from ..errors.errors import MissingField, dateInvalid, datePast, ApiError
class ValidateFields(BaseCommannd):
    def __init__(self, data, user_id):
        self.data = data
        self.user_id = user_id
    def execute(self):
        try:
            route_id = self.data['routeId']
            expire_at_str = self.data['expireAt']
            if not route_id or not expire_at_str:
                raise MissingField()
        except KeyError:
            raise MissingField()
        try:
            # Intentar varios formatos de fecha, incluido el formato completo con microsegundos
            if 'expireAt' in self.data and expire_at_str != None and not self.valid_date():
                raise datePast()
        except ValueError:
            # Si la fecha no coincide con ningún formato esperado, lanzar una excepción
            raise dateInvalid()
        dataPost= {
            "userId": self.user_id,
            "routeId": route_id,
            "expireAt": parser.parse(expire_at_str).date().isoformat()
        }
        return dataPost
    def valid_date(self):
        try:
            # If last character is a Z remove it
            expireAt = self.data['expireAt'][:-1] if self.data['expireAt'][-1] == 'Z' else self.data['expireAt']
            date_obj = parser.parse(expireAt).date()
            current_utc_datetime = datetime.utcnow().date()
            return date_obj > current_utc_datetime
        except:
            return False