import json
from faker import Faker
import random
faker = Faker("it_IT")


class BaseModel:
    def to_json(self):
        return json.loads(json.dumps(self.__dict__))

class Empty(BaseModel):
    pass


class Product(BaseModel):
    def __init__(self, id, name, type, attributes):
        self.id = id
        self.name = name
        self.type = type
        self.attributes = attributes


class Receiver(BaseModel):
    def __init__(self, name, city, province, address, zipcode, number, email, note, attributes):
        self.name = name
        self.city = city
        self.province = province
        self.address = address
        self.zipcode = zipcode
        self.number = number
        self.email = email
        self.note = note
        self.attributes = attributes


class Sender(BaseModel):
    def __init__(self, name, city, province, address, zipcode, number, attributes):
        self.name = name
        self.city = city
        self.province = province
        self.address = address
        self.zipcode = zipcode
        self.number = number
        self.attributes = attributes


class Location(BaseModel):
    def __init__(self, type, city, address, zipcode, nation, locationCode, attributes):
        self.type = type
        self.city = city
        self.address = address
        self.zipcode = zipcode
        self.nation = nation
        self.locationCode = locationCode
        self.attributes = attributes

class Carrier(BaseModel):
    def __init__(self, typeId, driverId, vehicleId):
        self.typeId = typeId
        self.driverId = driverId
        self.vehicleId = vehicleId

class FromTo(BaseModel):
    def __init__(self, type, address, zipcode, city, nation):
        self.type = type
        self.address = address
        self.zipcode = zipcode
        self.city = city
        self.nation = nation


class Transport(BaseModel):
    def __init__(self, id):
        self.id = id       


class Cluster(BaseModel):
    def __init__(self, id):
        self.id = id        




