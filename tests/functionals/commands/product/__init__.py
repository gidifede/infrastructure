from commands import BaseCommand, generate_cloudevent


class AcceptProduct(BaseCommand):

    def __init__(self, product=None, timestamp=None, attributes=None, receiver=None, sender=None, location=None):
        self.payload = {}
        if receiver is not None:
            self.payload['receiver'] = receiver.to_json()
        if sender is not None:
            self.payload['sender'] = sender.to_json()
        if location is not None:
            self.payload['location'] = location.to_json()
        if timestamp is not None:
            self.payload['timestamp'] = timestamp
        if product is not None:
            self.payload['product'] = product.to_json()
        if attributes is not None:
            self.payload['attributes'] = attributes
        self.api_path='product/v1/accept'
        
        self.validate()

        self.headers, self.body = generate_cloudevent(
            source='Logistic.PCL.UP.OMP',
            type='Logistic.PCL.Product.Accept.Accept',
            data=self.payload
        )

class StartProductProcessing(BaseCommand):

    def __init__(self, product=None, timestamp=None, location=None):
        self.payload = {}
        if location is not None:
            self.payload['location'] = location.to_json()
        if timestamp is not None:
            self.payload['timestamp'] = timestamp
        if product is not None:
            self.payload['product'] = product.to_json()
        
        self.api_path='product/v1/start_processing'

        self.validate()
        self.headers, self.body = generate_cloudevent(
            source='Logistic.PCL.UP.OMP',
            type='Logistic.PCL.Product.Processing.StartProcessing',
            data=self.payload
        )

class FailProductProcessing(BaseCommand):

    def __init__(self, product=None, timestamp=None, location=None, notes=None, attributes=None):
        self.payload = {}
        if location is not None:
            self.payload['location'] = location.to_json()
        if timestamp is not None:
            self.payload['timestamp'] = timestamp
        if product is not None:
            self.payload['product'] = product.to_json()
        if attributes is not None:
            self.payload['attributes'] = attributes
        if notes is not None:
            self.payload['notes'] = notes
        
        self.api_path='product/v1/fail_processing'

        self.validate()
        self.headers, self.body = generate_cloudevent(
            source='Logistic.PCL.UP.OMP',
            type='Logistic.PCL.Product.Error.FailProcessing',
            data=self.payload
        )

class MakeProductReadyToBeDelivered(BaseCommand):

        def __init__(self, product=None, timestamp=None, attributes=None):
            self.payload = {}
            if timestamp is not None:
                self.payload['timestamp'] = timestamp
            if product is not None:
                self.payload['product'] = product.to_json()
            if attributes is not None:
                self.payload['attributes'] = attributes

            self.api_path='product/v1/make_ready_to_be_delivered'

            self.validate()
            self.headers, self.body = generate_cloudevent(
                source='Logistic.PCL.UP.OMP',
                type='Logistic.PCL.Product.Delivery.MakeReadyToBeDelivered',
                data=self.payload
        )
            
class MakeProductWaitingToBeWithdrawn(BaseCommand):

        def __init__(self, product=None, timestamp=None, location=None):
            self.payload = {}
            if timestamp is not None:
                self.payload['timestamp'] = timestamp
            if product is not None:
                self.payload['product'] = product.to_json()
            if location is not None:
                self.payload['location'] = location.to_json()

            self.api_path='product/v1/make_waiting_to_be_withdrawn'

            self.validate()
            self.headers, self.body = generate_cloudevent(
                source='Logistic.PCL.UP.OMP',
                type='Logistic.PCL.Product.Delivery.MakeWaitingToBeWithdrawn',
                data=self.payload
        )

class StartProductDelivery(BaseCommand):

        def __init__(self, product=None, timestamp=None, carrier=None, attributes=None):
            self.payload = {}
            if carrier is not None:
                self.payload['carrier'] = carrier.to_json()
            if timestamp is not None:
                self.payload['timestamp'] = timestamp
            if product is not None:
                self.payload['product'] = product.to_json()
            if attributes is not None:
                self.payload['attributes'] = attributes
        
            self.api_path='product/v1/start_delivery'

            self.validate()
            self.headers, self.body = generate_cloudevent(
                source='Logistic.PCL.UP.OMP',
                type='Logistic.PCL.Product.Delivery.StartDelivery',
                data=self.payload
        )

