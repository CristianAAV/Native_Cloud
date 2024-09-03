from sqlalchemy.orm import sessionmaker
from sqlalchemy import create_engine
import os


class DataBase():
    def __init__(self):
        ...

    def url(self):
        db_user = os.environ['DB_USER']
        db_pass = os.environ['DB_PASSWORD']
        db_host = os.environ['DB_HOST']
        db_port = os.environ['DB_PORT']
        db_name = os.environ['DB_NAME']
        return f'postgresql://{db_user}:{db_pass}@{db_host}:{db_port}/{db_name}'

data_base = DataBase()
engine = create_engine(data_base.url())
Session = sessionmaker(bind=engine)
