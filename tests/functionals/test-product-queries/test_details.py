from commands.product import *
from queries.product import ProductDetails
from tests.functionals.plan import Plan
from commands.models import Receiver, Sender, Location, Product
import time
import utils

import pytest

@pytest.mark.functional
@pytest.mark.skip(reason="bug LOI-533")
def test_product_details(config,product_id):
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
        ProductDetails(
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
    results = plan.start_query()
    
    # Assert
    assert len(results) == 1
    assert results[0]['product'] == cmds[0].payload['product']
    assert results[0]['receiver'] == cmds[0].payload['receiver']
    assert results[0]['sender'] == cmds[0].payload['sender']

@pytest.mark.functional
# @pytest.mark.skip(reason="bug LOI-488")
def test_product_details_product_not_found(config,product_id):
    # Arrange  
    queries = [
        ProductDetails(
            id=product_id
        )
    ]
    plan = Plan(
        id=product_id,
        commands=[],
        queries=queries,
        base_api=config['BASE_API'],
        api_key=config['API_KEY']
    )

    # Execute
    results = plan.start_query()
    
    # Assert
    assert len(results) == 1
    assert results[0] == {}

def test_product_details_missing_api_key(config,product_id):
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
        ProductDetails(
            id=product_id
        )
    ]
    plan = Plan(
        id=product_id,
        commands=cmds,
        queries=queries,
        base_api=config['BASE_API']
        # missing api key
    )

    # Execute
    responses = plan.start_command()

    # Assert
    for response in responses:
        assert response.status_code == 403

def test_product_details_wrong_api_key(config,product_id):
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
        ProductDetails(
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
