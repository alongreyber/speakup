from app import app
from flask import render_template, request
import app.forms as Forms

# Temporary
import socket

@app.route('/submit_alignment', methods=['POST'])
def submit_alignment():
    alignment_form = Forms.Alignment(request.form)
    if(alignment_form.validate()):
        return 'Form validation success!'
    else:
        return 'Form validation failed'


@app.route('/alignment', methods=['GET'])
def alignment():
    alignment_form = Forms.Alignment(request.form)
    return render_template('upload_for_alignment.html',
            alignment_form=alignment_form)

@app.route('/ip')
def ip():
    ip_addr = socket.gethostbyname(socket.gethostname())
    return "Hello! My internal cluster IP is " + str(ip_addr)
