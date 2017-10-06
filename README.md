# SpeakUp
## It's time for your business to SpeakUp 
### Running the code
Install docker.
Run `docker-compose -f docker-compose.yml -f docker-compose.dev.yml up`
This will set up a development server with the web server exposed on port 5000 and gentle exposed on port 8765. The code is mounted and will live update on changes.
Alternatively, run `docker-compose -f docker-compose.yml -f docker-compose.prod.yml up`
This will start the same programs but will also start up an NGINX reverse proxy. The web server will then be accessible on port 80.