class CompleteProductDelivery(BaseCommand):

        def __init__(self, product=None, timestamp=None, carrier=None, attributes=None):
            self.payload = {}
            if carrier is not None:
                self.payload['carrier'] = carrier.to_json()
            if timestamp is not None:
                self.payload['timestamp'] = timestamp
            if product is not None:
                self.payload['product'] = product.to_json()
            if attributes is not None:
                self.payload['attributes'] = attributes        
            
            self.api_path='product/v1/complete_delivery'

            self.validate()
            self.headers, self.body = generate_cloudevent(
                source='Logistic.PCL.UP.OMP',
                type='Logistic.PCL.Product.Delivery.CompleteDelivery',
                data=self.payload
        )

class FailProductDelivery(BaseCommand):

        def __init__(self, product=None, timestamp=None, carrier=None, attributes=None, notes=None):
            self.payload = {}
            if carrier is not None:
                self.payload['carrier'] = carrier.to_json()
            if timestamp is not None:
                self.payload['timestamp'] = timestamp
            if product is not None:
                self.payload['product'] = product.to_json()
            if attributes is not None:
                self.payload['attributes'] = attributes 
            if notes is not None:
                self.payload['notes'] = notes        
            
            self.api_path='product/v1/fail_delivery'

            self.validate()
            self.headers, self.body = generate_cloudevent(
                source='Logistic.PCL.UP.OMP',
                type='Logistic.PCL.Product.Delivery.FailDelivery',
                data=self.payload
        )

class DestroyProduct(BaseCommand):
     
    def __init__(self, product=None, timestamp=None, location=None, notes=None):
        self.payload = {}
        if product is not None:
            self.payload['product'] = product.to_json()
        if location is not None:
            self.payload['location'] = location.to_json()
        if timestamp is not None:
            self.payload['timestamp'] = timestamp
        if notes is not None:
            self.payload['notes'] = notes
        
        self.api_path='product/v1/destroy'

        self.validate()
        self.headers, self.body = generate_cloudevent(
            source='Logistic.PCL.UP.OMP',
            type='Logistic.PCL.Product.Processing.Destroy',
            data=self.payload
        )

class EndProductTransport(BaseCommand):
     
    def __init__(self, product=None, carrier=None, location=None, attributes=None, timestamp=None):
        self.payload = {}
        if product is not None:
            self.payload['product'] = product.to_json()
        if carrier is not None:
            self.payload['carrier'] = carrier.to_json()
        if location is not None:
            self.payload['location'] = location.to_json()
        if timestamp is not None:
            self.payload['timestamp'] = timestamp
        if attributes is not None:
            self.payload['attributes'] = attributes
        
        self.api_path='product/v1/end_transport'

        self.validate()
        self.headers, self.body = generate_cloudevent(
            source='Logistic.PCL.UP.OMP',
            type='Logistic.PCL.Product.Transport.EndTransport',
            data=self.payload
        )

class StartProductTransport(BaseCommand):
     
    def __init__(self, product=None, carrier=None, from_=None, to=None, attributes=None, timestamp=None):
        self.payload = {}
        if product is not None:
            self.payload['product'] = product.to_json()
        if carrier is not None:
            self.payload['carrier'] = carrier.to_json()
        if from_ is not None:
            self.payload['from'] = from_.to_json()
        if to is not None:
            self.payload['to'] = to.to_json()
        if timestamp is not None:
            self.payload['timestamp'] = timestamp
        if attributes is not None:
            self.payload['attributes'] = attributes
        
        self.api_path='product/v1/start_transport'

        self.validate()
        self.headers, self.body = generate_cloudevent(
            source='Logistic.PCL.UP.OMP',
            type='Logistic.PCL.Product.Transport.StartTransport',
            data=self.payload
        )

class WithdrawProduct(BaseCommand):

    def __init__(self, product=None, timestamp=None, notes=None):
        self.payload = {}
        if timestamp is not None:
            self.payload['timestamp'] = timestamp
        if product is not None:
            self.payload['product'] = product.to_json()
        if notes is not None:
            self.payload['notes'] = notes
        
        self.api_path='product/v1/withdraw'

        self.validate()
        self.headers, self.body = generate_cloudevent(
            source='Logistic.PCL.UP.OMP',
            type='Logistic.PCL.Product.Delivery.Withdraw',
            data=self.payload
        )