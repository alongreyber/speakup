FROM python:latest

RUN mkdir code
WORKDIR /code

# Packages necessary for build but that should be removed
#ENV build_packages
#RUN apt update && apt install $build_packages

ENV run_packages openssl
RUN apt update && apt install $run_packages

#RUN apt remove $build_packages

ADD . /code
RUN pip install -r requirements.txt

# Default to our production web server
CMD gunicorn --bind 0.0.0.0:80 main:app
