import requests

class BaseQuery:
    
    api_path = None

    def send_query(self, base_api_url, query_param, api_key=None):

        headers = {
                'Content-type': 'application/json'
            }
        if api_key is not None:
            headers['x-api-key'] = api_key 

        response = requests.get(
            f"{base_api_url}{self.api_path}", params=query_param,
            headers=headers
        )
        if response.status_code not in [400, 500]:
            assert True
        else:
            print(f"Error in Query: {self.api_path}. Query params: {query_param}. Result: {response.status_code}, {response.content}")
            assert  False
        return response.json()