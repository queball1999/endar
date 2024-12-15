FROM python:3.8-slim-buster

WORKDIR /app

RUN apt-get update \
    && apt-get -y install libpq-dev gcc git\
    && pip3 install psycopg2

COPY requirements.txt requirements.txt
RUN pip3 install -r requirements.txt

COPY . .

CMD ["/bin/bash","run.sh"]
#CMD ["gunicorn", "--bind", "0.0.0.0:5000", "flask_app:app"]
