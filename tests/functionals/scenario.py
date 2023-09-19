
from commands.product import *
#from commands.product import AcceptProduct, StartProductProcessing, FailProductProcessing, MakeProductReadyToBeDelivered, StartProductDelivery, CompleteProductDelivery, DestroyProduct
from queries.product import ProductLocation, ProductStatus, ProductDetails, ProductTrace
from commands.models import Receiver, Sender, Location, Carrier
from datetime import datetime, timezone


# class ScenarioAcceptProduct:
#     def __init__(self, id: str):
#         self.cmds = []
#         self.queries = []
#         receiver = Receiver()
#         sender = Sender()
#         location = Location()
        
#         ap = AcceptProduct(
#             id=id,
#             timestamp=datetime.timestamp(datetime.now(timezone.utc)),
#             receiver=receiver,
#             sender=sender,
#             location=location
#         )

#         status = ProductStatus(
#             id = id
#         )

#         loc = ProductLocation(
#             id=id
#         )

#         details = ProductDetails(
#             id = id
#         )

        

#         self.cmds.append(ap)

#         self.queries.append(status)
#         self.queries.append(loc)
#         self.queries.append(details)
        

#     @property
#     def plan(self):
#         return self.cmds
    
#     @property
#     def plan_query(self):
#         return self.queries

# class ScenarioProductDelivered:
#     def __init__(self, id: str):
#         self.cmds = []
#         self.queries = []
#         receiver = Receiver()
#         sender = Sender()
#         location = Location()
#         carrier = Carrier()
#         ap = AcceptProduct(
#             id=id,
#             timestamp=datetime.timestamp(datetime.now(timezone.utc)),
#             receiver=receiver,
#             sender=sender,
#             location=location
#         )

#         spp = StartProductProcessing(
#             id=id,
#             timestamp=datetime.timestamp(datetime.now(timezone.utc)),
#             location=location
#         )

#         mprtbd = MakeProductReadyToBeDelivered(
#             id = id,
#             timestamp=datetime.timestamp(datetime.now(timezone.utc))
#         )

#         spd = StartProductDelivery(
#                 id = id,
#                 timestamp=datetime.timestamp(datetime.now(timezone.utc)),
#                 carrier = carrier
#         )

#         cpd = CompleteProductDelivery(
#                 id = id,
#                 timestamp=datetime.timestamp(datetime.now(timezone.utc)),
#                 carrier = carrier
#         )


#         status = ProductStatus(
#             id=id
#         )

#         trace = ProductTrace(
#             id = id
#         )



#         self.cmds.append(ap)
#         self.cmds.append(spp)
#         self.cmds.append(mprtbd)
#         self.cmds.append(spd)
#         self.cmds.append(cpd)

#         self.queries.append(status)
#         self.queries.append(trace)

  

#     @property
#     def plan(self):
#         return self.cmds
    
#     @property
#     def plan_query(self):
#         return self.queries


# class ScenarioProductDestroyed:

#     def __init__(self, id: str):
#         self.cmds = []
#         self.queries = []
#         receiver = Receiver()
#         sender = Sender()
#         location = Location()
        
#         ap = AcceptProduct(
#             id=id,
#             timestamp=datetime.timestamp(datetime.now(timezone.utc)),
#             receiver=receiver,
#             sender=sender,
#             location=location
#         )

#         spp = StartProductProcessing(
#             id=id,
#             timestamp=datetime.timestamp(datetime.now(timezone.utc)),
#             location=location
#         )

#         fpp = FailProductProcessing(
#             id=id,
#             timestamp=datetime.timestamp(datetime.now(timezone.utc)),
#             location=location
#         )

#         destroy = DestroyProduct(
#             id=id,
#             timestamp=datetime.timestamp(datetime.now(timezone.utc)),
#             location=location
#         )


#         status = ProductStatus(
#             id = id
#         )

#         loc = ProductLocation(
#             id=id
#         )

#         details = ProductDetails(
#             id = id
#         )

        

#         self.cmds.append(ap)
#         self.cmds.append(spp)
#         self.cmds.append(fpp)
#         self.cmds.append(destroy)

#         self.queries.append(status)
#         self.queries.append(loc)
#         self.queries.append(details)
        

#     @property
#     def plan(self):
#         return self.cmds
    
#     @property
#     def plan_query(self):
#         return self.queries
    


class ScenarioProductStatus:
    def __init__(self, id: str):
        self.cmds = []
        receiver = Receiver()
        sender = Sender()
        location = Location()
        
        ap = AcceptProduct(
            id=id,
            timestamp=datetime.timestamp(datetime.now(timezone.utc)),
            receiver=receiver,
            sender=sender,
            location=location
        )
    
        status = ProductStatus(
            id = id
        )

        self.cmds.append(ap)
        self.queries= status
        
    @property
    def plan(self):
        return self.cmds
    
    @property
    def plan_query(self):
        return self.queries
