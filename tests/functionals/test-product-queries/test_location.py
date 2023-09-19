from commands.product import *
from queries.product import ProductLocation
from tests.functionals.plan import Plan
from commands.models import Receiver, Sender, Location, Product, Carrier
import time
import utils 

import pytest

@pytest.mark.functional
@pytest.mark.skip(reason="bug LOI-502")
def test_product_location_one_command(config,product_id):
    # Arrange
    product = Product(product_id, 'Poste Delivery Business', 'BOX', [])
    receiver = Receiver('Mario Rossi', 'Roma', 'Roma', 'Corso Roma 23', '00144', '3333333333', 'mario@rossi.com', 'Carrier\'s note', [])
    sender = Sender('Luca Verdi', 'Milano', 'Milano', 'Corso Como 4', '12345', '2222222222', [])
    location = Location('UP', 'Roma', 'Viale Europa 190', '00144', 'Italia', '55Y90', [])

    cmds = [ 
        AcceptProduct(
            product=product,
            timestamp=utils.get_timestamp(),
            receiver=receiver,
            sender=sender,
            location=location,
            attributes=[]
        )
    ]
    
    queries = [
        ProductLocation(
            id=product_id
        )
    ]

    plan = Plan(
        id=product_id,
        commands=cmds,
        queries=queries,
        base_api=config['BASE_API'],
        api_key=config['API_KEY']
    )
    
    # Execute
    responses = plan.start_command()
    assert responses[0].status_code == 202
    time.sleep(5)
    results =  plan.start_query()

    # Assert
    assert len(results) == 1
    assert results[0]['product'] == cmds[0].payload['product']
    assert results[0]['location'] == cmds[0].payload['location']


@pytest.mark.functional
@pytest.mark.skip(reason="bug LOI-502")
def test_product_location_n_commands(config,product_id):
    # Arrange
    product = Product(product_id, 'Poste Delivery Business', 'BOX', [])
    receiver = Receiver('Mario Rossi', 'Roma', 'Roma', 'Corso Roma 23', '00144', '3333333333', 'mario@rossi.com', 'Carrier\'s note', [])
    sender = Sender('Luca Verdi', 'Milano', 'Milano', 'Corso Como 4', '12345', '2222222222', [])
    accepted_location = Location('UP', 'Roma', 'Viale Europa 190', '00144', 'Italia', '55Y90', [])
    processing_location = Location('UP', 'Bologna', 'Piazza Minghetti 4', '40124', 'Italia', '1928Z', [])

    cmds = [
        AcceptProduct(
             product=product,
             timestamp='2022-12-05T17:48:33.676829Z',
             receiver=receiver,
             sender=sender,
             location=accepted_location,
             attributes=[]
        ),
        StartProductProcessing(
             product=product,
             timestamp='2022-12-07T17:48:33.676829Z',
             location=processing_location
        ),
        StartProductDelivery(
             product=product,
             timestamp='2022-12-06T17:48:33.676829Z',
             attributes=[],
             carrier=Carrier('motorbike', 'Postman12345', 'FE764OU')
        )
    ]
    
    queries = [
        ProductLocation(
            id=product_id
        )
    ]
  
    plan = Plan(
        id=product_id,
        commands=cmds,
        queries=queries,
        base_api=config['BASE_API'],
        api_key=config['API_KEY']
    )

    
    # Execute
    responses = plan.start_command()
    for response in responses:
        assert response.status_code == 202
    time.sleep(5)
    results =  plan.start_query()

    # Assert
    assert len(results) == 1
    assert results[0]['product'] == cmds[0].payload['product']
    assert results[0]['location'] == cmds[1].payload['location']


