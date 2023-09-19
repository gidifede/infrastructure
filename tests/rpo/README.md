to start capturing logs, insert following log format in the application:

```
"REQUEST_JSON_ID|<timestamp>|<request_id>"
```

and enable `enable-log-inspection` in cdk.json

```
# python tests/rpo/load.py                 # sends 100 requests in sequence (sends a request after receiving a response from the previous one)
```

```
# python tests/rpo/load.py 1000 parallel   # sends 1000 requests in parallel (using threads)
```

to start a notebook:

```
pip install jupyter
```
```
jupyter notebook
```