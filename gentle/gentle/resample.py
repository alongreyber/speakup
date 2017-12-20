import os
import subprocess
import tempfile
import wave
from shutil import copyfile

from contextlib import contextmanager


from util.paths import get_binary

FFMPEG = get_binary("ffmpeg")

def resample(infile, outfile):
    if not os.path.isfile(infile):
        raise IOError("Not a file: %s" % infile)

    '''
    Use FFMPEG to convert a media file to a wav file sampled at 8K
    '''
    # Check if the file already meets our specifications
    try:
        wave_file = wave.open(infile, 'rb')
        if wave_file.getframerate() == 8000:
            copyfile(infile,outfile)
            return 0
    except wave.Error as err:
        pass
    return subprocess.call([FFMPEG,
                            '-loglevel', 'panic',
                            '-y',
                            '-i', infile,
                            '-ac', '1', '-ar', '8000',
                            '-acodec', 'pcm_s16le',
                            outfile])

@contextmanager
def resampled(infile):
    with tempfile.NamedTemporaryFile(suffix='.wav') as fp:
        if resample(infile, fp.name) != 0:
            raise RuntimeError("Unable to resample/encode '%s'" % infile)
        yield fp.name