@pytest.mark.functional
@pytest.mark.skip(reason="bug LOI-502")
def test_product_location_n_equal_commands(config,product_id):
   # Arrange
    product = Product(product_id, 'Poste Delivery Business', 'BOX', [])
    receiver = Receiver('Mario Rossi', 'Roma', 'Roma', 'Corso Roma 23', '00144', '3333333333', 'mario@rossi.com', 'Carrier\'s note', [])
    sender = Sender('Luca Verdi', 'Milano', 'Milano', 'Corso Como 4', '12345', '2222222222', [])

    cmds = [
        AcceptProduct(
             product=product,
             timestamp=utils.get_timestamp(),
             receiver=receiver,
             sender=sender,
             location=Location("UP", "Bologna", "Piazza Minghetti 4", "40124", "Italia", "1928Z", []),
             attributes=[]
        ),
        AcceptProduct(
             product=product,
             timestamp=utils.get_timestamp(),
             receiver=receiver,
             sender=sender,
             location=Location("UP", "Roma", "Viale Europa 190", "00144", "Italia", "55Y90", []),
             attributes=[]
        )
    ]
    
    queries = [
        ProductLocation(
            id=product_id
        )
    ]
  
    plan = Plan(
        id=product_id,
        commands=cmds,
        queries=queries,
        base_api=config['BASE_API'],
        api_key=config['API_KEY']
    )

    
    # Execute
    responses = plan.start_command()
    for response in responses:
        assert response.status_code == 202
    time.sleep(5)
    results =  plan.start_query()

    # Assert
    assert len(results) == 1
    assert results[0]['product'] == cmds[0].payload['product']
    assert results[0]['location'] == cmds[1].payload['location']

@pytest.mark.functional
@pytest.mark.skip(reason="bug LOI-504")
def test_product_command_no_location(config,product_id):
    # Arrange
    product = Product(product_id, 'Poste Delivery Business', 'BOX', [])
    carrier = Carrier('motorbike', 'Postman12345', 'FE764OU')

    cmds = [ 
        StartProductDelivery(
            product=product,
            carrier=carrier,
            timestamp=utils.get_timestamp(),
            attributes=[]
        )
    ]
    
    queries = [
        ProductLocation(
            id=product_id
        )
    ]

    plan = Plan(
        id=product_id,
        commands=cmds,
        queries=queries,
        base_api=config['BASE_API'],
        api_key=config['API_KEY']
    )
    
    # Execute
    responses = plan.start_command()
    assert responses[0].status_code == 202
    time.sleep(5)
    results =  plan.start_query()

    # Assert
    assert len(results) == 1
    assert results[0]['product'] == cmds[0].payload['product']
    assert results[0]['location'] == {}

@pytest.mark.functional
def test_product_location_missing_api_key(config,product_id):
    # Arrange
    product = Product(product_id, 'Poste Delivery Business', 'BOX', [])
    receiver = Receiver('Mario Rossi', 'Roma', 'Roma', 'Corso Roma 23', '00144', '3333333333', 'mario@rossi.com', 'Carrier\'s note', [])
    sender = Sender('Luca Verdi', 'Milano', 'Milano', 'Corso Como 4', '12345', '2222222222', [])
    location = Location('UP', 'Roma', 'Viale Europa 190', '00144', 'Italia', '55Y90', [])

    cmds = [ 
        AcceptProduct(
            product=product,
            timestamp=utils.get_timestamp(),
            receiver=receiver,
            sender=sender,
            location=location,
            attributes=[]
        )
    ]
    
    queries = [
        ProductLocation(
            id=product_id
        )
    ]

    plan = Plan(
        id=product_id,
        commands=cmds,
        queries=queries,
        base_api=config['BASE_API'],
        # missing api key
    )
    
    # Execute
    responses = plan.start_command()

    # Assert
    for response in responses:
        assert response.status_code == 403

@pytest.mark.functional
def test_product_location_wrong_api_key(config,product_id):
    # Arrange
    product = Product(product_id, 'Poste Delivery Business', 'BOX', [])
    receiver = Receiver('Mario Rossi', 'Roma', 'Roma', 'Corso Roma 23', '00144', '3333333333', 'mario@rossi.com', 'Carrier\'s note', [])
    sender = Sender('Luca Verdi', 'Milano', 'Milano', 'Corso Como 4', '12345', '2222222222', [])
    location = Location('UP', 'Roma', 'Viale Europa 190', '00144', 'Italia', '55Y90', [])

    cmds = [ 
        AcceptProduct(
            product=product,
            timestamp=utils.get_timestamp(),
            receiver=receiver,
            sender=sender,
            location=location,
            attributes=[]
        )
    ]
    
    queries = [
        ProductLocation(
            id=product_id
        )
    ]

    plan = Plan(
        id=product_id,
        commands=cmds,
        queries=queries,
        base_api=config['BASE_API'],
        api_key='fakeApiKey'
    )
    
    # Execute
    responses = plan.start_command()

    # Assert
    for response in responses:
        assert response.status_code == 403
