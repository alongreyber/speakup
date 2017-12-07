# To be run one directory up
docker build -t sass ./sass
for f in $(pwd)/web/app/static/sass/*
do
    # Mount the current directory in a correspondingly named folder in the container
    docker run -v $(pwd):$(pwd) -w $(pwd) sass $f $(pwd)/web/app/static/css/$(basename $f).css
done
