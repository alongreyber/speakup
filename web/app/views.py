from app import app
from flask import render_template
import app.forms as Forms

# Temporary
import socket

@app.route('/submit_gentle', methods=['GET', 'POST'])
def submit_gentle():
    alignment_form = Forms.Alignment()
    if alignment_form.validate_on_submit():
        return render_template('gentle_result.html')
    return render_template('submit_gentle.html',
            alignment_form=alignment_form)

@app.route('/angular')
def angular():
    return render_template('angular.html')

@app.route('/ip')
def ip():
    ip_addr = socket.gethostbyname(socket.gethostname())
    return "Hello! My IP is " + str(ip_addr)
