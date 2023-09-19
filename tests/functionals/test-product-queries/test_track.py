from commands.product import *
from queries.product import ProductTrack
from tests.functionals.plan import Plan
from commands.models import Receiver, Sender, Location, Carrier, Product
from datetime import datetime, timezone
import time
import utils

import pytest

@pytest.mark.functional
def test_product_track(config,product_id):
    # Arrange
    product = Product(product_id, 'Poste Delivery Business', 'BOX', [])
    numDuplicatedCmds = 2
    cmds = [
        AcceptProduct(
            product=product,
            timestamp=utils.get_timestamp(),
            receiver=Receiver('Mario Rossi', 'Roma', 'Roma', 'Corso Roma 23', '00144', '3333333333', 'mario@rossi.com', 'Carrier\'s note', []),
            sender=Sender('Luca Verdi', 'Milano', 'Milano', 'Corso Como 4', '12345', '2222222222', []),
            location=Location('UP', 'Roma', 'Viale Europa 190', '00144', 'Italia', '55Y90', []),
            attributes=[]
        ),
        AcceptProduct( # numDuplicatedCmds++
            product=product,
            timestamp=utils.get_timestamp(),
            receiver=Receiver('Mario Rossi', 'Roma', 'Roma', 'Corso Roma 23', '00144', '3333333333', 'mario@rossi.com', 'Carrier\'s note', []),
            sender=Sender('Luca Verdi', 'Sesto San Giovanni', 'Milano', 'Via Novara 23', '54321', '2222222222', []),
            location=Location('UP', 'Roma', 'Viale Europa 190', '00144', 'Italia', '55Y90', []),
            attributes=[]
        ),
        StartProductProcessing(
            product=product,
            timestamp='2000-12-07T17:48:33.676829Z',
            location=Location('CS', 'Roma', 'Via Fiumicino', '00144', 'Italia', '23Y90', [])
        ),
        MakeProductReadyToBeDelivered(
            product=product,
            timestamp=utils.get_timestamp(),
            attributes=[]
        ),
        MakeProductReadyToBeDelivered( # numDuplicatedCmds++
            product=product,
            timestamp=utils.get_timestamp(),
            attributes=[]
        ),
        StartProductDelivery(
            product=product,
            timestamp=utils.get_timestamp(),
            carrier = Carrier('motorbike', 'Postman12345', 'FE764OU'),
            attributes=[]
        ),
        CompleteProductDelivery(
            product=product,
            timestamp=utils.get_timestamp(),
            carrier = Carrier('motorbike', 'Postman12345', 'FE764OU'),
            attributes=[]
        )
    ]

    queries = [
        ProductTrack(
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
    assert results[0]['product']['id'] == product_id
    assert len(results[0]['trackList']) == 7     
    assert results[0]['trackList'][0]['description'] == 'WIP'
    assert results[0]['trackList'][1]['description'] == 'Accepted'
    assert results[0]['trackList'][2]['description'] == 'Accepted'
    assert results[0]['trackList'][3]['description'] == 'ReadyForDelivery'
    assert results[0]['trackList'][4]['description'] == 'ReadyForDelivery'
    assert results[0]['trackList'][5]['description'] == 'LastMileTransit'
    assert results[0]['trackList'][6]['description'] == 'Delivered'

@pytest.mark.functional
def test_product_track_missing_api_key(config,product_id):
    # Arrange
    product = Product(product_id, 'Poste Delivery Business', 'BOX', [])
    numDuplicatedCmds = 2
    cmds = [
        AcceptProduct(
            product=product,
            timestamp=utils.get_timestamp(),
            receiver=Receiver('Mario Rossi', 'Roma', 'Roma', 'Corso Roma 23', '00144', '3333333333', 'mario@rossi.com', 'Carrier\'s note', []),
            sender=Sender('Luca Verdi', 'Milano', 'Milano', 'Corso Como 4', '12345', '2222222222', []),
            location=Location('UP', 'Roma', 'Viale Europa 190', '00144', 'Italia', '55Y90', []),
            attributes=[]
        )
    ]

    queries = [
        ProductTrack(
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