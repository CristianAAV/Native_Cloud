from httmock import HTTMock, all_requests, response
from tests.utils.constants import STATIC_FAKE_UUID

@all_requests
def valide_auth(url, request):
  return response(200, { 'id': STATIC_FAKE_UUID }, {}, None, 5, request)

@all_requests
def failed_auth(url, request):
  return { 'status_code': 401 }

@all_requests
def forbidden_auth(url, request):
  return { 'status_code': 403 }