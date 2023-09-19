from queries import BaseQuery



class ProductLocation(BaseQuery):
    
    def __init__(self, id):
        self.id = id
        self.api_path = "query/product/v1/location"

class ProductStatus(BaseQuery):

    def __init__(self, id):
        self.id = id
        self.api_path = "query/product/v1/status"


class ProductDetails(BaseQuery):

    def __init__(self, id):
        self.id = id
        self.api_path = "query/product/v1/details"

class ProductTrack(BaseQuery):

    def __init__(self, id):
        self.id = id
        self.api_path = "query/product/v1/track"