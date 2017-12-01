from flask_wtf import FlaskForm
from wtforms import StringField, BooleanField, FileField, SubmitField
from wtforms.validators import DataRequired

class Alignment(FlaskForm):
    audio = FileField('Audio File', validators=[DataRequired()])
    transcript = FileField('Transcript Text File', validators=[DataRequired()])
    send = SubmitField('Send')
