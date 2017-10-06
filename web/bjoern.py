import bjoern
import os
import signal
from app import app

host = '0.0.0.0'
port = 5000
NUM_WORKERS = 2
worker_pids = []


bjoern.listen(app, host, port)
