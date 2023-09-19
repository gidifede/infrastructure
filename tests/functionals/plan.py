import json

from collections.abc import Sequence
from commands import BaseCommand
from queries import BaseQuery

class Plan():

    # Plan test definition
    def __init__(self, id, commands: Sequence[BaseCommand], queries: Sequence[BaseQuery], base_api: str, api_key:str = None):
        self.id = id
        self.commands = commands
        self.queries = queries
        self.base_api = base_api
        if api_key is not None:
            self.api_key = api_key
        super().__init__()

    # Commands send
    def start_command(self):
        responses = []
        for j in self.commands:
            responses.append(j.send(base_api_url=self.base_api, api_key=(self.api_key if hasattr(self, 'api_key') else None)))
        return responses

    # QueryHandler calling
    def start_query(self):
        param = {'id': self.id}
        results = []
        for query in self.queries:
            results.append(query.send_query(base_api_url=self.base_api, query_param = f'id={self.id}', api_key=(self.api_key if hasattr(self, 'api_key') else None)))
        return results