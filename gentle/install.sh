#!/bin/bash

set -e

git clone https://github.com/kaldi-asr/kaldi ext/kaldi
./install_deps.sh
(cd ext && ./install_kaldi.sh)
./install_models.sh
