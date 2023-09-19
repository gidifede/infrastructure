pip install --target ./package mysql-connector-python
cd package
zip -r ../lambda.zip .
cd ..
zip lambda.zip handler.py
rm -rf package