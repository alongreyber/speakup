from flask_wtf import Form
from wtforms import StringField, BooleanField, FileField
from wtforms.validators import DataRequired

class Alignment(Form):
    audio = FileField('Audio File', validators=[DataRequired()])
    transcript = FileField('Transcript Text File', validators=[DataRequired()])
