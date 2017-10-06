from app import app
from flask import render_template

@app.route('/')
def hello():
    return render_template('submit_gentle.html')

@app.route('/angular')
def angular():
    return render_template('angular.html')
