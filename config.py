import os
basedir = os.path.abspath(os.path.dirname(__file__))

class Config:
    UI_HOST = os.environ.get("UI_HOST")
    APP_NAME = os.environ.get("APP_NAME","App")
    APP_SUBTITLE = os.environ.get("APP_SUBTITLE","")
    CR_YEAR = os.environ.get("CR_YEAR","2022")
    VERSION = os.environ.get("VERSION","1.0.0")

    LOG_TYPE = os.environ.get("LOG_TYPE", "stream")
    LOG_LEVEL = os.environ.get("LOG_LEVEL", "WARNING")

    SECRET_KEY = os.environ.get('SECRET_KEY', 'change_secret_key')
    SQLALCHEMY_COMMIT_ON_TEARDOWN = True
    SQLALCHEMY_TRACK_MODIFICATIONS = False
    SQLALCHEMY_RECORD_QUERIES = True
    MAIL_SERVER = 'smtp.googlemail.com'
    MAIL_PORT = 587
    MAIL_USE_TLS = True
    MAIL_DEBUG = os.environ.get('MAIL_DEBUG',False)
    MAIL_USERNAME = os.environ.get('MAIL_USERNAME')
    MAIL_PASSWORD = os.environ.get('MAIL_PASSWORD')
    BASE_DIR = basedir
    ENABLE_SELF_REGISTRATION = os.environ.get("ENABLE_SELF_REGISTRATION",False)
    ENABLE_GOOGLE_AUTH = os.environ.get("ENABLE_GOOGLE_AUTH","0")
    DOC_LINK = os.environ.get("DOC_LINK","/")

    DEFAULT_EMAIL = os.environ.get("DEFAULT_EMAIL", "admin@example.com")
    DEFAULT_PASSWORD = os.environ.get("DEFAULT_PASSWORD", "admin")

    DEFAULT_TENANT_LABEL = "Default Tenant"
    DEFAULT_GROUP_LABEL = "Default Group"

    OAUTHLIB_RELAX_TOKEN_SCOPE = os.environ.get("OAUTHLIB_RELAX_TOKEN_SCOPE","1")
    os.environ['OAUTHLIB_RELAX_TOKEN_SCOPE'] = OAUTHLIB_RELAX_TOKEN_SCOPE
    GOOGLE_OAUTH_CLIENT_SECRET = os.environ.get("GOOGLE_OAUTH_CLIENT_SECRET")
    GOOGLE_OAUTH_CLIENT_ID = os.environ.get("GOOGLE_OAUTH_CLIENT_ID")

    @staticmethod
    def init_app(app):
        pass

class DevelopmentConfig(Config):
    DEBUG = True
    SQLALCHEMY_DATABASE_URI = os.environ.get('SQLALCHEMY_DATABASE_URI') or \
        "postgresql://db1:db1@postgres_db/db1"

class TestingConfig(Config):
    TESTING = True
    SQLALCHEMY_DATABASE_URI = os.environ.get('SQLALCHEMY_DATABASE_URI') or \
        "postgresql://db1:db1@postgres_db/db1"
    WTF_CSRF_ENABLED = False

config = {
    'development': DevelopmentConfig,
    'testing': TestingConfig,
    'default': DevelopmentConfig
}
