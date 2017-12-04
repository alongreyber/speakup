#!/bin/bash

set -e

# Specific branch of kaldi that gentle likes
(cd ext && git clone --recursive https://github.com/kaldi-asr/kaldi &&
    cd kaldi && git checkout 61510ca0d8d9b38096701227ee064a166816fcbe)
./install_deps.sh
(cd ext && ./install_kaldi.sh)

./install_models.sh
cd ext && make -j5 depend && make -j5
