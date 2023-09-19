from commands.product import *
from queries.product import ProductStatus
from tests.functionals.plan import Plan
from commands.models import Receiver, Sender, Location, Product
import time

import utils

import pytest

@pytest.mark.functional
def test_product_status_one_command(config,product_id):

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
        ProductStatus(id=product_id)
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
    assert results[0]['product']['id'] == product_id
    assert results[0]['status'] == 'Accepted'


@pytest.mark.functional
def test_product_status_n_commands(config,product_id):
    # Arrange
    product = Product(product_id, 'Poste Delivery Business', 'BOX', [])

    cmds = [
        StartProductProcessing(
            product=product,
            location=Location('CS', 'Piacenza', 'Via Roma 16', '32456', 'Italia', '77G22', []),
            timestamp='2022-12-07T17:48:33.676829Z'
        ),
        DestroyProduct(
            product=product,
            timestamp='2022-12-09T08:34:00.676829Z',
            location=Location('CS', 'Milano', 'Viale Something 125', '13276', 'Italia', '33Y33', []),
            notes='a note'
        ),
        AcceptProduct(
            product=product,
            timestamp='2022-12-06T11:24:00.676829Z',
            receiver=Receiver('Mario Rossi', 'Roma', 'Roma', 'Corso Roma 23', '00144', '3333333333', 'mario@rossi.com', 'Carrier\'s note', []),
            sender=Sender('Luca Verdi', 'Milano', 'Milano', 'Corso Como 4', '12345', '2222222222', []),
            location=Location('UP', 'Roma', 'Viale Europa 190', '00144', 'Italia', '55Y90', []),
            attributes=[]
        )
    ]
    # Expectation is that commands are ordered by timestamp generated on field (which is the timestamp field in the command body)
    expStatus = 'Destroyed'

    queries = [
        ProductStatus(
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
    assert results[0]['status'] == expStatus

@pytest.mark.functional
def test_product_status_missing_api_key(config,product_id):

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
        ProductStatus(id=product_id)
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
def test_product_status_wrong_api_key(config,product_id):

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
        ProductStatus(id=product_id)
    ]
    
    plan = Plan(
        id=product_id,
        commands=cmds,
        queries=queries,
        base_api=config['BASE_API'],
        api_key='fakeAPiKey'
    )

    # Execute
    responses = plan.start_command()

    # Assert
    for response in responses:
        assert response.status_code == 403