from commands import BaseCommand, generate_cloudevent


class Create(BaseCommand):

    def __init__(self, cluster=None,location=None, timestamp=None):
        self.payload = {}
        if cluster is not None:
            self.payload['cluster'] = cluster.to_json()
        if location is not None:
            self.payload['location'] = location.to_json()
        if timestamp is not None:
            self.payload['timestamp'] = timestamp
        self.api_path='cluster/v1/create'
        
        self.validate()

        self.headers, self.body = generate_cloudevent(
            source='Logistic.PCL.UP.OMP',
            type='Logistic.PCL.Cluster.Create.Create',
            data=self.payload
        )

class AddProduct(BaseCommand):

    def __init__(self, product=None, timestamp=None, location=None, cluster=None, mode=None):
        self.payload = {}
        if cluster is not None:
            self.payload['cluster'] = cluster.to_json()
        if location is not None:
            self.payload['location'] = location.to_json()
        if mode is not None:
            self.payload['mode'] = mode
        if timestamp is not None:
            self.payload['timestamp'] = timestamp
        if product is not None:
            self.payload['product'] = product.to_json()
        
        self.api_path='cluster/v1/add_product'

        self.validate()
        self.headers, self.body = generate_cloudevent(
            source='Logistic.PCL.UP.OMP',
            type='Logistic.PCL.Cluster.AddProduct.AddProduct',
            data=self.payload
        )

class Close(BaseCommand):

    def __init__(self, cluster=None, timestamp=None, location=None):
        self.payload = {}
        if location is not None:
            self.payload['location'] = location.to_json()
        if timestamp is not None:
            self.payload['timestamp'] = timestamp
        if cluster is not None:
            self.payload['cluster'] = cluster.to_json()
        
        self.api_path='cluster/v1/close'

        self.validate()
        self.headers, self.body = generate_cloudevent(
            source='Logistic.PCL.UP.OMP',
            type='Logistic.PCL.Cluster.Close.Close',
            data=self.payload
        )

class EndTransport(BaseCommand):

    def __init__(self, cluster=None,timestamp=None, location=None, transport=None):
        self.payload = {}
        if cluster is not None:
            self.payload['cluster'] = cluster.to_json()
        if location is not None:
            self.payload['location'] = location.to_json()
        if timestamp is not None:
            self.payload['timestamp'] = timestamp
        if transport is not None:
            self.payload['transport'] = transport.to_json()
        
        self.api_path='cluster/v1/end_transport'

        self.validate()
        self.headers, self.body = generate_cloudevent(
            source='Logistic.PCL.UP.OMP',
            type='Logistic.PCL.Cluster.EndTrasport.EndTransport',
            data=self.payload
        )

class Open(BaseCommand):

        def __init__(self, cluster=None, timestamp=None, location=None):
            self.payload = {}
            if timestamp is not None:
                self.payload['timestamp'] = timestamp
            if cluster is not None:
                self.payload['cluster'] = cluster.to_json()
            if location is not None:
                self.payload['location'] = location.to_json()

            self.api_path='cluster/v1/open'

            self.validate()
            self.headers, self.body = generate_cloudevent(
                source='Logistic.PCL.UP.OMP',
                type='Logistic.PCL.Cluster.Open.Open',
                data=self.payload
        )
            
class StartTransport(BaseCommand):

        def __init__(self, transport=None, timestamp=None, cluster=None, location=None):
            self.payload = {}
            if timestamp is not None:
                self.payload['timestamp'] = timestamp
            if cluster is not None:
                self.payload['cluster'] = cluster.to_json()
            if transport is not None:
                self.payload['transport'] = transport.to_json()
            if location is not None:
                self.payload['location'] = location.to_json()

            self.api_path='cluster/v1/start_transport'

            self.validate()
            self.headers, self.body = generate_cloudevent(
                source='Logistic.PCL.UP.OMP',
                type='Logistic.PCL.Cluster.StartTransport.StartTransport',
                data=self.payload
        )

