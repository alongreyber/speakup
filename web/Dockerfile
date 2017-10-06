FROM python:3.4-alpine
ADD . /code
WORKDIR /code

# Packages necessary for build but that should be removed
#ENV build_packages gcc build-base libev-dev
#RUN apk update && apk add $build_packages
RUN pip install -r requirements.txt
#RUN apk del $build_packages
CMD gunicorn --bind 0.0.0.0:5000 main:app
