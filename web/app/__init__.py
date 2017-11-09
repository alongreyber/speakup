from flask import Flask, render_template
from redis import Redis

app = Flask(__name__)
app.config.from_object('config')

redis = Redis(host='redis', port=6379)

from app import views
