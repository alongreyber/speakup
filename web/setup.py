from setuptools import setup

setup(
    name='app',
    packages=['app'],
    include_package_data=True,

    sass_manifests={
        'app': ('static/sass', 'static/css', '/static/css')
    },

    install_requires=[
        'flask',
    ],
)

