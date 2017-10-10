from app import app
from OpenSSL import SSL
context = SSL.Context(SSL.SSLv23_METHOD)
context.use_privatekey_file('app/self-signed-tmp.key')
context.use_certificate_file('app/self-signed-tmp.crt')

if __name__ == '__main__':
    app.run(host='0.0.0.0', debug=True,
            ssl_context='ad-hoc')
